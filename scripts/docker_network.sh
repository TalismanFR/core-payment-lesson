#!/bin/bash
if [[ $(docker network ls -f name=$1 -q) ]]; then
  echo "network $1 already created"
else
  docker network create $1;
fi