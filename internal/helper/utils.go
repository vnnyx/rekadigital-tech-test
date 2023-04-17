package helper

import (
	"math"

	"github.com/vnnyx/rekadigital-tech-test/internal/delivery/http/web"
	"github.com/vnnyx/rekadigital-tech-test/internal/entity"
)

type TransactionOptions struct {
	Limit        int64
	Page         int64
	Query        string
	CustomerName string
}

type Pagination struct {
	TotalRows int64
	Limit     int64
	Page      int64
	Rows      any
}

func (p *Pagination) ToPaginationDTO() *web.PaginationDTO {
	if p.TotalRows == 0 {
		return &web.PaginationDTO{
			TotalRows:   0,
			Limit:       0,
			CurrentPage: 0,
			TotalPages:  0,
			Rows:        []any{},
		}
	}
	totalPages := int(math.Ceil(float64(p.TotalRows) / float64(p.Limit)))
	transactionList := make([]*entity.Transaction, 0)

	switch rows := p.Rows.(type) {
	case []*entity.Transaction:
		transactionList = rows
	case []interface{}:
		for _, row := range rows {
			if transactionMap, ok := row.(map[string]interface{}); ok {
				t := &entity.Transaction{
					ID:           transactionMap["ID"].(string),
					CustomerName: transactionMap["CustomerName"].(string),
					Menu:         transactionMap["Menu"].(string),
					Price:        int64(transactionMap["Price"].(float64)),
					Qty:          int64(transactionMap["Qty"].(float64)),
					Payment:      transactionMap["Payment"].(string),
					Total:        int64(transactionMap["Total"].(float64)),
					CreatedAt:    int64(transactionMap["CreatedAt"].(float64)),
				}
				transactionList = append(transactionList, t)
			}
		}
	}

	transactionListDTO := make([]*web.TransactionDTO, len(transactionList))
	for i, t := range transactionList {
		transactionListDTO[i] = t.ToDTO()
	}

	return &web.PaginationDTO{
		TotalRows:   p.TotalRows,
		Limit:       p.Limit,
		CurrentPage: p.Page,
		TotalPages:  int64(totalPages),
		Rows:        transactionListDTO,
	}
}
