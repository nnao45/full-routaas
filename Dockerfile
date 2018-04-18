FROM golang:alpine AS build-env

MAINTAINER nnao45 <n4sekai5y@gmail.com>

RUN apk update && \
    apk upgrade && \
    apk add git

ADD . $GOPATH/src/github.com/nnao45/full-routaas
#COPY . /work
WORKDIR $GOPATH/src/github.com/nnao45/full-routaas
RUN ls
#COPY get-fullroute-mib.sh .
#COPY full-routaas.go .
#RUN ./get-fullroute-mib.sh
#RUN go build -o full-routaas full-routaas.go
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go build

FROM busybox
COPY --from=build-env full-routaas /usr/local/bin/full-routaas
ENTRYPOINT ["/usr/local/bin/full-routaas"]
