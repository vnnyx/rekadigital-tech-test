package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/vnnyx/rekadigital-tech-test/internal/delivery/http/web"
	"github.com/vnnyx/rekadigital-tech-test/internal/entity"
	"github.com/vnnyx/rekadigital-tech-test/internal/helper"
	"github.com/vnnyx/rekadigital-tech-test/internal/repository/customer"
	"github.com/vnnyx/rekadigital-tech-test/internal/repository/transaction"
	"github.com/vnnyx/rekadigital-tech-test/internal/validator"
)

type TransactionUCImpl struct {
	transactionRepository transaction.TransactionRepository
	customerRepository    customer.CustomerRepository
}

func NewTransactionUC(transactionRepository transaction.TransactionRepository, customerRepository customer.CustomerRepository) TransactionUC {
	return &TransactionUCImpl{
		transactionRepository: transactionRepository,
		customerRepository:    customerRepository,
	}
}

func (uc *TransactionUCImpl) CreateTransaction(ctx context.Context, req *web.TransactionCreateReq) (*web.TransactionDTO, error) {
	validator.CreateTransactionValidation(*req)
	var customerID string
	total := req.Price * req.Qty

	customer, err := uc.customerRepository.FindCustomerByName(req.CustomerName)
	if err != nil {
		customer = &entity.Customer{
			ID:   uuid.New().String(),
			Name: req.CustomerName,
		}
		err = uc.customerRepository.StoreCustomer(customer)
		if err != nil {
			return nil, err
		}
	}
	customerID = customer.ID

	transaction := &entity.Transaction{
		ID:           uuid.New().String(),
		CustomerID:   customerID,
		CustomerName: req.CustomerName,
		Menu:         req.Menu,
		Price:        req.Price,
		Qty:          req.Qty,
		Payment:      req.Payment,
		Total:        total,
	}
	err = uc.transactionRepository.StoreTransaction(transaction)
	if err != nil {
		return nil, err
	}

	return transaction.ToDTO(), nil
}

func (uc *TransactionUCImpl) GetAllTransaction(ctx context.Context, opt *helper.TransactionOptions) (*web.PaginationDTO, error) {
	got, err := uc.transactionRepository.GetAllTransaction(opt)
	if err != nil {
		return nil, err
	}
	return got.ToPaginationDTO(), nil
}
