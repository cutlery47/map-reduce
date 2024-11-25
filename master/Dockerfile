FROM golang:1.23.3-alpine3.20

COPY /master /master
COPY /mapreduce /mapreduce

WORKDIR /mapreduce
RUN go mod download && go mod verify

WORKDIR /master
RUN go mod download && go mod verify

CMD ["go", "run", "/master/cmd/main.go"]