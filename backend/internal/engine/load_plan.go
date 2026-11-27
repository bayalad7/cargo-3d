package engine

// LoadPlan es lo que devuelve el motor: dónde quedó cada unidad, qué no cupo y
// por qué, y las métricas. Corresponde a "LoadPlan" de
// contracts/engine.schema.json.
//
// # Cómo funciona
//
// Se crea con NewLoadPlan, que deja las listas vacías y listas para llenarse,
// y se llena con AddPlacement y AddUnplaced. Las métricas se calculan en el
// sprint 1.3, cuando exista el algoritmo.
//
// # Por qué así
//
// El contrato exige "placements" y "unplaced" siempre presentes como listas.
// En Go, un slice nil se escribe en JSON como null, no como []: un
// LoadPlan{} escrito tal cual no cumpliría el contrato. NewLoadPlan lo evita.
//
// # Nota de aprendizaje: receptor por puntero
//
// Los métodos de LoadPlan tienen receptor por puntero, (p *LoadPlan): reciben
// la dirección del plan, no una copia, y lo que cambian se queda en el plan.
// Con receptor por valor, (p LoadPlan), el append se haría sobre una copia y
// se perdería al regresar. Explicación completa en el capítulo 03.
type LoadPlan struct {
	// SchemaVersion es la versión del contrato con que se generó el plan.
	SchemaVersion string `json:"schemaVersion"`

	// ContainerID es el id del equipo de la entrada.
	ContainerID string `json:"containerId"`

	// Placements tiene una colocación por cada unidad acomodada. Sin omitempty:
	// el contrato la exige aunque esté vacía.
	Placements []Placement `json:"placements"`

	// Unplaced son las unidades que no se acomodaron, agrupadas por pieza y
	// motivo. Sin omitempty, por la misma razón.
	Unplaced []UnplacedItem `json:"unplaced"`

	// Metrics es el resumen del plan.
	Metrics Metrics `json:"metrics"`
}

// NewLoadPlan crea un plan vacío para el equipo containerID, con la versión del
// contrato del motor y las listas vacías (no nil), para que se escriban como []
// en JSON.
//
// # Nota de aprendizaje: funciones New
//
// El lenguaje no tiene constructores. Cuando crear un valor necesita algo más
// que el valor cero, la convención es una función NewTipo que lo regresa listo.
// Aquí el valor cero no sirve: sus slices son nil y saldrían como null.
//
// Regresa LoadPlan (un valor), no *LoadPlan: quien lo recibe decide dónde vive.
// Para usar sus métodos basta guardarlo en una variable: Go toma la dirección
// solo (plan.AddPlacement(...) es (&plan).AddPlacement(...)).
func NewLoadPlan(containerID string) LoadPlan {
	return LoadPlan{
		SchemaVersion: Version,
		ContainerID:   containerID,
		// []Placement{} es un slice vacío pero NO nil: tiene largo 0 y se
		// escribe en JSON como [].
		Placements: []Placement{},
		Unplaced:   []UnplacedItem{},
	}
}

// AddPlacement agrega la colocación de una unidad al plan.
func (p *LoadPlan) AddPlacement(pl Placement) {
	// append regresa el slice con el elemento al final. Hay que guardar lo que
	// regresa: si la lista ya no cabe en su espacio reservado, append crea una
	// nueva en otro lugar de la memoria.
	p.Placements = append(p.Placements, pl)
}

// AddUnplaced registra quantity unidades de la pieza itemID que no se
// acomodaron por reason.
//
// # Cómo funciona
//
// El contrato agrupa por pieza y motivo: si ya hay un grupo con ese itemID y
// ese reason, se le suma quantity; si no, se agrega uno nuevo. Así una pieza
// aparece más de una vez solo si sus unidades tienen motivos distintos.
//
// quantity debe ser al menos 1; con 0 o menos no se registra nada, porque un
// grupo con cero unidades no cumple el contrato.
func (p *LoadPlan) AddUnplaced(itemID string, quantity int64, reason UnplacedReason) {
	if quantity < 1 {
		return
	}

	// "for i := range p.Unplaced" recorre los índices. Se usa el índice, y no
	// "for _, u := range", porque u sería una COPIA del grupo: sumarle a u no
	// cambiaría la lista. p.Unplaced[i] es el elemento mismo.
	for i := range p.Unplaced {
		if p.Unplaced[i].ItemID == itemID && p.Unplaced[i].Reason == reason {
			p.Unplaced[i].Quantity += quantity
			return
		}
	}

	p.Unplaced = append(p.Unplaced, UnplacedItem{ItemID: itemID, Quantity: quantity, Reason: reason})
}

