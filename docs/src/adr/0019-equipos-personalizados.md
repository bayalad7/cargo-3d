# ADR-0019: Equipos personalizados por empresa

- **Estado:** Aceptado
- **Fecha:** 2026-11-25

## Contexto

Los equipos estándar (contenedores ISO, cajas de 48' y 53') tienen ficha pública y entran al catálogo canónico
`contracts/equipment.json`. En México, sin embargo, buena parte de la carga se mueve en equipos sin medida estándar:
camionetas de 1,5 y 3,5 t con caja hecha por un carrocero, rabones, tortones y cajas que cada empresa modifica (forro
de madera, rieles, techo reforzado, caja alargada o recortada). Con las medidas de catálogo el plan falla en el patio
(caso CS-026).

El motor ya no depende del tipo de equipo: recibe tramos (`sections`), puerta y `maxPayloadKg` ya resueltos
(ADR-0015 y ADR-0016). Lo que falta es dónde viven esas medidas y cómo se controlan.

## Decisión

Tres niveles, de lo general a lo particular:

1. **Equipo estándar** (`contracts/equipment.json`): ficha pública, con fuente y confianza. Lo mantiene el proyecto.
2. **Plantilla por clase**: camioneta 1,5 t, camioneta 3,5 t, rabón (C2), tortón (C3), caja seca, caja refrigerada,
   plataforma. Define qué datos hay que capturar y sus rangos válidos (ancho ≤ 2 600 mm y alto ≤ 4 250 mm de la
   NOM-012, configuración vehicular y número de llantas para el PBV).
3. **Equipo de la empresa**: parte de un equipo estándar o de una plantilla y guarda sus medidas propias. Cada medida
   lleva su **origen** (`medido`, `tarjeta`, `ficha`, `simulado`), quién y cuándo la capturó. Cambiar una medida crea
   una **versión nueva**: los planes guardan la versión con la que se calcularon.

La API arma la entrada del motor desde el equipo de la empresa: tramos, puerta y `maxPayloadKg` efectivo
(mínimo entre la capacidad del equipo y el PBV de la NOM-012 menos taras, por tipo de camino). El motor no cambia.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Estándar + plantilla + equipo de la empresa (elegida)** | Cubre cualquier caja sin perder el catálogo; historial de alteraciones | Más tablas y pantallas en la fase 2 |
| Solo catálogo estándar | Simple | No sirve para rígidos ni cajas modificadas |
| Medidas libres en cada plan | Sin alta previa | Se recapturan en cada viaje, sin control ni historial |

## Consecuencias

- **Positivas:** cualquier equipo se puede usar; las alteraciones quedan con fecha; el periodo de pruebas puede marcar
  qué medidas ya se midieron.
- **Negativas / costos:** tablas `equipment_templates`, `company_equipment` y `company_equipment_versions` en la
  sesión 2.3 y su pantalla de alta en la 2.4; el alta debe guiar la medición (dónde medir, el punto más angosto).
- **Qué tendríamos que hacer para revertirla:** dejar solo el catálogo estándar y pedir las medidas en cada plan.

Detalle en [Equipos personalizados](../conceptos/equipos-personalizados.md).
