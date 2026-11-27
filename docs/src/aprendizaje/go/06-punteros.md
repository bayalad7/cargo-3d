# 06 — Punteros: cuándo `*Container` y cuándo `Container`

Este capítulo usa el código del motor en `backend/internal/engine/`: los campos opcionales de `cargo_item.go` y
`container.go`, `load_plan.go` (`LoadPlan.AddPlacement`), `errors.go` (`FieldError`) y `validate.go`. Todo lo que se
afirma se probó con Go 1.27 con programas aparte antes de escribirlo (ADR-0022), y los tamaños de los tipos se midieron
con `unsafe.Sizeof`.

## 1. Un puntero guarda dónde vive un valor

```go
var p *int64      // p vale nil: no apunta a nada
n := int64(5)
p = &n            // &n es la dirección de n; p apunta a n
fmt.Println(*p)   // *p lee el valor al que apunta: 5
*p = 7            // *p también escribe: n ahora vale 7
```

- **`*T` en un tipo** es «puntero a `T`». **`&x`** da la dirección de `x`. **`*p`** sigue la dirección y llega al valor
  («desreferenciar»).
- El **valor cero** de un puntero es `nil`: no apunta a nada.
- **Desreferenciar un `nil` detiene el programa:** `*p` con `p == nil` da
  `panic: runtime error: invalid memory address or nil pointer dereference`. Por eso el motor pregunta antes:

```go
func (c CargoItem) IsStackable() bool {
	if c.Stackable == nil {      // no vino en el JSON: valor por omisión
		return true
	}
	return *c.Stackable          // vino: el valor que trae
}
```

- **Crear un puntero a un valor sin variable intermedia:** desde Go 1.26, `new` acepta una expresión. `new(false)`
  regresa un `*bool` que apunta a `false`, y `new(int64(400))` un `*int64` que apunta a 400. Las pruebas lo usan:
  `CargoItem{Stackable: new(true)}`.

## 2. Para qué usa punteros este proyecto

Hay dos razones reales. Los demás punteros del motor son consecuencia de ellas.

### Para distinguir «no vino» de «vino en cero»

`Stackable *bool`, `MaxLoadOnTopKG *int64`, `DeliverySequence *int64`, `Container.TareKG *int64` y los tres campos de
`Overhang` son punteros porque en el contrato **faltar** y **valer 0 o `false`** son cosas distintas. Un `bool` solo
tiene dos valores, y `false` ya significa «no se puede apilar»: no queda uno libre para decir «no se dijo». Un
`*bool` tiene tres: `nil` (no vino), `false` y `true`.

Lo mismo para estructuras completas: `Container.Door *Door` (`nil` = equipo sin puerta) y
`Section.Overhang *Overhang` (`nil` = no se permite rebasar).

`derefOrZero` (en `container.go`) es la forma corta de leer uno de estos opcionales cuando `nil` equivale a 0:

```go
func derefOrZero(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}
```

### Para que un método cambie el original

```go
func (p *LoadPlan) AddPlacement(pl Placement) {
	p.Placements = append(p.Placements, pl)
}
```

El receptor es `*LoadPlan`: `p` apunta al plan de quien llama, así que el `append` se queda en ese plan. Con un receptor
por valor (`func (p LoadPlan) AddPlacement…`), `p` sería una **copia** del plan: el `append` cambiaría la copia y el plan
original quedaría igual, sin error ni aviso. Lo comprobé:

```go
pl := NewPlan()
pl.AddCopy(1)   // receptor por valor
fmt.Println(len(pl.N)) // 0: se agregó a una copia
pl.Add(1)       // receptor puntero
fmt.Println(len(pl.N)) // 1
```

Sorpresa útil: `pl.Add(1)` funciona aunque `pl` no sea un puntero. Si `pl` es una variable, Go escribe `(&pl).Add(1)` por ti.

## 3. `Container` o `*Container`: cómo se elige

La misma pregunta sirve para el receptor de un método y para el parámetro de una función. **Estos son los tamaños del
motor** (medidos en Go 1.27, macOS arm64, bytes):

