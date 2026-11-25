# ADR-0012: Ramas `develop` y `main`, y versiones con etiquetas

- **Estado:** Aceptado
- **Fecha:** 2026-11-20

## Contexto

Hasta ahora todo se subía directo a `main`. Cuando haya clientes (empezando por el piloto con
la empresa piloto en la fase 2) se necesita distinguir entre "lo que se está construyendo" y
"lo que ya se entregó", para poder desplegar una versión estable y arreglarla sin arrastrar trabajo a medias.

## Decisión

El proyecto lo trabaja una sola persona, así que se usan **solo dos ramas**:

- **`develop`**: aquí se hace **todo** el trabajo, con commits directos y `git push`. Siempre debe compilar
  y pasar las pruebas (el CI lo revisa en cada push). No hay ramas de trabajo ni pull requests.
- **`main`**: solo versiones oficiales. Para liberar, se une `develop` a `main` (avance rápido, sin commits
  propios en `main`) y se etiqueta `vX.Y.Z` (SemVer); su sección del `CHANGELOG.md` pasa de "Sin publicar"
  a la versión con fecha.
- **Arreglo urgente de una versión liberada:** se corrige en `develop` y se libera una versión de parche
  con el mismo procedimiento.
- `main` sigue siendo la rama por defecto de GitHub. Las pruebas y la publicación viven en
  `.github/workflows/ci.yml` y `.github/workflows/release.yml` (ADR-0013).

Antes de la primera versión (`v0.1.0`, al cerrar la fase 1) `main` se queda en los cimientos.

> **Revisión 2026-11-20:** la primera versión de este ADR usaba ramas de trabajo (`area/descripcion`),
> `release/*`, `hotfix/*` y pull requests. Se simplificó porque solo una persona trabaja en el proyecto.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **`develop` + `main`, commits directos (elegida)** | `main` siempre es lo entregado; sin ceremonia para una persona | Sin revisión por pull request: el CI revisa después del push |
| Solo `main` (*trunk-based*) | Más simple | Lo entregado y lo en construcción se mezclan; requiere *feature flags* y CI maduro |
| `develop` + ramas de trabajo y pull requests | Cada cambio revisado antes de entrar | Demasiada ceremonia para una persona |

## Consecuencias

- **Positivas:** se puede desplegar `main` con confianza y arreglar una versión liberada sin sacar trabajo a medias.
- **Negativas / costos:** un error llega a `develop` antes de que el CI lo detecte; se corrige con otro commit.
- **Qué tendríamos que hacer para revertirla:** unir `develop` a `main` y borrar `develop`.
