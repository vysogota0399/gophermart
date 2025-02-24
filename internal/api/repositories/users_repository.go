package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/storage"

	"go.uber.org/zap"
)

type UsersRepository struct {
	strg Storage
	lg   *logging.ZapLogger
}

type Storage interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func NewUsersRepository(strg *storage.Storage, lg *logging.ZapLogger) *UsersRepository {
	return &UsersRepository{strg: strg.DB, lg: lg}
}

func (rep *UsersRepository) Find(ctx context.Context, id uint64) (*models.User, error) {
	var u models.User
	row := rep.strg.QueryRowContext(
		ctx,
		`select id, login, password, created_at from users where id = $1 limit 1`,
		id,
	)

	if err := row.Scan(&u.ID, &u.Login, &u.Password, &u.CreatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}

func (rep *UsersRepository) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	var u models.User
	row := rep.strg.QueryRowContext(
		ctx,
		`select id, login, password, created_at from users where login = $1 limit 1`,
		login,
	)

	if err := row.Scan(&u.ID, &u.Login, &u.Password, &u.CreatedAt); err != nil {
		return nil, err
	}

	rep.lg.DebugCtx(ctx, "user found", zap.Any("user", u))
	return &u, nil
}

var ErrLoginAlreadyExist = fmt.Errorf("login already exist")

func (rep *UsersRepository) Create(ctx context.Context, u *models.User) error {
	rep.lg.DebugCtx(ctx, "create user", zap.Any("user", *u))
	row := rep.strg.QueryRowContext(
		ctx,
		`
			insert into users(login, password)
			values ($1, $2) on conflict do nothing
			returning id
		`,
		u.Login,
		u.Password,
	)

	if err := row.Scan(&u.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrLoginAlreadyExist
		}

		return err
	}

	return nil
}
