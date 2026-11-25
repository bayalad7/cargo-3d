# CS-009 · Open top 20': máquina que solo entra por el techo

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Marítimo › Especial |
| **Equipo** | `20ot` — Contenedor 20' Open Top (catálogo, Hapag-Lloyd) |
| **Fase en que aplica** | 1 |
| **Motivos esperados** | `EXCEEDS_DOOR` (por la puerta) |

## Situación

La máquina mide 2,33 m de alto y el cliente tiene grúa en su planta.

## Equipo

| Dato | Valor |
|---|---|
| Interior | 5 895 × 2 350 × 2 340 mm (alto de los costados; 2 380 al centro) |
| Puerta | 2 338 × 2 280 mm |
| Abertura del techo | 5 338 × 2 230 mm |
| maxPayloadKg | 30 050 kg |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `maquina` | Torno CNC | 3 000 × 2 000 × 2 330 | 6 500 | 1 | "este lado arriba" (LWH, WLH) | no | — |

## Resultado esperado

- Por la puerta: `EXCEEDS_DOOR` (2 330 > 2 280).
- Por el techo (tramo del tamaño de la abertura, sin puerta): colocada.

## Por qué

- Mismo equipo, dos entradas en el catálogo: por puerta y por techo (ADR-0017).
- Se usa 2 340 de alto, el de los costados, que es el conservador.

## Qué validar con datos reales

- Confirmar que la grúa del patio o de la planta alcanza; si no, la máquina no se carga.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con datos simulados y fichas públicas del catálogo. |
