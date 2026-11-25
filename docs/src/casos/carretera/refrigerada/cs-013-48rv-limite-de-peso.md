# CS-013 · Caja refrigerada 48': se acaba el peso antes que el espacio

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Caja refrigerada |
| **Equipo** | `48rv` — Caja refrigerada 48' (catálogo Alzper; interior del repo) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `WEIGHT_LIMIT_REACHED` |
| **Ejemplo del contrato** | `contracts/examples/road/reefer-van/48rv-pallets` |

## Situación

Lácteos refrigerados; el cliente manda 24 pallets.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 14 173 × 2 463 × 2 628 mm (repo, C.H. Robinson 2017) |
| Puerta | 2 496 × 2 585 mm (repo) |
| maxPayloadKg | 25 000 kg (máximo Alzper, estructural) |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `pallet` | Pallet de lácteos | 1 200 × 1 000 × 1 800 | 1 100 | 24 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- 22 colocadas (24 200 kg), 2 con `WEIGHT_LIMIT_REACHED`.

## Por qué

- 23 × 1 100 = 25 300 > 25 000.
- Espacio sobra: con carriles de 1 200 y 1 000 caben 25 por piso.

## Qué validar con datos reales

- El alto interior de la 48' de Alzper no está confirmado (P-013).
- Útil legal: falta restar el equipo de frío y el tracto según el camino.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
