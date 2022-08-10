#!/usr/bin/env bash

LINES="${1:-8}"

conts=($(docker ps --filter=name="cantabular-metadata-pub" --format="{{.Names}}"))

for cont in "${conts[@]}"; do
    echo
    echo "$cont"
    docker logs $cont | tail "-$LINES"
done
