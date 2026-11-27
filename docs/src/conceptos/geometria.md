# Geometría de cajas

Todo lo que el motor decide («¿cabe aquí?», «¿choca con otra pieza?», «¿en qué se apoya?») se reduce a tres preguntas
sobre **cajas alineadas a los ejes** del equipo. Este capítulo las explica con dibujos. El código está en
`backend/internal/engine/box.go` y sus casos límite en `box_test.go`. Los ejes son los del
[sistema de coordenadas](sistema-coordenadas.md): `x` hacia la puerta, `y` hacia arriba y `z` hacia la derecha.

## La caja

Una `Box` guarda su **esquina mínima** (`XMM`, `YMM`, `ZMM`) y sus **medidas** sobre cada eje (`DXMM`, `DYMM`,
`DZMM`), igual que una colocación del plan (`Placement.Box()` convierte una en otra). «Alineada a los ejes» quiere
decir que sus caras son paralelas a las paredes del equipo: una pieza puede ir acostada, de lado o de punta (las seis
orientaciones), pero nunca en diagonal.

Todo es entero, en milímetros. Con enteros, «se tocan» es una igualdad exacta (`1000 == 1000`) y no hay errores de
redondeo que hagan parecer encimadas a dos cajas pegadas.

## Rangos semiabiertos: el inicio es suyo, el final no

Sobre cada eje, una caja ocupa el rango **semiabierto** `[inicio, inicio + medida)`: el milímetro del inicio es suyo
y el del final ya no. Una caja en `x = 0` de 1 000 mm ocupa del 0 al 999; el 1 000 ya está libre para la siguiente.

```
 eje x (mm)    0         500        1000       1500       2000
               |----------|----------|----------|----------|
 caja A        [=====================)
 caja B                              [=====================)
                                     ^
                                     A termina en 1000 (no es suyo);
                                     B empieza en 1000 (sí es suyo).
                                     Se tocan, no se enciman.
```

De ahí salen las dos comparaciones del motor:

| Pregunta | Comparación | Por qué |
|---|---|---|
| ¿Se enciman? | estricta, `<` | terminar donde la otra empieza es tocarse, no encimarse |
| ¿Está dentro? | `<=` | terminar justo en la pared es ocupar hasta el último milímetro, sin pasarse |

## `Intersects`: ¿se enciman?

Dos cajas se enciman si comparten volumen, y eso pasa **solo si sus rangos se traslapan en los tres ejes a la vez**.
Basta un eje donde estén separadas, o donde solo se toquen, para que no se encimen.

```
 Sobre un eje, A y B se traslapan si:   A.inicio < B.final   y   B.inicio < A.final

   se traslapan          se tocan              separadas por 1 mm
   [=====)               [=====)               [=====)
      [=====)                  [=====)                 [=====)
   sí                    no                    no
```

### Cara, arista y vértice

Las cajas pueden tocarse de tres formas, y ninguna es encimarse. Vistas desde arriba (`x` a la derecha, `z` hacia
abajo), con A de 1 000 mm por lado en el origen:

```
 por una cara             por una arista           por un vértice
 (comparten un plano)     (comparten una línea)    (comparten un punto)

 ┌─────┬─────┐            ┌─────┐                  ┌─────┐
 │  A  │  B  │            │  A  │                  │  A  │
 └─────┴─────┘            └─────┼─────┐            └─────┼─────┐
                                │  B  │                  │  B  │
                                └─────┘                  └─────┘
 B en x = 1000            B en x = 1000            igual que la arista, pero
 las dos en el piso       y z = 1000, en el piso   B además está arriba (y = 1000)
```

En una cara, los rangos se traslapan en dos ejes y se tocan en uno. En una arista, se traslapan en uno y se tocan en
dos. En un vértice, se tocan en los tres. En los tres casos hay al menos un eje donde **no** se traslapan, así que
`Intersects` da `false`.

