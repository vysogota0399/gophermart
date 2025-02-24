package balance

import (
	"context"
	"net/http"
	"strconv"

	"github.com/vysogota0399/gophermart_portal/internal/api"
	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
)

type WithdrawalsController struct {
	lg      *logging.ZapLogger
	billing WithdrawRpcService
}

type WithdrawRpcService interface {
	Withdraw(ctx context.Context, accountID int64, amount float64, orderNumber string) error
}

func NewWithdrawalsController(billing WithdrawRpcService, lg *logging.ZapLogger) *WithdrawalsController {
	return &WithdrawalsController{billing: billing, lg: lg}
}

func (c *WithdrawalsController) CreateRoutes(r *api.Router) []*api.Route {
	return []*api.Route{
		{
			Handler: c.createWithdrawHandler,
			Path:    "/api/user/balance/withdraw",
			Method:  http.MethodPost,
			Meta:    map[string]bool{api.AuthorizationRequire: true},
		},
	}
}

type createWithdrawInput struct {
	Number string  `json:"order" binding:"required,luhnablenumber"`
	Amount float64 `json:"sum" binding:"required,gte=0,lte=999999999999"`
}

func (cntr *WithdrawalsController) createWithdrawHandler(c *gin.Context) {
	var payload createWithdrawInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		cntr.lg.ErrorCtx(c, "invalid request params", zap.Error(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

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

	if err := cntr.billing.Withdraw(ctx, accountID, payload.Amount, payload.Number); err != nil {
		if status.Code(err) == codes.Aborted {
			s := status.Convert(err)
			for _, d := range s.Details() {
				switch info := d.(type) {
				case *errdetails.BadRequest:
					cntr.lg.ErrorCtx(ctx, "bad request", zap.Any("details", info))
				default:
					cntr.lg.ErrorCtx(ctx, "unexpected error type", zap.String("details", s.Message()))
				}
			}
			c.Status(http.StatusPaymentRequired)
			return
		}

		cntr.lg.ErrorCtx(ctx, "withdraw failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "contact the operator"})
		return
	}

	c.Status(http.StatusOK)
}
