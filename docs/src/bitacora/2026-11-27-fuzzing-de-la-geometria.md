# 27 nov 2026 · Fuzzing de la geometría

**Sesión:** 1.2.3 · **Área:** motor, docs

## Qué se hizo

- **Fuzzing nativo de Go** en `backend/internal/engine/box_fuzz_test.go`, una función `Fuzz…` por propiedad (con
  `go test` normal solo corren sus semillas, los casos límite de las pruebas de tabla):
  - `FuzzIntersectsIsSymmetric`: `Intersects` es simétrica con **todo** el rango de `int64` (la fórmula lo es aun con
    desbordamiento).
  - `FuzzBoxIntersectsItselfOnlyWithVolume`: una caja se encima consigo misma si y solo si tiene volumen.
  - `FuzzContainedInIsTransitive`: si A está en B y B en C, A está en C; las cajas se arman en cadena para que la
    premisa se cumpla (`boxInside`, con medida, posición y empuje de ±50 mm propios por eje).
  - `FuzzRotationKeepsVolume`: girar una pieza no cambia su volumen, ni el de cada candidata de `Candidates`.
  - `FuzzContactIsBoundedAndExcludesOverlap` (extra): `Contact` no es negativa ni pasa de la base más chica, y con
    contacto no hay traslape.
- **P-029, opción B:** `MaxGeometryMM` (100 000 mm) y `Box.InRange()` en `box.go`, con pruebas
  (`TestInRangeAcceptsOnlySmallNumbers` y `TestIntersectsIsWrongOutOfRange`, que muestra cómo `Intersects` responde mal
  con una coordenada cerca del tope de `int64`). La sección «Desbordamiento» de `conceptos/geometria.md` quedó
  actualizada y se agregó «Las propiedades que prueba el fuzzing».
- **Capítulo 08 de la guía de Go, parte 1** (`aprendizaje/go/08-pruebas-de-tabla-y-fuzzing.md`): pruebas de tabla,
  qué es el fuzzing, cómo se corre, qué hacer cuando falla, cómo se piensa una propiedad y trampas. Se comprobó con
  programas aparte: un solo `Fuzz…` por corrida, los tipos que admite `f.Fuzz`, que sin `-run '^$'` corren antes las
  pruebas normales y el resto de `%` con negativos. Para ver una falla real se rompió `Intersects` a propósito en una
  copia fuera del repo: Go la encontró en 0,12 s, la guardó en `testdata/fuzz/` y `go test` normal la repitió.
- **`CLAUDE.md`, «Comandos»:** los cinco comandos de fuzzing, con la nota de qué hacer cuando uno falla.
- **Revisión (`/revisar`):** cuatro hallazgos, todos corregidos:
  - `wrap` corría todo número por el inicio del rango y desviaba las semillas (30 000 llegaba como 1): ahora deja igual
    lo que ya está en rango; `TestWrapKeepsNumbersInRange` y `TestBoundedBoxKeepsTheSeeds` lo vigilan;
  - `boxInside` ahora tiene parámetros propios por eje;
  - los extremos de `int64` usan `math.MaxInt64` y `math.MinInt64`;
  - el capítulo ya no da una ruta de caché propia de esta computadora.

## Decisiones

- **P-029, opción B** (decidida con Benjamín): los límites viven en el código (`MaxGeometryMM`, `InRange`) y el máximo
  del contrato entra con el validador en la 1.4.2. Se evitó subir la versión del contrato ahora, porque los ejemplos y
  fixtures llevan la versión `0.1.0` y habría que tocar unos 90 archivos. Mientras tanto, quien reciba un plan de fuera
  llama a `InRange` antes de usar una caja. No hubo ADR nuevo.
- **Capítulo 08 en dos partes** (decidido con Benjamín): la parte 1 se escribe ahora, porque el comentario de
  `FuzzDecodeInput` ya prometía el tema para el sprint 1.2; la parte 2 (benchmarks) se completa en la 1.4.5, donde el
  prompt escribe ese capítulo.
- Cada propiedad atrapa solo ciertos errores: con un error puesto sobre `x`, `FuzzContactIsBoundedAndExcludesOverlap`
  corrió 30 s sin verlo y la simetría sí. Un `PASS` solo dice que no se encontró un contraejemplo en esa corrida.

## Verificación

- `gofmt -l backend` sin salida; `go vet ./backend/...` y `go test ./backend/... -count=1` sin errores.
- Cobertura de `internal/engine`: 96,9 %. `mdbook build docs` sin avisos. No se creó ningún `testdata/`.
- **Fuzzing, 1 minuto por propiedad, sin fallas** (`-fuzztime 1m`, de uno en uno, unos 300 000 ejecuciones por segundo):

  | Propiedad | Ejecuciones |
  |---|---|
  | `FuzzIntersectsIsSymmetric` | 17,9 millones |
  | `FuzzBoxIntersectsItselfOnlyWithVolume` | 18,2 millones |
  | `FuzzContainedInIsTransitive` | 16,0 millones |
  | `FuzzRotationKeepsVolume` | 15,9 millones |
  | `FuzzContactIsBoundedAndExcludesOverlap` | 16,7 millones |

  Tras corregir `wrap` y `boxInside` se repitieron los cuatro objetivos afectados (los de la tabla, salvo la simetría,
  que no usa `wrap`) y volvieron a pasar.

## Pendientes

- P-029: pasó de «Abiertos» a «Agendados» (1.4.2): máximos de `Placement` en el contrato y validación; escrito en el
  prompt de esa sesión.
- P-031 (nuevo, agendado 1.4.5): completar la parte 2 del capítulo 08.
- Sigue agendado P-030 (1.4.4), de la sesión anterior.
