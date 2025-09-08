# Debate Chatbot - Full-Stack Application

A sophisticated debate chatbot that supports multiple conversations, chat history, and topic management. Built with Go backend, Next.js frontend, and PostgreSQL/Redis storage options.

## 🚀 **Features**

### **Core Chatbot Features**
- `POST /chat` with `{ conversation_id|null, message }`
- Starts a new conversation when `conversation_id` is `null` (chooses topic + stance)
- Maintains coherent stance across turns
- Returns the **last 5 messages on both sides** (max 10 items), most recent last
- 30s response cap (25s LLM timeout + overhead)

### **Multi-Chat Support**
- **Chat History**: Browse and continue past conversations
- **Topic Management**: Organize debates by topic
- **Popular Topics**: See trending debate topics
- **Conversation Persistence**: All chats are saved and searchable
- **Responsive Design**: Works seamlessly on desktop and mobile devices

### **Database Options**
- **PostgreSQL**: Full-featured with complex queries and analytics
- **Redis**: Fast caching with fallback functionality
- **Hybrid Approach**: PostgreSQL for persistence, Redis for caching

## 🏗️ **Architecture**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   Database      │
│   (Next.js)     │◄──►│   (Go + Gin)    │◄──►│   (PostgreSQL)  │
│                 │    │                 │    │                 │
│ • Chat UI       │    │ • REST API      │    │ • Conversations │
│ • History       │    │ • CORS Support  │    │ • Messages      │
│ • Navigation    │    │ • OpenAI Bot    │    │ • Topics        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   Cache Layer   │
                       │   (Redis)       │
                       │                 │
                       │ • Session Cache │
                       │ • Fast Access   │
                       └─────────────────┘
```

## 🧰 **Tech Stack**

### **Backend**
- Go 1.23 + Gin Web Framework
- OpenAI Chat Completions API (configurable model)
- PostgreSQL with complex queries and analytics
- Redis for caching and session management
- ULID for unique, sortable identifiers
- Docker + docker-compose

### **Frontend**
- Next.js 15 with App Router
- TypeScript for type safety
- Tailwind CSS for styling
- shadcn/ui components
- Lucide React icons
- Responsive design

## 📊 **Database Schema**

### **PostgreSQL Tables**

```sql
-- Users (for future user management)
users (id, email, name, created_at)

-- Topics (debate subjects)
topics (id, name, description, category, created_at)

-- Conversations (chat sessions)
conversations (id, user_id, topic_id, topic_name, bot_stance, title, created_at, updated_at, message_count)

-- Messages (individual chat messages)
messages (id, conversation_id, role, content, created_at)
```

## 🛠️ **Setup Instructions**

### **Option 1: PostgreSQL (Recommended)**

1. **Start the enhanced Docker Compose:**
```bash
docker compose -f docker-compose.postgres.yml up --build -d
```

2. **Set environment variables:**
```bash
export OPENAI_API_KEY="your_openai_api_key_here"
```

3. **Access the application:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432

### **Option 2: Redis Only**

```bash
docker compose up --build -d
```

### **Option 3: Local Development**

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

3. **Alternative: Environment Variables**
   ```bash
   export OPENAI_API_KEY=sk-...      # required
   export OPENAI_MODEL=gpt-4o-mini   # optional (default: gpt-4o-mini)
   
   make install
   make run
   ```

The API listens on `http://localhost:8080`.

## 🔌 **API Endpoints**

### **Chat Endpoints**
- `POST /chat` - Send a message and get bot response
- `GET /healthz` - Liveness check
- `GET /readyz` - Readiness check (checks store)

### **Conversation Management**
- `GET /conversations` - List all conversations with pagination
- `GET /conversations/:id` - Get a specific conversation
- `GET /conversations/topics` - Get popular debate topics

### **Example API Calls**

```bash
# Send a chat message
curl -s localhost:8080/chat \
  -H 'content-type: application/json' \
  -d '{"conversation_id":null, "message":"No way the Earth is flat."}' | jq .

# List conversations
curl "http://localhost:8080/conversations?limit=10&offset=0"

# Get popular topics
curl "http://localhost:8080/conversations/topics?limit=5"

# Health check
curl http://localhost:8080/healthz
```

### **Example Response (shape)**
```json
{
  "conversation_id": "01K4KN1YB9DHWNAJDSYK3QJA2F",
  "message": [
    { "role": "user", "message": "No way the Earth is flat.", "ts": 1703123456789 },
    { "role": "bot",  "message": "Actually... [persuasive reply]", "ts": 1703123456790 }
  ]
}
```

## 🎨 **Frontend Features**

### **Pages**
- **Home** (`/`) - Landing page with features
- **Chat** (`/chat`) - Main debate interface with sidebar history
- **Responsive Design** - Mobile-friendly with collapsible sidebar

### **Chat Features**
- **Real-time Chat Interface**: Clean, responsive chat UI with message history
- **AI Debate Partner**: Connects to the Go backend for intelligent debate conversations
- **Topic Display**: Shows current debate topic and bot stance
- **Message History**: Full conversation history with auto-scroll
- **Continue Conversations**: Resume from conversation history via sidebar
- **Loading States**: Smooth loading indicators and error handling

