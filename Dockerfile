# Use Go 1.16 Alpine image
FROM golang:1.16-alpine

# Install git
RUN apk add --no-cache git

# Set working directory inside GOPATH
WORKDIR /go/src/app

# Copy source code
COPY . .

# Disable Go modules (old repo)
ENV GO111MODULE=off

# Fetch external dependencies
RUN go get github.com/go-sql-driver/mysql \
    github.com/gorilla/mux \
    github.com/pelletier/go-toml

# Build the server
RUN go build -o server

# Expose the server port
EXPOSE 8080

# Start the server
CMD ["./server"]
