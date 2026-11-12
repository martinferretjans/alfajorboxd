FROM golang:1.26

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1100

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
