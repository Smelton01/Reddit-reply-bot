# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY * ./
COPY /bot ./bot
COPY /scrape ./scrape

RUN go build -o streamer-bot

EXPOSE 8080

CMD [ "./streamer-bot" ]