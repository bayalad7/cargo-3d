# Contratos

Archivos que definen cómo se comunican las partes del sistema, **independientes del lenguaje**.
Si el motor o la API se reescriben en otro lenguaje, estos contratos no cambian.

| Archivo | Qué define | Quién lo usa |
|---|---|---|
| `openapi.yaml` | Rutas, peticiones y respuestas de la API | Backend Go (servidor), `packages/api-client` (generado), web y móvil |
| `engine.schema.json` | Entrada y salida del motor (`optimize(input) -> LoadPlan`) | Motor, validador, casos de prueba de `engine/fixtures/` |

Regla: **el contrato se cambia primero**, en su propio commit, y después el código.
Ver `docs/src/adr/0008-api-first-openapi.md`.

> `engine.schema.json` se escribe en el sprint 1.1, junto con los modelos del motor.
