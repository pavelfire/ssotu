package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	// Библиотека для миграций
	"github.com/golang-migrate/migrate/v4"
	// Драйвер для выполнения миграций sqlite 3
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	// Lhfqdth lkzz gjkextybz vbuhfwbq bp afqkjd
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to the storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to the migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "", "name of the migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage path is required")
	}
	if migrationsPath == "" {
		panic("migrations path is required")
	}

	absMigrations, err := filepath.Abs(migrationsPath)
	if err != nil {
		panic(fmt.Sprintf("migrations path: %v", err))
	}
	if info, err := os.Stat(absMigrations); err != nil {
		if os.IsNotExist(err) {
			panic("migrations path does not exist: " + absMigrations)
		}
		panic(fmt.Sprintf("migrations path: %v", err))
	} else if !info.IsDir() {
		panic("migrations path is not a directory: " + absMigrations)
	}
	entries, _ := os.ReadDir(absMigrations)
	if len(entries) == 0 {
		panic("migrations directory is empty (add at least one migration file): " + absMigrations)
	}

	m, err := migrate.New(
		"file://"+filepath.ToSlash(absMigrations),
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		} else {
			panic(err)
		}
	}

	fmt.Println("migrations applied successfully")
}
