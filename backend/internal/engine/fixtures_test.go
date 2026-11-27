package engine

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixturesDir es engine/fixtures/casos, relativo a esta carpeta.
var fixturesDir = filepath.Join("..", "..", "..", "engine", "fixtures", "casos")

// caseExpectation es el contenido de un archivo .expect.json: qué caso produce
// el fixture y qué resultado se espera del motor (engine/fixtures/README.md).
//
// Solo existe en las pruebas: no es parte del contrato del motor.
type caseExpectation struct {
	Case       string `json:"case"`       // id del caso: "CS-001"
	Doc        string `json:"doc"`        // archivo .md del caso, desde la raíz del repo
	Variant    string `json:"variant"`    // qué equipo, camino o dato usa esta variante
	Phase      string `json:"phase"`      // sesión o fase en que el caso aplica completo
	ValidInput bool   `json:"validInput"` // la entrada debe pasar DecodeInput
	Expected   struct {
		PlacedCount     int64          `json:"placedCount"`
		OutOfGaugeCount int64          `json:"outOfGaugeCount"`
		LoadedWeightKG  int64          `json:"loadedWeightKg"`
		Unplaced        []UnplacedItem `json:"unplaced"`
	} `json:"expected"`
	// Rules son las reglas que no se expresan con números, con la sesión desde
	// la que el motor debe cumplirlas.
	Rules []struct {
		From string `json:"from"`
		Text string `json:"text"`
	} `json:"rules,omitempty"`
}

// TestCaseFixturesAreValidAndConsistent recorre engine/fixtures/casos/ y, para
// cada input*.json, revisa que:
//
//   - pase DecodeInput (contrato y validación de la entrada);
//   - tenga su .expect.json, que se lee en modo estricto;
//   - las cuentas del caso cuadren: colocadas + no colocadas = unidades de la
//     entrada, el peso cargado esperado = peso total − peso de las no
//     colocadas, y no pasa de maxPayloadKg;
//   - cada pieza nombrada en las no colocadas exista;
//   - los motivos que no dependen del algoritmo (no cabe, no pasa por la
//     puerta, pesa más que el máximo) coincidan con la geometría de la entrada.
//
// Las cuentas no dicen DÓNDE va cada pieza (eso lo revisa el algoritmo desde el
// sprint 1.3), pero detectan un caso mal escrito antes de que llegue el
// algoritmo.
func TestCaseFixturesAreValidAndConsistent(t *testing.T) {
	inputs, err := filepath.Glob(filepath.Join(fixturesDir, "cs-*", "input*.json"))
	if err != nil {
		t.Fatalf("no se pudieron listar los fixtures: %v", err)
	}

	count := 0
	for _, inputPath := range inputs {
		// Los .expect.json también empiezan con "input": se saltan aquí.
		if strings.HasSuffix(inputPath, ".expect.json") {
			continue
		}
		count++

		name, err := filepath.Rel(fixturesDir, inputPath)
		if err != nil {
			t.Fatalf("ruta inesperada %s: %v", inputPath, err)
		}

		t.Run(filepath.ToSlash(name), func(t *testing.T) {
			in, err := DecodeInput(bytes.NewReader(readFile(t, inputPath)))
			if err != nil {
				t.Fatalf("el fixture no pasa la validación:\n%v", err)
			}

			// input-camino-b.json → input-camino-b.expect.json
			expectPath := strings.TrimSuffix(inputPath, ".json") + ".expect.json"
			exp := decodeStrict[caseExpectation](t, expectPath)
			checkExpectation(t, inputPath, in, exp)
		})
	}

	// Si la ruta estuviera mal, el ciclo no correría y la prueba pasaría sin
	// revisar nada.
	if count == 0 {
		t.Fatalf("no se encontró ningún fixture en %s", fixturesDir)
	}

	// Cada carpeta de caso debe tener al menos su input.json.
	folders, err := filepath.Glob(filepath.Join(fixturesDir, "cs-*"))
	if err != nil {
		t.Fatalf("no se pudieron listar las carpetas: %v", err)
	}
	for _, folder := range folders {
		if _, err := os.Stat(filepath.Join(folder, "input.json")); err != nil {
			t.Errorf("%s no tiene input.json", filepath.Base(folder))
		}
	}
}

