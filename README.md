# Debate Chatbot

A full-stack debate chatbot application with API key authentication, built with Go backend and Next.js frontend.

## 🚀 Features

- **AI-Powered Debates**: Chat with an AI that takes a stance and argues persuasively
- **API Key Authentication**: Secure access with X-API-Key header authentication
- **Chat History**: Browse and continue past conversations
- **Responsive Design**: Works on desktop and mobile
- **HTTPS Support**: Secure communication via CloudFront

## 🏗️ Architecture

```
Frontend (Next.js) → CloudFront (HTTPS) → ALB (HTTP) → ECS (Go API) → PostgreSQL/Redis
```

**HTTPS Solution**: CloudFront terminates HTTPS and communicates with ALB over HTTP internally. This prevents mixed content errors while keeping the setup simple.

## 🛠️ Quick Start

### Prerequisites
- Docker & Docker Compose
- OpenAI API key
- API key for authentication (provided separately)

### Local Development

1. **Clone and setup:**
```bash
git clone <repository-url>
cd debate
```

2. **Set environment variables:**
```bash
export OPENAI_API_KEY="your_openai_api_key_here"
```

3. **Run with PostgreSQL:**
```bash
docker compose -f docker-compose.postgres.yml up --build -d
```

4. **Access the application:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

### Production Deployment

1. **Deploy infrastructure:**
```bash
cd terraform/projects/backend
terraform apply
```

2. **Get your API key** from the administrator

3. **Access the application:**
- Frontend: https://d1iy0shli1spap.cloudfront.net
- Backend API: https://d13sbjy1c5yh6c.cloudfront.net

## 🔑 Authentication

The chatbot requires an API key to access. Contact your administrator to get your API key, then:

1. Visit the login page
2. Enter your API key
3. Start debating!

## 📱 Usage

1. **Login** with your API key
2. **Start a new chat** or continue an existing conversation
3. **Send messages** to engage in debates
4. **Browse history** to see past conversations

## 🧰 Tech Stack

- **Backend**: Go, Gin, OpenAI API, PostgreSQL, Redis
- **Frontend**: Next.js 15, TypeScript, Tailwind CSS, shadcn/ui
- **Infrastructure**: AWS ECS, ALB, CloudFront, RDS Aurora, ElastiCache
- **Security**: API key authentication, HTTPS, KMS encryption

## 📊 API Endpoints

- `POST /chat` - Send message and get bot response
- `GET /conversations` - List conversations
- `GET /conversations/:id` - Get specific conversation
- `GET /health` - Health check

## 🧪 Testing the API

### With curl

```bash
# Test health endpoint
curl -H "X-API-Key: your-api-key-here" \
  https://your-api-url/health

# Start a new conversation
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key-here" \
  -d '{"message": "Let'\''s debate about renewable energy!"}' \
  https://your-api-url/chat

# Continue existing conversation
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key-here" \
  -d '{"conversation_id": "01HZ1234567890", "message": "What about the costs?"}' \
  https://your-api-url/chat

# List conversations
curl -H "X-API-Key: your-api-key-here" \
  https://your-api-url/conversations
```

### With the test page

1. Open `frontend/test-auth.html` in your browser
2. Enter your API key and API URL
3. Click the test buttons to verify endpoints

## 🔧 Development

```bash
# Install dependencies
make install

# Run tests
make test

# Build and run
make run

# Clean up
make clean
```

## 📝 License

MIT License - see LICENSE file for details.
