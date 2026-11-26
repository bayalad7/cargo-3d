#!/usr/bin/env bash
# levantar.sh — Levanta el entorno de desarrollo en macOS (y Linux) con un solo comando.
#
# Qué hace, en orden:
#   1. Revisa que Docker esté corriendo; en macOS abre Docker Desktop u OrbStack y espera a que esté listo.
#   2. Crea .env desde .env.example si todavía no existe.
#   3. Apaga el contenedor de la API si estaba encendido (usa el mismo puerto 8080 que go run).
#   4. Levanta solo PostgreSQL (contenedor cargo-3d-db-1) y espera a que esté "healthy".
#   5. Corre la API con `go run` en esta terminal (Ctrl+C para detenerla).
#   6. Al detener la API, detiene también PostgreSQL (los datos se conservan en el volumen).
#
# Uso (desde cualquier carpeta):
#   ./scripts/servicios/levantar.sh             # todo lo anterior
#   ./scripts/servicios/levantar.sh --dejar-db  # al salir, PostgreSQL sigue encendido
#   o doble clic en scripts/servicios/levantar.command (macOS)
#
# Ver docs/src/desarrollo/entorno/servicios.md (uso) e instalacion.md (8. Verificación).

set -euo pipefail

DEJAR_DB=false
[[ "${1:-}" == "--dejar-db" ]] && DEJAR_DB=true

# Raíz del repo: dos carpetas arriba de este script (scripts/servicios/ → raíz).
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

step() { printf '\033[36m==> %s\033[0m\n' "$1"; }
fail() { printf '\033[31mERROR: %s\033[0m\n' "$1" >&2; exit 1; }
docker_ok() { docker info >/dev/null 2>&1; }

# --- 0. Herramientas instaladas ---
command -v docker >/dev/null || fail "Docker no está instalado (ver docs/src/desarrollo/entorno/instalacion.md)."
command -v go >/dev/null || fail "Go no está instalado (macOS: brew install go)."

# --- 1. Docker corriendo ---
step "Revisando Docker..."
if ! docker_ok; then
  if [[ "$(uname)" == "Darwin" ]]; then
    if [[ -d /Applications/OrbStack.app ]]; then app=OrbStack; else app=Docker; fi
    step "Abriendo $app (puede tardar un minuto)..."
    open -a "$app" || fail "No pude abrir $app. Ábrelo a mano."
  else
    fail "Docker no está corriendo. En Linux: sudo systemctl start docker"
  fi
  for _ in $(seq 1 60); do docker_ok && break; sleep 3; done
  docker_ok || fail "Docker no respondió en 3 minutos."
fi

# --- 2. Archivo .env ---
if [[ ! -f .env ]]; then
  cp .env.example .env
  printf '\033[33mSe creó .env desde .env.example. Cambia POSTGRES_PASSWORD cuando puedas.\033[0m\n'
fi

# Lee una variable de .env (CLAVE=valor), o devuelve el valor por defecto.
env_value() {
  local v
  v="$(grep -E "^[[:space:]]*$1[[:space:]]*=" .env | tail -n1 | cut -d= -f2- | tr -d '\r' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')" || true
  echo "${v:-$2}"
}
DB_USER="$(env_value POSTGRES_USER cargo3d)"
DB_NAME="$(env_value POSTGRES_DB cargo3d)"
DB_PASS="$(env_value POSTGRES_PASSWORD '')"
DB_PORT="$(env_value POSTGRES_PORT 5433)"

# --- 3. Contenedor de la API apagado (libera el puerto 8080) ---
docker compose --profile api rm --stop --force api >/dev/null 2>&1 || true

# --- 4. PostgreSQL ---
step "Levantando PostgreSQL (cargo-3d-db-1, puerto $DB_PORT)..."
docker compose up -d --wait db || fail "PostgreSQL no arrancó. ¿El puerto $DB_PORT está ocupado? Cámbialo en POSTGRES_PORT de .env."

# --- 6. Al salir (Ctrl+C o error) ---
cleanup() {
  if $DEJAR_DB; then
    step "PostgreSQL sigue encendido. Para apagarlo: ./scripts/servicios/detener.sh"
  else
    step "Deteniendo PostgreSQL (los datos se conservan)..."
    docker compose stop db >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

# --- 5. API con go run ---
export DATABASE_URL="postgres://${DB_USER}:${DB_PASS}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable"
export PORT="${PORT:-8080}"

printf '\n\033[32mListo. API en http://localhost:%s/health  (Ctrl+C para detener)\033[0m\n\n' "$PORT"
go run ./backend/cmd/api || true
