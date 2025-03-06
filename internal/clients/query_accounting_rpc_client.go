package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/config"
	"github.com/vysogota0399/gophermart_protos/gen/entities"
	query "github.com/vysogota0399/gophermart_protos/gen/queries/accounting"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QueryAccountingRpcClient struct {
	conn  *grpc.ClientConn
	query query.QueryAccountingClient
	lg    *logging.ZapLogger
}

func NewQueryAccountingRpcClient(lc fx.Lifecycle, cfg *config.Config, lg *logging.ZapLogger) (*QueryAccountingRpcClient, error) {
	conn, err := grpc.NewClient(cfg.QueryAccountingAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	lc.Append(
		fx.StopHook(func() error {
			return conn.Close()
		}),
	)

	return &QueryAccountingRpcClient{
		conn:  conn,
		query: query.NewQueryAccountingClient(conn),
		lg:    lg,
	}, nil
}

func (rpc *QueryAccountingRpcClient) GetBalance(ctx context.Context, accountID int64) (*models.Balance, error) {
	response, err := rpc.query.GetBalance(
		ctx,
		&query.GetBalanceParams{
			Account: &entities.Account{Id: accountID},
		})
	if err != nil {
		return nil, fmt.Errorf("query_accounting_rpc_client billing query balance error %w", err)
	}

	return &models.Balance{
		Current:   float64(response.Balance.Units) / 100,
		Withdrawn: float64(response.Withdrawn.Units) / 100,
	}, nil
}

func (rpc *QueryAccountingRpcClient) GetWithdrawals(ctx context.Context, accountID int64) ([]*models.Withdraw, error) {
	response, err := rpc.query.GetWithdrawals(
		ctx,
		&query.GetWithdrawalsParams{
			Account: &entities.Account{Id: accountID},
		})
	if err != nil {
		return nil, fmt.Errorf("query_accounting_rpc_client billing query balance error %w", err)
	}

	windrawals := []*models.Withdraw{}
	for _, w := range response.Withdrawals {
		windrawals = append(
			windrawals,
			&models.Withdraw{
				OrderNumber: w.OrderNumber,
				Sum:         float64(w.Sum.Units) / 100,
				ProcessedAt: w.ProcessedAt.AsTime().Format(time.RFC3339Nano),
			},
		)
	}

	return windrawals, nil
}
