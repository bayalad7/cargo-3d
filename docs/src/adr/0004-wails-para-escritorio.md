# ADR-0004: Wails v3 para la app de escritorio

- **Estado:** Aceptado (se implementa en la fase 3, solo si un cliente lo pide)
- **Fecha:** 2026-11-20

## Contexto

Algunos clientes pedirán trabajar sin internet, imprimir directo o conectarse a sistemas locales.
No queremos construir una segunda interfaz para eso.

## Decisión

`apps/desktop` usa **Wails v3**: envuelve `apps/web` en una ventana nativa con el WebView del sistema,
y el backend de la ventana está escrito en Go, así que puede usar el motor (`internal/engine`)
directamente, sin servidor y sin conexión.

**Precaución:** Wails saca versiones con mucha frecuencia; se fija una versión exacta en `go.mod`
y se actualiza a propósito, no por accidente.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Wails v3 (elegida)** | Go del lado nativo: el mismo motor corre local; WebView del sistema (instaladores pequeños) | Proyecto más joven que Electron; el WebView varía por sistema operativo |
| Tauri 2 | Muy maduro y ligero | Su lado nativo es Rust: agrega un segundo lenguaje de backend |
| Electron | El más maduro | Instaladores muy grandes; lado nativo en Node, no en Go |

## Consecuencias

- **Positivas:** la versión de escritorio recupera el cálculo sin conexión, usando el mismo motor en Go.
- **Negativas:** hay que probar el WebView de cada sistema operativo.
- **Revertir:** como la UI es web, se podría envolver con otra herramienta sin reescribirla.
