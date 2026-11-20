# ADR-0013: CI y publicación de versiones con GitHub Actions

- **Estado:** Aceptado
- **Fecha:** 2026-11-20

## Contexto

El flujo de ramas (ADR-0012) necesita que algo revise cada cambio y que publicar una versión sea
siempre igual. Esa configuración debe vivir **en el repositorio**, versionada y revisable en el historial,
y no en pantallas de configuración de GitHub que nadie recuerda haber cambiado.

## Decisión

Dos workflows de GitHub Actions en `.github/workflows/`:

| Archivo | Cuándo corre | Qué hace |
|---|---|---|
| `ci.yml` | Cada push a `develop` y `main` | `gofmt`, `go vet`, `go test -race`; construye la imagen Docker (amd64 y arm64) sin publicarla |
| `release.yml` | Al subir una etiqueta `vX.Y.Z` | Comprueba que la etiqueta esté en `main` y que el CHANGELOG tenga `## [X.Y.Z]`; repite las pruebas; publica la imagen multiplataforma (`linux/amd64` y `linux/arm64`) en `ghcr.io/<dueño>/<repo>/api` (`X.Y.Z`, `X.Y`, `latest`); crea el GitHub Release con las notas del CHANGELOG |

- `main` sigue siendo la **rama por defecto** de GitHub; el trabajo diario se hace en `develop`.
- La versión de Go en CI sale de `backend/go.mod` (una sola fuente de verdad, ADR-0011).
- Las acciones se fijan por versión mayor (`@v7`) y se actualizan junto con el ADR-0011.
- Todo corre en las máquinas de GitHub: no ocupa disco local (ver "Reglas de Docker" en `CLAUDE.md`).
- Cuando exista `apps/web` (sprint 1.6), `ci.yml` agrega un job de pnpm (lint, pruebas y build).

### Cómo publicar una versión

```bash
git switch develop && git pull
# En develop: CHANGELOG "## [Sin publicar]" -> "## [0.1.0] - AAAA-MM-DD"; commit y push.
git switch main && git pull
git merge --ff-only develop             # main queda igual que develop
git tag -a v0.1.0 -m "v0.1.0"
git push origin main v0.1.0             # el push de la etiqueta dispara release.yml
git switch develop                      # y se sigue trabajando en develop
```

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **GitHub Actions en el repo (elegida)** | Configuración versionada; ya está en GitHub; gratis para el volumen del proyecto | Atado a GitHub |
| Reglas solo en la configuración de GitHub | Sin archivos | Invisible en el repo; no queda historial |
| Otro servicio de CI (GitLab CI, CircleCI) | Más opciones | Otra cuenta y otra herramienta que aprender |

## Consecuencias

- **Positivas:** cada push se revisa solo; publicar es subir una etiqueta; las notas del release
  salen del CHANGELOG, así que no hay que escribirlas dos veces.
- **Negativas / costos:** un workflow **no puede** impedir un `git push` directo ni bloquear el botón de merge
  por sí solo. Para que un CI en rojo bloquee el merge hay que marcar sus checks como obligatorios
  (*Settings → Rules*); mientras no se haga, el CI avisa pero no bloquea.
- **Qué tendríamos que hacer para revertirla:** borrar `.github/workflows/`.
