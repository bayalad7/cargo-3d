package engine

import (
	"math"
	"testing"
)

// Pruebas de fuzzing de la geometría (sesión 1.2.3). Explicación completa en
// docs/src/aprendizaje/go/08-pruebas-de-tabla-y-fuzzing.md.
//
// Con "go test" normal, cada Fuzz... corre SOLO con sus semillas (las llamadas
// a f.Add) y con los casos guardados en testdata/fuzz/<NombreDelFuzz>/: es una
// prueba rápida de siempre. Para que Go genere entradas nuevas hay que pedirlo,
// de un objetivo por vez:
//
//	go test ./backend/internal/engine -run '^$' -fuzz '^FuzzIntersectsIsSymmetric$' -fuzztime 1m
//
// Los cinco comandos están en CLAUDE.md, sección "Comandos".

// wrap lleva v al rango [lo, hi] (los dos incluidos) sin importar qué número sea.
// Si v ya está en el rango, lo regresa igual.
//
// # Cómo funciona
//
// Go le da al fuzzing enteros de cualquier tamaño, incluidos los que desbordan.
// La geometría solo es confiable con números chicos (Box.InRange), así que casi
// todas las propiedades generan sus cajas con wrap. Un número que ya está en el
// rango pasa sin cambios: así las semillas (1 000, 1 200, 30 000…) llegan a la
// propiedad tal cual, y el fuzzing parte de los casos límite de verdad. Uno que
// se sale se reparte dentro del rango: lo + v % span, donde span es el tamaño
// del rango. No es una vuelta cíclica (wrap(30001, 1, 30000) da 2), pero cae
// siempre dentro y nunca desborda: calcular v - lo sí se desbordaría con
// math.MaxInt64.
//
// # Por qué así
//
// La primera versión daba siempre lo + v % span, que corre TODO número por lo:
// con lo = 1, wrap(30000, 1, 30000) daba 1 y la semilla "el máximo del contrato"
// probaba una pieza de 1 mm. TestWrapKeepsNumbersInRange lo evita.
//
// # Nota de aprendizaje: % con negativos
//
// En Go el resto de un número negativo es negativo (-7 % 5 == -2), por eso se
// le suma span cuando sale menor que 0. v % span no desborda ni con
// math.MinInt64, porque span siempre es positivo y chico.
func wrap(v, lo, hi int64) int64 {
	if lo <= v && v <= hi {
		return v
	}
	span := hi - lo + 1
	m := v % span
	if m < 0 {
		m += span
	}
	return lo + m
}

// boundedBox arma una caja dentro de los límites del contrato a partir de seis
// enteros cualquiera: coordenadas de −1 000 a 60 000 mm y medidas de −100 a
// 30 000 mm. Deja pasar medidas en 0 y negativas a propósito: son las "cajas sin
// volumen", que las tres funciones tienen que tratar bien.
func boundedBox(x, y, z, dx, dy, dz int64) Box {
	return Box{
		XMM: wrap(x, -1000, 60000), YMM: wrap(y, -1000, 60000), ZMM: wrap(z, -1000, 60000),
		DXMM: wrap(dx, -100, 30000), DYMM: wrap(dy, -100, 30000), DZMM: wrap(dz, -100, 30000),
	}
}

// boxSeeds devuelve las semillas de dos cajas (seis números cada una): los casos
// límite de las pruebas de tabla (se tocan por una cara, por una arista, por un
// vértice, encimadas 1 mm, separadas 1 mm, una dentro de otra, cruz en "+", caja
// plana y apoyo parcial).
func boxSeeds() [][12]int64 {
	return [][12]int64{
		{0, 0, 0, 1000, 1000, 1000, 0, 0, 0, 1000, 1000, 1000},          // idénticas
		{0, 0, 0, 1000, 1000, 1000, 1000, 0, 0, 1000, 1000, 1000},       // una cara
		{0, 0, 0, 1000, 1000, 1000, 1000, 1000, 0, 1000, 1000, 1000},    // una arista
		{0, 0, 0, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000}, // un vértice
		{0, 0, 0, 1000, 1000, 1000, 999, 0, 0, 1000, 1000, 1000},        // encimadas 1 mm
		{0, 0, 0, 1000, 1000, 1000, 1001, 0, 0, 1000, 1000, 1000},       // separadas 1 mm
		{0, 0, 0, 1000, 1000, 1000, 200, 200, 200, 100, 100, 100},       // una dentro de otra
		{0, 0, 400, 1000, 100, 200, 400, 0, 0, 200, 100, 1000},          // cruz en "+"
		{0, 0, 0, 1000, 1000, 1000, 200, 200, 200, 0, 100, 100},         // caja plana
		{0, 0, 0, 1200, 1000, 800, 400, 1000, 200, 1200, 1000, 800},     // apoyo parcial
	}
}

