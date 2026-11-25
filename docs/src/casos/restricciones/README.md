# Restricciones

Reglas que no son de medida: apilado, orden de entrega, peso por eje y orientación.

| Caso | Título | Subcategoría | Fase | Motivos esperados |
|---|---|---|---|---|
| [CS-021](apilado/cs-021-fragil-no-debajo.md) | Pantallas frágiles: nada encima | Apilado | 2.1 (`stackable` y `maxLoadOnTopKg`) | — |
| [CS-022](multidestino/cs-022-tres-destinos-lifo.md) | Tres entregas en una caja: lo último que se carga sale primero | Varios destinos | 2.1 (`deliverySequence`) | — |
| [CS-023](peso-por-eje/cs-023-rollos-de-acero-concentrados.md) | Rollos de acero: el peso no puede quedar en un solo extremo | Peso por eje | 2.1 (peso por eje y centro de gravedad) | — |
| [CS-024](orientacion/cs-024-este-lado-arriba.md) | «Este lado arriba»: acostada cabría, parada no | Orientación | 1 | `EXCEEDS_DIMENSIONS` |
