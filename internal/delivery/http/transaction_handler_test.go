package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	h "github.com/vnnyx/rekadigital-tech-test/internal/delivery/http"
	"github.com/vnnyx/rekadigital-tech-test/internal/delivery/http/web"
	"github.com/vnnyx/rekadigital-tech-test/internal/helper"
	mock_usecase "github.com/vnnyx/rekadigital-tech-test/internal/usecase/mocks"
)

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	type mockCreateTransaction struct {
		got *web.TransactionDTO
		err error
	}

	tests := []struct {
		name                  string
		request               string
		input                 *web.TransactionCreateReq
		mockCreateTransaction *mockCreateTransaction
		wantErrBinding        bool
		expectedStatus        int
		expectedResponse      string
		wantErr               bool
	}{
		{
			name: "Success create transaction",
			request: `{
				"name":"John Doe",
				"menu":"Pizza",
				"price": 10,
				"qty": 2,
				"payment":"Cash"
			}`,
			input: &web.TransactionCreateReq{
				CustomerName: "John Doe",
				Menu:         "Pizza",
				Price:        10,
				Qty:          2,
				Payment:      "Cash",
			},
			mockCreateTransaction: &mockCreateTransaction{
				got: &web.TransactionDTO{
					TransactionID: "transaction-id",
					CustomerName:  "John Doe",
					Menu:          "Pizza",
					Price:         10,
					Qty:           2,
					Payment:       "Cash",
					Total:         20,
				},
				err: nil,
			},
			wantErrBinding: false,
			expectedStatus: 201,
			expectedResponse: `{
				"code": 201,
				"status": "Created",
				"data": {
					"id": "transaction-id",
					"name": "John Doe",
					"menu": "Pizza",
					"price": 10,
					"qty": 2,
					"payment": "Cash",
					"total": 20
				},
				"error": null
			}`,
			wantErr: false,
		},
		{
			name: "Error binding request",
			request: `{
				"name":"John Doe",
				"menu":"Pizza",
				"price": 10,z
				"qty": 2,
				"payment":"Cash"
			}`,
			input:                 nil,
			mockCreateTransaction: nil,
			wantErrBinding:        true,
			expectedStatus:        400,
			expectedResponse: `{
				"code": 400,
				"status": "Bad Request",
				"data": null,
				"error": "Error binding json data"
			}`,
			wantErr: true,
		},
		{
			name: "Error create transaction usecase",
			request: `{
				"name":"John Doe",
				"menu":"Pizza",
				"price": 10,
				"qty": 2,
				"payment":"Cash"
			}`,
			input: &web.TransactionCreateReq{
				CustomerName: "John Doe",
				Menu:         "Pizza",
				Price:        10,
				Qty:          2,
				Payment:      "Cash",
			},
			mockCreateTransaction: &mockCreateTransaction{
				got: nil,
				err: errors.New("error create transaction usecase"),
			},
			wantErrBinding: false,
			expectedStatus: 500,
			expectedResponse: `{
				"code": 500,
				"status": "Internal Server Error",
				"data": null,
				"error": "error create transaction usecase"
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTransacntionUC := mock_usecase.NewMockTransactionUC(ctrl)

			if tt.mockCreateTransaction != nil && tt.input != nil {
				mockTransacntionUC.EXPECT().CreateTransaction(gomock.Any(), tt.input).
					Return(tt.mockCreateTransaction.got, tt.mockCreateTransaction.err).
					Times(1)
			}

			handler := h.NewTransactionHandler(mockTransacntionUC)

			e := echo.New()
			req := httptest.NewRequest("POST", "/rekadigital-api", strings.NewReader(tt.request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			handler.CreateTransaction(c)

			if rec.Body.String() != "" {
				assert.JSONEq(t, tt.expectedResponse, rec.Body.String())
			}
		})
	}
}

func TestTransactionHandler_GetAllTransaction(t *testing.T) {
	type mockGetAllTransaction struct {
		got *web.PaginationDTO
		err error
	}

	type request struct {
		limit    string
		page     string
		query    string
		customer string
	}

	tests := []struct {
		name                  string
		request               *request
		input                 *helper.TransactionOptions
		mockGetAllTransaction *mockGetAllTransaction
		expectedStatus        int
		expectedResponse      string
		wantErr               bool
	}{
		{
			name: "Success get all transaction",
			request: &request{
				limit:    "1",
				page:     "1",
				query:    "pizza",
				customer: "john",
			},
			input: &helper.TransactionOptions{
				Limit:        1,
				Page:         1,
				Query:        "pizza",
				CustomerName: "john",
			},
			mockGetAllTransaction: &mockGetAllTransaction{
				got: &web.PaginationDTO{
					TotalRows:   1,
					Limit:       1,
					CurrentPage: 1,
					TotalPages:  1,
					Rows: []*web.TransactionDTO{
						{
							TransactionID: "transaction-id",
							CustomerName:  "John Doe",
							Menu:          "Pizza",
							Price:         10,
							Qty:           2,
							Payment:       "Cash",
							Total:         20,
						},
					},
				},
				err: nil,
			},
			expectedStatus: 200,
			expectedResponse: `{
				"code": 200,
				"status": "OK",
				"data": {
					"total_rows": 1,
					"limit": 1,
					"current_page": 1,
					"total_pages": 1,
				"rows": [
						{
							"id": "transaction-id",
							"name": "John Doe",
							"menu": "Pizza",
							"price": 10,
							"qty": 2,
							"payment": "Cash",
							"total": 20
						}
					]
				},
				"error": null
			}`,
			wantErr: false,
		},
		{
			name: "Failed to get all transaction",
			request: &request{
				limit:    "1",
				page:     "1",
				query:    "pizza",
				customer: "john",
			},
			input: &helper.TransactionOptions{
				Limit:        1,
				Page:         1,
				Query:        "pizza",
				CustomerName: "john",
			},
			mockGetAllTransaction: &mockGetAllTransaction{
				got: nil,
				err: errors.New("error get all transaction usecase"),
			},
			expectedStatus: 500,
			expectedResponse: `{
				"code": 500,
				"status": "Internal Server Error",
				"data": null,
				"error": "error get all transaction usecase"
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTransacntionUC := mock_usecase.NewMockTransactionUC(ctrl)

			if tt.mockGetAllTransaction != nil {
				mockTransacntionUC.EXPECT().GetAllTransaction(gomock.Any(), gomock.Any()).
					Return(tt.mockGetAllTransaction.got, tt.mockGetAllTransaction.err).
					Times(1)
			}

			handler := h.NewTransactionHandler(mockTransacntionUC)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/rekadigital-api", nil)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetParamNames("query", "limit", "page", "customer")
			c.SetParamValues(tt.request.query, tt.request.limit, tt.request.page, tt.request.customer)

			err := handler.GetAllTransaction(c)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if rec.Body.String() != "" {
				assert.JSONEq(t, tt.expectedResponse, rec.Body.String())
			}

		})
	}
}
