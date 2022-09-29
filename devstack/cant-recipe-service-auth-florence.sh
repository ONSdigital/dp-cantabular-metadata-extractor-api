#!/usr/bin/env bash
#set -x

ARGS=$*

payload=$(./makerecp $ARGS -setalias)
curl -s -X POST https://publishing.dp.aws.onsdigital.uk/recipes -d "$payload" -H "accept: application/json" -H "Content-Type: application/json" -H "Authorization: $SERVICE_AUTH_TOKEN"
