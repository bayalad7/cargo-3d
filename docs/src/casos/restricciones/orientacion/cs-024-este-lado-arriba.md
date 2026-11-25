# CS-024 · «Este lado arriba»: acostada cabría, parada no

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Restricciones › Orientación |
| **Equipo** | `20dv` — Contenedor 20' Dry (catálogo) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `EXCEEDS_DIMENSIONS` |

## Situación

Un refrigerador industrial que no se puede acostar por el compresor.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 5 900 × 2 352 × 2 395 mm |
| Puerta | 2 340 × 2 292 mm |
| maxPayloadKg | 28 130 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `refrigerador` | Refrigerador industrial | 2 000 × 800 × 2 450 | 400 | 1 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- Con «este lado arriba»: `EXCEEDS_DIMENSIONS` (2 450 > 2 395).
- Variante con las 6 orientaciones permitidas: colocado acostado.

## Por qué

- La orientación sí se aplica desde la fase 1 (`allowedOrientations`).

## Qué validar con datos reales

- Preguntar en cada alta de pieza si se puede voltear: el error más caro es acostar algo que no se debe.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
