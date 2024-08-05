# Use the official Golang image for building the application
FROM golang:1.22.5 AS builder

# Set the working directory
WORKDIR /app

# Install system dependencies for protoc
RUN apt-get update && apt-get install -y unzip wget

# Install protoc
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v24.3/protoc-24.3-linux-x86_64.zip && \
    unzip protoc-24.3-linux-x86_64.zip -d /usr/local && \
    rm protoc-24.3-linux-x86_64.zip

# Install protoc-gen-go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Add Go binaries to PATH
ENV PATH="/root/go/bin:/usr/local/bin:${PATH}"

# Set environment variables for cross-compilation
# Adjust these to match the target architecture, e.g., 'linux/arm64' or 'linux/amd64'
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application using the Makefile to handle dependencies and proto generation
RUN make proto
# Cross-compile the binary for Linux
RUN go build -o bin/api-gateway ./cmd/...

# Use a smaller base image for the production stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary and config file from the builder stage
COPY --from=builder /app/bin/api-gateway .
COPY --from=builder /app/internal/config/config.json ./internal/config/config.json

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["./api-gateway"]
