# Plan de trabajo

> Versión del 2026-11-20. Si cambias el plan, actualiza este archivo en el mismo pull request
> y regístralo en la bitácora.

## Cómo usar esta guía

Este plan lleva Cargo 3D de los cimientos a un producto que cobra, en cuatro fases con criterios claros
de terminado, trabajando desde tu computadora con Claude Code.

**Regla de oro: no se pasa a la siguiente fase sin cumplir el criterio de terminado de la actual.**

1. **Preparación del entorno** y **Cómo trabajar con Claude Code**: una sola vez, antes de programar.
2. **Qué modelo usar en cada paso**: consúltala al empezar cada sprint.
3. **Fases 1 a 3**: objetivo, sprints de 1 a 2 semanas, entregables y criterio de terminado.
4. **Qué debes tener en cuenta** y **Riesgos**: antes de empezar y al cerrar cada fase.
5. **Después del plan**: qué sigue con el producto en el mercado.

Las duraciones suponen **10 a 15 horas por semana**. Si dedicas más o menos, ajusta las fechas, no el alcance.

## Preparación del entorno

El proyecto se desarrolla igual en **Windows, macOS y Linux** (también en Mac con Apple Silicon).
Elige la columna de tu sistema; los comandos del proyecto (`go`, `pnpm`, `docker compose`) son los mismos en los tres.

