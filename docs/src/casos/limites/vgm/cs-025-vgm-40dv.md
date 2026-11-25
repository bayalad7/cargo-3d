# CS-025 · VGM de un 40' Dry

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Límites › VGM |
| **Equipo** | `40dv` — Contenedor 40' Dry (catálogo, Hapag-Lloyd) |
| **Fase en que aplica** | 1.5 (la API reporta el VGM) |
| **Motivos esperados** | — |

## Situación

Exportación de café; la naviera pide el VGM antes del corte documental.

## Equipo

| Dato | Valor |
|---|---|
| Tara | 3 750 kg |
| Bruto máximo | 32 500 kg |
| maxPayloadKg | 28 750 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `saco` | Tarima de sacos de café | 1 200 × 1 000 × 1 150 | 1 050 | 24 | "este lado arriba" (LWH, WLH) | sí | — |

## Resultado esperado

- 24 tarimas colocadas (25 200 kg).
- VGM = tara + carga = 3 750 + 25 200 = 28 950 kg (más amarres y calzas, que se pesan aparte).

## Por qué

- SOLAS VI/2: el VGM es el peso bruto verificado; el motor da la carga y la tara viene en `tareKg`.

## Qué validar con datos reales

- El VGM oficial sale de báscula o del método 2 (sumar todo); el cálculo de la app es una estimación.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
