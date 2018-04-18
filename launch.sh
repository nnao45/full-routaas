#!/bin/bash

./get-fullroute-mib.sh
docker build -t nnao45/full-routaas . 
docker run -it --rm --privileged -p 179:179 nnao45/full-routaas:latest
