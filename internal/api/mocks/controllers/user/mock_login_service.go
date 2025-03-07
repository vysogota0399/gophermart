// Code generated by MockGen. DO NOT EDIT.
// Source: gophermart/internal/api/controllers/user (interfaces: LoginService)

// Package mocks is a generated GoMock package.
package user

import (
	context "context"
	models "github.com/vysogota0399/gophermart_portal/internal/api/models"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLoginService is a mock of LoginService interface.
type MockLoginService struct {
	ctrl     *gomock.Controller
	recorder *MockLoginServiceMockRecorder
}

// MockLoginServiceMockRecorder is the mock recorder for MockLoginService.
type MockLoginServiceMockRecorder struct {
	mock *MockLoginService
}

// NewMockLoginService creates a new mock instance.
func NewMockLoginService(ctrl *gomock.Controller) *MockLoginService {
	mock := &MockLoginService{ctrl: ctrl}
	mock.recorder = &MockLoginServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoginService) EXPECT() *MockLoginServiceMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockLoginService) Call(arg0 context.Context, arg1 http.ResponseWriter, arg2 *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Call indicates an expected call of Call.
func (mr *MockLoginServiceMockRecorder) Call(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockLoginService)(nil).Call), arg0, arg1, arg2)
}
