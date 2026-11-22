# Pendientes

Registro vivo de **todo lo que quedó abierto**, con lo necesario para retomarlo sin depender de la memoria de una
conversación. Lo actualiza `/cerrar-sesion` al final de cada sesión y lo lee `/iniciar-sesion` al empezar.

Reglas:

- Cada pendiente tiene un id `P-NNN` que **no se reutiliza**. Nunca se borra: al resolverse pasa a «Cerrados» con
  la fecha y cómo se resolvió.
- La bitácora de cada sesión solo cita los ids; el detalle vive aquí.
- Nada de datos de clientes (nombres, direcciones, placas, pedidos): si hace falta, se apunta al archivo de
  `docs/datos/` (que no se sube a git).
- **Quién:** «Benjamín» (algo que solo tú puedes conseguir o decidir), «Claude» (trabajo que se pide en una sesión)
  o «Decisión» (hay que elegir antes de seguir).
- **Cuándo:** la sesión de `prompts.md` antes de la cual tiene que estar resuelto, o «cuando se pueda».

## Abiertos

| Id | Quién | Cuándo | Qué falta | Por qué importa | Cómo retomarlo | Origen |
|---|---|---|---|---|---|---|
| P-001 | Benjamín | antes de 1.1.5 | Al menos 3 cargas reales, el catálogo de equipos que se usan y las reglas del patio | Sin cargas reales el contrato y los fixtures son sintéticos, y la puerta de la fase 1 pide «una carga real acomodada» | Pon los archivos (Excel, CSV o fotos) en `docs/datos/` y corre `/sesion 1.1.5` | 2026-11-20 · 1.1.1 |
| P-002 | Benjamín + Claude | 1.1.5 | Revisar el catálogo del cliente (`docs/datos/contenedores_dimensiones.sql`): tara en toneladas con etiqueta «kg», filas «High Cube 20'» y «High Cube 40'» con medidas de 40' y 45', 20' RF de 5 724 mm de largo, alturas de open top menores que las fichas | Si son medidas reales de sus unidades, mandan sobre las fichas públicas; si son errores, contaminan `equipment.json` | En la 1.1.5, decir fila por fila si es medida real o error. Lista completa en `docs/datos/investigacion-1.1.1-contrato-motor.md` §1 y `investigacion-equipos-faltantes.md` §2 | 1.1.1 |
| P-003 | Benjamín | antes de un caso `road/rigid/` | Medidas de rabón y tortón del patio: placa o tarjeta de circulación, interior de la caja (largo, ancho, alto, puerta) de 1–2 unidades y número de llantas del tortón | Ningún fabricante publica el interior; sin eso no hay ejemplo ni entrada de catálogo | Guárdalas en `docs/datos/` y pide «agrega los casos de rabón y tortón según `contracts/examples/README.md` §3» | 1.1.1 |
| P-004 | Benjamín | 1.1.5 | Confirmar valores supuestos: margen OOG del flat rack (hoy 600 arriba y 200 por lado), altura útil del 20' RF (2 158, de un fragmento), cargas útiles de caja seca 48' y 53' (30 000) y caja refrigerada 53' (24 900), altura de cama de plataformas (1 478) | Hoy están en los ejemplos marcados como «supuesto»; no pueden pasar al catálogo como definitivos | Con la naviera o las placas de las unidades; corrige `contracts/examples/` y `equipment.json` en la 1.1.5 | 1.1.1 |
| P-005 | Claude | antes de 2.1 | Comparar la Tabla A-2 (peso por eje) de la NOM-012-SCT-2-2017 contra el DOF (hoy viene de una fuente secundaria) y buscar si hay PROY-NOM-012 publicado en 2026 | El reparto de peso por eje del sprint 2.1 se programa con esa tabla | `/buscador NOM-012 tabla A-2 y reformas 2026` antes de `/sesion 2.1.x` | 1.1.1 |
| P-006 | Claude | cuando se pueda | Escribir los casos con datos ya listos: open top 40' (puerta y techo), flat rack 20', caja refrigerada 48' | Completa la cobertura por equipo | Pide «agrega los casos ⏳ con datos listos de `contracts/examples/README.md`»; los datos están en `docs/datos/investigacion-equipos-faltantes.md` §5 | 1.1.1 |
| P-007 | Benjamín | antes de 1.1.2 | Punto de control de la 1.1.1: explicar con tus palabras cada campo del contrato | Regla 9 de `CLAUDE.md`: no se avanza de módulo sin explicarlo | `/punto-de-control contrato del motor` | 1.1.1 |
| P-008 | Claude | 1.1.4 | La validación de la entrada debe: rechazar un equipo con `door` y `overhang` mayor que cero; revisar ids repetidos y tramos encimados; no tratar «no cabe» ni «puerta más grande que el interior» como error | Decidido en ADR-0016 y ADR-0017; el prompt de la 1.1.4 ya lo pide | `/sesion 1.1.4` y comprobar contra este pendiente | 1.1.1 |
| P-009 | Claude | 1.1.4 | Carpeta `contracts/examples/invalid/` con entradas mal formadas a propósito (una por cada error de validación) | Documenta y prueba cada error | Dentro de `/sesion 1.1.4` | 1.1.1 |
| P-010 | Claude | 1.4.2 | El validador independiente debe revisar lo que el esquema no puede: pieza dentro de su tramo y margen (incluido `zMm` negativo contra `overhang.leftMm`), `sectionId` existente, `outOfGauge` correcto, traslapes, apoyo, peso, motivos y métricas | Durante la 1.1.1 esa revisión la hizo un script temporal que **no** está en el repo | `/sesion 1.4.2`; la lista de reglas está en `contracts/examples/README.md` §3, paso 4 | 1.1.1 |
| P-011 | Claude | 1.1.3 y 1.1.5 | Documentos de conceptos: `docs/src/conceptos/sistema-coordenadas.md` (con dibujo) y `docs/src/conceptos/equipos-y-normas.md` (ISO 668/1496/6346, NOM-012, VGM, línea de carga del reefer, OOG) | Regla 5 de `CLAUDE.md` | El primero en `/sesion 1.1.3`; el segundo en `/sesion 1.1.5` con `docs/datos/investigacion-*.md` | 1.1.1 |
| P-012 | Decisión | cuando el motor acomode en varios tramos | Subir `sections.maxItems` y agregar ejemplos de escalonada, doble caída, lowboy, RGN y nodriza | Hoy están bloqueados (🔒) en `contracts/examples/README.md` | `/cambiar-contrato varios tramos` | 1.1.1 |
| P-013 | Benjamín | cuando se pueda | Fichas de equipos sin datos: 45' pallet wide (el cliente tiene la fila), 20' HC, 40' RF estándar, plataforma ISO, plataforma 40', extensible, Conestoga, camioneta 3,5 t | Para pasar de ⏳ a ✅ en la cobertura por equipo | Confirma cuáles se usan de verdad; los demás se quedan ⏳. Después `/buscador <equipos>` | 1.1.1 |
| P-014 | Claude | cuando se pueda | Subir al sistema oficial de claude.ai los componentes y tokens «propuesta» de Voxa (modal, toasts, barra masiva) | Regla de diseño: un valor se cambia en `tokens.css` y en el sistema oficial a la vez | Pedirlo en una sesión fuera del plan | 2026-11-21 · Voxa |
| P-015 | Claude | 1.6.1 | Conectar `@cargo-3d/ui` a `apps/web` | La web usa el sistema de diseño desde su primera pantalla | `/sesion 1.6.1` | 2026-11-21 · Voxa |
| P-016 | Claude | 3.2 y 3.5 | Tokens para móvil (React Native) e íconos de app para Wails y Expo | Las apps de escritorio y patio comparten la marca | `/sesion 3.2.x` y `/sesion 3.5.x` | 2026-11-21 · Voxa |
| P-017 | Benjamín | antes de 1.4.1 | Actualizar Node.js de 22 a 26 LTS (`fnm install 26`) | El proyecto pide Node 26; pnpm solo lo advierte, pero el validador en TypeScript arranca en la 1.4 | `fnm install 26 && fnm use 26`, luego `node -v` | 2026-11-21 · Voxa |
| P-018 | Claude | cuando se pueda | `.claude/skills/sesion/SKILL.md` línea 5 usa el atributo `arguments`, que el editor marca como no válido (los permitidos son `argument-hint`, `description`, `name`, `disable-model-invocation`…) | Puede que el argumento de `/sesion <número>` no se pase como se espera en futuras versiones | Cambiar `arguments` por `argument-hint` y probar `/sesion 1.1.2` | 2026-11-21 · 1.1.1 |

## Cerrados

| Id | Cerrado | Cómo se resolvió |
|---|---|---|
| — | — | — |
