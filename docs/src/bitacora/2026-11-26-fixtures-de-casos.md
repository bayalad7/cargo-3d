# 26 nov 2026 · Fixtures de los casos simulados

**Sesión:** 1.1.5 · **Área:** motor, docs

## Qué se hizo

- **36 fixtures de 26 casos** en `engine/fixtures/casos/cs-NNN/`: un `input.json` por caso, una variante por cada equipo,
  camino o dato que el caso compara (`input-camino-b.json`, `input-techo.json`, `input-caja-real.json`…) y un
  `.expect.json` por entrada con el caso, su fase, las cuentas esperadas (colocadas, sobredimensionadas, peso cargado, no
  colocadas con su motivo) y las reglas que no son números.
  - Los 21 casos de la fase 1.
  - Cinco de la fase 2 que el contrato 0.1.0 ya puede representar: CS-021, 022, 025, 026 y 027.
  - Quedan fuera CS-016 (falta la altura de cama), CS-023 (sin campos de peso por eje) y CS-029 (dos tramos, P-012).
- Se generaron con un script temporal, que no está en el repo, a partir de una sola tabla escrita desde los casos.
- **Prueba** (`fixtures_test.go`): cada entrada pasa `DecodeInput` y tiene su `.expect.json`, leído en modo estricto.
  Revisa las cuentas (unidades y pesos), que haya sobredimensionadas solo con margen y que los motivos que no dependen
  del algoritmo (no cabe, no pasa por la puerta, pesa más que la carga útil) coincidan con la geometría de la entrada.
  También revisa que ningún `.expect.json` quede sin entrada y que cada carpeta de caso tenga su `input.json`.
- **README de los fixtures:** formato del `.expect.json`, fuente de cada medida y tabla de los 36 fixtures.
- **Correcciones a los casos:** CS-010 (el módulo ancho pasa a solo `LWH`, porque con `WLH` sí cabía) y CS-026 (el panel
  pasa a 2 490 de ancho, porque con la puerta del repo uno de 2 500 no entraba ni en la caja de catálogo).
- **Revisión (`/revisar`):** el CS-026 era un error del plan de la sesión, y la prueba no lo detectó porque solo revisaba
  cuentas. Se agregó la revisión de la geometría y se comprobó que habría detectado el error.

## Decisiones

Siete decisiones aprobadas al inicio:
- CS-010 con solo `LWH`.
- `48rv` y `53rv` con las medidas de los ejemplos del repo, marcadas (P-020).
- `40rf` sin puerta, porque no está publicada (P-013).
- Puerta de 2 496 × 2 794 del repo en todas las cajas secas.
- CS-028 con las dos variantes, aunque el margen del permiso es por tramo (P-023).
- CS-009 por el techo en `xMm = 278`.
- Los casos de la fase 2 que el contrato ya puede representar se convierten desde ya.

## Verificación

- `gofmt -l backend` sin salida; `go vet ./backend/...` y `go test ./backend/...` sin errores; `mdbook build docs` sin
  errores.
- Se cambió a propósito el peso esperado de CS-013 y el ancho del panel de CS-026: la prueba detectó los dos y se
  regresaron.
- No se usó Docker. `node` no está instalado en esta computadora, así que no se corrió `pnpm examples:check`; esta
  sesión no tocó `contracts/examples/`.

## Pendientes

- P-028: casos que todavía no son fixtures o cuyas reglas se exigen en una sesión posterior.
- P-020: ahora también cubre las medidas de CS-013 y CS-014.
- Los resultados esperados salen de las cuentas de cada caso, no de un algoritmo: el prompt de la 1.3.3 ya los usa como
  criterio y ahora lo dice con más detalle.
