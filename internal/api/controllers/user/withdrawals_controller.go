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

type WithdrawalsController struct {
	lg    *logging.ZapLogger
	query QueryWithdrawalsRpcService
}

type QueryWithdrawalsRpcService interface {
	GetWithdrawals(ctx context.Context, accountID int64) ([]*models.Withdraw, error)
}

func NewWithdrawalsController(query QueryWithdrawalsRpcService, lg *logging.ZapLogger) *WithdrawalsController {
	return &WithdrawalsController{query: query, lg: lg}
}

func (c *WithdrawalsController) CreateRoutes(router *api.Router) []*api.Route {
	return []*api.Route{
		{
			Handler: c.showWithdrawHandler,
			Path:    "/api/user/withdrawals",
			Method:  http.MethodGet,
			Meta:    map[string]bool{api.AuthorizationRequire: true},
		},
	}
}

func (cntr *WithdrawalsController) showWithdrawHandler(c *gin.Context) {
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

	withdrawals, err := cntr.query.GetWithdrawals(ctx, accountID)
	if err != nil {
		cntr.lg.ErrorCtx(ctx, "get withdrawals errors", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "contact the operator"})
		return
	}

	if len(withdrawals) == 0 {
		cntr.lg.DebugCtx(ctx, "withdrawals not found", zap.Error(err))
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}
