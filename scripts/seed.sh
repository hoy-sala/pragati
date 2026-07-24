#!/usr/bin/env bash
# Run the Go seed command inside the api container
set -e

docker compose exec api /server seed
