FROM golang:1.10.0-alpine3.7

WORKDIR /protostub

RUN apk add --no-cache make

COPY . /protostub/

RUN make
