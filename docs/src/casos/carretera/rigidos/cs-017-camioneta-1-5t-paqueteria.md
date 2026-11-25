# CS-017 · Camioneta de 1,5 t: paquetería

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Camiones rígidos y camionetas |
| **Equipo** | Camioneta 1,5 t con caja seca — equipo de la empresa (**simulado**) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `WEIGHT_LIMIT_REACHED` |

## Situación

Reparto urbano en Manzanillo; las cajas pesan más de lo normal (herramienta).

## Equipo

| Dato | Valor |
|---|---|
| Interior | 3 000 × 1 800 × 1 800 mm (simulado) |
| Puerta | 1 700 × 1 700 mm (simulada) |
| maxPayloadKg | 1 500 kg (simulado) |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `caja` | Caja de paquetería | 600 × 400 × 400 | 30 | 60 | las 6 | sí | — |

## Resultado esperado

- 50 colocadas (1 500 kg), 10 con `WEIGHT_LIMIT_REACHED`.

## Por qué

- Por espacio caben al menos 80 (5 × 4 × 4); se acaba antes el peso: 50 × 30 = 1 500.

## Qué validar con datos reales

- Las camionetas casi siempre tienen la caja hecha por un carrocero: cada unidad se da de alta con sus medidas.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
