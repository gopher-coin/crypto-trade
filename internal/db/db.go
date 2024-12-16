package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool

func ConnectDB() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("DATABASE_URL is not set. Please configure it in your .env file\n")
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w\n", err)
	}

	pool, err := pgxpool.New(context.Background(), config.ConnString())
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w\n", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("failed to ping database: %w\n", err)
	}

	Pool = pool
	fmt.Println("Connected to the database!")

	err = runMigrations(pool)
	if err != nil {
		return fmt.Errorf("failed to run migration: %w\n", err)
	}
	return nil
}

func CloseDB() {
	if Pool != nil {
		Pool.Close()
		fmt.Println("Database connection pool closed!")
	}
}

func runMigrations(pool *pgxpool.Pool) error {
	migrationsPath := "./internal/db/migrations"

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w\n", err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(migrationsPath, file.Name())

			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w\n", file.Name(), err)
			}

			_, err = pool.Exec(context.Background(), string(content))
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w\n", file.Name(), err)
			}

			fmt.Printf("Migration %s applied successfully\n", file.Name())
		}
	}
	return nil
}
