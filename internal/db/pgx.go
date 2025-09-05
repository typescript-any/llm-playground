package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typescript-any/llm-playground/internal/config"
)

var pool *pgxpool.Pool

// Init initializes the global pgx pool
func Init(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	pool, err = pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Unable to connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("❌ Unable to ping database: %v", err)
	}

	log.Println("✅ Connected to database")
}

func GetPool() *pgxpool.Pool {
	return pool
}

func Close() {
	if pool != nil {
		pool.Close()
		log.Println("🛑 Database connection closed")
	}
}
