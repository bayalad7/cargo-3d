# Contratos

Archivos que definen cómo se comunican las partes del sistema, **independientes del lenguaje**.
Si el motor o la API se reescriben en otro lenguaje, estos contratos no cambian.

| Archivo | Qué define | Quién lo usa |
|---|---|---|
| `openapi.yaml` | Rutas, peticiones y respuestas de la API | Backend Go (servidor), `packages/api-client` (generado), web y móvil |
| `engine.schema.json` | Entrada y salida del motor (`optimize(EngineInput) -> LoadPlan`) | Motor, validador, casos de prueba de `engine/fixtures/` |
| `examples/` | Ejemplos de entrada y salida que validan contra `engine.schema.json` | Pruebas de ida y vuelta JSON ↔ Go (sesión 1.1.2) y documentación |

Regla: **el contrato se cambia primero**, en su propio commit, y después el código.
Ver `docs/src/adr/0008-api-first-openapi.md`.

## `engine.schema.json`

JSON Schema 2020-12. Las dos raíces están en `$defs`: `EngineInput` (lo que recibe el motor) y
`LoadPlan` (lo que devuelve). Cada campo trae su `description` en español; aquí solo van las convenciones.

### Convenciones

- **Nombres en camelCase**, igual que `openapi.yaml` (`lengthMm`, `maxPayloadKg`).
- **Unidades en el nombre:** todo lo que termina en `Mm` son milímetros y en `Kg` kilogramos, **siempre enteros**.
  Al capturar, las medidas y pesos de las piezas se redondean hacia arriba y las capacidades del equipo hacia abajo
  (ADR-0016).
  Regla corta: **redondea en tu contra**. El plan nunca debe prometer algo que no cabe.

  | Valor real | Qué es | Se guarda | En Go |
  |---|---|---|---|
  | puerta de 98,3" = 2 496,82 mm | equipo | 2 496 | `math.Floor` |
  | caja de 47,2" = 1 198,88 mm | pieza | 1 199 | `math.Ceil` |
  | carga útil de 26 580,7 kg | equipo | 26 580 | `math.Floor` |
  | pieza de 12,3 kg | pieza | 13 | `math.Ceil` |

- **Quién sabe qué** (ADR-0016):

  | Quién | Qué sabe |
  |---|---|
  | API o catálogo | NOM-012, tipo de camino, configuración del vehículo. **Calcula** `maxPayloadKg` |
  | Contrato | Que `maxPayloadKg` es un entero. Su `maximum` solo atrapa errores de captura; no es un límite legal |
  | Motor | Ese número, y que no debe pasarse |

  Si cambia la norma, se corrige la API; el contrato y el motor se quedan igual.
- **Coordenadas** (ADR-0015): origen en la esquina interior del fondo, abajo, a la izquierda **mirando del fondo
  hacia la puerta**; `x` = largo hacia la puerta, `y` = alto hacia arriba, `z` = ancho hacia la derecha.
  Es de mano derecha, igual que Three.js.

  ```
          y (arriba)
          │
   fondo  │
   (nariz)o──────────► x  (hacia la puerta)
         ╱
        z  (ancho)
  ```

- **El equipo es una lista de tramos** (`sections`), cada uno una caja. En la versión 0.1.0 solo se acepta uno.
- **Orientaciones:** 6 permutaciones (`LWH`, `WLH`, `LHW`, `HLW`, `WHL`, `HWL`); la 1.ª letra va sobre `x`,
  la 2.ª sobre `z` y la 3.ª sobre `y`.
- **Lo que no cupo no es un error:** sale en `unplaced` con su motivo (ADR-0016). Cada unidad lleva **un solo** motivo,
  el primero que se cumpla en el orden de `UnplacedReason`. Los errores de validación son
  solo para entradas mal formadas.
- **Carga sobredimensionada (OOG)** (ADR-0017): un tramo abierto (flat rack, plataforma) puede traer
  `overhang: { topMm, leftMm, rightMm }`, el margen que una pieza puede rebasar por arriba y a cada lado (nunca al
  frente ni atrás). La pieza que lo usa lleva `outOfGauge: true` y se cuenta en `metrics.outOfGaugeCount`. Con margen
  a la izquierda, `zMm` puede ser negativo, y `volumeUtilization` puede pasar de 1. Solo para tramos abiertos: un
  equipo con `door` no lleva margen. Los límites legales los fija quien llama al motor.
- **Límites de tamaño:** el esquema solo pone rangos físicos. Los límites de cuántas piezas y cuánto tiempo de
  cálculo acepta una petición los pone la API en el sprint 1.5; el esquema no es la única barrera.
- **Lo que el esquema no puede revisar:** en la entrada, ids repetidos y tramos encimados (validación del motor,
  sesión 1.1.4); en la salida, que cada pieza quede dentro de su tramo y de su margen, por ejemplo que un `zMm`
  negativo no pase de `−overhang.leftMm` (`packages/validator`, sprint 1.4). Una puerta más grande que el interior es válida (pasa en el reefer, con la línea de carga por
  debajo de la puerta): simplemente no limita.

### Versiones

