# ADR-0008: Diseño API-first con OpenAPI escrito primero

- **Estado:** Aceptado
- **Fecha:** 2026-11-20

## Contexto

Habrá varios clientes (web, escritorio, móvil y quizá ERPs de clientes). Con el backend en Go y las
pantallas en TypeScript, ya no se comparten tipos por compartir lenguaje: hace falta un contrato.

## Decisión

- El contrato vive en **`contracts/openapi.yaml`** y **se escribe primero**; el código se ajusta a él.
- `packages/api-client` y los tipos de `packages/types` se **generan** desde el contrato.
- En Go, los tipos de petición y respuesta se escriben a mano al inicio y se verifican contra el contrato
  con pruebas; si el contrato crece, se evalúa generarlos (oapi-codegen) en un ADR.
- Ningún cliente habla con la base de datos directamente.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **OpenAPI primero (elegida)** | Contrato independiente del lenguaje; documentación navegable; integraciones de terceros listas | Hay que mantener el archivo al día |
| OpenAPI generado desde el código | El código es la fuente | Go no tiene una forma estándar de hacerlo; el contrato queda en segundo plano |
| gRPC | Tipado fuerte de punta a punta | Incómodo desde el navegador; no es lo que esperan integraciones externas |

## Consecuencias

- **Positivas:** backend y pantallas se pueden trabajar por separado contra el mismo contrato.
- **Negativas:** disciplina de cambiar el contrato antes que el código (regla 10 de `CLAUDE.md`).
