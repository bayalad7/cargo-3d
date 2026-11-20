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
- Documentación (mdBook): introducción, convenciones, conceptos, ADR 0001–0012,
  guía de Go comparada con PHP, plan de trabajo y bitácora.
- ADR-0007: dónde vive el código de PostgreSQL (`backend/internal/store`).
- ADR-0011: versiones de todas las tecnologías (Node.js 26 LTS, pnpm 12, Turborepo 2.11, …).
- ADR-0012: ramas `develop` (trabajo diario) y `main` (versiones liberadas).
