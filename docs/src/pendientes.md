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
| P-001 | Benjamín | antes de 1.1.5 | Al menos 3 cargas reales, el catálogo de equipos que se usan y las reglas del patio | Sin cargas reales el contrato y los fixtures son sintéticos, y la puerta de la fase 1 pide «una carga real acomodada» | Llena los formatos de `docs/plantillas/patio/` (copia en `docs/datos/patio/`) o pon tus archivos en `docs/datos/`, y corre `/sesion 1.1.5` | 2026-11-20 · 1.1.1 |
| P-003 | Benjamín | antes de un caso `road/rigid/` | Medidas de rabón y tortón del patio: placa o tarjeta de circulación, interior de la caja (largo, ancho, alto, puerta) de 1–2 unidades y número de llantas del tortón | Ningún fabricante publica el interior; sin eso no hay ejemplo ni entrada de catálogo | Llena la hoja de medición §3 en `docs/datos/patio/` y pide «agrega los casos de rabón y tortón según `contracts/examples/README.md` §3» | 1.1.1 |
| P-004 | Benjamín | 1.1.5 | Confirmar en la unidad los valores que no tienen fuente: línea roja de los reefers (los ejemplos usan alto interior Hapag: 20' 2 159; 40' HC 2 425), carga útil **legal** de cajas secas (30 000 es capacidad estructural Fruehauf) y refrigeradas (25 000 / 30 000 es el máximo Alzper), altura de cama de plataformas (1 478), piso del flat rack 40' sobre el chasis (sin cifra pública), alto libre de la plataforma ISO 40' (2 383, derivado), ancho útil de la Conestoga (2 438) y tara del tracto (8 800). Ya se aplicó lo confirmado por fuente el 2026-11-23: chasis 1 219, 648 como altura plegada, 400 como límite interno, puerta del 45' pallet wide seco 2 416 × 2 585 y reefers con la serie Hapag | No pueden pasar al catálogo como definitivos | Llena la hoja de medición §4 en `docs/datos/sprint 1.1.1/patio/`. Antes de cerrarlo, cerrar también P-020 | 1.1.1 |
| P-013 | Benjamín | cuando se pueda | Medidas de los equipos que sí se usan y no tienen ficha completa: 40' reefer de altura estándar (CIMC da interior, tara y útil; faltan línea roja y puerta), plataforma de 40' (Fruehauf PM-HT2-40W da todo menos la altura de cama), caja refrigerada 48' (alto interior: Alzper publica 3 790, imposible), extensible (útil por largo) y camioneta de 3,5 t; y altura del piso de la plataforma ISO 20' | Para pasar de ⏳ a ✅ en la cobertura por equipo | Llena la hoja de medición §5 (copia en `docs/datos/sprint 1.1.1/patio/`) y pide «agrega los casos según `contracts/examples/README.md` §3» | 1.1.1 |
| P-020 | Claude | antes de cerrar P-004 | Verificar con ficha de fabricante los 108 valores «no verificados» de la auditoría del 2026-11-23 (interiores, puertas y taras de DV, HC, OT y FR) y decidir, ejemplo por ejemplo, si se alinea con `contracts/equipment.json`. Diferencias anotadas sin cambiar: 20dv, 20hc, 40dv, 40hc, 45hc, 45pw (interior), 20fr, 20ot, 40ot, 40pf (tara y útil), 53dv (ancho 2 565 contra 2 540 Fruehauf), 53fb (largo), 48rv y 53rv (interior C.H. Robinson contra Alzper y Hyundai) | Un ejemplo con otra serie no es un error, pero mezclar series en un mismo cálculo sí | Comparar el `container` de cada `input.json` con su equipo de `contracts/equipment.json` (interior, puerta, útil y tara) y pedir «verifica P-020 contra fichas de fabricante, sin cambiar valores sin fuente» | 2026-11-23 · catálogo |
| P-021 | Claude | cuando se pueda | Fotos del catálogo: la de la fila 20 (Bulk 20') devuelve 404; las de las filas 1 (403 a clientes automáticos) y 3 (sin respuesta) hay que revisarlas en un navegador | El catálogo se publica en el sitio con esas fotos | Buscar una foto real con la boca del techo visible para el Bulk 20' y abrir las filas 1 y 3 a mano; cambiar una foto solo avisando | 2026-11-23 · catálogo |

## Agendados

No están hechos. Cada uno ya está escrito en el prompt de su sesión (`prompts.md`), que es donde se hace.

| Id | Cuándo | Qué falta | Dónde quedó escrito | Origen |
|---|---|---|---|---|
| P-007 | 1.1.6 | Punto de control del contrato, **no aprobado** (dos intentos el 2026-11-21, 3 de 5 cada uno): falta quién calcula `maxPayloadKg` y hacia dónde se redondea | Prompts 1.1.2 (se explica sobre el código) y 1.1.6 (se vuelve a preguntar) | 1.1.1 |
| P-008 | 1.1.4 | La validación de la entrada debe: rechazar un equipo con `door` y `overhang` mayor que cero; revisar ids repetidos y tramos encimados; no tratar «no cabe» ni «puerta más grande que el interior» como error | Prompt 1.1.4; `contracts/examples/invalid/` ya trae los casos | 1.1.1 |
| P-010 | 1.4.2 | Pasar las reglas físicas a `packages/validator` y borrar el script provisional | Prompt 1.4.2; mientras, `scripts/contratos/check-examples.mjs` (`pnpm examples:check`) | 1.1.1 |
| P-012 | 2.1.6 | Subir `sections.maxItems` y agregar ejemplos de escalonada, doble caída, lowboy y RGN (la nodriza, de varios niveles, se decide ahí mismo) | Prompt 2.1.6 (sesión nueva) y plan, sprint 2.1 | 1.1.1 |
| P-015 | 1.6.1 | Conectar `@cargo-3d/ui` a `apps/web` | Prompt 1.6.1 | 2026-11-21 · Voxa |
| P-016 | 3.2.1 y 3.5.1 | Tokens para móvil (React Native) e íconos de app para Wails y Expo | Prompts 3.2.1 y 3.5.1 | 2026-11-21 · Voxa |

## Cerrados

| Id | Cerrado | Cómo se resolvió |
|---|---|---|
| P-002 | 2026-11-23 | Catálogo del cliente convertido con la columna «bruto» como carga útil y los renombres decididos; tabla privada en `docs/datos/sprint 1.1.1/catalogo-cliente-convertido.md`. Fila 9 (Flatrack Collapsible 20'): Benjamín dio la ficha de contenedoresmaritimos.eu (útil 27 110, tara 2 890, bruto 30 000), verificada el mismo día. No cambia `contracts/equipment.json`, que sigue con fuentes públicas |
| P-005 | 2026-11-21 | Tablas A-1 y A-2 leídas de las imágenes del DOF: coinciden con la fuente secundaria para C y T-S, y C-R y T-S-R tienen otros valores en ET y A. No se encontró proyecto de modificación publicado; la sesión 2.1.5 lo vuelve a revisar. Todo en `docs/src/conceptos/equipos-y-normas.md` |
| P-006 | 2026-11-21 | Casos `40ot-door`, `40ot-roof`, `20fr-in-gauge` y `48rv-pallets`. De paso se corrigió el margen de `40fr-oog`: 200 mm por lado pasaba de los 2 600 mm de la NOM-012 |
| P-009 | 2026-11-21 | `contracts/examples/invalid/`: 15 casos que rechaza el esquema y 2 que rechaza la validación del motor. «Tramos encimados» se agrega en la 2.1.6 |
| P-011 | 2026-11-21 | `docs/src/conceptos/sistema-coordenadas.md` y `docs/src/conceptos/equipos-y-normas.md` |
| P-014 | 2026-11-21 | Publicada la v5 del sistema oficial de Voxa (tokens `on-danger`, `scrim` y `shadow-float`; diálogo, aviso, barra masiva, botón `danger` y píldora `live`) y quitadas las marcas «PROPUESTA» del repo |
| P-017 | 2026-11-21 | Instalación de Node.js 26 LTS de forma manual (`fnm install 26`) |
| P-018 | 2026-11-21 | No había error: `arguments: [numero]` es lo que llena `$numero` en la skill y funciona; el aviso es del editor. El arreglo que se había propuesto la habría roto |
