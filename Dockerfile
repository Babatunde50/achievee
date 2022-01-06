FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# RUN go install github.com/cespare/reflex@latest

RUN go build -o api

EXPOSE 8081

CMD ./api