# CS-016 · Plataforma plana 40': la ficha no da la altura de la cama

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Plataforma |
| **Equipo** | `40fb` — Plataforma plana 40' Fruehauf PM-HT2-40W (catálogo) |
| **Fase en que aplica** | 2.3 (equipos de la empresa) |
| **Motivos esperados** | — (la entrada no se puede armar) |

## Situación

La empresa tiene una plataforma de 40' y la elige del catálogo.

## Equipo

| Dato | Valor |
|---|---|
| Piso | 12 125 × 2 457 mm (longitud y ancho de carga) |
| Altura de cama | **sin publicar** |
| Capacidad | 34 000 kg (estructural) |
| Tara | 5 650 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `viga` | Viga de acero | 12 000 × 600 × 800 | 3 000 | 6 | LWH, WLH | sí | — |

## Resultado esperado

- La API no puede calcular el alto útil (4 250 − cama) y pide la altura de la cama antes de optimizar.
- Con una cama **simulada** de 1 500 mm el alto útil sería 2 750 y entrarían 3 niveles de 800 (2 400): las 6 vigas en 2 carriles de 600.

## Por qué

- El motor no adivina medidas: sin alto útil no hay `sections` válidas.

## Qué validar con datos reales

- Es el ejemplo de por qué la app necesita equipos de la empresa con medidas propias (ADR-0019).

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
