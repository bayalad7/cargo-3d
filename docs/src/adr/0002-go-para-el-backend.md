# ADR-0002: Go para todo el backend (motor y API)

- **Estado:** Aceptado
- **Fecha:** 2026-11-20

## Contexto

El motor de acomodo resuelve una variante del *3D bin packing*, un problema NP-difícil que se ataca con
heurísticas que prueban muchas combinaciones: necesita buen rendimiento. Al mismo tiempo, Benjamín
mantendrá el sistema solo y viene de PHP, así que el lenguaje debe aprenderse rápido y leerse con facilidad.
La otra opción natural para un motor rápido es Rust, cuya curva de aprendizaje es muy alta.

## Decisión

Todo el backend (API y motor) se escribe en **Go**, en un solo módulo y un solo binario. El motor corre
en el servidor. El principio rector es: **Go para todo el backend, TypeScript para todas las pantallas.**

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Go (elegida)** | Lenguaje pequeño, se aprende en semanas; compilado y rápido; goroutines para la metaheurística en paralelo; un binario sin dependencias; recolector de basura (sin borrow checker) | Menos rápido que Rust; WebAssembly pesado con el compilador estándar; pocas librerías de optimización |
| Rust | El más rápido; WASM de primera clase | Curva de aprendizaje alta |
| Python | El más fácil; las mejores librerías de optimización | Lento en ciclos puros (requiere NumPy o Numba); tipado dinámico; despliegue más pesado; no corre bien en el navegador |
| C# / .NET | Rápido y tipado | No es más fácil que Go; ecosistema completo aparte y más ceremonia |
| PHP | Benjamín ya lo domina | No apto para cálculo intensivo ni paralelismo |

## Consecuencias

- **Positivas:** un solo lenguaje de backend, simple y didáctico; despliegue de un binario más PostgreSQL;
  la metaheurística puede usar todos los núcleos con goroutines; manejo de errores explícito.
- **Negativas:** el cálculo del plan ya no ocurre en el navegador: requiere conexión y el servidor paga
  el cómputo. Para cargas de cientos de piezas sigue siendo rápido.
- **Mitigación futura:** el motor no depende de HTTP ni de base de datos, así que podría compilarse a
  WebAssembly con TinyGo si algún día se necesita cálculo sin conexión (con la salvedad de que algunas
  opciones de tamaño de TinyGo no admiten goroutines).
- **Revertir:** los contratos (`contracts/`) y los casos de prueba (`engine/fixtures/`) son independientes
  del lenguaje; el motor podría reescribirse en otro lenguaje sin tocar las pantallas.
