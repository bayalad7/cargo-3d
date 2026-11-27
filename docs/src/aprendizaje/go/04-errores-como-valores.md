# 04 — Errores como valores

Este capítulo usa la validación de la sesión 1.1.4, en `backend/internal/engine/`: `errors.go` (las categorías y los
tipos de error), `validate.go` (`EngineInput.Validate`) y `decode.go` (`DecodeInput`). Todo lo que se afirma se probó
con Go 1.27.

## 1. Un error es un valor que se regresa

En Go, un error es un valor como cualquier otro. La función que puede fallar lo regresa **al final**, y quien la llama
lo revisa en la línea siguiente:

```go
in, err := engine.DecodeInput(r)
if err != nil {
	// algo falló: decidir aquí qué hacer
	return err
}
// a partir de aquí, in es válida
```

`error` es una interfaz con un solo método:

```go
type error interface {
	Error() string
}
```

Cualquier tipo con un método `Error() string` es un error. `nil` significa «no hubo error». El compilador no obliga a
revisarlo, pero `CLAUDE.md` sí: nunca se ignora un error con `_` sin un comentario que diga por qué.

## 2. Errores de categoría

```go
var (
	ErrOutOfRange = errors.New("valor fuera de rango")
	ErrDuplicateID = errors.New("id repetido")
	// …
)
```

`errors.New` crea un valor de error **único**. Guardado en una variable del paquete, con nombre `ErrAlgo` por
convención, sirve como categoría. Quien llama puede decidir según la categoría, por ejemplo así (la respuesta real de
la API se decide en la sesión 1.5):

```go
switch {
case errors.Is(err, engine.ErrMalformedJSON):
	// el JSON no se puede leer
case err != nil:
	// se pudo leer, pero no cumple las reglas
}
```

La lista completa está en `errors.go`; cada una dice cuándo aparece.

## 3. Envolver con `%w`

Un error de categoría no dice **dónde** ni **qué valor**. `fmt.Errorf` con `%w` arma un error nuevo con más contexto y
guarda el original adentro:

```go
return EngineInput{}, fmt.Errorf("%w: %w", ErrMalformedJSON, err)
```

El resultado envuelve dos errores: la categoría `ErrMalformedJSON` y la causa que dio `encoding/json`. Varios `%w` en un
mismo `Errorf` se pueden usar desde Go 1.20.

`%w` y `%v` escriben el mismo texto, pero solo `%w` guarda el error original:

```go
w := fmt.Errorf("leer: %w", ErrA)
v := fmt.Errorf("leer: %v", ErrA)

errors.Is(w, ErrA) // true: w envuelve a ErrA
errors.Is(v, ErrA) // false: v solo copió su texto
```

## 4. `errors.Is` y no comparar el texto

```go
err := in.Validate()
// err.Error() == "la entrada tiene 1 problema:\n- items[0].weightKg: es 0 y debe estar entre 1 y 100000"

errors.Is(err, ErrOutOfRange) // true
err == ErrOutOfRange          // false: err es otro valor que lo envuelve
err.Error() == "valor fuera de rango" // false: el texto trae la ruta y los números
```

`errors.Is(err, objetivo)` revisa `err`, después lo que `err` envuelve, y así hasta el fondo. Responde a la pregunta
«¿este error es, o contiene, esta categoría?».

Comparar el texto falla por tres razones:
- el mensaje trae datos que cambian en cada caso (la ruta, el valor, el límite);
- cambiar una coma del mensaje rompería a quien lo compara;
- el mensaje es para personas, en español; la categoría es para el programa.

La prueba `TestErrorsIsDoesNotDependOnText` lo deja a la vista.

## 5. Tipos de error propios y `errors.As`

Cuando quien recibe el error necesita **datos**, no solo la categoría, se usa un tipo propio:

```go
type FieldError struct {
	Path    string // "items[2].weightKg"
	Message string // "es 0 y debe estar entre 1 y 100000"
	Kind    error  // ErrOutOfRange
}

func (e *FieldError) Error() string { return e.Path + ": " + e.Message }
func (e *FieldError) Unwrap() error { return e.Kind }
```

- `Error()` lo vuelve un error.
- `Unwrap()` dice qué envuelve: su categoría. Por eso `errors.Is(fieldErr, ErrOutOfRange)` da `true`.
- El receptor es `*FieldError`, así que quien cumple la interfaz `error` es el **puntero**. Los errores se crean con
  `&FieldError{...}`.

`errors.As` busca dentro del error uno de cierto **tipo** y lo saca:

```go
var verr *engine.ValidationError
if errors.As(err, &verr) {
	for _, p := range verr.Problems {
		// p.Path dice qué campo marcar en la pantalla
	}
}
```

| | Pregunta | Se le pasa |
|---|---|---|
| `errors.Is(err, ErrOutOfRange)` | ¿Es de esta categoría? | Un **valor** de error |
| `errors.As(err, &verr)` | ¿Hay adentro un error de este **tipo**? Si sí, dámelo | La **dirección** de una variable de ese tipo |

