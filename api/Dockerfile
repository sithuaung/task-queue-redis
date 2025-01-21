FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN ls -la

WORKDIR /app/api

RUN go build -o api .

RUN ls -la

FROM alpine:latest

WORKDIR /app/api

COPY --from=builder /app/api/api /app/api/
RUN ls -la /app/api

CMD ["./api"]
