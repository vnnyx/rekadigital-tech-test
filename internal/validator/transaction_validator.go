package validator

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/vnnyx/rekadigital-tech-test/internal/delivery/http/web"
	"github.com/vnnyx/rekadigital-tech-test/internal/helper"
)

func CreateTransactionValidation(req web.TransactionCreateReq) {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.CustomerName, validation.Required),
		validation.Field(&req.Menu, validation.Required),
		validation.Field(&req.Payment, validation.Required),
		validation.Field(&req.Price, validation.Required),
		validation.Field(&req.Qty, validation.Required),
	)

	if err != nil {
		b, _ := json.Marshal(err)
		err = helper.ValidationError{
			Message: string(b),
		}
		if err != nil {
			panic(err)
		}
	}
}
