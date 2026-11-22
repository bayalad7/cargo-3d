# CLAUDE.md — Memoria del proyecto Cargo 3D

Este archivo es el punto de partida de cada sesión de trabajo con Claude.
Léelo completo antes de escribir código. Mantenlo actualizado al final de cada sesión.

## Qué es el proyecto

Cargo 3D es un software de **planeación de carga 3D** para contenedores marítimos,
cajas secas y plataformas, similar a EasyCargo, 3D Load Calculator, CargoWizz,
Cargo Optimizer Enterprise, Cube-IQ, CubeMaster, Goodloading, Load!, PackApp,
PackVol y SeaRates pero enfocado en el mercado mexicano (puerto de Manzanillo, Colima):
equipos y reglas de carga mexicanas (NOM-012), precio en pesos, soporte en español.

Eliges (o creas) tu espacio de carga, metes tus datos (peso, medidas, restricciones)
y el motor acomoda todo, los resultados se muestran en una visualización 3D en tiempo real
y las piezas se pueden mover o girar a mano si hace falta.

El corazón es 3D bin packing (empaquetado en contenedores), que es NP-hard.
Nadie lo resuelve de forma exacta; se usan heurísticas:

- Extreme Points o Skyline/Layer building: colocan cajas en “esquinas” disponibles,
  ordenadas por volumen o por base. Son rápidas y dan buenos resultados.
- Metaheurísticas encima de eso (algoritmos genéticos, búsqueda tabú, simulated annealing)
  para mejorar el orden y la rotación de las piezas.
- Restricciones reales, que es donde está el valor: estabilidad (que nada quede “flotando”),
  apilabilidad, carga máxima por eje, centro de gravedad, orden de descarga (LIFO por destino),
  piezas que no se pueden voltear y peso máximo del contenedor (VGM).

## Stack (decidido; ver `docs/src/adr/`)

| Capa | Tecnología | ADR |
|---|---|---|
| Motor de optimización | Go (`backend/internal/engine`), corre en el servidor | 0002 |
| API | Go con la librería estándar (`backend/internal/httpapi`), mismo binario que el motor | 0006 |
| Base de datos | PostgreSQL con pgx y sqlc | 0007 |
| Interfaz web | TypeScript + React + React Three Fiber (`apps/web`) | 0003 |
| Escritorio | Wails v3 (`apps/desktop`), reutiliza `apps/web` | 0004 |
| Móvil (patio) | React Native + Expo (`apps/mobile`) | 0005 |
| Contrato | API-first con OpenAPI (`contracts/openapi.yaml`) | 0008 |
| Contenedores | Docker: PostgreSQL en desarrollo, imagen del backend en producción | 0010 |
| Organización | Monorepo: Go workspace + pnpm + Turborepo | 0001 |
| Versiones | Última estable de cada tecnología (Go 1.27, PostgreSQL 19, Node 26 LTS…) | 0011 |
| Sistema de diseño | Voxa: Tailwind CSS v4 + Lucide en `packages/ui` (`@cargo-3d/ui`) | 0014 |

Principio rector: **Go para todo el backend, TypeScript para todas las pantallas.**

## Reglas de documentación (Acuerdos OBLIGATORIOS)

1. **Idioma:** identificadores de código (variables, funciones, tipos, paquetes, archivos de código) en **inglés**.
   Comentarios, documentación, ADRs, bitácora y mensajes al usuario en **español**.
2. **Go muy explicado al inicio:** comentar casi cada línea del backend, explicando
   conceptos de Go (paquetes, exportado con mayúscula, punteros, structs, interfaces,
   errores como valores, `defer`, goroutines, canales, `context`) la primera vez que aparecen.
   Reducir el detalle de forma gradual conforme el usuario avance; el nivel actual se registra en "Estado".
3. **Nada se da por terminado sin documentar:** todo lo exportado lleva comentario de documentación
   (empieza con el nombre del elemento, como pide `go doc`) y TSDoc en TypeScript,
   con secciones "Cómo funciona" y "Por qué así" cuando aplique.
4. **Cada decisión importante → un ADR** nuevo en `docs/src/adr/` (usar `plantilla.md`) y
   enlazarlo en `docs/src/SUMMARY.md`.
5. **Cada concepto grande → documento en `docs/src/conceptos/`** con diagrama si ayuda.
6. **Conceptos nuevos de Go → `docs/src/aprendizaje/go/`**, explicados con las
   estructuras reales del proyecto y comparados con PHP, no con ejemplos genéricos.
