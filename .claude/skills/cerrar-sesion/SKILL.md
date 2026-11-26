---
name: cerrar-sesion
description: Cierra una sesión de trabajo de Cargo 3D. Corre las verificaciones, actualiza bitácora, registro de pendientes, CHANGELOG y el Estado de CLAUDE.md, y propone el mensaje de commit.
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
3. **Bitácora.** Un archivo **por tema** en `docs/src/bitacora/<fecha de hoy>-<tema>.md` (tema corto en minúsculas,
   sin acentos, con guiones; por ejemplo `2026-11-21-contrato-del-motor.md`). Si en la sesión se trabajaron dos temas
   distintos, son dos archivos. Cada archivo sigue el formato de `docs/src/bitacora/README.md`:
   - título `# <día> <mes abreviado> <año> · <Tema>` y la línea `**Sesión:** <número o «Fuera del plan»> · **Área:** <áreas de `desarrollo/flujo/ramas-y-commits.md`>`;
   - secciones "Qué se hizo", "Decisiones", "Verificación" y "Pendientes";
   - agrega una fila arriba de la tabla de `docs/src/bitacora/README.md` (fecha, tema, sesión, área, ids) y la entrada
     en `docs/src/SUMMARY.md` bajo el mes que le toca, la más reciente primero (crea el mes si no existe);
   - si el tema ya tiene archivo de hoy, agrega una sección nueva sin borrar lo anterior.
   Escribe solo lo que realmente pasó en esta sesión; incluye la nota de Benjamín si dio una.
   En su sección "Pendientes" **solo cita los ids** del registro (`P-012: …` en una línea); el detalle va en el paso 3b.
3b. **Registro de pendientes** (`docs/src/pendientes.md`), **obligatorio en cada cierre** aunque no haya nada nuevo:
   - todo lo que quedó abierto en la sesión (tareas sin terminar, datos que faltan, valores supuestos, decisiones
     aplazadas, hallazgos de `/revisar` no corregidos, trabajo hecho con herramientas temporales que no quedó en el
     repo) entra como fila nueva en «Abiertos», con el siguiente id `P-NNN` libre y todas las columnas llenas:
     quién, cuándo (sesión de `prompts.md` antes de la cual debe resolverse), qué falta, por qué importa,
     **cómo retomarlo** (comando, sesión y archivo exactos) y origen;
   - lo que no se puede hacer antes de cierta sesión pasa a «Agendados» **solo si** ya quedó escrito en el prompt de
     esa sesión; no cuenta como hecho;
   - los pendientes que se resolvieron en la sesión pasan a «Cerrados» con la fecha y cómo se resolvieron;
     nunca se borra una fila ni se reutiliza un id;
   - sin datos de clientes: si hace falta, apunta al archivo de `docs/datos/`.
   Antes de seguir, revisa la conversación completa: nada de lo pendiente puede quedar solo en el chat.
4. **CHANGELOG.md.** Agrega en "[Sin publicar]" lo que cambió para el producto o el proyecto
   (Agregado / Cambiado / Corregido). No repitas detalles internos.
5. **CLAUDE.md, sección "Estado".** Actualiza fase, sprint, última sesión terminada, siguiente sesión sugerida
   (número de `docs/src/prompts.md`) y pendientes de Benjamín. No cambies otras secciones sin preguntar.
6. **Plan.** Si una tarea del checklist del sprint en `docs/src/plan.md` quedó realmente terminada, márcala
   con `[x]`. No marques tareas a medias.
7. **Mensaje de commit.** Propón un mensaje con el formato de `docs/src/desarrollo/flujo/ramas-y-commits.md`
   (`área: descripción en presente`, en español, con viñetas en el cuerpo).
   **No hagas el commit ni el push**: Benjamín los hace. Muéstrale los comandos:
   `git add -A`, `git commit` con el mensaje, y `git push` (todo se trabaja en `develop`).
8. Termina con un resumen de 3 líneas: qué se logró, qué quedó pendiente (ids nuevos y cerrados de
   `docs/src/pendientes.md`, y cuáles le tocan a Benjamín) y cuál es la siguiente sesión.
