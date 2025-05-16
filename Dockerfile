# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Create uploads directory for media files
RUN mkdir -p /app/uploads

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080
ENV MODULE=cart

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
