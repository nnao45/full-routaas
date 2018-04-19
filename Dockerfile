FROM golang:alpine AS builder

MAINTAINER nnao45 <n4sekai5y@gmail.com>

RUN apk update && \
    apk upgrade && \
    apk add git make

ADD . $GOPATH/src/app
WORKDIR $GOPATH/src/app

RUN make dep-install
RUN make dep
RUN make

RUN mkdir /app && \
    cp -a $GOPATH/src/app/full-routaas /app && \
    cp -a $GOPATH/src/app/config.tml   /app && \
    cp -a $GOPATH/src/app/rib* /app

FROM busybox

WORKDIR /root/
COPY --from=builder app .
CMD ["./full-routaas"]
