package engine

import (
	"fmt"          // armar mensajes con valores (Sprintf)
	"unicode/utf8" // contar caracteres, no bytes
)

// Límites del contrato (contracts/engine.schema.json). Si el esquema cambia un
// mínimo o un máximo, se cambia aquí en el mismo commit.
const (
	maxIDLength          = 64
	maxDescriptionLength = 200 // en caracteres
	maxSections          = 1   // v0.1.0 (P-012 lo sube)
)

// Validate revisa que la entrada cumpla el contrato antes de calcular. Regresa
// nil si todo está bien o un *ValidationError con TODOS los problemas, en el
// orden en que aparecen en la entrada.
//
// # Cómo funciona
//
// Revisa, en este orden: la versión del contrato, el equipo (id, descripción,
// tramos, puerta, carga útil, tara, que no haya puerta con margen y que los
// tramos no se encimen) y cada pieza (id, descripción, medidas, peso, cantidad,
// orientaciones y campos reservados). Cada regla agrega un FieldError y sigue.
//
// # Por qué así
//
//   - Corre antes de cualquier cálculo: con medidas fuera del contrato,
//     Volume() se desborda sin aviso.
//   - Una pieza que no cabe NO es error: sale en unplaced con
//     EXCEEDS_DIMENSIONS o EXCEEDS_DOOR. Una puerta más grande que el interior
//     tampoco: en el reefer pasa y la puerta simplemente no limita (ADR-0016).
//   - Lo que solo se ve al leer el JSON (decimales, campos desconocidos, null)
//     lo revisa DecodeInput, que después llama a Validate.
func (in EngineInput) Validate() error {
	// v junta los problemas. Es un valor local; sus métodos son de puntero
	// porque agregan a la lista.
	var v validator

	if in.SchemaVersion != Version {
		v.add("schemaVersion", ErrSchemaVersion, "es %s y el motor usa %q", short(in.SchemaVersion), Version)
	}

	v.container(in.Container)

	if len(in.Items) == 0 {
		v.add("items", ErrNoItems, "debe traer al menos una pieza")
	}
	// seen guarda los ids ya vistos para detectar repetidos. map[string]int
	// relaciona cada id con el índice donde apareció primero.
	seen := make(map[string]int, len(in.Items))
	for i, item := range in.Items {
		path := fmt.Sprintf("items[%d]", i)
		v.item(path, item)
		if first, ok := seen[item.ID]; ok {
			v.add(path+".id", ErrDuplicateID, "%s ya se usó en items[%d]", short(item.ID), first)
		} else {
			seen[item.ID] = i
		}
	}

	// Go: si se regresara v.result() de tipo *ValidationError directamente,
	// un puntero nil dentro de una interfaz error NO es igual a nil (ver
	// capítulo 04, "el nil con tipo"). Por eso se pregunta y se regresa el
	// nil de la interfaz de forma explícita.
	if len(v.problems) == 0 {
		return nil
	}
	return &ValidationError{Problems: v.problems}
}

// validator junta los problemas mientras se recorre la entrada. Es privado:
// solo existe dentro de Validate.
type validator struct {
	problems []*FieldError
}

// add registra un problema. format y args funcionan como en fmt.Sprintf.
//
// Go: "args ...any" es un parámetro variádico: recibe cero o más valores de
// cualquier tipo y dentro de la función es un slice []any. "args..." al
// llamar a Sprintf lo vuelve a pasar como lista.
func (v *validator) add(path string, kind error, format string, args ...any) {
	v.problems = append(v.problems, &FieldError{
		Path:    path,
		Message: fmt.Sprintf(format, args...),
		Kind:    kind,
	})
}

// rangeCheck registra un problema si value no está entre min y max (incluidos).
func (v *validator) rangeCheck(path string, value, min, max int64) {
	if value < min || value > max {
		v.add(path, ErrOutOfRange, "es %d y debe estar entre %d y %d", value, min, max)
	}
}

// optionalRange es rangeCheck para un campo opcional: si no vino (nil), no hay
// nada que revisar.
func (v *validator) optionalRange(path string, value *int64, min, max int64) {
	if value != nil {
		v.rangeCheck(path, *value, min, max)
	}
}

