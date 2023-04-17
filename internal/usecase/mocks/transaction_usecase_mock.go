// Code generated by MockGen. DO NOT EDIT.
// Source: transaction_usecase.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	web "github.com/vnnyx/rekadigital-tech-test/internal/delivery/http/web"
	helper "github.com/vnnyx/rekadigital-tech-test/internal/helper"
)

// MockTransactionUC is a mock of TransactionUC interface.
type MockTransactionUC struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionUCMockRecorder
}

// MockTransactionUCMockRecorder is the mock recorder for MockTransactionUC.
type MockTransactionUCMockRecorder struct {
	mock *MockTransactionUC
}

// NewMockTransactionUC creates a new mock instance.
func NewMockTransactionUC(ctrl *gomock.Controller) *MockTransactionUC {
	mock := &MockTransactionUC{ctrl: ctrl}
	mock.recorder = &MockTransactionUCMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionUC) EXPECT() *MockTransactionUCMockRecorder {
	return m.recorder
}

// CreateTransaction mocks base method.
func (m *MockTransactionUC) CreateTransaction(ctx context.Context, req *web.TransactionCreateReq) (*web.TransactionDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", ctx, req)
	ret0, _ := ret[0].(*web.TransactionDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockTransactionUCMockRecorder) CreateTransaction(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockTransactionUC)(nil).CreateTransaction), ctx, req)
}

// GetAllTransaction mocks base method.
func (m *MockTransactionUC) GetAllTransaction(ctx context.Context, opt *helper.TransactionOptions) (*web.PaginationDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTransaction", ctx, opt)
	ret0, _ := ret[0].(*web.PaginationDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTransaction indicates an expected call of GetAllTransaction.
func (mr *MockTransactionUCMockRecorder) GetAllTransaction(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTransaction", reflect.TypeOf((*MockTransactionUC)(nil).GetAllTransaction), ctx, opt)
}