`schemaVersion` (hoy `0.1.0`) es la misma versión que devuelve `/health` como `engineVersion`.
Agregar campos opcionales o relajar un límite sube el número menor; quitar o renombrar campos, o agregar
valores a un `enum` de salida, rompe a los clientes y se anuncia en el `CHANGELOG.md`.

### Validar los ejemplos

El proyecto no instala un validador de JSON Schema (no hay ADR que lo justifique). Para revisarlos a mano,
con Node instalado y fuera del repo:

```bash
npm i ajv@8          # en una carpeta temporal
```

```js
const Ajv = require("ajv/dist/2020");
const ajv = new Ajv({ strict: true, allErrors: true });
ajv.addSchema(require("./engine.schema.json"));
const validate = ajv.getSchema("urn:cargo-3d:engine:0.1.0#/$defs/EngineInput"); // o .../LoadPlan
console.log(validate(require("./examples/maritime/dry/20dv-heavy/input.json")), validate.errors);
```

Las **reglas físicas** (traslapes, apoyo, tramo y margen, puerta, peso, motivos y métricas) se revisan sin instalar
nada:

```bash
node scripts/contratos/check-examples.mjs   # o: pnpm examples:check
```

Es provisional: desde la sesión 1.4.2, `packages/validator` hace las dos revisiones.

### Ejemplos

Cada caso es una carpeta con `input.json` (un `EngineInput`) y `output.json` (un `LoadPlan`). Las carpetas siguen
la misma jerarquía que el catálogo de equipos (la categoría vive en el catálogo, no en el motor; ADR-0015):
`<modo>/<tipo de equipo>/<caso>/`. Todos son **sintéticos**, con medidas de fichas públicas, y cumplen las reglas
físicas: sin traslapes, todo dentro del equipo y por la puerta, nada flotando, peso respetado y métricas
correctas (y el margen OOG cuando lo hay). Las salidas son un plan válido escrito a mano, no lo que devolverá el
motor.

```
examples/
├─ maritime/                 marítimo (contenedores ISO)
│  ├─ bulk/                  granel seco cargado por la puerta
│  │  └─ 20bk-coffee-sacks/
│  ├─ dry/                   seco: estándar (DV, GP o ST) y high cube
│  │  ├─ 10dv-epal-pallets/
│  │  ├─ 20dv-heavy/
│  │  ├─ 20hc-tall-racks/
│  │  ├─ 40dv-gma-pallets/
│  │  ├─ 40hc-tall-pallets/
│  │  ├─ 45hc-gma-pallets/
│  │  └─ 45pw-epal-pallets/
│  ├─ reefer/                refrigerado
│  │  ├─ 20rf-epal-pallets/
│  │  ├─ 40rh-epal-pallets/
│  │  └─ 45pw-reefer-epal-pallets/
│  └─ special/               open top, flat rack y plataforma ISO
│     ├─ 20fr-in-gauge/
│     ├─ 20ot-door/
│     ├─ 20ot-roof/
│     ├─ 40fr-collapsible-steel/
│     ├─ 40fr-oog/
│     ├─ 40ot-door/
│     ├─ 40ot-roof/
│     └─ 40pf-wide-load/
├─ invalid/                  entradas rotas a propósito (schema/ y semantic/)
└─ road/                     autotransporte
   ├─ curtain/               lona corrediza (Conestoga)
   │  └─ 48cn-coils/
   ├─ dry-van/               caja seca
   │  ├─ 48dv-mixed/
   │  ├─ 53dv-contract-coverage/
   │  └─ 53dv-stacked/
   ├─ reefer-van/            caja refrigerada
   │  ├─ 48rv-pallets/
   │  └─ 53rv-pallets/
   └─ flatbed/               plataforma abierta
      ├─ 48fb-steel/
      └─ 53fb-long-pipes/
```

