# Casos de prueba del motor

Cargas de ejemplo en JSON, **independientes del lenguaje** del motor. Cada entrada sigue
`contracts/engine.schema.json`.

- `casos/`: los casos de uso simulados de `docs/src/casos/` (sesión 1.1.5, ADR-0020).
- `br/`: instancias académicas BR1 a BR15 (sprint 1.4).

## `casos/`

Una carpeta por caso, `cs-NNN/`, con:

- `input.json`: la entrada del motor (un `EngineInput`).
- `input-<variante>.json`: una entrada por cada equipo, camino o dato que el caso compara.
- `<entrada>.expect.json`: qué caso la produce y qué resultado se espera.

La prueba `TestCaseFixturesAreValidAndConsistent` (`backend/internal/engine/fixtures_test.go`) recorre todas las
entradas y revisa que:

- pasen `DecodeInput` (contrato y validación de la entrada, ADR-0023);
- tengan su `.expect.json`;
- las cuentas del caso cuadren: colocadas + no colocadas = unidades de la entrada; peso cargado = peso total − peso de
  las no colocadas; el peso cargado no pasa de `maxPayloadKg`; solo hay sobredimensionadas si algún tramo tiene margen.

Desde el sprint 1.3, el algoritmo se compara contra `expected`.

### Formato del `.expect.json`

```json
{
  "case": "CS-008",
  "doc": "docs/src/casos/maritimo/especial/cs-008-40fr-sobredimensionada.md",
  "variant": "40' Flat rack",
  "phase": "1",
  "validInput": true,
  "expected": {
    "placedCount": 1,
    "outOfGaugeCount": 1,
    "loadedWeightKg": 6000,
    "unplaced": [{ "itemId": "torre", "quantity": 1, "reason": "EXCEEDS_DIMENSIONS" }]
  },
  "rules": [{ "from": "2.1", "text": "…" }]
}
```

| Campo | Qué es |
|---|---|
| `case`, `doc` | El caso y su archivo `.md` |
| `variant` | Qué equipo, camino o dato usa esta entrada |
| `phase` | Sesión o fase en que el caso aplica completo (`1`, `1.5`, `2.1`, `2.3`, `2.4`) |
| `validInput` | La entrada debe pasar la validación |
| `expected` | Colocadas, sobredimensionadas, peso cargado y no colocadas con su motivo, como en el contrato |
| `rules` | Lo que no se dice con números, con la sesión desde la que el motor debe cumplirlo (opcional) |

### De dónde salen las medidas

- Por id de `contracts/equipment.json`, sin cambiarlo. La descripción del equipo de cada entrada dice la fuente.
- Lo que el catálogo no publica:
  - **Puerta de las cajas secas** (`53dry`, `48dry`): 2 496 × 2 794 del repo (C.H. Robinson 2017), en CS-011, 012, 021,
    022 y 026.
  - **Cajas refrigeradas** (`48rv`, `53rv`): interior y puerta de los ejemplos del repo (C.H. Robinson 2017), que no
    coinciden con el catálogo (P-020), en CS-013 y 014.
  - **Puerta del `40rf`**: no está publicada; la variante de CS-007 va sin puerta (P-013). No cambia el resultado: la
    tarima no cabe en el interior.
  - **Carga útil del `45pw`**: 29 000 kg, simulada.
- Equipos de la empresa **simulados** (ADR-0020): camionetas, rabón, tortón y la caja forrada de CS-026.
- **Open top por el techo** (CS-009): el tramo es la abertura de 5 338 × 2 230 en `xMm = 278`, `zMm = 60`. Centrada
  daría 278,5; queda medio milímetro descentrada, toda dentro del interior.

### Fixtures

