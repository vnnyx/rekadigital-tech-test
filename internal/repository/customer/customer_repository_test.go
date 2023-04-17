package customer_test

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
	"github.com/vnnyx/rekadigital-tech-test/internal/repository/customer"
)

func TestCustomerRepository_StoreCustomer(t *testing.T) {
	var mockDB *sql.DB
	var mockRedis *redis.Client
	var repo customer.CustomerRepository
	var mock sqlmock.Sqlmock

	query := `^INSERT INTO customer\(id, name, created_at\) VALUES\(\?,\?,\?\)$`

	customerData := &entity.Customer{
		ID:        "customer-id",
		Name:      "John Doe",
		CreatedAt: time.Now().Unix(),
	}

	args := []driver.Value{
		"customer-id",
		"John Doe",
		time.Now().Unix(),
	}

	tests := []struct {
		name     string
		customer *entity.Customer
		mockFunc func()
		wantErr  bool
	}{
		{
			name:     "Success store customer",
			customer: customerData,
			mockFunc: func() {
				mock.ExpectPrepare(query).ExpectExec().
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:     "Failed to prepare statement",
			customer: customerData,
			mockFunc: func() {
				mock.ExpectPrepare(query).WillReturnError(errors.New("failed to prepare statement"))
			},
			wantErr: true,
		},
		{
			name:     "Failed to exec statement",
			customer: customerData,
			mockFunc: func() {
				mock.ExpectPrepare(query).ExpectExec().
					WithArgs(args...).
					WillReturnError(errors.New("failed to exec statement"))
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

			repo = customer.NewCustomerRepository(mockDB, mockRedis)

			tt.mockFunc()

			err = repo.StoreCustomer(tt.customer)

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

func TestCustomerRepository_FindCustomerByName(t *testing.T) {
	var mockDB *sql.DB
	var mockRedis *redis.Client
	var repo customer.CustomerRepository
	var mock sqlmock.Sqlmock
	var rmock redismock.ClientMock

	query := `^SELECT id, name, created_at FROM customer WHERE name=\? LIMIT 1$`

	tests := []struct {
		name         string
		customerName string
		mockFunc     func()
		want         *entity.Customer
		wantErr      bool
	}{
		{
			name:         "Success get data from redis",
			customerName: "john",
			mockFunc: func() {
				cutomerData := &entity.Customer{
					ID:        "customer-id",
					Name:      "John Doe",
					CreatedAt: 1234,
				}
				customerByte, err := json.Marshal(cutomerData)
				if err != nil {
					t.Fatal(err)
				}
				rmock.ExpectGet(fmt.Sprintf("customer:%v", strings.ToLower("john"))).
					SetVal(string(customerByte))
			},
			want: &entity.Customer{
				ID:        "customer-id",
				Name:      "John Doe",
				CreatedAt: 1234,
			},
			wantErr: false,
		},
		{
			name:         "Success get data from main database",
			customerName: "john",
			mockFunc: func() {
				key := fmt.Sprintf("customer:%v", strings.ToLower("john"))
				cutomerData := &entity.Customer{
					ID:        "customer-id",
					Name:      "John Doe",
					CreatedAt: 1234,
				}
				customerByte, err := json.Marshal(cutomerData)
				if err != nil {
					t.Fatal(err)
				}
				rmock.ExpectGet(key).
					SetErr(errors.New("error"))

				mock.ExpectPrepare(query)

				rows := sqlmock.NewRows([]string{"id", "name", "create_at"}).
					AddRow(cutomerData.ID, cutomerData.Name, cutomerData.CreatedAt)
				mock.ExpectQuery(query).WithArgs("john").WillReturnRows(rows)

				rmock.ExpectSetEx(key, customerByte, 5*time.Minute).SetVal("OK")
			},
			want: &entity.Customer{
				ID:        "customer-id",
				Name:      "John Doe",
				CreatedAt: 1234,
			},
			wantErr: false,
		},
		{
			name:         "No customer found",
			customerName: "john",
			mockFunc: func() {
				key := fmt.Sprintf("customer:%v", strings.ToLower("john"))
				rmock.ExpectGet(key).
					SetErr(errors.New("error"))

				mock.ExpectPrepare(query)

				rows := sqlmock.NewRows([]string{"id", "name", "create_at"})
				mock.ExpectQuery(query).WithArgs("john").WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "Error unmarshalling JSON data",
			customerName: "john",
			mockFunc: func() {
				rmock.ExpectGet(fmt.Sprintf("customer:%v", strings.ToLower("john"))).
					SetVal("not-a-valid-json-string")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "Error when preparing statement",
			customerName: "john",
			mockFunc: func() {
				mock.ExpectPrepare(query).
					WillReturnError(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name:         "Query returns error",
			customerName: "john",
			mockFunc: func() {
				mock.ExpectPrepare(query).
					ExpectQuery().WithArgs("john").
					WillReturnError(fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name:         "Rows scan returns error",
			customerName: "john",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).
					AddRow("customer-id", "John Doe", "invalid-timestamp")

				mock.ExpectPrepare(query).
					ExpectQuery().WithArgs("john").
					WillReturnRows(rows)
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

			mockRedis, rmock = redismock.NewClientMock()

			tt.mockFunc()

			repo = customer.NewCustomerRepository(mockDB, mockRedis)

			got, err := repo.FindCustomerByName(tt.customerName)
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
