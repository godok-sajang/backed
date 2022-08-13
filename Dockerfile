# syntax=docker/dockerfile:1
FROM golang:1.18-alpine
WORKDIR /Users/ms/go/src/echo_sample
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /godok
CMD [ "/godok" ]
