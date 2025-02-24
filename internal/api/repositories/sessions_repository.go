package repositories

import (
	"context"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/storage"

	"go.uber.org/zap"
)

type SessionsRepository struct {
	strg Storage
	lg   *logging.ZapLogger
}

func NewSessionsRepository(strg *storage.Storage, lg *logging.ZapLogger) *SessionsRepository {
	return &SessionsRepository{strg: strg.DB, lg: lg}
}

func (rep *SessionsRepository) Find(ctx context.Context, sid uint64) (*models.Session, error) {
	row := rep.strg.QueryRowContext(
		ctx,
		`select id, sub, created_at, expired_at from sessions where id = $1 limit 1`,
		sid,
	)

	s := &models.Session{}
	if err := row.Scan(&s.ID, &s.Sub, &s.CreatedAt, &s.ExpiredAt); err != nil {
		return s, err
	}

	return s, nil
}

func (rep *SessionsRepository) Create(ctx context.Context, s *models.Session) error {
	rep.lg.DebugCtx(ctx, "create session", zap.Any("session", *s))
	row := rep.strg.QueryRowContext(
		ctx,
		`
			insert into sessions(sub, expired_at)
			values ($1, $2)
			returning id, created_at
		`,
		s.Sub,
		s.ExpiredAt,
	)
	if err := row.Scan(&s.ID, &s.CreatedAt); err != nil {
		return err
	}

	return nil
}
