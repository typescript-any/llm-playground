package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/typescript-any/llm-playground/internal/config"
)

func main() {
	// Flags: up, down, drop
	direction := flag.String("direction", "up", "migration direction: up, down, drop")
	steps := flag.Int("steps", 0, "number of steps for down (0 = all)")
	flag.Parse()

	cfg := config.LoadConfig()

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatalf("Migration init failed: %v", err)
	}

	switch *direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("Migration applied successfully")
	case "down":
		if *steps > 0 {
			if err := m.Steps(-*steps); err != nil {
				log.Fatalf("Migration down steps failed: %v", err)
			}
			fmt.Printf("Rolled back %d step(s)\n", *steps)
		} else {
			if err := m.Down(); err != nil {
				log.Fatalf("Migration down failed: %v", err)
			}
		}
	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("Migration drop failed: %v", err)
		}
		fmt.Println("Dropped everything")
	default:
		log.Fatalf("Unknown direction: %s", *direction)

	}
}
