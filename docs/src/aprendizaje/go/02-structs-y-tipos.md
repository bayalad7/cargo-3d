# 02 — `struct` y tipos numéricos

Este capítulo explica los primeros tipos reales del motor, en `backend/internal/engine/`. `Container` es el equipo y
`CargoItem` es la pieza. Los archivos son `container.go`, `cargo_item.go` e `input.go`; cada línea está comentada.
Aquí se junta lo importante. Todo lo que se afirma de Go se probó con Go 1.27.

## 1. Un `struct`: campos con nombre y tipo

```go
type Door struct {
	WidthMM  int64 `json:"widthMm"`
	HeightMM int64 `json:"heightMm"`
}
```

Un `struct` agrupa campos, cada uno con su nombre y su tipo fijo. Con él se arman todos los datos del motor.

- **Se crea escribiendo sus campos:** `Door{WidthMM: 2340, HeightMM: 2280}`. Los que no se escriben quedan en su valor
  cero.
- **No hay constructores en el lenguaje.** Si crear un valor necesita lógica (calcular o revisar algo), la convención
  es una función `NewDoor(...)` que lo regresa. Hoy el motor no la necesita.
- **`new(T)` existe, pero no es un constructor:** reserva un `T` en su valor cero y regresa un puntero `*T`. Desde Go
  1.26 también acepta una expresión: `new(int64(400))` regresa un `*int64` que apunta a 400. Las pruebas lo usan para
  llenar campos opcionales.
- **Mayúscula inicial = exportado.** `WidthMM` se ve fuera del paquete `engine`; `derefOrZero` no. La visibilidad es
  por **paquete**: cualquier archivo de `engine` usa `derefOrZero`, pero nadie de fuera. `encoding/json` solo ve los
  campos exportados.
- **Composición, no jerarquía:** un struct contiene otros. `EngineInput` contiene un `Container`, y el `Container`
  contiene sus `Section`.

**Valores cero.** Toda variable nace con un valor definido, pero ese valor no siempre se puede usar:

| Tipo | Valor cero |
|---|---|
| `int64` | `0` |
| `string` | `""` |
| `bool` | `false` |
| puntero (`*Door`), slice (`[]Section`), mapa (`map[string]int`) | `nil` |
| struct | cada campo en su valor cero |

Qué se puede hacer con un `nil`:

| Valor `nil` | Leer | Escribir |
|---|---|---|
| slice | `len(s)` da 0 y `range` no hace nada; `append` funciona | `s[0] = x` detiene el programa (índice fuera de rango) |
| mapa | `m["x"]` da el valor cero y `len(m)` da 0 | `m["x"] = 1` **detiene el programa** (`assignment to entry in nil map`) |
| puntero | `p == nil` se puede preguntar | `*p` o `p.Campo` **detienen el programa** (*panic*) |

Por eso el motor pregunta `!= nil` antes de leer un puntero (sección 4).

### 1.1. Qué se copia y qué se comparte

Un struct es un **valor**: `b := a` **copia** sus campos. La copia es superficial: un campo que es puntero o slice se
copia como referencia, y los dos structs siguen viendo los mismos datos.

Con el `Container` del proyecto:

```go
original := Container{
	ID:       "c-20dv",
	Sections: []Section{{ID: "main", LengthMM: 5898}},
	Door:     &Door{WidthMM: 2340, HeightMM: 2280},
}

copia := original
copia.ID = "c-copia"            // campo simple: solo cambia la copia
copia.Sections[0].LengthMM = 1  // slice: cambia el tramo de LOS DOS
copia.Door.WidthMM = 1          // puntero: cambia la puerta de LOS DOS

fmt.Println(original.ID, original.Sections[0].LengthMM, original.Door.WidthMM)
// c-20dv 1 1
```

| Campo | Al copiar el struct | ¿La copia puede cambiar el original? |
|---|---|---|
| `ID string`, `MaxPayloadKG int64` | Se copia el valor | No |
| `Sections []Section` | Se copia la referencia a la misma lista | **Sí**, en sus elementos |
| `Door *Door`, `TareKG *int64` | Se copia la dirección | **Sí**, a través del puntero |

Para una copia independiente, hay que copiar también lo que se comparte:

