# 20 nov 2026 · Docker y entorno de desarrollo

**Sesión:** Fase 0 · **Área:** infra

## Limpieza de Docker

- Regla nueva en `CLAUDE.md` ("Reglas de Docker"): todo lo que se cree en Docker para probar
  (contenedores, imágenes, volúmenes, redes, caché de build) se elimina al terminar, usando un proyecto
  de Compose aparte (`-p cargo3d-test`) para no tocar la base de datos de desarrollo.
  Motivo: el disco duro se llenó rápidamente con restos de Docker.

## Soporte multiplataforma

- `.gitattributes` con `* text=auto eol=lf`: mismo fin de línea en Windows y macOS.
- Dockerfile con `--platform=$BUILDPLATFORM` y `GOOS/GOARCH`: la imagen se compila de forma cruzada para
  `linux/amd64` y `linux/arm64`; CI y release construyen ambas.
- Plan: la preparación del entorno ya no exige WSL2; tabla de instalación para Windows nativo, macOS y Linux.

## Docker: un solo grupo de contenedores

- Objetivo: que al levantar el proyecto todo quede en un solo grupo `cargo-3d` (contenedores `db` y `api`),
  sin duplicados. Se fijó `name: cargo-3d` en Compose y se quitó el reinicio automático.
- Se evaluó meter PostgreSQL y la API en un solo contenedor y se descartó (ADR-0010).
- También: ICU `es-MX` en PostgreSQL, `.dockerignore` más estricto y scripts `db:down`, `docker:api`, `docker:clean`.
- Se probaron las imágenes alpine y se regresó a las normales (`postgres:19`, `golang:1.27`).
- El contenedor `db` nunca había arrancado: un PostgreSQL 18 instalado en Windows ocupa el puerto 5432.
  El puerto del proyecto pasa a 5433 (`.env.example` y `docker-compose.yml`). Se recreó el volumen y se
  verificó PostgreSQL 19.0 con ICU `es-MX` (la ñ se ordena después de la n).
- Plan: la sección de Docker explica el comando único, el grupo `cargo-3d` y cómo no duplicar recursos.

## Script para levantar el entorno

- `scripts/servicios/levantar.ps1` (Windows) y `levantar.sh` (macOS/Linux), con lanzadores de doble clic
  (`.cmd`, `.command`) y tareas de VS Code (`Ctrl+Shift+B`). Abren Docker si hace falta, levantan PostgreSQL,
  corren la API con `go run` y, al salir, apagan PostgreSQL. `detener` apaga todo el grupo.
- Probado en Windows: `/health` responde y PostgreSQL se detiene al cerrar la API. Falta probar en la MacBook.
