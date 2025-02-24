package user

import (
	"context"
	"net/http"

	"github.com/vysogota0399/gophermart_portal/internal/api"
	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthenticationService interface {
	Call(ctx context.Context, w http.ResponseWriter, u *models.User) error
}

type LoginController struct {
	lg          *logging.ZapLogger
	authService AuthenticationService
	Path        string
}

func NewLoginController(lg *logging.ZapLogger, authService AuthenticationService) *LoginController {
	return &LoginController{lg: lg, authService: authService, Path: "/api/user/login"}
}

func (cntr *LoginController) CreateRoutes(r *api.Router) []*api.Route {
	return []*api.Route{
		{
			Handler: cntr.loginHandler,
			Path:    cntr.Path,
			Method:  http.MethodPost,
		},
	}
}

type LoginInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (cntr *LoginController) loginHandler(c *gin.Context) {
	ctx := cntr.lg.WithContextFields(c, zap.String("actor", "loginController"))
	var payload LoginInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request params"})
		return
	}

	if err := cntr.authService.Call(ctx, c.Writer, &models.User{Login: payload.Login, Password: payload.Password}); err != nil {
		cntr.lg.ErrorCtx(ctx, "auth failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
