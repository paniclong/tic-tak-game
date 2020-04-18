FROM golang:1.14.2-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app

COPY . .

RUN go build -o main .

CMD ["./main"]
