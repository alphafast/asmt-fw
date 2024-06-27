// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/alphafast/asmt-fw/libs/domain/noti (interfaces: NotiAdapter)

// Package noti_mock is a generated GoMock package.
package noti_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/alphafast/asmt-fw/libs/domain/noti/model"
	gomock "github.com/golang/mock/gomock"
)

// MockNotiAdapter is a mock of NotiAdapter interface.
type MockNotiAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockNotiAdapterMockRecorder
}

// MockNotiAdapterMockRecorder is the mock recorder for MockNotiAdapter.
type MockNotiAdapterMockRecorder struct {
	mock *MockNotiAdapter
}

// NewMockNotiAdapter creates a new mock instance.
func NewMockNotiAdapter(ctrl *gomock.Controller) *MockNotiAdapter {
	mock := &MockNotiAdapter{ctrl: ctrl}
	mock.recorder = &MockNotiAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotiAdapter) EXPECT() *MockNotiAdapterMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockNotiAdapter) Send(arg0 context.Context, arg1 model.NotiRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockNotiAdapterMockRecorder) Send(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockNotiAdapter)(nil).Send), arg0, arg1)
}
