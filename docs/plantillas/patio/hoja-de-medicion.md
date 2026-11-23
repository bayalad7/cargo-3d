# Hoja de medición para el patio

Copia esta carpeta a `docs/datos/patio/` (que no se sube a git) y llénala **ahí**. Estos originales van siempre vacíos.

**Cómo medir:** en milímetros, entre los puntos más angostos (rieles, refuerzos, bisagras). Si dudas, anota la medida
más chica en equipos y la más grande en piezas (ADR-0016). Sin placas, nombres de clientes ni números de pedido.

## 1. Cargas reales (P-001): mínimo 3

- [ ] `cargas.csv`: una fila por carga.
- [ ] `cargas-piezas.csv`: una fila por tipo de pieza.
- [ ] Incluye al menos una carga que quedó justa y una donde algo **no** cupo.

## 2. Reglas del patio (P-001)

- [ ] ¿Qué nunca se pone abajo? ¿Qué nunca se encima?
- [ ] ¿Cuánto espacio se deja para montacargas, calzas y amarres?
- [ ] ¿Cómo se carga cada equipo: montacargas, rampa, grúa, por el costado?
- [ ] ¿En qué orden se carga cuando hay varios destinos?
- [ ] ¿Qué tipos de camino usan (ET, A, B, C, D) y con qué configuración (T3-S2, T3-S3…)?

## 3. Rabón y tortón (P-003): 1 o 2 unidades de cada uno

- [ ] Largo, ancho y alto interior de la caja; ancho y alto de la puerta.
- [ ] Tara y peso bruto vehicular de la tarjeta de circulación o de la placa del fabricante.
- [ ] Tortón: número de llantas del tándem (6 u 8).
- [ ] Distancia entre ejes, y del eje delantero al inicio de la caja (para el sprint 2.1).

## 4. Valores por confirmar (P-004)

| Dato | Hoy se usa | Medido |
|---|---|---|
| Reefer 20': altura del piso a la línea roja | 2 159 mm (alto interior Hapag; la línea roja de esa serie no se publica) | Aproximadamente 2,150 mm a 2,191 mm (2.15 m – 2.19 m / ~7' 1" – 7' 2"): Hapag-Lloyd: 2,159 mm (2.16 m), Maersk (referencia estándar 20'): 2,167 mm (2.17 m), |
| Reefer 40' HC: altura del piso a la línea roja | 2 425 mm (alto interior Hapag; Maersk publica 2 415 en otra serie) | Aproximadamente 2,380 mm a 2,425 mm (2.38 m – 2.43 m / ~7' 10" – 7' 11"): Hapag-Lloyd: 2,425 mm (2.43 m), Maersk (40' HC estándar): 2,415 mm (2.42 m), Maersk (unidades especiales como Starfresh / Super Freezer): entre 2,380 mm y 2,425 mm, |
| Plataforma 48': altura de la cama al suelo, cargada | 1 478 mm | Entre 1.40 m y 1.52 m (55" a 60" / ~4' 7" a 5' 0") |
| Plataforma 53': altura de la cama al suelo, cargada | 1 478 mm | Entre 1.40 m y 1.52 m (55" a 60" / ~4' 7" a 5' 0") |
| Chasis portacontenedor: altura al suelo | 1 219 mm (Triton y Conway) | Entre 1.12 m y 1.24 m (44" a 49") |
| Flat rack 40': altura del piso (sobre el chasis) | sin cifra pública (648 mm es la altura plegada, Textainer) | Entre 1.55 m y 1.72 m (5' 1" a 5' 8") |
| Plataforma ISO 40': alto libre sobre la plataforma montada en chasis | 2 383 mm, derivado (4 250 − 648 − 1 219) | Estándar (bajo el límite legal de 4,250 mm de la NOM-012-SCT-2 en México) es de 2,383 mm (2.38 m) |
| Flat rack 40': sobrealtura que acepta la naviera | 400 mm (límite interno; CMA CGM admite en principio 1 700) | En la práctica operativa marítima, las navieras suelen aceptar sin mayor complejidad cargas de hasta ~4.00 – 4.50 m de altura total (es decir, ~2.00 a 2.50 m de sobrealtura por encima de los testeros del equipo), siempre sujetas a aprobación previa caso por caso (OOG approval). |
| Caja seca 48': carga útil legal (tarjeta de circulación o NOM-012) | 30 000 kg es capacidad estructural Fruehauf, no útil legal | En México, la carga útil legal de una caja seca de 48 pies se ubica normalmente entre 26,000 kg y 30,500 kg (26 a 30.5 toneladas), dependiendo de si el semirremolque tiene 2 o 3 ejes (configuración T3-S2 o T3-S3), la tara del equipo y el tipo de camino. |
| Caja seca 53': carga útil legal (tarjeta de circulación o NOM-012) | 30 000 kg es capacidad estructural Fruehauf, no útil legal | En México, la carga útil legal de una caja seca de 53 pies se ubica normalmente entre 25,000 kg y 29,500 kg (25 a 29.5 toneladas) en su configuración estándar de 2 ejes (T3-S2), y puede llegar hasta ~35,000–36,000 kg si cuenta con 3 ejes (T3-S3) (aunque es menos común en 53'). |
| Caja refrigerada 48': carga útil | 25 000 kg (máximo Alzper, estructural) | En México, la carga útil legal de una caja refrigerada (thermo) de 48 pies se ubica normalmente entre 23,500 kg y 26,500 kg (23.5 a 26.5 toneladas) en su configuración estándar de 2 ejes (T3-S2). |
| Caja refrigerada 53': carga útil | 30 000 kg (máximo Alzper, estructural) | En México, la carga útil operativa de una caja refrigerada de 53 pies (configuración estándar T3-S2, 2 ejes en remolque) se ubica entre 22,500 kg y 25,500 kg (22.5 a 25.5 toneladas), con un máximo teórico legal de ~27,500 – 29,000 kg bajo la NOM-012-SCT-2-2017. |
| Tractocamión: tara | 8 800 kg | En México, la tara (peso vehicular en vacío) de un tractocamión quinta rueda estándar de 3 ejes (T3) se ubica normalmente entre 7,800 kg y 9,500 kg (7.8 a 9.5 toneladas), siendo 8,500 kg (8.5 t) el valor de referencia estándar para cálculos de carga útil bajo la NOM-012. |

## 5. Equipos sin ficha pública (P-013)

- [ ] 40' reefer de altura estándar (8'6"): altura a la línea roja y puerta (CIMC ya publica interior, tara y carga útil).
- [ ] Plataforma de 40': altura de la cama al suelo (Fruehauf PM-HT2-40W ya publica largo, ancho, tara y capacidad); confirmar el modelo de la unidad.
- [ ] Caja refrigerada 48': alto interior (Alzper publica 3 790 mm, imposible con su alto total de 4 200 mm).
- [ ] Plataforma extensible: largo cerrada y extendida, ancho, altura de la cama, carga útil en cada largo.
- [ ] Camioneta de 3,5 t: interior, puerta, carga útil.
- [ ] 45' pallet wide y Conestoga: confirmar puerta y ancho útil.
