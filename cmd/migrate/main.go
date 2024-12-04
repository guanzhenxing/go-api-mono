package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [up|down]")
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	migrationsPath := filepath.Join(dir, "scripts/migrations")
	dsn := "mysql://root:123456@tcp(localhost:3306)/go_api_mono?multiStatements=true"

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migration up completed successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migration down completed successfully")
	default:
		log.Fatal("Invalid command. Use 'up' or 'down'")
	}
}
