FROM golang:1.24.4-alpine AS builder 

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/main.go

# -------- Production Stage --------
FROM alpine:latest

# Set working dir in new container
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./server"]