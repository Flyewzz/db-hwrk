FROM golang:latest
ENV GO111MODULE=on

RUN mkdir -p /app

WORKDIR /app

COPY . /app

RUN go mod download
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build cmd/server/main.go

EXPOSE 5000
