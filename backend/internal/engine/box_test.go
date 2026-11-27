package engine

import (
	"math"
	"testing"
)

// unitCubeBox es una caja de 1000 mm por lado en el origen. Casi todos los casos de
// abajo se escriben contra ella para que solo cambie lo que se prueba.
var unitCubeBox = Box{XMM: 0, YMM: 0, ZMM: 0, DXMM: 1000, DYMM: 1000, DZMM: 1000}

// TestIntersectsOnlyWhenSharingVolume revisa que dos cajas se encimen solo si
// comparten volumen: tocarse por una cara, una arista o un vértice no cuenta,
// y 1 mm de traslape sí. Cada caso se revisa en los dos órdenes (A con B y B
// con A), porque Intersects tiene que ser simétrica.
func TestIntersectsOnlyWhenSharingVolume(t *testing.T) {
	tests := []struct {
		name string
		a, b Box
		want bool
	}{
		{
			name: "idénticas",
			a:    unitCubeBox, b: unitCubeBox,
			want: true,
		},
		{
			name: "una dentro de la otra",
			a:    unitCubeBox,
			b:    Box{XMM: 200, YMM: 200, ZMM: 200, DXMM: 100, DYMM: 100, DZMM: 100},
			want: true,
		},
		{
			name: "encimadas 1 mm sobre x",
			a:    unitCubeBox,
			b:    Box{XMM: 999, DXMM: 1000, DYMM: 1000, DZMM: 1000},
			want: true,
		},
		{
			name: "se tocan por una cara sobre x (una detrás de otra)",
			a:    unitCubeBox,
			b:    Box{XMM: 1000, DXMM: 1000, DYMM: 1000, DZMM: 1000},
			want: false,
		},
		{
			name: "se tocan por una cara sobre y (una encima de otra)",
			a:    unitCubeBox,
			b:    Box{YMM: 1000, DXMM: 1000, DYMM: 1000, DZMM: 1000},
			want: false,
		},
		{
			name: "se tocan por una cara sobre z (una al lado de otra)",
			a:    unitCubeBox,
			b:    Box{ZMM: 1000, DXMM: 1000, DYMM: 1000, DZMM: 1000},
			want: false,
		},
		{
			name: "se tocan por una arista",
			a:    unitCubeBox,
			b:    Box{XMM: 1000, YMM: 1000, DXMM: 1000, DYMM: 1000, DZMM: 1000},
			want: false,
		},
		{
			name: "se tocan por un vértice",
			a:    unitCubeBox,
			b:    Box{XMM: 1000, YMM: 1000, ZMM: 1000, DXMM: 1000, DYMM: 1000, DZMM: 1000},
			want: false,
		},
		{
			name: "separadas por 1 mm sobre x",
			a:    unitCubeBox,
			b:    Box{XMM: 1001, DXMM: 1000, DYMM: 1000, DZMM: 1000},
			want: false,
		},
		{
			name: "separadas por 1 mm sobre y",
			a:    unitCubeBox,
			b:    Box{YMM: 1001, DXMM: 1000, DYMM: 1000, DZMM: 1000},
			want: false,
		},
		{
			name: "separadas por 1 mm sobre z",
			a:    unitCubeBox,
			b:    Box{ZMM: 1001, DXMM: 1000, DYMM: 1000, DZMM: 1000},
			want: false,
		},
		{
			// Una barra larga sobre x y otra larga sobre z que se cruzan por el
			// centro: ninguna esquina de una queda dentro de la otra, pero sí
			// comparten volumen. Revisar solo esquinas fallaría aquí.
			name: "cruzadas en forma de +",
			a:    Box{XMM: 0, YMM: 0, ZMM: 400, DXMM: 1000, DYMM: 100, DZMM: 200},
			b:    Box{XMM: 400, YMM: 0, ZMM: 0, DXMM: 200, DYMM: 100, DZMM: 1000},
			want: true,
		},
		{
			// Traslapan en x y en z, pero sobre y una está arriba de la otra con
			// aire entre ellas: basta un eje separado.
			name: "traslapan en dos ejes y no en el tercero",
			a:    unitCubeBox,
			b:    Box{XMM: 500, YMM: 1500, ZMM: 500, DXMM: 1000, DYMM: 1000, DZMM: 1000},
			want: false,
		},
		{
			// Rebase por la izquierda (overhang.leftMm): z negativa.
			name: "con z negativa, encimadas",
			a:    Box{ZMM: -300, DXMM: 1000, DYMM: 1000, DZMM: 500},
			b:    unitCubeBox,
			want: true,
		},
		{
			name: "caja plana (alto 0) dentro de otra",
			a:    unitCubeBox,
			b:    Box{XMM: 200, YMM: 200, ZMM: 200, DXMM: 100, DYMM: 0, DZMM: 100},
			want: false,
		},
		{
			name: "caja con medida negativa dentro de otra",
			a:    unitCubeBox,
			b:    Box{XMM: 200, YMM: 200, ZMM: 200, DXMM: -100, DYMM: 100, DZMM: 100},
			want: false,
		},
		{
			name: "caja vacía (valor cero) consigo misma",
			a:    Box{}, b: Box{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// %+v imprime un struct con el nombre de cada campo.
			if got := tt.a.Intersects(tt.b); got != tt.want {
				t.Errorf("a.Intersects(b) = %t, se esperaba %t\na = %+v\nb = %+v", got, tt.want, tt.a, tt.b)
			}
			if got := tt.b.Intersects(tt.a); got != tt.want {
				t.Errorf("b.Intersects(a) = %t, se esperaba %t (no es simétrica)\na = %+v\nb = %+v", got, tt.want, tt.a, tt.b)
			}
		})
	}
}

