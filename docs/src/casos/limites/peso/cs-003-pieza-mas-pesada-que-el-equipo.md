# CS-003 · Pieza más pesada que todo el equipo

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Límites › Peso |
| **Equipo** | `20dv` — Contenedor 20' Dry (catálogo) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `EXCEEDS_MAX_WEIGHT` |

## Situación

El cliente pide un 20' para un transformador de 30 t porque «cabe de medidas».

## Equipo

| Dato | Valor |
|---|---|
| Interior | 5 900 × 2 352 × 2 395 mm |
| Puerta | 2 340 × 2 292 mm |
| maxPayloadKg | 28 130 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `transformador` | Transformador | 3 000 × 2 000 × 2 000 | 30 000 | 1 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- 0 colocadas; el transformador queda como `EXCEEDS_MAX_WEIGHT`.

## Por qué

- Cabe por medidas y por puerta (2 000 × 2 000 contra 2 340 × 2 292), así que el motivo no es de medida.
- 30 000 kg > 28 130 kg: no cabe ni con el contenedor vacío (ADR-0016, orden de los motivos).

## Qué validar con datos reales

- La app debería sugerir un equipo con más útil, por ejemplo el flat rack 40' plegable de 47 450 kg (`40fr-collapsible`).

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
