# Equipos y normas

Qué equipos acomoda Cargo 3D, de dónde salen sus límites y quién los aplica. El motor **no conoce ninguna norma**:
recibe medidas y un `maxPayloadKg` ya resueltos ([ADR-0016](../adr/0016-que-rechaza-el-motor.md)). Lo de este
documento lo usan el catálogo y la API.

Niveles de confianza: **oficial** (leído en la norma), **fabricante** (ficha publicada), **secundaria** (otra fuente)
y **supuesto** (sin fuente). Consultado el 2026-11-21.

## Contenedores marítimos

| Norma | Qué define |
|---|---|
| ISO 668 | Medidas exteriores y pesos brutos máximos |
| ISO 1496 | Pruebas y mínimos interiores y de puerta |
| ISO 6346 | Número del contenedor y código de tamaño y tipo (`22G1`, `45R1`…) |

(Los textos de ISO son de pago y no se leyeron; lo de aquí viene de fuentes secundarias y fichas.)

- **ST, DV y GP son el mismo contenedor seco estándar:** `22G1` en 20' y `42G1` en 40'. El `45G1` es el 40' high cube.
- **45':** la edición 2022 de ISO 6346 cambió el código de largo de `L` a `5` (`L5G1` → `55G1`, secundaria). El catálogo
  guarda los dos.
