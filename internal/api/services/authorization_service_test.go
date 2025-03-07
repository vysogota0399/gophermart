package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	mock "github.com/vysogota0399/gophermart_portal/internal/api/mocks/services"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/config"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizationService_Call(t *testing.T) {
	type fields struct {
		authorizer *mock.MockAuthorizer
	}
	type args struct {
		ctx     context.Context
		token   string
		session *models.Session
		request func(string) *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		prepare func(ctx context.Context, token string, session *models.Session, f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "when authorized",
			args: args{
				token:   "token",
				session: &models.Session{},
				ctx:     context.Background(),
				request: func(token string) *http.Request {
					b := []byte{}
					req, _ := http.NewRequestWithContext(context.Background(), "", "", bytes.NewReader(b))
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
					return req
				},
			},
			prepare: func(ctx context.Context, token string, session *models.Session, f *fields) {
				f.authorizer.EXPECT().Authorize(ctx, token).Return(session, nil)
			},
		},
		{
			name:    "when header has no key then error",
			wantErr: true,
			args: args{
				ctx: context.Background(),
				request: func(token string) *http.Request {
					b := []byte{}
					req, _ := http.NewRequestWithContext(context.Background(), "", "", bytes.NewReader(b))
					return req
				},
			},
			prepare: func(ctx context.Context, token string, session *models.Session, f *fields) {},
		},
		{
			name:    "when header has invalid token then error",
			wantErr: true,
			args: args{
				token: "invalid_token",
				ctx:   context.Background(),
				request: func(token string) *http.Request {
					b := []byte{}
					req, _ := http.NewRequestWithContext(context.Background(), "", "", bytes.NewReader(b))
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
					return req
				},
			},
			prepare: func(ctx context.Context, token string, session *models.Session, f *fields) {
				f.authorizer.EXPECT().Authorize(ctx, token).Return(nil, errors.New("error"))
			},
		},
	}
	lg, err := logging.NewZapLogger(&config.Config{})
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := tt.fields
			fields.authorizer = mock.NewMockAuthorizer(ctrl)
			tt.prepare(tt.args.ctx, tt.args.token, tt.args.session, &fields)

			actualSession, err := NewAuthorizationService(fields.authorizer, lg).Call(tt.args.ctx, tt.args.request(tt.args.token))
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.args.session, actualSession)
		})
	}
}