7. **Pruebas como documentación:** nombres de prueba que describan el comportamiento
   (ej. `TestDoesNotStackFragileUnderHeavy`), con comentario en español. Preferir pruebas de tabla.
8. **Bitácora y pendientes:** al cerrar cada sesión, agregar `docs/src/bitacora/AAAA-MM-DD.md`
   (qué se hizo, decisiones, ids de pendientes) y actualizar `CHANGELOG.md` y la sección "Estado".
   **Todo lo que quede abierto va a `docs/src/pendientes.md`** (id `P-NNN`, quién, cuándo, qué falta, por qué
   importa y cómo retomarlo): nada pendiente se queda solo en la conversación. Lo resuelto pasa a «Cerrados».
9. **Puntos de control:** al terminar cada módulo, pedir al usuario que lo explique
   con sus palabras antes de continuar.
10. **El contrato va primero:** cualquier cambio de la API empieza en `contracts/openapi.yaml`;
    cualquier cambio de entrada o salida del motor, en `contracts/engine.schema.json`.

## Reglas de Go para este proyecto

- Solo librería estándar salvo que un ADR justifique la dependencia.
- Los errores siempre se revisan (`if err != nil`); nunca se ignoran con `_` sin comentario.
- El motor (`internal/engine`) no importa nada de HTTP ni de base de datos.
- Solo `internal/store` habla con PostgreSQL; `httpapi` lo usa, `engine` nunca.
- Configuración solo por variables de entorno; nada fijo en el código.
- Antes de dar algo por terminado: `gofmt`, `go vet ./backend/...` y `go test ./backend/...` sin errores.

## Reglas de Docker (OBLIGATORIAS — el disco duro puede llenarse rápido)

Todo lo que se cree en Docker para correr o probar una implementación (contenedores, imágenes
construidas o descargadas, volúmenes, redes y caché de build) **se elimina al terminar la prueba**,
pase o falle. Nada se queda "por si acaso": imágenes, volúmenes y caché de build se acumulan
y pueden llenar el disco duro rápidamente si no se limpian.

1. **Antes de probar**, anotar lo que ya existe para no borrar recursos ajenos a la prueba:
   `docker system df` y `docker image ls`.
2. **Probar en un proyecto de Compose aparte** (`-p cargo3d-test`), para no tocar el volumen
   `db-data` con la base de datos de desarrollo:
   ```bash
   docker compose -p cargo3d-test up -d db # o: --profile api up -d --build
   ```
3. **Al terminar, limpiar todo lo creado:**
   ```bash
   docker compose -p cargo3d-test --profile api down -v --rmi all --remove-orphans
   docker image rm postgres:19 golang:1.27 gcr.io/distroless/static-debian13:nonroot # solo si no existían antes
   docker builder prune -f # caché de build (se regenera sola)
   docker system df # confirmar que quedó como al inicio
   ```
4. **Reportar** en el mensaje final qué se creó, qué se eliminó y cuánto espacio se liberó.
5. Nunca usar `docker system prune -a --volumes` ni borrar recursos que no se crearon en la sesión
   sin preguntar antes.

## Reglas de diseño (ver ADR-0014 y `docs/src/ux/sistema-diseno.md`) — OBLIGATORIO

Todas las pantallas (web, escritorio y móvil) usan el sistema de diseño **Voxa** de `packages/ui` (`@cargo-3d/ui`).

- **Fuente única de verdad:** `packages/ui/src/styles/tokens.css` (valores) y `components.css` (clases `vx-`).
  Un valor se cambia ahí **y** en el sistema oficial de claude.ai (`https://claude.ai/artifact/1Daeb4Bq3esqWZXj3KNk3a`),
  en el mismo cambio. Nunca en una pantalla.
- **Nada escrito a mano:** ningún color hex, tamaño de letra, radio ni sombra en una pantalla. Solo tokens
  y utilidades de Tailwind v4 del sistema (`bg-canvas`, `text-ink`, `text-h1`, `rounded-md`, `p-4`…).
  Las paletas de Tailwind (`bg-blue-500`) están desactivadas a propósito. Color nuevo → token nuevo con función.
- **Colores por función, 60 · 30 · 10:** superficies `canvas`/`card` (60 %), tinta `ink` (30 %), amarillo
  `accent` #FFC400 (10 %). En claro el amarillo **solo es relleno** con texto `on-accent` (nunca texto, borde ni
  blanco encima); para líneas y texto de acento, `accent-line`. **Un solo botón `primary` por vista.**
- **Estados** (`success`, `progress`, `warning`, `danger`, `delayed`): nunca amarillo y nunca solo color
  (texto + punto o ícono). Contraste mínimo 4.5:1 en texto, revisado en claro y oscuro.
