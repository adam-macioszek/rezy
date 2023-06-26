FROM golang:1.20-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz


FROM  alpine:3.17
WORKDIR /app
COPY --from=builder  /build/main .
COPY --from=builder /build/migrate ./migrate
COPY db/migration ./migration
COPY start.sh .
COPY app.env .


EXPOSE 8080

CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]