| Tipo | Bytes | Tipo | Bytes |
|---|---|---|---|
| `Point`, `Dimensions`, `OverhangValues` | 24 | `Placement` | 112 |
| `Box` | 48 | `CargoItem` | 120 |
| `Candidate` | 64 | `EngineInput` | 120 |
| `Section` | 72 | `LoadPlan` | 136 |
| `Container` | 80 | | |

Las cuatro preguntas, en orden:

1. **¿El método tiene que cambiar el valor?** Sí → puntero. (`LoadPlan.AddPlacement`, `LoadPlan.AddUnplaced`, los
   métodos del `validator`, que agregan a su lista.)
2. **¿Es un tipo que no se debe copiar?** Sí → puntero. En el motor no hay ninguno todavía; son los que llevan
   candados o recursos (`sync.Mutex`, un archivo abierto).
3. **¿Los demás métodos del tipo ya usan puntero?** Sí → puntero también, **aunque este no cambie nada**. Mezclar
   receptores en un mismo tipo confunde (ver la sección 4). El motor no mezcla: `LoadPlan` tiene solo métodos de
   puntero; `Box`, `CargoItem`, `Container`, `Section` y `Placement`, solo de valor.
4. **Si nada de lo anterior aplica → valor.** Es el caso de casi todo el motor: `Box.Intersects`, `CargoItem.Volume`,
   `Section.Margin`… Leen y calculan. El receptor es una **copia**: si el método cambia un campo directo (`c.LengthMM =
   0`), quien llama no se entera. **Pero la copia no protege lo que está detrás de un puntero, un slice o un mapa**: ahí
   copia solo la dirección (sección 5). Un método de `CargoItem` por valor podría escribir
   `c.AllowedOrientations[0] = …` o `*c.Stackable = true` y cambiar la pieza de quien llamó. Por eso `Orientations()`
   regresa `slices.Clone(c.AllowedOrientations)` y no el slice tal cual (capítulo 05). El receptor por valor hace
   inofensivo lo que se asigna directo, y nada más.

### Tamaño: ¿pasar un puntero es más rápido?

Un puntero mide 8 bytes (en un equipo de 64 bits) y copiar uno es más barato que copiar 120. Se midió en un programa
aparte, **que no está en el repo** (un struct con la forma de `CargoItem`, 1 000 piezas, una función que no se expande
en línea, Go 1.27, macOS arm64). Cada fila cambia una sola cosa:

| Forma | Tiempo por 1 000 llamadas | Por llamada |
|---|---|---|
| `for _, it := range s` y `f(it)` (por valor) | ≈ 1 920 ns | ≈ 1,9 ns |
| `for i := range s` y `f(s[i])` (por valor, sin la copia del `range`) | ≈ 1 805 ns | ≈ 1,8 ns |
| `for i := range s` y `f(&s[i])` (por puntero) | ≈ 740 ns | ≈ 0,74 ns |

Comparando las dos últimas (el mismo ciclo, solo cambia pasar el valor o su dirección), el puntero es unas 2,4 veces
más rápido, y la diferencia **absoluta** es poco más de 1 ns por llamada. La copia extra del `range` sumó apenas ≈ 0,1 ns:
el compilador no la duplica. Solo importa si la función hace casi nada y se llama millones de veces.
`EngineInput.Validate`, por ejemplo, arma textos y recorre listas en cada pieza: copiar 120 bytes es un ruido junto a
eso. Y un puntero tampoco es gratis: obliga a pensar quién más ve ese valor.

**Cómo usar el dato:** es un orden de magnitud, no una regla. Un micro-benchmark depende del compilador y del equipo, y
este no se puede repetir desde el repo; la medición reproducible, con el algoritmo real, es la de la sesión 1.4.4, y
ella es la que decide un cambio. Hasta entonces, se elige por las cuatro
preguntas y no por el tamaño.

## 4. Qué tipo cumple una interfaz: el conjunto de métodos

Con un receptor de puntero, **solo el puntero** cumple la interfaz. `error` pide un método `Error() string`; como
`FieldError` lo define sobre `*FieldError`, el compilador rechaza el valor:

