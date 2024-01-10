FROM golang:alpine

WORKDIR /app

COPY ../../go.mod .

RUN go mod download

# Копируем все остальное
COPY . .

RUN go build -o server cmd/server/server.go

CMD ["./server"]