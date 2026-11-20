# Visión general del sistema

## La idea central

El valor del producto está en el **motor de acomodo**. El motor vive dentro del
**backend en Go**, en el mismo binario que la API: la API recibe la carga, llama al motor como una
función normal y devuelve el plan. Las pantallas (web, escritorio, móvil) solo muestran y editan.

```
   apps/web (React + R3F)   apps/desktop (Wails)   apps/mobile (Expo)
            │                       │                      │
            │      packages/viewer-3d · packages/validator · packages/api-client
            └───────────────┬───────┴──────────────────────┘
                            │  HTTP + JSON  (contracts/openapi.yaml)
                            ▼
              ┌──────────── backend/ (Go) ─────────────┐
              │  cmd/api          arranque y apagado   │
              │  internal/httpapi rutas y validación   │
              │        │                  │            │
              │        ▼                  ▼            │
              │  internal/engine    internal/store     │
              │  motor de acomodo   acceso a datos     │
              │  (sin HTTP ni BD)   (pgx + sqlc)       │
              └───────────────────────────┬────────────┘
                                          │
                                          ▼
                                   PostgreSQL (Docker)
```

Wails es el caso especial: la app de escritorio puede llamar al motor en Go directamente,
sin servidor, para trabajar sin conexión.

## Recorrido de un plan de carga (fase 2)

1. El usuario importa su lista de piezas desde Excel en `apps/web`.
2. La web valida el formato con los tipos generados desde los contratos.
3. La web llama a `POST /plans/optimize` en la API.
4. `internal/httpapi` valida la petición y llama a `engine.Optimize(...)`.
5. El motor calcula el acomodo (la metaheurística usa varios núcleos con goroutines) y devuelve el plan.
6. La API guarda el plan en PostgreSQL (JSONB) y lo responde.
7. `packages/viewer-3d` lo dibuja; si el usuario mueve una pieza, `packages/validator` revisa en vivo.
8. En el patio, `apps/mobile` descarga el plan y guía la carga paso a paso, aun sin señal.

## Piezas compartidas

| Pieza | Qué contiene | Quién lo usa |
|---|---|---|
| `contracts/openapi.yaml` | Contrato de la API | Backend, generador del cliente |
| `contracts/engine.schema.json` | Entrada y salida del motor | Motor, validador, casos de prueba |
| `engine/fixtures/` | Cargas de prueba en JSON | Pruebas del motor (cualquier lenguaje) |
| `packages/types` | Tipos TypeScript generados | Web, escritorio, móvil |
| `packages/api-client` | Cliente generado desde OpenAPI | Web, escritorio, móvil |
| `packages/viewer-3d` | Componente 3D | Web, escritorio |
| `packages/validator` | Validador independiente de planes | Web (editor), móvil |

> Este documento se actualiza cada vez que cambia el flujo o aparece una pieza nueva.
