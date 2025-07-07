# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go.mod and go.sum files into /build
COPY go.mod go.sum ./

# Download dependencies and verify
RUN go mod download && go mod verify

# Copy the source code to build (exclude all files in  .dockerignore)
COPY . .

# Build the application
# CGO_ENABLED=0 => Disable
RUN CGO_ENABLED=0 GOOS=linux go build -o /app .

# Final stage
FROM alpine:3.19

WORKDIR /var/

# Copy the binary from the builder stage
#`builder` at line 2
COPY --from=builder /app .  

# Expose the port the app runs on
EXPOSE 3000

# Command to run the application
CMD ["./app"]
