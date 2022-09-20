#!/bin/bash

# XXX this breaks macOS but we need it on linux!
if ! [[ $OSTYPE =~ ^darwin.* ]]; then
    cd cantabular-import/minio && sudo chown -R 1001:1001 data
fi
