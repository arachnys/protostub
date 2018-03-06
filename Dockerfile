FROM golang:1.10.0-alpine3.7

WORKDIR /protostub

RUN apk add --no-cache make

COPY . /protostub/

RUN make

FROM alpine:3.7

WORKDIR /protostub

COPY --from=0 /protostub/bin/protostub /usr/local/bin
ENTRYPOINT ["/usr/local/bin/protostub"]
