package main

import (
	"github.com/vysogota0399/gophermart_portal/internal/api"
	"github.com/vysogota0399/gophermart_portal/internal/api/controllers/user"
	"github.com/vysogota0399/gophermart_portal/internal/api/controllers/user/balance"
	"github.com/vysogota0399/gophermart_portal/internal/api/iam"
	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/repositories"
	"github.com/vysogota0399/gophermart_portal/internal/api/services"
	"github.com/vysogota0399/gophermart_portal/internal/clients"
	"github.com/vysogota0399/gophermart_portal/internal/config"
	"github.com/vysogota0399/gophermart_portal/internal/storage"

	"go.uber.org/fx"
)

func main() {
	fx.New(CreateApp()).Run()
}

func CreateApp() fx.Option {
	return fx.Options(
		fx.Provide(
			logging.NewZapLogger,

			fx.Annotate(repositories.NewUsersRepository, fx.As(new(services.RegistrationServiceUsersRepository))),
			repositories.NewUsersRepository,
			repositories.NewSessionsRepository,
			storage.NewStorage,
			api.NewHTTPServer,
			iam.NewIam,

			fx.Annotate(clients.NewBillingOrdersRpcClient, fx.As(new(user.BillingOrdersRpcService))),
			fx.Annotate(clients.NewQueryOrdersRpcClient, fx.As(new(user.QueryOrdersRpcService))),
			fx.Annotate(clients.NewQueryAccountingRpcClient, fx.As(new(user.QueryBalanceRpcService))),
			fx.Annotate(clients.NewQueryAccountingRpcClient, fx.As(new(user.QueryWithdrawalsRpcService))),
			fx.Annotate(clients.NewBillingWithdrawRpcClient, fx.As(new(balance.WithdrawRpcService))),

			fx.Annotate(iam.NewIam, fx.As(new(services.SessionCreator)), fx.ResultTags(`name:"sessionCreator"`)),
			fx.Annotate(iam.NewIam, fx.As(new(services.Authenticator)), fx.ResultTags(`name:"authenticator"`)),
			fx.Annotate(iam.NewIam, fx.As(new(services.Authorizer)), fx.ResultTags(`name:"authorizer"`)),

			fx.Annotate(services.NewLoginService, fx.As(new(services.LoginProcessor)), fx.ParamTags(`name:"sessionCreator"`), fx.ResultTags(`name:"loginProcessor"`)),
			fx.Annotate(services.NewLoginService, fx.As(new(user.LoginService)), fx.ParamTags(`name:"sessionCreator"`)),
			fx.Annotate(services.NewAuthenticationService, fx.As(new(user.AuthenticationService)), fx.ParamTags(`name:"authenticator"`, `name:"loginProcessor"`)),
			fx.Annotate(services.NewAuthorizationService, fx.As(new(api.AuthorizationService)), fx.ParamTags(`name:"authorizer"`)),
			fx.Annotate(services.NewRegistrationService, fx.As(new(user.RegistrationService))),

			AsControllers(user.NewLoginController),
			AsControllers(user.NewRegistrtionController),
			AsControllers(user.NewOrdersController),
			AsControllers(user.NewBalanceController),
			AsControllers(user.NewWithdrawalsController),
			AsControllers(balance.NewWithdrawalsController),

			fx.Annotate(api.NewRouter, fx.ParamTags(`group:"contollers"`)),
		),
		fx.Invoke(startHTTPServer, checkDBConnection),
		fx.Supply(
			config.MustNewConfig(),
		),
	)
}

func startHTTPServer(*api.HTTPServer) {}

func checkDBConnection(*storage.Storage) {}

func AsControllers(f any, ants ...fx.Annotation) any {
	ants = append(ants, fx.ResultTags(`group:"contollers"`))
	ants = append(ants, fx.As(new(api.Controller)))

	return fx.Annotate(
		f,
		ants...,
	)
}
