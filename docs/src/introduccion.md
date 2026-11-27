# Introducción

Bienvenido a la documentación de **Cargo 3D**, software de planeación de carga 3D con el
backend completo en Go y las pantallas en TypeScript.

## Cómo está organizada esta documentación

| Sección | Para qué sirve | Cuándo leerla |
|---|---|---|
| **Plan de trabajo** | Fases, sprints, criterios de terminado y modelo sugerido por paso | Al empezar cada sprint |
| **Conceptos** | Cómo funciona cada parte del sistema | Para entender el panorama |
| **Decisiones (ADR)** | Por qué se eligió cada tecnología y qué alternativas había | Cuando te preguntes "¿por qué está hecho así?" |
| **Aprendizaje** | Go explicado con el código real del proyecto | Mientras aprendes |
| **Bitácora** | Qué se hizo en cada sesión y qué quedó pendiente | Para retomar trabajo |

La documentación del **código** se consulta con `go doc`, por ejemplo:

```bash
go doc ./backend/internal/engine
```

## Fases del proyecto

1. **Fase 0 — Cimientos** *(actual)*: estructura, Docker, servidor `/health` y documentación.
2. **Fase 1 — Prueba de concepto**: motor básico en Go + endpoint de optimización + visor 3D web.
3. **Fase 2 — MVP y piloto**: restricciones reales, PostgreSQL, PDF del plan, pruebas con cargas reales de una empresa piloto.
4. **Fase 3 — Producto y lanzamiento**: multiempresa, app de patio, cobro, beta con empresas del puerto.
5. **Fase 4 — Crecimiento**: integraciones, mejoras del motor y otros puertos.

El detalle está en el [Plan de trabajo](plan.md).
