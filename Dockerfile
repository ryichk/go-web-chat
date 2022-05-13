FROM golang:1.18.1-alpine

RUN apk update && apk add git

WORKDIR /go/src

ADD . /go/src
