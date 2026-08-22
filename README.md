# Xiaopuy AI

Personal AI assistant for owner and client conversations.

## Current scope

- Go-based AI core
- Provider abstraction
- 9router integration (OpenAI-compatible chat endpoint)
- Configurable private/client persona
- Minimal HTTP health endpoint
- Foundation for conversation, RAG, memory, and messaging channels

## Architecture

```text
WhatsApp / Telegram / Web
          |
     Channel Adapter
          |
      AI Core
      /  |  \
 Persona RAG Memory
          |
       9router
          |
    Model Rotation
```

## Development

Copy `.env.example` to `.env` and configure the 9router endpoint and credentials.

```bash
go run ./cmd/server
```

Health check:

```bash
curl http://localhost:8080/health
```

The 9router client currently expects an OpenAI-compatible endpoint at `/v1/chat/completions`. Adjust the provider contract if the private 9router API differs.
