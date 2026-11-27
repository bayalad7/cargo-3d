package engine

import (
	"errors"
	"strings"
	"testing"
)

// validInput regresa una entrada mínima que cumple el contrato: un 20' DV con
// puerta y una pieza. Cada caso de las tablas la copia y le cambia una cosa.
//
// Es una función, no una variable del paquete: cada llamada arma una entrada
// nueva, así que un caso que cambia Sections o Items no afecta a los demás
// (capítulo 02, «Qué se copia y qué se comparte»).
func validInput() EngineInput {
	return EngineInput{
		SchemaVersion: Version,
		Container: Container{
			ID:           "c-20dv",
			Sections:     []Section{{ID: "main", LengthMM: 5898, WidthMM: 2352, HeightMM: 2393}},
			Door:         &Door{WidthMM: 2340, HeightMM: 2280},
			MaxPayloadKG: 28200,
		},
		Items: []CargoItem{
			{ID: "caja", LengthMM: 1200, WidthMM: 1000, HeightMM: 1000, WeightKG: 500, Quantity: 2},
		},
	}
}

// TestValidateAcceptsValidInput revisa que la entrada mínima pase y que el
// resultado sea un nil de verdad (no un *ValidationError nil dentro de la
// interfaz error; capítulo 04).
func TestValidateAcceptsValidInput(t *testing.T) {
	if err := validInput().Validate(); err != nil {
		t.Fatalf("Validate() = %v, se esperaba nil", err)
	}
}

// TestValidateReportsEachProblem tiene un caso por cada regla: cada uno rompe
// una sola cosa y revisa la categoría (con errors.Is) y la ruta del campo.
func TestValidateReportsEachProblem(t *testing.T) {
	tests := []struct {
		name   string
		change func(in *EngineInput) // qué se rompe
		kind   error                 // categoría esperada
		path   string                // campo esperado
	}{
		{"versión distinta", func(in *EngineInput) { in.SchemaVersion = "0.2.0" }, ErrSchemaVersion, "schemaVersion"},
		{"id de equipo vacío", func(in *EngineInput) { in.Container.ID = "" }, ErrInvalidID, "container.id"},
		{"id con acento", func(in *EngineInput) { in.Container.ID = "camión" }, ErrInvalidID, "container.id"},
		{"id de 65 caracteres", func(in *EngineInput) { in.Container.ID = strings.Repeat("a", 65) }, ErrInvalidID, "container.id"},
		{"descripción de 201 caracteres", func(in *EngineInput) { in.Container.Description = strings.Repeat("ñ", 201) }, ErrTextTooLong, "container.description"},
		{"sin tramos", func(in *EngineInput) { in.Container.Sections = nil }, ErrSectionCount, "container.sections"},
		{"dos tramos", func(in *EngineInput) {
			in.Container.Sections = append(in.Container.Sections, Section{ID: "b", XMM: 6000, LengthMM: 100, WidthMM: 100, HeightMM: 100})
		}, ErrSectionCount, "container.sections"},
		{"tramo de largo cero", func(in *EngineInput) { in.Container.Sections[0].LengthMM = 0 }, ErrOutOfRange, "container.sections[0].lengthMm"},
		{"tramo de ancho negativo", func(in *EngineInput) { in.Container.Sections[0].WidthMM = -1 }, ErrOutOfRange, "container.sections[0].widthMm"},
		{"tramo más alto que la NOM", func(in *EngineInput) { in.Container.Sections[0].HeightMM = 4251 }, ErrOutOfRange, "container.sections[0].heightMm"},
		{"tramo con x negativa", func(in *EngineInput) { in.Container.Sections[0].XMM = -5 }, ErrOutOfRange, "container.sections[0].xMm"},
		{"margen negativo", func(in *EngineInput) {
			in.Container.Door = nil
			in.Container.Sections[0].Overhang = &Overhang{LeftMM: new(int64(-1))}
		}, ErrOutOfRange, "container.sections[0].overhang.leftMm"},
		{"puerta con margen", func(in *EngineInput) {
			in.Container.Sections[0].Overhang = &Overhang{TopMM: new(int64(100))}
		}, ErrDoorWithOverhang, "container.sections[0].overhang"},
		{"puerta de ancho cero", func(in *EngineInput) { in.Container.Door.WidthMM = 0 }, ErrOutOfRange, "container.door.widthMm"},
		{"carga útil cero", func(in *EngineInput) { in.Container.MaxPayloadKG = 0 }, ErrOutOfRange, "container.maxPayloadKg"},
		{"tara negativa", func(in *EngineInput) { in.Container.TareKG = new(int64(-1)) }, ErrOutOfRange, "container.tareKg"},
		{"sin piezas", func(in *EngineInput) { in.Items = nil }, ErrNoItems, "items"},
		{"pieza de alto cero", func(in *EngineInput) { in.Items[0].HeightMM = 0 }, ErrOutOfRange, "items[0].heightMm"},
		{"pieza de más de 30 m", func(in *EngineInput) { in.Items[0].LengthMM = 30001 }, ErrOutOfRange, "items[0].lengthMm"},
		{"peso negativo", func(in *EngineInput) { in.Items[0].WeightKG = -10 }, ErrOutOfRange, "items[0].weightKg"},
		{"cantidad cero", func(in *EngineInput) { in.Items[0].Quantity = 0 }, ErrOutOfRange, "items[0].quantity"},
		{"id de pieza repetido", func(in *EngineInput) { in.Items = append(in.Items, in.Items[0]) }, ErrDuplicateID, "items[1].id"},
		{"orientaciones vacías", func(in *EngineInput) { in.Items[0].AllowedOrientations = []Orientation{} }, ErrNoOrientations, "items[0].allowedOrientations"},
		{"orientación desconocida", func(in *EngineInput) {
			in.Items[0].AllowedOrientations = []Orientation{OrientationLWH, "XYZ"}
		}, ErrUnknownOrientation, "items[0].allowedOrientations[1]"},
		{"orientación repetida", func(in *EngineInput) {
			in.Items[0].AllowedOrientations = []Orientation{OrientationLWH, OrientationLWH}
		}, ErrDuplicateOrientation, "items[0].allowedOrientations[1]"},
		{"peso encima negativo", func(in *EngineInput) { in.Items[0].MaxLoadOnTopKG = new(int64(-1)) }, ErrOutOfRange, "items[0].maxLoadOnTopKg"},
		{"orden de entrega cero", func(in *EngineInput) { in.Items[0].DeliverySequence = new(int64(0)) }, ErrOutOfRange, "items[0].deliverySequence"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := validInput()
			tt.change(&in) // &in: la función recibe la dirección y cambia in

			err := in.Validate()
			if err == nil {
				t.Fatalf("Validate() = nil, se esperaba un error %v", tt.kind)
			}
			// errors.Is busca la categoría en todos los problemas.
			if !errors.Is(err, tt.kind) {
				t.Errorf("el error no es %v: %v", tt.kind, err)
			}
			// errors.As saca el *ValidationError para revisar la ruta.
			var verr *ValidationError
			if !errors.As(err, &verr) {
				t.Fatalf("el error no es un *ValidationError: %T", err)
			}
			if !hasProblem(verr, tt.path, tt.kind) {
				t.Errorf("no se reportó %v en %s:\n%v", tt.kind, tt.path, err)
			}
		})
	}
}

