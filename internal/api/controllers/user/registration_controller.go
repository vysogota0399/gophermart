package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/vysogota0399/gophermart_portal/internal/api"
	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/api/repositories"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RegistrationService interface {
	Call(ctx context.Context, u *models.User) error
}

type LoginService interface {
	Call(ctx context.Context, w http.ResponseWriter, u *models.User) error
}

type RegistrtionController struct {
	registrationService RegistrationService
	loginService        LoginService
	lg                  *logging.ZapLogger
	Path                string
}

func NewRegistrtionController(registrationService RegistrationService, lg *logging.ZapLogger, loginService LoginService) *RegistrtionController {
	return &RegistrtionController{registrationService: registrationService, lg: lg, loginService: loginService, Path: "/api/user/register"}
}

func (cntr *RegistrtionController) CreateRoutes(router *api.Router) []*api.Route {
	return []*api.Route{
		{
			Handler: cntr.reginstrationHandler,
			Path:    cntr.Path,
			Method:  http.MethodPost,
		},
	}
}

type RegistrationInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (cntr *RegistrtionController) reginstrationHandler(c *gin.Context) {
	ctx := cntr.lg.WithContextFields(c, zap.String("actor", "registration_controller"))

	var payload RegistrationInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request params"})
		return
	}

	user := &models.User{Login: payload.Login, Password: payload.Password}

	if err := cntr.registrationService.Call(ctx, user); err != nil {
		if errors.Is(err, repositories.ErrLoginAlreadyExist) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "cantact with operator"})
		return
	}

	if err := cntr.loginService.Call(ctx, c.Writer, user); err != nil {
		cntr.lg.ErrorCtx(ctx, "login failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
