package engine

import "testing"

// TestSectionVolumeIsLengthTimesWidthTimesHeight revisa que el volumen de un
// tramo sea largo × ancho × alto, incluso en los extremos del contrato.
//
// # Nota de aprendizaje: pruebas de tabla
//
// Una "prueba de tabla" es un slice de casos (struct anónimo con entrada y
// resultado esperado) y un ciclo que corre cada uno con t.Run. Agregar un caso
// es agregar una fila, y si uno falla, el reporte dice cuál por su nombre.
//
// # Nota de aprendizaje: new con una expresión
//
// Los campos opcionales son punteros, así que las tablas necesitan un puntero
// a un número: new(int64(400)) regresa un *int64 que apunta a 400. Desde Go
// 1.26, new acepta una expresión (antes solo un tipo, new(int64), que apunta a
// 0) y evita escribir una función auxiliar o una variable para tomar su
// dirección con &. new(int64(0)) significa "vino y vale 0", distinto de nil.
func TestSectionVolumeIsLengthTimesWidthTimesHeight(t *testing.T) {
	// "[]struct{ ... }{ ... }" crea un slice de un struct sin nombre: solo
	// existe dentro de esta prueba.
	tests := []struct {
		name    string
		section Section
		want    int64
	}{
		{
			name:    "interior de un 20' DV",
			section: Section{LengthMM: 5898, WidthMM: 2352, HeightMM: 2393},
			want:    5898 * 2352 * 2393, // 33 195 966 528 mm³ ≈ 33,2 m³
		},
		{
			name:    "tramo de 1 mm por lado",
			section: Section{LengthMM: 1, WidthMM: 1, HeightMM: 1},
			want:    1,
		},
		{
			name:    "máximo del contrato sin desbordar",
			section: Section{LengthMM: 30000, WidthMM: 3000, HeightMM: 4250},
			want:    382_500_000_000, // Go permite "_" para separar miles
		},
		{
			// La posición del tramo no cambia su volumen.
			name:    "la esquina no cuenta",
			section: Section{XMM: 100, YMM: 1478, ZMM: 5, LengthMM: 10, WidthMM: 20, HeightMM: 30},
			want:    6000,
		},
	}

	for _, tt := range tests {
		// t.Run crea una subprueba con nombre: si falla, el reporte dice cuál.
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.section.Volume(); got != tt.want {
				// %d imprime un entero en decimal.
				t.Errorf("Volume() = %d, se esperaba %d", got, tt.want)
			}
		})
	}
}

// TestContainerVolumeSumsSectionsWithoutOverhang revisa que el volumen del
// equipo sea la suma de sus tramos y que el margen de sobredimensión no cuente
// (es el denominador de volumeUtilization).
func TestContainerVolumeSumsSectionsWithoutOverhang(t *testing.T) {
	tests := []struct {
		name      string
		container Container
		want      int64
	}{
		{
			name:      "sin tramos es cero",
			container: Container{},
			want:      0,
		},
		{
			name: "un tramo",
			container: Container{Sections: []Section{
				{LengthMM: 12032, WidthMM: 2352, HeightMM: 2698},
			}},
			want: 12032 * 2352 * 2698,
		},
		{
			name: "el margen de overhang no suma volumen",
			container: Container{Sections: []Section{
				{
					LengthMM: 11652, WidthMM: 2208, HeightMM: 1955,
					Overhang: &Overhang{TopMM: new(int64(400)), LeftMM: new(int64(196)), RightMM: new(int64(196))},
				},
			}},
			want: 11652 * 2208 * 1955,
		},
		{
			// La v0.1.0 acepta un solo tramo, pero Volume ya suma varios para la
			// plataforma escalonada del sprint 2.1.
			name: "dos tramos se suman",
			container: Container{Sections: []Section{
				{LengthMM: 3000, WidthMM: 2500, HeightMM: 2600},
				{XMM: 3000, LengthMM: 13000, WidthMM: 2500, HeightMM: 3000},
			}},
			want: 3000*2500*2600 + 13000*2500*3000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.container.Volume(); got != tt.want {
				t.Errorf("Volume() = %d, se esperaba %d", got, tt.want)
			}
		})
	}
}

// TestSectionMarginDefaultsToZero revisa que un margen que no vino (o un lado
// que no vino) valga 0 al calcular, sin cambiar la entrada.
func TestSectionMarginDefaultsToZero(t *testing.T) {
	tests := []struct {
		name    string
		section Section
		want    OverhangValues
	}{
		{
			name:    "sin overhang no se puede rebasar",
			section: Section{},
			want:    OverhangValues{},
		},
		{
			name:    "overhang vacío es cero en todos lados",
			section: Section{Overhang: &Overhang{}},
			want:    OverhangValues{},
		},
		{
			name:    "solo los lados que vinieron",
			section: Section{Overhang: &Overhang{LeftMM: new(int64(196)), RightMM: new(int64(196))}},
			want:    OverhangValues{LeftMM: 196, RightMM: 196},
		},
		{
			name:    "un cero escrito es cero",
			section: Section{Overhang: &Overhang{TopMM: new(int64(0)), LeftMM: new(int64(150)), RightMM: new(int64(150))}},
			want:    OverhangValues{LeftMM: 150, RightMM: 150},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Dos structs sin punteros ni slices se comparan con == campo por campo.
			if got := tt.section.Margin(); got != tt.want {
				// %+v imprime el struct con los nombres de sus campos.
				t.Errorf("Margin() = %+v, se esperaba %+v", got, tt.want)
			}
		})
	}
}

// TestSectionBoxPutsHeightOnY revisa que la caja de un tramo lleve el largo
// sobre x, el alto sobre y y el ancho sobre z, desde su esquina mínima.
func TestSectionBoxPutsHeightOnY(t *testing.T) {
	s := Section{ID: "main", XMM: 100, YMM: 1200, ZMM: 50, LengthMM: 5898, WidthMM: 2352, HeightMM: 2393}
	want := Box{XMM: 100, YMM: 1200, ZMM: 50, DXMM: 5898, DYMM: 2393, DZMM: 2352}

	if got := s.Box(); got != want {
		t.Errorf("Box() = %+v, se esperaba %+v", got, want)
	}
}
