FROM golang:alpine as builder
RUN apk add --update musl-dev zeromq-dev gcc git make
WORKDIR /go/src/github.com/CanalTP/gormungandr
ADD $PWD /go/src/github.com/CanalTP/gormungandr
RUN make build

FROM alpine:latest
RUN apk --no-cache add libzmq
USER daemon:daemon
WORKDIR /
COPY --from=builder /go/src/github.com/CanalTP/gormungandr/schedules .
CMD ["./schedules"]
