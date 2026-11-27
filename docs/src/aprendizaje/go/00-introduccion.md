# Go para este proyecto

Esta guía **no** es un curso genérico de Go. Cada concepto se explica con el código real del
backend (`Container`, `CargoItem`, `LoadPlan`, los handlers de la API) en el momento en que aparece.

## Cómo usar esta guía

1. Cuando el código use un concepto nuevo, el comentario lo marcará con
   `// [aprendizaje: nombre-del-concepto]` y aquí habrá un capítulo con ese nombre.
2. Lee el capítulo, luego vuelve al código y léelo línea por línea.
3. Si algo no queda claro, pide `/explicar <archivo o concepto>`.

## Ideas clave de Go en este proyecto

Un resumen de lo que más se usa en el backend. Cada idea se explica a fondo en su capítulo, con el código real.

| Idea | Cómo se escribe | Qué hay que saber |
|---|---|---|
| Tipos con campos | `type Container struct { ... }` + métodos `func (c Container) Volume() int64` | Los datos se arman con structs y composición (un struct dentro de otro) |
| Copias | `b := a` con un struct | Copia los campos; los campos puntero y slice siguen compartiendo sus datos ([capítulo 02](02-structs-y-tipos.md)) |
| Interfaces | `type Constraint interface { ... }` | Un tipo cumple una interfaz con solo tener sus métodos; no se declara en ningún lado |
| Valor cero | `0`, `""`, `false`, `nil`, o devolver `(valor, ok bool)` | Toda variable nace con un valor definido, pero no siempre usable: escribir en un mapa `nil` o leer un puntero `nil` detiene el programa |
| Errores | `return nil, err` + `if err != nil` | Los errores son valores que se devuelven y se revisan donde ocurren |
| Visibilidad | Mayúscula = exportado, minúscula = no exportado | Está en el nombre y es por **paquete** |
| Paquetes | `package engine` + ruta de import | Un paquete = una carpeta ([capítulo 01b](01b-modulos-e-imports.md)) |
| Dependencias | `go.mod` (versiones) + `go.sum` (huellas) | `go.mod` fija las versiones; `go.sum` comprueba que lo descargado no cambió |
| Colecciones | `[]CargoItem` (slice) y `map[string]int` (mapa) | Tipos separados para listas y diccionarios |
| Limpieza al salir | `defer` | Corre al salir de la **función**, en orden inverso si hay varios; sus argumentos se calculan al escribir el `defer` |
| Concurrencia | `go funcion()` (goroutines) y canales | Incluida en el lenguaje y en la librería estándar |
| Tipos | `var x int64` o `x := 5` (tipo inferido) | Tipado estático: el compilador revisa todo antes de correr |
| Mínimo y máximo | `min(a, b)` y `max(a, b)` | Funciones integradas desde Go 1.21, sin importar nada; las usa `overlap` en `box.go` |
| Condiciones | `a && b`, `a \|\| b` | Se evalúan de izquierda a derecha y se detienen en cuanto se sabe el resultado (cortocircuito), como en `Box.Intersects` |
| Comparar structs | `got != want` | Dos structs se comparan con `==` si todos sus campos son comparables; un struct con un slice o un mapa no compila con `==` |
| Desbordamiento | `int64` cerca de 9,2 × 10¹⁸ | Una suma o un producto que se pasa del máximo da la vuelta sin aviso: ni error ni `panic` ([capítulo 02](02-structs-y-tipos.md); límites de `Box` en [Geometría de cajas](../../conceptos/geometria.md)) |

## Ruta de capítulos (se escriben conforme avanza el backend)

- [x] 01 — Primeros pasos: herramientas, módulos y el servidor `/health`
- [x] 01b — [Módulos e imports](01b-modulos-e-imports.md): por qué `github.com/...` se lee del disco
- [x] 02 — [`struct` y tipos numéricos](02-structs-y-tipos.md) (con `Container` y `CargoItem`)
- [x] 03 — [Métodos y receptores](03-metodos-y-receptores.md) (volumen, orientaciones, plan de carga)
- [x] 04 — [Errores como valores](04-errores-como-valores.md) (validación de la entrada)
- [ ] 05 — Slices, maps y ciclos (recorrer puntos extremos)
- [ ] 06 — Punteros: cuándo `*Container` y cuándo `Container`
- [ ] 07 — Interfaces (restricciones intercambiables)
- [ ] 08 — Pruebas de tabla, benchmarks y fuzzing
- [ ] 09 — Goroutines, canales y `context` (metaheurística en paralelo)
- [ ] 10 — HTTP con la librería estándar y JSON
- [ ] 11 — PostgreSQL con pgx y sqlc

Recursos para profundizar: el *Tour of Go* (go.dev/tour, tiene versión en español) y *Effective Go*.
