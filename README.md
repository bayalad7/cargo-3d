# Cargo 3D

Planeación de carga 3D para contenedores marítimos, cajas secas y plataformas,
pensada para el mercado mexicano.

> **Estado:** fase 0 — cimientos. Servidor con `/health`; el motor se construye en la fase 1.

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
cp .env.example .env          # cambia la contraseña
docker compose up -d db       # PostgreSQL
go run ./backend/cmd/api      # http://localhost:8080/health
go test ./backend/...
```

## Ramas

- `develop`: rama de trabajo diario. Cada cambio sale de aquí en su propia rama y vuelve por pull request.
- `main`: solo versiones oficiales liberadas, cada una con su etiqueta (`v0.1.0`, …).

Detalles en el [ADR-0012](docs/src/adr/0012-ramas-y-versiones.md).

## Documentación

Toda la documentación vive en [`docs/`](docs/src/SUMMARY.md): decisiones (ADR), conceptos,
guía de Go para quien viene de PHP, plan de trabajo y bitácora.
Para verla como sitio web: `mdbook serve docs --open` (requiere `cargo install mdbook` o el binario de mdBook).

## Requisitos

- Go 1.27 o superior
- Docker con Docker Compose
- Node.js 26 LTS y pnpm 12 (para las pantallas, desde la fase 1)

Versiones de todas las tecnologías: [ADR-0011](docs/src/adr/0011-versiones.md).

## Convenciones

Código en inglés, comentarios y documentación en español.
Si usas Docker para probar algo, elimina al terminar todo lo que creaste (ver "Reglas de Docker" en [CLAUDE.md](CLAUDE.md)).
Detalles en [docs/src/convenciones.md](docs/src/convenciones.md).
