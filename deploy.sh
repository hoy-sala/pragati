#!/bin/bash
set -e
cd "$(dirname "$0")"
git pull origin master
docker compose build --no-cache web
docker compose up -d
