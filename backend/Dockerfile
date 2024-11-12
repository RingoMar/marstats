FROM golang:1.16

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go mod tidy

RUN go mod download go.mongodb.org/mongo-driver
RUN go build .

EXPOSE 5284

CMD ["./mongoServer"]


# docker build -t mongoserver .
# docker run -p 5284:5284 --name mongoserver mongoserver