// TestValidateDoesNotRejectWhatIsNotAnError revisa lo que el contrato dice que
// NO es error de validación (ADR-0016 y ADR-0017).
func TestValidateDoesNotRejectWhatIsNotAnError(t *testing.T) {
	tests := []struct {
		name   string
		change func(in *EngineInput)
	}{
		{"puerta más grande que el interior", func(in *EngineInput) { in.Container.Door.HeightMM = 4000 }},
		{"pieza que no cabe en el equipo", func(in *EngineInput) { in.Items[0].LengthMM = 20000 }},
		{"pieza más pesada que la carga útil", func(in *EngineInput) { in.Items[0].WeightKG = 50000 }},
		{"puerta con margen escrito en cero", func(in *EngineInput) {
			in.Container.Sections[0].Overhang = &Overhang{TopMM: new(int64(0)), LeftMM: new(int64(0))}
		}},
		{"margen sin puerta", func(in *EngineInput) {
			in.Container.Door = nil
			in.Container.Sections[0].Overhang = &Overhang{TopMM: new(int64(400)), LeftMM: new(int64(196)), RightMM: new(int64(196))}
		}},
		{"pieza máxima del contrato", func(in *EngineInput) {
			in.Items[0] = CargoItem{ID: "max", LengthMM: 30000, WidthMM: 30000, HeightMM: 30000, WeightKG: 100000, Quantity: 10000}
		}},
		{"descripción de 200 caracteres con acentos", func(in *EngineInput) { in.Items[0].Description = strings.Repeat("ñ", 200) }},
		{"id de 64 caracteres", func(in *EngineInput) { in.Items[0].ID = strings.Repeat("a", 64) }},
		{"tara cero", func(in *EngineInput) { in.Container.TareKG = new(int64(0)) }},
		{"este lado arriba", func(in *EngineInput) {
			in.Items[0].AllowedOrientations = []Orientation{OrientationLWH, OrientationWLH}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := validInput()
			tt.change(&in)
			if err := in.Validate(); err != nil {
				t.Errorf("Validate() = %v, se esperaba nil", err)
			}
		})
	}
}

