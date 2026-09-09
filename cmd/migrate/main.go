package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/magdramazz/olx-api/internal/config"
)

func main() {
	cfg := config.MustLoad()
	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}
	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migration.up: %v", err)
		}
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("migration.down: %v", err)
		}
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
	fmt.Println("Migration completed successfully")
}
