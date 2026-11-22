# Changelog

Todos los cambios relevantes del proyecto. Formato basado en
[Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/); versiones con [SemVer](https://semver.org/lang/es/).
Cada versión liberada en `main` lleva su etiqueta de git (ADR-0012).

## [Sin publicar]

### Agregado
- Ejemplos del contrato: 8 equipos más (23 en total) y 17 entradas inválidas en `contracts/examples/invalid/`.
- `scripts/contratos/check-examples.mjs` (`pnpm examples:check`): revisa las reglas físicas de los ejemplos.
- Conceptos: sistema de coordenadas, y equipos y normas (Tablas A-2 y B-2 de la NOM-012 leídas del DOF).
- Formatos para levantar datos en el patio (`docs/plantillas/patio/`) y sesión 2.1.6 (equipos con varios tramos).
- Registro de pendientes `docs/src/pendientes.md` (id `P-NNN`, quién, cuándo y cómo retomarlo): `/cerrar-sesion`
  lo actualiza en cada cierre y `/iniciar-sesion` y `/sesion` lo consultan.
- Contrato del motor `contracts/engine.schema.json` v0.1.0 (sprint 1.1): entrada y salida con milímetros y
  kilogramos enteros, tramos del equipo, puerta opcional, orientaciones, motivos de piezas no colocadas y margen de
  carga sobredimensionada (OOG). 15 ejemplos sintéticos por tipo de equipo en `contracts/examples/`, con la lista
  de lo que falta por cubrir.
- ADR-0015 (geometría del contrato), ADR-0016 (qué rechaza el motor y qué reporta) y ADR-0017 (carga
  sobredimensionada y open top).
- Sistema de diseño Voxa (ADR-0014) en `packages/ui` (`@cargo-3d/ui`): tokens en claro y oscuro con
  `light-dark()` y Tailwind CSS v4 (`tokens.css`, fuente única de verdad), componentes `vx-` (`components.css`),
  logo Vóxel en SVG (componente `Logo` y 6 archivos), íconos con Lucide (`Icon`, `appIcons`) y componentes
  React de presentación: `Button`, `StatusBadge`, `FilterChip`, `ConfirmDialog`, `Toast`/`Toaster` y
  `BulkActionBar`. Vista previa en `packages/ui/preview/`.
- Documentación del sistema de diseño (`docs/src/ux/sistema-diseno.md`) y reglas de diseño en `CLAUDE.md`.
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
