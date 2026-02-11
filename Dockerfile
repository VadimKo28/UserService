FROM golang:1.25 AS builder

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /app/bin/app ./cmd
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /app/bin/migrate ./cmd/migrate/main.go

FROM alpine:3.20

RUN adduser -D -g '' app

WORKDIR /app

COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/bin/migrate /app/migrate
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migrations /app/migrations

USER app

EXPOSE 8081

CMD ["/app/app"]
