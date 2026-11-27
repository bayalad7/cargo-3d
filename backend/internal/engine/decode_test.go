package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

// TestDecodeInputAcceptsAllValidExamples revisa que los 27 input.json de
// contracts/examples/ se lean y pasen la validación.
func TestDecodeInputAcceptsAllValidExamples(t *testing.T) {
	for _, path := range findExamples(t, "input.json") {
		t.Run(exampleName(t, path), func(t *testing.T) {
			if _, err := DecodeInput(bytes.NewReader(readFile(t, path))); err != nil {
				t.Errorf("DecodeInput() = %v, se esperaba nil", err)
			}
		})
	}
}

// TestDecodeInputRejectsEachInvalidExample revisa que cada ejemplo de
// contracts/examples/invalid/ dé error y que sea el de SU categoría, no uno
// cualquiera.
func TestDecodeInputRejectsEachInvalidExample(t *testing.T) {
	// El nombre del archivo dice qué rompe cada ejemplo.
	expected := map[string]error{
		"schema/decimal-weight.json":        ErrMalformedJSON,
		"schema/duplicate-orientation.json": ErrDuplicateOrientation,
		"schema/empty-items.json":           ErrNoItems,
		"schema/empty-orientations.json":    ErrNoOrientations,
		"schema/invalid-id.json":            ErrInvalidID,
		"schema/missing-container.json":     ErrMissingField,
		"schema/negative-overhang.json":     ErrOutOfRange,
		"schema/negative-weight.json":       ErrOutOfRange,
		"schema/two-sections.json":          ErrSectionCount,
		"schema/unknown-field.json":         ErrMalformedJSON,
		"schema/unknown-orientation.json":   ErrUnknownOrientation,
		"schema/wrong-schema-version.json":  ErrSchemaVersion,
		"schema/zero-length.json":           ErrOutOfRange,
		"schema/zero-payload.json":          ErrOutOfRange,
		"schema/zero-quantity.json":         ErrOutOfRange,
		"semantic/door-with-overhang.json":  ErrDoorWithOverhang,
		"semantic/duplicate-item-id.json":   ErrDuplicateID,
	}

	// filepath.Glob regresa las rutas que coinciden con el patrón.
	paths, err := filepath.Glob(filepath.Join(examplesDir, "invalid", "*", "*.json"))
	if err != nil {
		t.Fatalf("no se pudieron listar los inválidos: %v", err)
	}
	// Si alguien agrega un ejemplo inválido, también tiene que decir qué error
	// espera: si no, la prueba falla.
	if len(paths) != len(expected) {
		t.Errorf("hay %d ejemplos inválidos y la tabla tiene %d", len(paths), len(expected))
	}

	for _, path := range paths {
		name, err := filepath.Rel(filepath.Join(examplesDir, "invalid"), path)
		if err != nil {
			t.Fatalf("ruta inesperada %s: %v", path, err)
		}
		name = filepath.ToSlash(name)

		t.Run(name, func(t *testing.T) {
			want, ok := expected[name]
			if !ok {
				t.Fatalf("el ejemplo %s no está en la tabla de errores esperados", name)
			}
			_, err := DecodeInput(bytes.NewReader(readFile(t, path)))
			if !errors.Is(err, want) {
				t.Errorf("se esperaba %v, llegó: %v", want, err)
			}
		})
	}
}

// TestDecodeInputRejectsWhatEncodingJSONWouldAccept revisa las trampas de
// encoding/json: cosas que leería sin avisar y que el contrato no permite.
func TestDecodeInputRejectsWhatEncodingJSONWouldAccept(t *testing.T) {
	tests := []struct {
		name  string
		edit  func(doc map[string]any) // cambia el documento antes de escribirlo
		raw   string                   // o un JSON escrito a mano
		kind  error
		inMsg string // texto que debe aparecer en el error
	}{
		{
			name:  "orientaciones en null",
			edit:  func(doc map[string]any) { item(doc)["allowedOrientations"] = nil },
			kind:  ErrMalformedJSON,
			inMsg: "items[0].allowedOrientations",
		},
		{
			name:  "puerta en null",
			edit:  func(doc map[string]any) { doc["container"].(map[string]any)["door"] = nil },
			kind:  ErrMalformedJSON,
			inMsg: "container.door",
		},
		{
			name:  "orden de entrega cero",
			edit:  func(doc map[string]any) { item(doc)["deliverySequence"] = 0 },
			kind:  ErrOutOfRange,
			inMsg: "items[0].deliverySequence",
		},
		{
			name:  "llave con otras mayúsculas",
			edit:  func(doc map[string]any) { delete(item(doc), "weightKg"); item(doc)["WeightKg"] = 500 },
			kind:  ErrMalformedJSON,
			inMsg: "items[0].WeightKg",
		},
		{
			name:  "llave de otro objeto",
			edit:  func(doc map[string]any) { item(doc)["maxPayloadKg"] = 500 },
			kind:  ErrMalformedJSON,
			inMsg: "items[0].maxPayloadKg",
		},
		{
			name:  "falta una esquina del tramo",
			edit:  func(doc map[string]any) { delete(section(doc), "xMm") },
			kind:  ErrMissingField,
			inMsg: "xMm",
		},
		{
			name:  "falta el peso",
			edit:  func(doc map[string]any) { delete(item(doc), "weightKg") },
			kind:  ErrMissingField,
			inMsg: "weightKg",
		},
		{
			name:  "llave repetida",
			raw:   `{"schemaVersion":"0.1.0","schemaVersion":"0.1.0","container":{},"items":[]}`,
			kind:  ErrMalformedJSON,
			inMsg: "repetida",
		},
		{
			name:  "entero escrito con decimal",
			edit:  func(doc map[string]any) { item(doc)["weightKg"] = 500.5 },
			kind:  ErrMalformedJSON,
			inMsg: "weightKg",
		},
		{name: "contenido de más", raw: validJSON(t) + ` {}`, kind: ErrMalformedJSON, inMsg: "después"},
		{name: "vacío", raw: ``, kind: ErrMalformedJSON, inMsg: "vacío"},
		{name: "una lista en lugar de un objeto", raw: `[]`, kind: ErrMalformedJSON},
		{name: "sintaxis rota", raw: `{"schemaVersion":`, kind: ErrMalformedJSON},
		{name: "anidado de más", raw: strings.Repeat("[", 100) + strings.Repeat("]", 100), kind: ErrMalformedJSON, inMsg: "niveles"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := tt.raw
			if tt.edit != nil {
				raw = editedJSON(t, tt.edit)
			}
			_, err := DecodeInput(strings.NewReader(raw))
			if !errors.Is(err, tt.kind) {
				t.Fatalf("se esperaba %v, llegó: %v", tt.kind, err)
			}
			if tt.inMsg != "" && !strings.Contains(err.Error(), tt.inMsg) {
				t.Errorf("el error no menciona %q: %v", tt.inMsg, err)
			}
		})
	}
}

