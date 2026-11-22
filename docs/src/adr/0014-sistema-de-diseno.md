# ADR-0014: Sistema de diseño Voxa con Tailwind CSS v4 y Lucide

- **Estado:** Aceptado
- **Fecha:** 2026-11-21

## Contexto

Todavía no existe ninguna pantalla (`apps/web`, `apps/desktop` y `apps/mobile` están vacías), pero ya hay un
sistema de diseño oficial en claude.ai («Voxa», `tokens.json` versión 4): nombre comercial, logo Vóxel,
paleta «amarillo grúa» 60 · 30 · 10 en claro y oscuro, tipografía y componentes base con prefijo `vx-`.
Conviene dejar la base lista antes del sprint 1.6 (primera interfaz), para que web y escritorio arranquen
con los mismos colores, tamaños y componentes, y para que nadie escriba colores a mano.

## Decisión

1. **Paquete `packages/ui` (`@cargo-3d/ui`)** con los estilos, el logo y los componentes base. Web y
   escritorio lo comparten (Wails reutiliza `apps/web`), y `viewer-3d` puede leer sus variables CSS.
2. **Fuente única de verdad en el código:** `src/styles/tokens.css` (valores) y `src/styles/components.css`
   (componentes `vx-`). La documentación (`docs/src/ux/sistema-diseno.md`) los explica; si no coinciden,
   manda el código. Cada cambio de un valor se hace también en el sistema oficial de claude.ai.
3. **Tailwind CSS v4** (configuración en CSS con `@theme`, sin `tailwind.config.js`). Se borran sus paletas,
   tamaños de texto, radios y sombras por omisión para que solo existan los del sistema.
4. **Temas con `light-dark()`:** cada color se escribe una vez; `data-theme="light"`/`"dark"` lo fuerza y, sin
   atributo, se sigue al sistema operativo.
5. **Íconos con Lucide** (`lucide-react`), trazo 1.5 px, tamaños 14/16/18, y un vocabulario del dominio
   (`appIcons`) para que cada concepto tenga un solo dibujo.
6. **Componentes React solo de presentación:** reciben propiedades y avisan con funciones; no guardan estado
   de negocio ni llaman a la API.
7. **Nombres de archivo en inglés** (`components.css`, `voxa-logo.svg`), por el ADR-0009. Los nombres de los
   tokens y de las clases son los del sistema oficial.
8. **«Voxa» es el nombre comercial** en la interfaz y la documentación de diseño, pero el proyecto, el módulo
   de Go, los paquetes `@cargo-3d/*` y el repositorio **no se renombran** hasta confirmar la marca ante el IMPI.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Tailwind CSS v4 + clases `vx-` (elegida)** | Utilidades para el diseño de pantallas; tokens como variables CSS que lee cualquiera (React, R3F, HTML); configuración en CSS | Hay que aprender Tailwind; dos formas de dar estilo (utilidades y clases) |
| Solo CSS propio (sin Tailwind) | Cero dependencias | Cada pantalla escribe su propio layout; más CSS que mantener |
| CSS-in-JS (styled-components, vanilla-extract) | Tipos en los estilos | Más peso; no sirve para la vista previa en HTML ni para el visor 3D |
| Librería de componentes (MUI, Chakra, shadcn/ui) | Muchos componentes listos | Su propia estética y tokens; pelear contra ella para lograr la del sistema |
| Tokens en `tokens.json` + generador | Una sola fuente para web y móvil | Paso de build extra; el usuario pidió que la fuente sea el CSS |
| Íconos: Heroicons, Phosphor, Tabler | Buenos sets | Lucide tiene más íconos de logística y trazo configurable a 1.5 px |

## Consecuencias

- **Positivas:** las pantallas nacen con la marca correcta en claro y oscuro; los colores que no son del
  sistema no existen en Tailwind; el logo cambia de tema solo; la vista previa permite revisar el sistema sin app.
- **Negativas / costos:** `light-dark()` requiere navegadores de 2024 en adelante (WebView2 y Safari actuales
  cumplen). React Native no lee CSS: en el sprint 3.2 habrá que conectar los tokens (por ejemplo con NativeWind).
  Los componentes que nacieron en el código (modal, toasts, barra masiva) se subieron al sistema oficial en su v5
  (2026-11-21).
- **Qué tendríamos que hacer para revertirla:** los tokens son variables CSS normales; quitar Tailwind es
  borrar los bloques `@theme` y reemplazar las utilidades de las pantallas por clases propias.
