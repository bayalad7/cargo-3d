# CS-006 · Reefer 20': tarima más alta que el interior

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Marítimo › Refrigerado |
| **Equipo** | `20rf` — Contenedor 20' Reefer (catálogo, Hapag-Lloyd) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `NO_SPACE_AVAILABLE`, `EXCEEDS_DIMENSIONS` |
| **Ejemplo del contrato** | `contracts/examples/maritime/reefer/20rf-epal-pallets` |

## Situación

Exportación de fruta a 2 °C; una tarima de flores viene más alta que el resto.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 5 450 × 2 280 × 2 159 mm (alto interior, no línea roja) |
| Puerta | 2 290 × 2 264 mm |
| maxPayloadKg | 29 140 kg |
| Tara | 2 860 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `epal` | Pallet EPAL de fruta | 1 200 × 800 × 2 000 | 900 | 12 | "este lado arriba" (LWH, WLH) | no | — |
| `tarima-alta` | Tarima de flores | 1 200 × 800 × 2 200 | 700 | 1 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- 10 pallets colocados, 2 con `NO_SPACE_AVAILABLE`.
- La tarima de flores con `EXCEEDS_DIMENSIONS`.

## Por qué

- Dos carriles: 1 200 + 800 = 2 000 ≤ 2 280; 5 450 / 800 = 6 y 5 450 / 1 200 = 4 → 10.
- 2 200 > 2 159.

## Qué validar con datos reales

- La línea roja real queda abajo del alto interior: si mide menos de 2 000 ningún pallet es válido. Se mide en el periodo de pruebas.
- El aire frío necesita espacio libre arriba de la carga: posible regla de la fase 2.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
