# 20 nov 2026 · Ramas, CI y flujo en develop

**Sesión:** Fase 0 · **Área:** infra

## Ramas y revisión general

- Se creó la rama `develop` para el trabajo diario; `main` queda para versiones liberadas (ADR-0012).
- Revisión de README, CHANGELOG, `docker-compose.yml`, plan y convenciones: diagrama con `internal/store`,
  sección de ramas, aviso de limpieza de Docker y CHANGELOG consolidado (aún no hay versiones liberadas).

## CI y publicación

- `main` se queda como rama por defecto de GitHub; la configuración vive en `.github/workflows/` (ADR-0013).
- `ci.yml`: a `main` solo entran `release/vX.Y.Z` o `hotfix/*`; `gofmt`, `go vet`, `go test -race` y build de la imagen.
- `release.yml`: con la etiqueta `vX.Y.Z` en `main`, publica la imagen en ghcr.io y crea el GitHub Release.
- Pendiente opcional: marcar los checks del CI como obligatorios en Settings → Rules (un workflow no puede hacerlo).

## Flujo simplificado: todo en develop

- Solo una persona trabaja en el proyecto: todo se hace con commits directos en `develop`, sin ramas de trabajo
  ni pull requests. `main` solo recibe versiones (`git merge --ff-only develop` y etiqueta `vX.Y.Z`).
- Se integró `infra/script-levantar` en `develop` (avance rápido) y se borró la rama; el PR #1 quedó como fusionado.
- Actualizados: `CLAUDE.md` (sección "Ramas"), ADR-0012, ADR-0013, `ci.yml` (sin revisión de ramas ni
  `pull_request`), convenciones, plan, README, `prompts.md` (sin campo "Rama"; sesión 1.7.1 ajustada al CI existente;
  Node 26 y pnpm 12 en 1.4.1), `flujo-claude-code.md` y las skills (`/revisar` y `/cerrar-sesion` comparan contra
  `origin/develop`).
- `CLAUDE.md` ya no menciona nombres; las skills `/explicar` y `/punto-de-control` dicen "desarrollador".
