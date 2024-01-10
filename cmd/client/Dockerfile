FROM golang:alpine

WORKDIR /app

COPY ../../go.mod .

RUN go mod download

# Копируем все остальное
COPY . .

RUN go build -o client cmd/client/client.go

CMD ["./client"]