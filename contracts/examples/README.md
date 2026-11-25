# Ejemplos del contrato del motor

Cada caso es una carpeta `<modo>/<tipo de equipo>/<caso>/` con `input.json` (un `EngineInput`) y `output.json`
(un `LoadPlan`), según `contracts/engine.schema.json`. Las carpetas siguen la jerarquía del catálogo de equipos.
Cada caso trae además `equipo.svg`: un dibujo del equipo generado desde su `input.json`, con la misma escala, ángulo y
paleta en todos (se ve en claro y oscuro). No es una fotografía: el repositorio es público y las fotos de
fabricantes no se pueden redistribuir. Se regenera con `pnpm examples:render`.

Este archivo es **la lista de lo que falta**: se actualiza en el mismo commit en que se agrega un caso.
Hay dos tipos de cobertura: por **equipo** (¿hay un ejemplo para cada tipo de equipo?) y por **contrato**
(¿cada campo y cada motivo aparece en algún ejemplo?).

Estados: ✅ hay ejemplo · ⏳ falta, se puede agregar ya (falta el dato con fuente) · 🔒 bloqueado por el contrato
(necesita varios tramos, `maxItems` > 1) · 🚫 fuera del alcance del motor por ahora.

## 1. Cobertura por equipo

### Marítimo (`maritime/`)

Códigos ISO 6346 verificados en fuente secundaria (Wikipedia, BIC); el resto, de memoria, se verifica al agregar el caso.
La ISO 6346:2022 cambió el código de largo 45' de `L` a `5`: el catálogo guarda los dos.

| Ruta | Equipo | Código ISO | Estado | Qué falta |
|---|---|---|---|---|
| `dry/10dv-*` | 10' estándar de 8'6" | 12G1 | ✅ `10dv-epal-pallets` | — (ficha Triton) |
| `dry/20dv-*` | 20' estándar (DV / GP / ST) | 22G1 | ✅ `20dv-heavy` | — |
| `dry/40dv-*` | 40' estándar (DV / GP / ST) | 42G1 | ✅ `40dv-gma-pallets` | — |
| `dry/40hc-*` | 40' high cube | 45G1 | ✅ `40hc-tall-pallets` | — |
| `dry/45hc-*` | 45' high cube | L5G1 (1995) / 55G1 (2022) | ✅ `45hc-gma-pallets` | — |
| `dry/45pw-*` | 45' pallet wide | LEG1 (DP World) | ✅ `45pw-epal-pallets` | Confirmar puerta y alto (fragmento de búsqueda) |
| `dry/20hc-*` | 20' high cube | 25G1 (de memoria) | ✅ `20hc-tall-racks` | Tara y carga útil de una placa (una sola fuente: Mainfreight) |
| `reefer/20rf-*` | 20' refrigerado | 22R1 | ✅ `20rf-epal-pallets` | Usa la unidad de ejemplo de Hapag-Lloyd (interior 5 450 × 2 280 × 2 159, puerta 2 290 × 2 264, tara 2 860, útil 29 140). El alto es alto interior, no línea roja: medir la línea roja en la unidad (hoja §4). El largo 5 724 del catálogo del cliente no coincide con ninguna ficha |
| `reefer/40rf-*` | 40' refrigerado estándar | 42R1 | ⏳ | CIMC publica interior, tara y útil; falta la puerta (P-013) |
| `reefer/40rh-*` | 40' refrigerado high cube | 45R1 | ✅ `40rh-epal-pallets` | — |
| `reefer/45pw-reefer-*` | 45' refrigerado HC pallet wide | LNR1 | ✅ `45pw-reefer-epal-pallets` | — (ficha Triton; el alto es alto interior, no línea roja) |
| `special/20ot-*` | Open top 20', por puerta y por techo | 22U1 | ✅ `20ot-door`, `20ot-roof` | — |
| `special/40ot-*` | Open top 40', por puerta y por techo | 42U1 | ✅ `40ot-door`, `40ot-roof` | — |
| `special/20fr-*` | Flat rack 20' | 22P1 (fijo) / 22P3 (plegable) | ✅ `20fr-in-gauge` | Carga útil supuesta (28 000) |
| `special/40fr-collapsible-*` | Flat rack 40' plegable | 42P8 | ✅ `40fr-collapsible-steel` | — (Textainer, serie 52,5 t: carga muy pesada) |
| `special/40fr-*` | Flat rack 40' | 42P1 / 42P3 | ✅ `40fr-oog` | El margen de arriba (400) es un límite interno conservador, no un tope de naviera (CMA CGM admite en principio 1 700; Maersk pide aprobación por reserva); el de los lados (196) sale de la NOM-012: (2 600 − 2 208) / 2 |
| `special/40pf-*` | Plataforma ISO 40' | 42P0 (de memoria) | ✅ `40pf-wide-load` | Alto libre derivado 2 383 = 4 250 − 648 (alto de la plataforma) − 1 219 (chasis, Triton y Conway), por medir en patio; carga útil por carretera (28 000), supuesta |
| `special/20pf-*` | Plataforma ISO 20' | 29P0 | ⏳ | Alto libre: falta la altura con fuente de un chasis de 20' |
| `bulk/20bk-*` | Granel seco 20', cargado por la puerta | 22B0 | ✅ `20bk-coffee-sacks` | — (ficha Triton; el granel suelto por las escotillas no es un problema de acomodo) |
| `special/tank-*` | Tanque | 22T… | 🚫 | Granel líquido, no se acomoda con cajas (se retiró del catálogo el 2026-11-25) |

