#!/usr/bin/env bash
# detener.sh — Apaga todo el grupo cargo-3d en macOS y Linux (PostgreSQL y, si estaba, la API en Docker).
# Los contenedores se eliminan, pero los datos se conservan en el volumen cargo-3d_db-data.
#
# Uso: ./scripts/servicios/detener.sh   o doble clic en scripts/servicios/detener.command (macOS)

set -euo pipefail
cd "$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

if ! docker info >/dev/null 2>&1; then
  echo "Docker no está corriendo: no hay nada que detener."
  exit 0
fi

printf '\033[36m==> Deteniendo el grupo cargo-3d...\033[0m\n'
docker compose --profile api down
printf '\033[32mListo. Los datos siguen en el volumen cargo-3d_db-data.\033[0m\n'
