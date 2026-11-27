# ADR-0023: Cómo el motor lee y valida la entrada

- **Estado:** Aceptado
- **Fecha:** 2026-11-26

## Contexto

El motor recibe un `EngineInput` en JSON (`contracts/engine.schema.json`). La sesión 1.1.4 tenía que rechazar las
entradas inválidas con errores claros y sin `panic`, incluidos los 17 ejemplos de `contracts/examples/invalid/`.

Al probar `encoding/json` con Go 1.27 aparecieron tolerancias que el contrato no permite y que serían errores
silenciosos:
- `null` se lee como «no vino»: `"allowedOrientations": null` daría las 6 orientaciones;
- una llave repetida se queda con el último valor;
- los nombres de campo se comparan sin distinguir mayúsculas, incluso con `DisallowUnknownFields`;
- un campo obligatorio que falta queda en su valor cero, que a veces es válido (`xMm` en 0).

Además, varios errores solo se ven al leer el texto (decimales, campos desconocidos, contenido de más) y no en el
struct ya lleno.

## Decisión

1. **Dos pasos.**
   - `DecodeInput(io.Reader)` vive en `internal/engine`. Recorre el JSON token por token, guiado por los tipos del
     motor: rechaza `null`, llaves repetidas, llaves que no son de ese objeto (escritas exactamente igual) y campos
     obligatorios que faltan. Después lo lee en modo estricto y llama a `Validate`.
   - `EngineInput.Validate()` revisa las reglas sobre el struct: rangos, ids, textos, tramos, puerta con margen y
     orientaciones.
2. **Los campos obligatorios salen de las etiquetas.** Un campo sin `omitempty` es obligatorio; coincide con los
   `required` del esquema, así que el contrato no se escribe dos veces. `schema_test.go` lee
   `contracts/engine.schema.json` y revisa que los campos, los obligatorios y los mínimos y máximos de `Validate`
   coincidan con el esquema.
3. **Errores por categoría.**
   - Las categorías son variables `ErrAlgo`, que la API distingue con `errors.Is`.
   - Cada problema es un `*FieldError` con la ruta del campo (nombres del contrato) y su categoría.
   - `Validate` junta **todos** en un `*ValidationError`. La lectura se detiene en el primero.
4. **Qué no es error** (como en el ADR-0016): una pieza que no cabe y una puerta más grande que el interior.
5. **Profundidad máxima** de 16 niveles al recorrer el JSON. El tamaño máximo lo pone quien llama (la API, sesión 1.5).
6. **`deliverySequence` pasa a `*int64`**, para que un 0 enviado se distinga de «no vino».
7. **P-027:** `"description": ""` se sigue aceptando y no regresa al volver a JSON. Es un texto que solo se muestra, y
   agregar `minLength` volvería inválidas entradas que hoy son válidas.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Recorrido guiado por tipos + `Validate` (elegida)** | Detecta lo que `encoding/json` tolera; mensajes con la ruta exacta; sin dependencias | Más código que un `Unmarshal` |
| Solo `Unmarshal` estricto + `Validate` | Simple | Deja pasar `null`, llaves repetidas, otras mayúsculas y obligatorios que faltan |
| Validar con una librería de JSON Schema | Aplica el esquema tal cual | Dependencia nueva (requiere su ADR); mensajes en inglés; no evita revisar las reglas semánticas en Go |
| `encoding/json/v2` | En Go 1.27 se importa sin `GOEXPERIMENT`; rechaza llaves repetidas, distingue mayúsculas y rechaza decimales en enteros (comprobado) | Acepta `null` como «no vino» y no avisa si falta un campo obligatorio (comprobado): el recorrido seguiría haciendo falta para eso. Cambiar de paquete se puede revisar cuando haga falta simplificar `scanJSON` |
| La lectura en `httpapi` | El motor solo recibe structs | Cada entrada (API, pruebas, línea de comandos) tendría que repetir las mismas revisiones |

## Consecuencias

- **Positivas:**
  - Los 17 inválidos se rechazan, cada uno con su categoría.
  - La API sabe qué responder con `errors.Is` y qué campo marcar con `errors.As`.
  - La prueba fuzz no encontró ninguna entrada que detenga el programa.
- **Negativas / costos:**
  - Dos lecturas del mismo texto.
  - El recorrido usa `reflect`.
  - `1.0` y `1e3` se rechazan aunque el esquema los considere enteros.
- **Qué tendríamos que hacer para revertirla:** quitar `scanJSON` y dejar `DecodeInput` con la lectura estricta y
  `Validate`, aceptando las tolerancias de `encoding/json`.
