FROM golang:1.22-alpine

ENV GOROOT /usr/local/go
ENV PIP_ROOT_USER_ACTION=ignore

RUN apk add --update bash git make gcc libc-dev binutils-gold && \
    apk add --update --no-cache python3-dev=3.10.14-r1 --repository=https://dl-cdn.alpinelinux.org/alpine/v3.17/main && \
    ln -sf python3 /usr/bin/python && \
    python3 -m ensurepip && \
    pip3 install --no-cache --upgrade pip setuptools pre-commit

CMD ["sleep", "infinity"]
