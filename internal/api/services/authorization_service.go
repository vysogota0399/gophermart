package services

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
)

type Authorizer interface {
	Authorize(ctx context.Context, token string) (*models.Session, error)
}
type AuthorizationService struct {
	lg  *logging.ZapLogger
	iam Authorizer
}

func NewAuthorizationService(iam Authorizer, lg *logging.ZapLogger) *AuthorizationService {
	return &AuthorizationService{lg: lg, iam: iam}
}

func (service *AuthorizationService) Call(ctx context.Context, req *http.Request) (*models.Session, error) {
	header := req.Header.Get("Authorization")
	if header == "" {
		return nil, fmt.Errorf("internal/api/services/authorization_service: auth header is empty error")
	}

	token := service.findToken(header)
	if token == "" {
		return nil, fmt.Errorf("internal/api/services/authorization_service: token not found in auth header error")
	}

	return service.iam.Authorize(ctx, token)
}

func (service *AuthorizationService) findToken(input string) string {
	re := regexp.MustCompile(`Bearer\s+(\S+)`)
	match := re.FindStringSubmatch(input)

	if match == nil || len(match) < 2 {
		return ""
	}

	return match[1]
}
