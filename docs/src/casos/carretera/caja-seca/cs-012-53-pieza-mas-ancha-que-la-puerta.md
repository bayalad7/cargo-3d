# CS-012 · Caja seca 53': cabe adentro pero no pasa por la puerta

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Caja seca |
| **Equipo** | `53dry` — Caja seca 53' (catálogo) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `EXCEEDS_DOOR` |

## Situación

Un tablero de 2,52 m de ancho: el interior tiene 2 540 pero la puerta 2 496.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 16 003 × 2 540 × 2 794 mm |
| Puerta | 2 496 × 2 794 mm (repo) |
| maxPayloadKg | 30 000 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `tablero` | Tablero eléctrico | 3 000 × 2 520 × 1 500 | 1 800 | 1 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- 0 colocadas; el tablero con `EXCEEDS_DOOR`.

## Por qué

- Parado a lo largo (LWH) cabe adentro: 2 520 ≤ 2 540. Girado (WLH) mide 3 000 de ancho: no cabe.
- La única orientación que cabe no pasa por la puerta: 2 520 > 2 496.

## Qué validar con datos reales

- Una plataforma de 48' o 53' resuelve el caso (carga por arriba con montacargas).

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
