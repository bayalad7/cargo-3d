# Equipos personalizados

Cómo maneja Cargo 3D los equipos que no tienen ficha estándar o que cada empresa modificó (ADR-0019).

## Los tres niveles

```
Equipo estándar  ──►  Plantilla por clase  ──►  Equipo de la empresa  ──►  Versión 1, 2, 3…
(contracts/            (qué se captura y          (medidas propias,          (cada alteración
 equipment.json)        rangos válidos)             origen de cada una)        con fecha)
```

| Nivel | Ejemplo | Quién lo mantiene |
|---|---|---|
| Estándar | `53dry`: caja seca 53' Fruehauf, interior 16 003 × 2 540 × 2 794 | El proyecto, con fuente y confianza |
| Plantilla | «Tortón C3»: pide interior, puerta, tara, PBV de tarjeta y número de llantas | El proyecto |
| Equipo de la empresa | «Tortón 07»: 9 480 × 2 430 × 2 590, 10 llantas, tara 9 620 kg | Cada empresa |

## Plantillas iniciales

| Plantilla | Configuración NOM-012 | Qué se captura |
|---|---|---|
| Camioneta 1,5 t | — (vehículo ligero) | Interior, puerta o cortina, carga útil de la tarjeta |
| Camioneta 3,5 t | — (vehículo ligero) | Interior, puerta o cortina, carga útil de la tarjeta |
| Rabón | C2 | Interior, puerta, tara, PBV de la tarjeta |
| Tortón | C3 | Interior, puerta, tara, PBV, **número de llantas (8 o 10)** |
| Caja seca / refrigerada | T3-S2, T3-S3 | Interior, puerta, tara de la caja; equipo de frío en refrigeradas |
| Plataforma | T3-S2, T3-S3 | Largo y ancho de la cama, **altura de la cama**, tara |
| Varios tramos | Escalonada, doble caída | Un tramo por piso (`sections[]`, sesión 2.1.6) |

Las medidas simuladas de los casos [CS-017](../casos/carretera/rigidos/cs-017-camioneta-1-5t-paqueteria.md) a
[CS-020](../casos/carretera/rigidos/cs-020-torton-cemento-llantas.md) sirven de ejemplo para estas plantillas.

## Reglas de validación al capturar

- Ancho ≤ 2 600 mm y alto total ≤ 4 250 mm (NOM-012); si el interior se acerca, se pide confirmar.
- El interior no puede ser mayor que el exterior ni la puerta mayor que el interior por más de un margen razonable
  (la puerta sí puede ser más ancha que el interior en algunas cajas).
- Medir **entre los puntos más angostos**: rieles, refuerzos, bisagras, forro y equipo de frío (ADR-0016: el equipo
  se redondea hacia abajo).
- Cada medida guarda su origen: `medido` (con fecha y quién), `tarjeta`, `ficha` o `simulado`. La app marca con el
  estado `warning` (texto + ícono, nunca solo color) las que no son `medido`.

## Alteraciones

Si una empresa forra, alarga o refuerza una caja, no se edita la medida: se crea una **versión nueva** del equipo con
la fecha y el motivo. Los planes guardan la versión con la que se calcularon, así un plan viejo sigue explicando por
qué se acomodó como se acomodó (caso [CS-026](../casos/operacion/equipo-alterado/cs-026-caja-forrada-por-dentro.md)).

## De la base de datos al motor

```
company_equipment_versions  ──API──►  EngineInput.container
  medidas, puerta, tara,               sections[], door,
  capacidad, llantas                   maxPayloadKg = mín(capacidad, PBV(config, camino) − taras)
```

El motor no sabe si el equipo es estándar o de la empresa: solo recibe números (ADR-0015).

## Cuándo se construye

| Sesión | Qué |
|---|---|
| 2.3 | Tablas `equipment_templates`, `company_equipment` y `company_equipment_versions`; carga del catálogo estándar |
| 2.4 | Pantalla de alta y edición con guía de medición; selector de equipo al crear un plan |
| 3.2 | Captura de medidas desde la app móvil, con foto |
