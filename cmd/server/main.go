package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/c0del1ar/xiaopuy-ai/internal/ai"
	"github.com/c0del1ar/xiaopuy-ai/internal/router9"
)

func main() {
	baseURL := os.Getenv("ROUTER9_BASE_URL")
	if baseURL == "" {
		log.Println("warning: ROUTER9_BASE_URL is not configured")
	}

	provider := router9.New(
		baseURL,
		os.Getenv("ROUTER9_API_KEY"),
		os.Getenv("ROUTER9_MODEL"),
	)

	agent := &ai.Agent{
		Provider: provider,
		Persona:  ai.DefaultPersona(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("POST /internal/test-reply", func(w http.ResponseWriter, r *http.Request) {
		response, err := agent.Reply(context.Background(), []ai.Message{
			{Role: "user", Content: "Hello"},
		}, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprint(w, response.Content)
	})

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("xiaopuy-ai listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
