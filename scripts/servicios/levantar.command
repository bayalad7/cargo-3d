#!/usr/bin/env bash
# Doble clic en Finder (macOS) para ejecutar levantar.sh en una ventana de Terminal.
exec "$(dirname "$0")/levantar.sh" "$@"