### Autotransporte (`road/`)

| Ruta | Equipo | Estado | Qué falta |
|---|---|---|---|
| `dry-van/48dv-*` | Caja seca 48' | ✅ `48dv-mixed` | Interior completo de un fabricante mexicano (hoy ancho Remosa + largo y puerta C.H. Robinson); los 30 000 kg son capacidad estructural Fruehauf, no útil legal NOM-012 |
| `dry-van/53dv-*` | Caja seca 53' | ✅ `53dv-stacked`, `53dv-contract-coverage` | — (puerta conservadora de C.H. Robinson, 2 496 × 2 794) |
| `reefer-van/53rv-*` | Caja refrigerada (thermo) 53' | ✅ `53rv-pallets` | Carga útil 30 000 = máximo Alzper 53' (estructural, no útil legal); interior C.H. Robinson 2017 por confirmar |
| `reefer-van/48rv-*` | Caja refrigerada (thermo) 48' | ✅ `48rv-pallets` | Carga útil 25 000 = máximo Alzper 48' (estructural, no útil legal); interior C.H. Robinson 2017 por confirmar |
| `flatbed/40fb-*` | Plataforma 40' | ⏳ | Altura de cama (Fruehauf PM-HT2-40W da largo, ancho, tara y capacidad) |
| `flatbed/48fb-*` | Plataforma 48' | ✅ `48fb-steel` | Altura de cama real (hoy 1 478 mm de un fabricante de EE. UU.) |
| `flatbed/53fb-*` | Plataforma 53' | ✅ `53fb-long-pipes` | Ídem |
| `flatbed/extendable-*` | Plataforma extensible | ⏳ | Largo según la extensión usada (un tramo por configuración) |
| `curtain/48cn-*` | Plataforma con lona corrediza (Conestoga) 48' | ✅ `48cn-coils` | Confirmar el ancho útil entre rieles (2 438) |
| `step-deck/*` | Escalonada (drop deck) | 🔒 | Dos tramos con distinta altura de piso |
| `double-drop/*` | Doble caída | 🔒 | Tres tramos (pozo central) |
| `lowboy/*`, `rgn/*` | Cama baja y cuello desmontable (RGN) | 🔒 | Varios tramos; la carga suele pedir permiso por sobredimensión |
| `rigid/rabon-*` | Camión rabón (C2) | ⏳ | **Medidas del patio:** placa o tarjeta de circulación y medir la caja de 1–2 unidades. Ya verificado: NOM-012 C2 = 19,0 t en ET/A |
| `rigid/torton-*` | Camión tortón (C3) | ⏳ | Ídem, más el número de llantas (C3 = 24,0 t con 8 llantas, 27,5 t con 10, en ET/A) |
| `rigid/van-3t-*` | Camioneta de 3,5 t | ⏳ | Medidas de la caja |
| `specialized/hot-shot-*` | Hot shot | ⏳ | Poco usado en México; confirmar antes de agregarlo |
| `specialized/car-hauler-*` | Nodriza | 🔒 | Varios niveles |
| `combination/full-*`, `b-train-*` | Full (T3-S2-R4) y tren B | 🚫 | Varios equipos en un viaje (sprint 2.4) |