// TestContainedInAllowsTouchingTheBorder revisa que una caja pegada al borde
// exacto del equipo quede dentro, y que pasarse 1 mm por cualquier lado la deje
// fuera.
func TestContainedInAllowsTouchingTheBorder(t *testing.T) {
	// space es el interior de un 20' estándar (5 898 × 2 393 × 2 352).
	space := Box{DXMM: 5898, DYMM: 2393, DZMM: 2352}

	tests := []struct {
		name         string
		inner, outer Box
		want         bool
	}{
		{
			name:  "una dentro de otra, sin tocar paredes",
			inner: Box{XMM: 100, YMM: 0, ZMM: 100, DXMM: 1200, DYMM: 1000, DZMM: 800},
			outer: space,
			want:  true,
		},
		{
			name:  "en la esquina del origen (toca fondo, piso y pared izquierda)",
			inner: Box{DXMM: 1200, DYMM: 1000, DZMM: 800},
			outer: space,
			want:  true,
		},
		{
			name:  "en la esquina máxima (toca puerta, techo y pared derecha)",
			inner: Box{XMM: 5898 - 1200, YMM: 2393 - 1000, ZMM: 2352 - 800, DXMM: 1200, DYMM: 1000, DZMM: 800},
			outer: space,
			want:  true,
		},
		{
			name:  "idéntica al espacio",
			inner: space,
			outer: space,
			want:  true,
		},
		{
			name:  "se pasa 1 mm por la puerta (x máxima)",
			inner: Box{XMM: 5898 - 1199, DXMM: 1200, DYMM: 1000, DZMM: 800},
			outer: space,
			want:  false,
		},
		{
			name:  "se pasa 1 mm por el techo (y máxima)",
			inner: Box{YMM: 2393 - 999, DXMM: 1200, DYMM: 1000, DZMM: 800},
			outer: space,
			want:  false,
		},
		{
			name:  "se pasa 1 mm por la pared derecha (z máxima)",
			inner: Box{ZMM: 2352 - 799, DXMM: 1200, DYMM: 1000, DZMM: 800},
			outer: space,
			want:  false,
		},
		{
			name:  "se pasa 1 mm por el fondo (x negativa)",
			inner: Box{XMM: -1, DXMM: 1200, DYMM: 1000, DZMM: 800},
			outer: space,
			want:  false,
		},
		{
			name:  "se pasa 1 mm por el piso (y negativa)",
			inner: Box{YMM: -1, DXMM: 1200, DYMM: 1000, DZMM: 800},
			outer: space,
			want:  false,
		},
		{
			name:  "se pasa 1 mm por la pared izquierda (z negativa)",
			inner: Box{ZMM: -1, DXMM: 1200, DYMM: 1000, DZMM: 800},
			outer: space,
			want:  false,
		},
		{
			name:  "más grande que el espacio",
			inner: Box{DXMM: 6000, DYMM: 1000, DZMM: 800},
			outer: space,
			want:  false,
		},
		{
			// Flat rack con 150 mm de margen por la izquierda: el espacio
			// permitido empieza en z = −150 y una pieza en z = −150 cabe.
			name:  "z negativa dentro de un espacio con margen a la izquierda",
			inner: Box{ZMM: -150, DXMM: 1200, DYMM: 1000, DZMM: 800},
			outer: Box{ZMM: -150, DXMM: 5638, DYMM: 2233, DZMM: 2438 + 300},
			want:  true,
		},
		{
			name:  "caja plana con coordenadas dentro",
			inner: Box{XMM: 100, YMM: 100, ZMM: 100, DXMM: 0, DYMM: 100, DZMM: 100},
			outer: space,
			want:  false,
		},
		{
			name:  "espacio sin volumen",
			inner: Box{},
			outer: Box{},
			want:  false,
		},
		{
			// Solo el espacio está mal (ancho 0): aunque la pieza sea real, no
			// puede quedar dentro de algo que no ocupa lugar. ContainedIn no lo
			// revisa aparte; sale solo de comparar los bordes.
			name:  "caja real en un espacio de ancho 0",
			inner: Box{DXMM: 1200, DYMM: 1000, DZMM: 800},
			outer: Box{DXMM: 5898, DYMM: 2393, DZMM: 0},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.inner.ContainedIn(tt.outer); got != tt.want {
				t.Errorf("ContainedIn() = %t, se esperaba %t\ninner = %+v\nouter = %+v", got, tt.want, tt.inner, tt.outer)
			}
		})
	}
}