- **Peso:** peso bruto = tara + carga útil. Los brutos que más se ven son 24 000, 30 480, 32 500 y 34 000 kg
  (45 000 en flat rack de 40'). Es la forma rápida de revisar un catálogo: si `tara + carga útil` no da un bruto
  conocido, algo está mal capturado.
- **La puerta es más chica que el interior** en los secos (20' DV: 2 393 de alto interior, 2 280 de puerta), por eso
  existe el motivo `EXCEEDS_DOOR`.
- **Reefer:** la altura útil llega a la **línea roja de carga**, no al techo: arriba tiene que circular el aire. Por
  eso su puerta puede ser más alta que su altura útil.
- **Open top:** se carga por la puerta o con grúa por el techo; son dos entradas del catálogo
  ([ADR-0017](../adr/0017-carga-sobredimensionada-y-open-top.md)).
- **Flat rack y plataforma:** admiten carga **OOG** (*out of gauge*), que rebasa los postes o el piso.

### VGM (SOLAS VI/2)

Todo contenedor lleno que se embarca necesita su peso bruto verificado. Método 1: pesar el contenedor cargado.
Método 2: sumar todo. Para el método 2 sirven `tareKg` y el peso del plan:

```
VGM = tara + piezas + pallets, calzas y amarres
```

El plan **no** incluye calzas ni amarres, así que por sí solo no es un VGM.

## Autotransporte: NOM-012-SCT-2-2017

[NOM-012-SCT-2-2017](https://dof.gob.mx/nota_detalle.php?codigo=5508944&fecha=26/12/2017&print=true), DOF del
26/12/2017, en vigor desde el 24/02/2018. Estado
[«Vigente» en PLATIICA](https://platiica.economia.gob.mx/normalizacion/nom-012-sct-2-2017/).

### Medidas (oficial)

- Ancho máximo 2 600 mm (§6.2.1.1) y alto máximo 4 250 mm con todo y carga (§6.2.1.2).
- Altura útil de una plataforma = 4 250 − altura de la cama.
- Margen lateral por carretera = (2 600 − ancho del tramo) / 2. En un flat rack de 2 208 mm son 196 mm por lado.

### Peso por eje, Tabla A-2, en toneladas (oficial)

Leída de la imagen publicada en el DOF. `C` camión, `T` tractocamión, `S` semirremolque, `R` remolque.

| Eje o conjunto | Configuración | ET y A | B | C | D |
|---|---|---|---|---|---|
| Sencillo, 2 llantas | todas | 6,5 | 6,0 | 5,5 | 5,0 |
| Sencillo, 4 llantas | C-R y T-S-R | 10,0 | 9,5 | 8,0 | 7,0 |
| | C y T-S | 11,0 | 9,5 | 8,0 | 7,0 |
| Motriz sencillo, 4 llantas | C-R y T-S-R | 11,0 | 10,5 | 9,0 | 8,0 |
| | C y T-S | 12,5 | 10,5 | 9,0 | 8,0 |
| Motriz tándem, 6 llantas | C-R y T-S-R | 15,0 | 13,0 | 11,5 | 11,0 |
| | C y T-S | 17,5 | 13,0 | 11,5 | 11,0 |
| Tándem, 8 llantas | C-R y T-S-R | 17,0 | 15,0 | 13,5 | 12,0 |
| | C y T-S | 19,0 | 15,0 | 13,5 | 12,0 |
| Motriz tándem, 8 llantas | C-R y T-S-R | 18,0 | 17,0 | 14,5 | 13,5 |
| | C y T-S | 21,0 | 17,0 | 14,5 | 13,5 |
| Trídem, 12 llantas | C-R y T-S-R | 23,5 | 22,5 | 20,0 | no autorizado |
| | C y T-S | 26,5 | 22,5 | 20,0 | no autorizado |

La diferencia entre configuraciones solo existe en caminos ET y A. La Tabla A-1 es la equivalente para autobuses.

### Peso bruto vehicular, Tabla B-2, en toneladas (oficial)

| Configuración | ET y A | B | C | D |
|---|---|---|---|---|
| C2 (rabón) | 19,0 | 16,5 | 14,5 | 13,0 |
| C3, 8 llantas | 24,0 | 19,0 | 17,0 | 16,0 |
| C3, 10 llantas (tortón típico) | 27,5 | 23,0 | 20,0 | 18,5 |
| T3-S2 | 46,5 | 38,0 | 33,5 | no autorizado |
| T3-S3 | 54,0 | 45,5 | 40,0 | no autorizado |
| T3-S2-R4 (full) | 66,5 | no autorizado | no autorizado | no autorizado |

**Comprobación cruzada** (ET y A): C2 = 6,5 + 12,5 = 19,0; tortón = 6,5 + 21,0 = 27,5; T3-S2 = 6,5 + 21,0 + 19,0 = 46,5.
En las combinaciones más largas el bruto es menor que la suma de los ejes, así que **se revisan los dos límites**.

### De la norma a `maxPayloadKg`

```
maxPayloadKg = mínimo entre
    la carga útil del equipo
    el peso bruto vehicular (configuración, camino) − tara del vehículo − tara del equipo
```

**Ejemplo:** caja de 53' (6 810 kg, fabricante) con tractocamión de 8 800 kg (supuesto) en T3-S2 por camino ET o A:
46 500 − 8 800 − 6 810 = 30 890 kg. Por camino B: 38 000 − 8 800 − 6 810 = 22 390 kg.

### Cambios a la vista

Al 2026-11-21 no se encontró ningún proyecto de modificación publicado. Una
[nota de septiembre de 2026](https://thelogisticsworld.com/transporte/camion-torton-capacidad-carga-medidsa-principales-usos-ventajas-mexico-sict/)
dice que el programa 2026 la contempla, y el 19/11/2026 salió un aviso de autorregulación que
[según T21](https://t21.us/self-regulation-of-road-transport-these-are-the-requirements-set-by-the-sict/) no cambia
pesos ni dimensiones. **Se vuelve a revisar el DOF en la sesión 2.1.5**, antes de programar los límites.

## Plataformas con varios tramos

La escalonada, la doble caída, la cama baja y el RGN tienen pisos a distinta altura: cada tramo lleva su propia altura
útil. El contrato 0.1.0 solo acepta un tramo; los demás llegan en la sesión 2.1.6.

## Pallets

Son piezas, no equipos. EPAL 1: 1 200 × 800 mm. GMA: 48" × 40" (1 219 × 1 016 mm).
