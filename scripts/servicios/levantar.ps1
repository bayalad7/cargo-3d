# levantar.ps1 — Levanta el entorno de desarrollo en Windows con un solo comando.
#
# Qué hace, en orden:
#   1. Revisa que Docker esté corriendo; si no, abre Docker Desktop y espera a que esté listo.
#   2. Crea .env desde .env.example si todavía no existe.
#   3. Apaga el contenedor de la API si estaba encendido (usa el mismo puerto 8080 que go run).
#   4. Levanta solo PostgreSQL (contenedor cargo-3d-db-1) y espera a que esté "healthy".
#   5. Corre la API con `go run` en esta ventana (Ctrl+C para detenerla).
#   6. Al detener la API, detiene también PostgreSQL (los datos se conservan en el volumen).
#
# Uso (desde cualquier carpeta):
#   .\scripts\servicios\levantar.ps1              # todo lo anterior
#   .\scripts\servicios\levantar.ps1 -DejarDb     # al salir, PostgreSQL sigue encendido
#   o doble clic en scripts\servicios\levantar.cmd
#
# Ver docs/src/desarrollo/entorno/servicios.md (uso) e instalacion.md (8. Verificación).

param(
    # Si se indica, PostgreSQL NO se detiene al cerrar la API.
    [switch]$DejarDb
)

# 'Continue' y no 'Stop': Windows PowerShell 5.1 convierte en error cualquier texto que un programa
# externo (docker, go) escriba en stderr, aunque sea solo informativo. Los fallos reales se revisan
# con $LASTEXITCODE después de cada comando.
$ErrorActionPreference = 'Continue'

# Raíz del repo: dos carpetas arriba de este script (scripts\servicios\ → raíz).
$Root = Resolve-Path (Join-Path $PSScriptRoot '..\..')
Set-Location $Root

function Write-Step([string]$Text) { Write-Host "==> $Text" -ForegroundColor Cyan }
function Stop-WithError([string]$Text) { Write-Host "ERROR: $Text" -ForegroundColor Red; exit 1 }

# Revisa si el motor de Docker responde. `docker info` falla si Docker Desktop está apagado.
function Test-Docker {
    & docker info 2>&1 | Out-Null
    return ($LASTEXITCODE -eq 0)
}

# --- 0. Herramientas instaladas ---
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) { Stop-WithError 'Docker no está instalado (ver docs/src/desarrollo/entorno/instalacion.md).' }
if (-not (Get-Command go -ErrorAction SilentlyContinue)) { Stop-WithError 'Go no está instalado (winget install GoLang.Go).' }

# --- 1. Docker corriendo ---
Write-Step 'Revisando Docker...'
if (-not (Test-Docker)) {
    $dockerDesktop = Join-Path $env:ProgramFiles 'Docker\Docker\Docker Desktop.exe'
    if (-not (Test-Path $dockerDesktop)) { Stop-WithError 'Docker no está corriendo y no encontré Docker Desktop. Ábrelo a mano.' }
    Write-Step 'Abriendo Docker Desktop (puede tardar un minuto)...'
    Start-Process $dockerDesktop
    $deadline = (Get-Date).AddMinutes(3)
    while (-not (Test-Docker)) {
        if ((Get-Date) -gt $deadline) { Stop-WithError 'Docker no respondió en 3 minutos.' }
        Start-Sleep -Seconds 3
    }
}

# --- 2. Archivo .env ---
if (-not (Test-Path '.env')) {
    Copy-Item '.env.example' '.env' -ErrorAction Stop
    Write-Host 'Se creó .env desde .env.example (desarrollo: PostgreSQL sin contraseña, ver ADR-0024).' -ForegroundColor Yellow
}

# Lee .env (líneas CLAVE=valor, sin comentarios) para armar la conexión de go run.
$envVars = @{}
foreach ($line in Get-Content '.env') {
    if ($line -match '^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)\s*$') { $envVars[$Matches[1]] = $Matches[2] }
}
function Get-EnvValue([string]$Name, [string]$Default) {
    if ($envVars.ContainsKey($Name) -and $envVars[$Name] -ne '') { return $envVars[$Name] }
    return $Default
}
$dbUser = Get-EnvValue 'POSTGRES_USER' 'cargo3d'
$dbName = Get-EnvValue 'POSTGRES_DB' 'cargo3d'
$dbPass = Get-EnvValue 'POSTGRES_PASSWORD' ''
$dbPort = Get-EnvValue 'POSTGRES_PORT' '5433'

# --- 3. Contenedor de la API apagado (libera el puerto 8080) ---
& docker compose --profile api rm --stop --force api 2>&1 | Out-Null

# --- 4. PostgreSQL ---
Write-Step "Levantando PostgreSQL (cargo-3d-db-1, puerto $dbPort)..."
& docker compose up -d --wait db
if ($LASTEXITCODE -ne 0) { Stop-WithError "PostgreSQL no arrancó. ¿El puerto $dbPort está ocupado? Cámbialo en POSTGRES_PORT de .env." }

# --- 5. API con go run ---
$env:DATABASE_URL = "postgres://${dbUser}:${dbPass}@localhost:${dbPort}/${dbName}?sslmode=disable"
if (-not $env:PORT) { $env:PORT = '8080' }

Write-Host ''
Write-Host "Listo. API en http://localhost:$($env:PORT)/health  (Ctrl+C para detener)" -ForegroundColor Green
Write-Host ''

try {
    & go run ./backend/cmd/api
}
finally {
    # --- 6. Al salir (Ctrl+C o error) ---
    if ($DejarDb) {
        Write-Step 'PostgreSQL sigue encendido. Para apagarlo: .\scripts\servicios\detener.ps1'
    }
    else {
        Write-Step 'Deteniendo PostgreSQL (los datos se conservan)...'
        & docker compose stop db 2>&1 | Out-Null
    }
}