### Ferroviario (`rail/`)

| Ruta | Equipo | Estado |
|---|---|---|
| `boxcar/*`, `gondola/*`, `double-stack/*` | Furgón, góndola, doble estiba | 🚫 futuro (fase 4) |

## 2. Cobertura del contrato

| Parte del contrato | Dónde aparece | Estado |
|---|---|---|
| Equipo con puerta / sin puerta | contenedores y cajas secas / plataformas | ✅ |
| Puerta más grande que el interior | `maritime/reefer/40rh-epal-pallets`, `road/dry-van/48dv-mixed` | ✅ |
| Puerta más baja que el interior | `road/reefer-van/53rv-pallets` | ✅ |
| Open top por puerta y por techo (dos entradas del catálogo) | `maritime/special/20ot-door`, `20ot-roof` | ✅ |
| Margen OOG (`overhang`), `outOfGauge`, `outOfGaugeCount` > 0, `zMm` negativo | `maritime/special/40fr-oog` | ✅ |
| `tareKg` presente / ausente | contenedores / caja seca y plataformas | ✅ |
| Orientaciones `LWH`, `WLH` | casi todos | ✅ |
| Orientaciones `LHW`, `HLW`, `WHL`, `HWL` | `road/dry-van/53dv-contract-coverage` | ✅ |
| `allowedOrientations` ausente (se permiten las 6) | `20dv-heavy`, `48fb-steel`, `53fb-long-pipes`, `53dv-contract-coverage` | ✅ |
| Piezas apiladas, con apoyo completo | `53dv-stacked`, `53dv-contract-coverage` | ✅ |
| `stackable`, `maxLoadOnTopKg`, `deliverySequence` | `53dv-contract-coverage` (y `stackable` en varios) | ✅ (el motor los aplica hasta el sprint 2.1) |
| `description` ausente en una pieza | `53dv-contract-coverage` | ✅ |
| Todo cabe (`unplaced` vacío) | `40dv-gma-pallets`, `53dv-contract-coverage` | ✅ |
| `EXCEEDS_DIMENSIONS` | `20dv-heavy`, `40rh-epal-pallets`, `48fb-steel`, `40ot-roof`, `20fr-in-gauge` y otros | ✅ |
| `EXCEEDS_DOOR` | `20dv-heavy`, `20ot-door`, `40ot-door`, `53rv-pallets` | ✅ |
| `EXCEEDS_MAX_WEIGHT` | `40hc-tall-pallets` | ✅ |
| `WEIGHT_LIMIT_REACHED` | `20dv-heavy`, `53dv-stacked`, `40fr-oog`, `48rv-pallets` | ✅ |
| `NO_SPACE_AVAILABLE` | `40hc-tall-pallets`, `45hc-gma-pallets`, `20rf-epal-pallets` | ✅ |
| `TIME_LIMIT_REACHED` | `53fb-long-pipes` (ilustrativo) | ✅ |
| Varios tramos (`sections` > 1) | — | 🔒 hasta que el contrato lo acepte |
| Entradas **inválidas** a propósito | `invalid/schema` (15) e `invalid/semantic` (2); ver §4 | ✅ (falta «tramos encimados», 🔒 hasta que haya varios tramos) |

