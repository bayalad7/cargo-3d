---
name: iniciar-sesion
description: Arranca una sesión de trabajo en Cargo 3D. Lee CLAUDE.md, el plan y la última bitácora, revisa el estado de git y resume dónde vamos antes de trabajar.
argument-hint: "[meta de la sesión, opcional]"
disable-model-invocation: true
---

# Iniciar sesión de trabajo

Vas a arrancar una sesión de trabajo en el proyecto Cargo 3D (backend en Go).
Responde siempre en español.

## Estado actual de git

- Rama actual: !`git branch --show-current || true`
- Cambios sin commit:
!`git status --short || true`
- Últimos commits:
!`git log --oneline -5 || true`
- Últimas bitácoras: !`ls -1 docs/src/bitacora 2>/dev/null | sort | tail -3 || true`

## Meta que indicó Benjamín

$ARGUMENTS

(Si está vacío, Benjamín no indicó meta: pregúntasela al final.)

## Qué hacer

1. Lee completos `CLAUDE.md`, `docs/src/plan.md`, `docs/src/pendientes.md` y la bitácora más reciente de
   `docs/src/bitacora/`.
2. Responde con un resumen corto, con estos apartados:
   - **Dónde vamos:** fase, sprint y última sesión terminada (según la sección "Estado" de `CLAUDE.md` y la bitácora).
   - **Pendientes:** de «Abiertos» (y los «Agendados» para la siguiente sesión) en `docs/src/pendientes.md`, primero los que **bloquean** la siguiente sesión
     (su columna «Cuándo» es esa sesión o una anterior), después los de Benjamín; cita cada uno con su id y su
     «Cómo retomarlo». No repitas toda la tabla.
   - **Siguiente sesión sugerida:** el número de sesión de `docs/src/prompts.md` que sigue (por ejemplo `1.1.3`),
     con su modelo sugerido.
3. Revisa y avisa claramente si algo no cuadra:
   - Todo se trabaja en `develop` (ver "Ramas" en `CLAUDE.md`). Si la rama actual no es `develop`, dile a
     Benjamín que cambie con `git switch develop`. No crees ramas.
   - Si hay cambios sin commit de una sesión anterior, pregúntale qué hacer con ellos antes de seguir.
   - Si `CLAUDE.md`, el plan y la bitácora se contradicen, señala la contradicción.
4. Si Benjamín dio una meta y no coincide con la siguiente sesión del plan, señálalo y pregúntale si sigue con su meta.
5. Termina con una sola pregunta: con qué sesión o meta seguimos. No escribas ni modifiques archivos en esta skill.