### El caso de la cruz

Una forma ingenua de revisar si dos cajas chocan es ver si alguna esquina de una queda dentro de la otra. Falla con
dos barras cruzadas:

```
 Vista desde arriba
                 ┌───┐
                 │ B │
         ┌───────┼───┼───────┐
         │   A   │▓▓▓│       │     ▓ volumen compartido
         └───────┼───┼───────┘
                 │   │
                 └───┘
```

Ninguna de las ocho esquinas de A queda dentro de B, ni al revés, pero comparten el cuadro del centro. Revisar los
tres rangos lo detecta: en `x` el rango de B queda dentro del de A, en `z` el de A dentro del de B, y en `y` las dos
están en el piso. Es uno de los casos de `box_test.go`.

## `ContainedIn`: ¿queda dentro?

Una caja queda dentro de otra si en los tres ejes empieza en o después del inicio y termina en o antes del final.
Pegada a la pared **sí** cuenta como dentro; pasarse 1 mm, no.

```
 Vista desde arriba de un 20' (5 898 × 2 352)

        x = 0                                       x = 5898
 z = 0    ┌───────────────────────────────────────────┐
          │┌─────┐                             ┌─────┐│
          ││  A  │                             │  B  ││
          │└─────┘                             └─────┘│
          │                                     ┌─────┼┐
          │                                     │  C  ││
          │                                     └─────┼┘
 z = 2352 └───────────────────────────────────────────┘

 A pegada al fondo y a la pared izquierda: dentro.
 B pegada a la puerta (termina en x = 5898): dentro.
 C termina en x = 5899: se pasa 1 mm, fuera.
```

Para un tramo abierto con margen (un flat rack con `overhang.leftMm`), el espacio permitido empieza en `z` negativa,
y una pieza con `zMm` negativa puede quedar dentro de él. `ContainedIn` no sabe de márgenes: solo compara dos cajas.
Quien la llama arma la caja del espacio permitido.

## `Contact`: ¿cuánto se apoya?

`top.Contact(below)` es el área, en mm², donde la cara de abajo de `top` descansa sobre la cara de arriba de
`below`. La estabilidad del sprint 2.1 la usará para saber qué parte de la base de una pieza tiene apoyo.

Primero se revisa la altura: hay contacto solo si `top` empieza justo donde termina `below`
(`top.YMM == below.YMM + below.DYMM`).

```
 Vista de lado

   apoyo                 1 mm de aire          encimadas 1 mm
   ┌──────┐              ┌──────┐              ┌──────┐
   │ top  │              │ top  │              │ top  │
   ├──────┤              └──────┘              ├▓▓▓▓▓▓┤  ▓ traslape
   │below │              ┌──────┐              │below │
   └──────┘              │below │              └──────┘
                         └──────┘
   hay contacto          0                     0 (es un error, no apoyo)
```

Si la altura coincide, el área es la parte común de las dos bases vistas desde arriba: el traslape en `x` por el
traslape en `z`.

```
 Vista desde arriba (x a la derecha, z hacia abajo)

 x:       0      400              1200     1600
 z = 0    ┌──────────────────────────┐
          │ below                    │
 z = 200  │      ┌───────────────────┼────────┐
          │      │▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓│        │
 z = 800  └──────┼───────────────────┘        │
                 │                      top   │
 z = 1000        └────────────────────────────┘

 below: 1 200 × 800 en el origen.   top: 1 200 × 800 en (400, 200).
 traslape en x: 1200 − 400 = 800    traslape en z: 800 − 200 = 600
 Contact = 800 × 600 = 480 000 mm² (la zona ▓)
```

Si en algún eje el traslape es 0 (las bases solo se tocan por una arista o un vértice), el producto da 0: no hay
apoyo real.

**El orden importa.** `Contact` no es simétrica: `below.Contact(top)` pregunta si `below` está encima de `top`, y da
0. La estabilidad pregunta «¿en qué se apoya esta pieza?», y el apoyo solo viene de abajo.

