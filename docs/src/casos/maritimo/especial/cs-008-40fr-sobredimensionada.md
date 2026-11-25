# CS-008 · Flat rack 40': carga más alta que los postes

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Marítimo › Especial |
| **Equipo** | `40fr` — Contenedor 40' Flatrack (catálogo) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `EXCEEDS_DIMENSIONS` |
| **Ejemplo del contrato** | `contracts/examples/maritime/special/40fr-oog` |

## Situación

Un proyecto de energía embarca un transformador y una sección de torre en un flat rack.

## Equipo

| Dato | Valor |
|---|---|
| Piso entre postes | 11 652 × 2 208 mm, alto de postes 1 955 mm |
| Margen (overhang) | arriba 400 mm (límite interno), lados 196 mm por lado (NOM-012) |
| maxPayloadKg | 28 000 kg (supuesto del repo) |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `transformador` | Transformador | 2 500 × 2 100 × 2 100 | 6 000 | 1 | "este lado arriba" (LWH, WLH) | no | — |
| `torre` | Sección de torre | 1 500 × 1 500 × 2 500 | 1 000 | 1 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- Transformador colocado como carga sobredimensionada (rebasa 145 mm arriba; `outOfGaugeCount` = 1).
- Torre con `EXCEEDS_DIMENSIONS`: 2 500 > 1 955 + 400.

## Por qué

- El margen de 400 mm es un límite interno de Cargo 3D, no de la naviera (ADR-0017 y auditoría G1).

## Qué validar con datos reales

- CMA CGM admite en principio hasta 1 700 mm arriba con aprobación por reserva: la torre se podría embarcar. La app debería pedir el margen aprobado en la reserva.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