// checkExpectation revisa que el .expect.json corresponda a su entrada.
func checkExpectation(t *testing.T, inputPath string, in EngineInput, exp caseExpectation) {
	t.Helper()

	// La carpeta del fixture es el id del caso en minúsculas: cs-001 ↔ CS-001.
	folder := filepath.Base(filepath.Dir(inputPath))
	if !strings.EqualFold(folder, exp.Case) {
		t.Errorf("el fixture está en %s pero dice ser %s", folder, exp.Case)
	}
	// El documento del caso existe (la ruta es desde la raíz del repo).
	if _, err := os.Stat(filepath.Join("..", "..", "..", exp.Doc)); err != nil {
		t.Errorf("no existe el documento del caso %s: %v", exp.Doc, err)
	}
	if !exp.ValidInput {
		t.Errorf("validInput es false, pero la entrada pasó la validación")
	}

	// Unidades y peso de la entrada, por pieza.
	units := make(map[string]int64)
	weight := make(map[string]int64)
	var totalUnits, totalWeight int64
	for _, it := range in.Items {
		units[it.ID] = it.Quantity
		weight[it.ID] = it.WeightKG
		totalUnits += it.Quantity
		totalWeight += it.Quantity * it.WeightKG
	}

	// Las no colocadas: que la pieza exista, que no pase de sus unidades y que
	// el motivo sea uno del contrato.
	var unplacedUnits, unplacedWeight int64
	perItem := make(map[string]int64)
	for _, u := range exp.Expected.Unplaced {
		if _, ok := units[u.ItemID]; !ok {
			t.Errorf("no colocada %q no es una pieza de la entrada", u.ItemID)
			continue
		}
		if !isKnownReason(u.Reason) {
			t.Errorf("motivo desconocido %q en %s", u.Reason, u.ItemID)
		}
		if u.Quantity < 1 {
			t.Errorf("no colocada %s con cantidad %d", u.ItemID, u.Quantity)
		}
		perItem[u.ItemID] += u.Quantity
		unplacedUnits += u.Quantity
		unplacedWeight += u.Quantity * weight[u.ItemID]
	}
	for id, n := range perItem {
		if n > units[id] {
			t.Errorf("%s: %d no colocadas, pero solo hay %d unidades", id, n, units[id])
		}
	}

	e := exp.Expected
	if e.PlacedCount+unplacedUnits != totalUnits {
		t.Errorf("colocadas (%d) + no colocadas (%d) = %d, pero la entrada tiene %d unidades",
			e.PlacedCount, unplacedUnits, e.PlacedCount+unplacedUnits, totalUnits)
	}
	if want := totalWeight - unplacedWeight; e.LoadedWeightKG != want {
		t.Errorf("peso cargado esperado %d, pero peso total − no colocadas = %d", e.LoadedWeightKG, want)
	}
	if e.LoadedWeightKG > in.Container.MaxPayloadKG {
		t.Errorf("peso cargado %d pasa de maxPayloadKg %d", e.LoadedWeightKG, in.Container.MaxPayloadKG)
	}

	// Sobredimensionadas: no más que las colocadas, y solo si algún tramo
	// tiene margen mayor que cero.
	if e.OutOfGaugeCount > e.PlacedCount {
		t.Errorf("outOfGaugeCount %d pasa de las colocadas %d", e.OutOfGaugeCount, e.PlacedCount)
	}
	if e.OutOfGaugeCount > 0 && !hasMargin(in.Container) {
		t.Errorf("hay %d sobredimensionadas, pero ningún tramo tiene margen", e.OutOfGaugeCount)
	}

	// Los motivos que no dependen del algoritmo (no cabe, no pasa por la
	// puerta, pesa más que el máximo) se pueden calcular con la entrada: se
	// revisa que el .expect.json diga lo mismo.
	for _, it := range in.Items {
		want := hardReason(in.Container, it)
		got := hardReasonsOf(exp.Expected.Unplaced, it.ID)
		switch {
		case want != "" && (len(got) != 1 || got[0].Reason != want || got[0].Quantity != it.Quantity):
			t.Errorf("%s: todas sus %d unidades deberían quedar con %s; el .expect.json dice %+v", it.ID, it.Quantity, want, got)
		case want == "" && len(got) > 0:
			t.Errorf("%s cabe, pasa por la puerta y no pesa de más, pero el .expect.json dice %+v", it.ID, got)
		}
	}

	for _, r := range exp.Rules {
		if r.From == "" || r.Text == "" {
			t.Errorf("regla incompleta: %+v", r)
		}
	}
}

