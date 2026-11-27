package engine

import (
	"bytes"         // leer bytes como si fueran un archivo
	"encoding/json" // convertir entre JSON y structs de Go
	"errors"        // comparar errores (errors.Is)
	"io"            // io.EOF: el error que marca el fin de los datos
	"io/fs"         // tipos para recorrer carpetas
	"os"            // leer archivos del disco
	"path/filepath" // armar rutas que funcionan en Windows, macOS y Linux
	"reflect"       // comparar valores de cualquier forma (DeepEqual)
	"testing"
)

// examplesDir es la carpeta de ejemplos del contrato, relativa a esta. Go corre
// las pruebas con la carpeta del paquete como directorio actual
// (backend/internal/engine), así que se sube tres niveles hasta la raíz.
var examplesDir = filepath.Join("..", "..", "..", "contracts", "examples")

// TestInputExamplesRoundTripUnchanged lee cada input.json de
// contracts/examples/ como EngineInput y lo vuelve a escribir como JSON: el
// resultado debe ser el mismo documento (mismos campos y valores; no importan
// el orden de las llaves ni los espacios).
//
// # Qué protege
//
//   - Que cada etiqueta json coincida con el nombre del contrato.
//   - Que no falte ningún campo: se lee en modo estricto y un campo desconocido
//     es falla.
//   - Que lo que no vino regrese sin venir y que un 0 o un false escritos
//     regresen escritos (por eso Stackable, TareKG y Overhang son punteros).
//
// Los output.json (LoadPlan) se prueban en la sesión 1.1.3, cuando exista ese
// tipo. Los ejemplos de invalid/ los usa la validación (sesión 1.1.4).
func TestInputExamplesRoundTripUnchanged(t *testing.T) {
	// Junta las rutas de todos los input.json fuera de invalid/.
	var paths []string

	// filepath.WalkDir recorre la carpeta y llama a la función por cada
	// archivo o subcarpeta. La función recibe la ruta, la entrada y un error.
	//
	// Go: una función se puede escribir en línea y pasar como valor (una
	// "función literal" o closure). Ve y puede modificar las variables de la
	// función que la contiene: por eso agrega a paths directamente.
	err := filepath.WalkDir(examplesDir, func(path string, d fs.DirEntry, err error) error {
		// Errores como valores: si no se pudo leer esta entrada, se regresa el
		// error y WalkDir se detiene. El error es un valor más que se revisa
		// con if (capítulo 04).
		if err != nil {
			return err
		}
		// Los inválidos no son EngineInput correctos: se saltan completos.
		if d.IsDir() && d.Name() == "invalid" {
			return filepath.SkipDir
		}
		if !d.IsDir() && d.Name() == "input.json" {
			paths = append(paths, path) // append regresa el slice con el elemento al final
		}
		return nil // nil = "sin error, sigue"
	})
	if err != nil {
		t.Fatalf("no se pudo recorrer %s: %v", examplesDir, err)
	}

	// Si la ruta estuviera mal, el ciclo no correría y la prueba pasaría sin
	// probar nada. Se exige al menos un ejemplo.
	if len(paths) == 0 {
		t.Fatalf("no se encontró ningún input.json en %s", examplesDir)
	}

	for _, path := range paths {
		// filepath.Rel da la ruta desde examplesDir, para que el nombre de la
		// subprueba sea corto: maritime/dry/20dv-heavy/input.json.
		name, err := filepath.Rel(examplesDir, path)
		if err != nil {
			t.Fatalf("ruta inesperada %s: %v", path, err)
		}

		t.Run(filepath.ToSlash(name), func(t *testing.T) {
			// os.ReadFile lee el archivo completo y regresa (bytes, error):
			// Go devuelve varios valores y el error siempre va al final.
			original, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("no se pudo leer: %v", err)
			}

			// De JSON a Go, en modo estricto: json.Unmarshal ignora en silencio
			// los campos que no conoce; un Decoder con DisallowUnknownFields los
			// convierte en error.
			decoder := json.NewDecoder(bytes.NewReader(original))
			decoder.DisallowUnknownFields()

			var input EngineInput
			// "&input" pasa la dirección de input para que Decode lo llene.
			if err := decoder.Decode(&input); err != nil {
				t.Fatalf("no se pudo leer como EngineInput: %v", err)
			}

			// Decode lee solo el PRIMER valor JSON y deja el resto sin leer: un
			// archivo "{...} basura" pasaría. Un segundo Decode debe encontrar
			// el fin del archivo (io.EOF); cualquier otra cosa es contenido de
			// más. errors.Is compara el error con io.EOF aunque venga envuelto
			// (se explica en el capítulo 04).
			if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
				t.Fatalf("hay contenido después del objeto JSON (err = %v)", err)
			}

			// De Go a JSON otra vez.
			encoded, err := json.Marshal(input)
			if err != nil {
				t.Fatalf("no se pudo escribir como JSON: %v", err)
			}

			// Se comparan como valores genéricos (mapas y listas), así no
			// importan el orden de las llaves ni la sangría.
			want := decodeGeneric(t, original)
			got := decodeGeneric(t, encoded)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("el JSON cambió al ir y volver\n original: %s\n regreso:  %s", compact(t, original), encoded)
			}
		})
	}
}

// decodeGeneric lee JSON sin struct: los objetos quedan como
// map[string]any y las listas como []any. encoding/json convierte TODOS los
// números a float64, que guarda enteros exactos solo hasta 2⁵³
// (≈ 9 × 10¹⁵); los valores del contrato son mucho menores, así que comparar
// así es seguro para estos archivos.
//
// Go: "any" es cualquier tipo. t.Helper() hace que, si algo falla aquí, el
// reporte señale la línea de la prueba que llamó y no esta.
func decodeGeneric(t *testing.T, data []byte) any {
	t.Helper()
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	return v
}

// compact quita espacios y saltos de línea de un JSON para mostrarlo en una
// sola línea en el mensaje de error.
func compact(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := json.Compact(&buf, data); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	return buf.Bytes()
}
