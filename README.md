# Kopi Challenge – Debate Chatbot API (Go)

A Go API for a debate chatbot that picks a topic, takes a side, and tries to persuade the user while staying on-topic and consistent. Uses OpenAI for text generation, Redis (optional) for conversation state.

## ✨ Features
- `POST /chat` with `{ conversation_id|null, message }`
- Starts a new conversation when `conversation_id` is `null` (chooses topic + stance)
- Maintains coherent stance across turns
- Returns the **last 5 messages on both sides** (max 10 items), most recent last
- 30s response cap (25s LLM timeout + overhead)

## 🧰 Tech
- Go 1.22 + Gin
- OpenAI Chat Completions API (configurable model)
- Redis (via Docker) or in-memory fallback
- Docker + docker-compose

## 🚀 Run

### **Local Development (Recommended)**

1. **Set up environment variables:**
   ```bash
   # Copy the example file
   cp .env.example .env
   
   # Edit .env and add your OpenAI API key
   # OPENAI_API_KEY=sk-your-actual-key-here
   ```

2. **Run with Docker Compose:**
   ```bash
   make install
   make run
   ```

### **Alternative: Environment Variables**
```bash
export OPENAI_API_KEY=sk-...      # required
export OPENAI_MODEL=gpt-4o-mini   # optional (default: gpt-4o-mini)

make install
make run
```

The API listens on `http://localhost:8080`.

### Health
- `GET /healthz` → liveness
- `GET /readyz` → readiness (checks store)

### Example Request
```bash
curl -s localhost:8080/chat \
  -H 'content-type: application/json' \
  -d '{"conversation_id":null, "message":"No way the Earth is flat."}' | jq .
```

### Example Response (shape)
```json
{
  "conversation_id": "b1e1...",
  "message": [
    { "role": "user", "message": "No way the Earth is flat." },
    { "role": "bot",  "message": "Actually... [persuasive reply]" }
  ]
}
```

## 🧪 Test
```bash
make test
```
(Uses a mock LLM engine; no network calls.)

## 🛠 Configuration (ENV)
- `OPENAI_API_KEY` (required) – your OpenAI key
- `OPENAI_MODEL` (optional) – default `gpt-4o-mini`
- `PORT` (optional) – default `8080`
- `REDIS_ADDR` (optional) – if omitted, in-memory store is used
- `REDIS_PASSWORD` (optional)

## 🧱 Notes
- Storage gracefully falls back to memory if Redis is not available.
- Conversation history capped defensively at 200 messages; responses return only last 10 (5 each side).
- Prompt enforces staying on the original topic and stance.

## 🔒 Security
- Never log the API key.
- Prefer setting env vars via a `.env` file or your CI/CD secrets manager.