| Caso | Equipo | Qué muestra | Motivos de `unplaced` |
|---|---|---|---|
| `maritime/dry/10dv-epal-pallets` | 10' DV 8'6" (12G1) | 8 EPAL en dos niveles; el noveno no cabe y un compresor de 2,35 m no pasa por la puerta | `NO_SPACE_AVAILABLE`, `EXCEEDS_DOOR` |
| `maritime/dry/20dv-heavy` | 20' DV / estándar (22G1) | Cajones de 7 t: se acaba el peso con el 18 % del volumen | `WEIGHT_LIMIT_REACHED`, `EXCEEDS_DOOR`, `EXCEEDS_DIMENSIONS` |
| `maritime/dry/20hc-tall-racks` | 20' HC (25G1) | Racks de 2,5 m que no caben en un 20' estándar; uno de 2,6 m cabe adentro pero no pasa por la puerta | `EXCEEDS_DOOR` |
| `maritime/dry/40dv-gma-pallets` | 40' DV / estándar (42G1) | 20 pallets GMA: un carril «como viene» y otro girado 90° | ninguno |
| `maritime/dry/40hc-tall-pallets` | 40' HC (45G1) | Tarimas de 2,5 m que solo caben en HC; dos carriles con orientaciones distintas | `NO_SPACE_AVAILABLE`, `EXCEEDS_MAX_WEIGHT` |
| `maritime/dry/45hc-gma-pallets` | 45' HC (L5G1 / 55G1) | 24 pallets GMA combinando orientaciones (11 + 13) | `NO_SPACE_AVAILABLE` |
| `maritime/dry/45pw-epal-pallets` | 45' HC pallet wide (LEG1) | 33 EPAL, tres a lo ancho: no caben en un 45' normal | `NO_SPACE_AVAILABLE` |
| `maritime/reefer/20rf-epal-pallets` | 20' refrigerado (22R1) | 10 pallets EPAL bajo la línea de carga (altura útil por confirmar) | `NO_SPACE_AVAILABLE`, `EXCEEDS_DIMENSIONS` |
| `maritime/reefer/40rh-epal-pallets` | 40' HC refrigerado (45R1) | Altura útil bajo la línea de carga (puerta más alta que el interior); la tarima de 2,5 m ya no cabe | `EXCEEDS_DIMENSIONS` |
| `maritime/reefer/45pw-reefer-epal-pallets` | 45' HC refrigerado pallet wide (LNR1) | 3 EPAL a lo ancho; se acaba el peso antes que el espacio | `WEIGHT_LIMIT_REACHED` |
| `maritime/bulk/20bk-coffee-sacks` | 20' bulk (22B0), por la puerta | 16 tarimas de sacos en dos niveles; un big bag cabe adentro pero no pasa por la puerta | `EXCEEDS_DOOR` |
| `maritime/special/20ot-door` | Open top 20' (22U1), por la puerta | Una máquina de 2,3 m cabe adentro pero no pasa por la puerta | `EXCEEDS_DOOR` |
| `maritime/special/20ot-roof` | Open top 20' (22U1), por el techo | La misma carga con grúa: sin puerta, tramo = abertura del techo centrada en el interior (`xMm` 109, `zMm` 49); todo cabe | ninguno |
| `maritime/special/40fr-oog` | Flat rack 40' (42P3) | Margen OOG: transformador sobre los postes y placa más ancha (`zMm` negativo), marcados `outOfGauge` | `WEIGHT_LIMIT_REACHED` |
| `maritime/special/40fr-collapsible-steel` | Flat rack 40' plegable (42P8) | Lingotes de 12 t: 47 450 kg de útil; un bloque más ancho que el piso no cabe | `WEIGHT_LIMIT_REACHED`, `EXCEEDS_DIMENSIONS` |
| `maritime/special/40ot-door` | Open top 40' (42U1), por la puerta | Cabe la viga de 11,9 m, pero la máquina de 2,3 m no pasa por la puerta | `EXCEEDS_DOOR` |
| `maritime/special/40ot-roof` | Open top 40' (42U1), por el techo | Al revés: cabe la máquina, pero la viga no pasa por la abertura del techo | `EXCEEDS_DIMENSIONS` |
| `maritime/special/20fr-in-gauge` | Flat rack 20' (22P3) | Sin margen OOG: un gabinete más alto que los postes se rechaza (compárese con `40fr-oog`) | `EXCEEDS_DIMENSIONS` |
| `maritime/special/40pf-wide-load` | Plataforma ISO 40' | Margen lateral de 81 mm (NOM-012): 2,59 m de ancho viaja como OOG y 2,7 m ya no | `EXCEEDS_DIMENSIONS` |
| `road/curtain/48cn-coils` | Conestoga 48' | Sin puerta rígida (se carga por los lados), pero con alto limitado por la lona | `WEIGHT_LIMIT_REACHED`, `EXCEEDS_DIMENSIONS` |
| `road/dry-van/48dv-mixed` | Caja seca 48' | Ancho conservador (2 489): una pieza de 2,5 m ya no cabe; la puerta es más ancha que el interior | `EXCEEDS_DIMENSIONS` |
| `road/dry-van/53dv-stacked` | Caja seca 53' | Tarimas en dos niveles, apoyadas completas | `WEIGHT_LIMIT_REACHED` |
| `road/dry-van/53dv-contract-coverage` | Caja seca 53' | Las 6 orientaciones, orientaciones por omisión, `maxLoadOnTopKg` y `deliverySequence` (lo que se entrega primero va junto a la puerta) | ninguno |
| `road/flatbed/48fb-steel` | Plataforma 48' | Sin puerta; altura útil = 4 250 − 1 478 mm de cama (NOM-012) | `EXCEEDS_DIMENSIONS` |
| `road/flatbed/53fb-long-pipes` | Plataforma 53' | Tubos de 16 m que no caben en una de 48' | `TIME_LIMIT_REACHED` (ilustrativo) |
| `road/reefer-van/48rv-pallets` | Caja refrigerada 48' | Se acaba el peso con espacio de sobra | `WEIGHT_LIMIT_REACHED` |
| `road/reefer-van/53rv-pallets` | Caja refrigerada 53' | La puerta (2 585) es más baja que el interior (2 628): un rack de 2,6 m no pasa | `EXCEEDS_DOOR` |

Qué equipos y qué partes del contrato faltan por cubrir, y el paso a paso para agregar un caso, está en
[`examples/README.md`](examples/README.md).

