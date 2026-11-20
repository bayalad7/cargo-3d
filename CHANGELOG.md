# Changelog

Todos los cambios relevantes del proyecto. Formato basado en
[Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/); versiones con [SemVer](https://semver.org/lang/es/).

## [Sin publicar]

### Cambiado
- Versión mínima de Go: 1.27 (`go.mod`, `go.work`, imagen `golang:1.27` del Dockerfile y documentación).
- PostgreSQL 17 → 19; el volumen de datos se monta en `/var/lib/postgresql` (cambio de la imagen desde la 18).
- Imagen final del backend: `distroless/static-debian13` (antes debian12).
- Node.js 26 LTS, pnpm 12.21 y Turborepo 2.11 (ADR-0011).

### Agregado
- Estructura del monorepo (Go workspace, pnpm workspaces, Turborepo).
- Backend en Go (motor y API en un binario, ADR-0002): servidor con `GET /health`, apagado ordenado y pruebas.
- Docker: `docker-compose.yml` (PostgreSQL 19, API con perfil), `backend/Dockerfile` multi-etapa, `.env.example`.
- Contratos en `contracts/` (OpenAPI) y `engine/fixtures/`.
- `CLAUDE.md` con el contexto y las reglas del proyecto.
- Documentación (mdBook): introducción, convenciones, conceptos, ADR 0001–0011,
  guía de Go comparada con PHP, plan de trabajo y bitácora.
