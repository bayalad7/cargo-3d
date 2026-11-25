# ADR-0020: Casos simulados en lugar de cargas reales durante la construcción

- **Estado:** Aceptado
- **Fecha:** 2026-11-25

## Contexto

El plan pedía cargas reales, las reglas de un patio y medidas de equipos tomadas en el patio antes de la sesión 1.1.5
(P-001, P-003, P-004 y P-013), y la puerta de la fase 1 era «una carga real acomodada». El proyecto lo desarrolla una
sola persona que no tiene patio, equipos ni cargas a la mano: la hoja de medición se quedó vacía y los CSV del patio
terminaron con copias de los ejemplos sintéticos.

## Decisión

Mientras se construye la app, las cargas reales se sustituyen por **casos de uso simulados** en `docs/src/casos/`,
clasificados por subcategoría (`casos/<categoría>/<subcategoría>/cs-NNN-<tema>.md`). Cada caso usa equipos del
catálogo con fuente o equipos simulados marcados como tales, y dice qué resultado se espera del motor. La sesión 1.1.5
los convierte en fixtures.

La validación con datos reales se hace en un **periodo de pruebas** que empieza con la versión v0.2.0 (MVP), con el
procedimiento de [Pruebas con cargas reales](../casos/pruebas-reales.md): los registros de cada empresa quedan en
`docs/datos/pruebas/` y lo aprendido regresa a los casos como `Validado en real`, `Corregido` o caso nuevo.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Casos simulados ahora, pruebas reales en v0.2.0 (elegida)** | No bloquea el desarrollo; los casos cubren a propósito los bordes | Los datos simulados pueden no parecerse a la operación real |
| Esperar cargas reales | Datos verdaderos desde el inicio | Bloquea la fase 1 sin fecha |
| Pruebas reales desde v0.1.0 | Antes | La v0.1.0 no guarda equipos ni planes: una empresa no la puede usar |

## Consecuencias

- **Positivas:** la fase 1 avanza; cada caso documenta el porqué de su resultado; los casos de la fase 2 sirven de
  especificación de las restricciones.
- **Negativas / costos:** el riesgo de «acomodos imposibles de cargar» (plan, riesgos) se descubre hasta la v0.2.0.
- **Cambios:** se cierran P-001 y P-003 como sustituidos; P-004 y P-013 conservan solo lo que se valida en el periodo de
  pruebas; las puertas de las fases 0 y 1 y el sprint 2.6 se reescriben; el nombre del cliente previsto se quita de
  los documentos que se suben a git.
- **Qué tendríamos que hacer para revertirla:** volver a pedir cargas reales antes de la 1.1.5 con la hoja de medición
  de `docs/plantillas/patio/`, que se conserva para el periodo de pruebas.
