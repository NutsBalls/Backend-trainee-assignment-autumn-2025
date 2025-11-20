package database

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func NewConn(databaseURL string) *pgxpool.Pool {
	log.Printf("Connecting to database with URL: %s", databaseURL)
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	SetupMigrations(pool)

	return pool
}

func SetupMigrations(pool *pgxpool.Pool) {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	migrationsDir := filepath.Join(wd, "migrations")
	log.Printf("Looking for migrations in: %s", migrationsDir)

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Printf("Migrations directory does not exist at %s, listing current directory:", migrationsDir)
		
		if files, err := os.ReadDir(wd); err == nil {
			log.Printf("Files in current directory (%s):", wd)
			for _, file := range files {
				log.Printf(" - %s", file.Name())
			}
		}
		
		log.Fatal("Migrations directory not found")
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	// Применяем миграции
	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}