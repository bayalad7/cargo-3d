# ADR-0016: Qué rechaza el motor y qué reporta (piezas no colocadas, peso efectivo y redondeo)

- **Estado:** Aceptado
- **Fecha:** 2026-11-21

## Contexto

Al diseñar `contracts/engine.schema.json` (sesión 1.1.1) aparecieron tres preguntas sobre los límites del motor:

1. Si una pieza no cabe en el equipo en ninguna orientación, ¿es un **error de la entrada** (se rechaza toda la
   petición, como decía el prompt de la sesión 1.1.4) o un **resultado** (se acomoda lo demás y se avisa)?
2. El peso máximo real depende de la carga útil del equipo **y** de la NOM-012-SCT-2-2017 según la configuración
   vehicular (T3-S2, T3-S3…) y el tipo de camino (ET, A, B, C, D). Por ejemplo, un T3-S2 baja de 46,5 t de peso
   bruto vehicular en camino ET/A a 33,5 t en camino C. ¿Quién calcula ese límite?
3. Las fichas de los equipos vienen en pulgadas y libras, con decimales, y las fuentes no coinciden entre sí
   (hasta 2 t de diferencia en carga útil). El contrato solo acepta milímetros y kilogramos enteros.

## Decisión

1. **Lo que no cabe es un resultado, no un error.** La pieza sale en `unplaced` con un motivo de `UnplacedReason`.
   Se aplica el primero que se cumpla, en este orden: `EXCEEDS_DIMENSIONS`, `EXCEEDS_DOOR`, `EXCEEDS_MAX_WEIGHT`,
   `WEIGHT_LIMIT_REACHED`, `NO_SPACE_AVAILABLE`, `TIME_LIMIT_REACHED`. Los errores de validación quedan solo para
   entradas mal formadas (medidas en cero, ids repetidos, versión distinta, tramos encimados…). Una puerta más
   grande que el interior **no** es error: en el reefer la puerta es más alta que la altura útil bajo la línea de
   carga, y en ese caso la puerta simplemente no limita.
2. **El motor recibe un solo `maxPayloadKg` efectivo**, ya calculado fuera del motor (API o catálogo): el menor entre
   la carga útil del equipo y lo que permite la NOM-012 para ese viaje. El motor no conoce la NOM.
3. **Redondeo conservador:** al capturar, las medidas interiores, las puertas y las capacidades del equipo se
   redondean **hacia abajo**; las medidas y los pesos de las piezas, **hacia arriba**. Un peso de menos de 1 kg
   se vuelve 1 kg. Así el redondeo nunca deja pasar algo que en realidad no cabe.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Pieza no colocada con motivo (elegida)** | Una pieza mal medida no tumba todo el plan; quien planea ve qué falló y por qué | El motor tiene que distinguir motivos y su orden |
| Error de validación | Avisa antes de calcular | Se pierde el plan completo por una sola pieza |
| **Peso efectivo calculado fuera (elegida)** | El motor no depende de una norma que cambia; la tabla de la NOM se prueba aparte | Quien llama al motor tiene que calcularlo bien |
| El motor calcula el límite con la NOM | Un solo lugar | Mete reglas legales al motor; cambia con cada reforma |
| Pesos en gramos (`weightG`) | Exacto para paquetería ligera | Rompe la regla de kilogramos enteros de `CLAUDE.md` y pide cambiar todo el backend |
| **kg enteros redondeados hacia arriba (elegida)** | Simple y seguro | Sobrestima la paquetería ligera: 1 000 sobres de 0,4 kg cuentan como 1 000 kg |

## Consecuencias

- **Positivas:** planes útiles aunque haya piezas imposibles; el motor sigue sin saber nada de leyes ni de unidades
  inglesas; ningún redondeo produce un plan imposible de cargar.
- **Negativas / costos:** la sesión 1.1.4 ya no trata «no cabe» como error de validación; el cálculo del peso
  efectivo por NOM-012 hay que programarlo y probarlo aparte (fase 2); con mucha paquetería ligera el motor puede
  rechazar por peso antes de tiempo.
- **Qué tendríamos que hacer para revertirla:** volver a agregar el error en la validación de la 1.1.4, o cambiar
  la unidad de peso en el contrato (versión mayor) y en todo el backend.
