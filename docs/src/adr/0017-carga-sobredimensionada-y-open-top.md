# ADR-0017: Carga sobredimensionada (OOG) y open top

- **Estado:** Aceptado
- **Fecha:** 2026-11-21

## Contexto

Al agregar ejemplos para los contenedores especiales (investigación de equipos faltantes, 2026-11-21) aparecieron dos
casos que el contrato del motor no representaba:

1. **Flat rack y plataformas ISO.** No tienen paredes laterales ni techo, y en la práctica se embarcan piezas más
   altas que los postes o más anchas que el piso como **carga sobredimensionada** (OOG, *out of gauge*), con recargo
   de la naviera o permiso en carretera. Con el contrato 0.1.0 una pieza así salía `EXCEEDS_DIMENSIONS`, aunque sí
   viaja.
2. **Open top.** Se puede cargar por la puerta, como un contenedor seco, o con grúa por el techo. Por el techo la
   puerta no limita, pero la abertura del techo es un poco más chica que el interior (ONE: 5 680 × 2 252 mm en el
   20' contra 5 898 × 2 350 de interior).

## Decisión

1. **Margen de sobredimensión por tramo.** `Section.overhang` (opcional) con `topMm`, `leftMm` y `rightMm`: cuánto
   puede rebasar una pieza el tramo por arriba y a cada lado. **No hay sobrelargo** (frente ni atrás): en el buque
   bloquea a los contenedores vecinos y en carretera cambia el largo total del vehículo.
2. **El margen lo fija quien llama al motor** (API o catálogo), igual que el peso efectivo (ADR-0016): ahí se aplican
   los límites legales (NOM-012: 4 250 mm de alto total y 2 600 mm de ancho) y las reglas de la naviera. El motor
   solo respeta los números que recibe.
3. **El plan avisa.** Cada colocación que rebasa el tramo nominal lleva `outOfGauge: true`, y `metrics.outOfGaugeCount`
   da el total, para que la interfaz advierta del recargo o del permiso. La ocupación de volumen se sigue calculando
   contra el volumen nominal del tramo.
4. **Con margen a la izquierda, `zMm` puede ser negativo** (hasta `−leftMm`). El esquema solo pone un mínimo general
   de −1 000 mm; el límite exacto de cada plan lo revisa el validador independiente (`packages/validator`, sprint 1.4).
5. **El open top son dos entradas del catálogo**, sin cambiar el contrato: «por puerta» (con `door` y el interior
   completo) y «por techo» (sin `door`, con el largo y el ancho de la abertura del techo). Quien planea elige cuál usar.
   En la entrada «por techo» el tramo **se coloca centrado** en el interior con `xMm` y `zMm` (en el 20' de ONE:
   `xMm = (5 898 − 5 680) / 2 = 109`, `zMm = (2 350 − 2 252) / 2 = 49`), para que las dos entradas usen el mismo
   marco de referencia y el visor dibuje las piezas respecto a las paredes reales.
6. **El margen es solo para tramos abiertos.** Un equipo con `door` (contenedor cerrado o caja) no puede traer
   `overhang` con `topMm`, `leftMm` o `rightMm` mayor que cero: sería cargar «a través» del techo o de las paredes.
   La validación de la entrada (sesión 1.1.4) lo rechaza como error.
7. **La ocupación de volumen puede pasar de 1** cuando hay piezas sobredimensionadas, porque se mide contra el volumen
   nominal del tramo; el contrato no le pone tope.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Margen por tramo con marca por pieza (elegida)** | Expresa flat rack y plataformas; cada pieza OOG queda identificada | Coordenadas negativas en z; un campo más en cada tramo |
| Un solo `overheightAllowanceMm` por equipo | Más simple | No cubre el sobreancho, que es el OOG más común en flat rack |
| Rechazar todo lo que rebase (`EXCEEDS_DIMENSIONS`) | Sin cambios | El plan dice «no cabe» de algo que sí se embarca |
| Margen también al frente y atrás | Cubre todo | Rara vez se permite; complica el motor sin un caso real |
| **Open top como dos entradas del catálogo (elegida)** | Sin cambio de contrato | El usuario tiene que elegir la forma de carga |
| Campo `loadAccess: door \| roof` en el contrato | Explícito | Más contrato para un caso que el catálogo ya resuelve; se puede agregar después |

## Consecuencias

- **Positivas:** flat rack, plataformas ISO y open top tienen ejemplos válidos; la carga sobredimensionada se ve en el
  plan en vez de desaparecer como «no cabe».
- **Negativas / costos:** el motor y el validador deben aceptar `zMm` negativo y calcular `outOfGauge`; quien llama
  tiene que fijar márgenes correctos (un margen mal puesto produce un plan ilegal en carretera).
- **Qué tendríamos que hacer para revertirla:** quitar `overhang`, `outOfGauge` y `outOfGaugeCount` del contrato (cambio
  incompatible), y los ejemplos `maritime/special/40fr-oog`.
