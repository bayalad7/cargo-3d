# CS-020 · Tortón (C3): 8 o 10 llantas cambian la carga

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Camiones rígidos y camionetas |
| **Equipo** | Tortón C3 con caja seca — equipo de la empresa (**simulado**) |
| **Fase en que aplica** | 1 (el límite legal lo calcula la API) |
| **Motivos esperados** | `WEIGHT_LIMIT_REACHED` |

## Situación

Material de construcción a una obra; la empresa tiene tortones de 8 y de 10 llantas.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 9 500 × 2 450 × 2 600 mm (simulado) |
| Puerta | 2 400 × 2 500 mm (simulada) |
| Tara del camión con caja | 9 500 kg (simulada) |
| Capacidad del equipo | 17 000 kg (simulada) |
| maxPayloadKg con 10 llantas (ET) | mín(17 000; 27 500 − 9 500 = 18 000) = **17 000** |
| maxPayloadKg con 8 llantas (ET) | mín(17 000; 24 000 − 9 500 = 14 500) = **14 500** |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `cemento` | Tarima de cemento (30 sacos de 50 kg) | 1 200 × 1 000 × 1 000 | 1 500 | 12 | "este lado arriba" (LWH, WLH) | sí | — |

## Resultado esperado

- 10 llantas: 11 colocadas (16 500 kg), 1 con `WEIGHT_LIMIT_REACHED`.
- 8 llantas: 9 colocadas (13 500 kg), 3 con `WEIGHT_LIMIT_REACHED`.

## Por qué

- El espacio sobra (28 lugares: 14 por piso y 2 niveles); manda el peso.
- PBV de la NOM-012: C3 = 27,5 t con 10 llantas y 24,0 t con 8 (ET y A).

## Qué validar con datos reales

- Por eso el alta de un tortón pide el número de llantas.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
