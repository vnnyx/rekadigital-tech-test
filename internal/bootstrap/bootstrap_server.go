package bootstrap

import (
	"log"
	"os"

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

func StartServer() {
	env := os.Getenv("ENV")

	cfg, err := config.New(env)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisablePrintStack: true}))
	e.Use(middleware.CORS())
	e.HTTPErrorHandler = http.CustomHTTPErrorHandler

	RunMigration(cfg)
	mysqlDB := infrastructure.NewMySQLDatabase(cfg)
	redisClient := infrastructure.NewRedisClient(cfg)

	customerRepo := customer.NewCustomerRepository(mysqlDB, redisClient)
	transactionRepo := transaction.NewTransactionRepository(mysqlDB, redisClient)
	transactionUC := usecase.NewTransactionUC(transactionRepo, customerRepo)

	transactionHaandler := http.NewTransactionHandler(transactionUC)

	r := routes.NewRoute(transactionHaandler, e)
	r.InitRoute()

	e.Logger.Fatal(e.Start(":3000"))
}
