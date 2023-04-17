package http

import (
	"net/http"

	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
	"github.com/vnnyx/rekadigital-tech-test/internal/delivery/http/web"
	"github.com/vnnyx/rekadigital-tech-test/internal/helper"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if validationError(err, c) {
		return
	}
	generalError(err, c)
}

func generalError(err error, ctx echo.Context) {
	_ = ctx.JSON(http.StatusInternalServerError, web.WebResponse{
		Code:   http.StatusInternalServerError,
		Status: http.StatusText(http.StatusInternalServerError),
		Data:   nil,
		Error:  err.Error(),
	})
}

func validationError(err error, ctx echo.Context) bool {
	_, ok := err.(helper.ValidationError)
	if ok {
		var obj interface{}
		_ = json.Unmarshal([]byte(err.Error()), &obj)
		_ = ctx.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: http.StatusText(http.StatusBadRequest),
			Data:   nil,
			Error:  obj,
		})
		return true
	}
	return false
}
