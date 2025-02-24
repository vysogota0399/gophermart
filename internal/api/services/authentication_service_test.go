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

func TestAuthenticationService_Call(t *testing.T) {
	type fields struct {
		authenticator *mock.MockAuthenticator
		logging       *mock.MockLoginProcessor
	}
	type args struct {
		ctx context.Context
		u   *models.User
		rw  http.ResponseWriter
	}
	tests := []struct {
		name       string
		fields     fields
		prepare    func(f *fields)
		args       args
		wantErr    bool
		wantHeader string
	}{
		{
			name: "when success authentication and login",
			args: args{
				u:   &models.User{},
				ctx: context.Background(),
				rw:  &customResponseWriter{header: http.Header{}},
			},
			prepare: func(f *fields) {
				f.authenticator.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(&models.User{}, nil)
				f.logging.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:    "when success authentication failed then error",
			wantErr: true,
			args: args{
				u:   &models.User{},
				ctx: context.Background(),
				rw:  &customResponseWriter{header: http.Header{}},
			},
			prepare: func(f *fields) {
				f.authenticator.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(&models.User{}, errors.New("error"))
			},
		},
		{
			name:    "when login failed then error",
			wantErr: true,
			args: args{
				u:   &models.User{},
				ctx: context.Background(),
				rw:  &customResponseWriter{header: http.Header{}},
			},
			prepare: func(f *fields) {
				f.authenticator.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(&models.User{}, nil)
				f.logging.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
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
			fields.authenticator = mock.NewMockAuthenticator(ctrl)
			fields.logging = mock.NewMockLoginProcessor(ctrl)
			tt.prepare(&fields)

			rw := &customResponseWriter{header: http.Header{}}

			err := NewAuthenticationService(fields.authenticator, fields.logging, lg).Call(tt.args.ctx, rw, tt.args.u)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
