---
name: buscador
description: Lanza el agente investigador-logistico para investigar y verificar con fuentes los datos logísticos (ISO, NOM-012, CTU Code, VGM, pallets) que necesita un sprint o tema de Cargo 3D, antes de escribir código.
argument-hint: "[sprint, sesión o tema, ej. 1.1 o 'pesos por eje NOM-012']"
disable-model-invocation: true
context: fork
agent: investigador-logistico
---

# Investigar datos logísticos

Tema o sprint que se pidió:

$ARGUMENTS

(Si está vacío, toma el sprint actual de la sección "Estado" de `CLAUDE.md` y dilo al inicio de tu respuesta.)

Esta skill es opcional. Si la información que pide el tema ya existe en `docs/datos/`, `contracts/` o
`engine/fixtures/` y está completa, dilo al inicio, investiga solo lo que falte y no repitas lo que ya hay.

Sigue tus instrucciones de agente completas (contexto del proyecto, investigación con fuentes, reglas de
conversión y entrega). No escribas archivos: entrega el paquete y el destino propuesto para cada pieza.
