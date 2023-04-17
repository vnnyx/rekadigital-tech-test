package customer

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/redis/go-redis/v9"
	"github.com/vnnyx/rekadigital-tech-test/infrastructure"
	"github.com/vnnyx/rekadigital-tech-test/internal/entity"
)

type CustomerRepositoryImpl struct {
	db    *sql.DB
	redis *redis.Client
}

func NewCustomerRepository(db *sql.DB, redis *redis.Client) CustomerRepository {
	return &CustomerRepositoryImpl{
		db:    db,
		redis: redis,
	}
}

func (r *CustomerRepositoryImpl) StoreCustomer(customer *entity.Customer) error {
	rctx, cancel := infrastructure.NewRedisContext()
	defer cancel()

	r.redis.FlushAll(rctx)

	ctx, cancel := infrastructure.NewMySQLContext()
	defer cancel()

	customer.CreatedAt = time.Now().Unix()

	args := []any{
		customer.ID,
		customer.Name,
		customer.CreatedAt,
	}

	query := "INSERT INTO customer(id, name, created_at) VALUES(?,?,?)"

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *CustomerRepositoryImpl) FindCustomerByName(name string) (*entity.Customer, error) {
	customer := new(entity.Customer)

	rctx, cancel := infrastructure.NewRedisContext()
	defer cancel()

	key := fmt.Sprintf("customer:%v", strings.ToLower(name))

	cache, err := r.redis.Get(rctx, key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cache), &customer)
		if err != nil {
			return nil, err
		}
		return customer, nil
	}

	ctx, cancel := infrastructure.NewMySQLContext()
	defer cancel()

	query := "SELECT id, name, created_at FROM customer WHERE name=? LIMIT 1"

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&customer.ID, &customer.Name, &customer.CreatedAt)
		if err != nil {
			return nil, err
		}
		val, err := json.Marshal(&customer)
		if err != nil {
			return nil, err
		}
		err = r.redis.SetEx(rctx, key, val, 5*time.Minute).Err()
		if err != nil {
			return nil, err
		}
		return customer, nil
	}
	return nil, fmt.Errorf("Customer with Name %v Not Found", name)
}
