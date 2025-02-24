package services

import (
	"context"
	"fmt"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"

	"go.uber.org/zap"
)

type RegistrationServiceUsersRepository interface {
	Create(ctx context.Context, u *models.User) error
}

type RegistrationService struct {
	lg  *logging.ZapLogger
	rep RegistrationServiceUsersRepository
}

func NewRegistrationService(rep RegistrationServiceUsersRepository, lg *logging.ZapLogger) *RegistrationService {
	return &RegistrationService{lg: lg, rep: rep}
}

func (service *RegistrationService) Call(ctx context.Context, u *models.User) error {
	ctx = service.lg.WithContextFields(ctx, zap.String("actor", "registrator"))
	if err := u.HashPwd(); err != nil {
		return fmt.Errorf("registration_service password hash failed error %w", err)
	}

	if err := service.rep.Create(ctx, u); err != nil {
		return fmt.Errorf("registration_service create user failed error %w", err)
	}

	return nil
}
