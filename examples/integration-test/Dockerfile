FROM golang:1.20.3-alpine
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build

CMD "./integration-test"