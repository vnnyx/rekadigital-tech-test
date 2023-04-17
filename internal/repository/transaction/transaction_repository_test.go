package transaction_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/vnnyx/rekadigital-tech-test/internal/entity"
	"github.com/vnnyx/rekadigital-tech-test/internal/helper"
	"github.com/vnnyx/rekadigital-tech-test/internal/repository/transaction"
)

func TestTransactionRepository_StoreTransaction(t *testing.T) {
	var mockDB *sql.DB
	var mockRedis *redis.Client
	var repo transaction.TransactionRepository
	var mock sqlmock.Sqlmock

	query := `^INSERT INTO transaction\(id, customer_id, menu, price, qty, payment, total, created_at\) VALUES\(\?,\?,\?,\?,\?,\?,\?,\?\)$`
	transactionData := &entity.Transaction{
		ID:         "transaction-id",
		CustomerID: "customer-id",
		Menu:       "Pizza",
		Price:      20,
		Qty:        2,
		Payment:    "Cash",
		Total:      40,
	}
	args := []driver.Value{
		"transaction-id",
		"customer-id",
		"Pizza",
		20,
		2,
		"Cash",
		40,
		time.Now().Unix(),
	}

	tests := []struct {
		name        string
		transaction *entity.Transaction
		mockFunc    func()
		wantErr     bool
	}{
		{
			name:        "Success store transaction",
			transaction: transactionData,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectPrepare(query)
				mock.ExpectExec(query).WithArgs(args...).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:        "Failed to begin tx",
			transaction: transactionData,
			mockFunc: func() {
				mock.ExpectBegin().WillReturnError(errors.New("failed to begin tx"))
			},
			wantErr: true,
		},
		{
			name:        "Failed to prepare statement",
			transaction: transactionData,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectPrepare(query).WillReturnError(errors.New("failed to prepare statement"))
			},
			wantErr: true,
		},
		{
			name:        "Failed to exec query",
			transaction: transactionData,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectPrepare(query).ExpectExec().
					WithArgs(args...).
					WillReturnError(errors.New("failed to exec query"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			mockDB, mock, err = sqlmock.New()
			if err != nil {
				t.Fatalf("error creating mock database: %s", err)
			}
			defer mockDB.Close()

			mockRedis, _ = redismock.NewClientMock()

			repo = transaction.NewTransactionRepository(mockDB, mockRedis)

			tt.mockFunc()

			err = repo.StoreTransaction(tt.transaction)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error result: got %v, want %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Failed to meet expectations: %v", err)
			}
		})
	}
	// Cleanup after all tests are done
	t.Cleanup(func() {
		if err := mockDB.Close(); err != nil {
			t.Errorf("Failed to close mock database: %v", err)
		}
	})
}

