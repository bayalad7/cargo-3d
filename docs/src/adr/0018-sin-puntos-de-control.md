# ADR-0018: Sin puntos de control en el ciclo de sprints

- **Estado:** Aceptado
- **Fecha:** 2026-11-25
- **Sustituye:** el punto 6 del [ADR-0009](0009-idioma-y-documentacion.md)

## Contexto

El ADR-0009 pedía un **punto de control** al terminar cada módulo: la skill `/punto-de-control` hacía cinco
preguntas, corregía las respuestas y registraba un veredicto en la bitácora; `/cerrar-sprint` no cerraba un sprint
sin ese veredicto. El plan lo pedía en 17 sesiones de `prompts.md`.

En la práctica cada punto de control consume muchos tokens (preguntas, seguimientos, correcciones y reintentos), y el
primero, el del contrato del motor (2026-11-21), se repitió dos veces sin aprobarse y llegó a bloquear el avance
(P-007). El aprendizaje de Go ya tiene otros apoyos que cuestan menos y se usan solo cuando hacen falta.

## Decisión

Se quitan los puntos de control del ciclo de sprints: se borra la skill `/punto-de-control`, `/cerrar-sprint` ya no
exige un veredicto y las sesiones «Punto de control y cierre» quedan como «Cierre». El aprendizaje sigue con los
comentarios detallados del código, la guía de Go y `/explicar`, que Benjamín pide cuando algo no queda claro.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Quitarlos (elegida)** | Ahorra tokens y no bloquea sprints | Nadie comprueba que el módulo se entendió |
| Dejarlos opcionales | Se conserva la herramienta | Sigue costando cuando se usa y se olvida cuando no |
| Una sola pregunta por sprint | Más barato | Ya no mide nada y sigue bloqueando el cierre |

## Consecuencias

- **Positivas:** sesiones de cierre más cortas y baratas; el cierre de un sprint depende solo de evidencia en el
  repositorio.
- **Negativas / costos:** la comprensión del código ya no se verifica; depende de que Benjamín pida `/explicar`.
- **Lo que se borró:** la skill, sus menciones en `CLAUDE.md`, `plan.md`, `prompts.md`, `flujo-claude-code.md` y la
  guía de Go, y las dos secciones de punto de control de la bitácora del 2026-11-21. P-007 se cerró como cancelado.
  Las mejoras al contrato que salieron de esos puntos de control (tabla «quién sabe qué» y ejemplos de redondeo en
  `contracts/README.md`) se conservan.
- **Qué tendríamos que hacer para revertirla:** recuperar `.claude/skills/punto-de-control/` del historial de git y
  volver a agregar los puntos de control a las sesiones de cierre de `prompts.md`.
