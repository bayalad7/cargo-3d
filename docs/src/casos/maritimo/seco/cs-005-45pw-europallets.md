# CS-005 · 45' pallet wide: 3 europallets a lo ancho

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Marítimo › Seco |
| **Equipo** | `45pw` — Contenedor 45' HC pallet wide (catálogo, PalletSpec) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `NO_SPACE_AVAILABLE` |

## Situación

Importación de Europa en europallets; el agente pregunta si conviene el 45' pallet wide en lugar del 45' HC normal.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 13 563 × 2 442 × 2 697 mm |
| Puerta | 2 416 × 2 585 mm |
| maxPayloadKg | 29 000 kg (**simulado**: la útil del seco PW no se publica) |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `epal` | Europallet EPAL 1 | 1 200 × 800 × 1 600 | 600 | 34 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- 33 colocadas (11 filas de 3), 1 no colocada con `NO_SPACE_AVAILABLE`.
- Referencia: en el 45' HC normal (`45hc`, 2 352 de ancho) no caben 3 a lo ancho (2 400); con dos carriles mixtos se llega a 27.

## Por qué

- 3 × 800 = 2 400 ≤ 2 442 y 13 563 / 1 200 = 11 filas → 33.
- No se pueden apilar: 2 × 1 600 = 3 200 > 2 697.

## Qué validar con datos reales

- La útil del seco pallet wide sale de la placa CSC de la unidad; aquí es simulada.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
