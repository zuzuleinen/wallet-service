package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"wallet-service/app"
	"wallet-service/handlers"
	"wallet-service/infra/db"
	"wallet-service/infra/pubsub"

	"github.com/apache/pulsar-client-go/pulsar"
)

type Config struct {
	Host   string
	Port   string
	DbName string
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, out io.Writer, getenv func(string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Init app config
	cfg := Config{
		Host:   withDefault(getenv("WALLET_HOST"), "0.0.0.0"),
		Port:   withDefault(getenv("WALLET_PORT"), "8081"),
		DbName: withDefault(getenv("WALLET_DB_NAME"), "dev.db"),
	}

	// Init database
	database, err := db.InitDatabase(cfg.DbName)
	if err != nil {
		return fmt.Errorf("connecting to database: %s", err)
	}
	defer func() {
		log.Println("stopping database")
		sqlDB, _ := database.DB()
		sqlDB.Close()
	}()

	// Init Pulsar Client
	pulsarClient, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: withDefault(getenv("PULSAR_CLIENT_URL"), "pulsar://localhost:6650"),
	})
	if err != nil {
		log.Fatalf("error creating pulsar client: %s", err)
	}
	defer pulsarClient.Close()

	producer, err := pulsarClient.CreateProducer(pulsar.ProducerOptions{
		Topic: pubsub.TopicTransactions,
	})
	if err != nil {
		log.Fatalf("error creating producer: %s", err)
	}

	// Init logger and WalletService
	logger := log.New(out, "", log.LstdFlags)
	ws := app.NewWalletService(producer, db.NewTransactionRepository(database), logger)

	// Start server
	srv := NewServer(ws, logger)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: srv,
	}
	go func() {
		log.Printf("server started on %s\n", httpServer.Addr)
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
		walletServiceShutdownCtx, cancelWallet := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancelWallet()
		ws.Stop(walletServiceShutdownCtx)
	}()
	wg.Wait()
	return nil
}

// withDefault returns want if not empty else will return default def
func withDefault(want, def string) string {
	if want != "" {
		return want
	}
	return def
}

func NewServer(ws *app.WalletService, logger *log.Logger) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, ws, logger)

	var handler http.Handler = mux
	return handler
}

func addRoutes(mux *http.ServeMux, ws *app.WalletService, logger *log.Logger) {
	mux.Handle("/health", handlers.HealthHandler(logger))
	mux.Handle("GET /wallet/{userId}", handlers.GetWalletHandler(ws))
	mux.Handle("POST /wallet/{userId}", handlers.CreateWalletHandler(ws, logger))
	mux.Handle("POST /add-funds/{userId}", handlers.AddFundsHandler(ws, logger))
	mux.Handle("POST /remove-funds/{userId}", handlers.RemoveFundsHandler(ws))
}