## 3. Cómo agregar un caso

1. **Elige la ruta** en las tablas de arriba (o crea una rama nueva, en inglés y sin acentos, por el ADR-0009).
   Nombre del caso: `<largo><tipo>-<qué muestra>`, por ejemplo `road/dry-van/48dv-mixed-boxes/`.
2. **Consigue los datos con fuente.** Medidas del equipo: ficha del fabricante o de la naviera, o la placa del
   equipo del cliente. Si hace falta investigar, corre `/buscador <equipo>`. Los datos del cliente se guardan en
   `docs/datos/` (no se sube a git) y aquí solo van medidas y pesos, **anonimizados**. Redondeo: capacidades del
   equipo hacia abajo y piezas hacia arriba (ADR-0016).
3. **Escribe `input.json`.** Un solo tramo `main` en `(0, 0, 0)`; en plataformas, `heightMm` = 4 250 − altura de la
   cama; en `description` del equipo indica que es sintético y de dónde salen las medidas.
4. **Escribe `output.json`**, un plan válido hecho a mano (no tiene que ser el mejor), que cumpla:
   - cada colocación dentro del tramo (o de su margen `overhang`, marcando `outOfGauge: true` y sumándola a
     `outOfGaugeCount`) y sin traslapes (caras que se tocan sí se permiten);
   - `dxMm`, `dyMm`, `dzMm` iguales a las medidas que da la orientación (1.ª letra sobre x, 2.ª sobre z, 3.ª sobre y);
   - orientaciones dentro de `allowedOrientations`; si hay puerta, `dz ≤ ancho` y `dy ≤ alto` de la puerta;
   - nada flotando: lo que tiene `yMm > 0` se apoya completo sobre otras piezas;
   - peso cargado ≤ `maxPayloadKg`; colocadas + no colocadas = `quantity` de cada pieza;
   - cada `reason` es el primero que se cumple en el orden del contrato;
   - métricas exactas (las ocupaciones con 4 decimales).
5. **Valida** contra el esquema con ajv (ver `contracts/README.md`). Las reglas físicas del paso 4 se revisan
   con `node scripts/contratos/check-examples.mjs` (provisional, hasta la sesión 1.4.2).
6. **Genera su imagen** con `pnpm examples:render` (escribe `equipo.svg` en cada caso; si el tipo de equipo es nuevo,
   agrégalo en `scripts/contratos/render-examples.mjs`).
7. **Actualiza este archivo** (estado ⏳ → ✅ y la tabla de cobertura del contrato) y la tabla de casos de
   `contracts/README.md`. Commit: `contracts: ejemplo <ruta>`.

## 4. Entradas inválidas (`invalid/`)

Cada archivo rompe **una sola regla** de la misma entrada mínima válida, y su nombre dice cuál. No se llaman
`input.json`, para que las pruebas que recorren los casos válidos no los tomen.

- `invalid/schema/`: los rechaza el propio esquema (medida en cero, campo desconocido, versión distinta…).
- `invalid/semantic/`: **pasan el esquema** y los rechaza la validación del motor (sesión 1.1.4):

| Archivo | Regla | De dónde sale |
|---|---|---|
| `duplicate-item-id.json` | Los ids de las piezas no se repiten | JSON Schema no lo puede expresar |
| `door-with-overhang.json` | Un equipo con puerta no lleva margen OOG | ADR-0017 |

Los **casos de uso** (simulados, ADR-0020) no van aquí sino en `docs/src/casos/`, y sus fixtures en `engine/fixtures/casos/` (sesión 1.1.5).
