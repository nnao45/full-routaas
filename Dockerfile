#FROM golang:alpine AS builder
FROM golang:alpine

MAINTAINER nnao45 <n4sekai5y@gmail.com>

RUN apk update && \
    apk upgrade && \
    apk add git


#COPY . go/src/github.com/nnao45/full-routaas/ 
#WORKDIR go/src/github.com/nnao45/full-routaas
ADD . /app
WORKDIR /app
RUN go get -d -v
#RUN ls
#COPY get-fullroute-mib.sh .
#COPY full-routaas.go .
#RUN ./get-fullroute-mib.sh
#RUN go build -o full-routaas full-routaas.go
#RUN go get -u github.com/golang/dep/cmd/dep
#RUN dep ensure

#RUN go build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo . 

#FROM ubuntu:16.04
FROM busybox
#FROM golang:alpine
#COPY --from=builder  go/src/github.com/nnao45/full-routaas /usr/local/bin/full-routaas
#COPY --from=builder  full-routaas /usr/local/bin/*
#RUN chmod 700 /bin
#ADD https://github.com/nnao45/full-routaas/raw/master/full-routaas /bin/full-routaas
#RUN chmod 700 /bin/full-routaas
#RUN ls /bin/
#ENTRYPOINT ["/bin/sh"]
WORKDIR /root/
COPY --from=0 /app .
#COPY --from=0 /work/full-routaas .
#COPY --from=0 /work/rib.20180417.0000 .
CMD ["./app"]
