#!/bin/ksh -e

service=$1

# note some services like zebedee and babbage don't have health checks like the
# go services.  The Train gives us some sort of response although it's mispelt.
 
typeset -A port
port["dp-api-router"]=23200
port["dp-dataset-api"]=22000
port["dp-download-service"]=23600
port["dp-frontend-router"]=20000
port["dp-import-api"]=21800
port["dp-import-cantabular-dataset"]=26100
port["dp-import-cantabular-dimension-options"]=26200
port["dp-publishing-dataset-controller"]=24000
port["dp-recipe-api"]=22300
port["florence"]=8081
port["the-train"]=8084

if [[ $service == "list" ]]; then
  for k in "${!port[@]}"; do
    echo $k
  done
  exit
fi

if [[ $service == "ports" ]]; then
  for k in "${!port[@]}"; do
    echo "$k ${port[$k]}"
  done
  exit
fi

if [[ -z $service ]]; then
  for k in "${!port[@]}"; do
    echo $k:${port[$k]}
    curl -s "http://localhost:${port[$k]}/health"| jq .
  done
else 
  echo $service
    curl -s "http://localhost:${port[$service]}/health"| jq .
  exit
fi
