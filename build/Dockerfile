# Use golang:1.19 as the builder stage
FROM golang:1.22.2-alpine AS builder

# Set working directory
WORKDIR /lcgodemo

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o lcgodemo cmd/main.go

# Use alpine base for final image
FROM alpine:latest

# Install dependencies for reading PDFs
RUN apk add poppler-utils

# Create a non-root user to run the app
RUN adduser -D apprunner

# Set working directory
WORKDIR /lcgodemo

# Copy the built binary from builder stage
COPY --from=builder /lcgodemo /lcgodemo

# Set user for the process
RUN chown -R apprunner:apprunner /lcgodemo

# Expose port 8080
EXPOSE 8080

# Run the application as the non-root user
USER apprunner

# Set the default command to run the app
CMD ["./lcgodemo"]
