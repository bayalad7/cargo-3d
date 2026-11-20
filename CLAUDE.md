# CLAUDE.md — Memoria del proyecto Cargo 3D

Este archivo es el punto de partida de cada sesión de trabajo con Claude.
Léelo completo antes de escribir código. Mantenlo actualizado al final de cada sesión.

## Qué es el proyecto

Cargo 3D es un software de **planeación de carga 3D** para contenedores marítimos,
cajas secas y plataformas, similar a EasyCargo pero enfocado en el mercado mexicano
(puerto de Manzanillo, Colima): equipos y reglas de carga mexicanas (NOM-012),
precio en pesos, soporte en español.

- **Dueño y mantenedor:** Benjamín (desarrollador; viene de PHP, está aprendiendo Go con este proyecto).
- **Primer cliente / laboratorio de pruebas:** Transportes Villa.

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

Principio rector: **Go para todo el backend, TypeScript para todas las pantallas.**

## Reglas de documentación (acuerdo con Benjamín — OBLIGATORIAS)

1. **Idioma:** identificadores de código (variables, funciones, tipos, paquetes, archivos de código) en **inglés**.
   Comentarios, documentación, ADRs, bitácora y mensajes al usuario en **español**.
2. **Go muy explicado al inicio:** comentar casi cada línea del backend, explicando
   conceptos de Go (paquetes, exportado con mayúscula, punteros, structs, interfaces,
   errores como valores, `defer`, goroutines, canales, `context`) la primera vez que aparecen.
   Reducir el detalle de forma gradual conforme Benjamín avance; el nivel actual se registra en "Estado".
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
8. **Bitácora:** al cerrar cada sesión, agregar `docs/src/bitacora/AAAA-MM-DD.md`
   (qué se hizo, decisiones, pendientes) y actualizar `CHANGELOG.md` y la sección "Estado".
9. **Puntos de control:** al terminar cada módulo, pedir a Benjamín que lo explique
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

## Plan de trabajo

El plan completo vive en `docs/src/plan.md`: fases, sprints, criterios de terminado
y el modelo sugerido para cada paso. Antes de empezar un sprint, léelo y confirma
con Benjamín en qué sprint estamos. No se pasa de fase sin cumplir su criterio de terminado.

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
apps/web | desktop | mobile
docs/                    mdBook: ADRs, conceptos, aprendizaje, bitácora
docker-compose.yml       PostgreSQL (y la API con --profile api)
```

## Comandos

```bash
cp .env.example .env                  # una vez; cambia la contraseña
docker compose up -d db               # PostgreSQL de desarrollo
go run ./backend/cmd/api              # API en http://localhost:8080/health
go test ./backend/...                 # pruebas
go vet ./backend/...                  # revisión estática
gofmt -l backend                      # archivos sin formato (debe salir vacío)
go doc ./backend/internal/engine      # documentación de un paquete
docker compose --profile api up --build   # API dentro de Docker
mdbook serve docs --open              # sitio de documentación
```

## Estado

- **Fase actual:** 0 — cimientos (estructura, Docker y documentación). Servidor con `/health`; motor vacío.
- **Nivel de comentarios Go:** máximo (casi cada línea).
- **Siguiente paso:** fase 1, sprint 1.1 (ver `docs/src/plan.md`) — modelos del motor y `contracts/engine.schema.json`.
- **Pendiente de Benjamín:** cargas reales de ejemplo, catálogo de equipos usados,
  reglas del patio.