```go
separada := original
separada.Sections = slices.Clone(original.Sections) // lista nueva con los mismos tramos (paquete slices)
puerta := *original.Door                            // copia del Door (solo si no es nil)
separada.Door = &puerta
```

Importa en el motor porque la metaheurística (sprint 2.2) probará muchos acomodos a partir del mismo equipo, y en
paralelo. Una copia que comparte un slice sin querer hace que un intento modifique el de otro. Los punteros se ven a
fondo en el capítulo 06 y los slices en el 05.

## 2. Por qué `int64` y no `float64`

Todas las medidas (mm) y pesos (kg) del motor son `int64`.

```go
a, b := 0.1, 0.2
fmt.Println(a + b)              // 0.30000000000000004
fmt.Println(int64(1200) + 1200) // 2400, siempre
```

Ojo: `fmt.Println(0.1 + 0.2)` escrito así imprime `0.3`. Go suma las **constantes** escritas tal cual con precisión
exacta al compilar y solo después las convierte a `float64`. El error de redondeo aparece cuando los valores ya son
variables `float64`, que es lo que pasaría con medidas leídas de un JSON.

Las razones:

1. **Nada de traslapes fantasma.** Dos tarimas pegadas, una en `x = 0..1200` y otra en `x = 1200..2400`, se
   **tocan** pero no se traslapan. Con `float64`, después de sumar varias medidas, el borde podría quedar en
   `1199,9999999` y el validador vería un traslape de una diezmilésima de milímetro. Con enteros la comparación
   `a.fin <= b.inicio` es exacta.
2. **El límite de peso es exacto.** 28 000 kg es 28 000, no 27 999,999… ni 28 000,000…1.
3. **Un milímetro alcanza.** Nadie mide una tarima en décimas de milímetro; el contrato decidió enteros (ADR-0016).
4. **No se desborda dentro del contrato.** `int64` llega a 9,2 × 10¹⁸. La pieza más grande del contrato mide
   30 000³ = 2,7 × 10¹³ mm³ y, por 10 000 unidades, 2,7 × 10¹⁷. Hay una prueba que lo confirma
   (`TestCargoItemVolumeTimesMaxQuantityFitsInInt64`). Fuera del contrato sí se desborda: un entero que pasa del tope
   **da la vuelta sin avisar** (una pieza de 10 000 000 mm por lado da `3875820019684212736`, un número sin sentido).
   Por eso la validación de la sesión 1.1.4 corre antes de cualquier cálculo.

¿Y por qué `int64` y no `int`? `int` mide 32 o 64 bits según la máquina; `int64` mide lo mismo en todas, incluso si
el motor se compila a WebAssembly.

Go nunca mezcla tipos numéricos por su cuenta: `int64(5) + 2.5` no compila. Para pasar de uno a otro se escribe la
conversión, por ejemplo `float64(volumen)`. Eso pasará en las métricas, que son fracciones.

## 3. Etiquetas `json`: el puente con el contrato

```go
LengthMM int64 `json:"lengthMm"`
```

- En Go el campo se llama `LengthMM`: la unidad va en el nombre y las siglas en mayúscula, como pide el estilo de Go
  (`ID`, `MM`, `KG`).
- En JSON se llama `lengthMm`, como dice `contracts/engine.schema.json`.
- La etiqueta es texto que el compilador ignora. `encoding/json` la lee mientras el programa corre, con el paquete
  `reflect`.

`,omitempty` significa: «si el campo está en su valor cero, no lo escribas al convertir a JSON». Se usa en los
opcionales.

## 4. «No vino» no es lo mismo que «vino en cero»

Esta fue la parte delicada de la sesión. El contrato tiene valores por omisión: `stackable` = `true`, el margen de
`overhang` = 0 y, si no viene `allowedOrientations`, se permiten las 6. La regla es que **lo que no vino debe
regresar sin venir**.

Ve qué pasaría con un `bool` normal:

```go
Stackable bool `json:"stackable,omitempty"` // ❌
```

La pieza llega con `"stackable": false` (no se le puede poner nada encima). El valor cero de `bool` es `false`, así que
`omitempty` lo omite al escribir. El JSON de regreso ya no dice `stackable` y, por omisión, **la pieza vuelve a ser
apilable**. Un error silencioso y peligroso.

