---
name: cambiar-contrato
description: Cambia primero el contrato (contracts/openapi.yaml o contracts/engine.schema.json) según la regla 10, y se detiene antes de tocar el código hasta que Benjamín apruebe.
argument-hint: "[descripción del cambio]"
disable-model-invocation: true
---

# Cambio de contrato: $ARGUMENTS

Responde siempre en español. Aplica la regla 10 de `CLAUDE.md`: **el contrato cambia primero, el código
después**. En esta skill solo se cambian archivos de `contracts/`.

## Estado de los contratos

!`ls -1 contracts contracts/examples 2>/dev/null || true`
!`git status --short contracts || true`

## Qué hacer

1. Si la descripción está vacía, pregunta qué cambio se necesita y detente.
2. Decide qué contrato cambia:
   - **API** (rutas, peticiones, respuestas, códigos de error): `contracts/openapi.yaml`.
   - **Motor** (entrada o salida de `Optimize`): `contracts/engine.schema.json`.
   - Pueden ser los dos.
3. **Propón el cambio sin escribirlo todavía:**
   - el fragmento nuevo o modificado del contrato;
   - un ejemplo de petición y respuesta (o de entrada y salida del motor);
   - si el cambio **rompe compatibilidad** (quita o renombra campos, cambia tipos, vuelve obligatorio algo
     opcional) o no;
   - la versión nueva siguiendo SemVer: cambio que rompe → versión mayor; campo o ruta nueva compatible →
     versión menor; corrección de descripción → parche. (Mientras la versión sea 0.x, un cambio que rompe
     sube la versión menor.)
   - la lista de lugares que habrá que cambiar después: código Go (`internal/httpapi`, `internal/engine`,
     la constante `engine.Version`), `packages/validator`, tipos generados, fixtures en `engine/fixtures/`
     y `contracts/examples/`.
4. Espera la aprobación de Benjamín.
5. Cuando apruebe:
   - escribe el cambio en el contrato y actualiza los ejemplos en `contracts/examples/`;
   - sube la versión del contrato (`info.version` en OpenAPI y la versión del esquema del motor);
   - valida que el YAML y el JSON sean válidos y que los ejemplos cumplan el esquema;
   - si existen paquetes generados desde el contrato, regenéralos con su script de pnpm.
6. **Detente antes de cambiar el código Go o TypeScript.** Resume qué cambió y la lista del paso 3, y
   pregunta si seguimos con la implementación en esta misma sesión o en otra.
