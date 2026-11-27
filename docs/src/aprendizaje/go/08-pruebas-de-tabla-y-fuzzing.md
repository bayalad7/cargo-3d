# 08 — Pruebas de tabla y fuzzing

Este capítulo va **en dos partes**:

- **Parte 1 (esta, sesión 1.2.3):** pruebas de tabla y fuzzing, con la geometría del motor: `box_test.go` y
  `box_fuzz_test.go` en `backend/internal/engine/`.
- **Parte 2 (sesión 1.4.5):** benchmarks con `testing.B` y la línea base de ocupación. Se agrega aquí cuando llegue.

Todo lo que se afirma se probó con Go 1.27 con programas aparte antes de escribirlo (ADR-0022).

## 1. Prueba de tabla: los casos que se nos ocurrieron

Una **prueba de tabla** es una lista de casos (un slice de structs) y un solo ciclo que los corre. Así se ve en
`TestContainedInAllowsTouchingTheBorder`:

```go
tests := []struct {
	name         string
	inner, outer Box
	want         bool
}{
	{name: "idéntica al espacio", inner: space, outer: space, want: true},
	{name: "se pasa 1 mm por la puerta (x máxima)", inner: Box{XMM: 5898 - 1199, /* … */}, outer: space, want: false},
	// …
}

for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		if got := tt.inner.ContainedIn(tt.outer); got != tt.want {
			t.Errorf("ContainedIn() = %t, se esperaba %t\ninner = %+v\nouter = %+v", got, tt.want, tt.inner, tt.outer)
		}
	})
}
```

- **El struct anónimo** (`[]struct{ … }{ … }`) se declara y se llena en el mismo lugar: no hace falta un tipo con nombre
  para algo que solo usa esta prueba.
- **`t.Run(nombre, función)`** corre cada caso como una **subprueba**: si una falla, las demás siguen y el mensaje dice
  cuál fue. Los espacios del nombre se vuelven guiones bajos (`idéntica_al_espacio`), y con eso se puede correr un solo
  caso: `go test ./backend/internal/engine -run 'TestContainedInAllowsTouchingTheBorder/idéntica'`.
- **`%+v`** imprime un struct con el nombre de cada campo. Un mensaje de error tiene que bastar para reproducir el caso
  sin abrir el código: por eso el error dice qué dio, qué se esperaba y con qué cajas.
- **La variable del ciclo** (`tt`) es nueva en cada vuelta desde Go 1.22 (capítulo 05), así que la función que se pasa a
  `t.Run` ve su propio caso.

**Qué no hace una tabla:** solo prueba lo que se le ocurrió a quien la escribió. Los casos que nadie pensó (una caja de
1 mm junto a una de 30 000 mm, tres cajas encadenadas, un número negativo) quedan sin probar. Para eso sirve el fuzzing.

## 2. Fuzzing: propiedades en lugar de casos

El **fuzzing** deja que Go invente las entradas. No se escribe «con estas dos cajas el resultado es `false`», sino una
**propiedad**: algo que tiene que ser cierto para *cualquier* entrada. Si Go encuentra una entrada que la rompe, la
guarda y avisa.

Las propiedades de la geometría (en `box_fuzz_test.go`, una función `Fuzz…` por propiedad):

| Función | Propiedad |
|---|---|
| `FuzzIntersectsIsSymmetric` | `a.Intersects(b)` y `b.Intersects(a)` dan siempre lo mismo |
| `FuzzBoxIntersectsItselfOnlyWithVolume` | una caja se encima consigo misma si y solo si tiene volumen |
| `FuzzContainedInIsTransitive` | si A está en B y B en C, A está en C |
| `FuzzRotationKeepsVolume` | girar una pieza no cambia su volumen (y cada candidata de `Candidates` ocupa ese volumen) |
| `FuzzContactIsBoundedAndExcludesOverlap` | `Contact` no es negativa ni pasa de la base más chica; con contacto no hay traslape |

### Cómo se escribe

```go
func FuzzIntersectsIsSymmetric(f *testing.F) {
	f.Add(int64(0), int64(0), /* … */)          // semillas

	f.Fuzz(func(t *testing.T, ax, ay, az, adx, ady, adz, bx, by, bz, bdx, bdy, bdz int64) {
		a := Box{XMM: ax, /* … */}
		b := Box{XMM: bx, /* … */}
		if a.Intersects(b) != b.Intersects(a) {
			t.Fatalf("no es simétrica\na = %+v\nb = %+v", a, b)
		}
	})
}
```

- **El nombre empieza con `Fuzz`** y recibe un `*testing.F` (no un `*testing.T`).
- **`f.Add(…)`** agrega una **semilla**: una entrada de partida, con los mismos tipos que los argumentos del cuerpo.
  Aquí son los casos límite de las pruebas de tabla (se tocan por una cara, por una arista, separadas 1 mm…).
