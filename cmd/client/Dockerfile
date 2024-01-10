FROM golang:alpine

WORKDIR /app

COPY client.go .

RUN go build -o client client.go

CMD ["./client"]