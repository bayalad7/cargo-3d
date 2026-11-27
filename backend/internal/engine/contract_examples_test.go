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
// Los ejemplos de invalid/ los usa la validación (sesión 1.1.4).
func TestInputExamplesRoundTripUnchanged(t *testing.T) {
	for _, path := range findExamples(t, "input.json") {
		t.Run(exampleName(t, path), func(t *testing.T) {
			// roundTrip[EngineInput] elige el tipo con que se lee el archivo
			// (ver la nota de genéricos en roundTrip).
			roundTrip[EngineInput](t, path)
		})
	}
}

// TestOutputExamplesRoundTripUnchanged hace lo mismo con cada output.json como
// LoadPlan. En especial revisa outOfGauge: si no viene, debe regresar sin
// venir (false + omitempty), y si viene, regresa true.
func TestOutputExamplesRoundTripUnchanged(t *testing.T) {
	for _, path := range findExamples(t, "output.json") {
		t.Run(exampleName(t, path), func(t *testing.T) {
			roundTrip[LoadPlan](t, path)
		})
	}
}

// TestOutputExamplesDimensionsMatchOrientation junta cada output.json con el
// input.json de su carpeta y revisa que las medidas de cada colocación (dxMm,
// dyMm, dzMm) sean las que da RotatedDimensions para su pieza y su orientación.
//
// Así se comprueba, con todos los ejemplos del contrato, que el código y
// docs/src/conceptos/sistema-coordenadas.md dicen lo mismo.
func TestOutputExamplesDimensionsMatchOrientation(t *testing.T) {
	for _, outputPath := range findExamples(t, "output.json") {
		t.Run(exampleName(t, outputPath), func(t *testing.T) {
			// filepath.Dir quita el nombre del archivo y deja su carpeta.
			inputPath := filepath.Join(filepath.Dir(outputPath), "input.json")
			input := decodeStrict[EngineInput](t, inputPath)
			plan := decodeStrict[LoadPlan](t, outputPath)

			// Un mapa del id de la pieza a la pieza, para buscarla rápido.
			// make crea el mapa vacío y listo para escribir (un mapa nil no
			// acepta escrituras).
			items := make(map[string]CargoItem, len(input.Items))
			for _, item := range input.Items {
				items[item.ID] = item
			}

			for _, pl := range plan.Placements {
				// Leer un mapa regresa (valor, ok): ok es false si la llave no
				// existe.
				item, ok := items[pl.ItemID]
				if !ok {
					t.Errorf("la colocación %s/%d apunta a una pieza que no existe", pl.ItemID, pl.Instance)
					continue // pasa a la siguiente colocación
				}

				want, ok := item.RotatedDimensions(pl.Orientation)
				if !ok {
					t.Errorf("%s/%d: orientación desconocida %q", pl.ItemID, pl.Instance, pl.Orientation)
					continue
				}
				// Dos Dimensions se comparan con == campo por campo.
				if got := pl.Dimensions(); got != want {
					t.Errorf("%s/%d en %s: medidas %+v, la orientación da %+v", pl.ItemID, pl.Instance, pl.Orientation, got, want)
				}
			}
		})
	}
}

// roundTrip lee el archivo como T en modo estricto, lo vuelve a escribir como
// JSON y compara el resultado con el original como valores.
//
// # Nota de aprendizaje: genéricos
//
// "[T any]" es un parámetro de TIPO: roundTrip funciona con cualquier tipo T,
// y quien la llama elige cuál con roundTrip[EngineInput] o roundTrip[LoadPlan].
// El compilador revisa cada uso con su tipo real. Sin genéricos habría que
// escribir la misma función dos veces. Los genéricos tienen su lugar en el
// capítulo 08 (pruebas); aquí solo evitan repetir código.
func roundTrip[T any](t *testing.T, path string) {
	t.Helper()

	value := decodeStrict[T](t, path)

	// De Go a JSON otra vez.
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("no se pudo escribir como JSON: %v", err)
	}

	original := readFile(t, path)

	// Se comparan como valores genéricos (mapas y listas), así no importan el
	// orden de las llaves ni la sangría.
	if !reflect.DeepEqual(decodeGeneric(t, encoded), decodeGeneric(t, original)) {
		t.Errorf("el JSON cambió al ir y volver\n original: %s\n regreso:  %s", compact(t, original), encoded)
	}
}