- **`f.Fuzz(func(t *testing.T, …))`** recibe el cuerpo. Los argumentos después de `t` son los que Go inventa.
- **Tipos permitidos:** `string`, `bool`, `[]byte`, `float32`, `float64` y los enteros (`int`, `int8` a `int64`, `uint`,
  `uint8` a `uint64`). Un struct como argumento no compila: el compilador dice
  `fuzzing arguments can only have the following types`. Por eso una caja entra como seis `int64`.
- **El número y el tipo de los valores de `f.Add` tienen que coincidir** con los argumentos: si no, tampoco compila
  (`wrong number of values in call to (*testing.F).Add`).

### Cómo funciona por dentro

1. Con **`go test` normal** (el que corres siempre) un `Fuzz…` corre **solo con sus semillas** y con los casos
   guardados en `testdata/fuzz/<NombreDelFuzz>/`. Es una prueba rápida de siempre.
2. Con **`-fuzz`**, Go parte de las semillas y las **muta** (cambia números, mezcla entradas) miles de veces por
   segundo. Se guía por la **cobertura**: si una entrada recorre una parte del código que ninguna otra había recorrido,
   la guarda como «interesante» y sigue mutando a partir de ella. Por eso no es azar puro.
3. Si una entrada rompe la propiedad, Go la **minimiza** (la reduce a lo más simple que sigue fallando), la guarda en
   `testdata/fuzz/` y se detiene.

## 3. Cómo se corre

```bash
go test ./backend/internal/engine -run '^$' -fuzz '^FuzzIntersectsIsSymmetric$' -fuzztime 1m
```

| Parte | Para qué |
|---|---|
| `-run '^$'` | Salta las pruebas normales. Sin esto, `go test -fuzz` **primero corre todas las pruebas normales** y recién después fuzzea: se pierde tiempo o se confunde un fallo normal con uno del fuzzing |
| `-fuzz '^NombreExacto$'` | Qué `Fuzz…` generar. **Solo uno por vez**: si el patrón coincide con dos, Go se niega (`will not fuzz, -fuzz matches more than one fuzz test`). Por eso el patrón lleva `^` y `$` |
| `-fuzztime 1m` | Cuánto tiempo. También se puede pedir un número de ejecuciones (`-fuzztime 1000x`). Sin esta opción corre hasta que se pare con Ctrl+C o falle |

Cómo se lee lo que imprime:

```
fuzz: elapsed: 39s, execs: 11748694 (304912/sec), new interesting: 3 (total: 15)
```

- `execs`: cuántas entradas ha probado, y a qué ritmo (aquí, unas 300 000 por segundo).
- `new interesting` (`total`): entradas que abrieron código nuevo. Cuando el número deja de subir, el fuzzing ya
  recorrió lo que puede alcanzar.
- Termina con `PASS` si no encontró nada en el tiempo pedido. **Un `PASS` no prueba que la propiedad sea cierta:**
  solo que no se encontró un contraejemplo en esa corrida.

Los cinco comandos, uno por propiedad, están en `CLAUDE.md`, sección «Comandos».

**Dónde guarda Go lo que va aprendiendo:** el corpus de entradas «interesantes» se guarda en la **caché de Go**
(`$(go env GOCACHE)/fuzz/…`; `go env GOCACHE` dice dónde queda en tu computadora), no en el repositorio. Solo las
entradas que **fallan** se escriben dentro del repo. La caché se borra con `go clean -fuzzcache`.

## 4. Qué hacer cuando el fuzzing encuentra una falla

Para ver una de verdad, se rompió `Intersects` a propósito **en una copia fuera del repo** (cambiando un `<` por `<=`
sobre el eje `z`) y se corrió `FuzzIntersectsIsSymmetric`. Go la encontró casi al instante:

```
--- FAIL: FuzzIntersectsIsSymmetric (0.12s)
    box_fuzz_test.go:92: no es simétrica: a.Intersects(b) = false, b.Intersects(a) = true
        a = {XMM:55 YMM:0 ZMM:0 DXMM:1000 DYMM:1000 DZMM:1000}
        b = {XMM:1000 YMM:984 ZMM:1000 DXMM:1000 DYMM:1000 DZMM:1000}

    Failing input written to testdata/fuzz/FuzzIntersectsIsSymmetric/6ac6aec8a83a70cd
    To re-run:
    go test -run=FuzzIntersectsIsSymmetric/6ac6aec8a83a70cd
```

El archivo guardado es texto, un valor por renglón, en el orden de los argumentos:

```
go test fuzz v1
int64(55)
int64(0)
…
```

Pasos:

1. **Lee el mensaje.** Dice cuál propiedad falló y con qué entrada. Aquí las dos cajas se tocan por una cara en `z`
   (la de `b` empieza en 1 000, donde termina la de `a`) y `Intersects` dio respuestas distintas según el orden.
2. **Repite el caso** sin fuzzing: `go test -run=FuzzIntersectsIsSymmetric/6ac6aec8a83a70cd`. Y `go test` normal
   también lo repite, porque lee todo `testdata/fuzz/`.
3. **Averigua si el error está en el código o en la propiedad.** Una propiedad mal pensada también falla, por ejemplo
   si promete algo que solo es cierto con números chicos y el fuzzing le da un `int64` enorme. Pregúntate: ¿esta
   entrada es una que el motor puede recibir?
