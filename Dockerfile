FROM golang:1.21.5

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]