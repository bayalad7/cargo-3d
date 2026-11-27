# 26 nov 2026 · Orientaciones, colocaciones y plan de carga

**Sesión:** 1.1.3 · **Área:** motor, docs

## Qué se hizo

- **Orientaciones** (`orientation.go`): `Orientation`, que se movió desde `cargo_item.go`, y `Dimensions`.
  - `CargoItem.RotatedDimensions(o)` regresa las medidas de la pieza girada y `ok = false` si la orientación no existe,
    sin detener el programa.
  - `CargoItem.Orientations()` regresa las seis solo si la lista no vino (`== nil`) y siempre una copia (parte de P-025).
  - El valor por omisión sale de la función `allOrientations()`, que arma una lista nueva en cada llamada, y no de una
    variable del paquete que otro archivo podría cambiar.
- **Colocación** (`placement.go`): `Placement` con `Volume()` y `Dimensions()`. `OutOfGauge` es un `bool` con `omitempty`:
  el contrato solo admite `true`, así que «no vino» y `false` significan lo mismo.
- **Plan de carga** (`load_plan.go`): `LoadPlan`, `UnplacedItem`, `UnplacedReason` con sus seis motivos y `Metrics`.
  - `NewLoadPlan` crea las listas vacías, porque un slice `nil` se escribe como `null` y el contrato exige `[]`.
  - `AddPlacement` y `AddUnplaced` usan receptor por puntero. `AddUnplaced` junta las unidades por pieza y motivo.
  - Las métricas existen como tipo; su cálculo es del sprint 1.3.
- **Pruebas:**
  - Tabla de las seis orientaciones y de orientaciones desconocidas.
  - El valor por omisión y la copia de `Orientations()`.
  - Ida y vuelta de los 27 `output.json`, con la prueba convertida en una función genérica `roundTrip[T any]` que también
    sirve para los `input.json`.
  - Las medidas de cada colocación de todos los ejemplos coinciden con `RotatedDimensions`.
  - `NewLoadPlan` escribe `[]` y el valor cero escribe `null`; agrupación de `AddUnplaced`.
  - Se cambió a propósito una fila de la tabla de orientaciones y fallaron la tabla y el ejemplo `53dv-contract-coverage`;
    después se regresó.
- **Documentación:** capítulo 03 de la guía (métodos y receptores) con ejemplos compilados o corridos con Go 1.27;
  `sistema-coordenadas.md` coincidía con el código y con el ejemplo que cita, y ahora dice dónde vive cada cosa en Go y qué
  prueba lo protege; referencias del capítulo 02 actualizadas.
- **Revisión (`/revisar`):** sin hallazgos bloqueantes ni importantes. Se corrigió el conteo de campos de `LoadPlan` en
  el capítulo 03, se agregó `orientation.go` al capítulo 02, `allOrientations` pasó a función y se quitó una predicción
  sin base sobre el sprint 1.3.

## Decisiones

- `OutOfGauge` como `bool` y no como puntero: con un contrato que solo admite `true`, no hay un `false` con significado.
- Receptor por puntero en todos los métodos de `LoadPlan` y por valor en los tipos que el motor solo lee (convención de
  no mezclar receptores en un mismo tipo).
- Las métricas se calculan en el sprint 1.3, como dice el plan.

## Verificación

- `gofmt -l backend` sin salida; `go vet ./backend/...` y `go test ./backend/...` sin errores; `mdbook build docs` sin
  errores.
- No se usó Docker.

## Pendientes

- P-025: la mitad de la 1.1.3 quedó hecha; falta rechazar la lista vacía y las orientaciones desconocidas en la 1.1.4.
- P-026: `AddUnplaced` ignora en silencio una cantidad menor a 1 (decidir en la 1.3.3).
- P-027: descripción vacía y una prueba que no ejercita código del motor.
- P-010: ahora incluye que solo el esquema rechaza un `"outOfGauge": false`.
