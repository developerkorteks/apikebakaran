# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install required packages
RUN apk --no-cache add ca-certificates curl bash vnstat

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create necessary directories
RUN mkdir -p /etc/apivpn /etc/xray /var/lib/scrz-prem

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]