// addBoxSeeds agrega a f las semillas de boxSeeds.
func addBoxSeeds(f *testing.F) {
	for _, s := range boxSeeds() {
		f.Add(s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8], s[9], s[10], s[11])
	}
}

// FuzzIntersectsIsSymmetric revisa que a.Intersects(b) y b.Intersects(a) den
// siempre lo mismo, con CUALQUIER int64 (sin límites).
//
// Aquí no se usa wrap: la fórmula de Intersects es la misma con las dos cajas
// intercambiadas, así que es simétrica incluso cuando la suma se desborda. Se
// prueba con todo el rango para comprobarlo, no para suponerlo.
func FuzzIntersectsIsSymmetric(f *testing.F) {
	addBoxSeeds(f)
	// Extremos de int64 que las semillas chicas no cubren.
	f.Add(int64(math.MaxInt64), int64(0), int64(0), int64(1000), int64(1000), int64(1000),
		int64(0), int64(0), int64(0), int64(1000), int64(1000), int64(1000))
	f.Add(int64(math.MinInt64), int64(0), int64(0), int64(math.MinInt64), int64(1), int64(1),
		int64(0), int64(0), int64(0), int64(1), int64(1), int64(1))

	f.Fuzz(func(t *testing.T, ax, ay, az, adx, ady, adz, bx, by, bz, bdx, bdy, bdz int64) {
		a := Box{XMM: ax, YMM: ay, ZMM: az, DXMM: adx, DYMM: ady, DZMM: adz}
		b := Box{XMM: bx, YMM: by, ZMM: bz, DXMM: bdx, DYMM: bdy, DZMM: bdz}

		if ab, ba := a.Intersects(b), b.Intersects(a); ab != ba {
			t.Fatalf("no es simétrica: a.Intersects(b) = %t, b.Intersects(a) = %t\na = %+v\nb = %+v", ab, ba, a, b)
		}
	})
}

// FuzzBoxIntersectsItselfOnlyWithVolume revisa que una caja se encime consigo
// misma si y solo si tiene volumen (las tres medidas positivas).
func FuzzBoxIntersectsItselfOnlyWithVolume(f *testing.F) {
	f.Add(int64(0), int64(0), int64(0), int64(1000), int64(1000), int64(1000)) // con volumen
	f.Add(int64(0), int64(0), int64(0), int64(1), int64(1), int64(1))          // de 1 mm
	f.Add(int64(0), int64(0), int64(0), int64(0), int64(1000), int64(1000))    // plana
	f.Add(int64(0), int64(0), int64(0), int64(1000), int64(-5), int64(1000))   // negativa
	f.Add(int64(-1000), int64(0), int64(60000), int64(30000), int64(30000), int64(30000))

	f.Fuzz(func(t *testing.T, x, y, z, dx, dy, dz int64) {
		a := boundedBox(x, y, z, dx, dy, dz)
		if !a.InRange() {
			t.Fatalf("boundedBox armó una caja fuera de rango: %+v", a)
		}

		want := a.DXMM > 0 && a.DYMM > 0 && a.DZMM > 0
		if got := a.Intersects(a); got != want {
			t.Fatalf("a.Intersects(a) = %t, se esperaba %t\na = %+v", got, want, a)
		}
	})
}

