# Use the official Go image for building the application
FROM golang:1.23.3-bookworm AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the application binary
RUN go build -o icmp_api

# Use a lightweight image for the final deployment
FROM debian:bookworm-slim

# Install required system tools for ICMP ping
RUN apt-get update && apt-get install -y \
    iputils-ping \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the application binary from the builder
COPY --from=builder /app/icmp_api /app/

# Expose the port the application listens on
EXPOSE 8080

# Set the environment variable for the authentication key (can be overridden at runtime)
ENV AUTH_KEY=default_key

# Run the application
CMD ["./icmp_api"]
