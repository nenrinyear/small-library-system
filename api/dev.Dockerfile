FROM golang:1.26-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations

EXPOSE 8080

CMD ["go", "run", "./cmd/server"]
