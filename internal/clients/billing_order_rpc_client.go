package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/config"
	billing "github.com/vysogota0399/gophermart_protos/gen/commands/create_order"
	"github.com/vysogota0399/gophermart_protos/gen/common"
	"github.com/vysogota0399/gophermart_protos/gen/entities"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type BillingOrdersRpcClient struct {
	conn    *grpc.ClientConn
	billing billing.CreateOrderServiceClient
	lg      *logging.ZapLogger
}

func NewBillingOrdersRpcClient(lc fx.Lifecycle, cfg *config.Config, lg *logging.ZapLogger) (*BillingOrdersRpcClient, error) {
	conn, err := grpc.NewClient(cfg.BillingOrdersAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	lc.Append(
		fx.StopHook(func() error {
			return conn.Close()
		}),
	)

	return &BillingOrdersRpcClient{
		conn:    conn,
		billing: billing.NewCreateOrderServiceClient(conn),
		lg:      lg,
	}, nil
}

type NewOrder struct {
	UUID       string
	Number     string
	UploadedAt time.Time
	AccountID  int64
}

func (rpc *BillingOrdersRpcClient) CreateOrder(ctx context.Context, in *NewOrder) error {
	ctx = rpc.lg.WithContextFields(ctx, zap.String("actor", "billing_rpc_client"))
	rpc.lg.DebugCtx(ctx, "remote billing call CreateOrder", zap.Any("order", in))

	message := &billing.CreateNewOrderParams{
		Uuid:       &common.Uuid{Value: in.UUID},
		Number:     in.Number,
		UploadedAt: timestamppb.New(in.UploadedAt),
		Account:    &entities.Account{Id: in.AccountID},
	}

	rpc.lg.DebugCtx(ctx, "create order", zap.Any("order", message))
	_, err := rpc.billing.Create(ctx, message)

	if err != nil {
		return fmt.Errorf("gophermart/internal/clients/billing_rpc_client billing CreateOrder error %w", err)
	}

	return nil
}
