# Build Stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/main.go

# Production Stage
FROM alpine:latest AS production

# Set working directory for clarity
WORKDIR /app

# Copy only the binary from the builder stage
COPY --from=builder /app/app .
COPY migrations /app/migrations

# Run the m
CMD ["./app"]