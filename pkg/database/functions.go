package database

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func GetConn() *sql.DB {

	conn, err := sql.Open("postgres", helpers.PgConnString)

	if err != nil {
		helpers.LogError(err.Error())
	}

	return conn
}

func RunMigrations() error {
	migrationsAuth := &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.Dir("./pkg/database/migrations/"),
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

	helpers.LogInfo(fmt.Sprintf("Applied %d migrations", n))
	return nil
}
