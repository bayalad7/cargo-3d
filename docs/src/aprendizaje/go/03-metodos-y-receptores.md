# 03 — Métodos y receptores

Este capítulo usa los tipos de la sesión 1.1.3, en `backend/internal/engine/`: `Orientation` y `Dimensions`
(`orientation.go`), `Placement` (`placement.go`) y `LoadPlan` (`load_plan.go`). Todo lo que se afirma se probó con
Go 1.27.

## 1. Qué es un método

```go
func (p Placement) Volume() int64 {
	return p.DXMM * p.DYMM * p.DZMM
}
```

Un método es una función con un **receptor**: el parámetro que va entre paréntesis antes del nombre, `(p Placement)`.

- `pl.Volume()` es exactamente `Placement.Volume(pl)`: el receptor es un parámetro más, escrito en otro lugar.
- Por convención el receptor se llama con la inicial del tipo (`p`, `c`) y con el mismo nombre en todos los métodos
  del tipo.
- Un método solo se puede declarar en el **mismo paquete** que su tipo. `engine` puede agregar métodos a
  `CargoItem`, pero no a `string` ni a un tipo de otro paquete.
- Un tipo no puede tener un campo y un método con el mismo nombre. Por eso el campo es `AllowedOrientations` y el
  método es `Orientations()`.

## 2. Receptor por valor

```go
func (c CargoItem) Orientations() []Orientation
func (p Placement) Volume() int64
```

Con `(c CargoItem)`, el método recibe una **copia** de la pieza. Lo que haga con sus campos simples no se ve afuera.

Pero la copia es superficial (capítulo 02, sección 1.1): `c.AllowedOrientations` es un slice y la copia **comparte**
sus elementos con la pieza original. Por eso `Orientations()` no regresa ese slice tal cual:

```go
func (c CargoItem) Orientations() []Orientation {
	if c.AllowedOrientations == nil {
		return allOrientations() // ya es una lista nueva
	}
	return slices.Clone(c.AllowedOrientations)
}
```

Sin `slices.Clone`, quien recibiera el resultado podría hacer `got[0] = OrientationHWL` y cambiar la entrada. La
prueba `TestOrientationsReturnsACopy` lo revisa.

El valor por omisión sale de una **función**, `allOrientations()`, que arma una lista nueva en cada llamada. Si fuera
una variable del paquete, cualquier archivo de `engine` podría cambiar sus elementos y cambiaría el valor por omisión
para todo el programa.

## 3. Receptor por puntero

```go
func (p *LoadPlan) AddPlacement(pl Placement) {
	p.Placements = append(p.Placements, pl)
}
```

Con `(p *LoadPlan)`, el método recibe la **dirección** del plan. Lo que cambia se queda en el plan.

Ve qué pasaría con receptor por valor:

```go
func (p LoadPlan) addByValue(pl Placement) {
	p.Placements = append(p.Placements, pl) // cambia la COPIA
}

plan := NewLoadPlan("c-20dv")
plan.addByValue(Placement{ItemID: "a"})
fmt.Println(len(plan.Placements)) // 0: la colocación se perdió
```

`append` regresa un slice nuevo (con otro largo y, si no cupo, en otro lugar de la memoria) y ese slice se guarda en
la copia. Al terminar el método la copia desaparece. `go vet` no avisa de este error: compila y corre sin quejarse.

### Go toma la dirección por ti

```go
plan := NewLoadPlan("c-20dv")   // plan es un LoadPlan, no un *LoadPlan
plan.AddPlacement(pl)           // Go lo convierte en (&plan).AddPlacement(pl)
```

Funciona porque `plan` es una variable: tiene una dirección. Con algo que no tiene dirección, no compila:

```go
NewLoadPlan("c-20dv").AddPlacement(pl) // ❌ no se puede tomar la dirección del resultado de una función
// error del compilador: cannot call pointer method AddPlacement on LoadPlan
```

Al revés también funciona: un método con receptor por valor se puede llamar desde un puntero (`pp.Volume()`); Go
copia el valor al que apunta. Si el puntero es `nil`, esa copia detiene el programa.

## 4. Cuándo usar cada uno

| Usa receptor por puntero (`*T`) cuando… | Usa receptor por valor (`T`) cuando… |
|---|---|
| el método **cambia** el valor: `AddPlacement`, `AddUnplaced` | el método solo **lee**: `Volume`, `Orientations`, `RotatedDimensions`, `IsStackable` |
| el struct es grande y copiarlo en cada llamada cuesta | el tipo es chico y se usa como un dato: `Placement`, `Dimensions`, `Orientation` |
| otro método del mismo tipo ya es por puntero | ningún método del tipo necesita cambiarlo |

La última fila es una convención de Go: **no mezclar** los dos receptores en un mismo tipo. Si un método de
`LoadPlan` necesita puntero, todos sus métodos lo usan. Así nadie tiene que recordar cuáles modifican el plan y
cuáles no.

En el proyecto quedó así:

