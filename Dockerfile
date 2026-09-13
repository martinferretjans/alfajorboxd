FROM golang:1.26

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
