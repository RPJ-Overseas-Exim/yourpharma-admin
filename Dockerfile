FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./tmp/main cmd/main.go

EXPOSE 7000

ENTRYPOINT ["./tmp/main"]
