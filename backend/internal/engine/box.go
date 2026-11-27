package engine

// Box es una caja alineada a los ejes del equipo: no puede quedar girada en
// diagonal, solo acostada, de lado o de punta (las seis orientaciones). Es la
// pieza básica de la geometría del motor: con ella se revisa si dos piezas se
// enciman, si una pieza cabe en el tramo y cuánto se apoya una sobre otra.
//
// # Cómo funciona
//
// (XMM, YMM, ZMM) es la esquina mínima, en coordenadas del equipo (x hacia la
// puerta, y hacia arriba, z hacia la derecha; ADR-0015). DXMM, DYMM y DZMM son
// sus medidas sobre cada eje. La caja ocupa el rango semiabierto
// [XMM, XMM + DXMM) sobre x, y lo mismo sobre y y z: el inicio es suyo, el
// final ya no. Por eso dos cajas que se tocan por una cara no se enciman.
// Los dibujos están en docs/src/conceptos/geometria.md.
//
// # Por qué así
//
// Los nombres y el orden de los campos son los de Placement y Dimensions, para
// pasar de uno a otro sin traducir ejes (Placement.Box). Todo es int64 en mm,
// como el resto del motor: con enteros, "se tocan" es una igualdad exacta y no
// hay errores de redondeo.
//
// Una caja con alguna medida en cero o negativa no tiene volumen: nunca se
// encima con nada, no se apoya en nada y no cuenta como contenida en otra. La
// validación de la entrada ya impide medidas así; esta regla es la segunda línea
// de defensa por si el algoritmo arma una caja mal.
//
// # Nota de aprendizaje: el valor cero de un struct
//
// Box{} es una caja válida para el compilador: todos sus campos valen 0, que es
// el valor cero de int64. Go no obliga a llenar los campos, así que una caja
// "olvidada" llega con medidas en cero. Es justo el caso que cubre la regla de
// arriba: Box{} no tiene volumen y ninguna función la trata como una caja real.
type Box struct {
	// XMM, YMM y ZMM son la esquina mínima. Pueden ser negativas: una pieza que
	// rebasa por la izquierda (overhang.leftMm) empieza en z < 0.
	XMM int64
	YMM int64
	ZMM int64

	// DXMM, DYMM y DZMM son las medidas sobre x, y y z.
	DXMM int64
	DYMM int64
	DZMM int64
}

// maxX, maxY y maxZ devuelven dónde termina la caja sobre cada eje: el primer
// milímetro que ya NO es suyo.
//
// Go: empiezan con minúscula, así que no se exportan; solo el paquete engine
// las ve. Son detalle interno de Box.
//
// Go: la suma de int64 no avisa si se pasa del máximo (≈ 9,2 × 10¹⁸): da la
// vuelta y sale un número negativo, sin error ni panic. Con una caja que pasa
// InRange (todos sus números de a lo más MaxGeometryMM en valor absoluto) la
// suma queda muy lejos de ese tope. El contrato todavía no pone máximo a las
// coordenadas de una colocación: un plan que llegue de fuera con xMm cerca del
// tope haría que maxX dé la vuelta e Intersects responda "no se enciman" aunque
// sí se enciman. Por eso, con datos de fuera, hay que revisar InRange ANTES de
// usar la caja (P-029; el máximo del contrato entra con el validador, sesión
// 1.4.2).
func (b Box) maxX() int64 { return b.XMM + b.DXMM }
func (b Box) maxY() int64 { return b.YMM + b.DYMM }
func (b Box) maxZ() int64 { return b.ZMM + b.DZMM }

// MaxGeometryMM es el valor absoluto más grande que puede tener cualquiera de
// los seis números de una Box para que la geometría sea confiable: 100 000 mm
// (100 m).
//
// # Por qué este número
//
// Cubre con holgura todo lo que permite el contrato: una pieza mide a lo más
// 30 000 mm por lado, una coordenada de un plan que arma el motor llega a lo
// más a x = 60 000 (un tramo que empieza en 30 000 y mide 30 000) y a lo menos a
// z = −1 000 (margen de rebase por la izquierda). Con 100 000, ninguna suma
// (a lo más 2 × 10⁵) ni área (a lo más 4 × 10¹⁰ mm²) se acerca al tope de int64
// (≈ 9,2 × 10¹⁸), así que no hay desbordamiento.
const MaxGeometryMM int64 = 100_000

