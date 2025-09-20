package main

import (
	"fmt"
	"log"

	"github.com/typescript-any/llm-playground/internal/app"
	"github.com/typescript-any/llm-playground/internal/config"
)

func main() {
	// Load config
	cfg := config.LoadConfig()
	fiberApp, pool := app.SetupApp(cfg)

	addr := fmt.Sprintf(":%s", cfg.Port)
	started := make(chan bool)

	go func() {
		started <- true // signal start attempt
		if err := fiberApp.Listen(addr); err != nil {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	<-started
	log.Printf("🚀 Server is running on http://localhost:%s", cfg.Port)

	app.GracefulShutdown(fiberApp, pool)
}
