FROM golang:1.22-alpine

WORKDIR /go/src/github.com/componego/componego

COPY ./../.. .

RUN go build -o /opt/app ./examples/hello-app/cmd/application

ENTRYPOINT ["/opt/app"]
