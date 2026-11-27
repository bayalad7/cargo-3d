package engine

import (
	"slices"
	"testing"
)

// TestRotatedDimensionsForEachOrientation revisa las 6 orientaciones con una
// caja de L = 1200, W = 800 y H = 600, contra la tabla de
// docs/src/conceptos/sistema-coordenadas.md: las letras van en x, z, y; los
// campos, en x, y, z.
func TestRotatedDimensionsForEachOrientation(t *testing.T) {
	box := CargoItem{LengthMM: 1200, WidthMM: 800, HeightMM: 600}

	tests := []struct {
		orientation Orientation
		want        Dimensions
	}{
		{OrientationLWH, Dimensions{DXMM: 1200, DYMM: 600, DZMM: 800}}, // como viene
		{OrientationWLH, Dimensions{DXMM: 800, DYMM: 600, DZMM: 1200}}, // girada sobre el eje vertical
		{OrientationLHW, Dimensions{DXMM: 1200, DYMM: 800, DZMM: 600}}, // acostada de lado
		{OrientationHLW, Dimensions{DXMM: 600, DYMM: 800, DZMM: 1200}}, // acostada de lado y girada
		{OrientationWHL, Dimensions{DXMM: 800, DYMM: 1200, DZMM: 600}}, // parada de punta
		{OrientationHWL, Dimensions{DXMM: 600, DYMM: 1200, DZMM: 800}}, // parada de punta y girada
	}

	for _, tt := range tests {
		// string(tt.orientation) convierte el tipo con nombre a string para
		// usarlo como nombre de la subprueba.
		t.Run(string(tt.orientation), func(t *testing.T) {
			got, ok := box.RotatedDimensions(tt.orientation)
			if !ok {
				t.Fatalf("RotatedDimensions(%s) no reconoce la orientación", tt.orientation)
			}
			if got != tt.want {
				t.Errorf("RotatedDimensions(%s) = %+v, se esperaba %+v", tt.orientation, got, tt.want)
			}

			// Girar no cambia el volumen.
			if v := got.DXMM * got.DYMM * got.DZMM; v != box.Volume() {
				t.Errorf("el volumen girado es %d, se esperaba %d", v, box.Volume())
			}
		})
	}
}

// TestRotatedDimensionsRejectsUnknownOrientation revisa que una orientación que
// no es una de las seis regrese ok = false y medidas en cero, sin detener el
// programa.
func TestRotatedDimensionsRejectsUnknownOrientation(t *testing.T) {
	box := CargoItem{LengthMM: 1200, WidthMM: 800, HeightMM: 600}

	// "lwh" en minúsculas tampoco vale: la comparación distingue mayúsculas.
	for _, o := range []Orientation{"", "XYZ", "lwh"} {
		t.Run("«"+string(o)+"»", func(t *testing.T) {
			got, ok := box.RotatedDimensions(o)
			if ok {
				t.Errorf("RotatedDimensions(%q) dio ok = true", o)
			}
			if got != (Dimensions{}) {
				t.Errorf("RotatedDimensions(%q) = %+v, se esperaba el valor cero", o, got)
			}
		})
	}
}

// TestOrientationsDefaultsToAllSixOnlyWhenNil revisa el valor por omisión de
// las orientaciones permitidas: las seis solo cuando la lista no vino (nil).
// Una lista vacía regresa vacía, porque es un error de la entrada (P-025).
func TestOrientationsDefaultsToAllSixOnlyWhenNil(t *testing.T) {
	tests := []struct {
		name    string
		allowed []Orientation
		want    []Orientation
	}{
		{
			name:    "no vino: las seis",
			allowed: nil,
			want:    []Orientation{OrientationLWH, OrientationWLH, OrientationLHW, OrientationHLW, OrientationWHL, OrientationHWL},
		},
		{
			name:    "este lado arriba",
			allowed: []Orientation{OrientationLWH, OrientationWLH},
			want:    []Orientation{OrientationLWH, OrientationWLH},
		},
		{
			name:    "lista vacía: vacía, no las seis",
			allowed: []Orientation{},
			want:    []Orientation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CargoItem{AllowedOrientations: tt.allowed}.Orientations()

			// slices.Equal compara largo y elementos. No distingue nil de
			// vacío, así que la lista vacía se revisa aparte.
			if !slices.Equal(got, tt.want) {
				t.Errorf("Orientations() = %v, se esperaba %v", got, tt.want)
			}
			if tt.allowed != nil && len(tt.allowed) == 0 && got == nil {
				t.Errorf("una lista vacía regresó nil; debe regresar vacía para que la validación la vea")
			}
		})
	}
}

// TestOrientationsReturnsACopy revisa que cambiar lo que regresa Orientations()
// no cambie la pieza ni el valor por omisión.
func TestOrientationsReturnsACopy(t *testing.T) {
	item := CargoItem{AllowedOrientations: []Orientation{OrientationLWH, OrientationWLH}}
	got := item.Orientations()
	got[0] = OrientationHWL
	if item.AllowedOrientations[0] != OrientationLWH {
		t.Errorf("cambiar el resultado cambió la pieza: %v", item.AllowedOrientations)
	}

	defaults := CargoItem{}.Orientations()
	defaults[0] = OrientationHWL
	if again := (CargoItem{}).Orientations(); again[0] != OrientationLWH {
		t.Errorf("cambiar el resultado cambió el valor por omisión: %v", again)
	}
}
