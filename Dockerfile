# Build Stage
FROM golang:1.17.1-alpine3.14 AS builder
WORKDIR /app
COPY . .
RUN go build -o app *.go

# Run Stage
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 4000
CMD ["/app/app"]