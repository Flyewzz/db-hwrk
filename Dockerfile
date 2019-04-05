FROM golang:1.12-stretch AS builder

WORKDIR /usr/src/app
COPY . .

RUN go mod download & go build cmd/server/main.go

RUN apt-get update && apt-get install -y postgresql-10

USER postgres

RUN service postgresql start &&\
    psql -f init/init.sql &&\
    service postgresql stop

COPY config/pg_hba.conf /etc/postgresql/10/main/pg_hba.conf
COPY config/postgresql.conf /etc/postgresql/10/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

COPY --from=builder /usr/src/app .
CMD service postgresql start && ./main
