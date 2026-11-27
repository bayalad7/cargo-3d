package engine

// "slices" es un paquete de la librería estándar con funciones para slices
// (copiar, buscar, ordenar). Clone hace una copia con sus propios elementos.
import "slices"

// Orientation dice cómo queda girada una pieza: las tres letras son las
// medidas de la pieza (L largo, W ancho, H alto) que quedan sobre los ejes del
// equipo, en este orden: 1.ª sobre x, 2.ª sobre z, 3.ª sobre y (vertical).
// La tabla completa está en docs/src/conceptos/sistema-coordenadas.md.
//
// # Nota de aprendizaje: tipos con nombre
//
// "type Orientation string" crea un tipo NUEVO cuyo valor por dentro es un
// string. El tipo NO limita los valores: Orientation("XYZ") compila, y un
// texto escrito tal cual en el código (una constante sin tipo, como "XYZ") se
// convierte solo a Orientation: si f pide un Orientation, f("XYZ") también
// compila. Lo que el compilador rechaza es pasar una VARIABLE de tipo string
// sin convertirla con Orientation(s). Tampoco encoding/json revisa el valor:
// lee cualquier texto. Por eso la validación (sesión 1.1.4) revisa que sea uno
// de los seis, y RotatedDimensions responde "no la conozco" en lugar de fallar.
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

// allOrientations regresa las seis, en el orden del contrato: el valor por
// omisión de allowedOrientations.
//
// Go: un slice no puede ser constante. Si fuera una variable del paquete,
// cualquier archivo de engine podría cambiar sus elementos y cambiaría el valor
// por omisión para todo el programa. Una función que arma una lista nueva en
// cada llamada evita ese riesgo: nadie comparte la lista.
func allOrientations() []Orientation {
	return []Orientation{
		OrientationLWH, OrientationWLH, OrientationLHW,
		OrientationHLW, OrientationWHL, OrientationHWL,
	}
}

// IsValid dice si o es una de las seis orientaciones del contrato. Distingue
// mayúsculas: "lwh" no es válida.
func (o Orientation) IsValid() bool {
	// slices.Contains recorre la lista y dice si o está en ella.
	return slices.Contains(allOrientations(), o)
}

// Dimensions son las medidas de una pieza ya girada, sobre cada eje del
// equipo, en mm. Son las mismas que Placement guarda en dxMm, dyMm y dzMm.
//
// Ojo con el orden: las letras de una orientación van en x, z, y; los campos de
// Dimensions van en x, y, z, como en el contrato. No son el mismo orden.
type Dimensions struct {
	DXMM int64 // medida sobre x (largo del equipo)
	DYMM int64 // medida sobre y (vertical)
	DZMM int64 // medida sobre z (ancho del equipo)
}

// RotatedDimensions devuelve las medidas de una unidad de la pieza girada con
// la orientación o. El segundo valor es false si o no es una de las seis; en ese
// caso Dimensions viene en cero y no debe usarse.
//
// # Cómo funciona
//
// Cada orientación dice qué medida de la pieza (L, W o H) queda sobre x, sobre
// z y sobre y. Por ejemplo, una caja de L = 1200, W = 800, H = 600 en WHL mide
// dx = 800 (W), dz = 600 (H) y dy = 1200 (L): queda parada de punta.
//
// # Por qué así
//
// Una orientación desconocida no detiene el programa: el motor nunca usa panic
// con datos de entrada. La validación (sesión 1.1.4) la rechaza antes; esto es
// la segunda línea de defensa.
//
// # Nota de aprendizaje: "valor, ok"
//
// Una función de Go puede regresar varios valores. Cuando algo puede no
// existir, la costumbre es regresar el valor y un bool llamado ok:
//
//	d, ok := item.RotatedDimensions(o)
//	if !ok {
//		// la orientación no existe
//	}
//
// Es la misma forma que usa Go para leer un mapa (v, ok := m[k]).
func (c CargoItem) RotatedDimensions(o Orientation) (Dimensions, bool) {
	// Nombres cortos para leer la tabla de un vistazo.
	l, w, h := c.LengthMM, c.WidthMM, c.HeightMM

	// "switch" compara o contra cada caso, de arriba abajo, y corre solo el
	// primero que coincide. En Go no hace falta "break": cada caso termina solo.
	// El orden de los campos es DXMM (1.ª letra), DYMM (3.ª), DZMM (2.ª).
	switch o {
	case OrientationLWH:
		return Dimensions{DXMM: l, DYMM: h, DZMM: w}, true
	case OrientationWLH:
		return Dimensions{DXMM: w, DYMM: h, DZMM: l}, true
	case OrientationLHW:
		return Dimensions{DXMM: l, DYMM: w, DZMM: h}, true
	case OrientationHLW:
		return Dimensions{DXMM: h, DYMM: w, DZMM: l}, true
	case OrientationWHL:
		return Dimensions{DXMM: w, DYMM: l, DZMM: h}, true
	case OrientationHWL:
		return Dimensions{DXMM: h, DYMM: l, DZMM: w}, true
	default:
		// Ninguna coincidió: Dimensions{} es el valor cero (todo en 0).
		return Dimensions{}, false
	}
}

// Orientations devuelve las orientaciones permitidas de la pieza: las que
// vinieron en AllowedOrientations o, si no vino la lista (nil), las seis.
//
// Siempre regresa una copia nueva: cambiar lo que regresa no cambia la pieza
// ni el valor por omisión.
//
// # Por qué así
//
// Solo nil significa "no vino" (P-025). Una lista vacía ([]) es un error de la
// entrada que rechaza la validación (sesión 1.1.4); aquí se regresa vacía, tal
// cual, para que nadie la confunda con "todas".
//
// El método no puede llamarse AllowedOrientations: Go no permite que un tipo
// tenga un campo y un método con el mismo nombre.
func (c CargoItem) Orientations() []Orientation {
	if c.AllowedOrientations == nil {
		return allOrientations() // ya es una lista nueva
	}

	// c es una copia de la pieza (receptor por valor), pero la copia comparte
	// los elementos del slice con la original. Sin Clone, quien recibiera el
	// resultado podría cambiar la entrada.
	return slices.Clone(c.AllowedOrientations)
}
