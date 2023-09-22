// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qiangyt/go-comm/v2 (interfaces: File)

// Package test is a generated GoMock package.
package test

import (
	url "net/url"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	models "github.com/goodsru/go-universal-network-adapter/models"
)

// MockFile is a mock of File interface.
type MockFile struct {
	ctrl     *gomock.Controller
	recorder *MockFileMockRecorder
}

// MockFileMockRecorder is the mock recorder for MockFile.
type MockFileMockRecorder struct {
	mock *MockFile
}

// NewMockFile creates a new mock instance.
func NewMockFile(ctrl *gomock.Controller) *MockFile {
	mock := &MockFile{ctrl: ctrl}
	mock.recorder = &MockFileMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFile) EXPECT() *MockFileMockRecorder {
	return m.recorder
}

// Credentials mocks base method.
func (m *MockFile) Credentials() *models.Credentials {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Credentials")
	ret0, _ := ret[0].(*models.Credentials)
	return ret0
}

// Credentials indicates an expected call of Credentials.
func (mr *MockFileMockRecorder) Credentials() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Credentials", reflect.TypeOf((*MockFile)(nil).Credentials))
}

// Dir mocks base method.
func (m *MockFile) Dir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dir")
	ret0, _ := ret[0].(string)
	return ret0
}

// Dir indicates an expected call of Dir.
func (mr *MockFileMockRecorder) Dir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dir", reflect.TypeOf((*MockFile)(nil).Dir))
}

// Download mocks base method.
func (m *MockFile) Download() (*models.RemoteFileContent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download")
	ret0, _ := ret[0].(*models.RemoteFileContent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockFileMockRecorder) Download() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockFile)(nil).Download))
}

// DownloadP mocks base method.
func (m *MockFile) DownloadP() *models.RemoteFileContent {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadP")
	ret0, _ := ret[0].(*models.RemoteFileContent)
	return ret0
}

// DownloadP indicates an expected call of DownloadP.
func (mr *MockFileMockRecorder) DownloadP() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadP", reflect.TypeOf((*MockFile)(nil).DownloadP))
}

// Name mocks base method.
func (m *MockFile) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockFileMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockFile)(nil).Name))
}

// Protocol mocks base method.
func (m *MockFile) Protocol() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Protocol")
	ret0, _ := ret[0].(string)
	return ret0
}

// Protocol indicates an expected call of Protocol.
func (mr *MockFileMockRecorder) Protocol() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Protocol", reflect.TypeOf((*MockFile)(nil).Protocol))
}

// Timeout mocks base method.
func (m *MockFile) Timeout() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Timeout")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// Timeout indicates an expected call of Timeout.
func (mr *MockFileMockRecorder) Timeout() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Timeout", reflect.TypeOf((*MockFile)(nil).Timeout))
}

// URL mocks base method.
func (m *MockFile) URL() *url.URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "URL")
	ret0, _ := ret[0].(*url.URL)
	return ret0
}

// URL indicates an expected call of URL.
func (mr *MockFileMockRecorder) URL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "URL", reflect.TypeOf((*MockFile)(nil).URL))
}

// Url mocks base method.
func (m *MockFile) Url() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Url")
	ret0, _ := ret[0].(string)
	return ret0
}

// Url indicates an expected call of Url.
func (mr *MockFileMockRecorder) Url() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Url", reflect.TypeOf((*MockFile)(nil).Url))
}