| Tipo | Receptor | Por qué |
|---|---|---|
| `CargoItem`, `Container`, `Section`, `Placement` | valor | Son la entrada o el resultado: el motor los lee y no los cambia |
| `LoadPlan` | puntero | Se va llenando con `AddPlacement` y `AddUnplaced` |

El tamaño no decidió nada aquí: copiar un `LoadPlan` copia dos encabezados de slice (`Placements` y `Unplaced`),
dos de `string` y las métricas, no las colocaciones.

## 5. `range` y la copia del elemento

`AddUnplaced` suma unidades a un grupo que ya existe. Esta versión **no funciona**:

```go
for _, u := range p.Unplaced {
	if u.ItemID == itemID && u.Reason == reason {
		u.Quantity += quantity // u es una COPIA del grupo: la lista no cambia
		return
	}
}
```

En `for _, u := range lista`, `u` es una copia de cada elemento. La versión del código recorre los **índices** y
cambia el elemento mismo:

```go
for i := range p.Unplaced {
	if p.Unplaced[i].ItemID == itemID && p.Unplaced[i].Reason == reason {
		p.Unplaced[i].Quantity += quantity
		return
	}
}
```

## 6. Regresar dos valores: «valor, ok»

```go
func (c CargoItem) RotatedDimensions(o Orientation) (Dimensions, bool)
```

Una orientación que no es una de las seis no detiene el programa. La función regresa las medidas y un `bool`:

```go
d, ok := item.RotatedDimensions(pl.Orientation)
if !ok {
	// orientación desconocida: d viene en cero y no se usa
}
```

Es la misma forma con que Go lee un mapa (`item, ok := items[id]`), y la prueba contra los ejemplos usa las dos. Se
usa `bool` porque solo hay una forma de fallar. Cuando hay varias y quien llama necesita saber cuál, se regresa un
`error` (capítulo 04).

## 7. `switch`

```go
switch o {
case OrientationLWH:
	return Dimensions{DXMM: l, DYMM: h, DZMM: w}, true
// …
default:
	return Dimensions{}, false
}
```

- Compara `o` con cada caso, de arriba abajo, y corre solo el primero que coincide.
- Cada caso termina solo; no hace falta nada para salir. Para seguir al siguiente caso existe `fallthrough`, que
  aquí no se usa.
- `default` corre si ninguno coincidió.

## 8. Un método como valor

Un método también se puede guardar en una variable, sin llamarlo, y llamarlo después.

Con receptor por valor, Go guarda **una copia del receptor en ese momento**. Si después el valor cambia, la función
guardada sigue viendo la copia vieja:

```go
item := CargoItem{}
f := item.IsStackable       // método guardado, todavía sin llamar
item.Stackable = new(false) // la pieza cambia
fmt.Println(f(), item.IsStackable())
// true false: f sigue viendo la pieza como era al guardarlo
```

Con receptor por puntero se guarda la dirección y la función ve los cambios. Hoy el motor no guarda métodos así.

## 9. `NewLoadPlan`: cuándo hace falta una función `New`

```go
func NewLoadPlan(containerID string) LoadPlan {
	return LoadPlan{
		SchemaVersion: Version,
		ContainerID:   containerID,
		Placements:    []Placement{},
		Unplaced:      []UnplacedItem{},
	}
}
```

El valor cero de `LoadPlan` no sirve: sus slices son `nil` y `encoding/json` los escribe como `null`, cuando el
contrato exige `[]`. Las dos pruebas lo dejan a la vista: `TestZeroLoadPlanEncodesNull` y
`TestNewLoadPlanEncodesEmptyListsAsArrays`.

`NewLoadPlan` regresa un `LoadPlan`, no un `*LoadPlan`: quien lo recibe lo guarda en una variable y Go toma su
dirección al llamar `AddPlacement`.

## 10. Nota breve: genéricos en las pruebas

`contract_examples_test.go` usa una función con un parámetro de tipo:

```go
func roundTrip[T any](t *testing.T, path string)

roundTrip[EngineInput](t, path)
roundTrip[LoadPlan](t, path)
```

`[T any]` deja que quien llama elija el tipo con que se lee el archivo. El compilador revisa cada uso con su tipo
real. Aquí solo evita escribir la misma prueba dos veces; los genéricos se ven con calma en el capítulo 08.

## Para practicar

1. Cambia `AddPlacement` a receptor por valor (`(p LoadPlan)`) y corre `go test ./backend/internal/engine`. ¿Qué
   prueba falla y por qué? Regrésalo.
2. Cambia el ciclo de `AddUnplaced` a `for _, u := range p.Unplaced` y corre las pruebas. ¿Qué caso falla?
3. Quita `slices.Clone` de `Orientations()` (regresa `c.AllowedOrientations`) y corre las pruebas. ¿Qué pasa?
4. Responde sin ver el capítulo: ¿por qué todos los métodos de `LoadPlan` usan puntero, aunque uno solo leyera?
