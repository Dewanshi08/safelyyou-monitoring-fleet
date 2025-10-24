# Stage 1 — Build the Go binary
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first — enables better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the app code
COPY . .

# Build binary (static build)
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./main.go

# Stage 2 — Run the binary
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy compiled binary from builder stage
COPY --from=builder /app/server .

# Copy devices.csv
COPY --from=builder /app/devices.csv .

# Expose port if your app listens on it
EXPOSE 6733

# Run the binary
CMD ["./server"]
