package transaction

import (
	"github.com/vnnyx/rekadigital-tech-test/internal/entity"
	"github.com/vnnyx/rekadigital-tech-test/internal/helper"
)

type TransactionRepository interface {
	StoreTransaction(transaction *entity.Transaction) error
	GetAllTransaction(opt *helper.TransactionOptions) (*helper.Pagination, error)
}
