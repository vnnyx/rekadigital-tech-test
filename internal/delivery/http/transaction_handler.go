package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/vnnyx/rekadigital-tech-test/internal/delivery/http/web"
	"github.com/vnnyx/rekadigital-tech-test/internal/helper"
	"github.com/vnnyx/rekadigital-tech-test/internal/usecase"
)

type TransactionHandler struct {
	transactionUC usecase.TransactionUC
}

func NewTransactionHandler(transactionUC usecase.TransactionUC) *TransactionHandler {
	return &TransactionHandler{
		transactionUC: transactionUC,
	}
}

func (t *TransactionHandler) CreateTransaction(ctx echo.Context) error {
	var req web.TransactionCreateReq
	err := ctx.Bind(&req)
	if err != nil {
		return err
	}

	res, err := t.transactionUC.CreateTransaction(ctx.Request().Context(), &req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, web.WebResponse{
		Code:   http.StatusCreated,
		Status: "Created",
		Data:   res,
	})
}

func (t *TransactionHandler) GetAllTransaction(ctx echo.Context) error {
	var opt helper.TransactionOptions
	opt.Query = ctx.QueryParam("query")
	opt.CustomerName = ctx.QueryParam("customer")
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	page, _ := strconv.Atoi(ctx.QueryParam("page"))

	opt.Limit = int64(limit)
	opt.Page = int64(page)

	res, err := t.transactionUC.GetAllTransaction(ctx.Request().Context(), &opt)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   res,
	})
}