`DecodeInput` envuelve la causa con `%w`, así que también se puede sacar el error original de `encoding/json`:
`TestDecodeInputErrorKeepsTheCause` saca un `*json.SyntaxError`.

## 6. Varios errores a la vez

```go
type ValidationError struct {
	Problems []*FieldError
}

func (e *ValidationError) Unwrap() []error { … }
```

Un error que envuelve **varios** regresa `[]error` en su `Unwrap`. `errors.Is` y `errors.As` recorren la lista
completa: `errors.Is(err, ErrDuplicateID)` encuentra el problema aunque sea el quinto.

La librería estándar trae una forma lista, `errors.Join`:

```go
err := errors.Join(ErrA, fmt.Errorf("contexto: %w", ErrB))
errors.Is(err, ErrB) // true
err.Error()          // "a\ncontexto: b"
```

El motor usa su propio `ValidationError` en lugar de `errors.Join` porque la API necesita la **lista** de
`FieldError` con sus rutas, y un tipo propio la deja a la mano con `errors.As`.

### Por qué todos juntos

Si `Validate` se detuviera en el primer problema, quien captura una carga con cinco errores tendría que corregir y
reenviar cinco veces. Con todos juntos, la pantalla los marca de una vez. `TestValidateReportsAllProblemsTogether`
revisa que lleguen todos y en el orden de la entrada.

La excepción es `DecodeInput`: con un JSON roto no se puede seguir leyendo, así que se detiene en el primer problema de
lectura. Si el JSON se pudo leer, la validación sí reporta todo.

## 7. La trampa del nil con tipo

```go
func bad() error {
	var v *ValidationError // un puntero nil
	return v
}

err := bad()
err == nil // false
```

Un valor de interfaz guarda dos cosas: el **tipo** y el **valor**. `err` guarda «tipo `*ValidationError`, valor
`nil`». Solo es igual a `nil` cuando no guarda nada, ni tipo ni valor. Quien recibe `err` creería que hubo un error,
y al llamar `err.Error()` el programa se detendría, porque el método lee un puntero `nil`.

Por eso `Validate` pregunta y regresa el `nil` de la interfaz de forma explícita:

```go
if len(v.problems) == 0 {
	return nil
}
return &ValidationError{Problems: v.problems}
```

`TestValidateAcceptsValidInput` revisa que `Validate()` dé un `nil` de verdad.

## 8. Lo que `encoding/json` deja pasar

`encoding/json` es tolerante a propósito. Para el contrato, varias de esas tolerancias serían errores silenciosos.
`DecodeInput` primero recorre el JSON token por token (`json.Decoder.Token`) y rechaza lo que se colaría:

| JSON | Qué haría `encoding/json` solo | Qué hace `DecodeInput` |
|---|---|---|
| `"allowedOrientations": null` | Lo deja en `nil`, igual que «no vino»: la pieza tendría las 6 orientaciones | `ErrMalformedJSON`: `null` no se permite |
| `"weightKg": 100, "weightKg": 50000` | Se queda con el último | `ErrMalformedJSON`: llave repetida |
| `"WeightKg": 500` | Lo acepta como `weightKg`: compara sin distinguir mayúsculas, incluso en modo estricto | `ErrMalformedJSON`: no es un campo del contrato |
| Falta `"xMm"` en un tramo | Lo deja en 0, que es válido | `ErrMissingField` |
| `"deliverySequence": 0` | Con `int64` y `omitempty` sería igual a «no vino» | El campo es `*int64`; `Validate` da `ErrOutOfRange` |

Para saber qué llaves tiene cada objeto y cuáles son obligatorias, el recorrido lee las etiquetas `json` de los
structs con el paquete `reflect`: un campo **sin** `omitempty` es obligatorio. Así el contrato no se escribe dos veces.

Al revés, `encoding/json` es más estricto que el esquema en un caso: `1.0` y `1e3` no se aceptan como enteros. Eso se
deja así: nunca deja pasar algo malo.

## 9. Nunca detener el programa

El motor no usa `panic` con datos de entrada: todo problema regresa como `error`. Para comprobarlo, `FuzzDecodeInput`
le da a `DecodeInput` millones de entradas generadas a partir de los ejemplos:

```bash
go test ./backend/internal/engine -run '^$' -fuzz FuzzDecodeInput -fuzztime 30s
```

Con `go test` normal corre solo las semillas. El fuzzing se ve completo en el capítulo 08.

## Para practicar

1. En `DecodeInput`, cambia los dos `%w` de la lectura estricta (después de `decoder.Decode`) por `%v` y corre las
   pruebas. ¿Cuáles fallan y por qué? Regrésalo.
2. Cambia el final de `Validate` por `return &ValidationError{Problems: v.problems}` sin el `if`. ¿Qué prueba falla?
3. En una prueba, valida una entrada con dos errores y saca con `errors.As` el `*ValidationError`. Imprime cada
   `p.Path`.
4. Responde sin ver el capítulo: ¿por qué `errors.Is(err, ErrOutOfRange)` y no `err.Error() == "..."`?
