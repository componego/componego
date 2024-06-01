FROM golang:1.22-alpine

ENV URL_SHORTENER_PORT=8080

ENV URL_SHORTENER_DB_USERNAME=secret_db_name
ENV URL_SHORTENER_DB_PASSWORD=secret_db_password

WORKDIR /go/src/github.com/componego/componego

COPY ./../.. .

WORKDIR /go/src/github.com/componego/componego/examples/url-shortener-app

# A simple example of how to prepare a configuration for different environments.
COPY ./config/config.json.example /opt/config/production.config.json

RUN go build -o /opt/app ./cmd/application

WORKDIR /opt

EXPOSE $URL_SHORTENER_PORT

ENTRYPOINT ["/opt/app"]
