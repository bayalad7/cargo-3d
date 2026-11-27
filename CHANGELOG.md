# Changelog

Todos los cambios relevantes del proyecto. Formato basado en
[Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/); versiones con [SemVer](https://semver.org/lang/es/).
Cada versión liberada en `main` lleva su etiqueta de git (ADR-0012).

## [Sin publicar]

### Agregado
- Candidatas por orientación en el motor: `CargoItem.Candidates` da las orientaciones permitidas que caben en un punto, sin repetidas (piezas cúbicas), y `Section.Box()` la caja del tramo; capítulo 05 de la guía de Go (sesión 1.2.2).
- Geometría de cajas en el motor: `Box` con `Intersects`, `ContainedIn` y `Contact`, 46 casos límite y el concepto `geometria.md` con dibujos (sesión 1.2.1).
- Sprint 1.1 cerrado: contrato, tipos, lectura y validación de la entrada del motor, fixtures de los casos y capítulos 02 a 04 de la guía.
- 36 fixtures de 26 casos simulados en `engine/fixtures/casos/`, cada uno con su `.expect.json`, y una prueba que los valida con `DecodeInput` y revisa que sus cuentas y sus motivos de dimensión, puerta y peso cuadren con la entrada (sesión 1.1.5).
- Lectura y validación de la entrada del motor: `DecodeInput` y `EngineInput.Validate`, con errores por categoría y la ruta de cada campo, que juntan todos los problemas; rechazan `null`, llaves repetidas y campos obligatorios que faltan (sesión 1.1.4, ADR-0023).
- Capítulo de aprendizaje `04-errores-como-valores.md` y una prueba que compara los tipos y los rangos del motor con `contracts/engine.schema.json`.
- Tipos de salida del motor `Placement`, `LoadPlan`, `UnplacedItem` y `Metrics`; orientaciones con `RotatedDimensions` y `Orientations`; y pruebas que pasan los 27 ejemplos de salida de JSON a Go y de regreso y revisan sus medidas contra la orientación (sesión 1.1.3).
- Capítulo de aprendizaje `03-metodos-y-receptores.md`; `sistema-coordenadas.md` dice ahora dónde vive cada concepto en el código.
- Tipos del motor `EngineInput`, `Container`, `Section`, `Overhang`, `Door`, `CargoItem` y `Orientation` en `backend/internal/engine`, con `Volume()`, `Margin()` e `IsStackable()`, y una prueba que pasa los 27 ejemplos de entrada de JSON a Go y de regreso (sesión 1.1.2).
- Capítulo de aprendizaje `02-structs-y-tipos.md` y ADR-0022: la guía y los comentarios explican Go solo con Go.
- Sección `docs/src/desarrollo/` (ADR-0021): instalación del entorno, servicios, conexión a PostgreSQL con pgAdmin y psql (`entorno/base-de-datos.md`) y guía de ramas, commits y versiones (`flujo/ramas-y-commits.md`).
- Capítulo de aprendizaje `docs/src/aprendizaje/go/01b-modulos-e-imports.md`: por qué los imports `github.com/...` se leen del disco.
- Casos de uso simulados en `docs/src/casos/` (29 casos por categoría, incluida carga sobredimensionada con permiso de la NOM-040) y procedimiento de pruebas con cargas reales
  para la v0.2.0 (ADR-0020).
- Equipos personalizados por empresa: plantilla por clase y versiones de cada equipo (ADR-0019, se construye en la fase 2).
- Imagen `equipo.svg` en cada ejemplo de `contracts/examples/`, generada con `pnpm examples:render`.
- Catálogo canónico de equipos `contracts/equipment.json` (26 equipos, con el ejemplo de cada uno en `examples`) y su página
  `docs/src/catalogos/equipos.md`.
- Ejemplos nuevos: 10' Dry, 45' reefer pallet wide, flat rack 40' plegable y 20' bulk (27 en total).
- Índice de la bitácora (`docs/src/bitacora/README.md`).
- Ejemplos del contrato: 8 equipos más (23 en total) y 17 entradas inválidas en `contracts/examples/invalid/`.
- Sistema oficial de Voxa v5: diálogo, aviso, barra masiva, botón `danger`, estado en vivo y tres tokens.
- `scripts/contratos/check-examples.mjs` (`pnpm examples:check`): revisa las reglas físicas de los ejemplos.
- Conceptos: sistema de coordenadas, y equipos y normas (Tablas A-2 y B-2 de la NOM-012 leídas del DOF).
- Formatos para levantar datos en el patio (`docs/plantillas/patio/`) y sesión 2.1.6 (equipos con varios tramos).
- Registro de pendientes `docs/src/pendientes.md` (id `P-NNN`, quién, cuándo y cómo retomarlo): `/cerrar-sesion`
  lo actualiza en cada cierre y `/iniciar-sesion` y `/sesion` lo consultan.
- Contrato del motor `contracts/engine.schema.json` v0.1.0 (sprint 1.1): entrada y salida con milímetros y
  kilogramos enteros, tramos del equipo, puerta opcional, orientaciones, motivos de piezas no colocadas y margen de
  carga sobredimensionada (OOG). 15 ejemplos sintéticos por tipo de equipo en `contracts/examples/`, con la lista
  de lo que falta por cubrir.
- ADR-0015 (geometría del contrato), ADR-0016 (qué rechaza el motor y qué reporta) y ADR-0017 (carga
  sobredimensionada y open top).
- Sistema de diseño Voxa (ADR-0014) en `packages/ui` (`@cargo-3d/ui`): tokens en claro y oscuro con
  `light-dark()` y Tailwind CSS v4 (`tokens.css`, fuente única de verdad), componentes `vx-` (`components.css`),
  logo Vóxel en SVG (componente `Logo` y 6 archivos), íconos con Lucide (`Icon`, `appIcons`) y componentes
  React de presentación: `Button`, `StatusBadge`, `FilterChip`, `ConfirmDialog`, `Toast`/`Toaster` y
  `BulkActionBar`. Vista previa en `packages/ui/preview/`.
- Documentación del sistema de diseño (`docs/src/ux/sistema-diseno.md`) y reglas de diseño en `CLAUDE.md`.
- Estructura del monorepo (Go workspace, pnpm workspaces, Turborepo).
- Backend en Go 1.27 (motor y API en un binario, ADR-0002): servidor con `GET /health`, apagado ordenado y pruebas.
- Docker: `docker-compose.yml` (PostgreSQL 19, API con perfil), `backend/Dockerfile` multi-etapa
  (`golang:1.27` → `distroless/static-debian13`), `.env.example`.
- Contratos en `contracts/` (OpenAPI) y `engine/fixtures/`.
- `CLAUDE.md` con el contexto y las reglas del proyecto, incluidas las reglas para limpiar Docker después de probar.
- Documentación (mdBook): introducción, convenciones, conceptos, ADR 0001–0013,
  guía de Go, plan de trabajo y bitácora.
- ADR-0007: dónde vive el código de PostgreSQL (`backend/internal/store`).
- ADR-0011: versiones de todas las tecnologías (Node.js 26 LTS, pnpm 12, Turborepo 2.11, …).
- ADR-0012: todo se trabaja directo en `develop` (sin ramas de trabajo ni pull requests); `main` solo recibe versiones.
- Docker Compose con nombre de proyecto fijo (`cargo-3d`: un solo grupo con `db` y `api`), sin reinicio automático,
  `.dockerignore` más estricto y scripts `db:down`, `docker:api`, `docker:clean`.
- `scripts/servicios/levantar` y `detener` (Windows y macOS/Linux) y tareas de VS Code: PostgreSQL en Docker
  y la API con `go run` en un solo paso.
- PostgreSQL con ICU `es-MX` para ordenar texto en español.
- PostgreSQL de Docker en el puerto 5433 de la computadora, para no chocar con un PostgreSQL instalado en el sistema.
- Plan: instrucciones de Docker con un solo comando y un solo grupo `cargo-3d`.
- Soporte de desarrollo en Windows, macOS y Linux: `.gitattributes` (LF para todos), guía de instalación por
  sistema en el plan e imagen Docker para `linux/amd64` y `linux/arm64` (compilación cruzada de Go).
- ADR-0013 y GitHub Actions: `ci.yml` (flujo de ramas, `gofmt`, `go vet`, `go test -race`, build de la imagen)
  y `release.yml` (imagen en ghcr.io y GitHub Release con notas del CHANGELOG al subir `vX.Y.Z`).
- Biblioteca de prompts por sesión para las fases 1 a 3 (`docs/src/prompts.md`), con puntos de control y cierre
  en cada sprint y bloque, publicación de versiones (v0.1.0, v0.2.0, v1.0.0) y prompts de apoyo.
- `docs/datos/` fuera de git (datos de clientes) y fixtures reales anonimizados.
- Skills de Claude Code del proyecto en `.claude/skills/` (iniciar-sesion, sesion, cerrar-sesion, punto-de-control, revisar, explicar, nuevo-adr, cambiar-contrato, cerrar-sprint).
- Guía del flujo de trabajo con Claude Code (`docs/src/flujo-claude-code.md`).

### Cambiado
- PostgreSQL de desarrollo sin contraseña (`trust`, solo en `127.0.0.1`); producción siempre con contraseña (ADR-0024).
- CS-010 (módulo ancho solo `LWH`) y CS-026 (panel de 2 490 mm de ancho): corregidos al convertirlos en fixtures.
- Regla 6 de `CLAUDE.md`, skill `/explicar` y guía de Go (capítulos 00, 01 y 01b): se explica Go solo con Go, sin comparaciones con otros lenguajes (ADR-0022).
- La bitácora es un archivo por tema (`AAAA-MM-DD-tema.md`), agrupada por mes en el sitio.
- Las cargas reales del plan se sustituyen por casos simulados hasta la v0.2.0; el sprint 2.6 pasa a ser de pruebas
  con cargas reales.
- Ejemplos del contrato corregidos con fuente (reefers serie Hapag, cajas refrigeradas Alzper, chasis 1 219 mm).

### Corregido
- Cuatro ejemplos tenían una descripción de más de 200 caracteres, que el esquema no permite.

### Eliminado
- Skill `/punto-de-control` y los puntos de control del ciclo de sprints (ADR-0018).
- Tanque ISO 20' y chasis portacontenedor del catálogo de equipos: no son espacios de carga para el motor.
