# Marítimo

Contenedores ISO que salen o llegan por el puerto.

| Caso | Título | Subcategoría | Fase | Motivos esperados |
|---|---|---|---|---|
| [CS-001](seco/cs-001-20dv-justo-por-espacio.md) | 20' Dry lleno justo por espacio | Seco | 1 | `NO_SPACE_AVAILABLE` |
| [CS-002](seco/cs-002-20dv-justo-por-peso.md) | 20' Dry lleno justo por peso | Seco | 1 | `WEIGHT_LIMIT_REACHED` |
| [CS-004](seco/cs-004-40hc-tarimas-altas.md) | Tarimas de 2,5 m: solo caben en High Cube | Seco | 1 | `EXCEEDS_DIMENSIONS` (en el 40' Dry) |
| [CS-005](seco/cs-005-45pw-europallets.md) | 45' pallet wide: 3 europallets a lo ancho | Seco | 1 | `NO_SPACE_AVAILABLE` |
| [CS-006](refrigerado/cs-006-20rf-tarima-sobre-alto-interior.md) | Reefer 20': tarima más alta que el interior | Refrigerado | 1 | `NO_SPACE_AVAILABLE`, `EXCEEDS_DIMENSIONS` |
| [CS-007](refrigerado/cs-007-40rf-contra-40rh.md) | Reefer 40' estándar contra High Cube | Refrigerado | 1 | `EXCEEDS_DIMENSIONS` (en el 40' RF) |
| [CS-008](especial/cs-008-40fr-sobredimensionada.md) | Flat rack 40': carga más alta que los postes | Especial | 1 | `EXCEEDS_DIMENSIONS` |
| [CS-009](especial/cs-009-20ot-carga-por-techo.md) | Open top 20': máquina que solo entra por el techo | Especial | 1 | `EXCEEDS_DOOR` (por la puerta) |
| [CS-010](especial/cs-010-40pf-carga-ancha.md) | Plataforma ISO 40': carga más ancha que el piso | Especial | 1 | `EXCEEDS_DIMENSIONS` |
