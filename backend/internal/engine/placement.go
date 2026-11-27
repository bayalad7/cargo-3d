package engine

// Placement dice dónde quedó UNA unidad de una pieza. Corresponde a "Placement"
// de contracts/engine.schema.json.
//
// # Cómo funciona
//
// (XMM, YMM, ZMM) es la esquina mínima de la pieza ya girada, en coordenadas
// del equipo: x hacia la puerta, y hacia arriba, z hacia la derecha mirando del
// fondo hacia la puerta (ADR-0015, docs/src/conceptos/sistema-coordenadas.md).
// La pieza ocupa de XMM a XMM + DXMM, de YMM a YMM + DYMM y de ZMM a ZMM + DZMM.
//
// # Por qué así
//
// DXMM, DYMM y DZMM repiten lo que ya dice Orientation a propósito: el visor
// no recalcula y el validador revisa que coincidan con
// CargoItem.RotatedDimensions.
type Placement struct {
	// ItemID es el id de la pieza en la entrada.
	ItemID string `json:"itemId"`

	// Instance dice qué unidad de la pieza es, de 0 a Quantity − 1.
	Instance int64 `json:"instance"`

	// SectionID es el tramo donde quedó.
	SectionID string `json:"sectionId"`

	// XMM, YMM y ZMM son la esquina mínima. ZMM puede ser negativa si la pieza
	// rebasa por la izquierda dentro de overhang.leftMm.
	XMM int64 `json:"xMm"`
	YMM int64 `json:"yMm"`
	ZMM int64 `json:"zMm"`

	// DXMM, DYMM y DZMM son las medidas de la pieza girada sobre cada eje.
	DXMM int64 `json:"dxMm"`
	DYMM int64 `json:"dyMm"`
	DZMM int64 `json:"dzMm"`

	// Orientation es cómo quedó girada.
	Orientation Orientation `json:"orientation"`

	// OutOfGauge es true si la pieza rebasa el tramo nominal (usa el margen de
	// overhang) y viaja como carga sobredimensionada.
	//
	// Aquí sí basta un bool con omitempty, sin puntero: el contrato solo admite
	// true ("const": true), así que "no vino" y false significan lo mismo, "la
	// pieza queda dentro del tramo". false es el valor cero de bool y omitempty
	// no lo escribe: regresa sin venir, como pide el contrato. Es el caso
	// contrario de CargoItem.Stackable, donde false sí es un dato.
	OutOfGauge bool `json:"outOfGauge,omitempty"`
}

// Volume devuelve el volumen que ocupa la pieza colocada, en mm³
// (DXMM × DYMM × DZMM). Es lo que suma loadedVolumeMm3 en las métricas.
// Igual que CargoItem.Volume, solo es confiable con medidas dentro del
// contrato.
func (p Placement) Volume() int64 {
	return p.DXMM * p.DYMM * p.DZMM
}

// Dimensions devuelve las medidas de la pieza colocada sobre cada eje, para
// compararlas con CargoItem.RotatedDimensions.
func (p Placement) Dimensions() Dimensions {
	return Dimensions{DXMM: p.DXMM, DYMM: p.DYMM, DZMM: p.DZMM}
}
