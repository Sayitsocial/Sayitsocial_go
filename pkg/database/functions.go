package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

// GetConn returns Conn to database
func GetConn() *sql.DB {
	for {
		conn, err := sql.Open("postgres", helpers.PgConnString)
		if err != nil {
			helpers.LogError(err.Error())
			time.Sleep(30 * time.Second)
		}
		return conn
	}
}

// RunMigrations runs all provided migrations
func RunMigrations() error {
	migrationsAuth := &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.Dir(filepath.Join(helpers.GetExecutableDirectory(), "/pkg/database/migrations/")),
	}

	err := doMigrate(migrationsAuth)
	if err != nil {
		return err
	}
	return nil
}

func doMigrate(migrations *migrate.HttpFileSystemMigrationSource) error {
	conn := GetConn()

	// _, err := migrate.Exec(conn, "postgres", migrations, migrate.Down)

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
