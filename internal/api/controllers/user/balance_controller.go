package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/vysogota0399/gophermart_portal/internal/api"
	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type BalanceController struct {
	lg    *logging.ZapLogger
	query QueryBalanceRpcService
}

type QueryBalanceRpcService interface {
	GetBalance(ctx context.Context, accountID int64) (*models.Balance, error)
}

func NewBalanceController(query QueryBalanceRpcService, lg *logging.ZapLogger) *BalanceController {
	return &BalanceController{query: query, lg: lg}
}

func (c *BalanceController) CreateRoutes(r *api.Router) []*api.Route {
	return []*api.Route{
		{
			Handler: c.showBalanceHandler,
			Path:    "/api/user/balance",
			Method:  http.MethodGet,
			Meta:    map[string]bool{api.AuthorizationRequire: true},
		},
	}
}

func (cntr *BalanceController) showBalanceHandler(c *gin.Context) {
	current_user := c.GetString(api.CurrentUserIDKey)

	ctx := cntr.lg.WithContextFields(
		c,
		zap.String("actor", "orders_controller"),
		zap.String(api.CurrentUserIDKey, current_user),
	)

	accountID, err := strconv.ParseInt(current_user, 10, 64)
	if err != nil {
		cntr.lg.ErrorCtx(ctx, "convert current_user id to int64 failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "contact the operator"})
		return
	}

	balance, err := cntr.query.GetBalance(c, accountID)
	if err != nil {
		cntr.lg.ErrorCtx(ctx, "get balance errors", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "contact the operator"})
		return
	}

	c.JSON(http.StatusOK, balance)
}