// id revisa un identificador: de 1 a 64 caracteres, solo A-Z, a-z, 0-9, punto,
// guion y guion bajo (el patrón ^[A-Za-z0-9._-]{1,64}$ del contrato).
func (v *validator) id(path, id string) {
	if id == "" {
		v.add(path, ErrInvalidID, "no puede estar vacío")
		return
	}
	// Recorrer un string con índice (id[i]) da bytes, no caracteres. Un acento
	// ocupa dos bytes y ninguno de los dos está en la lista permitida, así que
	// cualquier carácter fuera de la lista se detecta aquí.
	for i := 0; i < len(id); i++ {
		if !isIDByte(id[i]) {
			v.add(path, ErrInvalidID, "%s solo puede tener letras sin acentos, dígitos, punto, guion y guion bajo", short(id))
			return
		}
	}
	// Si llegó aquí, cada carácter mide un byte: len(id) (que cuenta bytes) es
	// también el número de caracteres.
	if len(id) > maxIDLength {
		v.add(path, ErrInvalidID, "mide %d caracteres y el máximo es %d", len(id), maxIDLength)
	}
}

// maxEcho es cuántos caracteres de un valor enviado por el usuario se repiten
// en un mensaje de error.
const maxEcho = 64

// short regresa el texto entre comillas, recortado a maxEcho caracteres con
// «…» si es más largo. Así un id de un millón de caracteres no produce un
// mensaje de error de un millón de caracteres.
//
// Go: convertir un string a []rune lo separa en caracteres (runas), así el
// corte nunca parte un carácter de varios bytes a la mitad.
func short(text string) string {
	runes := []rune(text)
	if len(runes) <= maxEcho {
		return fmt.Sprintf("%q", text)
	}
	return fmt.Sprintf("%q…", string(runes[:maxEcho]))
}

// isIDByte dice si b es un byte permitido en un id.
//
// Go: un byte es un número de 0 a 255; 'a' es la constante del carácter a
// (97). Comparar rangos de bytes es la forma directa de revisar letras ASCII.
func isIDByte(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9') ||
		b == '.' || b == '-' || b == '_'
}

// description revisa un texto libre: máximo 200 caracteres.
//
// utf8.RuneCountInString cuenta caracteres (runas). len() contaría bytes, y
// «Contenedor 20' DV con ñ» tiene más bytes que caracteres. El contrato cuenta
// caracteres.
func (v *validator) description(path, text string) {
	if n := utf8.RuneCountInString(text); n > maxDescriptionLength {
		v.add(path, ErrTextTooLong, "mide %d caracteres y el máximo es %d", n, maxDescriptionLength)
	}
}

// container revisa el equipo.
func (v *validator) container(c Container) {
	v.id("container.id", c.ID)
	v.description("container.description", c.Description)

	if n := len(c.Sections); n < 1 || n > maxSections {
		v.add("container.sections", ErrSectionCount, "trae %d tramos y la versión %s admite %d", n, Version, maxSections)
	}
	sectionIDs := make(map[string]int, len(c.Sections))
	for i, s := range c.Sections {
		path := fmt.Sprintf("container.sections[%d]", i)
		v.section(path, s)
		if first, ok := sectionIDs[s.ID]; ok {
			v.add(path+".id", ErrDuplicateID, "%s ya se usó en container.sections[%d]", short(s.ID), first)
		} else {
			sectionIDs[s.ID] = i
		}

		// Un margen mayor que cero solo tiene sentido en un tramo abierto.
		// Un margen escrito en 0 con puerta sí se acepta: no rebasa nada.
		if c.Door != nil {
			if m := s.Margin(); m.TopMM > 0 || m.LeftMM > 0 || m.RightMM > 0 {
				v.add(path+".overhang", ErrDoorWithOverhang, "un equipo con puerta no puede tener margen de sobredimensión (ADR-0017)")
			}
		}

		// Cada par de tramos se compara una sola vez: j empieza después de i.
		for j := i + 1; j < len(c.Sections); j++ {
			if sectionsOverlap(s, c.Sections[j]) {
				v.add(path, ErrSectionsOverlap, "se encima con container.sections[%d]", j)
			}
		}
	}

	if c.Door != nil {
		// Una puerta más grande que el interior NO es error (ADR-0016): solo
		// se revisan los límites del contrato.
		v.rangeCheck("container.door.widthMm", c.Door.WidthMM, 1, 3000)
		v.rangeCheck("container.door.heightMm", c.Door.HeightMM, 1, 4250)
	}

	v.rangeCheck("container.maxPayloadKg", c.MaxPayloadKG, 1, 200000)
	v.optionalRange("container.tareKg", c.TareKG, 0, 100000)
}

