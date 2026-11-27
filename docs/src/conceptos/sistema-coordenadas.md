# Sistema de coordenadas

Todo el proyecto (motor en Go, validador en TypeScript y visor en Three.js) usa **un solo sistema de coordenadas**.
Aquí se explica con dibujos; la decisión y sus alternativas están en el
[ADR-0015](../adr/0015-geometria-del-contrato-del-motor.md) y los campos exactos en `contracts/engine.schema.json`.

## Los tres ejes

Párate **dentro del equipo, pegado a la pared del fondo** (la nariz, el extremo opuesto a la puerta) y **mira hacia
la puerta**. El origen es la esquina de abajo a tu izquierda.

| Eje | Dirección | Medida del equipo |
|---|---|---|
| `x` | hacia el frente de ti, rumbo a la puerta | largo |
| `y` | hacia arriba | alto |
| `z` | hacia tu derecha | ancho |

```
 Vista desde arriba (y sale de la hoja)

            z = 0                 pared a tu izquierda
            ┌──────────────────────────────────────────┐
   fondo    │ o ───────────► x                         ╎  puerta
  (nariz)   │ │                                        ╎  (x máxima)
            │ ▼ z                                      ╎
            └──────────────────────────────────────────┘
            z máxima              pared a tu derecha

 Vista desde tu izquierda (z entra en la hoja)

   y máxima ┌──────────────────────────────────────────┐  techo
            │ ▲ y                                      ╎
            │ │                                        ╎  puerta
            │ o ───────────► x                         ╎
            └──────────────────────────────────────────┘  piso (y = 0)
```

**Cuidado con «izquierda».** Quien carga en el patio suele mirar desde la puerta hacia adentro y ve la pared de
`z = 0` a su **derecha**. En una caja o plataforma ya enganchada, esa pared es la del lado del **copiloto**.
Si la interfaz habla de izquierda o derecha, tiene que decir desde dónde se mira.

**Por qué así:** es un sistema de *mano derecha* con `y` hacia arriba, el mismo de Three.js. El visor no convierte
ejes; solo cambia la escala (1 mm = 0,001 m).

## Dónde está una pieza

Una colocación guarda la **esquina mínima** de la pieza ya girada (`xMm`, `yMm`, `zMm`) y sus medidas sobre cada
eje (`dxMm`, `dyMm`, `dzMm`). La pieza ocupa:

```
de xMm a xMm + dxMm      de yMm a yMm + dyMm      de zMm a zMm + dzMm
```

Dos piezas cuyas caras **se tocan** no se traslapan: la comparación es estricta (`<`), no `<=`. Cómo lo revisa
el motor (encimarse, quedar dentro y área de contacto) está en [Geometría de cajas](geometria.md).

**Ejemplo** (`contracts/examples/maritime/dry/20dv-heavy`): el segundo cajón está en `(0, 0, 1100)` y mide
`1400 × 1000 × 1100` sobre x, y, z. Termina en `z = 2200`, dentro de los 2 352 de ancho, y toca al primero (que
termina en `z = 1100`) sin encimarse.

En Three.js, `BoxGeometry` se dibuja **centrada**, así que el visor suma la mitad de cada medida:

```ts
mesh.position.set(x + dx / 2, y + dy / 2, z + dz / 2); // ya en metros
```

## Las 6 orientaciones

Cada orientación son tres letras: qué medida de la pieza (`L` largo, `W` ancho, `H` alto) queda sobre cada eje, en el
orden **x, z, y**. La última letra es siempre la vertical.

| Orientación | sobre x (`dx`) | sobre z (`dz`) | sobre y (`dy`) | Cómo se ve |
|---|---|---|---|---|
| `LWH` | L | W | H | como viene |
| `WLH` | W | L | H | girada 90° sobre el eje vertical |
| `LHW` | L | H | W | acostada de lado |
| `HLW` | H | L | W | acostada de lado y girada |
| `WHL` | W | H | L | parada de punta |
| `HWL` | H | W | L | parada de punta y girada |

«Este lado arriba» significa que `H` se queda vertical: solo `LWH` y `WLH`.

**Ejemplo:** una caja de L = 1200, W = 800, H = 600 en `WHL` mide `dx = 800`, `dz = 600`, `dy = 1200`.
Las seis juntas están en `contracts/examples/road/dry-van/53dv-contract-coverage`.

> **Trampa:** el orden de las letras es x, **z**, y; el de los campos es `dx`, **`dy`**, `dz`. No son el mismo.

## En el código del motor

| Concepto | Dónde vive en Go |
|---|---|
| Las 6 orientaciones | `Orientation` y sus constantes, en `backend/internal/engine/orientation.go` |
| Medidas de la pieza girada | `CargoItem.RotatedDimensions(o)`, que regresa `Dimensions{DXMM, DYMM, DZMM}` y `ok = false` si la orientación no es una de las seis |
| Orientaciones permitidas | `CargoItem.Orientations()`: las seis solo si `allowedOrientations` no vino |
| Una colocación | `Placement` en `placement.go`; `Placement.Dimensions()` regresa sus `dx`, `dy` y `dz` |

La prueba `TestOutputExamplesDimensionsMatchOrientation` revisa que, en todos los ejemplos de `contracts/examples/`,
las medidas de cada colocación sean las que da `RotatedDimensions`. Si esta tabla y el código dejan de coincidir, esa
prueba falla.

## Los tramos

El espacio de carga es una lista de **tramos** (`sections`), cada uno una caja con su esquina mínima y sus medidas.
En la versión 0.1.0 hay uno solo, casi siempre en `(0, 0, 0)`. Sus coordenadas y las de las piezas usan el **mismo
origen**, el del equipo. Por eso, en un open top cargado por el techo, el tramo empieza en `xMm = 109`,
`zMm = 49`: es la abertura del techo, centrada en el interior.

## Coordenadas negativas

Solo `zMm` puede ser negativa, y solo si el tramo tiene margen a la izquierda (`overhang.leftMm`): la pieza sobresale
por el costado de un flat rack o una plataforma y lleva `outOfGauge: true`
([ADR-0017](../adr/0017-carga-sobredimensionada-y-open-top.md)).
