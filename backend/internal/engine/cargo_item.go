package engine

// Orientation dice cómo queda girada una pieza: las tres letras son las
// medidas de la pieza (L largo, W ancho, H alto) que quedan sobre los ejes del
// equipo, en este orden: 1.ª sobre x, 2.ª sobre z, 3.ª sobre y (vertical).
//
// En esta sesión solo existe el tipo y sus seis valores, para que
// CargoItem.AllowedOrientations tenga tipo. Girar las medidas y resolver las
// orientaciones por omisión se hace en la sesión 1.1.3.
//
// # Nota de aprendizaje: tipos con nombre
//
// "type Orientation string" crea un tipo NUEVO cuyo valor por dentro es un
// string. El tipo NO limita los valores: Orientation("XYZ") compila, y un
// texto escrito tal cual en el código (una constante sin tipo, como "XYZ") se
// convierte solo a Orientation: si f pide un Orientation, f("XYZ") también
// compila. Lo que el compilador rechaza es pasar una VARIABLE de tipo string
// sin convertirla con Orientation(s). Tampoco encoding/json revisa el valor: lee cualquier texto.
// Por eso la validación (sesión 1.1.4) revisa que sea uno de los seis.
type Orientation string

// Las seis orientaciones del contrato.
//
// Go: "const ( ... )" agrupa constantes. Cada una tiene tipo Orientation.
// Las constantes se calculan al compilar y no se pueden cambiar.
const (
	OrientationLWH Orientation = "LWH" // como viene
	OrientationWLH Orientation = "WLH" // girada 90° sobre el eje vertical
	OrientationLHW Orientation = "LHW" // acostada de lado
	OrientationHLW Orientation = "HLW" // acostada de lado y girada
	OrientationWHL Orientation = "WHL" // parada de punta
	OrientationHWL Orientation = "HWL" // parada de punta y girada
)

// CargoItem es un tipo de pieza y cuántas unidades hay. Las medidas y el peso
// son de UNA unidad. «Como viene» (orientación LWH) significa largo sobre x,
// ancho sobre z y alto sobre y.
//
// # Cómo funciona
//
// Se lee del JSON tal cual. Los campos opcionales que tienen valor por omisión
// (allowedOrientations = las 6, stackable = true) quedan vacíos si no vinieron;
// el valor por omisión se resuelve con métodos al calcular (IsStackable y,
// desde la 1.1.3, las orientaciones permitidas).
//
// # Por qué así
//
// Si se llenaran al leer, al volver a JSON aparecerían campos que el usuario no
// escribió y el validador no podría distinguir "dijo que sí" de "no dijo nada".
//
// Redondeo (ADR-0016): las medidas y el peso de la pieza se capturan
// redondeados HACIA ARRIBA (1 199,2 mm → 1 200; 0,4 kg → 1 kg). Es lo contrario
// que el equipo: el equipo se achica y la pieza se agranda, así el redondeo
// nunca deja pasar algo que en realidad no cabe o pesa de más. Lo hace quien
// captura; el motor recibe enteros.
type CargoItem struct {
	// ID identifica la pieza, único dentro de la entrada.
	ID string `json:"id"`

	// Description es texto libre para mostrar. Opcional.
	Description string `json:"description,omitempty"`

	// LengthMM, WidthMM y HeightMM son el largo (L), ancho (W) y alto (H) de
	// una unidad, en mm, redondeados hacia arriba.
	LengthMM int64 `json:"lengthMm"`
	WidthMM  int64 `json:"widthMm"`
	HeightMM int64 `json:"heightMm"`

	// WeightKG es el peso de una unidad en kg, redondeado hacia arriba
	// (mínimo 1 kg). Mil sobres de 0,4 kg cuentan como 1 000 kg: sobrestima
	// a propósito (ver alternativas en ADR-0016).
	WeightKG int64 `json:"weightKg"`

	// Quantity es el número de unidades iguales.
	Quantity int64 `json:"quantity"`

	// AllowedOrientations son las orientaciones permitidas. Si es nil (no
	// vino), se permiten las 6. «Este lado arriba» = [LWH, WLH].
	//
	// Go: un slice nil + omitempty no se escribe al volver a JSON, así que "no
	// vino" regresa como "no vino".
	AllowedOrientations []Orientation `json:"allowedOrientations,omitempty"`

	// Stackable dice si se le puede poner algo encima. Si es nil (no vino),
	// vale true; léelo con IsStackable(). Reservado para el sprint 2.1: en la
	// fase 1 el motor lo acepta pero no lo aplica.
	//
	// Es *bool y no bool porque false es un valor con significado: con bool +
	// omitempty, un "stackable": false desaparecería al volver a JSON (false es
	// el valor cero de bool) y la pieza pasaría a ser apilable. Grave.
	Stackable *bool `json:"stackable,omitempty"`

	// MaxLoadOnTopKG es el peso máximo que soporta encima. 0 es válido (nada
	// encima), por eso es *int64. Reservado para el sprint 2.1.
	MaxLoadOnTopKG *int64 `json:"maxLoadOnTopKg,omitempty"`

	// DeliverySequence es el orden de entrega (1 = se descarga primero, va
	// junto a la puerta). El contrato no permite 0, así que el valor cero de
	// int64 sirve para "no vino" y basta omitempty, sin puntero. Reservado para
	// el sprint 2.1 (LIFO por destino).
	DeliverySequence int64 `json:"deliverySequence,omitempty"`
}

// Volume devuelve el volumen de UNA unidad en mm³ (largo × ancho × alto). Para
// el total del tipo de pieza, multiplica por Quantity.
//
// Solo con entrada validada (sesión 1.1.4): dentro de los límites del contrato
// el máximo es 30 000³ = 2,7 × 10¹³ mm³ por unidad y, aun por las 10 000
// unidades posibles (2,7 × 10¹⁷), cabe en int64 (≈ 9,2 × 10¹⁸). Con medidas
// fuera del contrato se desborda sin aviso: una pieza de 10 000 000 mm por
// lado da un volumen equivocado, no un error.
func (c CargoItem) Volume() int64 {
	return c.LengthMM * c.WidthMM * c.HeightMM
}

// IsStackable dice si se le puede poner algo encima: el valor de Stackable, o
// true si no vino (valor por omisión del contrato).
func (c CargoItem) IsStackable() bool {
	if c.Stackable == nil {
		return true
	}
	// "*c.Stackable" lee el bool al que apunta el puntero.
	return *c.Stackable
}
