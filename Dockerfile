FROM golang:1.22.0-alpine3.19 AS builder

WORKDIR /usr/local/src/
COPY go.mod go.sum ./
RUN apk update && apk add --no-cache bash postgresql-client
RUN go mod download

COPY . .

RUN go build -o ./bin/weather-bot-app cmd/main.go

FROM alpine

RUN apk update && apk add --no-cache bash postgresql-client
COPY --from=builder /usr/local/src/bin/weather-bot-app /
COPY --from=builder /usr/local/src/wait-for-postgres.sh /

RUN chmod +x /wait-for-postgres.sh

CMD ["./weather-bot-app"]