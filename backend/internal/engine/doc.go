// Package engine es el motor de planeación de carga 3D de Cargo 3D.
//
// Aquí vive TODA la lógica de acomodo: modelos (Container, CargoItem, Placement,
// LoadPlan), el algoritmo de puntos extremos y las restricciones. No sabe nada de
// HTTP, base de datos ni interfaz: recibe datos y devuelve un plan. Así se puede
// probar solo y, si algún día hace falta, compilar a WebAssembly con TinyGo.
//
// Los tipos siguen contracts/engine.schema.json campo por campo:
//
//   - EngineInput (input.go): la entrada completa.
//   - Container, Section, Overhang y Door (container.go): el equipo.
//   - CargoItem (cargo_item.go): las piezas.
//   - Orientation y Dimensions (orientation.go): cómo se gira una pieza.
//   - Placement (placement.go): dónde quedó una unidad.
//   - Box, con Intersects, ContainedIn y Contact (box.go): la geometría de
//     cajas alineadas a los ejes, ver docs/src/conceptos/geometria.md.
//   - LoadPlan, UnplacedItem y Metrics (load_plan.go): la salida.
//   - DecodeInput (decode.go) y EngineInput.Validate (validate.go): la lectura.
//   - Las categorías de error y FieldError (errors.go), ver ADR-0023.
//
// Falta el algoritmo (sprint 1.3). Ver docs/src/plan.md.
//
// Medidas en milímetros y pesos en kilogramos, siempre int64 y ya redondeados
// por quien llama (ADR-0016). Explicación en
// docs/src/aprendizaje/go/02-structs-y-tipos.md.
//
// # Nota de aprendizaje
//
//   - Un comentario justo encima de "package" documenta el paquete completo.
//   - Un comentario justo encima de una función o tipo exportado lo documenta.
//   - "Exportado" en Go = empieza con mayúscula (Container). Con minúscula
//     (container) solo se ve dentro del paquete. No hay palabras clave de
//     visibilidad: la decide la primera letra.
//   - "internal/" en la ruta significa que solo el código de este módulo puede
//     importar este paquete: nadie de fuera puede depender de él por accidente.
//
// Para ver la documentación: go doc ./internal/engine
package engine

// Version es la versión del contrato del motor (contracts/engine.schema.json).
// Si el formato de entrada o salida cambia, esta versión sube.
const Version = "0.1.0"
