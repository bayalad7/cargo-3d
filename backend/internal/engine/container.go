package engine

// Este archivo define el equipo (Container) y sus partes: tramos (Section),
// margen de sobredimensión (Overhang) y puerta (Door). Todo sigue
// contracts/engine.schema.json; ver también ADR-0015 (coordenadas), ADR-0016
// (peso efectivo y redondeo) y ADR-0017 (overhang).
//
// # Nota de aprendizaje: struct
//
// Un "struct" agrupa campos, cada uno con su nombre y su tipo fijo:
//
//	type Door struct { WidthMM int64; HeightMM int64 }
//
// Lo esencial:
//   - Un struct es un VALOR: b := a hace una COPIA de sus campos. La copia es
//     superficial: si un campo es un puntero o un slice, se copia la
//     referencia y las dos copias siguen viendo los mismos datos (ver "Qué se
//     copia y qué se comparte" en el capítulo 02).
//   - Door{} ya existe con todos sus campos en su "valor cero" (0 para
//     números, "" para texto, false para bool, nil para punteros, slices y
//     mapas). El lenguaje no tiene constructores: cuando crear un valor
//     necesita lógica, por convención se escribe una función NewDoor(...).
//     La función new(T) reserva un T en su valor cero y regresa un puntero a
//     él; no llama a nada.
//   - Leer a través de un puntero nil (por ejemplo *c.Door sin puerta) detiene
//     el programa con un panic, igual que escribir en un mapa nil. Por eso se
//     pregunta "!= nil" antes (ver derefOrZero).
//   - Un struct se arma juntando otros (composición): EngineInput contiene un
//     Container y el Container contiene sus Section.
//   - Mayúscula inicial = exportado (visible fuera del paquete). Por eso todos
//     los campos empiezan con mayúscula: encoding/json solo ve los exportados.
//
// # Nota de aprendizaje: etiquetas (tags)
//
// El texto entre comillas invertidas después de un campo es una "etiqueta":
// metadatos que el compilador ignora y que otros paquetes leen mientras el
// programa corre (con el paquete reflect). encoding/json lee `json:"nombre"`
// para saber cómo se llama el campo en JSON.
//   - `json:"widthMm"` → en JSON se llama widthMm (el contrato usa camelCase).
//   - `json:"door,omitempty"` → si el campo está en su valor cero, no se
//     escribe al convertir a JSON. Se usa para los campos opcionales.
//
// # Nota de aprendizaje: int64 y no float64
//
// Todas las medidas (mm) y pesos (kg) son int64: enteros con signo de 64 bits
// (de −9,2 × 10¹⁸ a 9,2 × 10¹⁸). Dos variables float64 que valen 0,1 y 0,2
// suman 0,30000000000000004 (escrito como constante, 0.1 + 0.2 da 0.3 porque
// Go suma las constantes exactas al compilar). Con decimales, dos cajas que se
// tocan podrían "traslaparse" por 0,0000001 mm, y un peso exacto podría pasarse
// del límite por redondeo interno. Con enteros, 1 200 + 1 200 = 2 400 siempre.
// Go no mezcla tipos numéricos por su cuenta: int64 y float64 no se suman sin
// escribir la conversión, por ejemplo float64(volumen).
// ¿Por qué int64 y no int? "int" cambia de tamaño según la máquina (32 o 64
// bits); int64 mide lo mismo en todos lados, incluso compilado a WebAssembly.

