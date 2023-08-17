package main

import (
	"collect/business/data/dbschema"
	"context"
	"fmt"
	"log"

	"collect/business/sys/database"
	"collect/foundation/config"

	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	cfg, err := config.LoadConfig("dev", ".")
	if err != nil {
		return err
	}

	dbConfig := cfg.GetDBConfig()

	if err := Migrate(dbConfig); err != nil {
		return err
	}

	if err := Seed(dbConfig); err != nil {
		return err
	}

	return nil
}

// Migrate creates the schema in the database.
func Migrate(cfg config.DatabaseConfigurations) error {
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 160*time.Second)
	defer cancel()

	if err := dbschema.DropAll(ctx, db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	if err := dbschema.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migrations complete")
	return nil
}

// Seed loads test data into the database.
func Seed(cfg config.DatabaseConfigurations) error {
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	if err := dbschema.Seed(ctx, db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("seed data complete")
	return nil
}
