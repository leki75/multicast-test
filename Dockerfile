FROM golang:1.17

WORKDIR /go/src/github.com/leki75/multicast-test

COPY . .
RUN go install .

ENTRYPOINT ["multicast-test"]
