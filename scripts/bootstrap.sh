#!/bin/bash

set -e -u -x -o pipefail

docker-compose up -d

echo "waiting for dynamodb to start"
while ! aws --endpoint-url http://localhost:8001 dynamodb list-tables; do sleep 1; done

aws --endpoint-url http://localhost:8001 dynamodb create-table \
    --table-name LastSecond-local-UserStack-UserTable \
    --billing-mode PAY_PER_REQUEST \
    --key-schema "[{\"AttributeName\":\"id\",\"KeyType\":\"HASH\"}]" \
    --attribute-definitions "[{\"AttributeName\":\"id\",\"AttributeType\":\"S\"}]"