// TestDecodeInputErrorKeepsTheCause revisa que el error de lectura envuelva
// la causa original con %w: errors.As encuentra el *json.SyntaxError.
func TestDecodeInputErrorKeepsTheCause(t *testing.T) {
	_, err := DecodeInput(strings.NewReader(`{"schemaVersion": }`))
	var syntax *json.SyntaxError
	if !errors.As(err, &syntax) {
		t.Fatalf("se esperaba un *json.SyntaxError dentro de: %v (%T)", err, err)
	}
	if !errors.Is(err, ErrMalformedJSON) {
		t.Errorf("el error también debe ser ErrMalformedJSON: %v", err)
	}
}

// FuzzDecodeInput prueba que ninguna entrada, por rara que sea, detenga el
// programa: DecodeInput siempre regresa (entrada, nil) o (vacío, error).
//
// Con `go test` corre solo las semillas (los ejemplos del contrato). Para que
// Go genere entradas nuevas a partir de ellas:
//
//	go test ./backend/internal/engine -run '^$' -fuzz FuzzDecodeInput -fuzztime 30s
//
// El tema completo de fuzzing es del sprint 1.2 (capítulo 08).
func FuzzDecodeInput(f *testing.F) {
	// Semillas: todos los ejemplos, válidos e inválidos, y algunos extremos.
	// findExamples y readFile reciben testing.TB, así que aceptan f.
	for _, path := range findExamples(f, "input.json") {
		f.Add(readFile(f, path))
	}
	invalid, err := filepath.Glob(filepath.Join(examplesDir, "invalid", "*", "*.json"))
	if err != nil {
		f.Fatalf("no se pudieron listar los inválidos: %v", err)
	}
	for _, path := range invalid {
		f.Add(readFile(f, path))
	}
	f.Add([]byte(``))
	f.Add([]byte(`null`))
	f.Add([]byte(`{"items":[{"allowedOrientations":null}]}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		in, err := DecodeInput(bytes.NewReader(data))
		if err != nil {
			return // cualquier error es válido; lo que no puede pasar es un panic
		}
		// Si se aceptó, debe cumplir la validación y poder calcularse sin
		// desbordar: el volumen de cada pieza es positivo.
		if verr := in.Validate(); verr != nil {
			t.Fatalf("DecodeInput aceptó una entrada que Validate rechaza: %v", verr)
		}
		for _, it := range in.Items {
			if it.Volume() <= 0 {
				t.Fatalf("volumen no positivo tras validar: %+v", it)
			}
		}
	})
}

// validJSON escribe validInput() como JSON.
func validJSON(t *testing.T) string {
	t.Helper()
	data, err := json.Marshal(validInput())
	if err != nil {
		t.Fatalf("no se pudo escribir la entrada válida: %v", err)
	}
	return string(data)
}

// editedJSON toma validInput() como documento genérico, le aplica edit y lo
// regresa como JSON. Así cada caso cambia exactamente un campo del texto, como
// lo mandaría un cliente.
func editedJSON(t *testing.T, edit func(doc map[string]any)) string {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal([]byte(validJSON(t)), &doc); err != nil {
		t.Fatalf("no se pudo leer la entrada válida: %v", err)
	}
	edit(doc)
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("no se pudo escribir el documento: %v", err)
	}
	return string(data)
}

// item regresa la primera pieza del documento genérico.
//
// Go: doc["items"].([]any) es una "aserción de tipo": dice "sé que este any
// guarda un []any". Si no fuera así, detendría la prueba con un panic; aquí el
// documento viene de validInput(), así que siempre lo es.
func item(doc map[string]any) map[string]any {
	return doc["items"].([]any)[0].(map[string]any)
}

// section regresa el primer tramo del documento genérico.
func section(doc map[string]any) map[string]any {
	return doc["container"].(map[string]any)["sections"].([]any)[0].(map[string]any)
}
