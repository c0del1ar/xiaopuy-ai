package main

import (
	"log"
	"net/http"
	"os"

	"github.com/c0del1ar/xiaopuy-ai/internal/ai"
	"github.com/c0del1ar/xiaopuy-ai/internal/chat"
	"github.com/c0del1ar/xiaopuy-ai/internal/router9"
)

func main() {
	provider := router9.New(
		os.Getenv("ROUTER9_BASE_URL"),
		os.Getenv("ROUTER9_API_KEY"),
		os.Getenv("ROUTER9_MODEL"),
	)

	agent := &ai.Agent{
		Provider: provider,
		Persona:  ai.DefaultPersona(),
	}
	chatService := &chat.Service{Agent: agent}
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
