# CS-014 · Caja refrigerada 53': puerta más baja que el interior

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Caja refrigerada |
| **Equipo** | `53rv` — Caja refrigerada 53' (interior y puerta del repo) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `EXCEEDS_DOOR` |
| **Ejemplo del contrato** | `contracts/examples/road/reefer-van/53rv-pallets` |

## Situación

Carga mixta con un rack alto de panadería.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 15 697 × 2 463 × 2 628 mm (repo) |
| Puerta | 2 496 × 2 585 mm (repo) |
| maxPayloadKg | 30 000 kg (máximo Alzper, estructural) |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `pallet` | Pallet de cajas de verdura | 1 200 × 1 000 × 1 800 | 900 | 26 | "este lado arriba" (LWH, WLH) | no | — |
| `rack` | Rack de pan | 1 200 × 1 000 × 2 600 | 500 | 1 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- 26 pallets colocados (23 400 kg).
- Rack con `EXCEEDS_DOOR`: 2 600 ≤ 2 628 adentro, pero 2 600 > 2 585 en la puerta.

## Por qué

- Pasa con frecuencia en cajas refrigeradas: el marco de la puerta y las cortinas quitan altura.

## Qué validar con datos reales

- Hyundai publica 2 680 de interior para su 53': con esa unidad el rack seguiría sin pasar si la puerta es igual. Se mide la puerta.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
