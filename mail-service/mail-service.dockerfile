FROM golang:1.21.5-alpine as builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o mailApp ./cmd/api

RUN chmod +x /app/mailApp

# ...
FROM alpine:latest

RUN mkdir /app
COPY templates /templates

COPY --from=builder /app /app

CMD ["/app/mailApp"]