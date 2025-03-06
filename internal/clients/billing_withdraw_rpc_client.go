package clients

import (
	"context"
	"fmt"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/config"
	billing "github.com/vysogota0399/gophermart_protos/gen/commands/withdraw"
	"github.com/vysogota0399/gophermart_protos/gen/entities"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BillingWithdrawRpcClient struct {
	conn    *grpc.ClientConn
	billing billing.WithdrawServiceClient
	lg      *logging.ZapLogger
}

func NewBillingWithdrawRpcClient(lc fx.Lifecycle, cfg *config.Config, lg *logging.ZapLogger) (*BillingWithdrawRpcClient, error) {
	conn, err := grpc.NewClient(cfg.BillingAccountingAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	lc.Append(
		fx.StopHook(func() error {
			return conn.Close()
		}),
	)

	return &BillingWithdrawRpcClient{
		conn:    conn,
		billing: billing.NewWithdrawServiceClient(conn),
		lg:      lg,
	}, nil
}

func (rpc *BillingWithdrawRpcClient) DoWithdraw(ctx context.Context, accountID int64, amount float64, orderNumber string) error {
	params := &billing.DoWithdrawParams{
		Account:     &entities.Account{Id: accountID},
		Amount:      &money.Money{Units: int64(amount * 100)},
		OrderNumber: orderNumber,
	}

	rpc.lg.DebugCtx(
		ctx,
		"do withdraw",
		zap.Any("params", params),
	)
	_, err := rpc.billing.DoWithdraw(ctx, params)
	if err != nil {
		return fmt.Errorf("billing_accounting_rpc_client billing billing balance error %w", err)
	}

	return nil
}
