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
| Reefer 20': altura del piso a la línea roja | 2 158 mm | |
| Reefer 40' HC: altura del piso a la línea roja | 2 425 mm | |
| Plataforma 48': altura de la cama al suelo, cargada | 1 478 mm | |
| Plataforma 53': altura de la cama al suelo, cargada | 1 478 mm | |
| Chasis portacontenedor: altura al suelo | 1 200 mm | |
| Flat rack 40': altura del piso (sobre el chasis) | 648 mm | |
| Flat rack 40': sobrealtura que acepta la naviera | 400 mm | |
| Caja seca 48': carga útil | 30 000 kg | |
| Caja seca 53': carga útil | 30 000 kg | |
| Caja refrigerada 48' y 53': carga útil | 24 900 kg | |
| Tractocamión: tara | 8 800 kg | |

## 5. Equipos sin ficha pública (P-013)

- [ ] 40' reefer de altura estándar (8'6"): interior, altura a la línea roja, puerta, tara, carga útil.
- [ ] Plataforma de 40': largo, ancho, altura de la cama, tara.
- [ ] Plataforma extensible: largo cerrada y extendida, ancho, altura de la cama, carga útil en cada largo.
- [ ] Camioneta de 3,5 t: interior, puerta, carga útil.
- [ ] 45' pallet wide y Conestoga: confirmar puerta y ancho útil.