4. **Arregla el código** (o la propiedad, o agrega una validación) y vuelve a correr el caso.
5. **Conserva el archivo de `testdata/` y súbelo a git.** Es una prueba de regresión: desde ahí, cada `go test` revisa
   ese caso. No se borra.

**Un detalle:** cuando lo que falla es una **semilla** (una llamada a `f.Add`), Go no escribe nada en `testdata/`: la
semilla ya está en el código. En una primera corrida con el error puesto sobre el eje `x`, la semilla de «se tocan por
una cara» lo atrapó antes de empezar a generar. Para ver el archivo guardado hizo falta un error que ninguna semilla
cubría, por eso el ejemplo de arriba usa `z`.

## 5. Cómo se piensa una propiedad

Lo que salió de escribir estas cinco:

- **Que la propiedad sea cierta por diseño.** «`Intersects` es simétrica» lo es porque la fórmula no cambia al
  intercambiar las cajas. No se fuzzea algo que *esperamos* que se cumpla casi siempre.
- **Cada propiedad solo atrapa ciertos errores.** Con el mismo error en `x`, `FuzzContactIsBoundedAndExcludesOverlap`
  corrió 30 segundos y **no vio nada**: `Contact` exige además que las caras toquen en `y`, y por ahí no se nota. La
  simetría sí lo vio. Cinco propiedades dan más cobertura que una, y ninguna lo da todo.
- **Acota la entrada a lo que el código promete.** La geometría solo es confiable con números chicos (`Box.InRange`,
  hasta ±100 000 mm; si no, la suma de `int64` da la vuelta). La función `wrap` deja igual un número que ya está en el
  rango del contrato y mete en él a uno que se sale (con el resto de una división, sin desbordar); así las semillas llegan tal cual. Una sola propiedad, la simetría, se prueba con **todo** `int64`, porque se
  cumple aun con desbordamiento.
- **Arma la entrada para que se cumpla la premisa.** «Si A está en B y B en C…» casi nunca se cumple con tres cajas al
  azar. `FuzzContainedInIsTransitive` arma las cajas en cadena (B dentro de C, A dentro de B) y les da un empujón de ±50 mm sobre cada eje
  para que a veces se salgan. Sin eso, la propiedad pasaría sin haberse probado.
- **Comprueba que tu generador no desvíe las semillas.** La primera versión de `wrap` corría todo número por el inicio
  del rango: la semilla «el máximo del contrato» (30 000) llegaba como 1, y «se tocan por una cara» quedaba separada
  100 mm. Las propiedades seguían siendo ciertas y el fuzzing pasaba, pero partía lejos de los casos que importaban. La
  revisión lo encontró; `TestWrapKeepsNumbersInRange` y `TestBoundedBoxKeepsTheSeeds` lo vigilan.
- **Las semillas importan.** Son lo primero que Go muta: las del contrato (1 200 × 800, 5 898 de largo) dirigen la
  búsqueda a números parecidos a los reales.

## 6. Trampas

- **`%` con negativos:** en Go el resto de un número negativo es negativo (`-7 % 5` da `-2`), por eso `wrap` le suma el
  tamaño del rango cuando sale menor que 0.
- **El negativo de `math.MinInt64`:** `-math.MinInt64` vuelve a ser `math.MinInt64` (no cabe en `int64`). Por eso
  `InRange` compara contra los dos lados del rango y no usa un valor absoluto hecho con `-v`.
- **Desbordamiento:** una suma o un producto de `int64` que se pasa del máximo da la vuelta sin aviso. La prueba
  `TestIntersectsIsWrongOutOfRange` muestra cómo `Intersects` responde mal con una coordenada cerca del tope.
- **Pruebas con datos que cambian:** una propiedad no puede depender de la hora, del orden de un mapa o de números al
  azar propios; si falla, tiene que poder repetirse con la misma entrada.

## Dónde verlo en el código

| Idea | Lugar |
|---|---|
| Prueba de tabla con subpruebas | `TestContainedInAllowsTouchingTheBorder` en `box_test.go` |
| Semillas y propiedad con todo `int64` | `FuzzIntersectsIsSymmetric` en `box_fuzz_test.go` |
| Entradas acotadas con `wrap` y `boundedBox` | `box_fuzz_test.go` |
| Premisa armada en cadena | `FuzzContainedInIsTransitive` y `boxInside` |
| Límite de geometría | `MaxGeometryMM` y `Box.InRange` en `box.go` |
| Fuzzing de datos que llegan de fuera | `FuzzDecodeInput` en `decode_test.go` |

## Parte 2: benchmarks (sesión 1.4.5)

Pendiente (la línea base se mide en la sesión 1.4.4 y el capítulo se completa en la 1.4.5). Aquí van `testing.B`, `-bench`, `-benchmem`, cómo leer los resultados y la línea base de ocupación y tiempo del
algoritmo. Incluye medir el costo del mapa `seen` de `Candidates` (P-030).
