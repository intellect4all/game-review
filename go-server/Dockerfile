# syntax=docker/dockerfile:1

FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify


COPY . .


RUN go build -o /build/game_review_api

EXPOSE 3000

CMD ["/build/game_review_api"]



