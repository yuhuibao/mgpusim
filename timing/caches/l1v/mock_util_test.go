// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/akita/util/v2/buffering (interfaces: Buffer)

package l1v

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBuffer is a mock of Buffer interface.
type MockBuffer struct {
	ctrl     *gomock.Controller
	recorder *MockBufferMockRecorder
}

// MockBufferMockRecorder is the mock recorder for MockBuffer.
type MockBufferMockRecorder struct {
	mock *MockBuffer
}

// NewMockBuffer creates a new mock instance.
func NewMockBuffer(ctrl *gomock.Controller) *MockBuffer {
	mock := &MockBuffer{ctrl: ctrl}
	mock.recorder = &MockBufferMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBuffer) EXPECT() *MockBufferMockRecorder {
	return m.recorder
}

// CanPush mocks base method.
func (m *MockBuffer) CanPush() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanPush")
	ret0, _ := ret[0].(bool)
	return ret0
}

// CanPush indicates an expected call of CanPush.
func (mr *MockBufferMockRecorder) CanPush() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanPush", reflect.TypeOf((*MockBuffer)(nil).CanPush))
}

// Capacity mocks base method.
func (m *MockBuffer) Capacity() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Capacity")
	ret0, _ := ret[0].(int)
	return ret0
}

// Capacity indicates an expected call of Capacity.
func (mr *MockBufferMockRecorder) Capacity() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capacity", reflect.TypeOf((*MockBuffer)(nil).Capacity))
}

// Clear mocks base method.
func (m *MockBuffer) Clear() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Clear")
}

// Clear indicates an expected call of Clear.
func (mr *MockBufferMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockBuffer)(nil).Clear))
}

// Peek mocks base method.
func (m *MockBuffer) Peek() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peek")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Peek indicates an expected call of Peek.
func (mr *MockBufferMockRecorder) Peek() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peek", reflect.TypeOf((*MockBuffer)(nil).Peek))
}

// Pop mocks base method.
func (m *MockBuffer) Pop() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pop")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Pop indicates an expected call of Pop.
func (mr *MockBufferMockRecorder) Pop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pop", reflect.TypeOf((*MockBuffer)(nil).Pop))
}

// Push mocks base method.
func (m *MockBuffer) Push(arg0 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Push", arg0)
}

// Push indicates an expected call of Push.
func (mr *MockBufferMockRecorder) Push(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockBuffer)(nil).Push), arg0)
}

// Size mocks base method.
func (m *MockBuffer) Size() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size")
	ret0, _ := ret[0].(int)
	return ret0
}

// Size indicates an expected call of Size.
func (mr *MockBufferMockRecorder) Size() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockBuffer)(nil).Size))
}
