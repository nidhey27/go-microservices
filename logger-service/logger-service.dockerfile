FROM golang:1.21.5-alpine as builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o loggerApp ./cmd/api

RUN chmod +x /app/loggerApp

# ...
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/loggerApp /app

CMD ["/app/loggerApp"]