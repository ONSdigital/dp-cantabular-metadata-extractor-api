#!/bin/bash

if [[ $OSTYPE =~ ^darwin.* ]]; then
    cd cantabular-import/minio && chmod 777 . && chmod 777 data
else 
    # linux
    cd cantabular-import/minio && sudo chown -R 1001:1001 data
fi
