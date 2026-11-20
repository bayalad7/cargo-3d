# Changelog

Todos los cambios relevantes del proyecto. Formato basado en
[Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/); versiones con [SemVer](https://semver.org/lang/es/).

## [Sin publicar]

### Agregado
- Estructura del monorepo (Go workspace, pnpm workspaces, Turborepo).
- Backend en Go (motor y API en un binario, ADR-0002): servidor con `GET /health`, apagado ordenado y pruebas.
- Docker: `docker-compose.yml` (PostgreSQL 17, API con perfil), `backend/Dockerfile` multi-etapa, `.env.example`.
- Contratos en `contracts/` (OpenAPI) y `engine/fixtures/`.
- `CLAUDE.md` con el contexto y las reglas del proyecto.
- Documentación (mdBook): introducción, convenciones, conceptos, ADR 0001–0010,
  guía de Go comparada con PHP, plan de trabajo y bitácora.
