#!/usr/bin/env bash
#set -x

ARGS=${*:--id TS009}

payload=$(./makerecp $ARGS)
curl -X POST http://localhost:22300/recipes -d "$payload" -H "accept: application/json" -H "Content-Type: application/json" -H "Authorization: $SERVICE_AUTH_TOKEN"
