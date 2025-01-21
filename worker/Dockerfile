FROM golang:1.23-alpine

WORKDIR /app

COPY . .

WORKDIR /app/worker

RUN go build -o worker .

WORKDIR /app/worker

CMD ["./worker"]
