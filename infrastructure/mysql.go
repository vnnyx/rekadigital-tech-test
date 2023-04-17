package infrastructure

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/vnnyx/rekadigital-tech-test/internal/config"
)

func NewMySQLDatabase(cfg *config.Config) *sql.DB {
	ctx, cancel := NewMySQLContext()
	defer cancel()

	sqlDB, err := sql.Open("mysql", cfg.Mysql.DSN)
	if err != nil {
		logrus.Fatal(err)
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cfg.Mysql.IdleMax)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(cfg.Mysql.PoolMax)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Mysql.MaxLifeTimeMinute) * time.Minute)

	//sqlDB.SetConnMaxIdleTime(time.Duration(mysqlMaxIdleTime) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Mysql.MaxIdleTimeMinute) * time.Minute)

	return sqlDB
}

func NewMySQLContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
