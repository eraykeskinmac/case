FROM golang:1.23-alpine

RUN apk add --no-cache tini

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

EXPOSE 3000

ENTRYPOINT ["/sbin/tini", "--", "./main"]
