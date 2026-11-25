# CS-018 · Camioneta de 3,5 t: línea blanca

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Camiones rígidos y camionetas |
| **Equipo** | Camioneta 3,5 t con caja seca — equipo de la empresa (**simulado**) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `EXCEEDS_DOOR` |

## Situación

Entrega de línea blanca a una tienda; la estufa industrial es más alta que la puerta.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 4 500 × 2 200 × 2 100 mm (simulado) |
| Puerta | 2 100 × 2 000 mm (simulada) |
| maxPayloadKg | 3 500 kg (simulado) |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `refri` | Refrigerador | 900 × 800 × 1 800 | 120 | 10 | "este lado arriba" (LWH, WLH) | no | — |
| `estufa` | Estufa industrial | 1 500 × 900 × 2 050 | 250 | 1 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- 10 refrigeradores colocados (1 200 kg).
- Estufa con `EXCEEDS_DOOR`: 2 050 ≤ 2 100 adentro, pero 2 050 > 2 000 en la puerta.

## Por qué

- 5 a lo largo (900) y 2 a lo ancho (800) en un solo nivel; no se apilan.

## Qué validar con datos reales

- Muchas camionetas tienen cortina en lugar de puertas: la cortina enrollada quita alto en la entrada.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
