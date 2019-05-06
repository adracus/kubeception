// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/adracus/kubeception/pkg/internal/mock/kubeception/intern/util (interfaces: AddToManager)

// Package util is a generated GoMock package.
package util

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	manager "sigs.k8s.io/controller-runtime/pkg/manager"
)

// MockAddToManager is a mock of AddToManager interface
type MockAddToManager struct {
	ctrl     *gomock.Controller
	recorder *MockAddToManagerMockRecorder
}

// MockAddToManagerMockRecorder is the mock recorder for MockAddToManager
type MockAddToManagerMockRecorder struct {
	mock *MockAddToManager
}

// NewMockAddToManager creates a new mock instance
func NewMockAddToManager(ctrl *gomock.Controller) *MockAddToManager {
	mock := &MockAddToManager{ctrl: ctrl}
	mock.recorder = &MockAddToManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAddToManager) EXPECT() *MockAddToManagerMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockAddToManager) Do(arg0 manager.Manager) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Do indicates an expected call of Do
func (mr *MockAddToManagerMockRecorder) Do(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockAddToManager)(nil).Do), arg0)
}
