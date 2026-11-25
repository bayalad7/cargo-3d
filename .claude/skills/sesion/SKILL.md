---
name: sesion
description: Ejecuta una sesión de la biblioteca docs/src/prompts.md por su número (por ejemplo 1.1.2). Siempre empieza proponiendo un plan y espera aprobación antes de escribir código.
argument-hint: "[número de sesión, ej. 1.1.2]"
arguments: [numero]
disable-model-invocation: true
---

# Ejecutar la sesión $numero

Responde siempre en español. Sigue todas las reglas de `CLAUDE.md`.

## Rama actual

!`git branch --show-current || true`

## Qué hacer

1. **Encuentra la sesión.** Abre `docs/src/prompts.md` y busca el encabezado que empieza exactamente con
   `### $numero ` (tres almohadillas, espacio, el número y un espacio). Lee completa esa sección, hasta el
   siguiente encabezado `###`, `##` o `#`.
   - Si no existe, dilo, muestra los números de sesión más cercanos que sí existen y detente.
   - Si `$numero` está vacío, pide el número y detente.
2. **Revisa los requisitos** antes de hacer nada:
   - **Pendientes:** busca en «Abiertos» y «Agendados» de `docs/src/pendientes.md` los que tienen esta sesión en «Cuándo»;
     inclúyelos en el plan (y, si se resuelven, se cierran en `/cerrar-sesion`).
   - **Necesitas:** confirma que lo que pide existe (archivos, sesiones previas terminadas, datos en `docs/datos/`).
     Si falta algo, dilo y detente.
   - **Rama:** todo se trabaja en `develop` (ver "Ramas" en `CLAUDE.md`); ignora el campo "Rama" de la sesión
     si existe. Si la rama actual no es `develop`, dile a Benjamín que cambie con `git switch develop` y espera.
   - **Modelo:** si el modelo sugerido es distinto al que estás usando, menciónalo en una línea
     (Benjamín lo cambia con `/model` y `/effort`).
3. **Si el prompt de la sesión es una o más skills** (por ejemplo `/cerrar-sprint ...`),
   no las imites: dile a Benjamín que las ejecute él, en orden, y detente.
4. **Propón el plan, sin modificar archivos todavía.** Toma el bloque "Prompt" de la sesión como la tarea y
   presenta:
   - qué vas a hacer, paso a paso;
   - qué archivos vas a crear o cambiar;
   - si cambia la API o el motor, el cambio de contrato en `contracts/` va primero (regla 10 de `CLAUDE.md`);
   - qué pruebas vas a escribir;
   - qué documentación vas a escribir o actualizar (guía de aprendizaje, conceptos, ADR);
   - qué conceptos nuevos de Go o TypeScript aparecerán y cómo los vas a explicar.
5. **Espera la aprobación explícita de Benjamín** ("sí", "adelante", o sus correcciones). No escribas código
   antes.
6. **Implementa en pasos chicos.** Después de cada paso importante, corre `gofmt -l backend`,
   `go vet ./backend/...` y `go test ./backend/...` (y las pruebas de TypeScript si tocaste paquetes de TS).
   Comenta el código con el nivel de detalle que indica la sección "Estado" de `CLAUDE.md`.
7. **Al terminar**, compara contra el campo "Terminada cuando" de la sesión y dile a Benjamín qué se cumplió,
   qué no y por qué. Recuérdale que puede correr `/revisar` y después `/cerrar-sesion`.
