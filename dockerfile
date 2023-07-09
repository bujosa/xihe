FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download

RUN go build -o main .

# This is optional, but good practice
EXPOSE 8080

CMD ["./main"]