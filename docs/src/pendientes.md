# Pendientes

Registro vivo de **todo lo que quedó abierto**, con lo necesario para retomarlo sin depender de la memoria de una
conversación. Lo actualiza `/cerrar-sesion` al final de cada sesión y lo lee `/iniciar-sesion` al empezar.

Reglas:

- Cada pendiente tiene un id `P-NNN` que **no se reutiliza**. Nunca se borra: al resolverse pasa a «Cerrados» con
  la fecha y cómo se resolvió.
- Lo que no se puede hacer antes de cierta sesión pasa a «Agendados», **solo si** ya está escrito en el prompt de esa
  sesión. Agendado no es hecho: pasa a «Cerrados» cuando la sesión lo termina.
- La bitácora de cada sesión solo cita los ids; el detalle vive aquí.
- Nada de datos de clientes (nombres, direcciones, placas, pedidos): si hace falta, se apunta al archivo de
  `docs/datos/` (que no se sube a git).
- **Quién:** «Benjamín» (algo que solo tú puedes conseguir o decidir), «Claude» (trabajo que se pide en una sesión)
  o «Decisión» (hay que elegir antes de seguir).
- **Cuándo:** la sesión de `prompts.md` antes de la cual tiene que estar resuelto, o «cuando se pueda».

## Abiertos

| Id | Quién | Cuándo | Qué falta | Por qué importa | Cómo retomarlo | Origen |
|---|---|---|---|---|---|---|
| P-007 | Benjamín | antes de 1.1.2 | Punto de control de la 1.1.1: explicar con tus palabras cada campo del contrato | Regla 9 de `CLAUDE.md`: no se avanza de módulo sin explicarlo | `/punto-de-control contrato del motor` | 1.1.1 |
| P-001 | Benjamín | antes de 1.1.5 | Al menos 3 cargas reales, el catálogo de equipos que se usan y las reglas del patio | Sin cargas reales el contrato y los fixtures son sintéticos, y la puerta de la fase 1 pide «una carga real acomodada» | Llena los formatos de `docs/plantillas/patio/` (copia en `docs/datos/patio/`) o pon tus archivos en `docs/datos/`, y corre `/sesion 1.1.5` | 2026-11-20 · 1.1.1 |
| P-003 | Benjamín | antes de un caso `road/rigid/` | Medidas de rabón y tortón del patio: placa o tarjeta de circulación, interior de la caja (largo, ancho, alto, puerta) de 1–2 unidades y número de llantas del tortón | Ningún fabricante publica el interior; sin eso no hay ejemplo ni entrada de catálogo | Llena la hoja de medición §3 en `docs/datos/patio/` y pide «agrega los casos de rabón y tortón según `contracts/examples/README.md` §3» | 1.1.1 |
| P-004 | Benjamín | 1.1.5 | Confirmar valores supuestos: sobrealtura OOG del flat rack (400 mm), altura útil de los reefers (20': 2 158; 40': 2 425), cargas útiles de cajas secas (30 000) y refrigeradas (24 900), altura de cama de plataformas (1 478), altura del chasis (1 200), ancho útil de la Conestoga (2 438) y puerta del 45' pallet wide. El margen lateral ya **no** es supuesto: sale de la NOM-012 | Están en los ejemplos marcados como «supuesto»; no pueden pasar al catálogo como definitivos | Llena la hoja de medición §4 en `docs/datos/patio/`; se corrigen `contracts/examples/` y `equipment.json` en la 1.1.5 | 1.1.1 |
| P-013 | Benjamín | cuando se pueda | Fichas o medidas de los equipos que sí se usan y no tienen ficha pública confiable: 40' reefer de altura estándar, plataforma de 40', extensible y camioneta de 3,5 t; y altura del piso de la plataforma ISO 20' | Para pasar de ⏳ a ✅ en la cobertura por equipo | Llena `docs/plantillas/patio/hoja-de-medicion.md` §5 (copia en `docs/datos/patio/`) y pide «agrega los casos según `contracts/examples/README.md` §3» | 1.1.1 |
| P-014 | Benjamín + Claude | cuando se pueda | Publicar la v5 del sistema oficial de Voxa (aprobada el 2026-11-21): tokens `on-danger`, `scrim` y `shadow-float`; diálogo, aviso, barra masiva, botón `danger` y píldora `live` | La publicación se bloqueó por una regla de permisos de lectura de Claude Code; el artefacto sigue en v4 y el repo conserva las marcas «PROPUESTA» | Permite la lectura del archivo cuando Claude Code lo pida y di «publica `docs/datos/voxa-v5/` en el sistema oficial». Antes se relee el artefacto; después se quitan las marcas «PROPUESTA» de `packages/ui`, `docs/src/ux/sistema-diseno.md` y el ADR-0014 | 2026-11-21 · Voxa |

