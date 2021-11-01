FROM golang:1.17 AS builder
WORKDIR /go/src/github.com/leki75/multicast-test
COPY . .
RUN go install .


FROM debian:buster
COPY --from=builder /go/bin/multicast-test /usr/bin/multicast-test
USER nobody
ENTRYPOINT ["multicast-test"]