// Container es el equipo donde se acomoda la carga: contenedor marítimo, caja
// seca o plataforma. Al motor solo le importan su geometría y sus límites; la
// categoría del equipo (reefer, flat rack…) vive en el catálogo
// (contracts/equipment.json), no aquí.
//
// # Cómo funciona
//
// El espacio útil son sus Sections (en la v0.1.0, exactamente una). La puerta,
// si existe, limita qué orientaciones pueden entrar. MaxPayloadKG es el único
// límite de peso que el motor conoce.
//
// # Por qué así
//
// Separar geometría (Sections) de categoría deja que el mismo motor acomode
// un contenedor de 20', una caja de 53' o una plataforma escalonada (varios
// tramos, sprint 2.1) sin saber qué son.
type Container struct {
	// ID identifica el equipo dentro de este cálculo; no es el número ISO 6346.
	// Go: por convención, las siglas van completas en mayúscula (ID, no Id).
	ID string `json:"id"`

	// Description es texto libre para mostrar, por ejemplo «Contenedor 20' DV».
	// Opcional: con omitempty, si viene vacío no se escribe de regreso. El
	// contrato no distingue "" de "no vino" para un texto que solo se muestra.
	Description string `json:"description,omitempty"`

	// Sections son los tramos del espacio de carga.
	Sections []Section `json:"sections"`

	// Door es la puerta trasera, en el extremo de x máxima. Opcional.
	//
	// Go: "*Door" es un PUNTERO a Door: guarda la dirección de un Door o nil
	// (nada). Aquí lo usamos para decir "no vino": una plataforma abierta no
	// tiene puerta y Door queda en nil. Con un Door normal (sin *) no habría
	// forma de distinguir "sin puerta" de "puerta de 0 × 0". Los punteros se
	// ven a fondo en el capítulo 06.
	Door *Door `json:"door,omitempty"`

	// MaxPayloadKG es el peso máximo de carga EFECTIVO, en kg, ya resuelto
	// FUERA del motor (ADR-0016).
	//
	// Quién lo calcula: la API (o el catálogo de equipos de la empresa), antes
	// de llamar al motor. Es el MENOR entre dos números:
	//   1. la carga útil del equipo (ficha del fabricante o placa CSC), y
	//   2. lo que permite la NOM-012-SCT-2 para la configuración vehicular
	//      (T3-S2, T3-S3…) y el tipo de camino (ET, A, B…) de ese viaje,
	//      menos la tara del tracto y del equipo.
	//
	// Por qué el motor no conoce la NOM-012: la norma cambia (hay proyectos
	// de modificación) y depende del viaje, no de la geometría. Si viviera
	// aquí, cada cambio de la norma obligaría a tocar y volver a probar el
	// algoritmo. Así el motor solo responde "¿cabe y no pasa de X kg?" y la
	// tabla de la NOM se prueba aparte (sesión 2.1.5): las reglas de negocio
	// que cambian quedan fuera del cálculo.
	//
	// Redondeo: este límite se redondea HACIA ABAJO al capturarlo (28 750,6 kg
	// → 28 750). Así el motor nunca cree que cabe más peso del real.
	MaxPayloadKG int64 `json:"maxPayloadKg"`

	// TareKG es la tara del equipo (su peso vacío), para calcular el VGM
	// (SOLAS VI/2). El motor no la usa para acomodar. Opcional.
	//
	// Es *int64 y no int64 porque 0 es un valor válido distinto de "no vino":
	// con int64 + omitempty, una tara escrita como 0 desaparecería al volver a
	// JSON. nil = no vino; un puntero a 0 = vino y vale 0.
	TareKG *int64 `json:"tareKg,omitempty"`
}

