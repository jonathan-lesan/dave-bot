FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -o ./main

CMD [ "./main" ]