// FuzzContainedInIsTransitive revisa que si A está dentro de B y B dentro de C,
// A esté dentro de C.
//
// # Cómo funciona
//
// Tres cajas al azar casi nunca quedan una dentro de otra, y la propiedad solo
// dice algo cuando se cumple la premisa. Por eso las cajas se arman en cadena:
// C es una caja con volumen; B se arma dentro de C y A dentro de B (boxInside),
// y a cada una se le da un empujón de -50 a 50 mm sobre cada eje para que a
// veces se salga por x, por y o por z. La propiedad se revisa sobre lo que
// realmente quedó (si A cabe en B y B cabe en C, entonces A cabe en C), no sobre
// lo que se pretendía.
func FuzzContainedInIsTransitive(f *testing.F) {
	// c: esquina y medidas (6). b y a: cada una lleva nueve números para
	// boxInside: medidas (3), posición dentro de la caja de afuera (3) y empuje
	// (3), todos en el orden x, y, z.
	f.Add(int64(0), int64(0), int64(0), int64(5898), int64(2393), int64(2352),
		int64(3000), int64(1000), int64(1000), int64(100), int64(0), int64(0), int64(0), int64(0), int64(0), // b
		int64(1000), int64(500), int64(500), int64(100), int64(0), int64(0), int64(0), int64(0), int64(0)) // a
	f.Add(int64(0), int64(0), int64(0), int64(1000), int64(1000), int64(1000), // idénticas
		int64(1000), int64(1000), int64(1000), int64(0), int64(0), int64(0), int64(0), int64(0), int64(0),
		int64(1000), int64(1000), int64(1000), int64(0), int64(0), int64(0), int64(0), int64(0), int64(0))
	f.Add(int64(0), int64(0), int64(0), int64(1000), int64(1000), int64(1000), // B se sale 1 mm por x
		int64(1000), int64(1000), int64(1000), int64(0), int64(0), int64(0), int64(1), int64(0), int64(0),
		int64(1000), int64(1000), int64(1000), int64(0), int64(0), int64(0), int64(0), int64(0), int64(0))
	f.Add(int64(0), int64(0), int64(0), int64(1000), int64(1000), int64(1000), // A se sale 1 mm por z
		int64(1000), int64(1000), int64(1000), int64(0), int64(0), int64(0), int64(0), int64(0), int64(0),
		int64(1000), int64(1000), int64(1000), int64(0), int64(0), int64(0), int64(0), int64(0), int64(1))

	f.Fuzz(func(t *testing.T,
		cx, cy, cz, cdx, cdy, cdz int64,
		bdx, bdy, bdz, bpx, bpy, bpz, bjx, bjy, bjz int64,
		adx, ady, adz, apx, apy, apz, ajx, ajy, ajz int64,
	) {
		// Go: el último parámetro lleva coma al final porque la lista ocupa
		// varios renglones; gofmt la exige.
		c := Box{
			XMM: wrap(cx, -1000, 60000), YMM: wrap(cy, -1000, 60000), ZMM: wrap(cz, -1000, 60000),
			DXMM: wrap(cdx, 1, 30000), DYMM: wrap(cdy, 1, 30000), DZMM: wrap(cdz, 1, 30000),
		}
		b := boxInside(c, [9]int64{bdx, bdy, bdz, bpx, bpy, bpz, bjx, bjy, bjz})
		a := boxInside(b, [9]int64{adx, ady, adz, apx, apy, apz, ajx, ajy, ajz})

		// Comprobación de la propia prueba: las cajas están en rango, así que
		// el resultado de ContainedIn es confiable.
		for _, box := range []Box{a, b, c} {
			if !box.InRange() {
				t.Fatalf("caja fuera de rango: %+v", box)
			}
		}

		if a.ContainedIn(b) && b.ContainedIn(c) && !a.ContainedIn(c) {
			t.Fatalf("no es transitiva\na = %+v\nb = %+v\nc = %+v", a, b, c)
		}
	})
}