## Agendados

No están hechos. Cada uno ya está escrito en el prompt de su sesión (`prompts.md`), que es donde se hace.

| Id | Cuándo | Qué falta | Dónde quedó escrito | Origen |
|---|---|---|---|---|
| P-002 | 1.1.5 | Convertir el catálogo del cliente con las correcciones ya decididas (tara restada dos veces, filas mal nombradas) y preguntar por la fila 9 | Prompt 1.1.5; tabla fila por fila en `docs/datos/revision-catalogo-contenedores.md` | 1.1.1 |
| P-008 | 1.1.4 | La validación de la entrada debe: rechazar un equipo con `door` y `overhang` mayor que cero; revisar ids repetidos y tramos encimados; no tratar «no cabe» ni «puerta más grande que el interior» como error | Prompt 1.1.4; `contracts/examples/invalid/` ya trae los casos | 1.1.1 |
| P-010 | 1.4.2 | Pasar las reglas físicas a `packages/validator` y borrar el script provisional | Prompt 1.4.2; mientras, `scripts/contratos/check-examples.mjs` (`pnpm examples:check`) | 1.1.1 |
| P-012 | 2.1.6 | Subir `sections.maxItems` y agregar ejemplos de escalonada, doble caída, lowboy y RGN (la nodriza, de varios niveles, se decide ahí mismo) | Prompt 2.1.6 (sesión nueva) y plan, sprint 2.1 | 1.1.1 |
| P-015 | 1.6.1 | Conectar `@cargo-3d/ui` a `apps/web` | Prompt 1.6.1 | 2026-11-21 · Voxa |
| P-016 | 3.2.1 y 3.5.1 | Tokens para móvil (React Native) e íconos de app para Wails y Expo | Prompts 3.2.1 y 3.5.1 | 2026-11-21 · Voxa |

## Cerrados

| Id | Cerrado | Cómo se resolvió |
|---|---|---|
| P-005 | 2026-11-21 | Tablas A-1 y A-2 leídas de las imágenes del DOF: coinciden con la fuente secundaria para C y T-S, y C-R y T-S-R tienen otros valores en ET y A. No se encontró proyecto de modificación publicado; la sesión 2.1.5 lo vuelve a revisar. Todo en `docs/src/conceptos/equipos-y-normas.md` |
| P-006 | 2026-11-21 | Casos `40ot-door`, `40ot-roof`, `20fr-in-gauge` y `48rv-pallets`. De paso se corrigió el margen de `40fr-oog`: 200 mm por lado pasaba de los 2 600 mm de la NOM-012 |
| P-009 | 2026-11-21 | `contracts/examples/invalid/`: 15 casos que rechaza el esquema y 2 que rechaza la validación del motor. «Tramos encimados» se agrega en la 2.1.6 |
| P-011 | 2026-11-21 | `docs/src/conceptos/sistema-coordenadas.md` y `docs/src/conceptos/equipos-y-normas.md` |
| P-017 | 2026-11-21 | Instalación de Node.js 26 LTS de forma manual (`fnm install 26`) |
| P-018 | 2026-11-21 | No había error: `arguments: [numero]` es lo que llena `$numero` en la skill y funciona; el aviso es del editor. El arreglo que se había propuesto la habría roto |
