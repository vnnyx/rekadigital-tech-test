package transaction

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/redis/go-redis/v9"
	"github.com/vnnyx/rekadigital-tech-test/infrastructure"
	"github.com/vnnyx/rekadigital-tech-test/internal/entity"
	"github.com/vnnyx/rekadigital-tech-test/internal/helper"
)

type TransactionRepositoryImpl struct {
	db    *sql.DB
	redis *redis.Client
}

func NewTransactionRepository(db *sql.DB, redis *redis.Client) TransactionRepository {
	return &TransactionRepositoryImpl{
		db:    db,
		redis: redis,
	}
}

func (r *TransactionRepositoryImpl) StoreTransaction(transaction *entity.Transaction) error {
	rctx, cancel := infrastructure.NewRedisContext()
	defer cancel()

	r.redis.FlushAll(rctx)

	ctx, cancel := infrastructure.NewMySQLContext()
	defer cancel()

	transaction.CreatedAt = time.Now().Unix()

	args := []any{
		transaction.ID,
		transaction.CustomerID,
		transaction.Menu,
		transaction.Price,
		transaction.Qty,
		transaction.Payment,
		transaction.Total,
		transaction.CreatedAt,
	}

	query := "INSERT INTO transaction(id, customer_id, menu, price, qty, payment, total, created_at) VALUES(?,?,?,?,?,?,?,?)"

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return err
}

func (r *TransactionRepositoryImpl) GetAllTransaction(opt *helper.TransactionOptions) (*helper.Pagination, error) {
	pagination := new(helper.Pagination)

	rctx, cancel := infrastructure.NewRedisContext()
	defer cancel()

	if opt.Limit <= 0 {
		opt.Limit = 10
	}
	if opt.Page <= 0 {
		opt.Page = 1
	}

	key := fmt.Sprintf("tx:%v:%v:%v:%v", opt.Limit, opt.Page, strings.ToLower(opt.Query), strings.ToLower(opt.CustomerName))

	cache, err := r.redis.Get(rctx, key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cache), &pagination)
		if err != nil {
			return nil, err
		}
		return pagination, nil
	}

	ctx, cancel := infrastructure.NewMySQLContext()
	defer cancel()

	offset := (opt.Page - 1) * opt.Limit
	args := []interface{}{}

	query := `SELECT SQL_CALC_FOUND_ROWS transaction.id, c.name, transaction.menu, transaction.price, transaction.qty, transaction.payment, transaction.total, transaction.created_at
        FROM transaction JOIN customer c on c.id = transaction.customer_id`

	if opt.Query != "" {
		num, err := strconv.Atoi(opt.Query)
		if err == nil {
			query += ` WHERE transaction.price=?`
			args = append(args, num)
		} else {
			query += ` WHERE transaction.menu LIKE ?`
			args = append(args, "%"+opt.Query+"%")
		}
		if opt.CustomerName != "" {
			query += ` AND c.Name LIKE ?`
			args = append(args, "%"+opt.CustomerName+"%")
		}
	} else {
		if opt.CustomerName != "" {
			query += ` WHERE c.Name LIKE ?`
			args = append(args, "%"+opt.CustomerName+"%")
		}
	}

	query += ` ORDER BY transaction.created_at DESC LIMIT ?, ?`
	args = append(args, offset, opt.Limit)

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listTransaction = make([]*entity.Transaction, 0)
	for rows.Next() {
		var t = new(entity.Transaction)
		err := rows.Scan(&t.ID, &t.CustomerName, &t.Menu, &t.Price, &t.Qty, &t.Payment, &t.Total, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		listTransaction = append(listTransaction, t)
	}

	var totalRows int64
	err = tx.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&totalRows)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	pagination = &helper.Pagination{
		TotalRows: totalRows,
		Limit:     opt.Limit,
		Page:      opt.Page,
		Rows:      listTransaction,
	}

	val, err := json.Marshal(&pagination)
	if err != nil {
		return nil, err
	}

	err = r.redis.SetEx(rctx, key, val, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return pagination, nil
}