// Volume devuelve el volumen NOMINAL del equipo en mm³: la suma del volumen
// de sus tramos, sin contar el margen de sobredimensión (overhang). Igual que
// Section.Volume, solo es confiable con entrada validada.
//
// # Por qué así
//
// Es el denominador de volumeUtilization en las métricas (contrato,
// "Metrics"). Por eso puede pasar de 1 la ocupación cuando hay piezas
// sobredimensionadas: ocupan un margen que aquí no se cuenta.
//
// # Nota de aprendizaje: métodos
//
// "func (c Container) Volume() int64" es un MÉTODO: una función atada a un
// tipo. Lo que va entre paréntesis antes del nombre, (c Container), es el
// "receptor": un parámetro normal que va antes del nombre del método, con el
// nombre que tú elijas (por convención, la inicial del tipo). c.Volume() es lo
// mismo que Container.Volume(c).
// Este receptor es "por valor": el método recibe una COPIA del Container. Si
// cambiara c.ID, el original no se enteraría. Pero la copia es superficial:
// c.Sections sigue apuntando a los MISMOS tramos y c.Door a la MISMA puerta,
// así que c.Sections[0].LengthMM = 0 o c.Door.WidthMM = 0 sí cambiarían el
// original. Volume solo lee, así que no importa aquí. Cuándo conviene un
// receptor por puntero se explica en el capítulo 03 (ver LoadPlan).
func (c Container) Volume() int64 {
	// "var total int64" declara una variable con su tipo y valor cero (0).
	// Equivale a "total := int64(0)".
	var total int64

	// "for _, s := range c.Sections" recorre el slice: en cada vuelta, s es
	// una copia del tramo. El "_" descarta el índice, que no usamos. Go exige
	// usar todas las variables locales declaradas; "_" es la forma de decir
	// "esto no lo quiero".
	for _, s := range c.Sections {
		total += s.Volume()
	}

	return total
}

// Section es un tramo del espacio de carga: una caja alineada a los ejes con
// su esquina mínima en coordenadas del equipo (x hacia la puerta, y hacia
// arriba, z hacia la derecha; ADR-0015) y sus medidas útiles.
//
// Redondeo (ADR-0016): las medidas interiores se capturan redondeadas HACIA
// ABAJO (2 352,8 mm → 2 352). Si el equipo mide un poco más de lo que dice,
// sobra un poco de holgura; nunca al revés, así ninguna pieza "cabe" en el
// cálculo y no en la realidad. Redondear es trabajo de quien captura (la API
// o la pantalla): el motor recibe enteros y no redondea nada.
type Section struct {
	// ID identifica el tramo, único dentro del equipo.
	ID string `json:"id"`

	// XMM, YMM y ZMM son la esquina mínima del tramo, en mm.
	// Go: el nombre sigue la regla de siglas en mayúscula (X + MM).
	XMM int64 `json:"xMm"`
	YMM int64 `json:"yMm"`
	ZMM int64 `json:"zMm"`

	// LengthMM es el largo útil sobre x.
	LengthMM int64 `json:"lengthMm"`
	// WidthMM es el ancho útil sobre z.
	WidthMM int64 `json:"widthMm"`
	// HeightMM es el alto útil sobre y. En una plataforma es 4 250 mm
	// (NOM-012) menos la altura de la cama; ese cálculo también se hace fuera.
	HeightMM int64 `json:"heightMm"`

	// Overhang es el margen de sobredimensión de un tramo abierto. Opcional:
	// si es nil, no se permite rebasar. Para leerlo sin preguntar por nil,
	// usa Margin().
	Overhang *Overhang `json:"overhang,omitempty"`
}

// Volume devuelve el volumen nominal del tramo en mm³ (largo × ancho × alto).
//
// Solo con entrada validada (sesión 1.1.4): dentro de los límites del contrato
// el máximo es 30 000 × 3 000 × 4 250 ≈ 3,8 × 10¹¹ mm³, muy por debajo del
// tope de int64 (≈ 9,2 × 10¹⁸). Con medidas fuera de esos límites la
// multiplicación se desborda: en Go un int64 que pasa del tope da la vuelta a
// un número equivocado SIN error ni aviso. Por eso la validación corre antes
// de cualquier cálculo.
func (s Section) Volume() int64 {
	return s.LengthMM * s.WidthMM * s.HeightMM
}

// Box devuelve la caja nominal del tramo: su esquina mínima y su largo, alto
// y ancho útiles, sin el margen de rebase (overhang). Es el espacio que se
// pasa a CargoItem.Candidates cuando la pieza no puede rebasar.
//
// Ojo con el orden: Section guarda largo, ancho y alto (LengthMM, WidthMM,
// HeightMM), y Box va en x, y, z, así que el alto (y) queda en medio.
func (s Section) Box() Box {
	return Box{
		XMM: s.XMM, YMM: s.YMM, ZMM: s.ZMM,
		DXMM: s.LengthMM, DYMM: s.HeightMM, DZMM: s.WidthMM,
	}
}

