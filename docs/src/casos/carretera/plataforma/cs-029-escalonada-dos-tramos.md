# CS-029 · Plataforma escalonada 53': la máquina alta solo cabe en el tramo bajo

| Campo | Valor |
|---|---|
| **Estado** | Simulado (2026-11-25) |
| **Categoría** | Carretera › Plataforma |
| **Equipo** | Plataforma escalonada (*step deck*) 53' — equipo de la empresa (**simulado**) |
| **Fase en que aplica** | 2.1.6 (varios tramos, P-012); el contrato 0.1.0 acepta un solo tramo |
| **Motivos esperados** | `EXCEEDS_DIMENSIONS` (si se modela con un solo tramo) |

## Situación

Una empresa mueve una prensa de 3,1 m de alto. En una plataforma normal no pasa del límite de 4,25 m; en la escalonada
sí, porque el piso de atrás está más bajo. Aprovecha el tramo alto de adelante para cajas de refacciones.

## Equipo

Dos tramos con distinta altura de piso. El origen es el fondo (adelante, junto al cuello) y la parte de atrás es
x = 16 155 (ADR-0015).

| Tramo | x (mm) | Largo (mm) | Ancho (mm) | Altura de la cama (mm) | `yMm` del tramo | Alto útil (mm) |
|---|---|---|---|---|---|---|
| `alto` (cuello) | 0 | 3 353 | 2 590 | 1 524 | 508 | 2 726 = 4 250 − 1 524 |
| `bajo` | 3 353 | 12 802 | 2 590 | 1 016 | 0 | 3 234 = 4 250 − 1 016 |

- Medidas típicas de una escalonada de 53' (11' arriba y 42' abajo, cama de 60" y 40"), **simuladas**.
- `yMm` se mide desde el piso más bajo: 1 524 − 1 016 = 508.
- maxPayloadKg: 22 000 kg (**simulado**).

## Piezas

| Id | Descripción | L × A × H (mm) | Peso (kg) | Cantidad | Orientaciones | Apilable | Entrega |
|---|---|---|---|---|---|---|---|
| `prensa` | Prensa hidráulica | 4 000 × 2 400 × 3 100 | 15 000 | 1 | "este lado arriba" (LWH, WLH) | no | — |
| `caja` | Caja de refacciones | 1 200 × 1 000 × 1 200 | 500 | 4 | "este lado arriba" (LWH, WLH) | sí | — |

## Resultado esperado

- **Con dos tramos** (sesión 2.1.6): la prensa queda en el tramo `bajo` (3 100 ≤ 3 234) y las 4 cajas caben en el
  tramo `alto` en un solo nivel (2 a lo largo: 2 400 ≤ 3 353; 2 a lo ancho: 2 000 ≤ 2 590), o en el bajo si conviene
  al peso por eje. Todo colocado:
  17 000 kg.
- **Modelada como un solo tramo** con el alto menor (2 726), como obliga hoy el contrato 0.1.0: la prensa queda con
  `EXCEEDS_DIMENSIONS`. Con el alto mayor (3 234) en un solo tramo, el plan sería imposible: una caja podría quedar
  sobre el cuello con más alto del permitido.
- Ninguna pieza puede quedar encima del escalón (cruzando los dos tramos).

## Por qué

- La plataforma escalonada existe justo para cargas altas: cada tramo tiene su propio alto útil (concepto «Plataformas
  con varios tramos» en `conceptos/equipos-y-normas.md`).
- 4 000 de largo de la prensa caben en los 12 802 del tramo bajo; 2 400 ≤ 2 590 de ancho.

## Qué validar con datos reales

- Altura real de las dos camas cargadas y largo de cada tramo en la unidad (alta del equipo con origen `medido`, ADR-0019).
- Si el patio aprovecha el cuello o lo deja libre por el peso en el eje del tractocamión.

## Historial

| Fecha | Estado | Nota |
|---|---|---|
| 2026-11-25 | Simulado | Caso creado para la sesión 2.1.6 (varios tramos). |
