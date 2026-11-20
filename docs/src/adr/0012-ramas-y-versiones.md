# ADR-0012: Ramas `develop` y `main`, y versiones con etiquetas

- **Estado:** Aceptado
- **Fecha:** 2026-11-20

## Contexto

Hasta ahora todo se subía directo a `main`. Cuando haya clientes (empezando por el piloto con
Transportes Villa en la fase 2) se necesita distinguir entre "lo que se está construyendo" y
"lo que ya se entregó", para poder desplegar una versión estable y arreglarla sin arrastrar trabajo a medias.

## Decisión

Un flujo tipo *Git Flow* simplificado:

- **`develop`**: rama de trabajo diario. Siempre debe compilar y pasar las pruebas.
- **Ramas de trabajo** `area/descripcion` (ej. `motor/modelos-basicos`), con las mismas áreas que los commits.
  Salen de `develop` y regresan a `develop` por pull request.
- **`main`**: solo versiones oficiales. Cada merge a `main` se etiqueta `vX.Y.Z` (SemVer) y su sección
  del `CHANGELOG.md` pasa de "Sin publicar" a la versión con fecha.
- **`release/vX.Y.Z`**: sale de `develop` para preparar una versión (CHANGELOG, ajustes finales);
  se une a `main` (con etiqueta) y de regreso a `develop`.
- **`hotfix/descripcion`**: sale de `main` para un arreglo urgente de una versión liberada;
  se une a `main` (nueva etiqueta de parche) y a `develop`.
- `main` sigue siendo la rama por defecto de GitHub.
- Las revisiones del flujo (qué ramas pueden entrar a `main`, pruebas, publicación) viven en
  `.github/workflows/ci.yml` y `.github/workflows/release.yml` (ADR-0013).

Antes de la primera versión (`v0.1.0`, al cerrar la fase 1) `main` se queda en los cimientos.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **`develop` + `main` (elegida)** | `main` siempre es lo entregado; fácil saber qué tiene el cliente | Un paso más (release) para publicar |
| Solo `main` (*trunk-based*) | Más simple | Lo entregado y lo en construcción se mezclan; requiere *feature flags* y CI maduro |
| Git Flow completo | Muy estructurado | Demasiada ceremonia para una persona |

## Consecuencias

- **Positivas:** se puede desplegar `main` con confianza y arreglar una versión liberada sin sacar trabajo a medias.
- **Negativas / costos:** hay que acordarse de unir `release/` y `hotfix/` también a `develop`.
- **Qué tendríamos que hacer para revertirla:** unir `develop` a `main`, borrar `develop` y quitar la revisión de ramas de `ci.yml`.
