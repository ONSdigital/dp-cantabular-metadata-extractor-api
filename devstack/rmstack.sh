#!/usr/bin/env bash

NAME="cantabular-metadata-pub-2021"

conts=($(docker ps -a --filter=name="$NAME" --format="{{.Names}}"))
for cont in "${conts[@]}"; do
    docker rm -f $cont
done

images=($(docker images "$NAME*" --format={{.Repository}}))
for image in "${images[@]}"; do
    docker rmi -f $image
done
