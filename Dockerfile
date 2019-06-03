
FROM golang:1.12.0-alpine3.9

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN mkdir /schedule-api

ADD . /schedule-api

WORKDIR /schedule-api

RUN go build -o main .

CMD ["/schedule-api/main"]