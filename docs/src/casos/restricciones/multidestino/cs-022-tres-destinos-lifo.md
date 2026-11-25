# CS-022 · Tres entregas en una caja: lo último que se carga sale primero

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Restricciones › Varios destinos |
| **Equipo** | `48dry` — Caja seca 48' (catálogo) |
| **Fase en que aplica** | 2.1 (`deliverySequence`) |
| **Motivos esperados** | — |

## Situación

Una caja reparte en Colima, Guadalajara y León, en ese orden.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 14 427 × 2 489 × 2 794 mm |
| maxPayloadKg | 30 000 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `colima` | Tarima para Colima | 1 200 × 1 000 × 1 400 | 600 | 8 | "este lado arriba" (LWH, WLH) | sí | 1 |
| `guadalajara` | Tarima para Guadalajara | 1 200 × 1 000 × 1 400 | 600 | 8 | "este lado arriba" (LWH, WLH) | sí | 2 |
| `leon` | Tarima para León | 1 200 × 1 000 × 1 400 | 600 | 8 | "este lado arriba" (LWH, WLH) | sí | 3 |

## Resultado esperado

- Las 24 tarimas caben.
- Las de Colima (entrega 1) quedan junto a la puerta; las de León, al fondo.
- Ninguna tarima de una entrega posterior bloquea a una anterior.

## Por qué

- `deliverySequence` = 1 se descarga primero (contrato 0.1.0, reservado para la 2.1).

## Qué validar con datos reales

- Preguntar si el patio carga por destino o por tamaño cuando hay prisa.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
