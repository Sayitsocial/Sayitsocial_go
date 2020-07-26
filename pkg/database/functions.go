package database

import (
	"database/sql"
	"fmt"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	_ "github.com/lib/pq"
	"github.com/markbates/pkger"
	"github.com/rubenv/sql-migrate"
)

const component = "database"

func GetConn() *sql.DB {

	conn, err := sql.Open("postgres", helpers.PgConnString)

	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	return conn
}

func RunMigrations() error {
	migrationsAuth := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/pkg/database/migrations/"),
	}

	err := doMigrate(migrationsAuth)
	if err != nil {
		return err
	}
	return nil
}

func doMigrate(migrations *migrate.HttpFileSystemMigrationSource) error {
	conn := GetConn()

	n, err := migrate.Exec(conn, "postgres", migrations, migrate.Up)

	if err != nil {
		return err
	}

	err = conn.Close()
	if err != nil {
		return err
	}

	helpers.LogInfo(fmt.Sprintf("Applied %d migrations", n), component)
	return nil
}
