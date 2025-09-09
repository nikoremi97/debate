# API Testing Guide

Simple commands to test your debate chatbot API with and without authentication.

## Quick Start

### 1. Test API Health
```bash
make test-api
```
Tests the basic health and readiness endpoints (no authentication required).

### 2. Test API Without Authentication
```bash
make test-api-no-auth
```
Tests protected endpoints without API key (should fail after authentication is enabled).

### 3. Test API With Authentication
```bash
make test-api-auth
```
Tests protected endpoints with API key (requires valid API key).

## Configuration

### Update API Key
Edit the `API_KEY` variable in the Makefile:
```makefile
API_KEY=your-actual-api-key-here
```

### Update API URL
Edit the `API_URL` variable in the Makefile:
```makefile
API_URL=https://your-api-domain.com
```

## What Each Command Tests

### `make test-api`
- ✅ `GET /health` - Basic health check
- ✅ `GET /ready` - Readiness check (includes database connectivity)

### `make test-api-no-auth`
- ❌ `POST /chat` - Should return 401 Unauthorized
- ❌ `GET /conversations` - Should return 401 Unauthorized

### `make test-api-auth`
- ✅ `POST /chat` - Should work with valid API key
- ✅ `GET /conversations` - Should work with valid API key

## Expected Results

### Before Authentication is Deployed
- All commands should work (API runs without authentication)

### After Authentication is Deployed
- `make test-api` - ✅ Works
- `make test-api-no-auth` - ❌ Returns 401 errors
- `make test-api-auth` - ✅ Works (with correct API key)

## Troubleshooting

### "Failed to connect" errors
- Check if the API URL is correct
- Verify the API is running and accessible
- Check network connectivity

### 401 Unauthorized errors (when expected)
- This is normal behavior when authentication is enabled
- Verify you're using the correct API key

### 401 Unauthorized errors (when not expected)
- Check if the API key is correct
- Verify the secret exists in AWS Secrets Manager
- Check ECS task logs for authentication errors

## Dependencies

- `curl` - For making HTTP requests
- `jq` - For JSON formatting (optional, for better output)

Install on macOS:
```bash
brew install curl jq
```

Install on Ubuntu/Debian:
```bash
sudo apt-get install curl jq
```
