#!/bin/bash

env GOOS=linux GOARCH=amd64 go build
docker build . -t cbv2
docker run -it cbv2 bash
