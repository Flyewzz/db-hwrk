FROM golang:1.12-stretch AS builder

WORKDIR /usr/src/app
COPY . .

RUN go mod download & go build cmd/server/main.go

FROM ubuntu:18.10
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y postgresql-10

USER postgres

WORKDIR app
COPY --from=builder /usr/src/app .

RUN service postgresql start &&\
    psql -f init/init.sql &&\
    service postgresql stop

RUN echo "listen_addresses = '*'\nsynchronous_commit = off\nfsync = off" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/10/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

EXPOSE 5000
CMD service postgresql start && ./main
