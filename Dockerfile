# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
#COPY go.mod go.sum .

# Download dependencies
# RUN go get .

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/blog .

# Copy the posts directory
COPY --from=builder /app/posts ./posts

# Copy the pictures directory
COPY --from=builder /app/pictures ./pictures

# Expose port (Railway will set PORT env var)
EXPOSE 8080

# Run the application
CMD ["./blog"]