// section revisa un tramo: su id, su esquina, sus medidas y su margen.
func (v *validator) section(path string, s Section) {
	v.id(path+".id", s.ID)
	v.rangeCheck(path+".xMm", s.XMM, 0, 30000)
	v.rangeCheck(path+".yMm", s.YMM, 0, 4250)
	v.rangeCheck(path+".zMm", s.ZMM, 0, 3000)
	v.rangeCheck(path+".lengthMm", s.LengthMM, 1, 30000)
	v.rangeCheck(path+".widthMm", s.WidthMM, 1, 3000)
	v.rangeCheck(path+".heightMm", s.HeightMM, 1, 4250)
	if s.Overhang != nil {
		v.optionalRange(path+".overhang.topMm", s.Overhang.TopMM, 0, 2000)
		v.optionalRange(path+".overhang.leftMm", s.Overhang.LeftMM, 0, 1000)
		v.optionalRange(path+".overhang.rightMm", s.Overhang.RightMM, 0, 1000)
	}
}

// sectionsOverlap dice si dos tramos ocupan espacio en común. Que solo se
// toquen (el final de uno es el inicio del otro) no es encimarse: la
// comparación es estricta (<), como en docs/src/conceptos/sistema-coordenadas.md.
//
// Es una versión mínima para la validación; la geometría completa llega en el
// sprint 1.2.
func sectionsOverlap(a, b Section) bool {
	return a.XMM < b.XMM+b.LengthMM && b.XMM < a.XMM+a.LengthMM &&
		a.YMM < b.YMM+b.HeightMM && b.YMM < a.YMM+a.HeightMM &&
		a.ZMM < b.ZMM+b.WidthMM && b.ZMM < a.ZMM+a.WidthMM
}

// item revisa una pieza.
func (v *validator) item(path string, it CargoItem) {
	v.id(path+".id", it.ID)
	v.description(path+".description", it.Description)
	v.rangeCheck(path+".lengthMm", it.LengthMM, 1, 30000)
	v.rangeCheck(path+".widthMm", it.WidthMM, 1, 30000)
	v.rangeCheck(path+".heightMm", it.HeightMM, 1, 30000)
	v.rangeCheck(path+".weightKg", it.WeightKG, 1, 100000)
	v.rangeCheck(path+".quantity", it.Quantity, 1, 10000)
	v.orientations(path+".allowedOrientations", it.AllowedOrientations)
	v.optionalRange(path+".maxLoadOnTopKg", it.MaxLoadOnTopKG, 0, 1000000)
	v.optionalRange(path+".deliverySequence", it.DeliverySequence, 1, 999)
}

// orientations revisa la lista de orientaciones permitidas.
//
// Solo nil significa "no vino" (se permiten las seis). Una lista vacía que no
// es nil viene de "allowedOrientations": [] y es un error (P-025): si se
// dejara pasar, al volver a JSON desaparecería y la pieza tendría las seis.
func (v *validator) orientations(path string, list []Orientation) {
	if list == nil {
		return
	}
	if len(list) == 0 {
		v.add(path, ErrNoOrientations, "no puede venir vacía; si se permiten todas, no la envíes")
		return
	}

	// map[Orientation]bool como conjunto: true = ya apareció.
	seen := make(map[Orientation]bool, len(list))
	for i, o := range list {
		p := fmt.Sprintf("%s[%d]", path, i)
		if !o.IsValid() {
			v.add(p, ErrUnknownOrientation, "%s no es una de LWH, WLH, LHW, HLW, WHL, HWL", short(string(o)))
			continue
		}
		if seen[o] {
			v.add(p, ErrDuplicateOrientation, "%s está repetida", short(string(o)))
			continue
		}
		seen[o] = true
	}
}
