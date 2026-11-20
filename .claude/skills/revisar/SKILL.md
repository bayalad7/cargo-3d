---
name: revisar
description: Revisión de código como ingeniero senior de los cambios de develop que aún no se suben, más los que no tienen commit. Entrega hallazgos priorizados sin modificar archivos.
argument-hint: "[enfoque opcional, ej. seguridad o pruebas]"
model: opus
effort: high
disable-model-invocation: true
---

# Revisión de código

Responde siempre en español. Actúa como un ingeniero senior que revisa los cambios antes de que se suban a `develop`.
**No modifiques archivos en esta skill**: solo revisa y reporta.

## Cambios a revisar

- Rama: !`git branch --show-current || true`
- Commits sin subir:
!`git log --oneline origin/develop..HEAD 2>/dev/null || true`
- Archivos cambiados en esos commits:
!`git diff --stat origin/develop...HEAD 2>/dev/null || true`
- Cambios sin commit:
!`git status --short || true`

Enfoque pedido por Benjamín (puede venir vacío): $ARGUMENTS

## Qué revisar

Lee el diff completo (`git diff origin/develop...HEAD` y `git diff` para lo que no tiene commit) y los archivos
necesarios para entenderlo. Revisa:

1. **Correctitud:** errores de lógica, casos límite (ceros, listas vacías, límites exactos, números grandes),
   errores sin revisar, posibles `panic`.
2. **Reglas del motor:** `internal/engine` no importa HTTP ni base de datos; unidades en mm y kg enteros;
   resultados deterministas.
3. **Contrato primero:** si cambió la API o la entrada o salida del motor, `contracts/` cambió también y
   coincide con el código.
4. **Pruebas:** que cubran el comportamiento nuevo y los casos límite; nombres que describan el
   comportamiento; pruebas de tabla cuando aplique.
5. **Documentación:** comentarios de todo lo exportado, nivel de comentario según "Estado" de `CLAUDE.md`,
   capítulos de la guía y documentos de conceptos que debían escribirse, ADR si hubo una decisión.
6. **Seguridad:** entradas sin validar, secretos en el código, datos de una empresa visibles para otra,
   límites de tamaño y tiempo.
7. **Concurrencia** (si hay goroutines): estado compartido sin protección, goroutines que nunca terminan,
   `context` ignorado.
8. **Simplicidad:** código más complejo de lo necesario para que Benjamín lo mantenga.

## Formato de la respuesta

Una lista de hallazgos ordenada de más grave a menos grave. Para cada uno:

- **Gravedad:** Bloqueante / Importante / Menor.
- **Dónde:** archivo y línea.
- **Qué pasa:** el problema y un ejemplo concreto de cuándo falla.
- **Cómo arreglarlo:** en una o dos frases.

Si no encuentras problemas en una categoría, no la menciones. Si no hay hallazgos bloqueantes, dilo
claramente. Termina preguntando cuáles hallazgos quiere Benjamín que se corrijan.
