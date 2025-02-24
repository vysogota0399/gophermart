package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/api/repositories"
	"github.com/vysogota0399/gophermart_portal/internal/config"

	"golang.org/x/crypto/bcrypt"
)

type Iam struct {
	TokenTTL  time.Duration
	SecretKey string
	sessRep   *repositories.SessionsRepository
	usersRep  *repositories.UsersRepository
	lg        *logging.ZapLogger
}

func NewIam(cfg *config.Config, sessRep *repositories.SessionsRepository, usersRep *repositories.UsersRepository, lg *logging.ZapLogger) *Iam {
	return &Iam{
		TokenTTL:  time.Duration(cfg.IamTokenTTL) * time.Minute,
		SecretKey: cfg.IamSecretKey,
		sessRep:   sessRep,
		usersRep:  usersRep,
		lg:        lg,
	}
}

func (i *Iam) Authenticate(ctx context.Context, u *models.User) (*models.User, error) {
	foundUser, err := i.usersRep.FindByLogin(ctx, u.Login)
	if err != nil {
		return nil, fmt.Errorf("internal/api/iam/iam: user not found %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(u.Password)); err != nil {
		return nil, fmt.Errorf("internal/api/iam/iam: invalid password %w", err)
	}

	return foundUser, nil
}

func (i *Iam) Login(ctx context.Context, u *models.User) (string, error) {
	session := &models.Session{
		Sub:       fmt.Sprintf("%d", u.ID),
		ExpiredAt: time.Now().Add(i.TokenTTL),
	}
	if err := i.sessRep.Create(ctx, session); err != nil {
		return "", fmt.Errorf("internal/api/iam/iam: create session failed error %w", err)
	}

	return i.buildJWTString(session)
}

func (i *Iam) Authorize(ctx context.Context, token string) (*models.Session, error) {
	claims, err := i.decode(token)
	if err != nil {
		return nil, fmt.Errorf("internal/api/iam/iam: authorize error %w", err)
	}

	sess, err := i.sessRep.Find(ctx, claims.Sid)
	if err != nil {
		return nil, fmt.Errorf("internal/api/iam/iam: session not found %w", err)
	}

	if time.Now().After(sess.ExpiredAt) {
		return nil, fmt.Errorf("internal/api/iam/iam: session expired")
	}

	return sess, err
}
