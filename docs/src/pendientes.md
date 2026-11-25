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
| P-004 | Benjamín | 2.6 (pruebas con cargas reales) | Confirmar en unidades reales los valores sin fuente: línea roja de los reefers (los ejemplos usan alto interior Hapag: 20' 2 159; 40' HC 2 425), carga útil **legal** de cajas secas y refrigeradas, altura de cama de plataformas (1 478), piso del flat rack 40' sobre el chasis, alto libre de la plataforma ISO 40' (2 383, derivado), ancho útil de la Conestoga (2 438) y tara del tracto (8 800). Lo confirmado por fuente ya se aplicó el 2026-11-23 | No pueden pasar al catálogo como definitivos | Mientras, los casos usan estos valores marcados como supuestos o simulados (ADR-0020). En el periodo de pruebas se miden con la hoja §4 de `docs/plantillas/patio/` y el registro de `docs/src/casos/pruebas-reales.md`. Antes de cerrarlo, cerrar P-020 | 1.1.1 |
| P-013 | Benjamín | 2.6 (pruebas con cargas reales) | Medidas de equipos sin ficha completa: 40' reefer estándar (línea roja y puerta), plataforma de 40' (altura de cama), caja refrigerada 48' (alto interior), extensible (útil por largo), camionetas, rabón y tortón (hoy simulados en CS-016 a CS-020) y piso de la plataforma ISO 20' | Los casos con equipos simulados pasan a «Validado en real» solo con medidas reales | En el periodo de pruebas, alta de cada equipo de la empresa con origen `medido` (ADR-0019) y hoja §5 de `docs/plantillas/patio/` | 1.1.1 |
| P-020 | Claude | antes de cerrar P-004 | Verificar con ficha de fabricante los 108 valores «no verificados» de la auditoría del 2026-11-23 (interiores, puertas y taras de DV, HC, OT y FR) y decidir, ejemplo por ejemplo, si se alinea con `contracts/equipment.json`. Diferencias anotadas sin cambiar: 20dv, 20hc, 40dv, 40hc, 45hc, 45pw (interior), 20fr, 20ot, 40ot, 40pf (tara y útil), 53dv (ancho 2 565 contra 2 540 Fruehauf), 53fb (largo), 48rv y 53rv (interior C.H. Robinson contra Alzper y Hyundai) | Un ejemplo con otra serie no es un error, pero mezclar series en un mismo cálculo sí | Comparar el `container` de cada `input.json` con su equipo de `contracts/equipment.json` (interior, puerta, útil y tara) y pedir «verifica P-020 contra fichas de fabricante, sin cambiar valores sin fuente» | 2026-11-23 · catálogo |
| P-021 | Claude | cuando se pueda | Fotos del catálogo: la de la fila 20 (Bulk 20') devuelve 404; las de las filas 1 (403 a clientes automáticos) y 3 (sin respuesta) hay que revisarlas en un navegador | El catálogo se publica en el sitio con esas fotos | Buscar una foto real con la boca del techo visible para el Bulk 20' y abrir las filas 1 y 3 a mano; cambiar una foto solo avisando | 2026-11-23 · catálogo |
| P-019 | Benjamín | antes de cambiar de computadora | Copiar `docs/datos/` y `.env` a la computadora nueva (no están en git) | Ahí viven las investigaciones privadas y la configuración local | Copiar ambas rutas a mano o a un respaldo cifrado; si ya se hizo, pasar a «Cerrados» | 2026-11-21 (estaba en la bitácora y en CLAUDE.md, faltaba en este registro) |
| P-022 | Benjamín | 2.6 | Conseguir al menos una empresa para el periodo de pruebas con cargas reales de la v0.2.0 y su permiso para usar datos anonimizados | Es la única forma de validar los casos simulados y el criterio de la fase 2 | Seguir `docs/src/casos/pruebas-reales.md`: ficha de la empresa, alta de equipos, reglas del patio y `registro.csv` en `docs/datos/pruebas/<empresa>/`; sesiones 2.6.3 y 2.6.4 | 2026-11-25 · casos |
| P-023 | Claude | 2.1.5 | Carga sobredimensionada con permiso (NOM-040): revisar en el DOF si la NOM-040-SCT-2-2012 sigue vigente (cita la NOM-012 de 2008), decidir cómo llega el permiso a la API (tipo de combinación, ancho y alto autorizados) y agregar al contrato un campo por pieza (por ejemplo `indivisible`) para que solo la carga indivisible use el margen del permiso | Hoy el margen es por tramo: una pieza divisible podría quedar fuera del ancho legal (CS-028) | Leer `conceptos/equipos-y-normas.md` §«Carga sobredimensionada en carretera» y el caso CS-028; el cambio de contrato con `/cambiar-contrato` | 2026-11-25 · casos |

## Agendados

No están hechos. Cada uno ya está escrito en el prompt de su sesión (`prompts.md`), que es donde se hace.

| Id | Cuándo | Qué falta | Dónde quedó escrito | Origen |
|---|---|---|---|---|
| P-008 | 1.1.4 | La validación de la entrada debe: rechazar un equipo con `door` y `overhang` mayor que cero; revisar ids repetidos y tramos encimados; no tratar «no cabe» ni «puerta más grande que el interior» como error | Prompt 1.1.4; `contracts/examples/invalid/` ya trae los casos | 1.1.1 |
| P-010 | 1.4.2 | Pasar las reglas físicas a `packages/validator` y borrar el script provisional; incluir la validación contra el esquema (el script actual no la hace: el 2026-11-25 se encontraron 4 descripciones de más de 200 caracteres) | Prompt 1.4.2; mientras, `scripts/contratos/check-examples.mjs` (`pnpm examples:check`) | 1.1.1 |
| P-012 | 2.1.6 | Subir `sections.maxItems` y agregar ejemplos de escalonada, doble caída, lowboy y RGN (la nodriza, de varios niveles, se decide ahí mismo) | Prompt 2.1.6 (sesión nueva) y plan, sprint 2.1 | 1.1.1 |
| P-015 | 1.6.1 | Conectar `@cargo-3d/ui` a `apps/web` | Prompt 1.6.1 | 2026-11-21 · Voxa |
| P-016 | 3.2.1 y 3.5.1 | Tokens para móvil (React Native) e íconos de app para Wails y Expo | Prompts 3.2.1 y 3.5.1 | 2026-11-21 · Voxa |

## Cerrados

| Id | Cerrado | Cómo se resolvió |
|---|---|---|
| P-001 | 2026-11-25 | Sustituido por los casos de uso simulados de `docs/src/casos/` (ADR-0020): no hay patio, cargas ni equipos a la mano. La validación con cargas reales pasó a P-022 (v0.2.0) |
| P-003 | 2026-11-25 | Sustituido por los casos simulados de rabón y tortón (CS-019 y CS-020) y la plantilla de equipos de la empresa (ADR-0019). Las medidas reales se capturan en el periodo de pruebas (P-013) |
| P-007 | 2026-11-25 | Cancelado: se quitaron los puntos de control (ADR-0018). Su tema técnico (quién calcula `maxPayloadKg` y hacia dónde se redondea) se explica sobre el código en la sesión 1.1.2 |
| P-002 | 2026-11-23 | Catálogo del cliente convertido con la columna «bruto» como carga útil y los renombres decididos; tabla privada en `docs/datos/sprint 1.1.1/catalogo-cliente-convertido.md`. Fila 9 (Flatrack Collapsible 20'): Benjamín dio la ficha de contenedoresmaritimos.eu (útil 27 110, tara 2 890, bruto 30 000), verificada el mismo día. No cambia `contracts/equipment.json`, que sigue con fuentes públicas |
| P-005 | 2026-11-21 | Tablas A-1 y A-2 leídas de las imágenes del DOF: coinciden con la fuente secundaria para C y T-S, y C-R y T-S-R tienen otros valores en ET y A. No se encontró proyecto de modificación publicado; la sesión 2.1.5 lo vuelve a revisar. Todo en `docs/src/conceptos/equipos-y-normas.md` |
| P-006 | 2026-11-21 | Casos `40ot-door`, `40ot-roof`, `20fr-in-gauge` y `48rv-pallets`. De paso se corrigió el margen de `40fr-oog`: 200 mm por lado pasaba de los 2 600 mm de la NOM-012 |
| P-009 | 2026-11-21 | `contracts/examples/invalid/`: 15 casos que rechaza el esquema y 2 que rechaza la validación del motor. «Tramos encimados» se agrega en la 2.1.6 |
| P-011 | 2026-11-21 | `docs/src/conceptos/sistema-coordenadas.md` y `docs/src/conceptos/equipos-y-normas.md` |
| P-014 | 2026-11-21 | Publicada la v5 del sistema oficial de Voxa (tokens `on-danger`, `scrim` y `shadow-float`; diálogo, aviso, barra masiva, botón `danger` y píldora `live`) y quitadas las marcas «PROPUESTA» del repo |
| P-017 | 2026-11-21 | Instalación de Node.js 26 LTS de forma manual (`fnm install 26`) |
| P-018 | 2026-11-21 | No había error: `arguments: [numero]` es lo que llena `$numero` en la skill y funciona; el aviso es del editor. El arreglo que se había propuesto la habría roto |
