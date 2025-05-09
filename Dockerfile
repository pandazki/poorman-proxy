# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to leverage Docker cache
COPY go.mod go.sum* ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Ensure secret.json is available for embedding
# The user building the image must ensure secret/secret.json exists in the build context
COPY secret/secret.json secret/secret.json

# Build the statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/proxy ./cmd/proxy/...

# Stage 2: Create the runtime image
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/proxy .

# Expose the port the application runs on
EXPOSE 8080

# Set the command to run the application
# Users can pass environment variables for secrets and outbound proxy configuration
# e.g., docker run -p 8080:8080 -e OPENAI_API_KEY="sk-..." -e OUTBOUND_PROXY_URL="socks5://..." my-proxy-image
CMD ["/app/proxy"] 