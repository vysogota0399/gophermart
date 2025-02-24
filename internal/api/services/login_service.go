package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
)

type SessionCreator interface {
	Login(ctx context.Context, u *models.User) (string, error)
}

type LoginService struct {
	lg  *logging.ZapLogger
	iam SessionCreator
}

func NewLoginService(iam SessionCreator, lg *logging.ZapLogger) *LoginService {
	return &LoginService{lg: lg, iam: iam}
}

func (service *LoginService) Call(ctx context.Context, w http.ResponseWriter, u *models.User) error {
	token, err := service.iam.Login(ctx, u)
	if err != nil {
		return fmt.Errorf("internal/api/services/authentication_service: authentication failed error %w", err)
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}