// TestValidateReportsAllProblemsTogether revisa que varios problemas lleguen
// juntos, en el orden de la entrada, y no solo el primero.
func TestValidateReportsAllProblemsTogether(t *testing.T) {
	in := validInput()
	in.SchemaVersion = "9.9.9"
	in.Container.MaxPayloadKG = 0
	in.Items[0].Quantity = 0
	in.Items[0].AllowedOrientations = []Orientation{}

	var verr *ValidationError
	if !errors.As(in.Validate(), &verr) {
		t.Fatalf("se esperaba un *ValidationError")
	}

	want := []string{"schemaVersion", "container.maxPayloadKg", "items[0].quantity", "items[0].allowedOrientations"}
	if len(verr.Problems) != len(want) {
		t.Fatalf("se esperaban %d problemas, hubo %d:\n%v", len(want), len(verr.Problems), verr)
	}
	for i, p := range verr.Problems {
		if p.Path != want[i] {
			t.Errorf("problema %d en %s, se esperaba %s", i, p.Path, want[i])
		}
	}
}

// TestSectionsOverlapTouchingIsNotOverlap revisa la regla de los tramos: dos
// que se tocan no se enciman; dos que comparten espacio sí.
func TestSectionsOverlapTouchingIsNotOverlap(t *testing.T) {
	a := Section{LengthMM: 3000, WidthMM: 2500, HeightMM: 2600}
	tests := []struct {
		name string
		b    Section
		want bool
	}{
		{"se tocan en x", Section{XMM: 3000, LengthMM: 1000, WidthMM: 2500, HeightMM: 2600}, false},
		{"separados en y", Section{YMM: 2600, LengthMM: 3000, WidthMM: 2500, HeightMM: 100}, false},
		{"comparten 1 mm", Section{XMM: 2999, LengthMM: 1000, WidthMM: 2500, HeightMM: 2600}, true},
		{"uno dentro de otro", Section{XMM: 100, YMM: 100, ZMM: 100, LengthMM: 10, WidthMM: 10, HeightMM: 10}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sectionsOverlap(a, tt.b); got != tt.want {
				t.Errorf("sectionsOverlap = %t, se esperaba %t", got, tt.want)
			}
			// La relación es simétrica: el orden no importa.
			if got := sectionsOverlap(tt.b, a); got != tt.want {
				t.Errorf("al revés, sectionsOverlap = %t, se esperaba %t", got, tt.want)
			}
		})
	}
}

// TestValidateReportsOverlappingSections revisa que dos tramos encimados se
// reporten (además de que la v0.1.0 solo admite uno).
func TestValidateReportsOverlappingSections(t *testing.T) {
	in := validInput()
	in.Container.Door = nil
	in.Container.Sections = append(in.Container.Sections, Section{ID: "b", XMM: 5000, LengthMM: 2000, WidthMM: 2352, HeightMM: 2393})

	err := in.Validate()
	if !errors.Is(err, ErrSectionsOverlap) || !errors.Is(err, ErrSectionCount) {
		t.Errorf("se esperaban ErrSectionsOverlap y ErrSectionCount: %v", err)
	}
}

// TestErrorsIsDoesNotDependOnText muestra por qué la API usa errors.Is y no el
// texto: el mensaje dice cosas concretas (valores, rutas) que cambian en cada
// caso, y la categoría no.
func TestErrorsIsDoesNotDependOnText(t *testing.T) {
	in := validInput()
	in.Items[0].WeightKG = 0
	err := in.Validate()

	if err.Error() == ErrOutOfRange.Error() {
		t.Fatalf("el texto completo no debería ser solo la categoría: %q", err.Error())
	}
	if !errors.Is(err, ErrOutOfRange) {
		t.Errorf("errors.Is debería encontrar ErrOutOfRange dentro de: %v", err)
	}
}

// hasProblem dice si verr trae un problema con esa ruta y esa categoría.
func hasProblem(verr *ValidationError, path string, kind error) bool {
	for _, p := range verr.Problems {
		if p.Path == path && errors.Is(p, kind) {
			return true
		}
	}
	return false
}

// TestValidateDoesNotEchoHugeValues revisa que un valor enorme enviado por el
// usuario se recorte en el mensaje: un id inválido de un millón de caracteres
// no produce un mensaje de un millón de caracteres.
func TestValidateDoesNotEchoHugeValues(t *testing.T) {
	in := validInput()
	in.Items[0].ID = strings.Repeat("é", 1_000_000)

	err := in.Validate()
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("se esperaba ErrInvalidID: %v", err)
	}
	if n := len(err.Error()); n > 500 {
		t.Errorf("el mensaje mide %d bytes; debería recortar el valor", n)
	}
	if !strings.Contains(err.Error(), "…") {
		t.Errorf("el mensaje debería marcar el recorte con «…»: %v", err)
	}
}
