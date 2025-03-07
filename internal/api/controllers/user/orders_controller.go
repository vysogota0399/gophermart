package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/vysogota0399/gophermart_portal/internal/api"
	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/clients"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
)

type OrdersController struct {
	lg      *logging.ZapLogger
	billing BillingOrdersRpcService
	query   QueryOrdersRpcService
}

type BillingOrdersRpcService interface {
	CreateOrder(ctx context.Context, in *clients.NewOrder) error
}

type QueryOrdersRpcService interface {
	OrdersCollection(ctx context.Context, accountID int64) ([]*models.Order, error)
}

func NewOrdersController(billing BillingOrdersRpcService, query QueryOrdersRpcService, lg *logging.ZapLogger) *OrdersController {
	return &OrdersController{billing: billing, query: query, lg: lg}
}

func (c *OrdersController) CreateRoutes(r *api.Router) []*api.Route {
	return []*api.Route{
		{
			Handler: c.createOrderHandler,
			Path:    "/api/user/orders",
			Method:  http.MethodPost,
			Meta:    map[string]bool{api.AuthorizationRequire: true},
		},
		{
			Handler: c.showOrdersHandler,
			Path:    "/api/user/orders",
			Method:  http.MethodGet,
			Meta:    map[string]bool{api.AuthorizationRequire: true},
		},
	}
}

type NewOrderInput struct {
	Number string `json:"number" binding:"required,luhnablenumber"`
}

func (cntr *OrdersController) createOrderHandler(c *gin.Context) {
	current_user := c.GetString(api.CurrentUserIDKey)
	ctx := cntr.lg.WithContextFields(
		c,
		zap.String("actor", "orders_controller"),
		zap.String(api.CurrentUserIDKey, current_user),
	)

	var payload NewOrderInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		cntr.lg.ErrorCtx(ctx, "invalid request params", zap.Error(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	accountID, err := strconv.ParseInt(current_user, 10, 64)
	if err != nil {
		cntr.lg.ErrorCtx(ctx, "convert current_user id to int64 failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "contact the operator"})
	}

	order := &clients.NewOrder{
		UUID:       uuid.NewV4().String(),
		Number:     payload.Number,
		UploadedAt: time.Now(),
		AccountID:  accountID,
	}

	if err := cntr.billing.CreateOrder(ctx, order); err != nil {
		if status.Code(err) == codes.AlreadyExists {
			cntr.lg.ErrorCtx(ctx, "create order failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("order with number %s already exists", order.Number)})
			return
		}

		cntr.lg.ErrorCtx(ctx, "create order failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "contact the operator"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (cntr *OrdersController) showOrdersHandler(c *gin.Context) {
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

	orders, err := cntr.query.OrdersCollection(c, accountID)
	if err != nil {
		cntr.lg.ErrorCtx(ctx, "search orders errors", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "contact the operator"})
		return
	}

	if len(orders) == 0 {
		cntr.lg.DebugCtx(ctx, "orders not found", zap.Error(err))
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, orders)
}
