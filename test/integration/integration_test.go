package integration

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	app := newTestServer()
	mysql.Exec("DELETE FROM transaction")
	tests := []struct {
		name             string
		request          string
		expectedStatus   int
		expectedResponse string
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
			expectedStatus: 201,
			expectedResponse: `{
				"code": 201,
				"status": "Created",
				"data": {
					"name": "John Doe",
					"menu": "Pizza",
					"price": 10,
					"qty": 2,
					"payment": "Cash",
					"total": 20
				},
				"error": null
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest("POST", "/rekadigital-api/transaction", strings.NewReader(tt.request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)

			var resBody map[string]any
			var resBodyByte []byte

			err := json.Unmarshal(rec.Body.Bytes(), &resBody)
			if err != nil {
				t.Error(err)
			}

			if resBody["data"] != nil {
				delete(resBody["data"].(map[string]any), "id")

				resBodyByte, err = json.Marshal(resBody)
				if err != nil {
					t.Error(err)
				}
			}

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if rec.Body.String() != "" {
				assert.JSONEq(t, tt.expectedResponse, string(resBodyByte))
			}
		})
	}
}

func TestGettAllTransaction(t *testing.T) {
	app := newTestServer()
	type request struct {
		query    string
		customer string
		limit    string
		page     string
	}
	tests := []struct {
		name             string
		request          *request
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "Success get all transaction",
			request: &request{
				query:    "pizza",
				customer: "john",
				limit:    "1",
				page:     "1",
			},
			expectedStatus: 200,
			expectedResponse: `{
				"code": 200,
				"status": "OK",
				"data": {
					"total_rows": 1,
					"limit": 1,
					"current_page": 1,
					"total_pages": 1
				},
				"error": null
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/rekadigital-api/transaction?limit=%v&query=%v&customer=%v&page=%v", tt.request.limit, tt.request.query, tt.request.customer, tt.request.page), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)

			var resBody map[string]any
			var resBodyByte []byte

			err := json.Unmarshal(rec.Body.Bytes(), &resBody)
			if err != nil {
				t.Error(err)
			}

			if resBody["data"] != nil {
				delete(resBody["data"].(map[string]any), "rows")

				resBodyByte, err = json.Marshal(resBody)
				if err != nil {
					t.Error(err)
				}
			}

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if rec.Body.String() != "" {
				assert.JSONEq(t, tt.expectedResponse, string(resBodyByte))
			}
		})
	}
}
