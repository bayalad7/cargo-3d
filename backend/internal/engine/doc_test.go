// Los archivos que terminan en _test.go solo se compilan al correr `go test`.
// Viven junto al código que prueban y pertenecen al mismo paquete.
package engine

// "testing" es el paquete estándar de pruebas: con él y `go test` basta, sin
// librerías externas.
import "testing"

// TestVersionIsDefined es una prueba de humo: confirma que el paquete compila y
// que las pruebas corren.
//
// Toda prueba en Go es una función que empieza con "Test" y recibe *testing.T,
// el objeto con el que se reportan fallas.
func TestVersionIsDefined(t *testing.T) {
	// Se compara con un if y la falla se reporta con t.Errorf (marca la prueba
	// como fallida y sigue) o t.Fatalf (y se detiene). Es la forma de Go: la
	// librería estándar no trae funciones de aserción.
	if Version == "" {
		t.Fatalf("la versión del contrato del motor no debe estar vacía")
	}
}
