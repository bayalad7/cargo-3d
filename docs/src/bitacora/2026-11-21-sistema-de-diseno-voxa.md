# 21 nov 2026 · Sistema de diseño Voxa

**Sesión:** Fuera del plan · **Área:** ui

Sesión fuera del plan de trabajo: no cambia el sprint actual (sigue el 1.1).

## Qué se hizo

- Se llevó al repositorio el sistema de diseño oficial de claude.ai («Voxa», `tokens.json` v4): logo Vóxel,
  paleta amarillo grúa 60 · 30 · 10 en claro y oscuro, tipografía Inter Tight + IBM Plex Mono, espaciado,
  radios, sombra y componentes `vx-`.
- Paquete nuevo `packages/ui` (`@cargo-3d/ui`):
  - `src/styles/tokens.css`: fuente única de verdad, con `light-dark()` y los bloques `@theme` de Tailwind CSS v4.
    Las paletas, tamaños de texto, radios y sombras de Tailwind quedaron desactivados.
  - `src/styles/components.css`: las clases oficiales y, como **propuesta**, el modal de confirmación, los toasts,
    la barra de acciones masivas, el botón destructivo y el estado en vivo.
  - `src/brand/`: el logo en 6 SVG limpios y el componente `Logo`, que cambia de tema solo.
  - `src/icons.ts` y `Icon`: íconos de Lucide con trazo 1.5 y el vocabulario `appIcons` del dominio.
  - Componentes React de presentación: `Button`, `StatusBadge`, `FilterChip`, `ConfirmDialog`, `Toast`,
    `Toaster`, `BulkActionBar` y `BulkAction`.
  - `preview/index.html`: vista previa en claro y oscuro lado a lado.
- Documentación: `docs/src/ux/sistema-diseno.md`, ADR-0014, reglas de diseño en `CLAUDE.md`, versiones de
  Tailwind CSS y Lucide en el ADR-0011.
- No se tocó ningún código de Go ni lógica de la app (todavía no hay pantallas).

## Decisiones

- Nombres de archivo en inglés y sin ñ (`components.css`, `sistema-diseno.md`), por el ADR-0009.
- «Voxa» es el nombre comercial en la interfaz; el proyecto **no** se renombra hasta confirmar la marca ante el IMPI.
- Alcance: CSS + componentes React solo de presentación, más la vista previa en HTML.
- pnpm 12 pide decidir qué dependencias ejecutan scripts al instalarse: `@parcel/watcher` (del modo `--watch`
  de Tailwind) quedó en `allowBuilds: false` en `pnpm-workspace.yaml`, porque ya trae binarios precompilados.

## Verificación

- `tsc --noEmit` (TypeScript 7.0.2, modo estricto) sin errores.
- `tailwindcss` 4.3.5 compila la vista previa; se revisó en captura de pantalla en los dos temas.

## Pendientes

1. Subir al sistema oficial de claude.ai los componentes y tokens marcados como «propuesta».
2. Conectar `@cargo-3d/ui` a `apps/web` en el sprint 1.6.
3. Tokens para móvil (React Native) en el sprint 3.2; íconos de app para Wails y Expo.
4. El Node.js de esta máquina es 22; el proyecto pide 26 LTS (pnpm solo lo advierte).

---
