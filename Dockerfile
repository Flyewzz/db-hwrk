FROM golang:latest AS builder
WORKDIR /usr/src/app

# Cache
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build project
COPY . .
RUN go build cmd/server/main.go

FROM ubuntu:18.10

# Install PostgreSQL
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y postgresql-10

# Open ports
EXPOSE 5000

USER postgres

# Copy app
WORKDIR app
COPY --from=builder /usr/src/app/main .
COPY --from=builder /usr/src/app/init init

# Init db
RUN service postgresql start &&\
    psql --file=init/db_init.sql &&\
    service postgresql stop

# Configure PostgreSQL
RUN sed -i 's/\([^\n]\+\)peer$/\1md5/' /etc/postgresql/10/main/pg_hba.conf
RUN cat init/pg_hba.conf >> /etc/postgresql/10/main/pg_hba.conf
RUN cat init/postgresql.conf >> /etc/postgresql/10/main/postgresql.conf

# Connect PostgreSQL
VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Run DB & app
CMD service postgresql start && ./main