## Cajas sin volumen

Una caja con alguna medida en cero o negativa no ocupa espacio. La regla del motor es que **nunca se encima con nada,
no se apoya en nada y no cuenta como contenida** en otra, aunque sus coordenadas queden adentro.

```
 Una "caja" de alto 0 dentro de otra:

 ┌──────────────┐
 │   ────────   │   Intersects → false
 │              │   ContainedIn → false
 └──────────────┘   Contact → 0
```

**Por qué así:**

- La validación de la entrada ya rechaza piezas con medidas en cero. Esta regla es la segunda línea de defensa: si un
  error del algoritmo arma una caja así, ninguna de las tres funciones la trata como una pieza real, y la colocación
  no pasa como válida.
- En Go, `Box{}` (el valor cero del struct, con todos los campos en 0) compila sin avisar. Con esta regla, una caja
  «olvidada» no se confunde con una caja en el origen.
- Con volumen positivo, toda caja se encima consigo misma. Es una de las propiedades que el fuzzing de la sesión
  1.2.3 revisa con entradas al azar.

## Desbordamiento

`inicio + medida` y el producto del área se calculan con `int64`. En Go, una suma o un producto de enteros que se pasa
del máximo (cerca de 9,2 × 10¹⁸) **da la vuelta sin avisar**: sale un número negativo, sin error ni `panic`.

Con los límites del contrato no puede pasar. La constante `MaxGeometryMM` (100 000 mm) fija el valor absoluto más
grande que puede tener cualquiera de los seis números de una `Box`, y `Box.InRange()` lo comprueba:

- una pieza mide a lo más 30 000 mm por lado, y en un plan que arma el motor lo más lejos que termina es
  `x = 60 000` (un tramo que empieza en 30 000 y mide 30 000);
- con números de a lo más 100 000, una suma da a lo más 2 × 10⁵ y un área a lo más 4 × 10¹⁰ mm², lejísimos del tope
  de `int64`;
- el traslape sobre un eje nunca pasa de la medida de la pieza más chica, así que el área de contacto de una pieza
  del contrato es a lo más 30 000 × 30 000 = 9 × 10⁸ mm², estén o no las cajas dentro de un tramo. Es uno de los
  casos de `box_test.go`.

**Cuidado con los planes que llegan de fuera.** El contrato todavía no pone máximo a las coordenadas de una
colocación (`xMm`, `yMm`). Si un plan editado a mano trae `xMm` cerca del tope de `int64`, `xMm + dxMm` da la vuelta
a un número negativo e `Intersects` responde «no se enciman» aunque sí se encimen. Quien reciba un plan de fuera
llama a `InRange` **antes** de usar la caja. El máximo en el contrato entra con el validador (sesión 1.4.2, P-029).

## Las propiedades que prueba el fuzzing

Las pruebas de tabla revisan los casos que se nos ocurrieron. El fuzzing (capítulo 08 de la guía de Go) revisa
**propiedades**: cosas que deben ser ciertas para *cualquier* par de cajas, y deja que Go busque un contraejemplo.
Viven en `box_fuzz_test.go`:

| Propiedad | Entradas | Por qué importa |
|---|---|---|
| `Intersects` es simétrica | todo el rango de `int64` | el algoritmo compara en cualquier orden; la fórmula es simétrica aun con desbordamiento |
| Una caja se interseca consigo misma si y solo si tiene volumen | dentro de `InRange` | es la regla de las cajas sin volumen |
| Si A está en B y B en C, A está en C | dentro de `InRange` | la estabilidad del sprint 2.1 encadena contenciones |
| Girar una pieza no cambia su volumen | piezas del contrato | si una orientación perdiera volumen, el motor «ganaría» espacio |
| `Contact` no es negativa ni mayor que la base más chica, y con contacto no hay traslape | dentro de `InRange` | el apoyo no puede sumar más de lo que hay ni coexistir con un choque |

