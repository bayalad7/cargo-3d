---
name: cerrar-sprint
description: Cierra un sprint de Cargo 3D. Verifica con evidencia el checklist del sprint en el plan, marca lo terminado, actualiza el Estado y dice si se puede pasar al siguiente sprint o fase.
argument-hint: "[número de sprint, ej. 1.1]"
arguments: [sprint]
model: opus
effort: high
disable-model-invocation: true
---

# Cerrar el sprint $sprint

Responde siempre en español. Aplica la regla de oro de `docs/src/plan.md`: no se pasa a la siguiente fase
sin cumplir el criterio de terminado de la actual.

## Contexto

- Fecha de hoy: !`date +%F`
- Rama actual: !`git branch --show-current || true`
- Cambios sin commit:
!`git status --short || true`

## Qué hacer

1. Si el número de sprint está vacío, pregúntalo y detente.
2. En `docs/src/plan.md`, busca el sprint `$sprint` (un encabezado `### Sprint $sprint` o, en la fase 3,
   la línea del bloque `$sprint`). Lee su checklist completo.
3. Corre `gofmt -l backend`, `go vet ./backend/...` y `go test ./backend/...` (y las pruebas de TypeScript
   que existan). Si algo falla, el sprint **no** se puede cerrar: dilo y detente.
4. Para **cada** tarea del checklist, busca evidencia concreta en el repositorio (archivo, prueba, capítulo,
   ADR, entrada de bitácora) y presenta una tabla:

   | Tarea | ¿Terminada? | Evidencia |
   | --- | --- | --- |

   Solo cuenta como terminada si la evidencia existe. Revisa también que el punto de control del sprint
   esté registrado como **Aprobado** en la bitácora.
5. Si falta algo, lista exactamente qué falta y qué sesión de `docs/src/prompts.md` lo cubre. No cierres el
   sprint y detente.
6. Si todo está terminado:
   - marca con `[x]` las tareas del sprint en `docs/src/plan.md`;
   - actualiza la sección "Estado" de `CLAUDE.md`: sprint siguiente y su primera sesión;
   - agrega a la bitácora de hoy una sección "Cierre del sprint $sprint" con lo logrado y las mediciones
     relevantes (ocupación, tiempos, número de pruebas);
   - si es el **último sprint de una fase**, revisa también el "Criterio de terminado" de la fase, punto por
     punto con evidencia, y di claramente si la fase está terminada. No la des por terminada si falta un punto.
7. Si `$sprint` es el 1.1, 1.2 o 1.4, recuerda revisar con Benjamín qué tan cómodo se siente con Go y si el
   ritmo del plan es realista; anota su respuesta en la bitácora.
8. Propón el mensaje de commit (formato de `docs/src/convenciones.md`). No hagas el commit tú.
