FROM golang:alpine as builder
RUN apk add --update musl-dev zeromq-dev gcc git make
WORKDIR /go/src/github.com/CanalTP/gormungandr
ADD $PWD /go/src/github.com/CanalTP/gormungandr
RUN make setup && make

FROM alpine:latest
RUN apk --no-cache add ca-certificates zeromq
WORKDIR /
COPY --from=builder /go/src/github.com/CanalTP/gormungandr/schedules .
CMD ["./schedules"]
