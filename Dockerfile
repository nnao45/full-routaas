FROM golang:alpine AS build-env

MAINTAINER nnao45 <n4sekai5y@gmail.com>

RUN apk update && \
    apk upgrade && \
    apk add git

ADD . /work
WORKDIR /work
RUN go build -o full-routaas full-routaas.go
RUN get-fullroute-mib.sh

FROM busybox
COPY --from=build-env /work/full-routaas /usr/local/bin/full-routaas
ENTRYPOINT ["/usr/local/bin/full-routaas"]
