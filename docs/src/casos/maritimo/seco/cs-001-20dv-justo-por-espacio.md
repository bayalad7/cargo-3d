# CS-001 · 20' Dry lleno justo por espacio

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Marítimo › Seco |
| **Equipo** | `20dv` — Contenedor 20' Dry (catálogo, Hapag-Lloyd) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `NO_SPACE_AVAILABLE` |

## Situación

Un exportador manda cajas de refacciones hechas a la medida del contenedor: 5 a lo largo, 2 a lo ancho y 2 de alto. Pide una caja de más «por si cabe».

## Equipo

| Dato | Valor |
|---|---|
| Interior | 5 900 × 2 352 × 2 395 mm |
| Puerta | 2 340 × 2 292 mm |
| maxPayloadKg | 28 130 kg (útil de la ficha) |
| Tara | 2 350 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `caja` | Caja de refacciones | 1 180 × 1 176 × 1 197 | 500 | 21 | las 6 | sí | — |

## Resultado esperado

- 20 colocadas, 1 no colocada con `NO_SPACE_AVAILABLE`.
- Ocupación de volumen de 99,96 %; peso 10 000 de 28 130 kg.

## Por qué

- 5 900 / 1 180 = 5, 2 352 / 1 176 = 2 y 2 395 / 1 197 = 2: caben exactamente 20.
- La caja 21 no cabe en ningún acomodo: el volumen de 21 cajas (34,88 m³) pasa el del contenedor (33,24 m³).
- Es el caso «justo»: cualquier error de redondeo en el motor rompe el acomodo.

## Qué validar con datos reales

- La unidad real puede tener unos milímetros menos que la ficha de Hapag: con 2 394 de alto ya no caben 2 niveles.
- Revisar que el patio deje holgura para el montacargas (regla del patio).

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
