#!/bin/bash

tag=$1
docker build -t 192.168.12.40:30002/fpga-cloud/task:$tag -f ./deploy/Dockerfile .
