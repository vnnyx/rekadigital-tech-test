package bootstrap

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/vnnyx/rekadigital-tech-test/internal/config"
)

func RunMigration(cfg *config.Config) {
	migration, err := migrate.New(cfg.Migration.Source, cfg.Migration.DSN)
	if err != nil {
		log.Fatal(err)
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
