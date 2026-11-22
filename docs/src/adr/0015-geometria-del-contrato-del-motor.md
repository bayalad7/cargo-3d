# ADR-0015: Geometría del contrato del motor (coordenadas, orientaciones y tramos)

- **Estado:** Aceptado
- **Fecha:** 2026-11-21

## Contexto

`contracts/engine.schema.json` (sesión 1.1.1) necesita fijar tres cosas que usan el motor en Go, el validador en
TypeScript y el visor en React Three Fiber, y que serían muy caras de cambiar después:

1. **El sistema de coordenadas.** El motor trabaja en milímetros enteros; el visor usa Three.js, donde `y` apunta
   hacia arriba y el sistema es de mano derecha. La frase «esquina del fondo, abajo, a la izquierda» es ambigua:
   la izquierda cambia según se mire desde el fondo o desde la puerta, y una de las dos lecturas da un sistema de
   mano izquierda.
2. **Cómo se dice que una pieza puede girar.** Hay piezas que no se pueden voltear («este lado arriba») y otras
   que solo pueden girar sobre el eje vertical.
3. **La forma del espacio de carga.** Un contenedor o una caja seca son una sola caja, pero las plataformas
   escalonadas, de doble caída, de cama baja y RGN tienen tramos con distinta altura de piso
   (investigación de la sesión 1.1.1 y la guía de equipos de plataforma de TRAFFIX).

## Decisión

1. **Coordenadas:** origen en la esquina interior del fondo (nariz), abajo, a la izquierda **para alguien parado en
   la pared del fondo mirando hacia la puerta**. `x` = largo, del fondo hacia la puerta; `y` = alto, hacia arriba;
   `z` = ancho, hacia la derecha de esa persona. Es de mano derecha y coincide con Three.js: el visor solo escala
   (1 mm = 0,001 m) y suma la mitad de cada medida, porque `BoxGeometry` está centrada.

   ```
           y (arriba)
           │
    fondo  │
    (nariz)o──────────► x  (hacia la puerta)
          ╱
         z  (ancho, a la derecha de quien mira hacia la puerta)
   ```

2. **Una colocación guarda la esquina mínima** de la pieza ya girada (`xMm`, `yMm`, `zMm`) y sus medidas sobre
   cada eje (`dxMm`, `dyMm`, `dzMm`). Las medidas son redundantes con la orientación a propósito: el visor no
   recalcula y el validador revisa que coincidan.
3. **Orientaciones:** un solo arreglo `allowedOrientations` con las 6 permutaciones de (L, W, H). Las letras dicen
   qué medida de la pieza queda sobre cada eje, en el orden **x, z, y**: `LWH` = como viene, `WLH` = girada 90° sobre
   el eje vertical, `LHW`/`HLW` = acostada de lado, `WHL`/`HWL` = parada de punta. Si el campo no viene, se permiten
   las 6. No hay un booleano `keepUpright` aparte, para no tener dos fuentes de verdad.
4. **El equipo es una lista de tramos** (`sections`), cada uno una caja alineada a los ejes con su esquina mínima y
   sus medidas útiles. En la versión 0.1.0 del contrato se acepta un solo tramo (`maxItems: 1`); subir ese límite
   cuando el motor sepa acomodar en varios es un cambio compatible.
5. **La puerta** es opcional y está en el extremo de `x` máxima. Una pieza entra a lo largo de `x`, así que pasa si
   alguna orientación permitida cumple `dz ≤ ancho de la puerta` y `dy ≤ alto de la puerta`.
6. **La categoría del equipo no entra al contrato.** La jerarquía (marítimo / autotransporte / …) vive en el
   catálogo; al motor solo le importan la geometría y los límites.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Izquierda vista desde el fondo, y arriba (elegida)** | Mano derecha; igual que Three.js, sin convertir ejes | Hay que explicar desde dónde se mira |
| Izquierda vista desde la puerta | Es como ve el equipo quien carga en el patio | Mano izquierda: el visor tiene que invertir `z` y es fácil equivocarse |
| `z` hacia arriba (convención de CAD y de muchos artículos de bin packing) | Común en la literatura | Habría que rotar todo en el visor |
| Guardar el centro de la pieza | Igual que Three.js | Medias de milímetro con medidas impares; deja de ser entero |
| Booleanos `keepUpright` / `rotatable` | Fáciles de capturar | No expresan casos intermedios y pueden contradecirse con la lista |
| Equipo como una sola caja | Más simple | No representa las plataformas escalonadas; cambiarlo después rompería el contrato |

## Consecuencias

- **Positivas:** el motor, el validador y el visor comparten la misma geometría sin conversiones; las plataformas
  escalonadas caben en el contrato sin romperlo; la interfaz traduce «este lado arriba» a `["LWH","WLH"]`.
- **Negativas / costos:** quien lee las coordenadas tiene que recordar desde dónde se mira; `dxMm`/`dyMm`/`dzMm`
  ocupan espacio de más en la salida.
- **Qué tendríamos que hacer para revertirla:** cambiar `engine.schema.json` a una versión mayor, el motor, el
  validador, el visor y todos los casos de `engine/fixtures/`.