// UnplacedItem son unidades de una pieza que no se acomodaron por un mismo
// motivo. Cada unidad lleva exactamente un motivo, aunque incumpla varios.
type UnplacedItem struct {
	// ItemID es el id de la pieza en la entrada.
	ItemID string `json:"itemId"`
	// Quantity es cuántas unidades no se acomodaron por este motivo.
	Quantity int64 `json:"quantity"`
	// Reason es el motivo.
	Reason UnplacedReason `json:"reason"`
}

// UnplacedReason dice por qué no se acomodó una unidad. Es uno solo: el primero
// que se cumple, en el orden de las constantes de abajo (contrato,
// "UnplacedReason"; ADR-0016).
type UnplacedReason string

// Los motivos del contrato, en el orden en que se revisan.
const (
	// ReasonExceedsDimensions: no cabe en ningún tramo, ni con su margen de
	// overhang, con ninguna orientación permitida.
	ReasonExceedsDimensions UnplacedReason = "EXCEEDS_DIMENSIONS"
	// ReasonExceedsDoor: cabe adentro, pero ninguna orientación que cabe pasa
	// por la puerta.
	ReasonExceedsDoor UnplacedReason = "EXCEEDS_DOOR"
	// ReasonExceedsMaxWeight: una sola unidad pesa más que maxPayloadKg.
	ReasonExceedsMaxWeight UnplacedReason = "EXCEEDS_MAX_WEIGHT"
	// ReasonWeightLimitReached: el peso que queda no alcanza para la unidad.
	ReasonWeightLimitReached UnplacedReason = "WEIGHT_LIMIT_REACHED"
	// ReasonNoSpaceAvailable: queda peso, pero no hay hueco válido.
	ReasonNoSpaceAvailable UnplacedReason = "NO_SPACE_AVAILABLE"
	// ReasonTimeLimitReached: el motor se detuvo por tiempo.
	ReasonTimeLimitReached UnplacedReason = "TIME_LIMIT_REACHED"
)

// Metrics es el resumen del plan. Los conteos, el peso y el volumen son
// enteros; las dos ocupaciones son fracciones (float64), la única parte del
// contrato con decimales.
type Metrics struct {
	// PlacedCount son las unidades acomodadas (largo de Placements).
	PlacedCount int64 `json:"placedCount"`
	// UnplacedCount son las unidades sin acomodar (suma de Unplaced[].Quantity).
	UnplacedCount int64 `json:"unplacedCount"`
	// OutOfGaugeCount son las colocadas con OutOfGauge.
	OutOfGaugeCount int64 `json:"outOfGaugeCount"`
	// LoadedWeightKG es la suma del peso de las unidades acomodadas.
	LoadedWeightKG int64 `json:"loadedWeightKg"`
	// LoadedVolumeMM3 es la suma de Placement.Volume(), en mm³.
	LoadedVolumeMM3 int64 `json:"loadedVolumeMm3"`

	// VolumeUtilization es LoadedVolumeMM3 entre Container.Volume(). Puede
	// pasar de 1 con carga sobredimensionada, que ocupa el margen.
	//
	// Go: float64 es un número de punto flotante de 64 bits. Aquí sí conviene:
	// es una fracción para mostrar, nunca se compara exacta (el validador usa
	// una tolerancia de 0,0001) y no decide si una pieza cabe.
	VolumeUtilization float64 `json:"volumeUtilization"`
	// WeightUtilization es LoadedWeightKG entre MaxPayloadKG, de 0 a 1.
	WeightUtilization float64 `json:"weightUtilization"`
}
