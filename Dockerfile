# Use an older Go version compatible with the repo
FROM golang:1.16-alpine

# Install git (needed if repo fetches modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy source code into the container
COPY . .

# Build the Go app
RUN go build -o server

# Expose the port the app listens on
EXPOSE 8080

# Start the server
CMD ["./server"]
