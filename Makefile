# Docker data
IMAGE_NAME:=debate-chatbot-backend
BETA=-beta
# AWS data
REGION=us-east-2
ACCOUNT_ID=215229807348

ECR_REPOSITORY=$(ACCOUNT_ID).dkr.ecr.$(REGION).amazonaws.com/$(IMAGE_NAME)

LOGIN:=aws ecr get-login-password --region $(REGION) | docker login --username AWS --password-stdin $(ECR_REPOSITORY)

# Binary file data
BINARY_NAME=debate-chatbot-backend

# API testing
API_URL=http://debate-chatbot-alb-1114917533.us-east-2.elb.amazonaws.com
API_KEY=your-secure-api-key-here

.PHONY: make install run down clean test staging login build_image generate_binary push_image clean_binary test-api test-api-auth test-api-no-auth deploy-frontend build-frontend

make:
	@echo "Available commands:"
	@echo "  make install        Install dependencies and tools"
	@echo "  make run            Run API + Redis in Docker"
	@echo "  make down           Stop containers"
	@echo "  make clean          Remove containers + volumes"
	@echo "  make test           Run unit tests"
	@echo "  make staging        Build and push to ECR (staging)"
	@echo "  make login          Login to ECR"
	@echo "  make build_image    Build Docker image for ECR"
	@echo "  make push_image     Push image to ECR"
	@echo "  make clean_binary   Remove binary file"
	@echo "  make test-api       Test API health and readiness"
	@echo "  make test-api-auth  Test API with authentication"
	@echo "  make test-api-no-auth Test API without authentication (should fail)"
	@echo "  make build-frontend Build the frontend application"
	@echo "  make deploy-frontend Deploy frontend to S3 and invalidate CloudFront"

install:
	@which docker >/dev/null 2>&1 || { echo "Docker not found. Install Docker: https://docs.docker.com/get-docker/"; exit 1; }
	@which docker-compose >/dev/null 2>&1 || { echo "docker-compose not found. Install: https://docs.docker.com/compose/install/"; exit 1; }
	@which go >/dev/null 2>&1 || { echo "Go not found. Install: https://go.dev/dl/"; exit 1; }
	@which aws >/dev/null 2>&1 || { echo "AWS CLI not found. Install: https://aws.amazon.com/cli/"; exit 1; }
	go mod tidy

run:
	docker-compose up --build

down:
	docker-compose down

clean:
	docker-compose down -v

test:
	go test ./...

staging:
	$(MAKE) login
	$(MAKE) build_image
	$(MAKE) push_image
	$(MAKE) clean_binary

generate_binary:
	go mod download
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -tags truora.nomock,truora.extlogs -ldflags="-s -w" -o $(BINARY_NAME) cmd/server/main.go
	export CGO_ENABLED=1

login:
	$(LOGIN)

build_image: generate_binary
	docker build --platform linux/amd64 -t $(IMAGE_NAME):latest --build-arg ACCOUNT_ID=$(ACCOUNT_ID) .
	docker tag $(IMAGE_NAME):latest $(ECR_REPOSITORY):latest

push_image:
	docker push $(ECR_REPOSITORY):latest

clean_binary:
	rm -f $(BINARY_NAME)

# API Testing Commands
test-api:
	@echo "=== Testing API Health ==="
	@echo "1. Health endpoint:"
	@curl -s -X GET "$(API_URL)/health" | jq . || echo "Failed to connect"
	@echo
	@echo "2. Ready endpoint:"
	@curl -s -X GET "$(API_URL)/ready" | jq . || echo "Failed to connect"
	@echo

test-api-auth:
	@echo "=== Testing API with Authentication ==="
	@echo "1. Chat endpoint with API key:"
	@curl -s -X POST "$(API_URL)/chat" \
		-H "Content-Type: application/json" \
		-H "X-API-Key: $(API_KEY)" \
		-d '{"message": "Hello, test message", "conversation_id": ""}' | jq . || echo "Failed to connect"
	@echo
	@echo "2. Conversations endpoint with API key:"
	@curl -s -X GET "$(API_URL)/conversations" \
		-H "X-API-Key: $(API_KEY)" | jq . || echo "Failed to connect"
	@echo

test-api-no-auth:
	@echo "=== Testing API without Authentication (should fail) ==="
	@echo "1. Chat endpoint without API key:"
	@curl -s -X POST "$(API_URL)/chat" \
		-H "Content-Type: application/json" \
		-d '{"message": "Hello, test message", "conversation_id": ""}' | jq . || echo "Failed to connect"
	@echo
	@echo "2. Conversations endpoint without API key:"
	@curl -s -X GET "$(API_URL)/conversations" | jq . || echo "Failed to connect"
	@echo

# Frontend Deployment Commands
build-frontend:
	@echo "=== Building Frontend Application ==="
	@cd frontend && npm install
	@cd frontend && npm run build
	@echo "Frontend build completed successfully"

deploy-frontend:
	@echo "=== Deploying Frontend to S3 and CloudFront ==="
	@./scripts/deploy-frontend.sh