- **Temas:** `light-dark()` en los tokens; `data-theme="light"|"dark"` fuerza el tema, sin atributo sigue al
  sistema. Casi nunca hace falta `dark:`.
- **Tipografía:** Inter Tight (interfaz) e IBM Plex Mono (pesos, medidas, coordenadas). Solo la escala del
  sistema: `display`, `h1`, `h2`, `body`, `label`, `small`, `pill`, `data`.
- **Íconos:** Lucide con `<Icon>` (trazo 1.5, tamaños 14/16/18) y el vocabulario `appIcons`; un concepto nuevo
  se agrega a `appIcons`, no en la pantalla.
- **Logo:** componente `<Logo />` (cambia de tema solo) o los SVG de `packages/ui/src/brand/`. No girar,
  recolorear, sombrear ni deformar. Mínimo 20 px el logotipo y 24 px la insignia.
- **Componentes solo de presentación:** `@cargo-3d/ui` no importa la API, el motor ni estado global; recibe
  propiedades y avisa con funciones. Cambiar el aspecto de una pantalla nunca toca su lógica, hooks, manejadores
  de eventos ni conexiones.
- **Voz:** español de México, directo; botones con verbo + objeto («Optimizar carga»); `24 310 kg`, `67,4 m³`.
- «Voxa» es el nombre comercial en la interfaz; el proyecto y los paquetes siguen como `cargo-3d` (ADR-0014).

## Datos de clientes (OBLIGATORIO)

- Los datos reales (cargas, catálogos, fotos, nombres de clientes) viven en `docs/datos/`, que **no se sube a git**.
- Lo que sí se sube (`engine/fixtures/real/`, ejemplos, pruebas, documentación) va **anonimizado**: sin nombres
  de clientes ni destinatarios, direcciones, placas ni números de pedido; solo medidas, pesos, cantidades y restricciones.
- Nunca copiar un archivo de `docs/datos/` fuera de esa carpeta tal cual.

## Ramas (ver ADR-0012) — OBLIGATORIO

El proyecto lo trabaja **una sola persona**, así que el flujo es simple:

- **Todo se trabaja directo en `develop`.** No se crean ramas de trabajo (`area/descripcion`, `release/*`,
  `hotfix/*`) ni pull requests. Cada cambio es un commit en `develop` y se sube con `git push`.
- **`main` solo recibe versiones oficiales:** cuando se libera una versión, se une `develop` a `main`
  y se etiqueta `vX.Y.Z` (eso dispara `release.yml`). Nunca se hacen commits directos en `main`.
- `main` es la rama por defecto de GitHub. Toda la configuración de CI y publicación vive en
  `.github/workflows/ci.yml` y `.github/workflows/release.yml` (ADR-0013); no en la configuración de GitHub.
- Si un prompt, skill o documento menciona crear una rama o un pull request, esta sección manda.

## Plan de trabajo

El plan completo vive en `docs/src/plan.md`: fases, sprints, criterios de terminado
y el modelo sugerido para cada paso. Antes de empezar un sprint, léelo y confirma
con el usuario en qué sprint estamos. No se pasa de fase sin cumplir su criterio de terminado.

## Cómo se trabaja con Claude Code

- Cada sesión de trabajo tiene su prompt en `docs/src/prompts.md`, identificado por número (`1.1.2`).
- Comandos del proyecto (skills en `.claude/skills/`): `/iniciar-sesion`, `/sesion <número>`, `/revisar`,
  `/cerrar-sesion`, `/punto-de-control <tema>`, `/cerrar-sprint <sprint>`, `/explicar <algo>`,
  `/nuevo-adr <título>`, `/cambiar-contrato <cambio>` y `/buscador <sprint o tema>` (lanza el agente
  `.claude/agents/investigador-logistico.md`: investiga datos logísticos con fuentes, no escribe archivos;
  es **opcional**, solo cuando el sprint necesita datos que todavía no se tienen).
  Guía completa en `docs/src/flujo-claude-code.md`.
- Siempre se propone un plan y se espera la aprobación del usuario antes de escribir código.
- Claude nunca hace commit ni push: propone el mensaje y el usuario los ejecuta.
- Si cambias una sesión de `prompts.md`, conserva el formato de encabezado `### <número> — <título>`.

## Estructura

