# CS-026 · Caja seca forrada por dentro: el catálogo ya no sirve

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Operación › Equipo alterado |
| **Equipo** | `53dry` modificada — equipo de la empresa (**simulado**) |
| **Fase en que aplica** | 2.3 (equipos de la empresa) |
| **Motivos esperados** | `EXCEEDS_DIMENSIONS` (con la caja real) |

## Situación

La empresa forró su caja con madera para proteger la carga; nadie actualizó las medidas.

## Equipo

| Dato | Valor |
|---|---|
| Interior de catálogo | 16 003 × 2 540 × 2 794 mm |
| Interior real (forro de madera y techo reforzado) | 15 950 × 2 480 × 2 700 mm (simulado) |
| maxPayloadKg | 30 000 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `panel` | Panel aislante (viaja acostado) | 6 000 × 2 490 × 300 | 900 | 4 | "este lado arriba" (LWH, WLH) | sí | — |

## Resultado esperado

- Con la caja de catálogo: 4 paneles colocados, acostados y apilados (2 490 ≤ 2 540, y pasan por la puerta de 2 496).
- Con la caja real: los 4 con `EXCEEDS_DIMENSIONS` (2 490 > 2 480). De canto sí cabrían, pero el panel no se puede parar.

## Por qué

- Un plan hecho con medidas de catálogo falla en el patio: por eso cada empresa registra sus unidades (ADR-0019).

## Qué validar con datos reales

- Al dar de alta una caja, preguntar si tiene forro, rieles o refuerzos y medir entre los puntos más angostos.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
| 2026-11-26 | Simulado | El panel pasa de 2 500 a 2 490 de ancho: con la puerta de la caja seca 53' (2 496 × 2 794, del repo) uno de 2 500 no entraba ni en la caja de catálogo (sesión 1.1.5). |