### **History Features**
- **Sidebar Navigation**: Left sidebar with conversation history
- **Mobile Responsive**: Collapsible sidebar with toggle button
- **Quick Access**: Continue any past conversation
- **Statistics**: Message count and timestamps
- **Smooth Scrolling**: Proper overflow handling and scrolling

## 🔧 **Configuration**

### **Environment Variables**

```bash
# Required
OPENAI_API_KEY=your_openai_api_key_here

# Database (PostgreSQL)
POSTGRES_URL=postgres://user:password@localhost:5432/debate?sslmode=disable

# Cache (Redis)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=optional

# Server
PORT=8080
OPENAI_MODEL=gpt-4o-mini
```

### **Frontend Configuration**

```typescript
// frontend/lib/config.ts
export const config = {
    apiUrl: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
    endpoints: {
        chat: "/chat",
        conversations: "/conversations",
        topics: "/conversations/topics",
        health: "/healthz",
    },
};
```

## 🚀 **Deployment Options**

### **1. Local Development**
```bash
# Backend
go run ./cmd/server

# Frontend
cd frontend && npm run dev
```

### **2. Docker Compose**
```bash
# PostgreSQL version (recommended)
docker compose -f docker-compose.postgres.yml up -d

# Redis version
docker compose up -d
```

### **3. Production (AWS ECS)**
- Use PostgreSQL RDS for database
- Use ElastiCache for Redis
- Deploy containers to ECS with load balancer

## 🧪 **Testing**

```bash
# Run backend tests
make test
go test ./...

# Run linting
golangci-lint run

# Test API endpoints
curl http://localhost:8080/healthz
curl http://localhost:8080/conversations

# Frontend tests
cd frontend && npm run lint
```

## 📝 **Usage Examples**

### **Starting a New Debate**
1. Go to http://localhost:3000/chat
2. Send a message like "Let's debate about climate change"
3. Bot will choose a topic and stance
4. Continue the conversation

### **Browsing History**
1. Use the left sidebar in the chat interface
2. See all past debates with topic and message count
3. Click any conversation to continue
4. Use the "New Chat" button to start fresh

### **Popular Topics**
1. Visit http://localhost:8080/conversations/topics
2. See trending debate subjects
3. Use for inspiration or analytics

## 📈 **Performance Considerations**

### **Database Optimization**
- **Indexes**: Optimized for conversation queries
- **Pagination**: Efficient large dataset handling
- **Caching**: Redis for frequently accessed data
- **ULID**: Sortable, URL-safe identifiers

### **API Optimization**
- **CORS**: Properly configured for frontend
- **Error Handling**: Comprehensive error responses
- **Rate Limiting**: Ready for production scaling
- **Transaction Management**: Proper database transactions

## 🔒 **Security Features**

- **CORS Configuration**: Secure cross-origin requests
- **Input Validation**: Sanitized user inputs
- **SQL Injection Protection**: Parameterized queries
- **Environment Variables**: Secure secret management
- **Never log API keys**: Secure credential handling

## 🧱 **Technical Notes**

- Storage gracefully falls back to memory if Redis is not available
- Conversation history capped defensively at 200 messages; responses return only last 10 (5 each side)
- Prompt enforces staying on the original topic and stance
- Uses ULID for unique, sortable conversation identifiers
- Comprehensive test coverage with unit tests and benchmarks

## 🔮 **Future Enhancements**

- **User Authentication**: Login and personal chat history
- **Topic Suggestions**: AI-powered topic recommendations
- **Analytics Dashboard**: Debate statistics and insights
- **Export Features**: Download conversation history
- **Real-time Collaboration**: Multiple users in same debate
- **Voice Integration**: Speech-to-text and text-to-speech

## 🐛 **Troubleshooting**

### **Common Issues**

1. **CORS Errors**: Ensure backend is running and CORS is configured
2. **Database Connection**: Check PostgreSQL/Redis is running
3. **OpenAI API**: Verify API key is set correctly
4. **Port Conflicts**: Ensure ports 3000, 8080, 5432, 6379 are free
5. **Chat Input Missing**: Check responsive layout and flex constraints

### **Debug Commands**

```bash
# Check running containers
docker ps

# View logs
docker compose logs api
docker compose logs postgres

# Test database connection
psql postgres://debate_user:debate_password@localhost:5432/debate

# Test Redis connection
redis-cli ping

# Check frontend build
cd frontend && npm run build
```

## 📚 **Project Structure**

```
debate/
├── cmd/                    # Application entry points
│   └── server/            # Main server application
├── internal/              # Private application code
│   ├── api/              # HTTP handlers and routes
│   ├── bot/              # AI bot logic and prompts
│   ├── models/           # Data structures
│   └── storage/          # Database implementations
├── frontend/             # Next.js frontend application
│   ├── app/             # App Router pages
│   ├── components/      # React components
│   └── lib/             # Utility functions
├── test/                # Integration tests
├── docker-compose.yml   # Redis setup
├── docker-compose.postgres.yml  # PostgreSQL setup
├── Dockerfile          # Backend container
├── Makefile           # Build and run commands
└── init.sql           # Database schema
```

## 📚 **Documentation**

- **API Documentation**: Available at `/healthz` endpoint
- **Database Schema**: See `init.sql` for full schema
- **Frontend Components**: Check `frontend/app/` directory
- **Backend Structure**: See `internal/` package organization

---

**Ready to debate? Start your first conversation at http://localhost:3000/chat!** 🎯