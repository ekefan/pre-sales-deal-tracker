FROM golang:1.23.1-alpine3.19 AS build
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/main .
COPY .env .
COPY internal/db/migrations ./internal/db/migrations
EXPOSE 8080
CMD ["/main"]
