# CS-028 · Plataforma 48': intercambiador de 3,2 m de ancho con permiso de la NOM-040

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Carga sobredimensionada |
| **Equipo** | `48fb` — Plataforma flatbed 48' (catálogo) en T3-S2 |
| **Fase en que aplica** | 1 (el margen lateral ya existe en el contrato); el permiso lo arma la API (sesiones 2.1.5 y 2.3) |
| **Motivos esperados** | `EXCEEDS_DIMENSIONS` (sin permiso) |

## Situación

Un fabricante de Altamira manda un intercambiador de calor de una sola pieza (no se puede desarmar) a una refinería.
Mide 3,2 m de ancho: más que los 2,6 m que permite la NOM-012. La empresa tramita el permiso especial de la SCT.

## Equipo

| Dato | Valor |
|---|---|
| Piso | 14 630 × 2 590 mm |
| Alto útil | 2 772 mm = 4 250 (NOM-012) − 1 478 de cama |
| Largo total de la combinación | 21 000 mm (tractocamión + plataforma 48', **simulado**) |
| maxPayloadKg | 30 000 kg (capacidad del repo; el límite legal de T3-S2 camino ET da 31 020 con taras **simuladas**) |
| Margen lateral sin permiso | 5 mm por lado: (2 600 − 2 590) / 2 (NOM-012) |
| Margen lateral con permiso tipo 2 | 355 mm por lado: (3 300 − 2 590) / 2 (NOM-040, Tabla B) |

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `intercambiador` | Intercambiador de calor (indivisible) | 12 000 × 3 200 × 2 500 | 22 000 | 1 | "este lado arriba" (LWH, WLH) | no | — |
| `refacciones` | Tarima de refacciones del mismo cliente | 1 200 × 1 000 × 1 000 | 600 | 2 | "este lado arriba" (LWH, WLH) | sí | — |

## Resultado esperado

- **Sin permiso** (margen de 5 mm): el intercambiador queda con `EXCEEDS_DIMENSIONS`; las 2 tarimas se colocan.
- **Con permiso tipo 2** (margen de 355 mm): el intercambiador se coloca centrado como sobredimensionado
  (`outOfGauge: true`, rebasa 305 mm por lado) y las 2 tarimas en los 2 630 mm libres del piso, **dentro** del ancho
  de la plataforma.
- La salida del plan debe decir qué exige el permiso (lo agrega la API, no el motor): combinación **tipo 2**, dos
  torretas ámbar (unidad motriz y parte posterior de la carga), dos letreros «EXCESO DE ANCHO» de 1,20 × 0,60 m
  (delante y atrás), **sin unidad piloto**, y el horario de circulación.

## Por qué

- NOM-040-SCT-2-2012, Tabla B: tipo 2 = largo hasta 28,00 m y ancho hasta 3,30 m, alto libre. 21 m y 3,2 m → tipo 2.
  Con 3,5 m de ancho sería tipo 3 (hasta 3,70 m) y llevaría **una** unidad piloto; con más de 3,70 m, tipo 5 y **dos**.
- Tabla D: tipo 2 lleva dos torretas, dos señales de advertencia y ninguna unidad piloto.
- Numeral 8.1.1: los elementos de sujeción que no rebasan 10 cm por lado no cuentan en el ancho total.
- Numeral 5.5: con permiso solo se circula de lunes a viernes (en caminos ET2, A2, B2, C y D de 6:00 a 18:30 y de
  0:00 a 6:00; de 18:30 a 23:59 solo con dos carros piloto) y sábados de 6:00 a 14:00. Domingos y días de descanso
  obligatorio, nada.
- El permiso es para objetos **indivisibles** (numeral 3): las tarimas de refacciones no pueden rebasar el ancho
  normal aunque viajen en el mismo equipo. Por eso el margen es por pieza, no por tramo.
- Alto: 1 478 + 2 500 = 3 978 ≤ 4 250, así que no excede por alto.

## Qué validar con datos reales

- Un permiso real de la SCT para un movimiento de este tipo: qué datos trae y si coinciden con la Tabla B.
- Si la NOM-040 sigue vigente o hay versión nueva (su texto cita la NOM-012 de 2008, no la de 2017): se revisa en la
  sesión 2.1.5 junto con el DOF (P-023).
- Hoy el contrato pone el margen por **tramo**: para que las tarimas no usen el margen del permiso hace falta un campo
  por pieza (por ejemplo `indivisible: true`). Es un cambio de contrato pendiente (P-023).

## Fuente

NOM-040-SCT-2-2012, DOF del 20 de marzo de 2013: numerales 3, 5.4, 5.5, 8.1.1, 8.5 y Tablas B y D, leídos del texto
publicado (https://diariooficial.gob.mx/nota_detalle.php?codigo=5292917&fecha=20/03/2013) el 2026-11-25.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado con la NOM-040 verificada (Tablas B y D) y equipo del catálogo. |
