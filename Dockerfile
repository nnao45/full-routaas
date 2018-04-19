FROM golang:alpine

MAINTAINER nnao45 <n4sekai5y@gmail.com>

RUN apk update && \
    apk upgrade && \
    apk add git

ADD . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo . 

FROM busybox

WORKDIR /root/
COPY --from=0 $GOPATH/src/app .
CMD ["./app"]
