#!/usr/bin/env bash

# exit on error
set -e

cd /root

echo "cloning git repo..."
git clone https://github.com/d-rk/checkin-system.git
cd checkin-system

echo "pulling docker images..."
docker compose pull

echo "starting docker containers..."
docker compose up -d
