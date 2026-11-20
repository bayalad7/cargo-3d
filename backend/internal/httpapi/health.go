// Package httpapi expone el backend por HTTP: rutas, validación de peticiones y
// respuestas JSON. Es la capa "API" del contrato contracts/openapi.yaml.
package httpapi

// Los imports de la librería estándar no necesitan instalarse.
import (
	"encoding/json" // convertir structs de Go a JSON y viceversa
	"net/http"      // servidor HTTP de la librería estándar

	// Importamos nuestro propio paquete del motor por su ruta de módulo.
	"github.com/bayalad7/cargo-3d/backend/internal/engine"
)

// HealthResponse es la respuesta de GET /health.
//
// Las "etiquetas" entre comillas invertidas (`json:"status"`) le dicen a
// encoding/json qué nombre usar en el JSON. Sin ellas usaría "Status" con mayúscula.
type HealthResponse struct {
	Status        string `json:"status"`
	EngineVersion string `json:"engineVersion"`
}

// NewRouter crea el enrutador con todas las rutas de la API.
//
// # Por qué así
//
// Desde Go 1.22, el enrutador de la librería estándar (http.ServeMux) entiende
// métodos y parámetros ("GET /health", "GET /plans/{id}"), así que no
// necesitamos un framework. Ver docs/src/adr/0006-api-con-libreria-estandar.md.
func NewRouter() *http.ServeMux {
	// := declara una variable e infiere su tipo (aquí, *http.ServeMux).
	// El * significa "puntero": una referencia al enrutador, no una copia.
	mux := http.NewServeMux()

	// Registramos la ruta. handleHealth es una función que pasamos como valor.
	mux.HandleFunc("GET /health", handleHealth)

	return mux
}

// handleHealth responde si el servidor está vivo. Docker y el balanceador lo usan
// para saber si el contenedor está sano.
//
// Empieza con minúscula: es privada de este paquete.
// w es donde escribimos la respuesta; r es la petición recibida (aquí no la usamos,
// por eso se nombra _ ).
func handleHealth(w http.ResponseWriter, _ *http.Request) {
	// Armamos la respuesta como un struct de Go.
	resp := HealthResponse{
		Status:        "ok",
		EngineVersion: engine.Version,
	}

	// Avisamos que respondemos JSON.
	w.Header().Set("Content-Type", "application/json")

	// Codificamos el struct a JSON directamente sobre la respuesta.
	// En Go los errores son valores que se devuelven: hay que revisarlos siempre.
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Si falló al escribir, respondemos error 500.
		http.Error(w, "no se pudo generar la respuesta", http.StatusInternalServerError)
	}
}
