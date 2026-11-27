package engine

import (
	"errors"  // crear y comparar errores
	"strconv" // convertir números a texto (Itoa)
	"strings" // armar el mensaje con varias líneas
)

// Errores de categoría de la entrada. La API los distingue con errors.Is para
// decidir qué responder, sin depender del texto del mensaje (ADR-0023).
//
// # Nota de aprendizaje: errores de categoría
//
// errors.New crea un valor de error único. Guardarlo en una variable del
// paquete (por convención con nombre ErrAlgo) permite preguntar después
// "¿este error es de esta categoría?" con errors.Is(err, ErrOutOfRange),
// aunque el error real venga envuelto dentro de otro con más detalle.
// Comparar el texto (err.Error() == "...") se rompería al cambiar una coma.
var (
	// ErrMalformedJSON: el JSON no se puede leer como EngineInput (sintaxis,
	// decimales donde va un entero, campo desconocido, null, llave repetida o
	// con otras mayúsculas, contenido después del objeto).
	ErrMalformedJSON = errors.New("el JSON no se puede leer")
	// ErrMissingField: falta un campo obligatorio del contrato.
	ErrMissingField = errors.New("falta un campo obligatorio")
	// ErrSchemaVersion: schemaVersion no es la del motor (Version).
	ErrSchemaVersion = errors.New("versión del contrato distinta")
	// ErrInvalidID: un id vacío, de más de 64 caracteres o con caracteres fuera
	// de A-Z, a-z, 0-9, punto, guion y guion bajo.
	ErrInvalidID = errors.New("id inválido")
	// ErrDuplicateID: dos piezas o dos tramos con el mismo id.
	ErrDuplicateID = errors.New("id repetido")
	// ErrOutOfRange: un número fuera del mínimo o el máximo del contrato
	// (medidas en cero o negativas, pesos, cantidades, márgenes…).
	ErrOutOfRange = errors.New("valor fuera de rango")
	// ErrTextTooLong: una descripción de más de 200 caracteres.
	ErrTextTooLong = errors.New("texto demasiado largo")
	// ErrNoItems: la lista de piezas está vacía.
	ErrNoItems = errors.New("no hay piezas")
	// ErrSectionCount: el equipo no trae exactamente un tramo (v0.1.0).
	ErrSectionCount = errors.New("número de tramos inválido")
	// ErrSectionsOverlap: dos tramos ocupan el mismo espacio.
	ErrSectionsOverlap = errors.New("tramos encimados")
	// ErrDoorWithOverhang: un equipo con puerta trae margen de
	// sobredimensión mayor que cero; el margen es solo para tramos abiertos
	// (ADR-0017).
	ErrDoorWithOverhang = errors.New("puerta con margen de sobredimensión")
	// ErrNoOrientations: allowedOrientations vino vacía ([]) (P-025).
	ErrNoOrientations = errors.New("sin orientaciones permitidas")
	// ErrUnknownOrientation: una orientación que no es una de las seis.
	ErrUnknownOrientation = errors.New("orientación desconocida")
	// ErrDuplicateOrientation: una orientación repetida en la lista.
	ErrDuplicateOrientation = errors.New("orientación repetida")
)

// FieldError es un problema de un campo de la entrada: dónde está, qué pasa y
// de qué categoría es.
//
// # Cómo funciona
//
// Path usa los nombres del contrato, con índices de lista desde 0
// (items[2].weightKg), para que la API y la pantalla señalen el campo exacto.
// Kind es uno de los errores de categoría de arriba; Unwrap lo regresa, así que
// errors.Is(fieldErr, ErrOutOfRange) da true.
//
// # Nota de aprendizaje: tipos de error propios
//
// En Go, error es una interfaz con un solo método: Error() string. Cualquier
// tipo que tenga ese método es un error. Un tipo propio puede llevar datos
// (aquí Path y Kind) que quien lo recibe saca con errors.As.
type FieldError struct {
	Path    string // campo, con nombres del contrato: "items[2].weightKg"
	Message string // qué pasa, en español, para mostrar
	Kind    error  // categoría: ErrOutOfRange, ErrInvalidID…
}

// Error regresa el problema en una línea: "items[2].weightKg: debe ser …".
//
// Go: el receptor es *FieldError (puntero), así que el que cumple la interfaz
// error es *FieldError, no FieldError. Por eso los errores se crean con
// &FieldError{...}.
func (e *FieldError) Error() string {
	if e.Path == "" {
		return e.Message
	}
	return e.Path + ": " + e.Message
}

// Unwrap regresa la categoría del problema, para que errors.Is y errors.As
// puedan verla.
func (e *FieldError) Unwrap() error {
	return e.Kind
}

// ValidationError junta todos los problemas de una entrada. Validate lo
// regresa cuando hay al menos uno.
//
// # Por qué así
//
// Si se reportara solo el primero, quien captura una carga con cinco errores
// tendría que corregir y reenviar cinco veces. Con todos juntos, la pantalla
// los marca de una vez.
type ValidationError struct {
	// Problems son los problemas, en el orden en que aparecen en la entrada.
	Problems []*FieldError
}

// Error regresa un resumen con un problema por línea.
func (e *ValidationError) Error() string {
	// strings.Builder arma un texto por partes sin crear un string nuevo en
	// cada suma.
	var b strings.Builder
	b.WriteString("la entrada tiene ")
	b.WriteString(strconv.Itoa(len(e.Problems)))
	if len(e.Problems) == 1 {
		b.WriteString(" problema:")
	} else {
		b.WriteString(" problemas:")
	}
	for _, p := range e.Problems {
		b.WriteString("\n- ")
		b.WriteString(p.Error())
	}
	return b.String()
}

// Unwrap regresa cada problema como error, para que errors.Is y errors.As
// revisen todos.
//
// Go: desde Go 1.20 un error puede envolver varios con Unwrap() []error.
// errors.Is(err, ErrDuplicateID) recorre la lista completa y sus categorías.
func (e *ValidationError) Unwrap() []error {
	// make crea un slice de largo 0 y espacio para len(e.Problems) elementos:
	// los append de abajo no tienen que buscar más memoria.
	errs := make([]error, 0, len(e.Problems))
	for _, p := range e.Problems {
		errs = append(errs, p)
	}
	return errs
}
