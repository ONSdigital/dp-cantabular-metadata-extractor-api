#!/bin/bash
# https://rmoff.net/2019/11/29/using-tcpdump-with-docker/
#docker run --tty --net=container:cantabular-import-journey_florence_1 tcpdump tcpdump -N -A 'port 8083'

docker build -t tcpdump - <<EOF 
FROM ubuntu 
RUN apt-get update && apt-get install -y tcpdump 
CMD tcpdump -i eth0 
EOF
