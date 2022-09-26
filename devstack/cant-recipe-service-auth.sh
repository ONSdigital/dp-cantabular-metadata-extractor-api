#!/usr/bin/env bash
#set -x

ARGS=$*

payload=$(./makerecp $ARGS -setalias)
curl -X POST http://localhost:22300/recipes -d "$payload" -H "accept: application/json" -H "Content-Type: application/json" -H "Authorization: $SERVICE_AUTH_TOKEN"
