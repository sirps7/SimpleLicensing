# Use Go 1.25 Alpine (latest stable)
FROM golang:1.25-alpine

# Install git
RUN apk add --no-cache git

# Set workdir
WORKDIR /app

# Copy the source code
COPY . .

# Enable Go modules (default in 1.25)
ENV GO111MODULE=on

# Fetch all dependencies
RUN go mod tidy

# Build the server
RUN go build -o server

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]
