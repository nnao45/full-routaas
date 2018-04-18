#!/bin/bash

./get-fullroute-mib.sh
docker run -it --rm --privileged -p 179:179 nnao45/full-routaas:latest
