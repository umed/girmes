FROM golang:alpine

COPY go.mod /app/src/
COPY go.sum /app/src/

WORKDIR /app/src

RUN go mod download