// InRange dice si los seis números de la caja tienen un valor absoluto de a lo
// más MaxGeometryMM. Con una caja que no cumple, Intersects, ContainedIn y
// Contact pueden dar un resultado equivocado por desbordamiento de int64.
//
// # Cómo funciona
//
// Revisa los seis números, sean medidas positivas, cero o negativas: una caja
// sin volumen pero de números chicos sí está en rango (las tres funciones ya
// saben qué hacer con ella; ver Box).
//
// # Por qué así
//
// El motor arma sus propias cajas dentro de estos límites. Un plan que llega de
// fuera (por ejemplo, tras mover una pieza a mano en el visor) puede traer
// cualquier número: quien lo recibe llama a InRange antes de usar la caja.
//
// # Nota de aprendizaje: valor absoluto sin math.Abs
//
// math.Abs trabaja con float64. Con enteros se compara contra los dos lados del
// rango (-MaxGeometryMM <= v && v <= MaxGeometryMM). No se escribe -v para sacar
// el valor absoluto: el negativo de math.MinInt64 no cabe en int64 y vuelve a
// salir math.MinInt64, así que "-v <= max" sería verdadero para el número más
// negativo que existe.
func (b Box) InRange() bool {
	// Go: una función anónima asignada a una variable (un "closure"). inRange
	// se usa seis veces abajo y dice una sola vez cuál es el rango.
	inRange := func(v int64) bool {
		return -MaxGeometryMM <= v && v <= MaxGeometryMM
	}
	return inRange(b.XMM) && inRange(b.YMM) && inRange(b.ZMM) &&
		inRange(b.DXMM) && inRange(b.DYMM) && inRange(b.DZMM)
}

// hasVolume dice si las tres medidas son positivas. Una caja plana (una medida
// en 0) o con una medida negativa no ocupa espacio.
func (b Box) hasVolume() bool {
	return b.DXMM > 0 && b.DYMM > 0 && b.DZMM > 0
}

// overlap devuelve cuántos milímetros comparten dos rangos semiabiertos
// [aMin, aMax) y [bMin, bMax) sobre un eje. Si no comparten nada (están
// separados o solo se tocan en un punto), devuelve 0.
//
// # Cómo funciona
//
// La parte común empieza en el mayor de los dos inicios y termina en el menor
// de los dos finales. Si ese final no pasa del inicio, no hay parte común.
//
//	a:   [=======)
//	b:        [=========)
//	          ^ max(inicios)
//	             ^ min(finales)
//	común:    [==)           = min(aMax, bMax) − max(aMin, bMin)
//
// # Nota de aprendizaje: min y max
//
// min y max son funciones integradas del lenguaje desde Go 1.21: no hay que
// importar nada. Aceptan dos o más valores del mismo tipo ordenable (enteros,
// flotantes, strings) y devuelven el menor o el mayor.
func overlap(aMin, aMax, bMin, bMax int64) int64 {
	// Go: "aMin, aMax, bMin, bMax int64" declara los cuatro parámetros con el
	// mismo tipo; el tipo se escribe una sola vez, al final.
	common := min(aMax, bMax) - max(aMin, bMin)
	if common <= 0 {
		return 0
	}
	return common
}

// Intersects dice si a y b se enciman, es decir, si comparten volumen.
// Tocarse por una cara, una arista o un vértice NO es encimarse. Es simétrica:
// a.Intersects(b) == b.Intersects(a).
//
// # Cómo funciona
//
// Dos cajas alineadas a los ejes se enciman solo si sus rangos se traslapan
// en los TRES ejes a la vez. Basta con un eje donde estén separadas (o solo se
// toquen) para que no se encimen. En cada eje la comparación es estricta (<):
// si a termina justo donde empieza b (a.maxX() == b.XMM), no se traslapan.
//
// # Por qué así
//
// Revisar los tres rangos detecta también el caso de dos cajas cruzadas en
// forma de "+", donde ninguna esquina de una queda dentro de la otra; revisar
// solo las esquinas fallaría ahí. Una caja sin volumen nunca se encima
// (ver Box).
//
// # Nota de aprendizaje: receptor por valor
//
// "(a Box)" recibe una copia de la caja. Box son seis int64, 48 bytes: copiarla
// es barato, y la función no puede cambiar la caja de quien la llama. Cuándo
// conviene un puntero (*Box) se ve en el capítulo 06 de la guía de Go.
func (a Box) Intersects(b Box) bool {
	if !a.hasVolume() || !b.hasVolume() {
		return false
	}
	// Go: && evalúa de izquierda a derecha y se detiene en el primer false
	// (evaluación en cortocircuito): si x ya no se traslapa, no revisa y ni z.
	return a.XMM < b.maxX() && b.XMM < a.maxX() &&
		a.YMM < b.maxY() && b.YMM < a.maxY() &&
		a.ZMM < b.maxZ() && b.ZMM < a.maxZ()
}

