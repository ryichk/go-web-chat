FROM golang:1.18.1-alpine

RUN apk update && apk add git
RUN apk add --repository https://git.alpinelinux.org/aports/tree/community/bzr?h=3.11-stable

WORKDIR /go/src

ADD . /go/src