La solución es un **puntero**:

```go
Stackable *bool `json:"stackable,omitempty"` // ✅
```

| JSON | `Stackable` en Go | De regreso a JSON |
|---|---|---|
| (no viene) | `nil` | (no viene) |
| `"stackable": true` | apunta a `true` | `"stackable": true` |
| `"stackable": false` | apunta a `false` | `"stackable": false` |

Un puntero (`*bool`) guarda la **dirección** de un valor que vive en otro lado, o `nil` si no apunta a nada. Dos
cosas que hay que tener presentes:

- si dos structs tienen el mismo puntero, un cambio a través de uno se ve en el otro (sección 1.1);
- leer `*p` con `p == nil` detiene el programa.

Los punteros se ven a fondo en el capítulo 06. En este capítulo solo se usan para decir «no vino».

En el proyecto el criterio quedó así:

| Campo | Tipo | Por qué |
|---|---|---|
| `Stackable` | `*bool` | `false` significa algo |
| `TareKG`, `MaxLoadOnTopKG` | `*int64` | `0` es válido y distinto de «no vino» |
| `Overhang.TopMM`, `LeftMM`, `RightMM` | `*int64` | Hay ejemplos con `"topMm": 0` escrito |
| `Door`, `Overhang` | `*Door`, `*Overhang` | Un objeto que puede no venir |
| `AllowedOrientations` | `[]Orientation` | Un slice ya puede ser `nil` |
| `Description`, `DeliverySequence` | `string`, `int64` | Su valor cero no es un dato útil (`deliverySequence` empieza en 1; una descripción vacía no se muestra) |

Por eso el valor por omisión **no se llena al leer el JSON, sino que se calcula con un método**:

```go
func (c CargoItem) IsStackable() bool {
	if c.Stackable == nil {
		return true // valor por omisión del contrato
	}
	return *c.Stackable
}
```

Lo mismo hace `Section.Margin()`, que regresa un `OverhangValues` sin punteros, listo para calcular. Las
orientaciones permitidas por omisión se resuelven en la sesión 1.1.3.

## 5. Métodos: `Volume()`

```go
func (s Section) Volume() int64 {
	return s.LengthMM * s.WidthMM * s.HeightMM
}
```

- `(s Section)` es el **receptor**: un parámetro normal que va antes del nombre del método. Por convención se llama
  con la inicial del tipo. `s.Volume()` es lo mismo que `Section.Volume(s)`.
- Este receptor es **por valor**: el método recibe una **copia** de la sección, con las reglas de la sección 1.1. Un
  cambio a un campo simple (`s.LengthMM = 0`) no se ve afuera. Pero si el receptor fuera `(c Container)`, un cambio a
  `c.Sections[0]` o a `*c.Door` **sí** se vería, porque la copia comparte esos datos.
- La diferencia con el receptor por puntero (`*Section`) es tema de la sesión 1.1.3 y del capítulo 03.

Qué mide cada uno:

- `Section.Volume()` es el volumen nominal del tramo.
- `Container.Volume()` suma el de sus tramos **sin el margen de `overhang`**. Es el denominador de
  `volumeUtilization`, y por eso esa ocupación puede pasar de 1 con carga sobredimensionada.
- `CargoItem.Volume()` es el de **una unidad**. Para el total del tipo de pieza se multiplica por `Quantity`.

## 6. Un tipo con nombre: `Orientation`

```go
type Orientation string

const (
	OrientationLWH Orientation = "LWH" // como viene
	OrientationWLH Orientation = "WLH" // girada 90° sobre el eje vertical
	// …
)
```

`type Orientation string` crea un tipo nuevo que por dentro es un `string`. El tipo **no limita los valores**: Go no
tiene una forma de decir «solo estos seis».

```go
func place(o Orientation) { /* … */ }

place(OrientationLWH)       // ✅ la forma correcta
place("XYZ")                // ⚠️ compila: un texto escrito tal cual se convierte solo
var s string = "LWH"
place(s)                    // ❌ no compila: s es una variable de tipo string
place(Orientation(s))       // ✅ compila: la conversión se escribe a mano
```

