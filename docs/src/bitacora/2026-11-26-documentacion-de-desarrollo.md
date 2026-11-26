# 26 nov 2026 · Documentación de desarrollo y pgAdmin

**Sesión:** Fuera del plan · **Área:** docs, infra

## Qué se hizo

- **Aprendizaje de Go:** capítulo `aprendizaje/go/01b-modulos-e-imports.md`, que responde por qué los imports
  `github.com/bayalad7/cargo-3d/...` se leen del disco y no de GitHub (el nombre del módulo es solo un prefijo;
  comparado con PSR-4 de Composer). Enlazado desde la introducción, «Primeros pasos» y `SUMMARY.md`.
- **Sección `docs/src/desarrollo/`** (ADR-0021), con una carpeta por propósito:
  - `entorno/instalacion.md`: la sección «Preparación del entorno» que vivía dentro de `plan.md`, más un paso
    opcional para pgAdmin. En `plan.md` quedó un enlace.
  - `entorno/servicios.md`: script del día a día, todo dentro de Docker, reglas para no llenar el disco y el
    proyecto de pruebas `cargo3d-test`.
  - `entorno/base-de-datos.md`: instalar pgAdmin, registrar el servidor (`127.0.0.1:5433`, base y usuario
    `cargo3d`), cadena de conexión, `psql` dentro del contenedor, conexión a la base de pruebas con otro puerto
    y tabla de errores comunes.
  - `flujo/ramas-y-commits.md`: ciclo de un cambio, formato y áreas de los commits, cómo liberar una versión,
    versión de parche y qué hacer si algo sale mal.
  - `flujo/claude-code.md` y `convenciones.md` se movieron desde la raíz de `docs/src/`.
- **Rutas actualizadas** en `CLAUDE.md`, `README.md`, `SUMMARY.md`, `prompts.md`, ADR-0009, el índice de la
  bitácora, las skills `/cerrar-sesion` y `/cerrar-sprint` y los comentarios de `scripts/servicios/levantar.*`.

## Decisiones

- `docs/desarrollo/` pasó a `docs/src/desarrollo/`: mdBook solo publica lo que está dentro de `docs/src/`.
- `prompts.md` se queda en la raíz junto a `plan.md`: lo leen las skills y el agente por su ruta (ADR-0021).
- pgAdmin se instala como aplicación de escritorio, no como contenedor, para no sumar imágenes al disco.
- Las entradas antiguas de la bitácora y el ADR-0018 conservan los nombres viejos de los archivos movidos
  (`flujo-claude-code.md`, `convenciones.md`): describen cómo estaban las cosas en su fecha.

## Verificación

- `gofmt -l backend` sin salida; `go vet ./backend/...` y `go test ./backend/...` sin errores (no cambió código Go).
- `mdbook build docs` sin errores; revisión de todos los enlaces `.md` internos y de las anclas nuevas: ninguno roto.
- No se usó Docker en la sesión.
- **No se probó** la instalación de pgAdmin ni su conexión con PostgreSQL 19: ver P-024.

## Pendientes

- P-024: probar la instalación de pgAdmin y la conexión descritas en `entorno/base-de-datos.md`.
