# CS-010 · Plataforma ISO 40': carga más ancha que el piso

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Marítimo › Especial |
| **Equipo** | `40pf` — Contenedor 40' Platform (catálogo) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `EXCEEDS_DIMENSIONS` |
| **Ejemplo del contrato** | `contracts/examples/maritime/special/40pf-wide-load` |

## Situación

Módulos prefabricados que rebasan un poco el piso de la plataforma.

## Equipo

| Dato | Valor |
|---|---|
| Piso | 12 190 × 2 438 mm |
| Alto libre | 2 383 mm, derivado: 4 250 − 648 − 1 219 |
| Margen (overhang) | 81 mm por lado: (2 600 − 2 438) / 2 |
| maxPayloadKg | 28 000 kg (supuesto del repo) |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `modulo` | Módulo de casa prefabricada | 5 000 × 2 590 × 2 200 | 9 000 | 2 | "este lado arriba" (LWH, WLH) | no | — |
| `modulo-ancho` | Módulo ancho | 2 000 × 2 700 × 1 000 | 1 000 | 1 | solo como viene (LWH) | no | — |

## Resultado esperado

- 2 módulos colocados como sobredimensionados (76 mm por lado).
- Módulo ancho con `EXCEEDS_DIMENSIONS`: 2 700 > 2 438 + 81 + 81.

## Por qué

- El ancho máximo en carretera es 2 600 mm (NOM-012); de ahí sale el margen.
- El módulo ancho solo viaja como viene: girado (`WLH`) mediría 2 000 de ancho y sí cabría en el piso.

## Qué validar con datos reales

- El alto libre de 2 383 es derivado: se mide la plataforma montada en el chasis (P-004).

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
| 2026-11-26 | Simulado | El módulo ancho pasa a solo `LWH`, como el ejemplo `40pf-wide-load`: con `WLH` medía 2 000 de ancho y sí cabía, así que su motivo no podía ser `EXCEEDS_DIMENSIONS` (sesión 1.1.5). |
