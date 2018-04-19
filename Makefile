GO15VENDOREXPERIMENT=1

NAME	 := full-routaas
TARGET	 := $(NAME)
VERSION  := 1.0.0
DIST_DIRS := find * -type d -exec

SRCS	:= $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.version=$(VERSION)\" -extldflags \"-static\""

$(TARGET): $(SRCS)
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o $(NAME)

.PHONY: mibupdate clean upde dep dep-install build run launch

mibupdate:
	./get-fullroute-mib.sh 

clean:
	rm -f full-routaas
	rm -f rib*

upde:
	dep ensure -update

dep:
	dep ensure

dep-install:
	go get -u github.com/golang/dep/cmd/dep
	
build:
	docker build -t nnao45/full-routaas .

run:
	docker run -it --rm --privileged -p 179:179 nnao45/full-routaas:latest

launch:
	./launch.sh
