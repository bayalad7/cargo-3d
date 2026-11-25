# CS-011 · Caja seca 53': la misma carga por camino ET y por camino B

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Caja seca |
| **Equipo** | `53dry` — Caja seca 53' Fruehauf (catálogo) en T3-S2 |
| **Fase en que aplica** | 1 (el límite legal lo calcula la API) |
| **Motivos esperados** | `WEIGHT_LIMIT_REACHED` (camino B) |

## Situación

Distribución de abarrotes de Manzanillo a Guadalajara; parte de la ruta es camino tipo B.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 16 003 × 2 540 × 2 794 mm |
| Puerta | 2 496 × 2 794 mm (repo, C.H. Robinson 2017) |
| Capacidad estructural | 30 000 kg |
| Tara caja / tracto | 6 680 kg / 8 800 kg (tracto **simulado**) |
| maxPayloadKg camino ET | mín(30 000; 46 500 − 8 800 − 6 680 = 31 020) = **30 000** |
| maxPayloadKg camino B | mín(30 000; 38 000 − 8 800 − 6 680 = 22 520) = **22 520** |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `tarima` | Tarima de abarrotes | 1 200 × 1 000 × 1 300 | 450 | 52 | "este lado arriba" (LWH, WLH) | sí | — |

## Resultado esperado

- Camino ET: 52 colocadas en 2 niveles (23 400 kg).
- Camino B: 50 colocadas (22 500 kg) y 2 con `WEIGHT_LIMIT_REACHED`.

## Por qué

- Por piso caben al menos 26 tarimas (2 a lo ancho, 13 a lo largo) y se apilan 2 niveles (2 600 ≤ 2 794).
- El motor solo recibe `maxPayloadKg`; el tipo de camino lo resuelve la API (ADR-0016, punto 2).

## Qué validar con datos reales

- La tara real del tracto sale de su tarjeta de circulación o de la báscula.
- La ruta manda: el tramo más restrictivo fija la carga del viaje completo.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
