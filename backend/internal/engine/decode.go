package engine

import (
	"bytes"         // leer bytes como si fueran un archivo
	"encoding/json" // leer JSON
	"errors"        // comparar errores (errors.Is)
	"fmt"           // envolver errores con %w
	"io"            // io.Reader, io.ReadAll, io.EOF
	"reflect"       // leer las etiquetas json de los structs
	"slices"        // ordenar la lista de campos faltantes
	"strconv"       // convertir índices a texto
	"strings"       // cortar la etiqueta json en la coma
)

// maxJSONDepth es la profundidad máxima de objetos y listas anidados. El
// contrato llega a 5 (raíz → container → sections → tramo → overhang); el
// margen es amplio. Sin límite, un JSON con millones de corchetes anidados
// agotaría la pila de la goroutine en el recorrido recursivo de abajo.
const maxJSONDepth = 16

// DecodeInput lee una entrada del motor en JSON y la valida. Regresa la
// entrada y nil si cumple el contrato, o un error que se clasifica con
// errors.Is:
//
//   - ErrMalformedJSON si el JSON no se puede leer como EngineInput, o
//     ErrMissingField si falta un campo obligatorio. Se detiene en el primer
//     problema: con un JSON roto no se puede seguir leyendo.
//   - un *ValidationError con todos los problemas si se pudo leer, pero no
//     cumple las reglas (ver Validate).
//
// # Cómo funciona
//
//  1. Recorre el JSON pieza por pieza (scanJSON), guiado por los tipos del
//     motor, y rechaza lo que encoding/json aceptaría sin avisar: null, llaves
//     repetidas, llaves con otras mayúsculas ("WeightKg" en lugar de
//     "weightKg"), llaves que no son de ese objeto y campos obligatorios que
//     faltan (un xMm que falta se leería como 0, que es válido). También
//     rechaza el contenido después del objeto.
//  2. Lo lee en modo estricto: un decimal donde va un entero (1.5, 1.0 o 1e3)
//     o una llave en el objeto equivocado es error.
//  3. Llama a Validate.
//
// # Por qué así
//
// Ver ADR-0023. encoding/json está hecho para ser tolerante; el contrato del
// motor es estricto, y lo que encoding/json tolera aquí sería un error
// silencioso. El caso más grave: "allowedOrientations": null se leería como
// "no vino" y la pieza quedaría con las seis orientaciones.
//
// DecodeInput no limita el tamaño de lo que lee: eso lo hace quien llama (la
// API, con un límite de bytes por petición, sesión 1.5).
func DecodeInput(r io.Reader) (EngineInput, error) {
	// io.ReadAll lee todo lo que tenga r. Hacen falta dos lecturas del mismo
	// texto (el recorrido y la lectura estricta), así que se guarda completo.
	data, err := io.ReadAll(r)
	if err != nil {
		// %w envuelve: el error resultante "contiene" a ErrMalformedJSON y a
		// err. errors.Is(resultado, ErrMalformedJSON) da true. Desde Go 1.20 se
		// pueden usar varios %w en un mismo Errorf.
		return EngineInput{}, fmt.Errorf("%w: no se pudo leer la entrada: %w", ErrMalformedJSON, err)
	}

	if err := scanJSON(data); err != nil {
		return EngineInput{}, err
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var in EngineInput
	if err := decoder.Decode(&in); err != nil {
		return EngineInput{}, fmt.Errorf("%w: %w", ErrMalformedJSON, err)
	}

	if err := in.Validate(); err != nil {
		return EngineInput{}, err
	}
	return in, nil
}

// fieldInfo describe un campo de un struct del contrato para el recorrido.
type fieldInfo struct {
	typ      reflect.Type // tipo del campo en Go
	required bool         // true si no tiene omitempty
}

// contractFields regresa los campos de un struct del contrato por su nombre en
// JSON, leídos de las etiquetas json.
//
// Un campo SIN omitempty es obligatorio: en los structs del motor eso
// coincide exactamente con la lista "required" de cada objeto en
// contracts/engine.schema.json. Así el contrato no se escribe dos veces: si un
// campo nuevo llega a CargoItem con su etiqueta, el recorrido ya lo conoce.
//
// # Nota de aprendizaje: reflect
//
// El paquete reflect deja que el programa revise sus propios tipos mientras
// corre: qué campos tiene un struct y qué dicen sus etiquetas. Es lo que usa
// encoding/json por dentro. reflect.TypeFor[EngineInput]() da la descripción
// del tipo EngineInput (con un parámetro de tipo, como las pruebas genéricas
// del capítulo 03).
func contractFields(t reflect.Type) map[string]fieldInfo {
	fields := make(map[string]fieldInfo, t.NumField())
	// "for i := range n" (con n entero) repite n veces, con i de 0 a n − 1.
	for i := range t.NumField() {
		field := t.Field(i)
		// Tag.Get("json") da "lengthMm" o "door,omitempty". strings.Cut corta
		// en la primera coma: lo de antes es el nombre, lo de después las
		// opciones, separadas por comas.
		name, options, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "" || name == "-" {
			continue
		}
		fields[name] = fieldInfo{typ: field.Type, required: !isOptional(options)}
	}
	return fields
}

// isOptional dice si las opciones de una etiqueta json ("omitempty",
// "omitempty,string"…) marcan el campo como opcional: con omitempty u omitzero
// el campo puede no venir.
func isOptional(options string) bool {
	// strings.Split corta en cada coma; slices.ContainsFunc revisa si alguna
	// parte cumple la condición.
	return slices.ContainsFunc(strings.Split(options, ","), func(o string) bool {
		return o == "omitempty" || o == "omitzero"
	})
}

// elemStruct quita punteros y listas hasta llegar al tipo de adentro: de
// *Door da Door; de []Section, Section. Regresa nil si no es un struct.
func elemStruct(t reflect.Type) reflect.Type {
	for t != nil && (t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice) {
		t = t.Elem() // Elem da el tipo al que apunta o el de los elementos
	}
	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}
	return t
}

