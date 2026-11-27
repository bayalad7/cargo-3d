package engine

import (
	"encoding/json"
	"errors"
	"maps"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// schemaPath es contracts/engine.schema.json, relativo a esta carpeta.
var schemaPath = filepath.Join("..", "..", "..", "contracts", "engine.schema.json")

// schemaDef es lo que estas pruebas leen de cada objeto de $defs: sus
// propiedades (con mínimo y máximo si los tiene) y cuáles son obligatorias.
type schemaDef struct {
	Properties map[string]struct {
		Type    any    `json:"type"` // "integer", "string"…; puede faltar si es $ref
		Minimum *int64 `json:"minimum"`
		Maximum *int64 `json:"maximum"`
	} `json:"properties"`
	Required []string `json:"required"`
}

// schemaDefs lee los objetos de $defs del esquema del motor.
func schemaDefs(t *testing.T) map[string]schemaDef {
	t.Helper()
	var schema struct {
		Defs map[string]schemaDef `json:"$defs"`
	}
	if err := json.Unmarshal(readFile(t, schemaPath), &schema); err != nil {
		t.Fatalf("no se pudo leer el esquema: %v", err)
	}
	return schema.Defs
}

// contractTypes relaciona cada objeto del esquema con su struct de Go.
var contractTypes = map[string]reflect.Type{
	"EngineInput":  reflect.TypeFor[EngineInput](),
	"Container":    reflect.TypeFor[Container](),
	"Section":      reflect.TypeFor[Section](),
	"Overhang":     reflect.TypeFor[Overhang](),
	"Door":         reflect.TypeFor[Door](),
	"CargoItem":    reflect.TypeFor[CargoItem](),
	"LoadPlan":     reflect.TypeFor[LoadPlan](),
	"Placement":    reflect.TypeFor[Placement](),
	"UnplacedItem": reflect.TypeFor[UnplacedItem](),
	"Metrics":      reflect.TypeFor[Metrics](),
}

// TestStructsMatchSchemaFieldsAndRequired revisa que cada struct del motor
// tenga exactamente las propiedades de su objeto en el esquema y que sus
// campos sin omitempty sean exactamente los "required".
//
// DecodeInput depende de esto: usa las etiquetas json para saber qué llaves
// acepta cada objeto y cuáles son obligatorias (ADR-0023). Si alguien agrega
// omitempty a un campo obligatorio, o el esquema agrega un campo, esta prueba
// falla en lugar de que la validación se afloje sin avisar.
func TestStructsMatchSchemaFieldsAndRequired(t *testing.T) {
	defs := schemaDefs(t)

	for name, typ := range contractTypes {
		t.Run(name, func(t *testing.T) {
			def, ok := defs[name]
			if !ok {
				t.Fatalf("el esquema no tiene $defs.%s", name)
			}
			fields := contractFields(typ)

			// maps.Keys da las llaves de un mapa; slices.Sorted las ordena
			// para comparar listas.
			gotNames := slices.Sorted(maps.Keys(fields))
			wantNames := slices.Sorted(maps.Keys(def.Properties))
			if !slices.Equal(gotNames, wantNames) {
				t.Errorf("campos de %s:\n Go:      %v\n esquema: %v", typ, gotNames, wantNames)
			}

			var gotRequired []string
			for fieldName, info := range fields {
				if info.required {
					gotRequired = append(gotRequired, fieldName)
				}
			}
			slices.Sort(gotRequired)
			wantRequired := slices.Sorted(slices.Values(def.Required))
			if !slices.Equal(gotRequired, wantRequired) {
				t.Errorf("obligatorios de %s (sin omitempty):\n Go:      %v\n esquema: %v", typ, gotRequired, wantRequired)
			}
		})
	}
}

// TestValidateRangesMatchSchema revisa que los mínimos y máximos de Validate
// sean los del esquema: para cada propiedad entera de la entrada, el mínimo y
// el máximo pasan, y un valor justo afuera da ErrOutOfRange en esa ruta.
func TestValidateRangesMatchSchema(t *testing.T) {
	defs := schemaDefs(t)

	// Dónde vive cada objeto dentro de la entrada de prueba, como ruta del
	// contrato y como forma de llegar a él en el documento genérico.
	locations := []struct {
		def  string
		path string
		get  func(doc map[string]any) map[string]any
	}{
		{"EngineInput", "", func(doc map[string]any) map[string]any { return doc }},
		{"Container", "container", func(doc map[string]any) map[string]any { return doc["container"].(map[string]any) }},
		{"Door", "container.door", func(doc map[string]any) map[string]any {
			return doc["container"].(map[string]any)["door"].(map[string]any)
		}},
		{"Section", "container.sections[0]", section},
		{"Overhang", "container.sections[0].overhang", func(doc map[string]any) map[string]any {
			// El margen solo es válido sin puerta (ADR-0017).
			delete(doc["container"].(map[string]any), "door")
			overhang := map[string]any{}
			section(doc)["overhang"] = overhang
			return overhang
		}},
		{"CargoItem", "items[0]", item},
	}

	for _, loc := range locations {
		for prop, spec := range defs[loc.def].Properties {
			if spec.Type != "integer" {
				continue // solo números enteros
			}
			path := joinKey(loc.path, prop)

			// Cada caso: el valor y si debe dar ErrOutOfRange.
			type rangeCase struct {
				value int64
				bad   bool
			}
			var cases []rangeCase
			if spec.Minimum != nil {
				cases = append(cases, rangeCase{*spec.Minimum, false}, rangeCase{*spec.Minimum - 1, true})
			}
			if spec.Maximum != nil {
				cases = append(cases, rangeCase{*spec.Maximum, false}, rangeCase{*spec.Maximum + 1, true})
			}

			for _, c := range cases {
				t.Run(path+"="+strconv.FormatInt(c.value, 10), func(t *testing.T) {
					raw := editedJSON(t, func(doc map[string]any) { loc.get(doc)[prop] = c.value })
					_, err := DecodeInput(strings.NewReader(raw))

					var verr *ValidationError
					reported := errors.As(err, &verr) && hasProblem(verr, path, ErrOutOfRange)
					if c.bad && !reported {
						t.Errorf("%s = %d está fuera del esquema y no se reportó: %v", path, c.value, err)
					}
					if !c.bad && reported {
						t.Errorf("%s = %d está dentro del esquema y se reportó: %v", path, c.value, err)
					}
				})
			}
		}
	}
}
