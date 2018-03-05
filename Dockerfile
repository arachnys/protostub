FROM golang:1.9.2-alpine3.7

RUN apk update && apk add make

ADD . /protostub/
WORKDIR /protostub

RUN make