```go
var _ error = E{}
// cannot use E{} (value of struct type E) as error value in variable declaration:
// E does not implement error (method Error has pointer receiver)

var _ error = &E{}   // compila
```

(`E` es un struct de prueba igual a `FieldError` en esto.) Que `FieldError` use receptor de puntero es una **convención
de los tipos de error**, no una necesidad: `Error()` no cambia nada, y con receptor de valor también funcionaría. Se
hace así porque los errores se regresan y se comparan como `*FieldError` (`nil` quiere decir «sin error») y porque
`errors.As` busca ese tipo: quien use el motor espera `*FieldError`. El mismo efecto explica por qué no se mezclan receptores en
un tipo: con `func (T) A()` y `func (*T) B()`, el valor `T` solo tiene el método `A`, y `*T` tiene los dos. Entonces `T{}`
cumple una interfaz que pide solo `A`, y **no** una que pida `A` y `B` (`T does not implement AB (method B has pointer
receiver)`), mientras `&T{}` cumple las dos. Con todos los receptores iguales, `T` y `*T` no se comportan distinto. Por eso el motor crea sus errores con `&FieldError{…}` y
guarda `[]*FieldError` en `ValidationError.Problems`: el slice guarda los mismos valores que `errors.Is` y `errors.As`
reconocen como errores.

## 5. Los punteros dentro de un struct se copian como punteros

Copiar un struct copia sus campos. Si un campo es un puntero, se copia **la dirección**: la copia y el original
apuntan al mismo valor.

```go
a := Item{Tare: new(int64(10))}
b := a            // copia de a
*b.Tare = 99
fmt.Println(*a.Tare)       // 99: a cambió aunque solo se escribió por b
fmt.Println(a.Tare == b.Tare) // true: es la misma dirección
```

En el motor: `for _, s := range c.Sections` recibe en `s` una **copia** de cada `Section` (72 bytes), pero su campo
`Overhang` sigue apuntando al `Overhang` de la entrada. Leerlo está bien (`Section.Margin` solo lee). Escribir por esa
copia cambiaría la entrada de quien llamó al motor. Es la misma idea del capítulo 02, ahora con el campo que causa el
efecto.

Al comparar punteros: `x == y` pregunta si **apuntan al mismo lugar**, no si valen lo mismo. Dos
`&Item{}` distintos son distintos (`false`); para comparar los valores hay que seguir los punteros (`*x == *y`).

## 6. Trampas

- **Puntero `nil`.** Desreferenciarlo detiene el programa (sección 1). Una pregunta `if p == nil` antes de usar un
  opcional cuesta una línea. Un **método** con receptor de puntero sí puede llamarse con un puntero `nil`, y funciona
  mientras no lea nada a través de él:

  ```go
  func (n *Named) Name() string {
  	if n == nil { return "(nil)" }
  	return n.name
  }
  var nm *Named
  nm.Name()   // "(nil)": no hay panic
  ```

- **Puntero a un elemento de un slice, y `append`.** Si el `append` necesita más espacio, crea un arreglo nuevo, y el
  puntero se queda apuntando al arreglo viejo:

  ```go
  s := make([]int, 1, 1)   // sin espacio de sobra
  ps := &s[0]
  s = append(s, 2)         // s se muda a un arreglo nuevo
  *ps = 77
  fmt.Println(s[0])        // 0: se escribió en el arreglo viejo
  ```

  Con espacio de sobra (`make([]int, 1, 4)`) el mismo código sí cambia `s[0]` (77). Que funcione o no depende de la
  capacidad, así que es un error que aparece solo a veces. El motor lo evita al recorrer con índices
  (`p.Unplaced[i]`, no un puntero que se guarde).

- **No se puede tomar la dirección de todo.** Un método de puntero sobre un valor que no está en una variable no compila:

  ```go
  NewLoadPlan("x").AddPlacement(pl)
  // cannot call pointer method AddPlacement on LoadPlan
  ```

  Hay que guardar el plan primero (`plan := NewLoadPlan("x")`). Tampoco se puede `&m["a"]` si `m` es un mapa:
  `cannot take address of m["a"] (map index expression of struct type P)`.