// hardReason calcula el motivo de una pieza que no depende del algoritmo, en
// el orden del contrato (ADR-0016), o "" si la pieza podría colocarse:
//
//   - EXCEEDS_DIMENSIONS: ninguna orientación permitida cabe en ningún tramo,
//     ni con su margen;
//   - EXCEEDS_DOOR: alguna cabe, pero ninguna de las que caben pasa por la
//     puerta (dz ≤ ancho y dy ≤ alto de la puerta);
//   - EXCEEDS_MAX_WEIGHT: una sola unidad pesa más que maxPayloadKg.
//
// Es una versión mínima para revisar los casos; la del motor llega en el
// sprint 1.3.
func hardReason(c Container, it CargoItem) UnplacedReason {
	fitsInside, passesDoor := false, false
	for _, o := range it.Orientations() {
		d, ok := it.RotatedDimensions(o)
		if !ok {
			continue
		}
		for _, s := range c.Sections {
			m := s.Margin()
			if d.DXMM <= s.LengthMM && d.DYMM <= s.HeightMM+m.TopMM && d.DZMM <= s.WidthMM+m.LeftMM+m.RightMM {
				fitsInside = true
				if c.Door == nil || (d.DZMM <= c.Door.WidthMM && d.DYMM <= c.Door.HeightMM) {
					passesDoor = true
				}
			}
		}
	}
	switch {
	case !fitsInside:
		return ReasonExceedsDimensions
	case !passesDoor:
		return ReasonExceedsDoor
	case it.WeightKG > c.MaxPayloadKG:
		return ReasonExceedsMaxWeight
	}
	return ""
}

// hardReasonsOf regresa los grupos de no colocadas de la pieza id cuyo motivo
// no depende del algoritmo.
func hardReasonsOf(unplaced []UnplacedItem, id string) []UnplacedItem {
	var out []UnplacedItem
	for _, u := range unplaced {
		if u.ItemID != id {
			continue
		}
		switch u.Reason {
		case ReasonExceedsDimensions, ReasonExceedsDoor, ReasonExceedsMaxWeight:
			out = append(out, u)
		}
	}
	return out
}

// isKnownReason dice si r es uno de los seis motivos del contrato.
func isKnownReason(r UnplacedReason) bool {
	switch r {
	case ReasonExceedsDimensions, ReasonExceedsDoor, ReasonExceedsMaxWeight,
		ReasonWeightLimitReached, ReasonNoSpaceAvailable, ReasonTimeLimitReached:
		return true
	}
	return false
}

// hasMargin dice si algún tramo del equipo tiene margen de sobredimensión.
func hasMargin(c Container) bool {
	for _, s := range c.Sections {
		if m := s.Margin(); m.TopMM > 0 || m.LeftMM > 0 || m.RightMM > 0 {
			return true
		}
	}
	return false
}

// TestEveryExpectationHasItsInput revisa el lado contrario: ningún
// .expect.json queda sin su input (por ejemplo, después de renombrar uno).
func TestEveryExpectationHasItsInput(t *testing.T) {
	expects, err := filepath.Glob(filepath.Join(fixturesDir, "cs-*", "*.expect.json"))
	if err != nil {
		t.Fatalf("no se pudieron listar los .expect.json: %v", err)
	}
	for _, p := range expects {
		input := strings.TrimSuffix(p, ".expect.json") + ".json"
		if _, err := os.Stat(input); err != nil {
			t.Errorf("%s no tiene su entrada %s", p, filepath.Base(input))
		}
	}
}
