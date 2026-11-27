# 26 nov 2026 · Candidatas por orientación

**Sesión:** 1.2.2 · **Área:** motor, docs

## Qué se hizo

- **`CargoItem.Candidates(at, space)`** en `backend/internal/engine/candidates.go`: para una pieza y un punto, regresa
  una caja por cada orientación permitida que queda dentro de `space` (`Box.ContainedIn`), con la esquina mínima en
  `at`. Si dos orientaciones dan la misma caja (piezas cúbicas o con dos lados iguales), solo queda la primera de la
  lista; si no cabe ninguna, regresa `nil`. Va con dos tipos nuevos: `Point` (el punto donde se prueba la esquina; los
  puntos extremos del sprint 1.3 serán una lista de ellos) y `Candidate` (orientación y caja).
- **`Section.Box()`** en `container.go`: la caja nominal del tramo, sin margen de rebase.
- **Pruebas** en `candidates_test.go`: 12 casos de tabla (sin lista, «este lado arriba», cúbica, cúbica con otra
  lista, dos lados iguales, tubo largo en un 20' y en un 40', pegada a la puerta, poco alto libre, punto fuera,
  orientación desconocida y lista vacía) y una prueba de que con «este lado arriba» el alto queda siempre vertical.
  En cada candidata se revisa la esquina, las medidas contra `RotatedDimensions`, que quede dentro del espacio y que
  no haya medidas repetidas. `Section.Box()` tiene su prueba en `container_test.go`.
- **Capítulo 05 de la guía de Go** (`aprendizaje/go/05-slices-maps-y-ciclos.md`): slices (`len`, `cap`, `make`,
  `append`, `nil` contra vacío, arreglo de fondo compartido, parámetro variádico sin argumentos), mapas como conjunto
  con llave struct, mapa `nil` y orden aleatorio, y ciclos (`range` da copias, una variable por vuelta desde Go 1.22).
  Todo se comprobó con un programa en Go 1.27 antes de escribirlo (ADR-0022).
- **Concepto:** sección «Candidatas» en `conceptos/geometria.md`, con un dibujo de las orientaciones que caben con
  poco alto libre. `doc.go` lista los tipos nuevos.
- **Revisión (`/revisar`):** cinco hallazgos, todos corregidos:
  - se borró una prueba que no podía fallar y el capítulo ahora cita la que sí revisa `slices.Clone`
    (`TestOrientationsReturnsACopy`);
  - se rehízo el dibujo de las candidatas, que mostraba las cajas al revés;
  - `Section.Box()` pasó a `container.go`;
  - se quitó una rama de `append` que no hacía falta;
  - se quitó del capítulo una cifra que cambia en cada corrida.

## Decisiones

- `Candidates` recibe el espacio como una `Box` y no como un `Section`: quien llama decide si pasa el tramo nominal o
  el tramo con margen de rebase (eso lo arma el algoritmo en el sprint 1.3). No cambia el contrato; son tipos internos.
- Cuando dos orientaciones dan la misma caja, se queda la **primera** según el orden de `allowedOrientations`
  (o del contrato si no vino). El orden lo da el slice, nunca un mapa, para que el plan sea determinista.
- «Este lado arriba» no lleva regla aparte: es `allowedOrientations = ["LWH", "WLH"]`, como ya decía el contrato.
- Una orientación desconocida se salta sin `panic`; la validación la rechaza antes (segunda línea de defensa).
- No hubo ADR nuevo: ninguna decisión cambia la arquitectura.

## Verificación

- `gofmt -l backend` sin salida; `go vet ./backend/...` y `go test ./backend/... -count=1` sin errores.
- Cobertura de `internal/engine`: 96,8 %.
- `mdbook build docs` sin avisos.

## Pendientes

- P-030: medir el costo del mapa `seen` de `Candidates` (sesión 1.4.4).
- Sigue abierto P-029 (desbordamiento de `int64` con planes de fuera; sesión 1.2.3).
