FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y gcc libc6-dev
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o scheduler main.go

FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/scheduler /app/scheduler
COPY --from=builder /app/web /app/web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

EXPOSE 7540

CMD ["/app/scheduler"]