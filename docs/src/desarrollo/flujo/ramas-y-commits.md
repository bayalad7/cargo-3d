# Ramas, commits y versiones

Guía práctica del día a día. La decisión y sus razones están en el
[ADR-0012](../../adr/0012-ramas-y-versiones.md); la revisión automática y la publicación, en el
[ADR-0013](../../adr/0013-ci-y-release.md).

## Las dos ramas

| Rama | Para qué | Quién escribe en ella |
|---|---|---|
| `develop` | **Todo** el trabajo diario. Siempre debe compilar y pasar las pruebas | Commits directos + `git push` |
| `main` | Solo versiones oficiales, cada una con su etiqueta `vX.Y.Z` | Nadie directo: solo recibe `develop` al liberar |

No hay ramas de trabajo (`area/descripcion`, `release/*`, `hotfix/*`) ni pull requests: el proyecto lo
trabaja una sola persona. Si un prompt, skill o documento dice lo contrario, manda esta regla (`CLAUDE.md`).

## Ciclo de un cambio

```text
editar ──► revisar en local ──► commit en develop ──► git push ──► CI en GitHub (ci.yml)
```

1. **Revisar en local** antes del commit (lo mismo que corre el CI):

   ```bash
   gofmt -l backend          # debe salir vacío
   go vet ./backend/...
   go test ./backend/...
   pnpm examples:check       # si tocaste contracts/
   mdbook build docs         # si tocaste docs/
   ```

2. **Commit** en `develop`. Claude propone el mensaje; tú haces el commit y el push.
3. **CI:** `.github/workflows/ci.yml` corre en cada push (formato, `go vet`, `go test -race` y la construcción de
   la imagen Docker). Si queda en rojo, se corrige con **otro commit** en `develop`; no se reescribe la historia.

## Mensajes de commit

Formato: **`área: descripción en presente`**, en español y en minúsculas.

| Área | Cuándo |
|---|---|
| `motor` | `backend/internal/engine` |
| `api` | `backend/internal/httpapi`, `backend/cmd/api` |
| `contratos` | `contracts/` y sus scripts |
| `web`, `escritorio`, `movil` | `apps/web`, `apps/desktop`, `apps/mobile` |
| `ui` | `packages/ui` (sistema de diseño Voxa) |
| `docs` | `docs/`, `README.md`, `CHANGELOG.md`, `CLAUDE.md` |
| `infra` | Docker, `scripts/`, `.github/`, configuración del monorepo |

Ejemplos reales del historial:

```text
motor: agrega modelo Container
contratos: aclara el motivo único y quién calcula el peso
docs: capitulo de aprendizaje sobre modulos e imports de Go
```

Reglas:

- **Un tema por commit.** Si el mensaje necesita «y además…», probablemente son dos commits.
- La descripción dice **qué cambia**, no qué archivo se tocó («agrega validación de VGM», no «edita validate.go»).
- Si el cambio viene de un ADR o de un pendiente, se menciona: `docs: quita los puntos de control (ADR-0018)`.
- Nunca se suben `.env`, `docs/datos/` ni archivos con datos reales de clientes.

## Liberar una versión

Se hace al cumplir el criterio de terminado de una fase (por ejemplo, `v0.1.0` al cerrar la fase 1).

1. En `develop`, en `CHANGELOG.md`, cambia `## [Sin publicar]` por `## [X.Y.Z] - AAAA-MM-DD` y deja una sección
   `## [Sin publicar]` vacía arriba. Haz el commit (`docs: cierra la version X.Y.Z en el changelog`) y el push.
2. Une `develop` a `main` sin commits propios en `main` (avance rápido) y etiqueta:

   ```bash
   git switch main && git pull
   git merge --ff-only develop
   git push origin main
   git tag -a vX.Y.Z -m "vX.Y.Z"
   git push origin vX.Y.Z
   git switch develop
   ```

3. La etiqueta dispara `.github/workflows/release.yml`: revisa que la etiqueta esté en `main` y que el
   `CHANGELOG.md` tenga la sección `## [X.Y.Z]`, vuelve a correr las pruebas, publica la imagen de la API en
   GitHub Container Registry y crea el GitHub Release con las notas del changelog.

Números de versión con [SemVer](https://semver.org/lang/es/): `MAYOR.MENOR.PARCHE`. Mientras sea `0.x`, un cambio
que rompe el contrato sube el `MENOR`.

## Arreglo urgente de una versión liberada

No hay `hotfix/*`: se corrige en `develop`, se hace commit y push, y se libera una **versión de parche**
(`v0.1.1`) con el mismo procedimiento.

## Si algo sale mal

| Situación | Qué hacer |
|---|---|
| El CI quedó en rojo | Otro commit en `develop` que lo arregle |
| El último commit **todavía no** tiene push y quieres corregir el mensaje | `git commit --amend` |
| El commit ya tiene push | No reescribas la historia (`--amend`, `rebase`, `push --force`): haz otro commit o `git revert <commit>` |
| `git merge --ff-only` falla al liberar | `main` tiene algo que `develop` no; revisa con `git log develop..main` antes de seguir |
