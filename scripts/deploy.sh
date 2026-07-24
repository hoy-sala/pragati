#!/usr/bin/env bash
set -euo pipefail

APP_DIR="/opt/apps/pragati"
cd "$APP_DIR"

if [ ! -f .env ]; then
  echo "Creating .env from .env.example..."
  cp .env.example .env
  echo ">>> Edit .env with your actual secrets before continuing!"
  exit 1
fi

echo "Pulling latest code..."
git pull origin master

echo "Building and starting services..."
docker compose up -d --build

echo "Waiting for API to be ready..."
for i in $(seq 1 30); do
  if curl -sf http://localhost:5050/health > /dev/null 2>&1; then
    echo "API is ready."
    break
  fi
  sleep 1
done

echo "Running seed (safe — skips if data exists)..."
docker compose exec api /server seed || true

echo ""
echo "Deployment complete!"
echo "  API:  http://<server-ip>:5050"
echo "  Web:  http://<server-ip>:3000"
echo ""
echo "  Admin login: admin@pragati.edu / pragati123"
