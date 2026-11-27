package engine

import "testing"

// TestCargoItemVolumeIsForOneUnit revisa que el volumen de la pieza sea el de
// UNA unidad: la cantidad no se multiplica.
func TestCargoItemVolumeIsForOneUnit(t *testing.T) {
	tests := []struct {
		name string
		item CargoItem
		want int64
	}{
		{
			name: "tarima EPAL",
			item: CargoItem{LengthMM: 1200, WidthMM: 800, HeightMM: 1500, Quantity: 1},
			want: 1_440_000_000, // 1,44 m³
		},
		{
			name: "la cantidad no cambia el volumen de una unidad",
			item: CargoItem{LengthMM: 1200, WidthMM: 800, HeightMM: 1500, Quantity: 20},
			want: 1_440_000_000,
		},
		{
			name: "pieza de 1 mm",
			item: CargoItem{LengthMM: 1, WidthMM: 1, HeightMM: 1, Quantity: 1},
			want: 1,
		},
		{
			// 30 000³ = 2,7 × 10¹³; por 10 000 unidades = 2,7 × 10¹⁷, todavía
			// lejos del tope de int64 (≈ 9,2 × 10¹⁸).
			name: "máximo del contrato sin desbordar",
			item: CargoItem{LengthMM: 30000, WidthMM: 30000, HeightMM: 30000, Quantity: 10000},
			want: 27_000_000_000_000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.Volume(); got != tt.want {
				t.Errorf("Volume() = %d, se esperaba %d", got, tt.want)
			}
		})
	}
}

// TestCargoItemVolumeTimesMaxQuantityFitsInInt64 revisa que el volumen total
// del peor caso del contrato (pieza máxima × cantidad máxima) no desborde.
// En Go un int64 que se pasa del tope "da la vuelta" a negativo sin avisar.
func TestCargoItemVolumeTimesMaxQuantityFitsInInt64(t *testing.T) {
	item := CargoItem{LengthMM: 30000, WidthMM: 30000, HeightMM: 30000, Quantity: 10000}

	total := item.Volume() * item.Quantity
	if total != 270_000_000_000_000_000 {
		t.Fatalf("volumen total = %d, se esperaba 2,7 × 10¹⁷ (¿desbordó?)", total)
	}
}

// TestCargoItemIsStackableDefaultsToTrue revisa que una pieza sin "stackable"
// sea apilable y que un false escrito se respete.
func TestCargoItemIsStackableDefaultsToTrue(t *testing.T) {
	tests := []struct {
		name string
		item CargoItem
		want bool
	}{
		{name: "no vino: se puede apilar", item: CargoItem{}, want: true},
		{name: "vino true", item: CargoItem{Stackable: new(true)}, want: true},
		{name: "vino false", item: CargoItem{Stackable: new(false)}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.IsStackable(); got != tt.want {
				// %t imprime un bool.
				t.Errorf("IsStackable() = %t, se esperaba %t", got, tt.want)
			}
		})
	}
}
