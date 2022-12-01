FROM golang:1.17 as builder

RUN apt update

RUN wget --no-check-certificate https://github.com/apple/foundationdb/releases/download/7.0.0/foundationdb-clients_7.0.0-1_amd64.deb
RUN dpkg -i foundationdb-clients_7.0.0-1_amd64.deb
RUN apt update
RUN apt install -y dnsutils

WORKDIR /fdb-exporter

COPY go.mod /fdb-exporter

RUN go mod download

COPY . /fdb-exporter

RUN go build -o fdb-exporter

FROM ubuntu:20.04

RUN apt update
RUN apt install wget -y

RUN wget --no-check-certificate https://github.com/apple/foundationdb/releases/download/7.0.0/foundationdb-clients_7.0.0-1_amd64.deb
RUN dpkg -i foundationdb-clients_7.0.0-1_amd64.deb
RUN apt update
RUN apt install -y dnsutils

COPY --from=builder /fdb-exporter/fdb-exporter /fdb-exporter

CMD /fdb-exporter