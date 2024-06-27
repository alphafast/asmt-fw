// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/alphafast/asmt-fw/libs/domain/noti (interfaces: NotiRepository)

// Package noti_mock is a generated GoMock package.
package noti_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/alphafast/asmt-fw/libs/domain/noti/model"
	gomock "github.com/golang/mock/gomock"
)

// MockNotiRepository is a mock of NotiRepository interface.
type MockNotiRepository struct {
	ctrl     *gomock.Controller
	recorder *MockNotiRepositoryMockRecorder
}

// MockNotiRepositoryMockRecorder is the mock recorder for MockNotiRepository.
type MockNotiRepositoryMockRecorder struct {
	mock *MockNotiRepository
}

// NewMockNotiRepository creates a new mock instance.
func NewMockNotiRepository(ctrl *gomock.Controller) *MockNotiRepository {
	mock := &MockNotiRepository{ctrl: ctrl}
	mock.recorder = &MockNotiRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotiRepository) EXPECT() *MockNotiRepositoryMockRecorder {
	return m.recorder
}

// FindUserNotification mocks base method.
func (m *MockNotiRepository) FindUserNotification(arg0 context.Context, arg1 string) (*model.NotiUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserNotification", arg0, arg1)
	ret0, _ := ret[0].(*model.NotiUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserNotification indicates an expected call of FindUserNotification.
func (mr *MockNotiRepositoryMockRecorder) FindUserNotification(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserNotification", reflect.TypeOf((*MockNotiRepository)(nil).FindUserNotification), arg0, arg1)
}

// GetNotifyResults mocks base method.
func (m *MockNotiRepository) GetNotifyResults(arg0 context.Context, arg1 string) ([]model.NotiResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotifyResults", arg0, arg1)
	ret0, _ := ret[0].([]model.NotiResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotifyResults indicates an expected call of GetNotifyResults.
func (mr *MockNotiRepositoryMockRecorder) GetNotifyResults(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotifyResults", reflect.TypeOf((*MockNotiRepository)(nil).GetNotifyResults), arg0, arg1)
}

// UpsertNotifyResult mocks base method.
func (m *MockNotiRepository) UpsertNotifyResult(arg0 context.Context, arg1 model.NotiResult) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertNotifyResult", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertNotifyResult indicates an expected call of UpsertNotifyResult.
func (mr *MockNotiRepositoryMockRecorder) UpsertNotifyResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertNotifyResult", reflect.TypeOf((*MockNotiRepository)(nil).UpsertNotifyResult), arg0, arg1)
}
