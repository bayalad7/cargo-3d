package engine

import (
	"slices"
	"testing"
)

// Interiores de los ejemplos del contrato (contracts/examples/maritime/dry/):
// 20dv-heavy y 40dv-gma-pallets. Sobre x, y, z: largo, alto y ancho.
var (
	space20DV = Box{DXMM: 5898, DYMM: 2393, DZMM: 2352}
	space40DV = Box{DXMM: 12032, DYMM: 2393, DZMM: 2352}
)

// TestCandidatesKeepAllowedOrientationsThatFit revisa qué candidatas salen de
// una pieza en un punto: solo las orientaciones permitidas, solo las que caben
// en el espacio, sin medidas repetidas y en el orden de la lista.
func TestCandidatesKeepAllowedOrientationsThatFit(t *testing.T) {
	// box1200 es una caja de 1 200 × 800 × 600 (L × W × H): sus tres medidas son
	// distintas, así que sus seis orientaciones dan seis cajas distintas.
	box1200 := CargoItem{ID: "box", LengthMM: 1200, WidthMM: 800, HeightMM: 600, WeightKG: 100, Quantity: 1}

	// withOrientations regresa una copia de item con otra lista de
	// orientaciones. item llega por valor, así que cambiarlo aquí no cambia el
	// original.
	withOrientations := func(item CargoItem, os ...Orientation) CargoItem {
		item.AllowedOrientations = os
		return item
	}

	cube := CargoItem{ID: "cube", LengthMM: 1000, WidthMM: 1000, HeightMM: 1000, WeightKG: 100, Quantity: 1}
	pipe := CargoItem{ID: "pipe", LengthMM: 12000, WidthMM: 300, HeightMM: 300, WeightKG: 500, Quantity: 1}

	tests := []struct {
		name  string
		item  CargoItem
		at    Point
		space Box
		// want son las orientaciones esperadas, en orden. nil = ninguna.
		want []Orientation
	}{
		{
			name:  "sin lista y con espacio de sobra: las seis, en el orden del contrato",
			item:  box1200,
			space: space20DV,
			want:  []Orientation{OrientationLWH, OrientationWLH, OrientationLHW, OrientationHLW, OrientationWHL, OrientationHWL},
		},
		{
			// "Este lado arriba": solo se gira sobre el eje vertical, el alto
			// siempre es H (la prueba de abajo lo revisa).
			name:  "este lado arriba: solo los dos giros sobre el eje vertical",
			item:  withOrientations(box1200, OrientationLWH, OrientationWLH),
			space: space20DV,
			want:  []Orientation{OrientationLWH, OrientationWLH},
		},
		{
			name:  "cúbica: las seis orientaciones dan la misma caja, queda una",
			item:  cube,
			space: space20DV,
			want:  []Orientation{OrientationLWH},
		},
		{
			name:  "cúbica con otra lista: queda la primera de la lista",
			item:  withOrientations(cube, OrientationWLH, OrientationLWH),
			space: space20DV,
			want:  []Orientation{OrientationWLH},
		},
		{
			// L = W: LWH y WLH son la misma caja, igual LHW y WHL, igual HLW
			// y HWL. Quedan tres.
			name:  "dos lados iguales: quedan tres",
			item:  CargoItem{ID: "sq", LengthMM: 1200, WidthMM: 1200, HeightMM: 1000, WeightKG: 100, Quantity: 1},
			space: space20DV,
			want:  []Orientation{OrientationLWH, OrientationLHW, OrientationHLW},
		},
		{
			name:  "pieza larga en un 20': no cabe en ninguna orientación",
			item:  pipe,
			space: space20DV,
			want:  nil,
		},
		{
			// Solo cabe acostada a lo largo (12 000 sobre x). LWH y LHW dan la
			// misma caja porque W = H.
			name:  "pieza larga en un 40': una sola, a lo largo",
			item:  pipe,
			space: space40DV,
			want:  []Orientation{OrientationLWH},
		},
		{
			// A 1 000 mm de la puerta solo caben las de dx ≤ 1 000:
			// WLH y WHL (dx = W = 800), HLW y HWL (dx = H = 600).
			name:  "pegada a la puerta: solo las que miden 1 000 o menos sobre x",
			item:  box1200,
			at:    Point{XMM: 5898 - 1000},
			space: space20DV,
			want:  []Orientation{OrientationWLH, OrientationHLW, OrientationWHL, OrientationHWL},
		},
		{
			// Con 700 mm al techo solo caben las de dy ≤ 700: las dos con el
			// alto H = 600 hacia arriba.
			name:  "poco alto libre: solo las que miden 700 o menos sobre y",
			item:  box1200,
			at:    Point{YMM: 2393 - 700},
			space: space20DV,
			want:  []Orientation{OrientationLWH, OrientationWLH},
		},
		{
			name:  "punto fuera del espacio: ninguna",
			item:  box1200,
			at:    Point{XMM: 6000},
			space: space20DV,
			want:  nil,
		},
		{
			name:  "orientación desconocida en la lista: se salta sin panic",
			item:  withOrientations(box1200, Orientation("XYZ"), OrientationLWH),
			space: space20DV,
			want:  []Orientation{OrientationLWH},
		},
		{
			// La validación rechaza una lista vacía antes; aquí no inventa
			// orientaciones.
			//
			// Go: withOrientations(box1200) sin más argumentos NO sirve aquí:
			// un parámetro variádico (os ...Orientation) sin argumentos vale
			// nil, y nil significa "no vino la lista" (las seis). Para una
			// lista vacía hay que pasar un slice vacío con "...".
			name:  "lista vacía: ninguna",
			item:  withOrientations(box1200, []Orientation{}...),
			space: space20DV,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.Candidates(tt.at, tt.space)

			// Go: un slice no se puede comparar con ==; se arma la lista de
			// orientaciones y se compara con slices.Equal.
			var gotOrientations []Orientation
			for _, c := range got {
				gotOrientations = append(gotOrientations, c.Orientation)
			}
			if !slices.Equal(gotOrientations, tt.want) {
				t.Fatalf("orientaciones = %v, se esperaba %v", gotOrientations, tt.want)
			}

			// Sin candidatas tiene que regresar nil, no un slice vacío.
			if tt.want == nil && got != nil {
				t.Errorf("sin candidatas se esperaba nil, salió %#v", got)
			}

			seen := map[Dimensions]bool{}
			for _, c := range got {
				d := Dimensions{DXMM: c.Box.DXMM, DYMM: c.Box.DYMM, DZMM: c.Box.DZMM}

				// La esquina es el punto probado.
				if c.Box.XMM != tt.at.XMM || c.Box.YMM != tt.at.YMM || c.Box.ZMM != tt.at.ZMM {
					t.Errorf("%s: esquina (%d, %d, %d), se esperaba %+v", c.Orientation, c.Box.XMM, c.Box.YMM, c.Box.ZMM, tt.at)
				}
				// Las medidas son las de RotatedDimensions.
				if want, _ := tt.item.RotatedDimensions(c.Orientation); d != want {
					t.Errorf("%s: medidas %+v, se esperaba %+v", c.Orientation, d, want)
				}
				// Queda dentro del espacio.
				if !c.Box.ContainedIn(tt.space) {
					t.Errorf("%s: la caja %+v se sale del espacio", c.Orientation, c.Box)
				}
				// No hay dos candidatas con las mismas medidas.
				if seen[d] {
					t.Errorf("%s: medidas repetidas %+v", c.Orientation, d)
				}
				seen[d] = true
			}
		})
	}
}

// TestThisSideUpKeepsHeightVertical revisa que con "este lado arriba" el alto
// de la pieza (H) siempre quede sobre y, el eje vertical.
func TestThisSideUpKeepsHeightVertical(t *testing.T) {
	item := CargoItem{
		ID: "fragile", LengthMM: 1200, WidthMM: 800, HeightMM: 600, WeightKG: 50, Quantity: 1,
		AllowedOrientations: []Orientation{OrientationLWH, OrientationWLH},
	}

	got := item.Candidates(Point{}, space20DV)
	if len(got) != 2 {
		t.Fatalf("se esperaban 2 candidatas, salieron %d", len(got))
	}
	for _, c := range got {
		if c.Box.DYMM != item.HeightMM {
			t.Errorf("%s: alto %d, se esperaba H = %d", c.Orientation, c.Box.DYMM, item.HeightMM)
		}
	}
}
