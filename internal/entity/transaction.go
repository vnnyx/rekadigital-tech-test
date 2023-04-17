package entity

import "github.com/vnnyx/rekadigital-tech-test/internal/delivery/http/web"

type Transaction struct {
	ID           string
	CustomerID   string
	CustomerName string
	Menu         string
	Price        int64
	Qty          int64
	Payment      string
	Total        int64
	CreatedAt    int64
}

func (t *Transaction) ToDTO() *web.TransactionDTO {
	return &web.TransactionDTO{
		TransactionID: t.ID,
		CustomerName:  t.CustomerName,
		Menu:          t.Menu,
		Price:         t.Price,
		Qty:           t.Qty,
		Payment:       t.Payment,
		Total:         t.Total,
	}
}
