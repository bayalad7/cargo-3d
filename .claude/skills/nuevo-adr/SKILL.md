---
name: nuevo-adr
description: Crea un nuevo registro de decisión de arquitectura (ADR) en docs/src/adr/ a partir de la plantilla, con el número siguiente, y lo enlaza en el índice.
argument-hint: "[título de la decisión]"
disable-model-invocation: true
---

# Nuevo ADR: $ARGUMENTS

Responde siempre en español (regla 4 de `CLAUDE.md`).

## ADR existentes

!`ls -1 docs/src/adr 2>/dev/null | sort || true`

## Qué hacer

1. Si el título está vacío, pide el título y detente.
2. Calcula el número siguiente: el mayor número de 4 dígitos de la lista anterior, más uno.
3. Reúne la información. Si en la conversación no está claro alguno de estos puntos, pregúntalo antes
   de escribir (todas las preguntas en un solo mensaje):
   - **Contexto:** qué problema nos obliga a decidir y qué restricciones hay.
   - **Decisión:** qué se decidió.
   - **Alternativas:** al menos 2, con lo que tienen a favor y en contra.
   - **Consecuencias:** positivas, negativas y qué haría falta para revertirla.
   Si la decisión depende de información actual (precios, versiones, normas), búscala y cita las fuentes
   en el ADR.
4. Crea `docs/src/adr/NNNN-titulo-en-minusculas-con-guiones.md` siguiendo exactamente
   `docs/src/adr/plantilla.md`, con fecha de hoy y estado **Propuesto**.
5. Agrega el enlace en `docs/src/SUMMARY.md`, en la sección de ADR, en orden numérico.
6. Si el ADR reemplaza a otro, cambia el estado del anterior a "Reemplazado por ADR-NNNN" y enlázalos.
7. Muestra el ADR completo y pregunta si se acepta. Solo cuando Benjamín diga que sí, cambia el estado a
   **Aceptado**. Si la decisión cambia el stack, actualiza también la tabla "Stack" de `CLAUDE.md`.