// Margin devuelve el margen de sobredimensión del tramo ya resuelto: si
// Overhang no vino, un margen en cero (no se permite rebasar); si vino, sus
// valores, con 0 en los lados que no vinieron.
//
// # Por qué así
//
// El contrato dice que lo que no viene debe regresar sin venir. Por eso el
// valor por omisión no se escribe al leer el JSON (se perdería la diferencia
// entre "no vino" y "vino en 0"), sino que se calcula aquí, al usarlo.
//
// Go: Margin regresa OverhangValues por valor (una copia), no un puntero: quien
// lo usa no puede cambiar la entrada por accidente.
func (s Section) Margin() OverhangValues {
	// "if s.Overhang == nil" pregunta si el puntero no apunta a nada. Si no
	// vino, el valor cero de OverhangValues ya es "todo en 0".
	if s.Overhang == nil {
		return OverhangValues{}
	}

	// Un struct se crea escribiendo sus campos por nombre; los que no se
	// escriben quedan en su valor cero. derefOrZero (abajo) convierte cada *int64 en int64.
	return OverhangValues{
		TopMM:   derefOrZero(s.Overhang.TopMM),
		LeftMM:  derefOrZero(s.Overhang.LeftMM),
		RightMM: derefOrZero(s.Overhang.RightMM),
	}
}

// Overhang es el margen de sobredimensión (carga OOG, out of gauge) tal como
// viene en el JSON: cuánto puede rebasar una pieza el tramo por arriba y a cada
// lado. No hay margen al frente ni atrás (ADR-0017). Los límites legales
// (NOM-012: 4 250 mm de alto total y 2 600 mm de ancho) los aplica quien llama
// al motor al fijar estos valores.
//
// Cada campo es *int64 porque el contrato permite escribir 0 a propósito
// (hay ejemplos con "topMm": 0) y eso debe regresar igual. Para calcular, usa
// Section.Margin(), que ya resuelve los nil.
type Overhang struct {
	// TopMM es cuánto puede rebasar por arriba de HeightMM.
	TopMM *int64 `json:"topMm,omitempty"`
	// LeftMM es cuánto puede rebasar por la izquierda (z menor que el inicio).
	LeftMM *int64 `json:"leftMm,omitempty"`
	// RightMM es cuánto puede rebasar por la derecha (z mayor que el final).
	RightMM *int64 `json:"rightMm,omitempty"`
}

// OverhangValues es el margen de sobredimensión ya resuelto, sin nil: lo que
// el algoritmo usa para calcular. No tiene etiquetas json porque no viaja en
// el contrato; solo existe dentro del motor.
type OverhangValues struct {
	TopMM   int64 // margen por arriba, en mm
	LeftMM  int64 // margen por la izquierda, en mm
	RightMM int64 // margen por la derecha, en mm
}

// Door es la puerta trasera del equipo, en el extremo de x máxima. Una pieza
// entra a lo largo de x, así que debe cumplir dz ≤ WidthMM y dy ≤ HeightMM en
// alguna orientación permitida (eso se revisa en el sprint 1.3).
//
// Puede ser más grande que el interior (en el reefer la puerta es más alta que
// la altura útil bajo la línea de carga): en ese caso no limita y no es error.
// Se captura redondeada HACIA ABAJO, igual que el interior.
type Door struct {
	// WidthMM es el ancho libre de la puerta, sobre z.
	WidthMM int64 `json:"widthMm"`
	// HeightMM es el alto libre de la puerta, sobre y.
	HeightMM int64 `json:"heightMm"`
}

// derefOrZero devuelve el valor al que apunta p, o 0 si p es nil.
//
// Go: empieza con minúscula, así que no se exporta: la usa cualquier archivo
// del paquete engine, pero nadie de fuera. "*p" lee el
// valor al que apunta el puntero ("desreferenciar"); hacerlo con p == nil
// detendría el programa con un panic, por eso se pregunta antes.
func derefOrZero(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}
