# Carretera

Cajas, plataformas y camiones del autotransporte mexicano (NOM-012).

| Caso | Título | Subcategoría | Fase | Motivos esperados |
|---|---|---|---|---|
| [CS-011](caja-seca/cs-011-53-dos-niveles-y-tipo-de-camino.md) | Caja seca 53': la misma carga por camino ET y por camino B | Caja seca | 1 (el límite legal lo calcula la API) | `WEIGHT_LIMIT_REACHED` (camino B) |
| [CS-012](caja-seca/cs-012-53-pieza-mas-ancha-que-la-puerta.md) | Caja seca 53': cabe adentro pero no pasa por la puerta | Caja seca | 1 | `EXCEEDS_DOOR` |
| [CS-013](refrigerada/cs-013-48rv-limite-de-peso.md) | Caja refrigerada 48': se acaba el peso antes que el espacio | Caja refrigerada | 1 | `WEIGHT_LIMIT_REACHED` |
| [CS-014](refrigerada/cs-014-53rv-rack-mas-alto-que-la-puerta.md) | Caja refrigerada 53': puerta más baja que el interior | Caja refrigerada | 1 | `EXCEEDS_DOOR` |
| [CS-015](plataforma/cs-015-53fb-tubos-largos.md) | Plataforma 53': atados de tubo de 12 m | Plataforma | 1 | `EXCEEDS_DIMENSIONS` |
| [CS-016](plataforma/cs-016-40fb-sin-altura-de-cama.md) | Plataforma plana 40': la ficha no da la altura de la cama | Plataforma | 2.3 (equipos de la empresa) | — (la entrada no se puede armar) |
| [CS-017](rigidos/cs-017-camioneta-1-5t-paqueteria.md) | Camioneta de 1,5 t: paquetería | Camiones rígidos y camionetas | 1 | `WEIGHT_LIMIT_REACHED` |
| [CS-018](rigidos/cs-018-camioneta-3-5t-linea-blanca.md) | Camioneta de 3,5 t: línea blanca | Camiones rígidos y camionetas | 1 | `EXCEEDS_DOOR` |
| [CS-019](rigidos/cs-019-rabon-tipo-de-camino.md) | Rabón (C2): la ruta baja la carga permitida | Camiones rígidos y camionetas | 1 (el límite legal lo calcula la API) | `WEIGHT_LIMIT_REACHED` (camino D) |
| [CS-020](rigidos/cs-020-torton-cemento-llantas.md) | Tortón (C3): 8 o 10 llantas cambian la carga | Camiones rígidos y camionetas | 1 (el límite legal lo calcula la API) | `WEIGHT_LIMIT_REACHED` |
| [CS-028](sobredimensionada/cs-028-plataforma-carga-ancha-con-permiso.md) | Plataforma 48': intercambiador de 3,2 m de ancho con permiso de la NOM-040 | Carga sobredimensionada | 1 (el permiso lo arma la API) | `EXCEEDS_DIMENSIONS` (sin permiso) |
| [CS-029](plataforma/cs-029-escalonada-dos-tramos.md) | Plataforma escalonada 53': la máquina alta solo cabe en el tramo bajo | Plataforma | 2.1.6 (varios tramos) | `EXCEEDS_DIMENSIONS` (con un solo tramo) |
