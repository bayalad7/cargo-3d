# 05 — Slices, mapas y ciclos

Este capítulo usa las candidatas de la sesión 1.2.2, en `backend/internal/engine/candidates.go`
(`CargoItem.Candidates`), y `CargoItem.Orientations` de `orientation.go`. Todo lo que se afirma se probó con Go 1.27
con un programa aparte antes de escribirlo.

`Candidates` responde: «si pongo esta pieza con su esquina en este punto, ¿en qué orientaciones cabe?». Para eso
recorre una lista (slice), recuerda lo que ya vio (mapa) y arma otra lista con el resultado. Son las tres ideas del
capítulo.

## 1. Un slice es una vista de un arreglo

`[]Orientation` es un **slice**: una lista de largo variable. Por dentro son tres datos:

- un apuntador a un **arreglo de fondo**, donde viven los elementos;
- `len`, cuántos elementos tiene;
- `cap`, cuántos caben en el arreglo de fondo antes de tener que pedir uno nuevo.

```go
errs := make([]error, 0, len(e.Problems)) // errors.go, ValidationError.Unwrap
```

`make([]T, largo, capacidad)` crea el arreglo de fondo. Aquí el largo es 0 (todavía no hay nada) y la capacidad es el
número de problemas: se sabe de antemano cuántos elementos van a entrar, así que los `append` de después nunca tienen
que pedir más memoria. Cuando no se sabe, o son pocos, basta empezar con un slice `nil` y dejar que `append` crezca
solo, como hace `Candidates`.

### `append` regresa el slice: hay que guardarlo

```go
candidates = append(candidates, Candidate{Orientation: o, Box: box})
```

`append` agrega al final y **regresa** el slice resultante. Si el elemento cabía en la capacidad, usa el mismo arreglo
de fondo; si no, crea uno más grande, copia todo y regresa un slice que apunta al nuevo. Por eso siempre se escribe
`s = append(s, x)`: si se tira lo que regresa, el elemento se pierde.

### Slice `nil` contra slice vacío

```go
var candidates []Candidate // nil: sin arreglo de fondo
empty := []Candidate{}      // vacío, pero no nil
```

Los dos miden 0, `for range` no da ninguna vuelta sobre ellos y `append` funciona igual en los dos. La diferencia
es que `candidates == nil` es `true` y `empty == nil` es `false`. En el motor importa en dos lugares:

- `Candidates` regresa `nil` cuando no cabe ninguna orientación: no aparta memoria si no hace falta. La prueba revisa
  que sea `nil` y no un slice vacío.
- En `CargoItem.AllowedOrientations`, `nil` quiere decir **«no vino la lista»** (las seis) y una lista vacía es un
  error de la entrada (P-025). No son lo mismo.

### Trampa: dos slices que comparten el arreglo de fondo

Copiar un slice (`b := a`, pasarlo a una función, guardarlo en un campo) copia los tres datos de adentro, **no** los
elementos. Los dos apuntan al mismo arreglo:

```go
s := []int{1, 2, 3}
t := s[:2]   // t ve los dos primeros de s
t[0] = 99
// s[0] ahora vale 99
```

Con `append` es más tramposo, porque depende de la capacidad:

```go
base := make([]int, 3, 6) // largo 3, lugar para 6
x := append(base, 10)     // cabe: escribe en base[3]
y := append(base, 20)     // cabe: escribe en el MISMO base[3]
// x[3] y y[3] valen los dos 20
```

Por eso `Orientations` regresa `slices.Clone(c.AllowedOrientations)` y no el slice tal cual: `c` es una copia de la
pieza (receptor por valor), pero su slice sigue apuntando a los elementos de la entrada. Sin `Clone`, quien cambiara
la lista que regresa cambiaría la entrada. `TestOrientationsReturnsACopy` (en `orientation_test.go`) lo revisa:
cambia lo que regresa `Orientations()` y comprueba que la pieza siga igual.

### Trampa: un parámetro variádico sin argumentos es `nil`

En la prueba, `withOrientations(item CargoItem, os ...Orientation)` recibe las orientaciones como **parámetro
variádico**: se le pueden pasar cero, una o varias, y adentro `os` es un `[]Orientation`.

```go
withOrientations(box1200, OrientationLWH, OrientationWLH) // os = [LWH WLH]
withOrientations(box1200)                                 // os = nil (¡no vacío!)
withOrientations(box1200, []Orientation{}...)             // os = [] (vacío, no nil)
```

Sin argumentos, `os` es `nil`, y para el motor `nil` es «no vino la lista»: salen las seis orientaciones, no ninguna.
Este error salió en la primera versión de la prueba «lista vacía». Para pasar un slice ya hecho se escribe `s...`; ese
slice **se comparte**, no se copia.

