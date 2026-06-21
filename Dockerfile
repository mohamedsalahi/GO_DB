# Frontend build stage
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci

COPY frontend/ .
RUN npx vite build

# Go build stage
FROM golang:1.23-alpine AS go-builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copy the built frontend into the Go build context
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always 2>/dev/null || echo 'dev')" \
    -o /app/api ./cmd/api

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=go-builder /app/api .
COPY --from=go-builder /app/.env.example .env
COPY --from=go-builder /app/migrations ./migrations
COPY --from=go-builder /app/frontend/dist ./frontend/dist

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./api"]
