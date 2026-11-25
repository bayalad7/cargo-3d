# CS-021 · Pantallas frágiles: nada encima

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Restricciones › Apilado |
| **Equipo** | `53dry` — Caja seca 53' (catálogo) |
| **Fase en que aplica** | 2.1 (`stackable` y `maxLoadOnTopKg`) |
| **Motivos esperados** | — |

## Situación

Carga mixta de electrónica y ferretería; en el patio nunca se pone peso sobre pantallas.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 16 003 × 2 540 × 2 794 mm |
| maxPayloadKg | 30 000 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `pantalla` | Caja de pantallas | 1 200 × 1 000 × 1 300 | 300 | 20 | "este lado arriba" (LWH, WLH) | no (nada encima) | — |
| `ferreteria` | Tarima de ferretería | 1 200 × 1 000 × 1 300 | 1 200 | 20 | "este lado arriba" (LWH, WLH) | sí | — |

## Resultado esperado

- Las 40 tarimas caben (26 por piso × 2 niveles).
- Ninguna tarima queda encima de una pantalla: las pantallas van arriba o solas.
- En la fase 1 el motor acepta `stackable` pero no lo aplica: este caso falla a propósito hasta la sesión 2.1.

## Por qué

- Es la regla del patio más común: «lo pesado abajo, lo frágil arriba».

## Qué validar con datos reales

- Preguntar a cada empresa qué nunca va abajo y qué nunca se encima.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