| Caso | Entrada | Variante | Fase | Colocadas | No colocadas |
|---|---|---|---|---|---|
| [CS-001](../../docs/src/casos/maritimo/seco/cs-001-20dv-justo-por-espacio.md) | `cs-001/input.json` | 20' Dry | 1 | 20 | 1 NO_SPACE_AVAILABLE |
| [CS-002](../../docs/src/casos/maritimo/seco/cs-002-20dv-justo-por-peso.md) | `cs-002/input.json` | 20' Dry | 1 | 10 | 1 WEIGHT_LIMIT_REACHED |
| [CS-003](../../docs/src/casos/limites/peso/cs-003-pieza-mas-pesada-que-el-equipo.md) | `cs-003/input.json` | 20' Dry | 1 | 0 | 1 EXCEEDS_MAX_WEIGHT |
| [CS-004](../../docs/src/casos/maritimo/seco/cs-004-40hc-tarimas-altas.md) | `cs-004/input-40dv.json` | 40' Dry | 1 | 0 | 22 EXCEEDS_DIMENSIONS |
| [CS-004](../../docs/src/casos/maritimo/seco/cs-004-40hc-tarimas-altas.md) | `cs-004/input.json` | 40' HC | 1 | 22 | — |
| [CS-005](../../docs/src/casos/maritimo/seco/cs-005-45pw-europallets.md) | `cs-005/input.json` | 45' HC pallet wide | 1 | 33 | 1 NO_SPACE_AVAILABLE |
| [CS-006](../../docs/src/casos/maritimo/refrigerado/cs-006-20rf-tarima-sobre-alto-interior.md) | `cs-006/input.json` | 20' Reefer | 1 | 10 | 2 NO_SPACE_AVAILABLE, 1 EXCEEDS_DIMENSIONS |
| [CS-007](../../docs/src/casos/maritimo/refrigerado/cs-007-40rf-contra-40rh.md) | `cs-007/input-40rf.json` | 40' Reefer estándar | 1 | 0 | 20 EXCEEDS_DIMENSIONS |
| [CS-007](../../docs/src/casos/maritimo/refrigerado/cs-007-40rf-contra-40rh.md) | `cs-007/input.json` | 40' Reefer HC | 1 | 20 | — |
| [CS-008](../../docs/src/casos/maritimo/especial/cs-008-40fr-sobredimensionada.md) | `cs-008/input.json` | 40' Flat rack | 1 | 1 | 1 EXCEEDS_DIMENSIONS |
| [CS-009](../../docs/src/casos/maritimo/especial/cs-009-20ot-carga-por-techo.md) | `cs-009/input-techo.json` | 20' Open top por el techo | 1 | 1 | — |
| [CS-009](../../docs/src/casos/maritimo/especial/cs-009-20ot-carga-por-techo.md) | `cs-009/input.json` | 20' Open top por la puerta | 1 | 0 | 1 EXCEEDS_DOOR |
| [CS-010](../../docs/src/casos/maritimo/especial/cs-010-40pf-carga-ancha.md) | `cs-010/input.json` | 40' Platform | 1 | 2 | 1 EXCEEDS_DIMENSIONS |
| [CS-011](../../docs/src/casos/carretera/caja-seca/cs-011-53-dos-niveles-y-tipo-de-camino.md) | `cs-011/input-camino-b.json` | Camino B (22 520 kg, tracto simulado) | 1 | 50 | 2 WEIGHT_LIMIT_REACHED |
| [CS-011](../../docs/src/casos/carretera/caja-seca/cs-011-53-dos-niveles-y-tipo-de-camino.md) | `cs-011/input.json` | Camino ET (30 000 kg) | 1 | 52 | — |
| [CS-012](../../docs/src/casos/carretera/caja-seca/cs-012-53-pieza-mas-ancha-que-la-puerta.md) | `cs-012/input.json` | Caja seca 53' | 1 | 0 | 1 EXCEEDS_DOOR |
| [CS-013](../../docs/src/casos/carretera/refrigerada/cs-013-48rv-limite-de-peso.md) | `cs-013/input.json` | Caja refrigerada 48' | 1 | 22 | 2 WEIGHT_LIMIT_REACHED |
| [CS-014](../../docs/src/casos/carretera/refrigerada/cs-014-53rv-rack-mas-alto-que-la-puerta.md) | `cs-014/input.json` | Caja refrigerada 53' | 1 | 26 | 1 EXCEEDS_DOOR |
| [CS-015](../../docs/src/casos/carretera/plataforma/cs-015-53fb-tubos-largos.md) | `cs-015/input.json` | Plataforma 53' | 1 | 8 | 1 EXCEEDS_DIMENSIONS |
| [CS-017](../../docs/src/casos/carretera/rigidos/cs-017-camioneta-1-5t-paqueteria.md) | `cs-017/input.json` | Camioneta 1,5 t (simulada) | 1 | 50 | 10 WEIGHT_LIMIT_REACHED |
| [CS-018](../../docs/src/casos/carretera/rigidos/cs-018-camioneta-3-5t-linea-blanca.md) | `cs-018/input.json` | Camioneta 3,5 t (simulada) | 1 | 10 | 1 EXCEEDS_DOOR |
| [CS-019](../../docs/src/casos/carretera/rigidos/cs-019-rabon-tipo-de-camino.md) | `cs-019/input-camino-d.json` | Camino D (6 000 kg) | 1 | 8 | 6 WEIGHT_LIMIT_REACHED |
| [CS-019](../../docs/src/casos/carretera/rigidos/cs-019-rabon-tipo-de-camino.md) | `cs-019/input.json` | Camino ET (10 000 kg) | 1 | 14 | — |
| [CS-020](../../docs/src/casos/carretera/rigidos/cs-020-torton-cemento-llantas.md) | `cs-020/input-8-llantas.json` | 8 llantas (14 500 kg) | 1 | 9 | 3 WEIGHT_LIMIT_REACHED |
| [CS-020](../../docs/src/casos/carretera/rigidos/cs-020-torton-cemento-llantas.md) | `cs-020/input.json` | 10 llantas (17 000 kg) | 1 | 11 | 1 WEIGHT_LIMIT_REACHED |
| [CS-021](../../docs/src/casos/restricciones/apilado/cs-021-fragil-no-debajo.md) | `cs-021/input.json` | Caja seca 53' | 2.1 | 40 | — |
| [CS-022](../../docs/src/casos/restricciones/multidestino/cs-022-tres-destinos-lifo.md) | `cs-022/input.json` | Caja seca 48' | 2.1 | 24 | — |
| [CS-024](../../docs/src/casos/restricciones/orientacion/cs-024-este-lado-arriba.md) | `cs-024/input-6-orientaciones.json` | Las 6 orientaciones | 1 | 1 | — |
| [CS-024](../../docs/src/casos/restricciones/orientacion/cs-024-este-lado-arriba.md) | `cs-024/input.json` | Este lado arriba | 1 | 0 | 1 EXCEEDS_DIMENSIONS |
| [CS-025](../../docs/src/casos/limites/vgm/cs-025-vgm-40dv.md) | `cs-025/input.json` | 40' Dry | 1.5 | 24 | — |
| [CS-026](../../docs/src/casos/operacion/equipo-alterado/cs-026-caja-forrada-por-dentro.md) | `cs-026/input-caja-real.json` | Caja forrada por dentro (simulada) | 2.3 | 0 | 4 EXCEEDS_DIMENSIONS |
| [CS-026](../../docs/src/casos/operacion/equipo-alterado/cs-026-caja-forrada-por-dentro.md) | `cs-026/input.json` | Caja de catálogo | 2.3 | 4 | — |
| [CS-027](../../docs/src/casos/operacion/cambios-en-patio/cs-027-tarima-mas-alta-de-lo-declarado.md) | `cs-027/input-medido.json` | Plan corregido (una tarima medida en el patio) | 2.4 | 40 | — |
| [CS-027](../../docs/src/casos/operacion/cambios-en-patio/cs-027-tarima-mas-alta-de-lo-declarado.md) | `cs-027/input.json` | Plan original (declarado) | 2.4 | 40 | — |
| [CS-028](../../docs/src/casos/carretera/sobredimensionada/cs-028-plataforma-carga-ancha-con-permiso.md) | `cs-028/input-con-permiso.json` | Con permiso tipo 2 (margen 355 mm por lado) | 1 | 3 | — |
| [CS-028](../../docs/src/casos/carretera/sobredimensionada/cs-028-plataforma-carga-ancha-con-permiso.md) | `cs-028/input.json` | Sin permiso (margen 5 mm por lado) | 1 | 2 | 1 EXCEEDS_DIMENSIONS |

### Casos que todavía no son fixtures

| Caso | Por qué | Cuándo |
|---|---|---|
| [CS-016](../../docs/src/casos/carretera/plataforma/cs-016-40fb-sin-altura-de-cama.md) | La entrada no se puede armar: falta la altura de cama del `40fb` | 2.3 (equipos de la empresa) |
| [CS-023](../../docs/src/casos/restricciones/peso-por-eje/cs-023-rollos-de-acero-concentrados.md) | El contrato 0.1.0 no tiene campos de peso por eje ni centro de gravedad | 2.1 |
| [CS-029](../../docs/src/casos/carretera/plataforma/cs-029-escalonada-dos-tramos.md) | Necesita dos tramos y el contrato 0.1.0 acepta uno (P-012) | 2.1.6 |
