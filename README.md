# Cargo 3D

Planeación de carga 3D para contenedores marítimos, cajas secas y plataformas,
pensada para el mercado mexicano.

> **Estado:** fase 1 — prueba de concepto, sprint 1.1. Servidor con `/health`, contrato del motor v0.1.0 y
> catálogo de equipos listos; el motor en Go empieza en la sesión 1.1.2. Mientras no hay cargas reales se trabaja
> con casos de uso simulados ([ADR-0020](docs/src/adr/0020-casos-simulados.md)).

## Arquitectura en una frase

Un **backend en Go** (API y motor de acomodo en el mismo binario) con **PostgreSQL**,
y **pantallas en TypeScript** para web, escritorio y móvil, conectadas por un
**contrato OpenAPI**.

```
   apps/web (React + R3F)   apps/desktop (Wails)   apps/mobile (Expo)
            │                       │                      │
            └───────────────┬───────┴──────────────────────┘
                            │  HTTP + JSON (contracts/openapi.yaml)
                            ▼
              ┌──────── backend/ (Go, un binario) ────────┐
              │  internal/httpapi  →  internal/engine     │
              │          └────────→  internal/store       │
              └─────────────────────┬─────────────────────┘
                                    │
                               PostgreSQL
                     (Docker en desarrollo y producción)
```

## Arranque rápido

```bash
cp .env.example .env                          # cambia la contraseña
docker compose --profile api up -d --build    # grupo cargo-3d: PostgreSQL (5433) + API (8080/health)
docker compose --profile api down             # apagar todo
go test ./backend/...
```

Para programar el backend: `Ctrl+Shift+B` en VS Code o `scripts/servicios/levantar` (`.ps1`/`.cmd` en Windows,
`.sh`/`.command` en macOS y Linux). Levanta PostgreSQL en Docker y la API con `go run`.

## Contrato, equipos y casos de uso

| Qué | Dónde |
|---|---|
| Entrada y salida del motor (JSON Schema) | [`contracts/engine.schema.json`](contracts/engine.schema.json) |
| API (OpenAPI) | [`contracts/openapi.yaml`](contracts/openapi.yaml) |
| Catálogo de equipos (26, con fuente, confianza y su ejemplo) | [`contracts/equipment.json`](contracts/equipment.json) · [página con fotos](docs/src/catalogos/equipos.md) |
| Ejemplos por equipo, cada uno con su dibujo `equipo.svg` | [`contracts/examples/`](contracts/examples/README.md) |
| Casos de uso simulados (marítimo, carretera, restricciones, límites, operación) | [`docs/src/casos/`](docs/src/casos/README.md) |
| Equipos personalizados por empresa | [ADR-0019](docs/src/adr/0019-equipos-personalizados.md) |

```bash
pnpm examples:check     # revisa las reglas físicas de los ejemplos
pnpm examples:render    # vuelve a dibujar equipo.svg de cada ejemplo
```

## Ramas

- `develop`: aquí se trabaja todo, con commits directos (sin ramas de trabajo ni pull requests).
- `main`: solo versiones oficiales liberadas, cada una con su etiqueta (`v0.1.0`, …).

El CI (`.github/workflows/ci.yml`) revisa cada push; subir una etiqueta `vX.Y.Z` en `main`
publica la versión (`.github/workflows/release.yml`). Detalles en el [ADR-0012](docs/src/adr/0012-ramas-y-versiones.md)
y el [ADR-0013](docs/src/adr/0013-ci-y-release.md).

## Documentación

Toda la documentación vive en [`docs/`](docs/src/SUMMARY.md): decisiones (ADR), conceptos, catálogo de equipos,
casos de uso, guía de Go para quien viene de PHP, plan de trabajo, [pendientes](docs/src/pendientes.md) y
[bitácora](docs/src/bitacora/README.md) (un archivo por tema, con índice por fecha y área).
Para verla como sitio web: `mdbook serve docs --open` (requiere mdBook: `winget install Rustlang.mdBook` en Windows, `brew install mdbook` en macOS).

## Requisitos

Funciona en Windows, macOS (incluida Apple Silicon) y Linux. Guía de instalación por sistema en el
[instalación del entorno](docs/src/desarrollo/entorno/instalacion.md); para ver la base de datos con pgAdmin, [Base de datos](docs/src/desarrollo/entorno/base-de-datos.md).

- Go 1.27 o superior
- Docker con Docker Compose
- Node.js 26 LTS y pnpm 12 (scripts de `contracts/` y, desde el sprint 1.6, las pantallas)

Versiones de todas las tecnologías: [ADR-0011](docs/src/adr/0011-versiones.md).

## Convenciones

Código en inglés, comentarios y documentación en español.
Si usas Docker para probar algo, elimina al terminar todo lo que creaste (ver "Reglas de Docker" en [CLAUDE.md](CLAUDE.md)).
Detalles en [convenciones](docs/src/desarrollo/convenciones.md) y [ramas y commits](docs/src/desarrollo/flujo/ramas-y-commits.md); todo lo de desarrollo, en [docs/src/desarrollo/](docs/src/desarrollo/README.md).