// decodeStrict lee el archivo como T: falla si trae un campo que T no tiene o
// si queda contenido después del objeto.
func decodeStrict[T any](t *testing.T, path string) T {
	t.Helper()

	// De JSON a Go, en modo estricto: json.Unmarshal ignora en silencio los
	// campos que no conoce; un Decoder con DisallowUnknownFields los convierte
	// en error.
	decoder := json.NewDecoder(bytes.NewReader(readFile(t, path)))
	decoder.DisallowUnknownFields()

	// "var value T" declara una variable del tipo que eligió quien llamó, en su
	// valor cero. "&value" pasa su dirección para que Decode la llene.
	var value T
	if err := decoder.Decode(&value); err != nil {
		t.Fatalf("no se pudo leer %s: %v", path, err)
	}

	// Decode lee solo el PRIMER valor JSON y deja el resto sin leer: un archivo
	// "{...} basura" pasaría. Un segundo Decode debe encontrar el fin del
	// archivo (io.EOF); cualquier otra cosa es contenido de más. errors.Is
	// compara el error con io.EOF aunque venga envuelto (capítulo 04).
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		t.Fatalf("hay contenido después del objeto JSON en %s (err = %v)", path, err)
	}

	return value
}

// findExamples regresa las rutas de todos los archivos llamados fileName en
// contracts/examples/, sin entrar a invalid/. Falla si no encuentra ninguno:
// con una ruta equivocada, la prueba pasaría sin probar nada.
func findExamples(t testing.TB, fileName string) []string {
	t.Helper()

	var paths []string

	// filepath.WalkDir recorre la carpeta y llama a la función por cada archivo
	// o subcarpeta, en orden alfabético.
	//
	// Go: una función se puede escribir en línea y pasar como valor (una
	// "función literal" o closure). Ve y puede modificar las variables de la
	// función que la contiene: por eso agrega a paths directamente.
	err := filepath.WalkDir(examplesDir, func(path string, d fs.DirEntry, err error) error {
		// Errores como valores: si no se pudo leer esta entrada, se regresa el
		// error y WalkDir se detiene.
		if err != nil {
			return err
		}
		// Los inválidos no cumplen el contrato: se saltan completos.
		if d.IsDir() && d.Name() == "invalid" {
			return filepath.SkipDir
		}
		if !d.IsDir() && d.Name() == fileName {
			paths = append(paths, path)
		}
		return nil // nil = "sin error, sigue"
	})
	if err != nil {
		t.Fatalf("no se pudo recorrer %s: %v", examplesDir, err)
	}
	if len(paths) == 0 {
		t.Fatalf("no se encontró ningún %s en %s", fileName, examplesDir)
	}
	return paths
}

// exampleName da la ruta del ejemplo desde contracts/examples/, con /, para que
// el nombre de la subprueba sea corto: maritime/dry/20dv-heavy/input.json.
func exampleName(t *testing.T, path string) string {
	t.Helper()
	name, err := filepath.Rel(examplesDir, path)
	if err != nil {
		t.Fatalf("ruta inesperada %s: %v", path, err)
	}
	return filepath.ToSlash(name)
}

// readFile lee el archivo completo o detiene la prueba.
//
// Go: testing.TB es la interfaz que cumplen *testing.T (pruebas), *testing.B
// (benchmarks) y *testing.F (fuzz). Con ella la misma función auxiliar sirve
// para las tres.
func readFile(t testing.TB, path string) []byte {
	t.Helper()
	// os.ReadFile regresa (bytes, error): el error siempre va al final.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("no se pudo leer %s: %v", path, err)
	}
	return data
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
