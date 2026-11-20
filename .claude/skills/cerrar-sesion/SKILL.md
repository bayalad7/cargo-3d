---
name: cerrar-sesion
description: Cierra una sesión de trabajo de Cargo 3D. Corre las verificaciones, actualiza bitácora, CHANGELOG y el Estado de CLAUDE.md, y propone el mensaje de commit.
argument-hint: "[nota opcional para la bitácora]"
model: sonnet
effort: medium
disable-model-invocation: true
---

# Cerrar la sesión de trabajo

Responde siempre en español.

## Contexto

- Fecha de hoy: !`date +%F`
- Rama actual: !`git branch --show-current || true`
- Archivos cambiados:
!`git status --short || true`
- Commits de develop que aún no se suben:
!`git log --oneline origin/develop..HEAD 2>/dev/null || true`
- Resumen de cambios sin commit:
!`git diff --stat 2>/dev/null || true`

Nota de Benjamín para la bitácora (puede venir vacía): $ARGUMENTS

## Qué hacer

1. **Verificaciones.** Corre, en este orden:
   - `gofmt -l backend` (no debe listar archivos; si lista, corre `gofmt -w` sobre ellos y avisa);
   - `go vet ./backend/...`;
   - `go test ./backend/...`;
   - si existen paquetes de TypeScript con cambios: sus pruebas y su lint con pnpm.

   **Si alguna prueba o `go vet` falla, detente aquí:** explica la falla y pregunta si la arreglamos ahora o
   se cierra la sesión con el pendiente anotado. No ocultes fallas.
2. **Documentación del código.** Revisa que todo lo exportado que se agregó o cambió en esta sesión tenga
   comentario de documentación, como exige `CLAUDE.md`. Si falta, agrégalo.
3. **Bitácora.** En `docs/src/bitacora/<fecha de hoy>.md`:
   - si el archivo no existe, créalo con las secciones "Qué se hizo", "Decisiones" y "Pendientes", y enlázalo
     en `docs/src/SUMMARY.md` (la entrada más reciente va arriba de la lista de Bitácora);
   - si ya existe, agrega una sección nueva con el número de la sesión, sin borrar lo anterior.
   Escribe solo lo que realmente pasó en esta sesión; incluye la nota de Benjamín si dio una.
4. **CHANGELOG.md.** Agrega en "[Sin publicar]" lo que cambió para el producto o el proyecto
   (Agregado / Cambiado / Corregido). No repitas detalles internos.
5. **CLAUDE.md, sección "Estado".** Actualiza fase, sprint, última sesión terminada, siguiente sesión sugerida
   (número de `docs/src/prompts.md`) y pendientes de Benjamín. No cambies otras secciones sin preguntar.
6. **Plan.** Si una tarea del checklist del sprint en `docs/src/plan.md` quedó realmente terminada, márcala
   con `[x]`. No marques tareas a medias.
7. **Mensaje de commit.** Propón un mensaje con el formato de `docs/src/convenciones.md`
   (`área: descripción en presente`, en español, con viñetas en el cuerpo).
   **No hagas el commit ni el push**: Benjamín los hace. Muéstrale los comandos:
   `git add -A`, `git commit` con el mensaje, y `git push` (todo se trabaja en `develop`).
8. Termina con un resumen de 3 líneas: qué se logró, qué quedó pendiente y cuál es la siguiente sesión.
