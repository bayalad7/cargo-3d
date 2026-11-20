// Los archivos que terminan en _test.go solo se compilan al correr `go test`.
// Viven junto al código que prueban y pertenecen al mismo paquete.
package engine

// "testing" es el paquete estándar de pruebas: Go no necesita PHPUnit ni
// librerías externas para probar.
import "testing"

// TestVersionIsDefined es una prueba de humo: confirma que el paquete compila y
// que las pruebas corren.
//
// Toda prueba en Go es una función que empieza con "Test" y recibe *testing.T,
// el objeto con el que se reportan fallas.
func TestVersionIsDefined(t *testing.T) {
	// En Go no hay "assert": se compara con un if y se reporta la falla con
	// t.Errorf (marca la prueba como fallida y sigue) o t.Fatalf (y se detiene).
	if Version == "" {
		t.Fatalf("la versión del contrato del motor no debe estar vacía")
	}
}
