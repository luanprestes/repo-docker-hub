#!/bin/sh
set -e

cd "$(dirname "$0")"

echo "[deploy] Baixando nova imagem..."
docker compose pull api

echo "[deploy] Reiniciando container..."
docker compose up -d --force-recreate api

echo "[deploy] Concluído!"
