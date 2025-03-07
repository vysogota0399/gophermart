package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
)

type AuthenticationService struct {
	lg           *logging.ZapLogger
	iam          Authenticator
	loginService LoginProcessor
}

type LoginProcessor interface {
	Call(ctx context.Context, w http.ResponseWriter, u *models.User) error
}

type Authenticator interface {
	Authenticate(ctx context.Context, u *models.User) (*models.User, error)
}

func NewAuthenticationService(iam Authenticator, loginService LoginProcessor, lg *logging.ZapLogger) *AuthenticationService {
	return &AuthenticationService{lg: lg, iam: iam, loginService: loginService}
}

func (service *AuthenticationService) Call(ctx context.Context, w http.ResponseWriter, u *models.User) error {
	authorizedUser, err := service.iam.Authenticate(ctx, u)
	if err != nil {
		return fmt.Errorf("internal/api/services/authentication_service: authentication failed error %w", err)
	}

	return service.loginService.Call(ctx, w, authorizedUser)
}