// scanJSON recorre data token por token, guiado por los tipos de EngineInput,
// y regresa el primer problema que encoding/json dejaría pasar, o nil.
func scanJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := scanValue(decoder, "", 1, reflect.TypeFor[EngineInput]()); err != nil {
		return err
	}
	// Después del objeto solo puede venir el fin del texto.
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return &FieldError{Message: "hay contenido después del objeto JSON", Kind: ErrMalformedJSON}
	}
	return nil
}

// scanValue lee un valor completo (número, texto, objeto o lista) en path y
// revisa lo que tenga adentro. t es el tipo de Go que le corresponde según el
// contrato (nil si no se sabe) y depth la profundidad actual.
//
// # Nota de aprendizaje: json.Decoder.Token
//
// Token() regresa el JSON en pedazos: un json.Delim para cada { } [ ], un
// string para cada texto (incluidas las llaves), float64 para los números,
// bool, o nil para null. Así se ve el JSON tal como vino, antes de convertirlo
// a un struct.
func scanValue(decoder *json.Decoder, path string, depth int, t reflect.Type) error {
	if depth > maxJSONDepth {
		return &FieldError{Path: displayPath(path), Message: fmt.Sprintf("tiene más de %d niveles anidados", maxJSONDepth), Kind: ErrMalformedJSON}
	}

	token, err := decoder.Token()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return &FieldError{Path: displayPath(path), Message: "el JSON está vacío o incompleto", Kind: ErrMalformedJSON}
		}
		return fmt.Errorf("%w: %w", ErrMalformedJSON, err)
	}

	// "switch tok := token.(type)" pregunta de qué tipo es el valor guardado en
	// la interfaz y, en cada caso, tok ya tiene ese tipo.
	switch tok := token.(type) {
	case nil:
		return &FieldError{Path: displayPath(path), Message: "null no está permitido; si el campo es opcional, no lo envíes", Kind: ErrMalformedJSON}
	case json.Delim:
		switch tok {
		case '{':
			return scanObject(decoder, path, depth, elemStruct(t))
		case '[':
			var elem reflect.Type
			if t != nil && t.Kind() == reflect.Slice {
				elem = t.Elem()
			}
			return scanArray(decoder, path, depth, elem)
		}
	}
	// Números, textos y bools no tienen nada adentro. Si el tipo no coincide
	// (un texto donde va un número), lo reporta la lectura estricta.
	return nil
}

// scanObject revisa las llaves de un objeto: que sean de ese objeto en el
// contrato, escritas exactamente igual, que ninguna se repita y que no falte
// ninguna obligatoria. st es el struct que le corresponde (nil si no se sabe;
// la lectura estricta reportará el tipo equivocado).
func scanObject(decoder *json.Decoder, path string, depth int, st reflect.Type) error {
	var fields map[string]fieldInfo
	if st != nil {
		fields = contractFields(st)
	}
	seen := make(map[string]bool)

	// More dice si quedan elementos antes del } que cierra.
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrMalformedJSON, err)
		}
		// Dentro de un objeto, el token antes de cada valor es su llave: un
		// string. "key, ok := token.(string)" saca el string de la interfaz;
		// ok es false si no lo era.
		key, ok := token.(string)
		if !ok {
			return &FieldError{Path: displayPath(path), Message: "llave inesperada", Kind: ErrMalformedJSON}
		}
		keyPath := joinKey(path, key)

		var fieldType reflect.Type
		if fields != nil {
			info, known := fields[key]
			if !known {
				return &FieldError{Path: keyPath, Message: "no es un campo de este objeto en el contrato (los nombres distinguen mayúsculas)", Kind: ErrMalformedJSON}
			}
			fieldType = info.typ
		}
		if seen[key] {
			return &FieldError{Path: keyPath, Message: "la llave está repetida", Kind: ErrMalformedJSON}
		}
		seen[key] = true

		if err := scanValue(decoder, keyPath, depth+1, fieldType); err != nil {
			return err
		}
	}
	// Lee el } que cierra.
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("%w: %w", ErrMalformedJSON, err)
	}

	// Los obligatorios que no vinieron, en orden alfabético para que el
	// mensaje sea siempre el mismo (recorrer un mapa no tiene orden fijo).
	var missing []string
	for name, info := range fields {
		if info.required && !seen[name] {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		slices.Sort(missing)
		return &FieldError{Path: displayPath(path), Message: "faltan campos obligatorios: " + strings.Join(missing, ", "), Kind: ErrMissingField}
	}
	return nil
}

// scanArray revisa cada elemento de una lista. elem es el tipo de los
// elementos (nil si no se sabe).
func scanArray(decoder *json.Decoder, path string, depth int, elem reflect.Type) error {
	for i := 0; decoder.More(); i++ {
		if err := scanValue(decoder, path+"["+strconv.Itoa(i)+"]", depth+1, elem); err != nil {
			return err
		}
	}
	// Lee el ] que cierra.
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("%w: %w", ErrMalformedJSON, err)
	}
	return nil
}

// joinKey arma la ruta de un campo: "items[0]" + "weightKg" = "items[0].weightKg".
func joinKey(path, key string) string {
	if path == "" {
		return key
	}
	return path + "." + key
}

// displayPath regresa path, o "(raíz)" si es la raíz del documento.
func displayPath(path string) string {
	if path == "" {
		return "(raíz)"
	}
	return path
}
