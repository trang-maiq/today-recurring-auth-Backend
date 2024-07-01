# Use the official Golang image as a builder
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install necessary dependencies
RUN apk update && apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached    if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a minimal image for the final container
FROM alpine:latest

# Install necessary dependencies for running the app
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Command to run the executable
CMD ["./main"]

# Expose port if needed
EXPOSE 4242