// boxInside arma una caja con volumen que, sin el empuje, queda dentro de outer.
// Si outer ya es la caja más chica posible (1 mm por eje) queda igual a outer.
//
// v lleva nueve números, cada grupo en el orden x, y, z:
//   - v[0:3] son las medidas de la caja, de 1 a la medida de outer sobre ese eje;
//   - v[3:6] son la posición de su esquina dentro del espacio que sobra, de 0 a
//     (medida de outer − medida de la caja);
//   - v[6:9] son el empuje sobre cada eje, de −50 a 50 mm: puede sacar la caja
//     de outer por cualquiera de los tres ejes.
//
// Cada eje tiene sus propios números: el fuzzer puede mover una caja en z sin
// moverla en x.
func boxInside(outer Box, v [9]int64) Box {
	dx := wrap(v[0], 1, outer.DXMM)
	dy := wrap(v[1], 1, outer.DYMM)
	dz := wrap(v[2], 1, outer.DZMM)
	return Box{
		XMM:  outer.XMM + wrap(v[3], 0, outer.DXMM-dx) + wrap(v[6], -50, 50),
		YMM:  outer.YMM + wrap(v[4], 0, outer.DYMM-dy) + wrap(v[7], -50, 50),
		ZMM:  outer.ZMM + wrap(v[5], 0, outer.DZMM-dz) + wrap(v[8], -50, 50),
		DXMM: dx, DYMM: dy, DZMM: dz,
	}
}

// FuzzRotationKeepsVolume revisa que girar una pieza nunca cambie su volumen:
// en las seis orientaciones, largo × ancho × alto da lo mismo. También revisa
// que cada candidata de Candidates ocupe ese mismo volumen.
func FuzzRotationKeepsVolume(f *testing.F) {
	f.Add(int64(1200), int64(800), int64(600))
	f.Add(int64(1000), int64(1000), int64(1000))
	f.Add(int64(30000), int64(30000), int64(30000)) // el máximo del contrato
	f.Add(int64(1), int64(1), int64(1))
	f.Add(int64(12000), int64(300), int64(300))

	f.Fuzz(func(t *testing.T, l, w, h int64) {
		item := CargoItem{
			ID:       "fuzz",
			LengthMM: wrap(l, 1, 30000), WidthMM: wrap(w, 1, 30000), HeightMM: wrap(h, 1, 30000),
			WeightKG: 1, Quantity: 1,
		}
		want := item.Volume()

		for _, o := range allOrientations() {
			d, ok := item.RotatedDimensions(o)
			if !ok {
				t.Fatalf("%s: RotatedDimensions dijo que no la conoce", o)
			}
			if got := d.DXMM * d.DYMM * d.DZMM; got != want {
				t.Fatalf("%s: volumen %d, se esperaba %d\nitem = %+v", o, got, want, item)
			}
		}

		// Un espacio de 30 000 mm por lado cabe cualquier orientación, así que
		// tiene que salir al menos una candidata y todas con el mismo volumen.
		space := Box{DXMM: 30000, DYMM: 30000, DZMM: 30000}
		candidates := item.Candidates(Point{}, space)
		if len(candidates) == 0 {
			t.Fatalf("ninguna candidata en un espacio de 30 000 mm por lado\nitem = %+v", item)
		}
		for _, c := range candidates {
			if got := c.Box.DXMM * c.Box.DYMM * c.Box.DZMM; got != want {
				t.Fatalf("candidata %s: volumen %d, se esperaba %d\nitem = %+v", c.Orientation, got, want, item)
			}
		}
	})
}

// FuzzContactIsBoundedAndExcludesOverlap revisa dos cosas del área de apoyo:
//   - nunca es negativa ni mayor que la base (x × z) de la caja más chica de
//     las dos, y es 0 si alguna no tiene volumen;
//   - si hay contacto (área mayor que 0), las cajas no se enciman.
func FuzzContactIsBoundedAndExcludesOverlap(f *testing.F) {
	addBoxSeeds(f)

	f.Fuzz(func(t *testing.T, ax, ay, az, adx, ady, adz, bx, by, bz, bdx, bdy, bdz int64) {
		top := boundedBox(ax, ay, az, adx, ady, adz)
		below := boundedBox(bx, by, bz, bdx, bdy, bdz)

		got := top.Contact(below)
		if got < 0 {
			t.Fatalf("Contact negativa: %d\ntop = %+v\nbelow = %+v", got, top, below)
		}

		hasVolume := func(b Box) bool { return b.DXMM > 0 && b.DYMM > 0 && b.DZMM > 0 }
		if !hasVolume(top) || !hasVolume(below) {
			if got != 0 {
				t.Fatalf("una caja sin volumen tuvo contacto: %d\ntop = %+v\nbelow = %+v", got, top, below)
			}
			return
		}

		topBase, belowBase := top.DXMM*top.DZMM, below.DXMM*below.DZMM
		if limit := min(topBase, belowBase); got > limit {
			t.Fatalf("Contact = %d, pasa de la base más chica (%d)\ntop = %+v\nbelow = %+v", got, limit, top, below)
		}
		if got > 0 && top.Intersects(below) {
			t.Fatalf("hay contacto (%d) y a la vez se enciman\ntop = %+v\nbelow = %+v", got, top, below)
		}
	})
}

