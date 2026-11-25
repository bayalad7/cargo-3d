# CS-027 · Llega una tarima más alta de lo que se declaró

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Operación › Cambios en el patio |
| **Equipo** | `40hc` — Contenedor 40' HC (catálogo) |
| **Fase en que aplica** | 2.4 (volver a optimizar en la web) |
| **Motivos esperados** | — |

## Situación

El cliente declaró tarimas de 1,30 m; en el patio una mide 1,45 m y ya no caben dos de alto en esa columna (2 750 > 2 700).

## Equipo

| Dato | Valor |
|---|---|
| Interior | 12 032 × 2 350 × 2 700 mm |
| Puerta | 2 340 × 2 597 mm |
| maxPayloadKg | 28 600 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `tarima` | Tarima declarada | 1 200 × 1 000 × 1 300 | 500 | 39 | "este lado arriba" (LWH, WLH) | sí | — |
| `tarima-real` | Tarima medida en el patio | 1 200 × 1 000 × 1 450 | 500 | 1 | "este lado arriba" (LWH, WLH) | sí | — |

## Resultado esperado

- Plan original (40 tarimas de 1,30 m): caben en 2 niveles.
- Plan corregido: la tarima de 1,45 m no se puede encimar con otra (1 450 + 1 300 = 2 750 > 2 700) y va sola en su columna. Las 40 siguen cabiendo: con dos carriles hay 22 lugares por piso.
- Se registra la diferencia (declarado contra medido) para el periodo de pruebas.

## Por qué

- ADR-0016: las piezas se redondean hacia arriba; el error viene de la declaración, no del motor.

## Qué validar con datos reales

- Medir en el patio una muestra de cada tipo de pieza antes de cargar.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
