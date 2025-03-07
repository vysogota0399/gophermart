package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/vysogota0399/gophermart_portal/internal/api"
	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	mocks "github.com/vysogota0399/gophermart_portal/internal/api/mocks/controllers/user"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/api/repositories"
	"github.com/vysogota0399/gophermart_portal/internal/config"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegistrtionController_reginstrationHandler(t *testing.T) {
	type fields struct {
		registrationService *mocks.MockRegistrationService
		loginService        *mocks.MockLoginService
		authService         *mocks.MockAuthorizationService
	}
	type want struct {
		status   int
		response string
		headers  map[string]string
	}
	tests := []struct {
		name    string
		payload string
		fields  fields
		prepare func(f *fields)
		want    want
	}{
		{
			name:    "when invalid params then return 400",
			payload: "{}",
			prepare: func(f *fields) {},
			want: want{
				status:   http.StatusBadRequest,
				response: `{"error": "invalid request params"}`,
				headers:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name:    "when user already exists then return 400",
			payload: `{"login": "test", "password": "secret"}`,
			prepare: func(f *fields) {
				f.registrationService.EXPECT().Call(gomock.Any(), gomock.Any()).Return(repositories.ErrLoginAlreadyExist)
			},
			want: want{
				status:   http.StatusBadRequest,
				response: `{"error": "user already exists"}`,
				headers:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name:    "when service return error then return 500",
			payload: `{"login": "test", "password": "secret"}`,
			prepare: func(f *fields) {
				f.registrationService.EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("fatal error"))
			},
			want: want{
				status:   http.StatusInternalServerError,
				response: `{"error": "cantact with operator"}`,
				headers:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name:    "when login failed then return 401",
			payload: `{"login": "test", "password": "secret"}`,
			prepare: func(f *fields) {
				f.registrationService.EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
				f.loginService.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("fatal error"))
			},
			want: want{
				status:   http.StatusUnauthorized,
				response: `{"error": "unauthorized"}`,
				headers:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name:    "when login succeeded then return 200",
			payload: `{"login": "test", "password": "secret"}`,
			prepare: func(f *fields) {
				f.registrationService.EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
				f.loginService.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, w http.ResponseWriter, u *models.User) error {
						w.Header().Add("Authorization", "Bearer test")
						return nil
					},
				)
			},
			want: want{
				status:   http.StatusOK,
				response: `{}`,
				headers: map[string]string{
					"Content-Type":  "application/json",
					"Authorization": "Bearer test",
				},
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
			fields.registrationService = mocks.NewMockRegistrationService(ctrl)
			fields.loginService = mocks.NewMockLoginService(ctrl)
			fields.authService = mocks.NewMockAuthorizationService(ctrl)
			tt.prepare(&fields)

			cntr := NewRegistrtionController(fields.registrationService, lg, fields.loginService)
			router := api.NewRouter([]api.Controller{cntr}, fields.authService, lg)
			srv := api.NewTestHTTPServer("testserver", router)
			assert.NoError(t, err)
			defer srv.Srv.Close()

			req := resty.New().R()
			req.Method = http.MethodPost
			req.SetBody(tt.payload)
			req.URL = fmt.Sprintf("%s%s", srv.Srv.URL, cntr.Path)
			req.Header.Add("Content-Type", "application/json")

			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.want.status)
			assert.JSONEq(t, string(resp.Body()), tt.want.response)
			for k, v := range tt.want.headers {
				assert.Equal(t, v, resp.Header().Get(k))
			}
		})
	}
}
