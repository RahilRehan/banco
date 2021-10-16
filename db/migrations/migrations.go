package migration

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrations(dbSource, mPath string) error {
	fmt.Println("running migrations...")
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		return fmt.Errorf("migrations: cannot open sql source: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrations: cannot create migrate with instance: %v", err)
	}

	srcPath, err := getSourcePath(mPath)
	if err != nil {
		return fmt.Errorf("cannot create source path: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(srcPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrations: cannot create migrate new instance: %v", err)
	}
	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("no change")
			return nil
		}
		return fmt.Errorf("migrations %v", err)
	}

	fmt.Println("completed running migrations...")

	return nil
}

func getSourcePath(directory string) (string, error) {
	cutSet := "file://"
	directory = strings.TrimLeft(directory, cutSet)

	absPath, err := filepath.Abs(directory)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", cutSet, absPath), nil
}
