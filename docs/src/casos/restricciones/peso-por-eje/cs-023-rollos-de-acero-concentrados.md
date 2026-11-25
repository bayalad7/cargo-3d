# CS-023 · Rollos de acero: el peso no puede quedar en un solo extremo

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Restricciones › Peso por eje |
| **Equipo** | `53dry` — Caja seca 53' (catálogo) en T3-S2 |
| **Fase en que aplica** | 2.1 (peso por eje y centro de gravedad) |
| **Motivos esperados** | — |

## Situación

Tres rollos de 8 t; el montacarguista los deja juntos al fondo para acabar rápido.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 16 003 × 2 540 × 2 794 mm |
| maxPayloadKg | 30 000 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `rollo` | Rollo de lámina | 1 500 × 1 500 × 1 500 | 8 000 | 3 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- Los 3 rollos caben por peso (24 000 kg) y por espacio.
- El plan los reparte a lo largo para no pasar el límite de los ejes; juntos al fondo sobrecargarían el tándem del tracto.

## Por qué

- La NOM-012 limita el peso por eje además del bruto; el cálculo por eje entra en la sesión 2.1.

## Qué validar con datos reales

- Pesar en báscula por eje en el periodo de pruebas y comparar con el cálculo.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
