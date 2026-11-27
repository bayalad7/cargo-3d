# ADR-0022: La guía y los comentarios explican Go solo con Go

- **Estado:** Aceptado
- **Fecha:** 2026-11-26

## Contexto

El [ADR-0009](0009-idioma-y-documentacion.md) pedía una guía de aprendizaje de Go comparada con otro lenguaje, y la
regla 6 de `CLAUDE.md` repetía la idea. La revisión de la sesión 1.1.2 mostró el costo. Varias comparaciones
enseñaban algo falso de Go:

- «un receptor por valor no puede modificar el original»: sí puede, a través de un slice o un puntero;
- «un struct es como una clase»: un struct se copia al asignarlo;
- equivalencias de herramientas y de imports que no lo eran.

Otras forzaban un parecido que no existe. Benjamín decidió que el proyecto se escriba y se explique solo en Go, sin
traer ideas de otros lenguajes.

## Decisión

El código, sus comentarios, la guía de aprendizaje, las skills y los prompts explican Go **solo en términos de Go**:
- no se compara con otros lenguajes;
- cada concepto se explica con el código real del proyecto y con el comportamiento que Go tiene de verdad;
- si hay duda, el comportamiento se comprueba con un programa o una prueba antes de escribirlo;
- el código se escribe como Go idiomático y actual (por ejemplo, `new(int64(400))` desde Go 1.26 y `slices.Clone`).

Sustituye la parte de la comparación del punto 3 del ADR-0009.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Solo Go (elegida)** | Nada que desaprender; las explicaciones describen lo que el compilador hace | Al principio cuesta más relacionar un concepto nuevo con algo conocido |
| Comparar con otro lenguaje diciendo dónde se rompe la comparación | Da un punto de partida conocido | Sigue sembrando el modelo mental de otro lenguaje; cada comparación es un lugar más donde equivocarse |
| Comparar sin matices (lo de antes) | Rápido de leer | Enseñó cosas falsas de Go |

## Consecuencias

- **Positivas:**
  - Lo que se aprende en la guía es lo que pasa al compilar y correr.
  - Las trampas de Go (copias superficiales, `nil`, valores cero, desbordamiento, `defer` con `os.Exit`) se explican
    por sí mismas.
- **Negativas / costos:**
  - Se reescribieron los capítulos 00, 01, 01b y 02, los comentarios del backend, la skill `/explicar` y los prompts que
    pedían comparar.
  - Los registros históricos (bitácoras anteriores a esta decisión y el ADR-0002, que eligió el lenguaje del
    backend) se conservan como estaban en su fecha.
- **Qué tendríamos que hacer para revertirla:** volver a escribir la regla 6 de `CLAUDE.md` y el punto 3 del ADR-0009, y
  agregar las comparaciones capítulo por capítulo.
