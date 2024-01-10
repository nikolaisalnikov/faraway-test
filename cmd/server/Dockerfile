FROM golang:alpine

WORKDIR /app

COPY server.go .

RUN go build -o server server.go

CMD ["./server"]