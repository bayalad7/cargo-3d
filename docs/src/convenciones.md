# Convenciones

## Idioma

| Qué | Idioma | Ejemplo |
|---|---|---|
| Nombres de variables, funciones, tipos, paquetes | Inglés | `func FindBestPoint(c Container) (Point, bool)` |
| Comentarios en el código | Español | `// Descartamos puntos fuera del contenedor` |
| Documentación (comentarios de Go, TSDoc, `docs/`) | Español | `// FindBestPoint busca el mejor punto disponible…` |
| Mensajes de commit | Español | `motor: agrega modelo Container` |
| Textos que ve el usuario final | Español (i18n después) | `"Contenedor lleno"` |

**Por qué:** el inglés en el código es el estándar de la industria, facilita usar
librerías y contratar después. El español en la documentación hace que todo se entienda sin fricción.

## Nombres en Go

- **Exportado** (visible fuera del paquete) = empieza con mayúscula: `Container`, `Optimize`.
- **Privado** del paquete = empieza con minúscula: `findCandidates`.
- Paquetes en minúsculas, una palabra, sin guiones: `engine`, `httpapi`.
- Siglas en mayúsculas completas: `ID`, `HTTP`, `JSON` (`PlanID`, no `PlanId`).
- Nombres cortos para variables de alcance corto (`c` para un contenedor en un ciclo), descriptivos para lo demás.
- La documentación de un elemento exportado empieza con su nombre: `// Container representa…`.

## Nombres en TypeScript

`camelCase` para funciones y variables, `PascalCase` para tipos y componentes.

## Unidades

En todo el motor las medidas son **milímetros** y los pesos **kilogramos**, como enteros (`int64`)
para evitar errores de redondeo. Cuando haya ambigüedad, la unidad va en el nombre: `LengthMM`, `WeightKG`.

## Errores

- Se devuelven como último valor: `func Optimize(in Input) (LoadPlan, error)`.
- Se revisan siempre: `if err != nil { return ..., err }`.
- Se envuelven con contexto: `fmt.Errorf("leyendo carga %s: %w", name, err)`.
- Los mensajes de error van en español y en minúsculas, sin punto final (convención de Go).

## Glosario (español → código)

| Concepto | En código |
|---|---|
| Contenedor / caja seca / plataforma | `Container` |
| Pieza, bulto | `CargoItem` |
| Colocación de una pieza | `Placement` |
| Plan de carga | `LoadPlan` |
| Restricción | `Constraint` |
| Punto extremo | `ExtremePoint` |
| Apilable | `Stackable` |
| Destino / parada | `Stop` |

## Nivel de comentarios en Go

- **Etapa 1 (actual):** casi cada línea explicada; cada concepto de Go nuevo se explica la primera vez
  que aparece y se enlaza a la guía de aprendizaje.
- **Etapa 2:** se explican bloques y decisiones, no cada línea.
- **Etapa 3:** comentarios solo donde el "por qué" no es obvio (nivel profesional estándar).

El cambio de etapa lo decide Benjamín y se registra en `CLAUDE.md` y en la bitácora.

## Ramas

Ver ADR-0012.

| Rama | Para qué |
|---|---|
| `develop` | Todo el trabajo diario, con commits directos y `git push`. Siempre compila y pasa las pruebas |
| `main` | Solo versiones liberadas: se une `develop` y se etiqueta `vX.Y.Z` |

No se crean ramas de trabajo ni pull requests: el proyecto lo trabaja una sola persona.

## Commits

Formato: `área: descripción en presente`.
Áreas: `motor`, `api`, `web`, `escritorio`, `movil`, `contratos`, `docs`, `infra`, `ui`.
