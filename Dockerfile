# Use an older Go version compatible with the repo
FROM golang:1.16-alpine

# Install git
RUN apk add --no-cache git

# Set working directory inside GOPATH
WORKDIR /go/src/app

# Copy source code into container
COPY . .

# Force GOPATH mode to build the old project
ENV GO111MODULE=off

# Build the Go app
RUN go build -o server

# Expose the port the app listens on
EXPOSE 8080

# Start the server
CMD ["./server"]
