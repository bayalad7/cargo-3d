# CS-002 · 20' Dry lleno justo por peso

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Marítimo › Seco |
| **Equipo** | `20dv` — Contenedor 20' Dry (catálogo, Hapag-Lloyd) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `WEIGHT_LIMIT_REACHED` |

## Situación

Piezas fundidas muy densas: el peso se acaba con el contenedor casi vacío.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 5 900 × 2 352 × 2 395 mm |
| Puerta | 2 340 × 2 292 mm |
| maxPayloadKg | 28 130 kg |
| Tara | 2 350 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `cajon` | Cajón de piezas fundidas | 1 000 × 1 000 × 1 000 | 2 813 | 11 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- 10 colocadas (28 130 kg, 100 % del peso), 1 no colocada con `WEIGHT_LIMIT_REACHED`.
- Ocupación de volumen de solo 30,1 %.

## Por qué

- 10 × 2 813 = 28 130 kg = carga útil exacta.
- El cajón 11 sí cabe por medidas y por puerta, pero ya no por peso.

## Qué validar con datos reales

- Con 10 cajones de 2,8 t el piso trabaja muy concentrado: en la fase 2 entra la distribución de peso a lo largo (CTU Code).
- El VGM sería 2 350 + 28 130 = 30 480 kg, igual al bruto máximo de la ficha: justo en el límite. Revisar la placa CSC de la unidad real.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
