#!/usr/bin/env bash
# Doble clic en Finder (macOS) para ejecutar detener.sh en una ventana de Terminal.
exec "$(dirname "$0")/detener.sh" "$@"
