FROM golang:1.12-stretch AS builder

WORKDIR /usr/src/app
COPY . .

#RUN go mod download
RUN go build cmd/server/main.go

FROM ubuntu:18.10
ENV DEBIAN_FRONTEND=noninteractive
EXPOSE 5000

RUN apt-get update && apt-get install -y postgresql-10

#RUN locale-gen ru_RU.UTF-8

USER postgres

RUN service postgresql start &&\
    psql --command "CREATE USER forum WITH SUPERUSER PASSWORD 'forum';" &&\
    createdb -O forum forum &&\
    #createdb -T template0 -l ru_RU.UTF-8 -O forum forum &&\
    service postgresql stop

WORKDIR app
COPY --from=builder /usr/src/app .

RUN echo "listen_addresses = '*'\nsynchronous_commit = off\nfsync = off" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/10/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

CMD service postgresql start && ./main
