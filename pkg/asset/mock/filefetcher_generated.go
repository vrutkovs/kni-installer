// Code generated by MockGen. DO NOT EDIT.
// Source: ./filefetcher.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	asset "github.com/metalkube/kni-installer/pkg/asset"
	reflect "reflect"
)

// MockFileFetcher is a mock of FileFetcher interface
type MockFileFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockFileFetcherMockRecorder
}

// MockFileFetcherMockRecorder is the mock recorder for MockFileFetcher
type MockFileFetcherMockRecorder struct {
	mock *MockFileFetcher
}

// NewMockFileFetcher creates a new mock instance
func NewMockFileFetcher(ctrl *gomock.Controller) *MockFileFetcher {
	mock := &MockFileFetcher{ctrl: ctrl}
	mock.recorder = &MockFileFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileFetcher) EXPECT() *MockFileFetcherMockRecorder {
	return m.recorder
}

// FetchByName mocks base method
func (m *MockFileFetcher) FetchByName(arg0 string) (*asset.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchByName", arg0)
	ret0, _ := ret[0].(*asset.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchByName indicates an expected call of FetchByName
func (mr *MockFileFetcherMockRecorder) FetchByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchByName", reflect.TypeOf((*MockFileFetcher)(nil).FetchByName), arg0)
}

// FetchByPattern mocks base method
func (m *MockFileFetcher) FetchByPattern(pattern string) ([]*asset.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchByPattern", pattern)
	ret0, _ := ret[0].([]*asset.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchByPattern indicates an expected call of FetchByPattern
func (mr *MockFileFetcherMockRecorder) FetchByPattern(pattern interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchByPattern", reflect.TypeOf((*MockFileFetcher)(nil).FetchByPattern), pattern)
}
