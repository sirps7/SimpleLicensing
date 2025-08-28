# Dockerfile
FROM golang:1.25-alpine

RUN apk add --no-cache git

WORKDIR /go/src/github.com/sirps7/SimpleLicensing
COPY . .

RUN go mod tidy
RUN go build -o server

EXPOSE 8080

CMD ["./server"]
