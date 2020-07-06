package database

import (
	"database/sql"
	"fmt"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/router"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/markbates/pkger"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
)

const component = "database"

func GetConn(databasePath string) *sql.DB {

	conn, err := sql.Open("sqlite3", databasePath)

	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	return conn
}

func RunMigrations() error {
	authDatabasepath := router.GetDatabase("auth")

	migrationsAuth := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/pkg/database/migrations/auth"),
	}

	err := doMigrate(migrationsAuth, authDatabasepath)
	if err != nil {
		return err
	}
	return nil
}

func doMigrate(migrations *migrate.HttpFileSystemMigrationSource, databasePath string) error {
	conn := GetConn(databasePath)

	n, err := migrate.Exec(conn, "sqlite3", migrations, migrate.Up)

	if err != nil {
		return err
	}

	err = conn.Close()
	if err != nil {
		return err
	}

	helpers.LogInfo(fmt.Sprintf("Applied %d migrations in %s", n, databasePath), component)
	return nil
}
