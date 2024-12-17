#!/bin/bash -eux

pushd dp-cantabular-metadata-extractor-api
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.1
  make lint
popd