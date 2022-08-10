#!/usr/bin/env bash

NAME="cantabular-metadata-pub"

conts=($(docker ps -a --filter=name="$NAME" --format="{{.Names}}"))
for cont in "${conts[@]}"; do
    docker rm -f $cont
done

images=($(docker images "$NAME*" --format={{.Repository}}))
for image in "${conts[@]}"; do
    docker rmi -f $image
done
