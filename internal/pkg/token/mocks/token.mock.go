// Code generated by MockGen. DO NOT EDIT.
// Source: ./token.go

// Package tkmocks is a generated GoMock package.
package tkmocks

import (
	context "context"
	reflect "reflect"

	token "github.com/ogreks/meeseeks-box/internal/pkg/token"
	gomock "go.uber.org/mock/gomock"
)

// MockToken is a mock of Token interface.
type MockToken[T token.Type, V token.Val] struct {
	ctrl     *gomock.Controller
	recorder *MockTokenMockRecorder[T, V]
}

// MockTokenMockRecorder is the mock recorder for MockToken.
type MockTokenMockRecorder[T token.Type, V token.Val] struct {
	mock *MockToken[T, V]
}

// NewMockToken creates a new mock instance.
func NewMockToken[T token.Type, V token.Val](ctrl *gomock.Controller) *MockToken[T, V] {
	mock := &MockToken[T, V]{ctrl: ctrl}
	mock.recorder = &MockTokenMockRecorder[T, V]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockToken[T, V]) EXPECT() *MockTokenMockRecorder[T, V] {
	return m.recorder
}

// CreateToken mocks base method.
func (m *MockToken[T, V]) CreateToken(ctx context.Context, v ...V) (T, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range v {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateToken", varargs...)
	ret0, _ := ret[0].(T)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockTokenMockRecorder[T, V]) CreateToken(ctx interface{}, v ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, v...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockToken[T, V])(nil).CreateToken), varargs...)
}

// RefreshToken mocks base method.
func (m *MockToken[T, V]) RefreshToken(ctx context.Context, token *T, v ...V) (T, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, token}
	for _, a := range v {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RefreshToken", varargs...)
	ret0, _ := ret[0].(T)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockTokenMockRecorder[T, V]) RefreshToken(ctx, token interface{}, v ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, token}, v...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockToken[T, V])(nil).RefreshToken), varargs...)
}

// Store mocks base method.
func (m *MockToken[T, V]) Store() token.Store[T] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store")
	ret0, _ := ret[0].(token.Store[T])
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockTokenMockRecorder[T, V]) Store() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockToken[T, V])(nil).Store))
}

// Validate mocks base method.
func (m *MockToken[T, V]) Validate(token T) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockTokenMockRecorder[T, V]) Validate(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockToken[T, V])(nil).Validate), token)
}
