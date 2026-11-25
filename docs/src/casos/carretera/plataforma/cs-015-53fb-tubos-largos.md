# CS-015 · Plataforma 53': atados de tubo de 12 m

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Plataforma |
| **Equipo** | `48fb` (versión 53') — Plataforma flatbed |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `EXCEEDS_DIMENSIONS` |

## Situación

Tubería de acero para una obra; un tramo viene de 16,5 m.

## Equipo

| Dato | Valor |
|---|---|
| Piso | 16 154 × 2 590 mm |
| Alto útil | 2 772 mm = 4 250 (NOM-012) − 1 478 de cama |
| maxPayloadKg | 30 000 kg (repo) |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `atado` | Atado de tubo | 12 000 × 1 200 × 600 | 2 000 | 8 | las 6 | sí | — |
| `tubo-largo` | Tubo de 16,5 m | 16 500 × 400 × 400 | 900 | 1 | las 6 | sí | — |

## Resultado esperado

- 8 atados colocados: 2 a lo ancho y 4 niveles (16 000 kg).
- Tubo de 16,5 m con `EXCEEDS_DIMENSIONS`.

## Por qué

- 2 × 1 200 = 2 400 ≤ 2 590 y 4 × 600 = 2 400 ≤ 2 772.
- 16 500 > 16 154 y el contrato 0.1.0 no tiene vuelo trasero.

## Qué validar con datos reales

- La NOM-012 permite algo de vuelo trasero con señalamiento: posible regla de la fase 2.
- Amarres y cadenas: cuánto espacio dejan es regla del patio.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
