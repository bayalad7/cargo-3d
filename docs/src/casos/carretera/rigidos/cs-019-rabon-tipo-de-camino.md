# CS-019 · Rabón (C2): la ruta baja la carga permitida

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Camiones rígidos y camionetas |
| **Equipo** | Rabón C2 con caja seca — equipo de la empresa (**simulado**) |
| **Fase en que aplica** | 1 (el límite legal lo calcula la API) |
| **Motivos esperados** | `WEIGHT_LIMIT_REACHED` (camino D) |

## Situación

Reparto a tiendas de la sierra de Colima; el último tramo es camino tipo D.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 7 300 × 2 450 × 2 500 mm (simulado) |
| Puerta | 2 400 × 2 400 mm (simulada) |
| Tara del camión con caja | 7 000 kg (simulada) |
| Capacidad del equipo | 10 000 kg (simulada) |
| maxPayloadKg camino ET | mín(10 000; 19 000 − 7 000 = 12 000) = **10 000** |
| maxPayloadKg camino D | mín(10 000; 13 000 − 7 000 = 6 000) = **6 000** |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `tarima` | Tarima de abarrotes | 1 200 × 1 000 × 1 100 | 700 | 14 | "este lado arriba" (LWH, WLH) | sí | — |

## Resultado esperado

- Camino ET: 14 colocadas (9 800 kg).
- Camino D: 8 colocadas (5 600 kg) y 6 con `WEIGHT_LIMIT_REACHED`.

## Por qué

- Por piso caben 12 (2 a lo ancho, 6 a lo largo) y se apilan 2 niveles (2 200 ≤ 2 500).
- PBV de la NOM-012, tabla B-2: C2 = 19,0 t (ET) y 13,0 t (D).

## Qué validar con datos reales

- La tara y el PBV salen de la tarjeta de circulación de cada rabón.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
