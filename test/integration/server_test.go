package integration

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vnnyx/rekadigital-tech-test/infrastructure"
	"github.com/vnnyx/rekadigital-tech-test/internal/config"
	"github.com/vnnyx/rekadigital-tech-test/internal/delivery/http"
	"github.com/vnnyx/rekadigital-tech-test/internal/delivery/http/routes"
	"github.com/vnnyx/rekadigital-tech-test/internal/repository/customer"
	"github.com/vnnyx/rekadigital-tech-test/internal/repository/transaction"
	"github.com/vnnyx/rekadigital-tech-test/internal/usecase"
)

var (
	mysql *sql.DB
)

func newTestServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisablePrintStack: true}))
	e.Use(middleware.CORS())
	e.HTTPErrorHandler = http.CustomHTTPErrorHandler

	cfg, err := config.New("test")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	mysql = infrastructure.NewMySQLDatabase(cfg)
	redisClient := infrastructure.NewRedisClient(cfg)

	customerRepo := customer.NewCustomerRepository(mysql, redisClient)
	transactionRepo := transaction.NewTransactionRepository(mysql, redisClient)
	transactionUC := usecase.NewTransactionUC(transactionRepo, customerRepo)

	transactionHaandler := http.NewTransactionHandler(transactionUC)

	r := routes.NewRoute(transactionHaandler, e)
	r.InitRoute()

	return e
}
