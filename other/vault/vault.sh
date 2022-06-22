#!/bin/sh
echo $PWD
docker run -d --rm \
  --name server \
  -v ${PWD}/vault.json:/vault/config/vault.json \
  -p "8200:8200" \
  -e VAULT_ADDR=http://0.0.0.0:8200 \
  -e VAULT_API_ADDR=http://0.0.0.0:8200 \
  -e VAULT_ADDRESS=http://0.0.0.0:8200 \
  --cap-add "IPC_LOCK" \
  --network=$1 \
  vault:1.10.4 \
  vault server -config=/vault/config/vault.json
