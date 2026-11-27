package engine

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

// TestNewLoadPlanEncodesEmptyListsAsArrays revisa que un plan recién creado se
// escriba con "placements": [] y "unplaced": [], nunca null: el contrato exige
// listas.
func TestNewLoadPlanEncodesEmptyListsAsArrays(t *testing.T) {
	data, err := json.Marshal(NewLoadPlan("c-20dv"))
	if err != nil {
		t.Fatalf("no se pudo escribir como JSON: %v", err)
	}
	text := string(data)

	for _, want := range []string{`"schemaVersion":"0.1.0"`, `"containerId":"c-20dv"`, `"placements":[]`, `"unplaced":[]`} {
		// strings.Contains dice si un texto contiene a otro.
		if !strings.Contains(text, want) {
			t.Errorf("falta %s en %s", want, text)
		}
	}
	if strings.Contains(text, "null") {
		t.Errorf("el plan vacío trae null: %s", text)
	}
}

// TestZeroLoadPlanEncodesNull documenta por qué existe NewLoadPlan: el valor
// cero de LoadPlan tiene slices nil y encoding/json los escribe como null.
func TestZeroLoadPlanEncodesNull(t *testing.T) {
	data, err := json.Marshal(LoadPlan{})
	if err != nil {
		t.Fatalf("no se pudo escribir como JSON: %v", err)
	}
	if !strings.Contains(string(data), `"placements":null`) {
		t.Errorf("se esperaba placements null en el valor cero: %s", data)
	}
}

// TestAddPlacementKeepsOrder revisa que AddPlacement agregue al final y que el
// cambio se quede en el plan (receptor por puntero).
func TestAddPlacementKeepsOrder(t *testing.T) {
	plan := NewLoadPlan("c-20dv")
	plan.AddPlacement(Placement{ItemID: "a", Instance: 0})
	plan.AddPlacement(Placement{ItemID: "a", Instance: 1})

	if len(plan.Placements) != 2 {
		t.Fatalf("se esperaban 2 colocaciones, hay %d", len(plan.Placements))
	}
	if plan.Placements[1].Instance != 1 {
		t.Errorf("la segunda colocación debía ser la unidad 1: %+v", plan.Placements[1])
	}
}

// TestAddUnplacedGroupsByItemAndReason revisa que las unidades de una misma
// pieza con un mismo motivo queden en un solo grupo, y que un motivo distinto
// abra otro.
func TestAddUnplacedGroupsByItemAndReason(t *testing.T) {
	tests := []struct {
		name string
		adds []UnplacedItem // lo que se registra, en orden
		want []UnplacedItem
	}{
		{
			name: "mismo motivo se suma",
			adds: []UnplacedItem{
				{ItemID: "pallet", Quantity: 1, Reason: ReasonNoSpaceAvailable},
				{ItemID: "pallet", Quantity: 2, Reason: ReasonNoSpaceAvailable},
			},
			want: []UnplacedItem{{ItemID: "pallet", Quantity: 3, Reason: ReasonNoSpaceAvailable}},
		},
		{
			name: "otro motivo abre otro grupo",
			adds: []UnplacedItem{
				{ItemID: "pallet", Quantity: 2, Reason: ReasonNoSpaceAvailable},
				{ItemID: "pallet", Quantity: 3, Reason: ReasonWeightLimitReached},
			},
			want: []UnplacedItem{
				{ItemID: "pallet", Quantity: 2, Reason: ReasonNoSpaceAvailable},
				{ItemID: "pallet", Quantity: 3, Reason: ReasonWeightLimitReached},
			},
		},
		{
			name: "otra pieza abre otro grupo",
			adds: []UnplacedItem{
				{ItemID: "pallet", Quantity: 1, Reason: ReasonExceedsDoor},
				{ItemID: "rack", Quantity: 1, Reason: ReasonExceedsDoor},
			},
			want: []UnplacedItem{
				{ItemID: "pallet", Quantity: 1, Reason: ReasonExceedsDoor},
				{ItemID: "rack", Quantity: 1, Reason: ReasonExceedsDoor},
			},
		},
		{
			name: "cero unidades no se registra",
			adds: []UnplacedItem{{ItemID: "pallet", Quantity: 0, Reason: ReasonExceedsDoor}},
			want: []UnplacedItem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := NewLoadPlan("c-20dv")
			for _, u := range tt.adds {
				plan.AddUnplaced(u.ItemID, u.Quantity, u.Reason)
			}
			if !slices.Equal(plan.Unplaced, tt.want) {
				t.Errorf("Unplaced = %+v, se esperaba %+v", plan.Unplaced, tt.want)
			}
		})
	}
}

// TestPlacementVolumeIsDxDyDz revisa que el volumen de una colocación use las
// medidas ya giradas.
func TestPlacementVolumeIsDxDyDz(t *testing.T) {
	pl := Placement{DXMM: 800, DYMM: 1200, DZMM: 600}
	if got, want := pl.Volume(), int64(800*1200*600); got != want {
		t.Errorf("Volume() = %d, se esperaba %d", got, want)
	}
}
