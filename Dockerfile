FROM ubuntu:16.04 as builder
RUN apt update && apt install -y libzmq3-dev gcc git make wget pkg-config
RUN wget https://dl.google.com/go/go1.10.2.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.10.2.linux-amd64.tar.gz
WORKDIR /root/go/src/github.com/CanalTP/gormungandr
ADD $PWD /root/go/src/github.com/CanalTP/gormungandr
RUN PATH="/usr/local/go/bin:/root/go/bin:$PATH" make setup build

FROM ubuntu:16.04
RUN apt update && apt install -y libzmq3-dev curl tzdata
USER daemon:daemon
WORKDIR /
COPY --from=builder /root/go/src/github.com/CanalTP/gormungandr/schedules .
HEALTHCHECK --interval=10s --timeout=3s CMD curl -f http://localhost:8080/status || exit 1
CMD ["./schedules"]
