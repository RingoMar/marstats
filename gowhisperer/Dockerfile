
FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod download

RUN go mod tidy

RUN go build .

CMD ["./gowhisperer"]


