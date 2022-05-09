#!/bin/bash -eux

pushd dp-cantabular-metadata-extractor-api
  make build
  cp build/dp-cantabular-metadata-extractor-api Dockerfile.concourse ../build
popd