func TestTransactionRepository_GetAllTransaction(t *testing.T) {
	var mockDB *sql.DB
	var mockRedis *redis.Client
	var repo transaction.TransactionRepository
	var mock sqlmock.Sqlmock
	var rmock redismock.ClientMock

	tests := []struct {
		name     string
		opt      *helper.TransactionOptions
		mockFunc func()
		want     *helper.Pagination
		wantErr  bool
	}{
		{
			name: "Success get all transaction from redis",
			opt: &helper.TransactionOptions{
				Limit:        -1,
				Page:         -1,
				Query:        "pizza",
				CustomerName: "john",
			},
			mockFunc: func() {
				pagination := &helper.Pagination{
					TotalRows: 1,
					Limit:     10,
					Page:      1,
					Rows: []*entity.Transaction{
						{
							ID:           "transaction-id",
							CustomerName: "John Doe",
							Menu:         "Pizza",
							Price:        10,
							Qty:          2,
							Payment:      "Cash",
							Total:        20,
							CreatedAt:    123,
						},
					},
				}
				paginationByte, err := json.Marshal(&pagination)
				if err != nil {
					t.Fatal(err)
				}
				rmock.ExpectGet(fmt.Sprintf("tx:%v:%v:%v:%v", 10, 1, strings.ToLower("pizza"), strings.ToLower("john"))).
					SetVal(string(paginationByte))
			},
			want: &helper.Pagination{
				TotalRows: 1,
				Limit:     10,
				Page:      1,
				Rows: []any{
					map[string]any{
						"ID":           "transaction-id",
						"CustomerName": "John Doe",
						"CustomerID":   "",
						"Menu":         "Pizza",
						"Price":        float64(10),
						"Qty":          float64(2),
						"Payment":      "Cash",
						"Total":        float64(20),
						"CreatedAt":    float64(123),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success get all transaction from main database",
			opt: &helper.TransactionOptions{
				Limit:        10,
				Page:         1,
				Query:        "pizza",
				CustomerName: "john",
			},
			mockFunc: func() {
				key := fmt.Sprintf("tx:%v:%v:%v:%v", 10, 1, strings.ToLower("pizza"), strings.ToLower("john"))
				query := `^SELECT SQL_CALC_FOUND_ROWS transaction.id, c.name, transaction.menu, transaction.price, transaction.qty, transaction.payment, transaction.total, transaction.created_at
        				FROM transaction JOIN customer c on c.id = transaction.customer_id WHERE transaction.menu LIKE \? AND c.Name LIKE \? ORDER BY transaction.created_at DESC LIMIT \?, \?$`
				pagination := &helper.Pagination{
					TotalRows: 1,
					Limit:     10,
					Page:      1,
					Rows: []*entity.Transaction{
						{
							ID:           "transaction-id",
							CustomerName: "John Doe",
							Menu:         "Pizza",
							Price:        10,
							Qty:          2,
							Payment:      "Cash",
							Total:        20,
							CreatedAt:    123,
						},
					},
				}
				args := []driver.Value{
					"%" + "pizza" + "%",
					"%" + "john" + "%",
					0,
					10,
				}
				paginationByte, err := json.Marshal(&pagination)
				if err != nil {
					t.Fatal(err)
				}
				rmock.ExpectGet(key).
					SetErr(errors.New("error"))

				mock.ExpectBegin()
				mock.ExpectPrepare(query)

				rows := sqlmock.NewRows([]string{
					"id",
					"name",
					"menu",
					"price",
					"qty",
					"payment",
					"total",
					"created_at",
				}).AddRow(
					"transaction-id",
					"John Doe",
					"Pizza",
					10,
					2,
					"Cash",
					20,
					123,
				)
				foundRows := sqlmock.NewRows([]string{
					"FOUND_ROWS()",
				}).AddRow(1)
				mock.ExpectQuery(query).WithArgs(args...).WillReturnRows(rows)
				mock.ExpectQuery(`^SELECT FOUND_ROWS\(\)$`).WillReturnRows(foundRows)
				mock.ExpectCommit()

				rmock.ExpectSetEx(key, paginationByte, 5*time.Minute).SetVal("OK")
			},
			want: &helper.Pagination{
				TotalRows: 1,
				Limit:     10,
				Page:      1,
				Rows: []*entity.Transaction{
					{
						ID:           "transaction-id",
						CustomerName: "John Doe",
						CustomerID:   "",
						Menu:         "Pizza",
						Price:        10,
						Qty:          2,
						Payment:      "Cash",
						Total:        20,
						CreatedAt:    123,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error begin tx",
			opt: &helper.TransactionOptions{
				Limit:        10,
				Page:         1,
				Query:        "pizza",
				CustomerName: "john",
			},
			mockFunc: func() {
				key := fmt.Sprintf("tx:%v:%v:%v:%v", 10, 1, strings.ToLower("pizza"), strings.ToLower("john"))
				rmock.ExpectGet(key).
					SetErr(errors.New("error"))
				mock.ExpectBegin().WillReturnError(errors.New("error begin tx"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error prepare statement",
			opt: &helper.TransactionOptions{
				Limit:        10,
				Page:         1,
				Query:        "",
				CustomerName: "",
			},
			mockFunc: func() {
				query := `^SELECT SQL_CALC_FOUND_ROWS transaction.id, c.name, transaction.menu, transaction.price, transaction.qty, transaction.payment, transaction.total, transaction.created_at
        			FROM transaction JOIN customer c on c.id = transaction.customer_id ORDER BY transaction.created_at DESC LIMIT \?, \?$`
				key := fmt.Sprintf("tx:%v:%v:%v:%v", 10, 1, strings.ToLower(""), strings.ToLower(""))
				rmock.ExpectGet(key).
					SetErr(errors.New("error"))
				mock.ExpectBegin()
				mock.ExpectPrepare(query).WillReturnError(errors.New("error prepare statement"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error execute statement",
			opt: &helper.TransactionOptions{
				Limit:        10,
				Page:         1,
				Query:        "",
				CustomerName: "",
			},
			mockFunc: func() {
				query := `^SELECT SQL_CALC_FOUND_ROWS transaction.id, c.name, transaction.menu, transaction.price, transaction.qty, transaction.payment, transaction.total, transaction.created_at
        			FROM transaction JOIN customer c on c.id = transaction.customer_id ORDER BY transaction.created_at DESC LIMIT \?, \?$`
				key := fmt.Sprintf("tx:%v:%v:%v:%v", 10, 1, strings.ToLower(""), strings.ToLower(""))
				rmock.ExpectGet(key).
					SetErr(errors.New("error"))
				mock.ExpectBegin()
				mock.ExpectPrepare(query)
				mock.ExpectQuery(query).WillReturnError(errors.New("error execute query"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error execute second statement",
			opt: &helper.TransactionOptions{
				Limit:        10,
				Page:         1,
				Query:        "2",
				CustomerName: "",
			},
			mockFunc: func() {
				query := `^SELECT SQL_CALC_FOUND_ROWS transaction.id, c.name, transaction.menu, transaction.price, transaction.qty, transaction.payment, transaction.total, transaction.created_at
        FROM transaction JOIN customer c on c.id = transaction.customer_id WHERE transaction.price=\? ORDER BY transaction.created_at DESC LIMIT \?, \?$`
				key := fmt.Sprintf("tx:%v:%v:%v:%v", 10, 1, strings.ToLower("2"), strings.ToLower(""))
				rmock.ExpectGet(key).
					SetErr(errors.New("error"))
				mock.ExpectBegin()
				mock.ExpectPrepare(query)
				mock.ExpectQuery(query).WillReturnRows(mock.NewRows([]string{
					"id",
					"name",
					"menu",
					"price",
					"qty",
					"payment",
					"total",
					"created_at",
				}))
				mock.ExpectQuery(`^SELECT FOUND_ROWS\(\)$`).WillReturnError(errors.New("error execute second statement"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error scan rows",
			opt: &helper.TransactionOptions{
				Limit:        10,
				Page:         1,
				Query:        "",
				CustomerName: "john",
			},
			mockFunc: func() {
				query := `^SELECT SQL_CALC_FOUND_ROWS transaction.id, c.name, transaction.menu, transaction.price, transaction.qty, transaction.payment, transaction.total, transaction.created_at
        FROM transaction JOIN customer c on c.id = transaction.customer_id WHERE c.Name LIKE \? ORDER BY transaction.created_at DESC LIMIT \?, \?$`
				key := fmt.Sprintf("tx:%v:%v:%v:%v", 10, 1, strings.ToLower(""), strings.ToLower("john"))
				rmock.ExpectGet(key).
					SetErr(errors.New("error"))
				mock.ExpectBegin()
				mock.ExpectPrepare(query)
				rows := sqlmock.NewRows([]string{
					"id",
					"name",
					"menu",
					"price",
					"qty",
					"payment",
					"total",
					"created_at",
				}).AddRow(
					"transaction-id",
					"John Doe",
					"Pizza",
					10,
					2,
					"Cash",
					20,
					123,
				).RowError(0, errors.New("error scan rows"))
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			mockDB, mock, err = sqlmock.New()
			if err != nil {
				t.Fatalf("error creating mock database: %s", err)
			}
			defer mockDB.Close()

			mockRedis, rmock = redismock.NewClientMock()

			tt.mockFunc()

			repo = transaction.NewTransactionRepository(mockDB, mockRedis)

			got, err := repo.GetAllTransaction(tt.opt)

			if (err != nil) != tt.wantErr {
				t.Errorf("repository.FindCustomerByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repository.FindCustomerByName() = %v, want %v", got, tt.want)
			}

		})
	}
	// Cleanup after all tests are done
	t.Cleanup(func() {
		if err := mockDB.Close(); err != nil {
			t.Errorf("Failed to close mock database: %v", err)
		}
	})
}