// TestContactIsTheSharedAreaOfTouchingFaces revisa que el área de contacto
// sea la parte de la base de top que descansa sobre la cara de arriba de
// below, y 0 cuando no hay apoyo real (aire, traslape, solo una arista o el
// orden al revés).
func TestContactIsTheSharedAreaOfTouchingFaces(t *testing.T) {
	// pallet es una tarima EPAL cargada de 1 200 × 1 000 × 800 en el piso.
	pallet := Box{DXMM: 1200, DYMM: 1000, DZMM: 800}

	tests := []struct {
		name       string
		top, below Box
		want       int64
	}{
		{
			name:  "apoyo completo: misma base, una encima de otra",
			top:   Box{YMM: 1000, DXMM: 1200, DYMM: 1000, DZMM: 800},
			below: pallet,
			want:  1200 * 800,
		},
		{
			name:  "apoyo parcial: recorrida 400 en x y 200 en z",
			top:   Box{XMM: 400, YMM: 1000, ZMM: 200, DXMM: 1200, DYMM: 1000, DZMM: 800},
			below: pallet,
			want:  (1200 - 400) * (800 - 200),
		},
		{
			name:  "la de arriba es más grande: el apoyo es la cara de la de abajo",
			top:   Box{XMM: -100, YMM: 1000, ZMM: -100, DXMM: 1400, DYMM: 500, DZMM: 1000},
			below: pallet,
			want:  1200 * 800,
		},
		{
			name:  "la de arriba es más chica: el apoyo es toda su base",
			top:   Box{XMM: 100, YMM: 1000, ZMM: 100, DXMM: 400, DYMM: 300, DZMM: 300},
			below: pallet,
			want:  400 * 300,
		},
		{
			name:  "1 mm de aire entre las caras",
			top:   Box{YMM: 1001, DXMM: 1200, DYMM: 1000, DZMM: 800},
			below: pallet,
			want:  0,
		},
		{
			name:  "encimadas 1 mm en y: es traslape, no contacto",
			top:   Box{YMM: 999, DXMM: 1200, DYMM: 1000, DZMM: 800},
			below: pallet,
			want:  0,
		},
		{
			name:  "a la misma altura pero sin traslape en x (al lado)",
			top:   Box{XMM: 1300, YMM: 1000, DXMM: 1200, DYMM: 1000, DZMM: 800},
			below: pallet,
			want:  0,
		},
		{
			name:  "solo se tocan por una arista",
			top:   Box{XMM: 1200, YMM: 1000, DXMM: 1200, DYMM: 1000, DZMM: 800},
			below: pallet,
			want:  0,
		},
		{
			name:  "solo se tocan por un vértice",
			top:   Box{XMM: 1200, YMM: 1000, ZMM: 800, DXMM: 1200, DYMM: 1000, DZMM: 800},
			below: pallet,
			want:  0,
		},
		{
			name:  "orden al revés: la de abajo no se apoya en la de arriba",
			top:   pallet,
			below: Box{YMM: 1000, DXMM: 1200, DYMM: 1000, DZMM: 800},
			want:  0,
		},
		{
			name:  "lado a lado tocándose por una cara vertical",
			top:   Box{XMM: 1200, DXMM: 1200, DYMM: 1000, DZMM: 800},
			below: pallet,
			want:  0,
		},
		{
			name:  "caja plana encima",
			top:   Box{YMM: 1000, DXMM: 1200, DYMM: 0, DZMM: 800},
			below: pallet,
			want:  0,
		},
		{
			name:  "la de abajo es plana",
			top:   Box{YMM: 0, DXMM: 1200, DYMM: 1000, DZMM: 800},
			below: Box{DXMM: 1200, DYMM: 0, DZMM: 800},
			want:  0,
		},
		{
			// Lo más grande que permite el contrato: una pieza mide a lo más
			// 30 000 mm por lado, y el traslape no pasa de la pieza más chica.
			// Sin suponer que estén dentro de un tramo, el área máxima es
			// 30 000 × 30 000 = 9 × 10⁸ mm². Cabe de sobra en int64.
			name:  "máximo del contrato sin desbordar",
			top:   Box{XMM: 30000, YMM: 30000, ZMM: -1000, DXMM: 30000, DYMM: 30000, DZMM: 30000},
			below: Box{XMM: 30000, YMM: 0, ZMM: -1000, DXMM: 30000, DYMM: 30000, DZMM: 30000},
			want:  900_000_000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.top.Contact(tt.below); got != tt.want {
				t.Errorf("Contact() = %d, se esperaba %d\ntop   = %+v\nbelow = %+v", got, tt.want, tt.top, tt.below)
			}
		})
	}
}

