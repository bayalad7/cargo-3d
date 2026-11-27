# ADR-0021: Organización de la documentación por propósito

- **Estado:** Aceptado
- **Fecha:** 2026-11-26

## Contexto

La documentación creció de forma orgánica en `docs/src/`. Lo que sirve para **programar** quedó repartido:
la instalación del entorno y el uso de Docker eran una sección larga dentro de `plan.md`, mientras que
`flujo-claude-code.md` y `convenciones.md` estaban sueltos en la raíz. Al llegar la necesidad de documentar
la conexión a PostgreSQL con pgAdmin no había un lugar obvio donde ponerla, y el plan mezclaba «qué se va a
construir» con «cómo se prepara la computadora».

Restricciones:

- mdBook solo publica lo que está dentro de `docs/src/`.
- `plan.md`, `prompts.md` y `pendientes.md` los leen las skills (`/iniciar-sesion`, `/sesion`,
  `/cerrar-sesion`, `/cerrar-sprint`) y el agente `investigador-logistico` por su ruta.

## Decisión

Cada carpeta de `docs/src/` responde **una** pregunta. Todo lo de desarrollo va en `docs/src/desarrollo/`:

| Carpeta o archivo | Pregunta que responde | Ejemplos |
|---|---|---|
| raíz: `introduccion.md`, `plan.md`, `prompts.md`, `pendientes.md` | ¿Qué se construye, en qué orden y qué falta? (gestión del proyecto) | Fases, sesiones, `P-NNN` |
| `desarrollo/` | ¿Cómo preparo la computadora y cómo trabajo día a día? | Instalación, servicios, base de datos, git, Claude Code, convenciones |
| `desarrollo/entorno/` | Preparar y operar el entorno local | `instalacion.md`, `servicios.md`, `base-de-datos.md` |
| `desarrollo/flujo/` | El proceso de trabajo | `claude-code.md`, `ramas-y-commits.md` |
| `conceptos/` | ¿Cómo funciona el producto por dentro? | Bin packing, coordenadas, normas |
| `catalogos/` | Datos de referencia | Equipos |
| `casos/` | ¿Qué debe resolver el motor? | Casos simulados, pruebas reales |
| `adr/` | ¿Por qué se decidió así? | Una decisión por archivo |
| `ux/` | ¿Cómo se ve y se comporta la interfaz? | Sistema de diseño Voxa |
| `aprendizaje/` | ¿Qué necesito aprender para entender el código? | Go explicado con el código real (ADR-0022) |
| `bitacora/` | ¿Qué pasó y cuándo? | Un archivo por tema |

Una página nueva se ubica por la pregunta que responde, se registra en `SUMMARY.md` y, si es de desarrollo, en
la tabla de `desarrollo/README.md`.

`prompts.md` se queda en la raíz junto a `plan.md`: es gestión del proyecto y lo leen las skills por su ruta.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Carpetas por propósito, `desarrollo/` con `entorno/` y `flujo/` (elegida)** | Lugar obvio para cada página nueva; el plan queda solo como plan | Mover archivos obliga a actualizar enlaces una vez |
| Dejar todo como estaba y agregar `base-de-datos.md` en la raíz | Sin mover nada | La raíz sigue creciendo sin orden; el entorno seguiría dentro del plan |
| `docs/desarrollo/` fuera de `docs/src/` | Separado del sitio | mdBook no lo publica; quedaría fuera del índice y de la búsqueda |
| Mover también `prompts.md` a `desarrollo/flujo/` | Todo el flujo junto | Cambiar seis archivos de `.claude/` sin ganar claridad; el prompt de cada sesión es parte del plan |

## Consecuencias

- **Positivas:** quien prepara una computadora nueva lee una sola sección; `plan.md` queda más corto;
  los procedimientos (pgAdmin, liberar una versión) tienen un lugar fijo.
- **Negativas / costos:** se actualizaron las rutas en `CLAUDE.md`, `README.md`, `SUMMARY.md`, las skills
  `cerrar-sesion` y `cerrar-sprint`, `prompts.md`, ADR-0009, el índice de la bitácora y los scripts de
  `scripts/servicios/`. Las entradas antiguas de la bitácora y el ADR-0018 conservan los nombres viejos
  (`flujo-claude-code.md`, `convenciones.md`) porque describen cómo estaban las cosas en su fecha.
- **Qué tendríamos que hacer para revertirla:** mover los archivos de `desarrollo/` de regreso a la raíz,
  devolver la sección «Preparación del entorno» a `plan.md` y revertir las rutas.
