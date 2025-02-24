package services

import (
	"context"
	"errors"
	"testing"

	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	mock "github.com/vysogota0399/gophermart_portal/internal/api/mocks/services"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/config"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegistrationService_Call(t *testing.T) {
	type fields struct {
		rep *mock.MockRegistrationServiceUsersRepository
	}
	type args struct {
		ctx context.Context
		u   *models.User
	}
	tests := []struct {
		name    string
		fields  fields
		prepare func(ctx context.Context, f *fields, u *models.User)
		args    args
		wantErr bool
	}{
		{
			name: "when user created then no errors",
			prepare: func(ctx context.Context, f *fields, u *models.User) {
				f.rep.EXPECT().Create(gomock.Any(), u).Return(nil)
			},
			args: args{
				ctx: context.Background(),
				u:   &models.User{},
			},
		},
		{
			name: "when user creation failed then return error",
			prepare: func(ctx context.Context, f *fields, u *models.User) {
				f.rep.EXPECT().Create(gomock.Any(), u).Return(errors.New("error"))
			},
			args: args{
				ctx: context.Background(),
				u:   &models.User{},
			},
			wantErr: true,
		},
	}

	lg, err := logging.NewZapLogger(&config.Config{})
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := tt.fields
			fields.rep = mock.NewMockRegistrationServiceUsersRepository(ctrl)
			tt.prepare(tt.args.ctx, &fields, tt.args.u)

			res := NewRegistrationService(fields.rep, lg).Call(tt.args.ctx, tt.args.u)
			assert.Equal(t, tt.wantErr, res != nil)
		})
	}
}
