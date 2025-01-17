FROM golang:1.23.3-alpine3.20

COPY /worker /worker
COPY /mapreduce /mapreduce

COPY .env /worker/

WORKDIR /mapreduce
RUN go mod download && go mod verify

WORKDIR /worker
RUN go mod download && go mod verify

CMD ["go", "run", "/worker/cmd/main.go"]