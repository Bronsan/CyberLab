# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Runtime stage — distroless for minimum attack surface
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy example config (user must provide their own config.yaml)
COPY config/config.example.yaml ./config/config.yaml

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./server"]
