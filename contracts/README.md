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
- **Lo que no cupo no es un error:** sale en `unplaced` con su motivo (ADR-0016). Los errores de validación son
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

Desde el sprint 1.4, `packages/validator` hace esta revisión y además las reglas físicas (sin traslapes,
todo dentro del equipo, peso respetado).

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
│  ├─ dry/                   seco: estándar (DV, GP o ST) y high cube
│  │  ├─ 20dv-heavy/
│  │  ├─ 40dv-gma-pallets/
│  │  ├─ 40hc-tall-pallets/
│  │  └─ 45hc-gma-pallets/
│  ├─ reefer/                refrigerado
│  │  ├─ 20rf-epal-pallets/
│  │  └─ 40rh-epal-pallets/
│  └─ special/               open top y flat rack
│     ├─ 20ot-door/
│     ├─ 20ot-roof/
│     └─ 40fr-oog/
└─ road/                     autotransporte
   ├─ dry-van/               caja seca
   │  ├─ 48dv-mixed/
   │  ├─ 53dv-contract-coverage/
   │  └─ 53dv-stacked/
   ├─ reefer-van/            caja refrigerada
   │  └─ 53rv-pallets/
   └─ flatbed/               plataforma abierta
      ├─ 48fb-steel/
      └─ 53fb-long-pipes/
```

| Caso | Equipo | Qué muestra | Motivos de `unplaced` |
|---|---|---|---|
| `maritime/dry/20dv-heavy` | 20' DV / estándar (22G1) | Cajones de 7 t: se acaba el peso con el 18 % del volumen | `WEIGHT_LIMIT_REACHED`, `EXCEEDS_DOOR`, `EXCEEDS_DIMENSIONS` |
| `maritime/dry/40dv-gma-pallets` | 40' DV / estándar (42G1) | 20 pallets GMA: un carril «como viene» y otro girado 90° | ninguno |
| `maritime/dry/40hc-tall-pallets` | 40' HC (45G1) | Tarimas de 2,5 m que solo caben en HC; dos carriles con orientaciones distintas | `NO_SPACE_AVAILABLE`, `EXCEEDS_MAX_WEIGHT` |
| `maritime/dry/45hc-gma-pallets` | 45' HC (L5G1 / 55G1) | 24 pallets GMA combinando orientaciones (11 + 13) | `NO_SPACE_AVAILABLE` |
| `maritime/reefer/20rf-epal-pallets` | 20' refrigerado (22R1) | 10 pallets EPAL bajo la línea de carga (altura útil por confirmar) | `NO_SPACE_AVAILABLE`, `EXCEEDS_DIMENSIONS` |
| `maritime/reefer/40rh-epal-pallets` | 40' HC refrigerado (45R1) | Altura útil bajo la línea de carga (puerta más alta que el interior); la tarima de 2,5 m ya no cabe | `EXCEEDS_DIMENSIONS` |
| `maritime/special/20ot-door` | Open top 20' (22U1), por la puerta | Una máquina de 2,3 m cabe adentro pero no pasa por la puerta | `EXCEEDS_DOOR` |
| `maritime/special/20ot-roof` | Open top 20' (22U1), por el techo | La misma carga con grúa: sin puerta, tramo = abertura del techo centrada en el interior (`xMm` 109, `zMm` 49); todo cabe | ninguno |
| `maritime/special/40fr-oog` | Flat rack 40' (42P3) | Margen OOG: transformador sobre los postes y placa más ancha (`zMm` negativo), marcados `outOfGauge` | `WEIGHT_LIMIT_REACHED` |
| `road/dry-van/48dv-mixed` | Caja seca 48' | Ancho conservador (2 489): una pieza de 2,5 m ya no cabe; la puerta es más ancha que el interior | `EXCEEDS_DIMENSIONS` |
| `road/dry-van/53dv-stacked` | Caja seca 53' | Tarimas en dos niveles, apoyadas completas | `WEIGHT_LIMIT_REACHED` |
| `road/dry-van/53dv-contract-coverage` | Caja seca 53' | Las 6 orientaciones, orientaciones por omisión, `maxLoadOnTopKg` y `deliverySequence` (lo que se entrega primero va junto a la puerta) | ninguno |
| `road/flatbed/48fb-steel` | Plataforma 48' | Sin puerta; altura útil = 4 250 − 1 478 mm de cama (NOM-012) | `EXCEEDS_DIMENSIONS` |
| `road/flatbed/53fb-long-pipes` | Plataforma 53' | Tubos de 16 m que no caben en una de 48' | `TIME_LIMIT_REACHED` (ilustrativo) |
| `road/reefer-van/53rv-pallets` | Caja refrigerada 53' | La puerta (2 585) es más baja que el interior (2 628): un rack de 2,6 m no pasa | `EXCEEDS_DOOR` |

Qué equipos y qué partes del contrato faltan por cubrir, y el paso a paso para agregar un caso, está en
[`examples/README.md`](examples/README.md).

