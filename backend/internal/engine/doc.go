// Package engine es el motor de planeación de carga 3D de Cargo 3D.
//
// Aquí vive TODA la lógica de acomodo: modelos (Container, CargoItem, Placement,
// LoadPlan), el algoritmo de puntos extremos y las restricciones. No sabe nada de
// HTTP, base de datos ni interfaz: recibe datos y devuelve un plan. Así se puede
// probar solo y, si algún día hace falta, compilar a WebAssembly con TinyGo.
//
// Este paquete todavía está vacío a propósito: la estructura y las decisiones se
// documentaron primero. La implementación empieza en la fase 1, sprint 1.1
// (ver docs/src/plan.md).
//
// # Nota de aprendizaje
//
//   - Un comentario justo encima de "package" documenta el paquete completo.
//   - Un comentario justo encima de una función o tipo exportado lo documenta.
//   - "Exportado" en Go = empieza con mayúscula (Container). Con minúscula
//     (container) es privado del paquete. No hay palabras public/private como en PHP.
//   - "internal/" en la ruta significa que solo el código de este módulo puede
//     importar este paquete: nadie de fuera puede depender de él por accidente.
//
// Para ver la documentación: go doc ./internal/engine
package engine

// Version es la versión del contrato del motor (contracts/engine.schema.json).
// Si el formato de entrada o salida cambia, esta versión sube.
const Version = "0.1.0"
