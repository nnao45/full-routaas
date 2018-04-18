FROM golang:alpine

MAINTAINER nnao45 <n4sekai5y@gmail.com>

RUN apk update && \
    apk upgrade && \
    apk add git

ADD . /app
WORKDIR /app
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo . 

FROM busybox

WORKDIR /root/
COPY --from=0 /app .
CMD ["./app"]
