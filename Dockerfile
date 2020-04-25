FROM golang:1.14.2-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main
