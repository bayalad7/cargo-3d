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

**Si usas Windows:** instala **WSL2 con Ubuntu** y haz todo dentro de WSL.

### 1. Git y GitHub

1. Crea en GitHub el repositorio vacío `cargo-3d` (privado).
2. Entra a la carpeta `cargo-3d` y súbelo: `git remote add origin https://github.com/bayalad7/cargo-3d.git` y `git push -u origin main`.
3. Protege `main` en GitHub (Settings → Branches): nada entra sin pull request.
4. Si usas otro usuario u otro nombre de repo, cambia la ruta del módulo en `backend/go.mod` y en los imports (búscala con `grep -r bayalad7`).

### 2. Go

1. Instala Go 1.27 o superior desde go.dev/dl.
2. Verifica: `go version`.

### 3. Docker

Docker Desktop (o Docker Engine en Linux) con Docker Compose. Copia `.env.example` a `.env` y cambia la contraseña.

### 4. Node.js y pnpm (desde el sprint 1.4: validador en TypeScript)

Node 26 LTS con fnm o nvm. Desde Node 25, Corepack ya no viene incluido: instala pnpm con
`npm install -g pnpm@12`. pnpm lee el campo `packageManager` de `package.json` y usa
exactamente la versión del proyecto.

### 5. mdBook (sitio de documentación)

Descarga el binario de mdBook desde su página de releases, o instálalo con `cargo install mdbook` si tienes Rust.

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
| `docker compose up -d db` y `docker compose ps` | El servicio `db` como `healthy` |
| `go run ./backend/cmd/api` y abrir `localhost:8080/health` | `{"status":"ok","engineVersion":"0.1.0"}` |
| `docker compose --profile api up --build` | La API respondiendo en `localhost:8080/health` desde Docker |
| `mdbook serve docs --open` | El sitio de documentación |

## Cómo trabajar con Claude Code

Tú diriges y revisas; Claude Code escribe, documenta y prueba. Cada sesión tiene una meta pequeña,
vive en su propia rama (creada desde `main`) y termina con la bitácora actualizada.

### Ciclo de una sesión

1. **Arranque.** `git checkout -b motor/modelos-basicos` y `claude`. Pide: *"Lee CLAUDE.md y la última bitácora y dime dónde vamos."*
2. **Plan.** Modo de planeación (`Shift+Tab`) para todo lo que no sea trivial.
3. **Contrato primero.** Si cambia la API o el motor, se edita `contracts/` antes que el código.
4. **Implementación en pasos chicos**, cada uno con pruebas de tabla y documentación.
5. **Verificación:** `gofmt -l backend` vacío, `go vet ./backend/...` y `go test ./backend/...` sin errores.
6. **Revisión tuya:** lee el diff completo; pregunta lo que no entiendas: *"Explícame esta línea como si viniera de PHP."*
7. **Punto de control:** si se terminó un módulo, explícalo con tus palabras.
8. **Cierre:** *"Actualiza la bitácora, el CHANGELOG y la sección Estado de CLAUDE.md."* Commit, push y pull request.

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
   │  Backend Go con /health, Docker, contratos, 11 ADR y guía de Go.
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

- [ ] GitHub Actions: `gofmt`, `go vet`, `go test`, build de la imagen Docker y build de la web.
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

- [ ] Repo `cargo-3d` en GitHub y `main` protegida.
- [ ] Las verificaciones del entorno pasan.
- [ ] Ejercicio del capítulo 01 de la guía de Go hecho.
- [ ] Leíste `CLAUDE.md` y los ADR 0002, 0006 y 0010.
- [ ] Al menos 3 cargas reales de Transportes Villa, el catálogo de equipos y las reglas del patio.

Primer mensaje sugerido para Claude Code (con `/model opus`):

```
Lee CLAUDE.md, docs/src/plan.md y la bitácora más reciente. Vamos a empezar
la fase 1, sprint 1.1: modelos y contrato del motor en Go. Te paso las cargas
reales en engine/fixtures/real/. Antes de programar, propón un plan con el
esquema JSON, los archivos, las pruebas y los capítulos de la guía.
```