| Herramienta | Windows (nativo) | macOS | Linux |
| --- | --- | --- | --- |
| Gestor de paquetes | `winget` (viene con Windows) | [Homebrew](https://brew.sh) | el de tu distribución |
| Git | `winget install Git.Git` (incluye Git Bash) | `xcode-select --install` o `brew install git` | `apt install git` |
| Go | `winget install GoLang.Go` | `brew install go` | binario de go.dev/dl |
| Docker | Docker Desktop (motor WSL2) | Docker Desktop u OrbStack (más ligero) | Docker Engine + plugin Compose |
| Node.js | `winget install Schniz.fnm`, luego `fnm install 26` | `brew install fnm`, luego `fnm install 26` | fnm o nvm |
| mdBook | `winget install Rustlang.mdBook` | `brew install mdbook` | binario de sus releases en GitHub |

Notas:

- **Windows:** se trabaja nativo con PowerShell o Git Bash; WSL2 es opcional (Docker Desktop lo usa por dentro).
  Si la ruta de tu usuario tiene espacios o acentos, clona el repo en una ruta corta como `C:\cargo-3d`.
- **Fin de línea:** `.gitattributes` y `.editorconfig` fijan LF para todos. Si clonaste antes de que existiera
  `.gitattributes`, corre una vez `git rm -r --cached . -q && git reset --hard` (con todo ya en commit).
- **Mac con Apple Silicon:** las imágenes de Docker del proyecto existen para ARM (`arm64`) y corren nativas.

### 1. Git y GitHub

1. Crea en GitHub el repositorio vacío `cargo-3d` (privado).
2. Entra a la carpeta `cargo-3d` y súbelo: `git remote add origin https://github.com/bayalad7/cargo-3d.git` y `git push -u origin main`.
3. Crea la rama `develop` desde `main` (`main` sigue como rama por defecto). El CI y la publicación
   ya están en `.github/workflows/` (ADR-0012 y ADR-0013).
4. Opcional: en Settings → Rules marca los checks del CI como obligatorios para que un CI en rojo bloquee el merge.
5. Si usas otro usuario u otro nombre de repo, cambia la ruta del módulo en `backend/go.mod` y en los imports (búscala con `grep -r bayalad7`).

### 2. Go

1. Instala Go 1.27 o superior desde go.dev/dl.
2. Verifica: `go version`.

### 3. Docker

Docker Desktop (Windows o macOS), OrbStack (macOS) o Docker Engine (Linux), con Docker Compose.

1. En la raíz del repo, copia `.env.example` a `.env` y cambia la contraseña
   (PowerShell: `Copy-Item .env.example .env`; macOS/Linux: `cp .env.example .env`).
2. Levanta todo con **un solo comando**, siempre desde la raíz del repo:

   ```bash
   docker compose --profile api up -d --build
   ```

   En Docker Desktop aparece **un solo grupo** con dos contenedores:

   ```
   ▾ cargo-3d
       ├─ cargo-3d-db-1    postgres:19    (127.0.0.1:5433)
       └─ cargo-3d-api-1   cargo-3d-api   (127.0.0.1:8080)
   ```

3. Apaga todo con `docker compose --profile api down` (los datos se conservan en el volumen `cargo-3d_db-data`).

Reglas para no crear contenedores, imágenes ni volúmenes de más:

- **Usa siempre `docker compose` desde la raíz del repo**, nunca `docker run` ni `docker build` sueltos:
  esos crean contenedores e imágenes fuera del grupo `cargo-3d`.
- El nombre del grupo está fijo en `docker-compose.yml` (`name: cargo-3d`), así que es el mismo aunque
  la carpeta se llame distinto o estés en otra computadora; volver a correr el comando reutiliza lo que ya existe.
- Incluye siempre `--profile api` al levantar **y** al apagar; sin él, el contenedor de la API queda fuera del comando.
- Los contenedores no arrancan solos al abrir Docker: solo corren cuando los levantas.
- **Puerto 5433:** PostgreSQL de Docker usa el 5433 de tu computadora para no chocar con un PostgreSQL
  instalado en el sistema (que ocupa el 5432). Si el 5433 también está ocupado, cambia `POSTGRES_PORT` en `.env`.

### 4. Node.js y pnpm (desde el sprint 1.4: validador en TypeScript)

Node 26 LTS con fnm o nvm. Desde Node 25, Corepack ya no viene incluido: instala pnpm con
`npm install -g pnpm@12`. pnpm lee el campo `packageManager` de `package.json` y usa
exactamente la versión del proyecto.

### 5. mdBook (sitio de documentación)

`winget install Rustlang.mdBook` en Windows, `brew install mdbook` en macOS o el binario de sus releases en GitHub en Linux.
Después de instalarlo, **abre una terminal nueva** para que reconozca el comando `mdbook`.

### 6. Editor

| Extensión de VS Code | Para qué |
| --- | --- |
| Go (oficial, usa gopls) | Autocompletado, errores en vivo, pruebas y depuración. Indispensable. |
| Error Lens | Errores en la misma línea. |
| Docker | Administrar contenedores desde el editor. |
| ESLint y Prettier | TypeScript (desde la fase 1). |
| YAML | Editar `openapi.yaml` y `docker-compose.yml` con validación. |

### 7. Claude Code

Instálalo con la guía oficial en [docs.claude.com](https://docs.claude.com), abre una terminal en la raíz
del repo y ejecuta `claude`. Lee `CLAUDE.md` en automático.

### 8. Verificación

| Comando | Debe mostrar |
| --- | --- |
| `go test ./backend/...` | `ok` en `engine` y en `httpapi` |
| `go vet ./backend/...` | Nada (sin errores) |
| `docker compose --profile api up -d --build` y `docker compose --profile api ps` | `cargo-3d-db-1` como `healthy` y `cargo-3d-api-1` como `running` |
| Abrir `localhost:8080/health` | `{"status":"ok","engineVersion":"0.1.0"}` (la API desde Docker) |
| `docker compose --profile api down` | Los dos contenedores se detienen y se eliminan; el volumen se conserva |
| `mdbook serve docs --open` | El sitio de documentación |

Mientras programas el backend conviene otra forma, dentro del **mismo grupo** `cargo-3d`.
**Un script lo hace todo** (abre Docker si está cerrado, crea `.env` si falta, apaga el contenedor de la API,
levanta PostgreSQL, espera a que esté listo y corre la API con `go run`; al detenerla con `Ctrl+C`, apaga PostgreSQL):

| Cómo | Windows | macOS / Linux |
| --- | --- | --- |
| VS Code | `Ctrl+Shift+B` (tarea "Levantar servicios") | `Cmd+Shift+B` / `Ctrl+Shift+B` |
| Doble clic | `scripts\servicios\levantar.cmd` | `scripts/servicios/levantar.command` (Finder) |
| Terminal | `.\scripts\servicios\levantar.ps1` | `./scripts/servicios/levantar.sh` |
| Dejar PostgreSQL encendido al salir | `levantar.ps1 -DejarDb` | `levantar.sh --dejar-db` |
| Apagar todo | `detener.ps1` / `detener.cmd` o la tarea "Detener servicios" | `detener.sh` / `detener.command` |

En macOS, la primera vez que hagas doble clic en un `.command`, macOS puede pedir permiso:
clic derecho → **Abrir**. Lo que hace el script por dentro, paso a paso:

- `docker compose up -d db` levanta solo PostgreSQL (el mismo contenedor `cargo-3d-db-1`, no uno nuevo)
  y la API se corre con `go run ./backend/cmd/api`, que recompila en segundos y se depura desde el editor.
- No mezcles las dos al mismo tiempo: `go run` y el contenedor de la API usan el puerto 8080.
  Antes de `go run`, apaga el grupo con `docker compose --profile api down`.
- Las dos formas usan el mismo volumen `cargo-3d_db-data`: la base de datos es la misma.
- Si quieres liberar la imagen de la API y la caché de construcción: `docker compose --profile api down --rmi local`
  y `docker builder prune -f`.

## Cómo trabajar con Claude Code

Tú diriges y revisas; Claude Code escribe, documenta y prueba. Cada sesión tiene una meta pequeña,
vive en su propia rama (creada desde `develop`) y termina con la bitácora actualizada.

### Ciclo de una sesión

1. **Arranque.** `git switch develop`, `git pull` y `git switch -c motor/modelos-basicos`; después `claude`. Pide: *"Lee CLAUDE.md y la última bitácora y dime dónde vamos."*
2. **Plan.** Modo de planeación (`Shift+Tab`) para todo lo que no sea trivial.
3. **Contrato primero.** Si cambia la API o el motor, se edita `contracts/` antes que el código.
4. **Implementación en pasos chicos**, cada uno con pruebas de tabla y documentación.
5. **Verificación:** `gofmt -l backend` vacío, `go vet ./backend/...` y `go test ./backend/...` sin errores.
6. **Revisión tuya:** lee el diff completo; pregunta lo que no entiendas: *"Explícame esta línea como si viniera de PHP."*
7. **Punto de control:** si se terminó un módulo, explícalo con tus palabras.
8. **Cierre:** *"Actualiza la bitácora, el CHANGELOG y la sección Estado de CLAUDE.md."* Commit, push y pull request hacia `develop`.

### Reglas que no se negocian

- **No aceptes código que no entiendes.**
- **Nada sin pruebas ni documentación.**
- **Ninguna dependencia nueva sin ADR.**
- **Sesiones cortas** y `/clear` entre tareas distintas.
- **El repo es la memoria.**

## Qué modelo usar en cada paso

Regla general: **Opus para pensar, Sonnet para construir, Haiku para lo mecánico**, y Fable solo cuando
un problema difícil se atore. Se cambia con `/model opus`, `/model sonnet`, `/model haiku`, `/model fable`
u `/model opusplan` (Opus planea, Sonnet ejecuta), y el esfuerzo con `/effort low|medium|high|xhigh|max`.
Fable puede cobrarse aparte con créditos de uso según tu plan.

| Paso | Modelo | Por qué |
| --- | --- | --- |
| Preparación del entorno | `sonnet` / `medium` | Comandos e instalación |
| 1.1 Modelos y contrato del motor | `opus` / `medium` | Define la base y los primeros capítulos de Go |
| 1.2 Geometría | `opus` / `high` | Los casos límite son donde nacen los bugs |
| 1.3 Puntos extremos | `opusplan` / `high` | Diseño del algoritmo con Opus; implementación con Sonnet |
| 1.4 Validador TS y banco de pruebas | `opus` / `high` | Es la red de seguridad |
| 1.5 Endpoint de optimización | `opusplan` / `medium` | Contrato y diseño con Opus; handlers con Sonnet |
| 1.6 Visor web | `sonnet` / `medium` | React y Three.js |
| 1.7 CI, Docker y demo | `sonnet` / `medium` | Configuración |
| 2.1 Restricciones (interfaces) | `opusplan` / `high` | Las restricciones interactúan entre sí |
| 2.2 Metaheurística con goroutines | `opus` / `xhigh`, o `fable` si se estanca | Lo más difícil del proyecto; además, concurrencia |
| 2.3 PostgreSQL, sqlc y autenticación | `opusplan` / `high` | Modelo de datos y seguridad |
| 2.4 La web como producto | `sonnet` / `medium` | Pantallas |
| 2.5 PDF del plan | `sonnet` / `medium` | Maquetación |
| 2.6 Piloto y despliegue con Docker | `sonnet` / `high` | Sin margen de error |
| 3.1 Multiempresa y permisos | `opus` / `high` | Un error expone datos entre clientes |
| 3.2 App móvil | `opusplan` / `medium` | Sincronización sin conexión con Opus |
| 3.3 Cobro y CFDI | `opus` / `high` | Dinero y obligaciones fiscales |
| 3.4 Seguridad y operación | `opus` / `max` en la revisión de seguridad | Buscar vulnerabilidades |
| 3.5 Escritorio con Wails | `sonnet` / `medium` | Reutiliza la web y el motor |
| 3.6 Salida al mercado | `opus` / `medium` | Textos y documentos legales (que revisa un abogado) |
| Explicarte un error o concepto de Go | `opus` / `medium` | |
| Revisar un diff antes del PR | `opus` / `high` | |
| Bitácora, CHANGELOG, Estado | `sonnet` / `low` | |
| Commits, formato, renombres | `haiku` / `low` | |

## Hoja de ruta

Del arranque de la fase 1 al lanzamiento: aproximadamente de 8 a 11 meses con 10 a 15 horas por semana.

```
Fase 0 · Cimientos ............................................ Terminada
   │  Backend Go con /health, Docker, contratos, 13 ADR, CI y guía de Go.
   ◇  Puerta: entorno verificado y cargas reales de Transportes Villa reunidas
   ▼
Fase 1 · Prueba de concepto ........................ SIGUIENTE · 6 a 8 semanas
   │  Motor de puntos extremos en Go, validador TS, endpoint y visor 3D web.
   ◇  Puerta: una carga real acomodada en menos de 2 s y sin errores del validador
   ▼
Fase 2 · MVP y piloto ........................................ 10 a 14 semanas
   │  Restricciones, metaheurística paralela, PostgreSQL, PDF y 4 semanas en Villa.
   ◇  Puerta: 4 semanas de uso real sin volver al método anterior
   ▼
Fase 3 · Producto y lanzamiento .................................. 4 a 6 meses
   │  Multiempresa, app de patio, cobro con CFDI y beta con empresas del puerto.
   ◇  Puerta: 3 empresas, además de Transportes Villa, pagando suscripción
   ▼
Fase 4 · Crecimiento ................................................ Continua
```

## Fase 1 — Prueba de concepto (6 a 8 semanas)

Al terminar, la web manda una carga real de Transportes Villa a la API, el motor en Go la acomoda
en segundos y la web la muestra en 3D.

### Sprint 1.1 — Modelos y contrato del motor

- [ ] `contracts/engine.schema.json`: entrada y salida del motor.
- [ ] Tipos `Container`, `CargoItem`, `Orientation`, `Placement` y `LoadPlan` en `internal/engine`, en milímetros y kilogramos (`int64`).
- [ ] Validación de entrada que devuelve errores claros (`error`), nunca `panic`.
- [ ] Primeras cargas reales en `engine/fixtures/real/` y catálogo de equipos.
- [ ] Guía de Go: capítulos 02 a 04.
- [ ] Punto de control: explicas los modelos con tus palabras.

### Sprint 1.2 — Geometría

- [ ] Cajas alineadas a los ejes: intersección, contención y contacto entre caras.
- [ ] Las 6 orientaciones y el filtro de orientaciones permitidas.
- [ ] Pruebas de tabla con casos límite y **fuzzing nativo de Go** (`go test -fuzz`).
- [ ] Guía de Go: capítulos 05 y 06.

### Sprint 1.3 — Primer algoritmo: puntos extremos

- [ ] Ordenamiento de piezas como estrategia intercambiable.
- [ ] Colocación por puntos extremos con la mejor posición y orientación.
- [ ] Restricción de peso máximo.
- [ ] Métricas: ocupación de volumen y de peso, piezas que no cupieron.
- [ ] Concepto del algoritmo documentado con diagrama.

### Sprint 1.4 — Validador independiente y banco de pruebas

- [ ] `packages/validator` en TypeScript: sin traslapes, todo dentro del equipo, peso respetado.
- [ ] Banco de pruebas con BR1 a BR15 en `engine/fixtures/br/` y las cargas reales.
- [ ] Benchmarks con `testing.B` y línea base de ocupación registrada en la bitácora.
- [ ] Guía de Go: capítulos 07 y 08.

### Sprint 1.5 — Endpoint de optimización

- [ ] `POST /plans/optimize` primero en `contracts/openapi.yaml`, luego en `internal/httpapi`.
- [ ] Límites de entrada (número de piezas, tamaño del cuerpo) y tiempo máximo de cálculo con `context`.
- [ ] `packages/types` y `packages/api-client` generados desde los contratos.
- [ ] Guía de Go: capítulo 10.

### Sprint 1.6 — Visor web mínimo

- [ ] `apps/web` con Vite, React y React Three Fiber.
- [ ] Captura de equipo y piezas (formulario y pegar desde Excel), llamada a la API.
- [ ] Visor 3D en `packages/viewer-3d`.

### Sprint 1.7 — Integración continua, Docker y demo

- [x] GitHub Actions para Go y Docker: `ci.yml` (`gofmt`, `go vet`, `go test -race`, build de la imagen) y `release.yml` (ADR-0013).
- [ ] Job de la web en `ci.yml` (lint, pruebas y build con pnpm).
- [ ] Demo publicada: backend en un contenedor y web en hosting estático.

### Criterio de terminado

1. Una carga real de Transportes Villa se acomoda y se muestra en menos de 2 segundos (incluida la red).
2. El validador no encuentra errores en ningún plan del banco de pruebas.
3. La ocupación en BR quedó medida y los acomodos reales le parecen razonables a quien carga en el patio.
4. Tú puedes explicar el algoritmo y el flujo completo, de la web al motor y de regreso.

## Fase 2 — MVP y piloto (10 a 14 semanas)

Al terminar, Transportes Villa planea sus cargas reales con Cargo 3D cada semana.

### Sprint 2.1 — Restricciones del patio

- [ ] Apoyo y estabilidad, apilabilidad, orden de descarga por destino.
- [ ] Reparto de peso por eje según la NOM-012 vigente (verifica la versión antes de programarla).
- [ ] Cada restricción es una **interfaz** de Go intercambiable, con pruebas y su sección en conceptos.

### Sprint 2.2 — Metaheurística en paralelo

- [ ] Algoritmo genético o búsqueda local sobre los puntos extremos.
- [ ] Varios intentos en paralelo con goroutines, con tiempo límite por `context`.
- [ ] Comparación contra la línea base de la fase 1.
- [ ] Guía de Go: capítulo 09.

### Sprint 2.3 — PostgreSQL y autenticación

- [ ] Paquete `backend/internal/store` con `migrations/`, `queries/` y `sqlcgen/` (estructura en ADR-0007).
- [ ] Migraciones SQL (herramienta elegida en un ADR) y consultas con sqlc + pgx.
- [ ] Empresas, usuarios, catálogos y planes (JSONB).
- [ ] Autenticación y separación de datos por empresa, con pruebas.
- [ ] Guía de Go: capítulo 11.

### Sprint 2.4 — La web como producto

- [ ] Importación de Excel con mapeo de columnas.
- [ ] Editor manual con validación en vivo (`packages/validator`).
- [ ] Historial, duplicar planes y varios equipos por carga.

### Sprint 2.5 — El plan en papel

- [ ] PDF del plan (vista 3D, paso a paso por capas, lista de piezas, pesos por eje). Dónde se genera (web o Go) se decide en un ADR.

### Sprint 2.6 — Piloto con Transportes Villa

- [ ] Despliegue con Docker Compose en un VPS: API, PostgreSQL, HTTPS y respaldos diarios.
- [ ] Capacitación de 1 hora y 4 semanas de uso real, midiendo tiempo antes y después.

### Criterio de terminado

1. Cuatro semanas de uso real sin volver al método anterior.
2. Ningún plan fue imposible de cargar.
3. Números reales para vender: tiempo ahorrado y ocupación lograda.
4. Restauraste la base de datos al menos una vez desde un respaldo.

## Fase 3 — Producto y lanzamiento (4 a 6 meses)

Al terminar, empresas además de Transportes Villa pagan cada mes.

- **3.1 Multiempresa y permisos:** roles, invitaciones, registro de actividad, límites por plan.
- **3.2 App móvil de patio (Expo):** plan del día, paso a paso, fotos de evidencia, modo sin conexión.
- **3.3 Cobro y CFDI:** suscripciones en pesos (ADR comparando Stripe, Mercado Pago y Conekta) y CFDI por un PAC.
- **3.4 Operación:** monitoreo de errores, ambientes separados, revisión de seguridad, plan de recuperación probado.
- **3.5 Escritorio con Wails (solo si un cliente lo pide):** reutiliza la web y corre el motor local sin conexión.
- **3.6 Salida al mercado:** página de producto, aviso de privacidad y términos revisados por un abogado, beta con 5 a 10 empresas del puerto, soporte por WhatsApp con horario.

### Criterio de terminado

1. Al menos 3 empresas, además de Transportes Villa, pagando.
2. Un mes sin caídas que afecten a un cliente.
3. La app de patio en uso real en al menos un cliente.
4. Cobro y facturación sin intervención manual.

## Qué debes tener en cuenta

### Técnico

- **El motor es el producto.** Más tiempo al motor que a las pantallas.
- **El validador manda.** Ningún plan sale sin pasar por él.
- **El contrato va primero.** Cambiar `contracts/` antes que el código evita que backend y pantallas se desincronicen.
- **El servidor paga el cálculo.** Pon límites de piezas y de tiempo por petición desde el sprint 1.5.
- **Datos reales antes que inventados.** Cada restricción se prueba primero con una carga real.
- **Unidades fijas:** milímetros y kilogramos enteros en todo el backend.

### Aprendizaje

- Go se aprende rápido, pero **la concurrencia (goroutines y canales) merece calma**: es la parte donde
  aparecen los errores difíciles. Se introduce hasta el sprint 2.2, con pruebas que usen `go test -race`.
- Escribe tú una función pequeña por sprint, con Claude Code como revisor.

### Negocio, legal y seguridad

- Vende el resultado (tiempo y espacio ahorrados), cobra desde la beta, precio mensual en pesos por empresa.
- Registra la marca ante el IMPI, aviso de privacidad, términos que dejen claro que el plan es una sugerencia.
- Secretos solo en `.env`; verificación en dos pasos en GitHub, hosting y pasarela de pagos.

### Tiempo

- Bloquea horas fijas a la semana y protégelas; sesiones cortas y frecuentes.
- Usa las vacaciones escolares para los sprints pesados.

## Riesgos y cómo mitigarlos

| Riesgo | Señal de alerta | Mitigación |
| --- | --- | --- |
| Acomodos imposibles de cargar | El patio cambia el plan sobre la marcha | Probar con cargas reales; convertir cada cambio del patio en restricción |
| Errores de concurrencia | Resultados distintos con la misma entrada; fallas raras | `go test -race`, semilla fija para lo aleatorio, concurrencia solo en la metaheurística |
| Costo de servidor por cálculos grandes | Peticiones lentas o servidor saturado | Límites por petición, tiempo máximo con `context`, cola de trabajos si hace falta |
| Dependencia de la conexión | Quejas de patios sin señal | El móvil descarga planes ya calculados; Wails calcula local en escritorio |
| Alcance que crece | Ideas nuevas antes de cerrar la fase | Lista para la siguiente fase; la actual no cambia |
| Pérdida de datos | Respaldo nunca restaurado | Respaldos diarios y restauración de prueba mensual |

## Después del plan: fase 4

- **Integraciones:** Carta Porte, API pública para ERPs, VGM, etiquetas QR por pieza.
- **Motor:** centro de gravedad, piezas irregulares, planeación de flota, cálculo por lotes.
- **Opcional:** compilar el motor con TinyGo a WebAssembly para cálculo en el navegador, si los clientes lo piden.
- **Mercado:** otros puertos (Lázaro Cárdenas, Veracruz, Altamira) y alianzas con agencias aduanales.
- **Equipo:** primero soporte y capacitación; después un desarrollador.

## Checklist para la primera sesión

- [ ] Repo `cargo-3d` en GitHub con `main` y `develop`, y el CI en verde.
- [ ] Las verificaciones del entorno pasan.
- [ ] Ejercicio del capítulo 01 de la guía de Go hecho.
- [ ] Leíste `CLAUDE.md` y los ADR 0002, 0006, 0010, 0012 y 0013.
- [ ] Al menos 3 cargas reales de Transportes Villa, el catálogo de equipos y las reglas del patio.

Primer mensaje sugerido para Claude Code (con `/model opus`):

```
Lee CLAUDE.md, docs/src/plan.md y la bitácora más reciente. Vamos a empezar
la fase 1, sprint 1.1: modelos y contrato del motor en Go. Te paso las cargas
reales en engine/fixtures/real/. Antes de programar, propón un plan con el
esquema JSON, los archivos, las pruebas y los capítulos de la guía.
```
