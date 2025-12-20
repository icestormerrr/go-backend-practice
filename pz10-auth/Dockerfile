FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o server ./cmd/server

EXPOSE ${APP_PORT}

CMD ["./server"]