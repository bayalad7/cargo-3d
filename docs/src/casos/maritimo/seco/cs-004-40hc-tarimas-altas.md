# CS-004 · Tarimas de 2,5 m: solo caben en High Cube

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Marítimo › Seco |
| **Equipo** | `40hc` — Contenedor 40' HC (catálogo) contra `40dv` — 40' Dry |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `EXCEEDS_DIMENSIONS` (en el 40' Dry) |

## Situación

La misma carga se cotiza en un 40' Dry y en un 40' HC.

## Equipo

| Dato | Valor |
|---|---|
| 40' HC interior | 12 032 × 2 350 × 2 700 mm |
| 40' HC puerta | 2 340 × 2 597 mm |
| 40' Dry interior | 12 032 × 2 352 × 2 395 mm |
| maxPayloadKg | 28 600 kg (HC) · 28 750 kg (Dry) |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `tarima-alta` | Tarima de rollos de tela | 1 200 × 1 000 × 2 500 | 900 | 22 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- En el 40' HC: 22 colocadas en dos carriles con orientaciones distintas (19 800 kg).
- En el 40' Dry: 0 colocadas; las 22 con `EXCEEDS_DIMENSIONS`.

## Por qué

- 2 500 > 2 395 (alto del Dry) y la tarima no se puede acostar.
- En el HC: un carril de 1 200 de ancho (tarimas de 1 000 a lo largo: 12) y otro de 1 000 (tarimas de 1 200 a lo largo: 10) suman 2 200 ≤ 2 350 y dan 22.

## Qué validar con datos reales

- Confirmar la puerta de la unidad real: con menos de 2 500 de alto libre la tarima ya no entra.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
