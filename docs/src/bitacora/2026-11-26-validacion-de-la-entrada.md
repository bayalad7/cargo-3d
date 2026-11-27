# 26 nov 2026 · Lectura y validación de la entrada

**Sesión:** 1.1.4 · **Área:** motor, docs

## Qué se hizo

- **Lectura** (`decode.go`): `DecodeInput(io.Reader)` en dos pasos.
  1. Recorre el JSON token por token, guiado por los tipos del motor, y rechaza lo que `encoding/json` deja pasar:
     `null`, llaves repetidas, llaves con otras mayúsculas o de otro objeto, campos obligatorios que faltan, contenido
     después del objeto y más de 16 niveles anidados.
  2. Lo lee en modo estricto y llama a `Validate`.
- **Validación** (`validate.go`): `EngineInput.Validate()` junta todos los problemas, en el orden de la entrada: rangos
  del contrato, ids (formato y repetidos), descripciones de más de 200 caracteres, número de tramos y tramos encimados,
  puerta con margen y orientaciones vacías, desconocidas o repetidas. Una pieza que no cabe y una puerta más grande que
  el interior no son error.
- **Errores** (`errors.go`): 14 categorías `Err…` (se agregó `ErrMissingField`), `*FieldError` con la ruta del campo y
  `*ValidationError` con `Unwrap() []error`. Los valores enviados se recortan a 64 caracteres en los mensajes.
- `DeliverySequence` pasó a `*int64` para poder rechazar un 0; se agregó `Orientation.IsValid()`.
- **Pruebas:**
  - Un caso por cada regla, lo que no es error y varios problemas a la vez.
  - Los 27 ejemplos válidos pasan y los 17 inválidos dan cada uno su categoría.
  - Cada trampa de `encoding/json`, la causa original dentro del error y el valor enorme recortado.
  - `schema_test.go` lee `contracts/engine.schema.json` y compara campos, obligatorios y mínimos y máximos contra el
    código; se comprobó cambiando a propósito un `omitempty` y un máximo.
  - Fuzz: 4,8 millones de entradas en 30 segundos sin detener el programa.
- **Documentación:** capítulo 04 de la guía (errores como valores) y ADR-0023.
- **Revisión (`/revisar`):** sin hallazgos bloqueantes. Se corrigió una razón falsa del ADR (`json/v2` sí existe en Go
  1.27; no cubre `null` ni obligatorios), se agregó la prueba contra el esquema, se endureció la lectura de `omitempty`,
  se recortaron los mensajes y se quitaron los códigos HTTP del capítulo, que todavía no están decididos.

## Decisiones

- ADR-0023: lectura en dos pasos en `internal/engine`, rechazo de `null` y llaves repetidas, obligatorios tomados de las
  etiquetas y errores por categoría con la ruta del campo.
- P-027, opción A: la descripción vacía se acepta y no regresa al volver a JSON.
- `1.0` y `1e3` se rechazan como enteros aunque el esquema los considere enteros: es más estricto y no deja pasar nada
  malo.

## Verificación

- `gofmt -l backend` sin salida; `go vet ./backend/...` y `go test ./backend/...` sin errores; `mdbook build docs` sin
  errores.
- No se usó Docker.

## Pendientes

- P-008, P-025 y P-027: cerrados en esta sesión.
- El prompt 1.5.2 ahora pide leer con `DecodeInput`, limitar el cuerpo con `http.MaxBytesReader` y sacar los `FieldError`
  con `errors.As`; ya estaba agendado, no hay id nuevo.
