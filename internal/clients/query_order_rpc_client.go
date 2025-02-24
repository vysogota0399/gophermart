package clients

import (
	"context"
	"fmt"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/config"
	query "github.com/vysogota0399/gophermart_protos/gen/queries/orders"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QueryOrdersRpcClient struct {
	conn  *grpc.ClientConn
	query query.QueryOrdersClient
	lg    *logging.ZapLogger
}

func NewQueryOrdersRpcClient(lc fx.Lifecycle, cfg *config.Config, lg *logging.ZapLogger) (*QueryOrdersRpcClient, error) {
	conn, err := grpc.NewClient(cfg.QueryOrdersAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	lc.Append(
		fx.StopHook(func() error {
			return conn.Close()
		}),
	)

	return &QueryOrdersRpcClient{
		conn:  conn,
		query: query.NewQueryOrdersClient(conn),
		lg:    lg,
	}, nil
}

func (rpc *QueryOrdersRpcClient) OrdersCollection(ctx context.Context, accountID int64) ([]*models.Order, error) {
	ctx = rpc.lg.WithContextFields(ctx, zap.String("actor", "query_orders_rpc_client"))
	rpc.lg.DebugCtx(ctx, "remote query call SearchOrders", zap.Any("search_parmas", accountID))

	response, err := rpc.query.OrdersCollection(
		ctx,
		&query.QueryOrdersParams{
			AccountId: accountID,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("query_order_rpc_client: search orders error %w", err)
	}

	orders := []*models.Order{}

	for _, o := range response.Orders {
		orders = append(
			orders,
			&models.Order{
				Number:     o.Number,
				State:      o.State,
				Accrual:    o.Accrual,
				UploadedAt: o.UploadedAt,
			},
		)
	}

	return orders, nil
}
