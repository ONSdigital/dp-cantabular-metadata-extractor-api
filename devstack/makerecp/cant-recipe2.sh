#!/usr/bin/env bash
# TODO merge

#set -x

ID=$1

url=http://localhost:8082/login
token=$(curl -s -d "{\"email\":\"florence@magicroundabout.ons.gov.uk\",\"password\":\"$FLORENCE_WEB_PW\"}" $url)

payload=$(go run . -id $ID)
curl -X POST http://localhost:22300/recipes -d "$payload" -H "accept: application/json" -H "Content-Type: application/json" -H "X-Florence-Token: $token"
