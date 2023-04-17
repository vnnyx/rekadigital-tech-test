package usecase

import (
	"context"

	"github.com/vnnyx/rekadigital-tech-test/internal/delivery/http/web"
	"github.com/vnnyx/rekadigital-tech-test/internal/helper"
)

type TransactionUC interface {
	CreateTransaction(ctx context.Context, req *web.TransactionCreateReq) (*web.TransactionDTO, error)
	GetAllTransaction(ctx context.Context, opt *helper.TransactionOptions) (*web.PaginationDTO, error)
}
