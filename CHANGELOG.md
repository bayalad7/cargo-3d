# Changelog

Todos los cambios relevantes del proyecto. Formato basado en
[Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/); versiones con [SemVer](https://semver.org/lang/es/).
Cada versión liberada en `main` lleva su etiqueta de git (ADR-0012).

## [Sin publicar]

### Agregado
- Estructura del monorepo (Go workspace, pnpm workspaces, Turborepo).
- Backend en Go 1.27 (motor y API en un binario, ADR-0002): servidor con `GET /health`, apagado ordenado y pruebas.
- Docker: `docker-compose.yml` (PostgreSQL 19, API con perfil), `backend/Dockerfile` multi-etapa
  (`golang:1.27` → `distroless/static-debian13`), `.env.example`.
- Contratos en `contracts/` (OpenAPI) y `engine/fixtures/`.
- `CLAUDE.md` con el contexto y las reglas del proyecto, incluidas las reglas para limpiar Docker después de probar.
- Documentación (mdBook): introducción, convenciones, conceptos, ADR 0001–0013,
  guía de Go comparada con PHP, plan de trabajo y bitácora.
- ADR-0007: dónde vive el código de PostgreSQL (`backend/internal/store`).
- ADR-0011: versiones de todas las tecnologías (Node.js 26 LTS, pnpm 12, Turborepo 2.11, …).
- ADR-0012: todo se trabaja directo en `develop` (sin ramas de trabajo ni pull requests); `main` solo recibe versiones.
- Docker Compose con nombre de proyecto fijo (`cargo-3d`: un solo grupo con `db` y `api`), sin reinicio automático,
  `.dockerignore` más estricto y scripts `db:down`, `docker:api`, `docker:clean`.
- `scripts/servicios/levantar` y `detener` (Windows y macOS/Linux) y tareas de VS Code: PostgreSQL en Docker
  y la API con `go run` en un solo paso.
- PostgreSQL con ICU `es-MX` para ordenar texto en español.
- PostgreSQL de Docker en el puerto 5433 de la computadora, para no chocar con un PostgreSQL instalado en el sistema.
- Plan: instrucciones de Docker con un solo comando y un solo grupo `cargo-3d`.
- Soporte de desarrollo en Windows, macOS y Linux: `.gitattributes` (LF para todos), guía de instalación por
  sistema en el plan e imagen Docker para `linux/amd64` y `linux/arm64` (compilación cruzada de Go).
- ADR-0013 y GitHub Actions: `ci.yml` (flujo de ramas, `gofmt`, `go vet`, `go test -race`, build de la imagen)
  y `release.yml` (imagen en ghcr.io y GitHub Release con notas del CHANGELOG al subir `vX.Y.Z`).
- Biblioteca de prompts por sesión para las fases 1 a 3 (`docs/src/prompts.md`), con puntos de control y cierre
  en cada sprint y bloque, publicación de versiones (v0.1.0, v0.2.0, v1.0.0) y prompts de apoyo.
- `docs/datos/` fuera de git (datos de clientes) y fixtures reales anonimizados.
- Skills de Claude Code del proyecto en `.claude/skills/` (iniciar-sesion, sesion, cerrar-sesion, punto-de-control, revisar, explicar, nuevo-adr, cambiar-contrato, cerrar-sprint).
- Guía del flujo de trabajo con Claude Code (`docs/src/flujo-claude-code.md`).