## Candidatas: en qué orientaciones cabe una pieza

`CargoItem.Candidates(at, space)` (en `candidates.go`) junta la geometría con las
[orientaciones](sistema-coordenadas.md): para una pieza y un punto `at`, regresa una caja por cada orientación
permitida que queda dentro de `space` (`ContainedIn`), con la esquina mínima en `at`. Es la pregunta que el algoritmo
del sprint 1.3 hará en cada punto extremo.

```
 Vista de lado: tarima de 1 200 × 800 × 600 (L × W × H) con 700 mm libres al techo.
 Todas empiezan en at y crecen hacia arriba.

                                              ┌───┐
                                              │   │
                                ┌───┐         │   │
 techo ─────────────────────────┼───┼─────────┼───┼──── y = 2393
                                │   │         │   │
   700    ┌──────────┐ ┌──────┐ │LHW│         │WHL│
   libres │   LWH    │ │ WLH  │ │800│         │1200
          │   600    │ │ 600  │ │   │         │   │
 at ──────┴──────────┴─┴──────┴─┴───┴─────────┴───┴──── y = 1693

           caben (sobran 100)    no caben: atraviesan el techo
           LWH, WLH: 600 de alto   LHW, HLW: 800 · WHL, HWL: 1 200

 Quedan LWH y WLH: las dos con H = 600 hacia arriba.
```

Las reglas que salen de los casos reales:

- **«Este lado arriba»** es `allowedOrientations = ["LWH", "WLH"]`: solo se gira sobre el eje vertical y el alto de
  la pieza siempre queda sobre `y`. No hace falta una regla aparte.
- **Piezas con medidas iguales.** Si dos orientaciones dan la misma caja (una pieza cúbica da seis iguales; una con
  `L = W` da tres pares), solo queda **la primera** según el orden de la lista. El algoritmo no prueba dos veces lo
  mismo.
- **Piezas largas.** Un tubo de 12 000 × 300 × 300 no cabe en ninguna orientación de un 20' (5 898 de largo); en un 40'
  (12 032) solo cabe acostado a lo largo, y como `W = H`, LWH y LHW son la misma caja: queda una candidata.
- **Orden determinista.** Las candidatas salen en el orden de `allowedOrientations` (o del contrato si no vino). El
  orden nunca depende de un mapa, que en Go cambia en cada ejecución (capítulo 05 de la guía de Go).
- **El espacio lo decide quien llama.** `Section.Box()` da la caja nominal del tramo, sin margen de rebase. El
  espacio con el margen de `overhang` lo arma el algoritmo en el sprint 1.3.

## Casos límite que prueba el código

| Función | Casos |
|---|---|
| `Intersects` | idénticas; una dentro de otra; encimadas 1 mm; se tocan por una cara (en x, y y z), por una arista y por un vértice; separadas 1 mm en cada eje; cruzadas en «+»; traslapan en dos ejes y no en el tercero; z negativa; sin volumen. Todos en los dos órdenes |
| `ContainedIn` | dentro sin tocar; en la esquina del origen; en la esquina máxima; idéntica; se pasa 1 mm por cada una de las seis caras; más grande; z negativa en un espacio con margen; sin volumen; caja real en un espacio de ancho 0 |
| `Contact` | apoyo completo; parcial; la de arriba más grande; más chica; 1 mm de aire; encimadas 1 mm; al lado; arista; vértice; orden al revés; lado a lado; caja plana arriba o abajo; máximo del contrato (9 × 10⁸ mm²) |
| `Candidates` | sin lista (las seis); «este lado arriba»; cúbica; cúbica con otra lista; dos lados iguales; pieza larga en un 20' y en un 40'; pegada a la puerta; poco alto libre; punto fuera; orientación desconocida; lista vacía |
