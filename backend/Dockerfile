FROM golang:1.23-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY internal/db/migrations ./internal/db/migrations

EXPOSE 8080 8080
CMD [ "/app/main" ]