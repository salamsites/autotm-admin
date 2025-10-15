# Build stage
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Copy go.mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/autotm-admin ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/autotm-admin .
COPY config.yml .

# Create logs directory
RUN mkdir -p /home/user/Desktop/sada/autotm-admin/logs

# Expose port
EXPOSE 8006

# Run the binary
CMD ["./autotm-admin"]