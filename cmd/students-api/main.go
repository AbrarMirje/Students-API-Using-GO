package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AbrarMirje/Students-API-Using-GO/internal/config"
	"github.com/AbrarMirje/Students-API-Using-GO/internal/http/handlers/student"
)

func main() {
	// Load config
	cfg := config.MustLoad()
	// Database setup
	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())
	// Setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.Addr))

	// Terminating server gracefully
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	slog.Info("Shutting down the server")
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown successfully!")
}