// TestPlacementBoxCopiesCornerAndDimensions revisa que la caja de una
// colocación tenga su esquina mínima y sus medidas, sin cruzar ejes. Usa el
// segundo cajón de contracts/examples/maritime/dry/20dv-heavy.
func TestPlacementBoxCopiesCornerAndDimensions(t *testing.T) {
	p := Placement{
		ItemID: "crate", XMM: 0, YMM: 0, ZMM: 1100,
		DXMM: 1400, DYMM: 1000, DZMM: 1100, Orientation: OrientationLWH,
	}
	want := Box{XMM: 0, YMM: 0, ZMM: 1100, DXMM: 1400, DYMM: 1000, DZMM: 1100}

	// Go: dos structs se pueden comparar con == si todos sus campos son
	// comparables (aquí, todos son int64).
	if got := p.Box(); got != want {
		t.Errorf("Box() = %+v, se esperaba %+v", got, want)
	}
}

// TestInRangeAcceptsOnlySmallNumbers revisa que InRange acepte una caja solo
// si sus seis números tienen valor absoluto de a lo más MaxGeometryMM, y que
// rechace el desbordamiento por cualquiera de los dos lados, incluido el
// número más negativo que cabe en int64.
func TestInRangeAcceptsOnlySmallNumbers(t *testing.T) {
	tests := []struct {
		name string
		box  Box
		want bool
	}{
		{name: "caja del contrato", box: Box{XMM: 60000, YMM: 0, ZMM: -1000, DXMM: 30000, DYMM: 30000, DZMM: 30000}, want: true},
		{name: "caja vacía (todo en 0)", box: Box{}, want: true},
		{name: "caja sin volumen de números chicos", box: Box{DXMM: -5, DYMM: 0, DZMM: 10}, want: true},
		{name: "justo en el límite positivo", box: Box{XMM: MaxGeometryMM}, want: true},
		{name: "justo en el límite negativo", box: Box{ZMM: -MaxGeometryMM}, want: true},
		{name: "1 mm arriba del límite en x", box: Box{XMM: MaxGeometryMM + 1}, want: false},
		{name: "1 mm abajo del límite en z", box: Box{ZMM: -MaxGeometryMM - 1}, want: false},
		{name: "medida de más de 100 000", box: Box{DYMM: MaxGeometryMM + 1}, want: false},
		{name: "medida negativa fuera de rango", box: Box{DXMM: -MaxGeometryMM - 1}, want: false},
		{name: "coordenada cerca del tope de int64", box: Box{XMM: math.MaxInt64 - 500, DXMM: 1000}, want: false},
		{name: "el número más negativo de int64", box: Box{YMM: math.MinInt64}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.box.InRange(); got != tt.want {
				t.Errorf("InRange() = %t, se esperaba %t\nbox = %+v", got, tt.want, tt.box)
			}
		})
	}
}

// TestIntersectsIsWrongOutOfRange documenta POR QUÉ existe InRange: con una
// coordenada cerca del tope de int64, xMm + dxMm da la vuelta a un número
// negativo e Intersects dice "no se enciman" aunque las dos cajas son
// idénticas. Si algún día Intersects se vuelve segura con cualquier número,
// esta prueba falla y hay que revisar la regla de InRange.
func TestIntersectsIsWrongOutOfRange(t *testing.T) {
	huge := Box{XMM: math.MaxInt64 - 500, DXMM: 1000, DYMM: 1000, DZMM: 1000}

	if huge.InRange() {
		t.Fatal("la caja de prueba debería estar fuera de rango")
	}
	// Idénticas: tendrían que encimarse. Salen "separadas" por la vuelta de int64.
	if huge.Intersects(huge) {
		t.Error("se esperaba el resultado equivocado (false) que causa el desbordamiento; ¿ya es segura con números enormes?")
	}
}
