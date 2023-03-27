# syntax=docker/dockerfile:1

FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

CMD ["swag", "init"]

RUN go build -o /build/game_review_api

EXPOSE 3000

CMD ["/build/game_review_api"]



