package usecase_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/vnnyx/rekadigital-tech-test/internal/delivery/http/web"
	"github.com/vnnyx/rekadigital-tech-test/internal/entity"
	"github.com/vnnyx/rekadigital-tech-test/internal/helper"
	cmock "github.com/vnnyx/rekadigital-tech-test/internal/repository/customer/mocks"
	tmock "github.com/vnnyx/rekadigital-tech-test/internal/repository/transaction/mocks"
	"github.com/vnnyx/rekadigital-tech-test/internal/usecase"
)

func TestTransactionUsecase_CreateTransaction(t *testing.T) {
	type args struct {
		ctx context.Context
		req *web.TransactionCreateReq
	}
	type mockFindCustomerByName struct {
		got *entity.Customer
		err error
	}
	tests := []struct {
		name                   string
		args                   *args
		mockFindCustomerByName *mockFindCustomerByName
		errorStoreCustomer     error
		errorStoreTransaction  error
		wantErr                bool
		want                   *web.TransactionDTO
	}{
		{
			name: "Create new transaction successfully",
			args: &args{
				ctx: context.TODO(),
				req: &web.TransactionCreateReq{
					CustomerName: "John Doe",
					Menu:         "Pizza",
					Price:        10,
					Qty:          2,
					Payment:      "Cash",
				},
			},
			mockFindCustomerByName: &mockFindCustomerByName{
				got: &entity.Customer{
					ID:   "customer-id",
					Name: "John Doe",
				},
				err: nil,
			},
			errorStoreCustomer:    nil,
			errorStoreTransaction: nil,
			wantErr:               false,
			want: &web.TransactionDTO{
				TransactionID: "transaction-id",
				CustomerName:  "John Doe",
				Menu:          "Pizza",
				Price:         10,
				Qty:           2,
				Payment:       "Cash",
				Total:         20,
			},
		},
		{
			name: "Create new transaction successfully with new customer",
			args: &args{
				ctx: context.TODO(),
				req: &web.TransactionCreateReq{
					CustomerName: "John Doe",
					Menu:         "Pizza",
					Price:        10,
					Qty:          2,
					Payment:      "Cash",
				},
			},
			mockFindCustomerByName: &mockFindCustomerByName{
				got: nil,
				err: errors.New("customer not found"),
			},
			errorStoreCustomer:    nil,
			errorStoreTransaction: nil,
			wantErr:               false,
			want: &web.TransactionDTO{
				TransactionID: "transaction-id",
				CustomerName:  "John Doe",
				Menu:          "Pizza",
				Price:         10,
				Qty:           2,
				Payment:       "Cash",
				Total:         20,
			},
		},
		{
			name: "Error when storing transaction",
			args: &args{
				ctx: context.TODO(),
				req: &web.TransactionCreateReq{
					CustomerName: "John Doe",
					Menu:         "Pizza",
					Price:        10,
					Qty:          2,
					Payment:      "Cash",
				},
			},
			mockFindCustomerByName: &mockFindCustomerByName{
				got: &entity.Customer{
					ID:   "customer-id",
					Name: "John Doe",
				},
				err: nil,
			},
			errorStoreCustomer:    nil,
			errorStoreTransaction: errors.New("failed to store transaction"),
			wantErr:               true,
			want:                  nil,
		},
		{
			name: "Error when storing customer",
			args: &args{
				ctx: context.TODO(),
				req: &web.TransactionCreateReq{
					CustomerName: "John Doe",
					Menu:         "Pizza",
					Price:        10,
					Qty:          2,
					Payment:      "Cash",
				},
			},
			mockFindCustomerByName: &mockFindCustomerByName{
				got: nil,
				err: errors.New("customer not found"),
			},
			errorStoreCustomer:    errors.New("failed to store customer"),
			errorStoreTransaction: nil,
			wantErr:               true,
			want:                  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTransactionRepo := tmock.NewMockTransactionRepository(ctrl)
			mockCustomerRepo := cmock.NewMockCustomerRepository(ctrl)

			uc := usecase.NewTransactionUC(mockTransactionRepo, mockCustomerRepo)

			if tt.mockFindCustomerByName != nil {
				mockCustomerRepo.EXPECT().FindCustomerByName(tt.args.req.CustomerName).Return(tt.mockFindCustomerByName.got, tt.mockFindCustomerByName.err)
			}
			if tt.mockFindCustomerByName.err != nil {
				mockCustomerRepo.EXPECT().StoreCustomer(gomock.Any()).Do(func(customer *entity.Customer) {
					customer.ID = "customer-id"
				}).Return(tt.errorStoreCustomer)
			}
			if tt.errorStoreCustomer == nil {
				mockTransactionRepo.EXPECT().StoreTransaction(gomock.Any()).Do(func(transaction *entity.Transaction) {
					transaction.CustomerID = "customer-id"
					transaction.ID = "transaction-id"
				}).Return(tt.errorStoreTransaction)
			}
			got, err := uc.CreateTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.CreateTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("usecase.CreateTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionUsecase_GetAllTransaction(t *testing.T) {
	type args struct {
		ctx context.Context
		opt *helper.TransactionOptions
	}
	type mockGetAllTransaction struct {
		got *helper.Pagination
		err error
	}
	tests := []struct {
		name                  string
		args                  *args
		mockGetAllTransaction *mockGetAllTransaction
		wantErr               bool
		want                  *web.PaginationDTO
	}{
		{
			name: "Get all transaction successfully",
			args: &args{
				ctx: context.TODO(),
				opt: &helper.TransactionOptions{
					Limit:        10,
					Page:         1,
					Query:        "Pizza",
					CustomerName: "John Doe",
				},
			},
			mockGetAllTransaction: &mockGetAllTransaction{
				got: &helper.Pagination{
					TotalRows: 1,
					Limit:     int64(10),
					Page:      1,
					Rows: []*entity.Transaction{
						{
							ID:           "transaction-id",
							CustomerID:   "customer-id",
							CustomerName: "John Doe",
							Menu:         "Pizza",
							Price:        10,
							Qty:          2,
							Payment:      "Cash",
							Total:        20,
							CreatedAt:    time.Now().Unix(),
						},
					},
				},
				err: nil,
			},
			wantErr: false,
			want: &web.PaginationDTO{
				TotalRows:   int64(1),
				Limit:       10,
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
		},
		{
			name: "Get all transaction successfully",
			args: &args{
				ctx: context.TODO(),
				opt: &helper.TransactionOptions{
					Limit:        10,
					Page:         1,
					Query:        "Pizza",
					CustomerName: "John Doe",
				},
			},
			mockGetAllTransaction: &mockGetAllTransaction{
				got: nil,
				err: errors.New("error when get all transaction from DB"),
			},
			wantErr: true,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTransactionRepo := tmock.NewMockTransactionRepository(ctrl)
			mockCustomerRepo := cmock.NewMockCustomerRepository(ctrl)

			uc := usecase.NewTransactionUC(mockTransactionRepo, mockCustomerRepo)

			if tt.mockGetAllTransaction != nil {
				mockTransactionRepo.EXPECT().GetAllTransaction(tt.args.opt).Return(tt.mockGetAllTransaction.got, tt.mockGetAllTransaction.err)
			}

			got, err := uc.GetAllTransaction(tt.args.ctx, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.GetAllTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// check with ignore pointer address
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("usecase.GetAllTransaction() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