```
backend/                 Go — módulo github.com/bayalad7/cargo-3d/backend
  cmd/api/               punto de entrada del servidor (binario "api")
  internal/engine/       motor: modelos, algoritmo, restricciones (sin HTTP ni BD)
  internal/httpapi/      rutas y respuestas HTTP
  internal/store/        acceso a PostgreSQL (sprint 2.3; ver ADR-0007): único paquete que importa pgx
    migrations/          cambios al esquema en SQL, versionados
    queries/             consultas SQL escritas a mano (entrada de sqlc)
    sqlcgen/             Go generado por sqlc — no editar a mano
  sqlc.yaml              configuración de sqlc (sprint 2.3)
  Dockerfile             imagen de producción
contracts/               openapi.yaml y engine.schema.json (independientes del lenguaje)
engine/fixtures/         casos de prueba del motor en JSON
packages/types           TS — tipos generados desde los contratos
packages/api-client      TS — cliente generado desde OpenAPI
packages/viewer-3d       TS — visor 3D reutilizable
packages/validator       TS — validador independiente de planes
packages/ui              TS — sistema de diseño Voxa (ADR-0014): tokens, Tailwind v4, logo, íconos y componentes
  src/styles/            tokens.css (fuente única), components.css (clases vx-), fonts.css, index.css (entrada)
  src/brand/             logo Vóxel: Logo.tsx y los SVG
  preview/               vista previa en HTML (claro y oscuro)
apps/web | desktop | mobile
docs/                    mdBook: plan, prompts, ADRs, conceptos, aprendizaje, bitácora; plantillas/ (formatos vacíos para el patio)
.claude/skills/          comandos de Claude Code del proyecto
.claude/agents/          subagentes del proyecto (investigador-logistico)
docker-compose.yml       PostgreSQL (y la API con --profile api)
scripts/contratos/        check-examples.mjs: reglas físicas de los ejemplos (provisional, hasta la sesión 1.4.2)
scripts/servicios/        levantar y detener el entorno de desarrollo (.ps1/.cmd Windows, .sh/.command macOS/Linux)
.vscode/tasks.json       tareas "Levantar servicios" (Ctrl+Shift+B) y "Detener servicios"
.github/workflows/       ci.yml (revisión de cada cambio) y release.yml (publicación con etiqueta vX.Y.Z)
```

## Comandos

```bash
cp .env.example .env                  # una vez; cambia la contraseña
docker compose --profile api up -d --build   # grupo cargo-3d: PostgreSQL (127.0.0.1:5433) + API
docker compose --profile api down     # apagar el grupo (el volumen se conserva)
docker compose up -d db               # solo PostgreSQL, para trabajar con go run
scripts/servicios/levantar.ps1 | .sh  # PostgreSQL + go run en un paso (Ctrl+Shift+B en VS Code); detener.ps1 | .sh apaga
go run ./backend/cmd/api              # API en http://localhost:8080/health
go test ./backend/...                 # pruebas
go vet ./backend/...                  # revisión estática
gofmt -l backend                      # archivos sin formato (debe salir vacío)
go doc ./backend/internal/engine      # documentación de un paquete
docker compose --profile api up --build   # API dentro de Docker
node scripts/contratos/check-examples.mjs   # revisa contracts/examples/ (o pnpm examples:check)
mdbook serve docs --open              # sitio de documentación
pnpm --filter @cargo-3d/ui lint       # revisa los tipos del sistema de diseño
pnpm --filter @cargo-3d/ui build:preview   # compila la vista previa (abrir packages/ui/preview/index.html)
```

## Estado

- **Fase actual:** 1 — prueba de concepto, sprint 1.1 (sesión 1.1.1 terminada el 2026-11-21). Servidor con `/health`; motor vacío.
- **Nivel de comentarios Go:** máximo (casi cada línea).
- **Contrato del motor:** `contracts/engine.schema.json` v0.1.0, 23 ejemplos válidos y 17 inválidos en `contracts/examples/` (ADR-0015, 0016 y 0017).
  Escrito con datos sintéticos; la sesión 1.1.5 lo ajusta con las cargas reales.
- **Siguiente paso:** sesión 1.1.2 (`Container` y `CargoItem` en Go, siguiendo el contrato). Falta el punto de control de
  la 1.1.1 (explicar los campos del contrato).
- **Sistema de diseño:** base lista en `packages/ui` (2026-11-21, fuera del plan); se conecta a `apps/web` en el sprint 1.6.
- **Pendientes:** registro completo en `docs/src/pendientes.md` (P-001 a P-019). Los de Benjamín más urgentes:
  P-019 (copiar `docs/datos/` y `.env` antes de cambiar de computadora), P-007 (punto de control, antes de 1.1.2)
  y P-001 (cargas reales y reglas del patio, antes de 1.1.5).
