FROM golang:1.23rc2-alpine3.19

WORKDIR /backend

COPY . .

RUN go mod tidy

RUN go build -o main .

CMD ["./main"]