# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server ./cmd/server

# Runtime stage
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /
COPY --from=build /bin/server /server
ENV PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
