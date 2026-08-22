package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/c0del1ar/xiaopuy-ai/internal/ai"
	"github.com/c0del1ar/xiaopuy-ai/internal/chat"
	"github.com/c0del1ar/xiaopuy-ai/internal/config"
	"github.com/c0del1ar/xiaopuy-ai/internal/router9"
	"github.com/c0del1ar/xiaopuy-ai/internal/storage/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := config.LoadDotEnv(".env"); err != nil {
		log.Fatalf("load .env: %v", err)
	}

	provider := router9.New(
		os.Getenv("ROUTER9_BASE_URL"),
		os.Getenv("ROUTER9_API_KEY"),
		os.Getenv("ROUTER9_MODEL"),
	)

	agent := &ai.Agent{
		Provider: provider,
		Persona:  ai.DefaultPersona(),
	}

	var repository chat.Repository
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			cancel()
			log.Fatalf("create PostgreSQL pool: %v", err)
		}
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			cancel()
			log.Fatalf("ping PostgreSQL: %v", err)
		}
		if err := postgres.Migrate(ctx, pool); err != nil {
			pool.Close()
			cancel()
			log.Fatalf("migrate PostgreSQL: %v", err)
		}
		cancel()
		defer pool.Close()
		repository = postgres.New(pool)
		log.Println("PostgreSQL persistence enabled")
	} else {
		repository = chat.NewMemoryRepository()
		log.Println("DATABASE_URL is not set; using in-memory persistence")
	}

	chatService := &chat.Service{
		Agent:      agent,
		Repository: repository,
	}
	chatHandler := &chat.HTTPHandler{Service: chatService}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("POST /v1/chat/reply", chatHandler.ReplyHTTP)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("xiaopuy-ai listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