Un texto escrito tal cual (`"XYZ"`) es una **constante sin tipo**: Go le asigna el tipo que pida el lugar donde se usa.
Lo único que el tipo con nombre evita es pasar por descuido una **variable** `string` (por ejemplo, un texto que
llegó de otro lado) sin convertirla. Tampoco `encoding/json` revisa el valor: lee `"XYZ"` sin quejarse. Por eso la
validación de la sesión 1.1.4 revisa que sea una de las seis.

## 7. Dos decisiones que no son de Go, pero se ven en el código

### Quién calcula `MaxPayloadKG` y por qué el motor no conoce la NOM-012

`MaxPayloadKG` es el peso máximo **efectivo**. Lo calcula la API (o el catálogo de equipos de la empresa) **antes** de
llamar al motor, y es el **menor** de dos números:

1. la carga útil del equipo (ficha del fabricante o placa CSC);
2. lo que permite la NOM-012-SCT-2 para esa configuración vehicular (T3-S2, T3-S3…) y ese tipo de camino, menos las
   taras.

El motor no conoce la NOM porque la norma cambia y depende del **viaje**, no de la geometría. Si viviera en el motor,
cada ajuste de la norma obligaría a tocar el algoritmo y a volver a probarlo. Así, el motor solo responde «¿cabe y no
pasa de X kg?» y la tabla de la NOM se prueba aparte, en la sesión 2.1.5 (ADR-0016). Además cumple la regla de
`CLAUDE.md`: `internal/engine` no depende de nada que no sea geometría y peso.

### Hacia dónde se redondea cada valor

El motor **no redondea nada**: recibe enteros. Quien captura (la pantalla o la API) redondea así:

| Valor | Hacia | Ejemplo | Por qué |
|---|---|---|---|
| Interior del equipo, puerta, `maxPayloadKg` | **abajo** | 2 352,8 mm → 2 352 | El equipo nunca parece más grande de lo que es |
| Medidas de la pieza | **arriba** | 1 199,2 mm → 1 200 | La pieza nunca parece más chica |
| Peso de la pieza | **arriba**, mínimo 1 kg | 0,4 kg → 1 kg | Nunca parece más ligera |

El equipo se achica y la pieza se agranda: el redondeo solo puede quitar holgura, nunca dejar pasar algo que en la
realidad no cabe. El costo es sobrestimar la paquetería ligera, como mil sobres de 0,4 kg contados como 1 000 kg;
esa alternativa está en ADR-0016.

## 8. Cómo se prueba que coincide con el contrato

`contract_examples_test.go` toma los 27 `input.json` de `contracts/examples/` y hace esto con cada uno:

1. Lo lee como `EngineInput` en **modo estricto** (`DisallowUnknownFields`): si el JSON trae un campo que el struct no
   tiene, falla. `json.Unmarshal` normal lo ignoraría en silencio.
2. Revisa que no quede nada después del objeto: un segundo `Decode` debe llegar al fin del archivo (`io.EOF`).
3. Lo vuelve a escribir con `json.Marshal`.
4. Compara el original con el resultado **como valores** (`map[string]any` y `[]any`), así no importan el orden de
   las llaves ni los espacios. Al leer así, `encoding/json` convierte todos los números a `float64`, que guarda enteros
   exactos hasta 2⁵³ (≈ 9 × 10¹⁵); los del contrato son mucho menores.

Si alguien cambia `Stackable` a `bool` normal, fallan todos los ejemplos con alguna pieza `"stackable": false` (hay
14 piezas así). Los `output.json` se prueban igual en la sesión 1.1.3, cuando exista `LoadPlan`.

```bash
go test ./backend/internal/engine -run RoundTrip -v
```

## Para practicar

1. Cambia `TareKG *int64` a `TareKG int64` (con `omitempty`) y corre las pruebas. ¿Falla algún ejemplo? ¿Fallaría
   si alguno trajera `"tareKg": 0`? Regrésalo.
2. Escribe en una prueba `Door{}` y muestra sus campos con `t.Logf("%+v", Door{})`. ¿Qué valores salen?
3. Copia un `Container` con `b := a`, cambia `b.Sections[0].LengthMM` y revisa `a`. Después haz la copia con
   `slices.Clone` y repite.
4. Responde sin ver el capítulo: ¿por qué `int64` y no `float64`? ¿Quién redondea, y hacia dónde, la carga útil?
