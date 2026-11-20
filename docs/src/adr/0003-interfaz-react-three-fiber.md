# ADR-0003: TypeScript + React + React Three Fiber para la interfaz

- **Estado:** Aceptado
- **Fecha:** 2026-11-20

## Contexto

La interfaz debe mostrar el plan de carga en 3D interactivo (girar, acercar, mover piezas),
además de formularios, tablas e importación de Excel. La misma interfaz se reutiliza en escritorio (Wails).

## Decisión

- **TypeScript** en modo estricto para todas las interfaces.
- **React** con **Vite** como herramienta de build.
- **Three.js** a través de **React Three Fiber (R3F)** para el 3D, en el paquete reutilizable `packages/viewer-3d`.
- La interfaz pide los planes a la API en Go; el editor manual valida en vivo con `packages/validator`.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **React + R3F (elegida)** | El ecosistema 3D web más grande; la escena 3D como componentes; conocimiento reutilizable en React Native | React tiene conceptos propios que aprender |
| Svelte / Vue + Three.js | Más ligeros | Ecosistema 3D menor; sin puente a móvil nativo |
| HTMX + plantillas de Go | Todo en Go | No sirve para un visor 3D interactivo |
| Babylon.js | Motor 3D completo | Más pesado; comunidad menor |

## Consecuencias

- **Positivas:** el visor 3D es un componente que se usa igual en web y escritorio.
- **Negativas:** hay que aprender React y TypeScript además de Go.
- **Revertir:** el visor está aislado en su paquete; el resto de la UI se podría migrar sin tocar el backend.
