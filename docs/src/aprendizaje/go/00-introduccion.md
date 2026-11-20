# Go para este proyecto

Esta guía **no** es un curso genérico de Go. Cada concepto se explica con el código real del
backend (`Container`, `CargoItem`, `LoadPlan`, los handlers de la API) en el momento en que aparece.

## Cómo usar esta guía

1. Cuando el código use un concepto nuevo, el comentario lo marcará con
   `// [aprendizaje: nombre-del-concepto]` y aquí habrá un capítulo con ese nombre.
2. Lee el capítulo, luego vuelve al código y léelo línea por línea.
3. En cada punto de control explicas el módulo con tus palabras.

## Comparación rápida con PHP

| En PHP | En Go | Diferencia clave |
|---|---|---|
| `class Container { ... }` | `type Container struct { ... }` + métodos `func (c Container) Volume() int64` | No hay clases ni herencia: structs con métodos |
| `interface Constraint` + `implements` | `type Constraint interface { ... }` | Las interfaces se cumplen solas: no se escribe `implements` |
| `$x = null;` | Valor cero (`0`, `""`, `nil`) o devolver `(valor, ok bool)` | Toda variable nace con un valor útil |
| `throw new Exception()` / `try/catch` | `return nil, err` + `if err != nil` | Los errores son valores que se devuelven y se revisan |
| `public` / `private` | Mayúscula = exportado, minúscula = privado | La visibilidad está en el nombre |
| `namespace App\Engine;` | `package engine` + ruta de import | Un paquete = una carpeta |
| `composer.json` / `composer install` | `go.mod` / `go mod tidy` | Dependencias del módulo |
| Arreglos asociativos | `map[string]int` y slices `[]CargoItem` | Tipos separados para listas y diccionarios |
| `finally` | `defer` | Ejecuta algo al salir de la función |
| Procesos o colas para paralelismo | `go funcion()` (goroutines) y canales | Concurrencia incluida en el lenguaje |
| Variables sin tipo | `var x int64` o `x := 5` (tipo inferido) | Tipado estático; el compilador revisa todo |

## Ruta de capítulos (se escriben conforme avanza el backend)

- [x] 01 — Primeros pasos: herramientas, módulos y el servidor `/health`
- [ ] 02 — `struct` y tipos numéricos (con `Container` y `CargoItem`)
- [ ] 03 — Métodos y receptores (volumen, orientaciones)
- [ ] 04 — Errores como valores (cuando una pieza no cabe)
- [ ] 05 — Slices, maps y ciclos (recorrer puntos extremos)
- [ ] 06 — Punteros: cuándo `*Container` y cuándo `Container`
- [ ] 07 — Interfaces (restricciones intercambiables)
- [ ] 08 — Pruebas de tabla, benchmarks y fuzzing
- [ ] 09 — Goroutines, canales y `context` (metaheurística en paralelo)
- [ ] 10 — HTTP con la librería estándar y JSON
- [ ] 11 — PostgreSQL con pgx y sqlc

Recursos para profundizar: el *Tour of Go* (go.dev/tour, tiene versión en español) y *Effective Go*.