## 2. Un mapa como conjunto

```go
seen := make(map[Dimensions]bool, len(orientations))
```

`map[Llave]Valor` guarda pares llave–valor y encuentra una llave sin recorrer todo. Con valores `bool` sirve como
**conjunto**: «¿ya vi estas medidas?».

```go
if seen[d] {   // leer una llave que no está regresa el valor cero: false
	continue
}
seen[d] = true
```

- **Llaves struct.** La llave es `Dimensions`, un struct. Se puede porque todos sus campos son comparables (`int64`):
  dos `Dimensions` con los mismos tres números son la misma llave. Un struct con un slice adentro no puede ser llave.
- **Saber si existe.** `v, ok := seen[d]` regresa además `ok`, que dice si la llave estaba. Con un conjunto de `bool`
  no hace falta: «no está» y «está en false» se tratan igual.
- **Trampa: mapa `nil`.** `var m map[string]int` es un mapa `nil`. Leerlo funciona (regresa el valor cero) y
  `len(m)` es 0, pero **escribir en él detiene el programa** (`panic: assignment to entry in nil map`). Por eso
  `seen` se crea con `make`.

### Para qué sirve aquí: piezas cúbicas

Una caja de 1 000 × 1 000 × 1 000 mide lo mismo en las seis orientaciones. Sin el conjunto, `Candidates` regresaría
seis candidatas idénticas y el algoritmo probaría seis veces lo mismo. Con el conjunto queda una:

| Pieza | Orientaciones que dan cajas distintas | Candidatas |
|---|---|---|
| 1 200 × 800 × 600 (tres medidas distintas) | las seis | 6 |
| 1 200 × 1 200 × 1 000 (L = W) | LWH, LHW, HLW | 3 |
| 1 000 × 1 000 × 1 000 (cúbica) | LWH | 1 |
| tubo 12 000 × 300 × 300 en un 40' (W = H) | LWH (LHW da la misma) | 1 |

Se queda **la primera** de cada grupo según el orden de la lista: con `allowedOrientations = ["WLH", "LWH"]`, la
cúbica sale como WLH.

### Trampa: el orden de un mapa cambia en cada ejecución

Recorrer un mapa con `for k := range m` da las llaves en un orden que Go cambia a propósito en cada ejecución. En una
prueba, 50 recorridos de un mapa de 20 llaves dieron más de diez órdenes distintos (el número cambia en cada corrida). El motor tiene que dar el **mismo plan con la misma
entrada**, así que nunca decide recorriendo un mapa: el orden de las candidatas lo da el slice de orientaciones, y el
mapa solo responde «¿ya la vi?».

## 3. Ciclos

Go tiene una sola palabra para ciclos: `for`.

```go
for _, o := range orientations {
	d, ok := c.RotatedDimensions(o)
	if !ok {
		continue
	}
	// …
}
```

- **`range` sobre un slice** da en cada vuelta el índice y una **copia** del elemento. `_` descarta el índice, que
  aquí no se usa.
- **`continue`** salta a la siguiente vuelta; `break` saldría del ciclo.
- **`range` sobre un entero** (desde Go 1.22): `for i := range 3` da 0, 1 y 2.

### Trampa: la variable de `range` es una copia

```go
for _, d := range ds {
	d.DXMM = 50   // cambia la copia, no ds
}
for i := range ds {
	ds[i].DXMM = 50   // cambia el elemento de ds
}
```

Para cambiar los elementos de un slice dentro de un ciclo hay que usar el índice.

### Cada vuelta tiene su propia variable

Desde Go 1.22, la variable del `for` es **nueva en cada vuelta**. Si una función guardada en una vuelta usa `i`, ve
el valor de esa vuelta:

```go
var fs []func() int
for i := range 3 {
	fs = append(fs, func() int { return i })
}
// fs[0](), fs[1](), fs[2]() dan 0, 1 y 2
```

El módulo del backend pide Go 1.27 en `go.mod`, así que esta regla siempre aplica en el proyecto.

## Dónde verlo en el código

| Idea | Lugar |
|---|---|
| `make` con capacidad | `ValidationError.Unwrap` en `errors.go` |
| `append` sobre un slice `nil` | `CargoItem.Candidates` en `candidates.go` |
| Mapa como conjunto con llave struct | `seen` en `CargoItem.Candidates` |
| `slices.Clone` para no compartir el arreglo de fondo | `CargoItem.Orientations` en `orientation.go`, probado por `TestOrientationsReturnsACopy` |
| Parámetro variádico y `s...` | `withOrientations` en `candidates_test.go` |
| `slices.Equal` para comparar slices (no se puede con `==`) | `TestCandidatesKeepAllowedOrientationsThatFit` |
