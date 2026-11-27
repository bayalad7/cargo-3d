package engine

// Point es un punto en coordenadas del equipo, en mm (x hacia la puerta, y
// hacia arriba, z hacia la derecha; ADR-0015). Es donde el algoritmo prueba
// poner la esquina mínima de una pieza: los puntos extremos del sprint 1.3
// serán una lista de Point.
type Point struct {
	XMM int64
	YMM int64
	ZMM int64
}

// Candidate es una forma posible de poner una pieza en un punto: con qué
// orientación y qué caja ocuparía.
type Candidate struct {
	// Orientation es cómo queda girada la pieza.
	Orientation Orientation

	// Box es la caja que ocuparía, con su esquina mínima en el punto probado y
	// las medidas de CargoItem.RotatedDimensions con esa orientación.
	Box Box
}

// Candidates devuelve las cajas que puede ocupar una unidad de la pieza con su
// esquina mínima en at, una por cada orientación permitida que queda dentro de
// space. Si dos orientaciones dan la misma caja (piezas cúbicas o con dos lados
// iguales), solo regresa la primera. Si ninguna cabe, regresa nil.
//
// # Cómo funciona
//
//  1. Recorre c.Orientations() en su orden: el del contrato si no vino la
//     lista, o el de allowedOrientations si vino.
//  2. Para cada orientación saca las medidas giradas y arma la caja en at.
//  3. Descarta la caja si se sale de space (Box.ContainedIn) o si ya salió
//     otra con las mismas medidas.
//
// # Por qué así
//
// El resultado es determinista: el orden lo da el slice de orientaciones,
// nunca el mapa de medidas vistas, porque Go recorre un mapa en un orden
// distinto en cada ejecución. El mapa solo responde "¿ya la vi?".
//
// space es una Box y no un Section para que quien llama decida el espacio: el
// tramo nominal (Section.Box) o, en el sprint 1.3, el tramo con su margen de
// rebase. "Este lado arriba" no necesita nada aparte: allowedOrientations =
// ["LWH", "WLH"] deja solo los giros sobre el eje vertical.
//
// Una orientación desconocida se salta sin panic: la validación la rechaza
// antes, y esto es la segunda línea de defensa.
//
// [aprendizaje: slices-maps-y-ciclos] — capítulo 05 de la guía de Go.
func (c CargoItem) Candidates(at Point, space Box) []Candidate {
	orientations := c.Orientations()

	// Go: "var x []T" declara un slice nil: no tiene arreglo de fondo y mide 0.
	// append lo crece solo (aparta memoria con el primer elemento), así que si
	// ninguna orientación cabe se regresa así, nil, sin apartar nada.
	var candidates []Candidate

	// seen guarda las medidas de las candidatas que ya salieron.
	//
	// Se crea en cada llamada; en el sprint 1.3 Candidates se llama por cada
	// punto y cada pieza. Si los benchmarks de la sesión 1.4.4 muestran que
	// cuesta, con a lo más seis candidatas basta revisar la lista una por una
	// (P-030).
	//
	// Go: un mapa con valores bool funciona como conjunto. Dimensions puede ser
	// llave porque todos sus campos son comparables (int64): dos Dimensions con
	// los mismos tres números son la misma llave. make crea el mapa vacío;
	// escribir en un mapa nil (sin make) detiene el programa con panic.
	seen := make(map[Dimensions]bool, len(orientations))

	// Go: "for _, o := range orientations" recorre el slice en orden; en cada
	// vuelta o es una COPIA del elemento. El _ descarta el índice, que aquí no
	// se usa.
	for _, o := range orientations {
		d, ok := c.RotatedDimensions(o)
		if !ok {
			// Go: continue salta a la siguiente vuelta del ciclo.
			continue
		}

		// Leer una llave que no está regresa el valor cero (false): no hace
		// falta preguntar con "v, ok :=" para saber si ya se vio.
		if seen[d] {
			continue
		}

		box := Box{
			XMM: at.XMM, YMM: at.YMM, ZMM: at.ZMM,
			DXMM: d.DXMM, DYMM: d.DYMM, DZMM: d.DZMM,
		}
		if !box.ContainedIn(space) {
			continue
		}

		// Se marca como vista solo si cupo. Si no cupo, tampoco cabría otra
		// orientación con las mismas medidas, así que da igual; marcarla aquí
		// deja claro que seen son "las que ya están en la lista".
		seen[d] = true

		// Go: append regresa el slice que resulta, y hay que guardarlo: si no
		// cabía en la capacidad, append crea un arreglo nuevo y el slice viejo
		// no ve el elemento agregado.
		candidates = append(candidates, Candidate{Orientation: o, Box: box})
	}

	return candidates
}
