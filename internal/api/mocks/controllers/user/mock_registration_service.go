// Code generated by MockGen. DO NOT EDIT.
// Source: gophermart/internal/api/controllers/user (interfaces: RegistrationService)

// Package mocks is a generated GoMock package.
package user

import (
	context "context"
	models "github.com/vysogota0399/gophermart_portal/internal/api/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRegistrationService is a mock of RegistrationService interface.
type MockRegistrationService struct {
	ctrl     *gomock.Controller
	recorder *MockRegistrationServiceMockRecorder
}

// MockRegistrationServiceMockRecorder is the mock recorder for MockRegistrationService.
type MockRegistrationServiceMockRecorder struct {
	mock *MockRegistrationService
}

// NewMockRegistrationService creates a new mock instance.
func NewMockRegistrationService(ctrl *gomock.Controller) *MockRegistrationService {
	mock := &MockRegistrationService{ctrl: ctrl}
	mock.recorder = &MockRegistrationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegistrationService) EXPECT() *MockRegistrationServiceMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockRegistrationService) Call(arg0 context.Context, arg1 *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Call indicates an expected call of Call.
func (mr *MockRegistrationServiceMockRecorder) Call(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockRegistrationService)(nil).Call), arg0, arg1)
}