// ContainedIn dice si inner queda completamente dentro de outer. Pegada al
// borde también cuenta como dentro: una pieza que termina justo en la pared
// del equipo cabe.
//
// # Cómo funciona
//
// En cada eje, inner tiene que empezar en o después del inicio de outer y
// terminar en o antes de su final. Aquí la comparación es <=, no <: en un rango
// semiabierto, terminar en el mismo milímetro que outer significa ocupar hasta
// el último milímetro de outer, sin pasarse.
//
// # Por qué así
//
// Una caja sin volumen no cuenta como contenida (devuelve false), aunque sus
// coordenadas queden dentro: así un error del algoritmo que arme una pieza de
// medida 0 no pasa como colocación válida. Una caja contenida en sí misma da
// true.
//
// No hace falta revisar aparte que outer tenga volumen: si inner mide más de 0
// sobre un eje y queda entre los bordes de outer, outer también mide más de 0
// sobre ese eje. Un espacio plano o con medida negativa nunca contiene nada.
func (inner Box) ContainedIn(outer Box) bool {
	if !inner.hasVolume() {
		return false
	}
	return outer.XMM <= inner.XMM && inner.maxX() <= outer.maxX() &&
		outer.YMM <= inner.YMM && inner.maxY() <= outer.maxY() &&
		outer.ZMM <= inner.ZMM && inner.maxZ() <= outer.maxZ()
}

// Contact devuelve el área, en mm², donde la cara de abajo de top se apoya en
// la cara de arriba de below. Si top no está justo encima de below, o solo se
// tocan por una arista o un vértice, devuelve 0. La usa la estabilidad del
// sprint 2.1 para saber qué parte de la base de una pieza tiene apoyo.
//
// # Cómo funciona
//
//  1. Las caras están a la misma altura solo si top empieza donde termina
//     below: top.YMM == below.maxY(). Con 1 mm de aire, o si se enciman en y,
//     no hay contacto.
//  2. Vistas desde arriba, las dos bases son rectángulos en x y z. El área de
//     contacto es el traslape en x por el traslape en z (overlap).
//
// # Por qué así
//
// No es simétrica: below.Contact(top) pregunta si below está encima de top, y
// da 0. El orden importa porque la estabilidad pregunta "¿en qué se apoya esta
// pieza?", y el apoyo solo viene de abajo.
//
// El área cabe en int64 sin problema: el traslape sobre un eje nunca pasa de
// la medida de la pieza más chica, y una pieza mide a lo más 30 000 mm por lado
// en el contrato. Así que el área es a lo más 30 000 × 30 000 = 9 × 10⁸ mm²,
// aunque las cajas no estén dentro de ningún tramo.
func (top Box) Contact(below Box) int64 {
	if !top.hasVolume() || !below.hasVolume() {
		return 0
	}
	if top.YMM != below.maxY() {
		return 0
	}
	// Si un traslape es 0 (solo se tocan por una arista), el producto da 0 y
	// no hace falta un caso aparte.
	return overlap(top.XMM, top.maxX(), below.XMM, below.maxX()) *
		overlap(top.ZMM, top.maxZ(), below.ZMM, below.maxZ())
}

// Box devuelve la caja que ocupa la pieza colocada, para usar Intersects,
// ContainedIn y Contact sobre las colocaciones de un plan.
func (p Placement) Box() Box {
	return Box{
		XMM: p.XMM, YMM: p.YMM, ZMM: p.ZMM,
		DXMM: p.DXMM, DYMM: p.DYMM, DZMM: p.DZMM,
	}
}
