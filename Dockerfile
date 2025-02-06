FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .env

RUN go build -o /out/app ./cmd/main.go

CMD ["/out/app"]