package services

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	mock "github.com/vysogota0399/gophermart_portal/internal/api/mocks/services"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/config"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type customResponseWriter struct {
	http.ResponseWriter
	header http.Header
}

func (w *customResponseWriter) Header() http.Header {
	return w.header
}

func TestLoginService_Call(t *testing.T) {
	type fields struct {
		service *mock.MockSessionCreator
	}
	type args struct {
		ctx context.Context
		u   *models.User
	}
	tests := []struct {
		name       string
		fields     fields
		prepare    func(ctx context.Context, u *models.User, f *fields)
		args       args
		wantErr    bool
		wantHeader string
	}{
		{
			name: "when success authentication",
			args: args{
				u:   &models.User{},
				ctx: context.Background(),
			},
			prepare: func(ctx context.Context, u *models.User, f *fields) {
				f.service.EXPECT().Login(ctx, u).Return("valid_token", nil)
			},
			wantHeader: "Bearer valid_token",
		},
		{
			name: "when failed authentication",
			args: args{
				u:   &models.User{},
				ctx: context.Background(),
			},
			wantErr: true,
			prepare: func(ctx context.Context, u *models.User, f *fields) {
				f.service.EXPECT().Login(ctx, u).Return("", errors.New("error"))
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
			fields.service = mock.NewMockSessionCreator(ctrl)
			tt.prepare(tt.args.ctx, tt.args.u, &fields)

			rw := &customResponseWriter{header: http.Header{}}

			err := NewLoginService(fields.service, lg).Call(tt.args.ctx, rw, tt.args.u)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantHeader, rw.header.Get("Authorization"))
		})
	}
}
