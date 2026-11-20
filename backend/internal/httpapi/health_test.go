package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest" // servidor y peticiones de mentira para probar sin red
	"testing"
)

// TestHealthReturnsOK verifica que GET /health responde 200 con status "ok".
func TestHealthReturnsOK(t *testing.T) {
	// Creamos una petición falsa y un "grabador" que captura la respuesta.
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// Pasamos la petición por el enrutador real, sin levantar un servidor.
	NewRouter().ServeHTTP(rec, req)

	// Revisamos el código HTTP.
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}

	// Decodificamos el JSON de la respuesta a nuestro struct.
	var body HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("la respuesta no es JSON válido: %v", err)
	}

	if body.Status != "ok" {
		t.Errorf("se esperaba status \"ok\", se obtuvo %q", body.Status)
	}
}

// TestHealthRejectsPost verifica que la ruta solo acepta GET.
// El enrutador estándar responde 405 (método no permitido) por nosotros.
func TestHealthRejectsPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()

	NewRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("se esperaba 405, se obtuvo %d", rec.Code)
	}
}
