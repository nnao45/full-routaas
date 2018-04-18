FROM golang:alpine AS builder

MAINTAINER nnao45 <n4sekai5y@gmail.com>

RUN apk update && \
    apk upgrade && \
    apk add git


COPY . go/src/github.com/nnao45/full-routaas/ 
WORKDIR go/src/github.com/nnao45/full-routaas
RUN go get -d -v
#RUN ls
#COPY get-fullroute-mib.sh .
#COPY full-routaas.go .
#RUN ./get-fullroute-mib.sh
#RUN go build -o full-routaas full-routaas.go
#RUN go get -u github.com/golang/dep/cmd/dep
#RUN dep ensure

RUN go build

FROM busybox
#COPY --from=builder  go/src/github.com/nnao45/full-routaas /usr/local/bin/full-routaas
COPY --from=builder  full-routaas /usr/local/bin/*
ENTRYPOINT ["/usr/local/bin/full-routaas"]
