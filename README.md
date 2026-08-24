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

## Website knowledge ingestion

With PostgreSQL/pgvector, `ROUTER9_EMBEDDING_MODEL`, and `INGEST_ALLOWED_DOMAINS` configured, crawl trusted website content into RAG with:

```bash
curl -X POST http://localhost:8080/v1/ingest/crawl \
  -H 'Content-Type: application/json' \
  -d '{"seed_url":"https://aryakun.id/"}'
```

The crawler is restricted to `INGEST_ALLOWED_DOMAINS`; set this only to domains you trust as knowledge sources.
