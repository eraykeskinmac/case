FROM golang:1.23-alpine

RUN apk add --no-cache tini

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

EXPOSE 3000

ENTRYPOINT ["/sbin/tini", "--", "./main"]
