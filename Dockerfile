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

WORKDIR /app
RUN cp -ar $GOPATH/src/app /app

FROM busybox

WORKDIR /root/
COPY --from=builder app .
CMD ["./app/full-routaas"]
