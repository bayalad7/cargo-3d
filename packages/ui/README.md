# @cargo-3d/ui — Sistema de diseño Voxa

Tokens, estilos (Tailwind CSS v4), logo, íconos (Lucide) y componentes base de todas las pantallas.
Documentación completa: [`docs/src/ux/sistema-diseno.md`](../../docs/src/ux/sistema-diseno.md). Decisión: ADR-0014.

```css
@import "@cargo-3d/ui/styles.css"; /* fuentes + Tailwind + tokens + componentes */
```

```tsx
import { Button, Logo, StatusBadge, appIcons } from "@cargo-3d/ui";
```

| Archivo | Qué es |
|---|---|
| `src/styles/tokens.css` | **Fuente única de verdad**: colores claro/oscuro, tipografía, espaciado, radios, sombras |
| `src/styles/components.css` | Clases `vx-` (las del sistema oficial, más `.vx-data`) |
| `src/brand/` | Logo Vóxel: `Logo.tsx` y los SVG con colores fijos |
| `src/icons.ts` | `appIcons`: un ícono de Lucide por concepto del dominio |
| `src/components/` | Componentes React solo de presentación |
| `preview/` | Vista previa en HTML: `pnpm --filter @cargo-3d/ui build:preview` y abrir `preview/index.html` |

Reglas: ningún valor escrito a mano en una pantalla; los componentes no importan la API ni el motor.