// TestWrapKeepsNumbersInRange revisa que wrap deje igual un número que ya está en
// el rango (para que las semillas lleguen a la propiedad tal cual) y que meta
// dentro del rango a uno que se sale, incluso con los extremos de int64.
func TestWrapKeepsNumbersInRange(t *testing.T) {
	tests := []struct {
		name            string
		v, lo, hi, want int64
	}{
		{name: "el máximo del contrato se queda", v: 30000, lo: 1, hi: 30000, want: 30000},
		{name: "un número en medio se queda", v: 1200, lo: 1, hi: 30000, want: 1200},
		{name: "el inicio del rango se queda", v: 1, lo: 1, hi: 30000, want: 1},
		{name: "cero se queda en un rango que lo incluye", v: 0, lo: -1000, hi: 60000, want: 0},
		{name: "un negativo dentro del rango se queda", v: -100, lo: -100, hi: 30000, want: -100},
		// Fuera del rango, el resultado es lo + v % span: cae dentro, no es una
		// vuelta cíclica del rango.
		{name: "1 arriba del fin cae dentro del rango", v: 30001, lo: 1, hi: 30000, want: 2},
		{name: "cero en un rango que empieza en 1 cae en el inicio", v: 0, lo: 1, hi: 30000, want: 1},
		{name: "un negativo en un rango positivo cae dentro", v: -7, lo: 1, hi: 30000, want: 29994},
		{name: "rango de un solo valor", v: 12345, lo: 7, hi: 7, want: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wrap(tt.v, tt.lo, tt.hi); got != tt.want {
				t.Errorf("wrap(%d, %d, %d) = %d, se esperaba %d", tt.v, tt.lo, tt.hi, got, tt.want)
			}
		})
	}

	// Cualquier int64, incluidos los extremos, cae dentro del rango.
	for _, v := range []int64{math.MinInt64, math.MinInt64 + 1, -1, math.MaxInt64 - 1, math.MaxInt64} {
		if got := wrap(v, -1000, 60000); got < -1000 || got > 60000 {
			t.Errorf("wrap(%d, -1000, 60000) = %d, quedó fuera del rango", v, got)
		}
	}
}

// TestBoundedBoxKeepsTheSeeds revisa que boundedBox no cambie ninguna semilla de
// boxSeeds: son números del contrato, así que tienen que llegar a las
// propiedades tal cual (un caso "se tocan por una cara" tiene que seguir
// tocándose).
func TestBoundedBoxKeepsTheSeeds(t *testing.T) {
	for i, s := range boxSeeds() {
		wantA := Box{XMM: s[0], YMM: s[1], ZMM: s[2], DXMM: s[3], DYMM: s[4], DZMM: s[5]}
		wantB := Box{XMM: s[6], YMM: s[7], ZMM: s[8], DXMM: s[9], DYMM: s[10], DZMM: s[11]}

		if got := boundedBox(s[0], s[1], s[2], s[3], s[4], s[5]); got != wantA {
			t.Errorf("semilla %d, caja A: boundedBox dio %+v, se esperaba %+v", i, got, wantA)
		}
		if got := boundedBox(s[6], s[7], s[8], s[9], s[10], s[11]); got != wantB {
			t.Errorf("semilla %d, caja B: boundedBox dio %+v, se esperaba %+v", i, got, wantB)
		}
	}
}
