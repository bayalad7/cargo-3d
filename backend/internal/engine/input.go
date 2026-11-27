package engine

// EngineInput es lo que recibe el motor: un equipo y la lista de piezas por
// acomodar. Corresponde a "EngineInput" de contracts/engine.schema.json.
//
// # Cómo funciona
//
// Se llena directo desde el JSON con encoding/json: cada campo lleva una
// etiqueta con el nombre exacto que usa el contrato. No se completa nada al
// leer: lo que no vino se queda sin venir y los valores por omisión se
// resuelven con métodos (ver CargoItem.IsStackable y Section.Margin).
//
// # Por qué así
//
// El contrato es la fuente de verdad (regla 9 de CLAUDE.md). Si el JSON que
// entra es el mismo que sale, la web, el validador y el motor ven exactamente
// los mismos datos.
type EngineInput struct {
	// SchemaVersion es la versión del contrato con que se escribió la entrada.
	// Debe ser igual a Version; lo revisa la validación (sesión 1.1.4).
	//
	// Go: "string" es texto inmutable (sus bytes no se pueden cambiar; se crea
	// otro string).
	SchemaVersion string `json:"schemaVersion"`

	// Container es el equipo donde se acomoda la carga.
	//
	// Go: un campo puede ser otro struct completo (no un puntero): el
	// Container vive DENTRO de EngineInput, como un objeto anidado en JSON.
	Container Container `json:"container"`

	// Items son los tipos de pieza, cada uno con su cantidad.
	//
	// Go: "[]CargoItem" es un slice, una lista de tamaño variable de CargoItem,
	// todos del mismo tipo. Ojo al copiarlo: b := a NO copia los elementos,
	// los comparte, y cambiar b[0] cambia a[0] (para copiarlos: slices.Clone).
	// Se explica a fondo en el capítulo 05.
	Items []CargoItem `json:"items"`
}
