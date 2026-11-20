# ADR-0011: Versiones de las tecnologías

- **Estado:** Aceptado
- **Fecha:** 2026-11-20

## Contexto

El proyecto arranca desde cero, así que no hay código heredado que obligue a usar versiones viejas.
Algunas tecnologías ya están fijadas en archivos del repo (Go, PostgreSQL, pnpm, Turborepo); otras
(React, Expo, Wails…) se instalarán en sprints futuros. Conviene dejar escrito qué versión se espera de
cada una para que todas las piezas encajen y nadie instale, por ejemplo, React 18 por accidente.

## Decisión

Usar la **última versión estable** de cada tecnología al momento de instalarla, y para Node.js la
**última LTS** (no la "Current"). Versiones vigentes al 2026-11-20:

| Tecnología | Versión | Dónde se fija | Se instala en |
|---|---|---|---|
| Go | 1.27 (1.27.4) | `backend/go.mod`, `go.work`, `backend/Dockerfile` (`golang:1.27`) | Fase 0 |
| PostgreSQL | 19 (ICU `es-MX`) | `docker-compose.yml` (`postgres:19`) | Fase 0 |
| Imagen final del backend | `distroless/static-debian13:nonroot` | `backend/Dockerfile` | Fase 0 |
| Node.js | 26 LTS | `package.json` (`engines`) | Sprint 1.4 |
| pnpm | 12.21 | `package.json` (`packageManager`) | Sprint 1.4 |
| Turborepo | 2.11 | `package.json` (`devDependencies`) | Sprint 1.4 |
| TypeScript | 7.0 | `package.json` de cada app o paquete | Sprint 1.4 |
| Vite | 8 | `apps/web` | Sprint 1.6 |
| React | 19.3 | `apps/web`, `apps/mobile` | Sprint 1.6 |
| Three.js / React Three Fiber / drei | 0.187 / 9.8 / 10.7 | `packages/viewer-3d` | Sprint 1.6 |
| Expo (React Native 0.88) | SDK 58 | `apps/mobile` | Sprint 3.2 |
| Wails | v3.1 (estable, no las *nightly*) | `apps/desktop` | Sprint 3.5 |
| pgx | v5.11 | `backend/go.mod` | Sprint 2.3 |
| sqlc | 1.31 | herramienta de desarrollo | Sprint 2.3 |
| mdBook | 0.5 | herramienta de desarrollo | Fase 0 |

Reglas:

1. Se fija la **versión mayor** (y la menor cuando la herramienta lo pide); los parches se aceptan solos.
2. Cuando una tecnología de la tabla se instala por primera vez, se usa la estable más nueva de ese día
   y se actualiza esta tabla si cambió.
3. Subir una **versión mayor** después de instalada se hace en su propio pull request, con pruebas
   pasando y una línea en el `CHANGELOG.md`.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Última estable (elegida)** | Más tiempo de soporte; sin migraciones al poco tiempo de empezar | Las versiones `.0` pueden traer errores tempranos |
| Penúltima versión mayor | Más probada | Empezar un proyecto nuevo ya con una actualización pendiente |
| No fijar versiones | Menos mantenimiento | Builds no reproducibles; "en mi máquina sí funciona" |

## Consecuencias

- **Positivas:** todo el equipo (y Claude) sabe qué versión usar; el proyecto arranca con soporte largo.
- **Negativas / costos:** PostgreSQL 19 y TypeScript 7 son recientes; si aparece un error grave, se baja
  a la versión anterior (PostgreSQL 18, TypeScript 6) y se anota aquí.
- **Notas de cambios incompatibles ya resueltos:**
  - Desde PostgreSQL 18, la imagen de Docker guarda los datos en `/var/lib/postgresql/<versión>/docker`;
    el volumen se monta en `/var/lib/postgresql`.
  - Desde Node.js 25, Corepack ya no viene incluido: pnpm se instala con `npm install -g pnpm@12`.
- **Qué tendríamos que hacer para revertirla:** cambiar las versiones en los archivos de la columna
  "Dónde se fija" y actualizar esta tabla.
