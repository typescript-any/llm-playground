package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typescript-any/llm-playground/internal/config"
	"github.com/typescript-any/llm-playground/internal/db"
	"github.com/typescript-any/llm-playground/internal/handler"
	"github.com/typescript-any/llm-playground/internal/repository"
	"github.com/typescript-any/llm-playground/internal/routes"
	service "github.com/typescript-any/llm-playground/internal/services"
)

func SetupApp(cfg *config.Config) (*fiber.App, *pgxpool.Pool) {
	db.Init(cfg)
	pool := db.GetPool()
	if pool == nil {
		log.Fatal("❌ DB pool is nil! Did you call db.Init?")
	}

	convRepo := repository.NewConversationRepo(pool)
	convService := service.NewConversationService(convRepo)
	convHandler := handler.NewConversationHandler(convService)

	app := fiber.New()
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendString("ok") })
	routes.RegisterConversationRoutes(app, convHandler)

	return app, pool
}

func GracefulShutdown(app *fiber.App, pool *pgxpool.Pool) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("❌ Server shutdown failed: %v", err)
	}

	if pool != nil {
		pool.Close()
	}

	log.Println("🛑 Server exited gracefully")
}
