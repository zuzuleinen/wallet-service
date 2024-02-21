package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"wallet-service/application"
	"wallet-service/handlers"
	"wallet-service/infrastructure"
)

type Config struct {
	Host   string
	Port   string
	DbName string
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg := Config{
		Host:   "localhost",
		Port:   "8080",
		DbName: "test.db",
	}

	// Init database
	db, err := infrastructure.InitDatabase(cfg.DbName)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Println("stopping database")
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	ws := application.NewWalletService(infrastructure.NewTransactionRepository(db))

	srv := NewServer(ws)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: srv,
	}
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Println("got interrupt signal. shutdown gracefully")
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

func NewServer(ws *application.WalletService) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, ws)

	var handler http.Handler = mux
	return handler
}

func addRoutes(mux *http.ServeMux, ws *application.WalletService) {
	mux.Handle("/health", handlers.HealthHandler())
	mux.Handle("GET /wallet/{userId}", handlers.GetWalletHandler(ws))
	mux.Handle("POST /wallet/{userId}", handlers.CreateWalletHandler(ws))
	mux.Handle("POST /add-funds/{userId}", handlers.AddFundsHandler(ws))
	mux.Handle("POST /remove-funds/{userId}", handlers.RemoveFundsHandler(ws))
}
