package customer

import (
	"github.com/vnnyx/rekadigital-tech-test/internal/entity"
)

type CustomerRepository interface {
	StoreCustomer(customer *entity.Customer) error
	FindCustomerByName(name string) (*entity.Customer, error)
}
