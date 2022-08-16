#!/usr/bin/env bash

LINES="${1:-8}"

conts=($(docker ps --filter=name="cantabular-metadata-pub" --format="{{.Names}}"))

if [[ ! -d logs ]]; then
    mkdir logs
fi

for cont in "${conts[@]}"; do
    echo
    echo "$cont"
    docker logs $cont > logs/$cont.log
done
