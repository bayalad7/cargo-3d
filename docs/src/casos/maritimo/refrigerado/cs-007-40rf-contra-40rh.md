# CS-007 · Reefer 40' estándar contra High Cube

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Marítimo › Refrigerado |
| **Equipo** | `40rf` — 40' Reefer 8'6" (CIMC) contra `40rh` — 40' Reefer HC (Hapag-Lloyd) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `EXCEEDS_DIMENSIONS` (en el 40' RF) |

## Situación

La naviera ofrece reefer estándar más barato; la carga ya viene paletizada a 2,3 m.

## Equipo

| Dato | Valor |
|---|---|
| 40' RF interior | 11 583 × 2 290 × 2 252 mm (puerta no publicada) |
| 40' RH interior | 11 599 × 2 290 × 2 425 mm |
| 40' RH puerta | 2 290 × 2 557 mm |
| maxPayloadKg | 28 200 kg (RF) · 29 580 kg (RH) |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `tarima` | Tarima de carne congelada | 1 200 × 1 000 × 2 300 | 1 000 | 20 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- En el 40' RH: 20 colocadas (20 000 kg).
- En el 40' RF: 0 colocadas, las 20 con `EXCEEDS_DIMENSIONS`.

## Por qué

- 2 300 > 2 252 (alto interior del RF) y la tarima no se acuesta.
- En el RH: carriles de 1 200 y 1 000 (2 200 ≤ 2 290): 11 + 9 = 20.

## Qué validar con datos reales

- La puerta del 42R1 no está publicada (P-013): se mide en la unidad.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
