# detener.ps1 — Apaga todo el grupo cargo-3d en Windows (PostgreSQL y, si estaba, la API en Docker).
# Los contenedores se eliminan, pero los datos se conservan en el volumen cargo-3d_db-data.
#
# Uso: .\scripts\servicios\detener.ps1   o doble clic en scripts\servicios\detener.cmd

# 'Continue' y no 'Stop': Windows PowerShell 5.1 convierte en error cualquier texto que un programa
# externo (docker, go) escriba en stderr, aunque sea solo informativo. Los fallos reales se revisan
# con $LASTEXITCODE después de cada comando.
$ErrorActionPreference = 'Continue'
Set-Location (Resolve-Path (Join-Path $PSScriptRoot '..\..'))

& docker info 2>&1 | Out-Null
if ($LASTEXITCODE -ne 0) { Write-Host 'Docker no está corriendo: no hay nada que detener.' -ForegroundColor Yellow; exit 0 }

Write-Host '==> Deteniendo el grupo cargo-3d...' -ForegroundColor Cyan
& docker compose --profile api down
Write-Host 'Listo. Los datos siguen en el volumen cargo-3d_db-data.' -ForegroundColor Green