- **Puntero a la variable de un ciclo.** Desde Go 1.22 cada vuelta tiene su propia variable, así que
  `for i := range 3 { ptrs = append(ptrs, &i) }` da tres punteros distintos (0, 1, 2). Con un `go.mod` que declare
  `go 1.21` o menos, el mismo ciclo (escrito con `for i := 0; i < 3; i++`) da tres punteros **iguales**, y los tres
  leen 3. Lo que manda es la versión de `go.mod`, no la del compilador instalado. El `go.mod` del backend pide Go 1.27,
  así que esta trampa no existe aquí.

- **Regresar la dirección de una variable local es seguro.** `func local() *int { x := 42; return &x }` funciona: Go
  mantiene `x` viva mientras haya alguien que apunte a ella.

## 7. Revisión del motor (sesión 1.2.4)

Se revisó cada puntero y cada struct grande que se copia, con las cuatro preguntas. **No se cambió código.**

| Dónde | Qué es | Veredicto |
|---|---|---|
| `Stackable`, `MaxLoadOnTopKG`, `DeliverySequence`, `TareKG`, los tres de `Overhang` | punteros de campos opcionales | **Necesarios.** «No vino» y «vino en 0» son cosas distintas en el contrato |
| `Container.Door`, `Section.Overhang` | puntero a estructura opcional | **Necesarios.** `nil` quiere decir «no hay» |
| `LoadPlan.AddPlacement`, `AddUnplaced` | receptor de puntero | **Necesario.** Cambian el plan; con receptor de valor se perdería lo agregado |
| `FieldError`, `ValidationError` (y `[]*FieldError`) | receptor de puntero | **Correcto.** Sigue la convención de los errores; con receptor de valor serían 48 bytes por copia sin ganar nada |
| métodos del `validator` | receptor de puntero | **Necesarios.** Agregan a `problems` |
| `NewLoadPlan` regresa `LoadPlan` (136 bytes) y no `*LoadPlan` | valor | **Correcto.** Se crea una vez; `LoadPlan` solo tiene métodos de puntero y Go toma la dirección solo |
| `Box` (48), `Section` (72), `Container` (80), `Placement` (112) como receptores y parámetros | valor | **Correcto.** Solo leen; son chicos |
| `validator.item(path, it CargoItem)` y `for i, item := range in.Items` en `Validate` | copian 120 bytes por pieza | **Sin problema.** La validación hace mucho más trabajo por pieza (textos, rangos, ids); la copia es ruido |
| `CargoItem.Candidates` y `RotatedDimensions` (receptor de valor, 120 bytes) | copian al llamarse | **Sin problema hoy.** Son los únicos candidatos a una llamada muy repetida: el algoritmo del sprint 1.3 las llamará por cada punto extremo y cada pieza. Se mide en la 1.4.4 (junto con el mapa `seen`, P-030) antes de decidir; cambiar solo uno de los métodos de `CargoItem` a puntero rompería la regla de no mezclar receptores |
| `EngineInput.Validate` (receptor de valor, 120 bytes) | copia | **Sin problema.** Se llama una vez por petición |

**Resultado:** no se encontró ningún puntero innecesario ni una copia que convenga cambiar hoy. El único lugar que
conviene volver a mirar con datos es `CargoItem.Candidates` en el sprint 1.3 y 1.4.

## Dónde verlo en el código

| Idea | Lugar |
|---|---|
| Opcional con `*bool` y su lectura | `CargoItem.Stackable` e `IsStackable` en `cargo_item.go` |
| `new(true)`, `new(int64(400))` | `cargo_item_test.go` y `container_test.go` |
| Opcional que se lee con `nil` = 0 | `derefOrZero` en `container.go` |
| Receptor de puntero que cambia el original | `LoadPlan.AddPlacement` en `load_plan.go` |
| Solo el puntero cumple `error` | `FieldError.Error` en `errors.go` |
| Copia de `Section` que comparte `Overhang` | `Container.Volume` y `Validate` (`container.go`, `validate.go`) |
| Recorrer por índice para cambiar el original | `LoadPlan.AddUnplaced` en `load_plan.go` |
