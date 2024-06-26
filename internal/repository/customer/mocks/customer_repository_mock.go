// Code generated by MockGen. DO NOT EDIT.
// Source: customer_repository.go

// Package mock_customer is a generated GoMock package.
package mock_customer

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/vnnyx/rekadigital-tech-test/internal/entity"
)

// MockCustomerRepository is a mock of CustomerRepository interface.
type MockCustomerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCustomerRepositoryMockRecorder
}

// MockCustomerRepositoryMockRecorder is the mock recorder for MockCustomerRepository.
type MockCustomerRepositoryMockRecorder struct {
	mock *MockCustomerRepository
}

// NewMockCustomerRepository creates a new mock instance.
func NewMockCustomerRepository(ctrl *gomock.Controller) *MockCustomerRepository {
	mock := &MockCustomerRepository{ctrl: ctrl}
	mock.recorder = &MockCustomerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCustomerRepository) EXPECT() *MockCustomerRepositoryMockRecorder {
	return m.recorder
}

// FindCustomerByName mocks base method.
func (m *MockCustomerRepository) FindCustomerByName(name string) (*entity.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCustomerByName", name)
	ret0, _ := ret[0].(*entity.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCustomerByName indicates an expected call of FindCustomerByName.
func (mr *MockCustomerRepositoryMockRecorder) FindCustomerByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCustomerByName", reflect.TypeOf((*MockCustomerRepository)(nil).FindCustomerByName), name)
}

// StoreCustomer mocks base method.
func (m *MockCustomerRepository) StoreCustomer(customer *entity.Customer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreCustomer", customer)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreCustomer indicates an expected call of StoreCustomer.
func (mr *MockCustomerRepositoryMockRecorder) StoreCustomer(customer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreCustomer", reflect.TypeOf((*MockCustomerRepository)(nil).StoreCustomer), customer)
}
