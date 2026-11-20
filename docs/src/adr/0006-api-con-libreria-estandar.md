# ADR-0006: API en Go con la librería estándar

- **Estado:** Aceptado
- **Fecha:** 2026-11-20

## Contexto

La API maneja usuarios, empresas, catálogos, historial de planes, importaciones y, más adelante, cobros.
En Go hay muchos frameworks web, pero cada dependencia es algo más que aprender y mantener.

## Decisión

- La API usa **`net/http` de la librería estándar**. Desde Go 1.22, `http.ServeMux` entiende métodos y
  parámetros de ruta (`GET /plans/{id}`), que era la razón principal para usar un framework.
- API y motor viven en el **mismo binario**: la API llama al motor como una función, sin red de por medio.
- Logs con `log/slog` (estándar), configuración por variables de entorno y apagado ordenado con `context`.
- Si un patrón se repite mucho (middleware de autenticación, validación), se escribe como función propia
  antes de traer una librería; agregar una dependencia requiere un ADR.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Librería estándar (elegida)** | Nada que instalar; estable por años; lo que se aprende sirve en cualquier proyecto Go | Algo más de código para middleware |
| chi | Muy ligero y compatible con `net/http` | Una dependencia más sin necesidad clara hoy |
| Gin / Echo / Fiber | Populares, con muchas utilidades | Estilos propios; Fiber ni siquiera usa `net/http` |
| Servicio de motor separado | Escala por separado | Dos servicios que operar; innecesario al inicio |

## Consecuencias

- **Positivas:** cero dependencias en el servidor; el código se parece al de cualquier tutorial oficial de Go.
- **Negativas:** middleware y validación se escriben a mano (pocas funciones).
- **Revertir:** chi es compatible con `net/http`; se podría adoptar sin reescribir los handlers.
