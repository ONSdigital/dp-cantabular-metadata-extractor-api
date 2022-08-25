#!/usr/bin/env bash

#set -x

url=http://localhost:8082/login
token=$(curl -s -d "{\"email\":\"florence@magicroundabout.ons.gov.uk\",\"password\":\"$FLORENCE_WEB_PW\"}" $url)

payload='{
  "alias": "Testing for metadata demo v3",
  "cantabular_blob": "dp_synth_dataset",
  "format": "cantabular_table",
  "id": "la2e031b-3064-427d-8fed-4b35c99bf1a0",
  "output_instances": [
  {

    "code_lists": [
    {
      "href": "http://localhost:22400/code-lists/ltla",
      "id": "ltla",
      "is_hierarchy": false,
      "name": "ltla",
      "is_cantabular_geography": true,
      "is_cantabular_default_geography": true
    },
    {
      "href": "http://localhost:22400/code-lists/sex",
      "id": "sex",
      "is_hierarchy": false,
      "name": "sex",
      "is_cantabular_geography": false,
      "is_cantabular_default_geography": false
    }
    ],
    "dataset_id": "TS009",
    "editions": [
    "2021"
      ],
      "title": "Testing for metadata demo v3"
    }
  ]
}'

curl -X POST http://localhost:22300/recipes -d "$payload" -H "accept: application/json" -H "Content-Type: application/json" -H "X-Florence-Token: $token"
