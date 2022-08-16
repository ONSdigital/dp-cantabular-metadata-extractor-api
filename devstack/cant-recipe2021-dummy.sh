#!/usr/bin/env bash
#set -x
 
# This is RM003 ("Dataset_Mnemonic_2011": "LC4402EW") but with "ladcd" rather than "oa"
# since "oa" doesn't work!

url=http://localhost:8082/login
token=$(curl -s -d "{\"email\":\"florence@magicroundabout.ons.gov.uk\",\"password\":\"$FLORENCE_WEB_PW\"}" $url)

payload='{
  "alias": "Testing for metadata demo v3",
  "cantabular_blob": "dummy_data_households",
  "format": "cantabular_table",
  "id": "la2e031b-3064-427d-8fed-4b35c99bf1a0",
  "output_instances": [
    {
       
        "code_lists": [
           {
            "href": "http://localhost:22400/code-lists/ladcd",
            "id": "ladcd",
            "is_hierarchy": false,
            "name": "Region",
            "is_cantabular_geography": true,
            "is_cantabular_default_geography": true
          },

          {
            "href": "http://localhost:22400/code-lists/accommodation_type_5a",
            "id": "accommodation_type_5a",
            "is_hierarchy": false,
            "name": "Accommodation type (5 categories)",
            "is_cantabular_geography": false,
             "is_cantabular_default_geography": false
          },
          {
            "href": "http://localhost:22400/code-lists/heating_type_3a",
            "id": "heating_type_3a",
            "is_hierarchy": false,
            "name": "Type of central heating in household (3 categories)",
            "is_cantabular_geography": false,
             "is_cantabular_default_geography": false
          },
          {
            "href": "http://localhost:22400/code-lists/hh_tenure_5a",
            "id": "hh_tenure_5a",
            "is_hierarchy": false,
            "name": "Tenure of household (5 categories)",
            "is_cantabular_geography": false,
            "is_cantabular_default_geography": false
          }
      ],
      "dataset_id": "RM003",
      "editions": [
        "2021"
      ],
      "title": "Testing for metadata demo v3"
    }
  ]
}'

curl -X POST http://localhost:22300/recipes -d "$payload" -H "accept: application/json" -H "Content-Type: application/json" -H "X-Florence-Token: $token"
