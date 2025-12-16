FROM golang:1.24.1 AS builder

# Install tesseract and leptonica for OCR support
RUN apt-get update && apt-get install -y \
    libtesseract-dev \
    libleptonica-dev \
    tesseract-ocr \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the API server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-server ./cmd/api/main.go

# Production image
FROM debian:bullseye-slim AS production

# Install CA certificates, curl for health checks, and OCR dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    libtesseract4 \
    libleptonica5 \
    tesseract-ocr \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/api-server .

EXPOSE 8080

CMD ["./api-server"]