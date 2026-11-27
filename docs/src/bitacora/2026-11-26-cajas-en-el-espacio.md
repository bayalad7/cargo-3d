# 26 nov 2026 · Cajas en el espacio

**Sesión:** 1.2.1 · **Área:** motor, docs

## Qué se hizo

- **`Box`** en `backend/internal/engine/box.go`: caja alineada a los ejes (esquina mínima y medidas, en mm `int64`),
  con los mismos nombres de campo que `Placement`. Funciones:
  - `Intersects`: se enciman solo si los rangos se traslapan en los tres ejes; tocarse por una cara, una arista o un
    vértice no cuenta. Es simétrica.
  - `ContainedIn`: pegada al borde cuenta como dentro; pasarse 1 mm no.
  - `Contact`: área en mm² donde la cara de abajo de una caja se apoya en la de arriba de otra; no es simétrica.
  - `Placement.Box()`: pasa una colocación del plan a caja.
- **Pruebas de tabla** en `box_test.go`: 46 casos, entre ellos todos los que pedía la sesión (cara, arista, vértice,
  una dentro de otra, borde exacto del equipo, separadas 1 mm) más la cruz en «+», `z` negativa por el rebase y el
  área de contacto máxima del contrato. `Intersects` se prueba en los dos órdenes.
- **Concepto** `docs/src/conceptos/geometria.md`, con dibujos en texto, enlazado en `SUMMARY.md` y desde
  `sistema-coordenadas.md`. `doc.go` lista `Box` entre los tipos del paquete.
- **Guía de Go:** cuatro filas nuevas en la tabla de `aprendizaje/go/00-introduccion.md` (`min` y `max`, cortocircuito
  de `&&` y `||`, comparar structs con `==` y desbordamiento). Los capítulos 05 y 06 tocan en las sesiones 1.2.2 y 1.2.4.
- **Revisión (`/revisar`):** seis hallazgos, todos corregidos. El de más peso: el límite con el que se descartaba el
  desbordamiento de `int64` no era exacto (una pieza puede llegar a `x = 60 000` y el área de contacto a 9 × 10⁸ mm²),
  y el contrato no pone máximo a las coordenadas de una colocación. Quedó como P-029.

## Decisiones

- **Una caja sin volumen** (alguna medida ≤ 0) nunca se encima con nada, no se apoya en nada y no cuenta como
  contenida. Es la segunda línea de defensa tras la validación de la entrada. No lleva ADR propio: la regla de
  «tocarse no es encimarse» ya estaba en el ADR-0015 y esta queda explicada en `geometria.md`.
- `Box` no cambia el contrato: es un tipo interno del motor.
- `ContainedIn` no revisa aparte que el espacio tenga volumen: sale solo de comparar los bordes.
- Las comprobaciones de Go (`min` y `max`, tamaño de `Box` de 48 bytes, vuelta de `int64`) se hicieron con un programa
  antes de escribirlas, como pide el ADR-0022.

## Verificación

- `gofmt -l backend` sin salida; `go vet ./backend/...` y `go test ./backend/... -count=1` sin errores.
- Cobertura de `internal/engine`: 96,7 % (antes 96,4 %).
- `mdbook build docs` sin avisos.

## Pendientes

- P-029: proteger la geometría del desbordamiento de `int64` con planes que llegan de fuera (sesión 1.2.3).
