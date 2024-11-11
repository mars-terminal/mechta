// Code generated by MockGen. DO NOT EDIT.
// Source: shortener.go
//
// Generated by this command:
//
//	mockgen -source=shortener.go -destination shortener_mock.gen.go -package service
//

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	domain "github.com/mars-terminal/mechta/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockShortener is a mock of Shortener interface.
type MockShortener struct {
	ctrl     *gomock.Controller
	recorder *MockShortenerMockRecorder
	isgomock struct{}
}

// MockShortenerMockRecorder is the mock recorder for MockShortener.
type MockShortenerMockRecorder struct {
	mock *MockShortener
}

// NewMockShortener creates a new mock instance.
func NewMockShortener(ctrl *gomock.Controller) *MockShortener {
	mock := &MockShortener{ctrl: ctrl}
	mock.recorder = &MockShortenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShortener) EXPECT() *MockShortenerMockRecorder {
	return m.recorder
}

// CreateShortLink mocks base method.
func (m *MockShortener) CreateShortLink(ctx context.Context, cmd CreateLinkCMD) (domain.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShortLink", ctx, cmd)
	ret0, _ := ret[0].(domain.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateShortLink indicates an expected call of CreateShortLink.
func (mr *MockShortenerMockRecorder) CreateShortLink(ctx, cmd any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShortLink", reflect.TypeOf((*MockShortener)(nil).CreateShortLink), ctx, cmd)
}

// DeleteLink mocks base method.
func (m *MockShortener) DeleteLink(ctx context.Context, shortURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLink", ctx, shortURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLink indicates an expected call of DeleteLink.
func (mr *MockShortenerMockRecorder) DeleteLink(ctx, shortURL any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLink", reflect.TypeOf((*MockShortener)(nil).DeleteLink), ctx, shortURL)
}

// GetLinkStatistics mocks base method.
func (m *MockShortener) GetLinkStatistics(ctx context.Context, shortURL string) (domain.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinkStatistics", ctx, shortURL)
	ret0, _ := ret[0].(domain.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinkStatistics indicates an expected call of GetLinkStatistics.
func (mr *MockShortenerMockRecorder) GetLinkStatistics(ctx, shortURL any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinkStatistics", reflect.TypeOf((*MockShortener)(nil).GetLinkStatistics), ctx, shortURL)
}

// GetLinks mocks base method.
func (m *MockShortener) GetLinks(ctx context.Context) ([]domain.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinks", ctx)
	ret0, _ := ret[0].([]domain.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinks indicates an expected call of GetLinks.
func (mr *MockShortenerMockRecorder) GetLinks(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinks", reflect.TypeOf((*MockShortener)(nil).GetLinks), ctx)
}

// RedirectLink mocks base method.
func (m *MockShortener) RedirectLink(ctx context.Context, shortURL string) (domain.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RedirectLink", ctx, shortURL)
	ret0, _ := ret[0].(domain.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RedirectLink indicates an expected call of RedirectLink.
func (mr *MockShortenerMockRecorder) RedirectLink(ctx, shortURL any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RedirectLink", reflect.TypeOf((*MockShortener)(nil).RedirectLink), ctx, shortURL)
}
