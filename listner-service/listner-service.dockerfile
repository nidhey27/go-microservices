FROM golang:1.21.5-alpine as builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o listnerApp .

RUN chmod +x /app/listnerApp

# ...
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/listnerApp /app

CMD ["/app/listnerApp"]