FROM golang:alpine AS builder

MAINTAINER nnao45 <n4sekai5y@gmail.com>

RUN apk update && \
    apk upgrade && \
    apk add git make

ADD . $GOPATH/src/app
#RUN ln -s app $GOPATH/src/app
WORKDIR $GOPATH/src/app
#RUN go get -u github.com/golang/dep/cmd/dep
RUN make dep-install
#RUN dep ensure
RUN make dep
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo . 
RUN make
WORKDIR /app
RUN cp -ar $GOPATH/src/app /app

FROM busybox

WORKDIR /root/
#COPY --from=builder $GOPATH/src/app .
COPY --from=builder app .
CMD ["./app"]
