# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 设置国内 Go 代理
ENV GOPROXY=https://goproxy.cn,direct

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application statically (important for alpine)
# CGO_ENABLED=0 prevents linking against C libraries
# -ldflags="-w -s" reduces binary size (optional)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/main ./cmd/api/main.go

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/main /app/main

# Expose the port the application runs on (as defined in .env or default 8080)
EXPOSE 8080

# Command to run the application
CMD ["/app/main"] 