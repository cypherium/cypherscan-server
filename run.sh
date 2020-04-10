#!/usr/bin/env bash
docker rm -f syncer query
docker run \
 -e "EXECUTION_TIMEOUT=0" \
 -e "NODES_URLS=ws://40.117.112.213:8546" \
 -e "DYNAMODB_REGION=us-east-2" \
 -e "AWS_ACCESS_KEY_ID=AKIAJYWTBXV3Z2HWLE3Q" \
 -e "AWS_SECRET_ACCESS_KEY=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" \
 -e "REGION=us-east-2" \
 -e "RECENT_TTL_DURATION_IN_SECONDS=36000000" \
 --name syncer \
 -d scan

docker run \
 -e "EXECUTION_TIMEOUT=0" \
 -e "NODES_URLS=ws://40.117.112.213:8546" \
 -e "DYNAMODB_REGION=us-east-2" \
 -e "AWS_ACCESS_KEY_ID=AKIAJYWTBXV3Z2HWLE3Q" \
 -e "AWS_SECRET_ACCESS_KEY=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" \
 -e "REGION=us-east-2" \
 -e "RECENT_TTL_DURATION_IN_SECONDS=36000000" \
 -p 8000:8000 \
 --name query \
 -d scan