# Use Go 1.25 Alpine
FROM golang:1.25-alpine

# Install git
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Enable Go modules
ENV GO111MODULE=on

# Download dependencies
RUN go mod tidy

# Build the server
RUN go build -o server

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]
