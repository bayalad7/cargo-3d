# Sistema de diseño Voxa

Este documento describe el sistema de diseño de la marca **Voxa**: el logo **Vóxel** (el píxel en 3D),
los colores por función en claro y oscuro, la tipografía, los tokens, los componentes y las reglas para usarlos.
Aplica a todas las pantallas: web (`apps/web`), escritorio (`apps/desktop`) y móvil (`apps/mobile`).

> **La fuente única de verdad es el código:** `packages/ui/src/styles/tokens.css` (valores) y
> `packages/ui/src/styles/components.css` (componentes). Este documento los explica; si algo no coincide,
> manda el código. La referencia visual oficial es el sistema «Voxa» en claude.ai
> (`https://claude.ai/artifact/1Daeb4Bq3esqWZXj3KNk3a`, `tokens.json` versión 4). Decisión: [ADR-0014](../adr/0014-sistema-de-diseno.md).

## Marca

- **Nombre comercial:** Voxa, de «vóxel», el píxel en 3D. El proyecto sigue llamándose `cargo-3d` en el
  código y en el repositorio hasta confirmar la marca ante el IMPI (ver ADR-0014).
- **Personalidad:** industrial, precisa y del puerto. El amarillo y negro de la señalización de seguridad de
  las grúas, con la limpieza de un tablero moderno.
- **Voz:** español de México, directo y técnico. «Acomoda 312 cajas respetando el peso por eje», no
  «¡Optimiza tu logística al máximo!». Los botones son verbo + objeto («Optimizar carga», «Liberar plan»).
  Cifras con separador de miles por espacio (`24 310 kg`) y coma decimal (`67,4 m³`).

## Logo: el Vóxel

Un cubo de 2 × 2 × 2 piezas en isométrico con una pieza resaltada (la que el motor acaba de colocar), dentro
de una insignia amarilla. Junto a ella, la palabra VOXA trazada con líneas geométricas de grosor constante.

| Pieza | Componente | Archivo (colores fijos) | Uso |
|---|---|---|---|
| Logotipo | `<Logo />` | `brand/voxa-logo.svg`, `voxa-logo-dark.svg` | Encabezado de la app, documentos, sitio web |
| Insignia | `<Logo variant="badge" />` | `brand/voxa-badge.svg` | Ícono de app (Wails, Expo), avatar; 24 px o más |
| Isotipo | `<Logo variant="mark" />` | `brand/voxa-mark.svg`, `voxa-mark-dark.svg` | Visor 3D, pantallas de carga, marcas de agua |
| Favicon | `<Logo variant="favicon" />` | `brand/voxa-favicon.svg` | Pestaña del navegador y barra de tareas (16–20 px) |

`<Logo />` es un SVG en línea que toma sus colores de los tokens `--logo-*`, así que **cambia de tema solo**.
Los `.svg` de archivo llevan colores fijos: elige el normal sobre fondos claros y el `-dark` sobre fondos
oscuros (o usa las clases `.vx-logo-light` / `.vx-logo-dark`, que se muestran según el tema).

**Reglas**

- **Construcción:** insignia de 64 con radio 14; la palabra mide 40 de alto (5/8 de la insignia) y va
  separada 22 (un tercio).
- **Color:** insignia siempre amarilla con el cubo en tinta y el vóxel en blanco. La palabra en `ink`
  (`#111111`), o en `#F4F4EF` sobre fondos oscuros.
- **Zona de respeto:** un cuarto del alto de la insignia en todos los lados.
- **Tamaño mínimo:** logotipo 20 px de alto; insignia 24 px; abajo de eso, el favicon.
- **No:** girar el cubo, cambiar el color de la insignia, poner el logotipo amarillo sobre blanco,
  agregar sombras, degradados o contornos, ni separar la palabra de la insignia en otra proporción.

## Color: 60 · 30 · 10

Los colores se nombran por **función**, no por tono: `ink` es «tinta» en los dos temas aunque en oscuro sea
casi blanca. Cada token se define una vez con `light-dark(claro, oscuro)`.

| Proporción | Token | Claro | Oscuro | Dónde va |
|---|---|---|---|---|
| **60 %** superficies | `canvas` | `#EFEFEA` | `#0F0F0E` | Fondo de la app, detrás de todo |
| | `card` | `#FFFFFF` | `#1A1A18` | Tarjetas, paneles, buscador, navegación inactiva |
| | `line` | `#DEDED7` | `#2E2E2B` | Bordes y separadores (nunca el único indicio de un control) |
| **30 %** tinta | `ink` | `#111111` | `#F4F4EF` | Texto, cifras, íconos, navegación activa |
| | `ink-muted` | `#5C5C57` | `#A6A69E` | Etiquetas, fechas, texto secundario |
| | `on-ink` | `#FFFFFF` | `#111111` | Texto sobre un fondo `ink` |
| **10 %** amarillo grúa | `accent` | `#FFC400` | `#FFC400` | Insignia, botón principal, filtro activo, progreso |
| | `on-accent` | `#111111` | `#111111` | Texto sobre `accent` (11.8:1) |
| | `accent-soft` | `#FFF6D6` | `#433914` | Fondo de una fila o tarjeta seleccionada |
| | `accent-line` | `#8A6500` | `#FFC400` | Bordes y texto de acento sobre `card` o `canvas` |
| | `track` | `#2A2A27` | `#3A3A36` | Riel de las barras de progreso |
| | `focus` | `#111111` | `#FFC400` | Anillo de foco del teclado |

### Reglas del amarillo

El amarillo da apenas 1.6:1 sobre blanco, así que en el tema claro **solo se usa como relleno**:

1. **Texto encima siempre en tinta** (`on-accent`). Nunca blanco sobre amarillo.
2. **Barras de progreso sobre riel de tinta** (`track`, 9.0:1), como la señalización de seguridad.
3. **Bordes, líneas y texto de acento con `accent-line`**: dorado oscuro en claro (5.3:1), amarillo en oscuro.
4. **Un solo botón principal amarillo por vista.** La navegación activa va en tinta, no en amarillo.
5. En el tema oscuro el amarillo luce de verdad (10.9:1 sobre las tarjetas).

### Estados

Ninguno usa amarillo, para no confundirse con la marca, y **todos llevan texto y punto o ícono, nunca solo
color**. Cada uno tiene tres tokens: el color puro (punto, borde), `-soft` (fondo) e `-ink` (texto, ≥4.5:1 en
ambos temas).

| Estado | Tokens | Utilidad de Tailwind | Ejemplo |
|---|---|---|---|
| Listo | `--status-success*` | `bg-success-soft text-success-ink` | Plan liberado, CG en rango |
| Cargando | `--status-progress*` | `bg-progress-soft text-progress-ink` | El patio está acomodando la carga |
| Cerca del límite | `--status-warning*` | `bg-warning-soft text-warning-ink` | Peso por eje arriba del 90 % |
| Excede | `--status-danger*` | `bg-danger-soft text-danger-ink` | Excede VGM o la NOM-012, CG fuera de rango |
| Demorado | `--status-delayed*` | `bg-delayed-soft text-delayed-ink` | Equipo o carga con retraso |

## Tipografía

**Inter Tight** para la interfaz (compacta y geométrica, como el logotipo) e **IBM Plex Mono** para pesos,
medidas y coordenadas (cifras alineadas). Ambas con licencia OFL, desde Google Fonts (`fonts.css`).

| Estilo | Utilidad | Tamaño / alto de línea | Peso | Uso |
|---|---|---|---|---|
| display | `text-display` | 44 / 48, −0.02em | 600 | Cifras grandes de indicadores |
| h1 | `text-h1` | 28 / 34, −0.01em | 600 | Título de pantalla |
| h2 | `text-h2` | 18 / 24 | 600 | Título de tarjeta o panel |
| body | `text-body` | 14 / 20 | 400 | Texto corrido (el del `body`) |
| label | `text-label` | 13 / 16 | 500 | Botones, navegación, filtros |
| small | `text-small` | 12 / 16 | 400 | Ayudas y etiquetas, en `ink-muted` |
| pill | `text-pill` | 11 / 16 | 500 | Píldoras de estado |
| data | `font-mono text-data` | 12 / 16 | 500 | `24 310 kg · 67,4 m³` |

Las escalas de Tailwind (`text-sm`, `text-xl`…) están **desactivadas**: solo existen las de esta tabla.

## Forma

- **Espaciado** en pasos de 4 px. En Tailwind, `p-1` = 4 px, `p-2` = 8, `p-3` = 12, `p-4` = 16, `p-6` = 24,
  `p-8` = 32. En CSS: `--space-1` … `--space-8`. Íconos y texto: 4; dentro de un control: 8; tarjetas: 12;
  paneles: 16; entre secciones: 24; márgenes de página: 32.
- **Radios:** `rounded-sm` 6 (píldoras, barras), `rounded-md` 10 (tarjetas, botones), `rounded-lg` 16
  (paneles, insignia), `rounded-full`.
- **Sombra:** una sola y casi imperceptible, `shadow-card`. La jerarquía la dan las superficies, no las
  sombras. `shadow-float` (propuesta) es solo para lo que flota: modal, toast, barra masiva.
- **Movimiento** (propuesta): `--motion-fast` 120 ms y `--motion-base` 200 ms, sin rebotes. Todo se apaga con
  `prefers-reduced-motion`.

## Íconos

[Lucide](https://lucide.dev) (`lucide-react`) con trazo de **1.5 px**, esquinas redondeadas y color
`currentColor` (por regla, `ink`). Tres tamaños: **14** (dentro de botones), **16** (junto al texto) y
**18** (indicadores, toasts, modales).

```tsx
import { Icon, appIcons } from "@cargo-3d/ui";

<Icon icon={appIcons.weight} />              // decorativo: se oculta a lectores de pantalla
<Icon icon={appIcons.view3d} label="Ver en 3D" /> // solo, sin texto al lado: lleva nombre
```

`appIcons` (en `packages/ui/src/icons.ts`) es el **vocabulario de íconos del dominio**: un dibujo por concepto
(`container`, `dryVan`, `weight`, `axleLoad`, `centerOfGravity`, `stacking`, `unloadOrder`, `optimize`,
`releasePlan`…). Un concepto nuevo se agrega ahí, no en la pantalla, para que «peso» no termine con tres
dibujos distintos. `statusIcons` da el ícono de cada estado.

**Íconos de las apps:** el favicon de la web es `brand/voxa-favicon.svg`. Los de escritorio (Wails, `.ico`) y
móvil (Expo, `.png` de 1024 px) se generarán a partir de `brand/voxa-badge.svg` cuando existan esas apps.

## Temas claro y oscuro

- Sin atributo, se sigue la preferencia del sistema operativo.
- `data-theme="light"` o `data-theme="dark"` en `<html>` (o en cualquier elemento) lo fuerzan.
  Se pueden anidar: un panel oscuro dentro de una página clara funciona.
- Como los tokens cambian solos, **casi nunca hace falta la variante `dark:`** de Tailwind. Si una pantalla la
  necesita mucho, falta un token.

## Componentes

Todos están en `components.css` con el prefijo `vx-` y solo usan tokens. Los componentes React están en
`packages/ui/src/` y son **solo de presentación**: reciben datos por propiedades y avisan con funciones
(`onConfirm`, `onDismiss`…). No guardan estado de negocio, no llaman a la API ni saben de WebSockets.

| Componente | Clases | React | Estado |
|---|---|---|---|
| Logo | `.vx-logo-light` / `-dark` | `Logo` | Oficial |
| Botón | `.vx-btn` + `--primary` `--secondary` `--ghost` `--sm` | `Button` | Oficial |
| Botón destructivo | `.vx-btn--danger` | `Button variant="danger"` | Propuesta |
| Navegación | `.vx-nav`, `.vx-nav__item`, `--active`, `.vx-search`, `.vx-round` | — | Oficial |
| Píldora de estado | `.vx-pill--success/progress/warning/danger/delayed` | `StatusBadge` | Oficial |
| Estado en tiempo real | `.vx-pill--live` | `StatusBadge live announce` | Propuesta |
| Filtro | `.vx-chip`, `--active` / `aria-pressed` | `FilterChip` | Oficial |
| Indicador | `.vx-kpi` y sus partes, `.vx-delta--up/down` | — | Oficial |
| Tarjeta de carga | `.vx-card`, `--selected`, `.vx-progress` | — | Oficial |
| Panel | `.vx-panel` | — | Oficial |
| Modal de confirmación | `.vx-dialog` y sus partes | `ConfirmDialog` | Propuesta |
| Notificación | `.vx-toaster`, `.vx-toast--<estado>` | `Toaster`, `Toast` | Propuesta |
| Acciones masivas | `.vx-bulkbar`, `.vx-bulkbar__action` | `BulkActionBar`, `BulkAction` | Propuesta |

«Propuesta» quiere decir que el componente se diseñó con los tokens oficiales pero todavía no está en el
sistema de claude.ai; hay que agregarlo allá para que diseño y código no se separen.

### Botones

- `primary` (amarillo): **uno solo por vista**, la acción principal.
- `secondary` (tinta): acción fuerte pero no principal.
- `ghost` (tarjeta con borde): acciones de apoyo. Es el valor por omisión de `Button`.
- `danger` (propuesta): **solo** para confirmar algo destructivo dentro de un `ConfirmDialog`.
- Deshabilitado con el atributo `disabled`, nunca solo con color. El anillo de foco nunca se quita.

### Badges de estado en tiempo real

`<StatusBadge status="progress" live announce />`: `live` hace latir el punto mientras el estado sigue
cambiando; `announce` pone `role="status"` para que el cambio se lea en voz. Usa `announce` solo en el estado
principal de la vista, no en cada fila. El componente no sabe de dónde viene el estado: la pantalla se lo pasa.

### Modal de confirmación

`<ConfirmDialog>` usa el `<dialog>` nativo (foco atrapado, Esc, velo). Pregunta corta como título
(«¿Eliminar el plan?»), qué pasará y si se puede deshacer como descripción, y verbo + objeto en el botón
(«Eliminar plan»). Con `tone="danger"` el foco empieza en «Cancelar».

### Notificaciones

`<Toaster>` va una sola vez en la raíz de la app, desde el inicio (aunque vacío), y adentro los `<Toast>`.
El color del estado va en el borde izquierdo y el ícono; el texto siempre en tinta. Los `danger` se anuncian
de inmediato (`role="alert"`). Cuánto duran y la cola de avisos los decide la app.

### Barra de acciones masivas

Aparece cuando hay elementos seleccionados (`count > 0`). Va en tinta para no competir con el botón
amarillo. Muestra el contador («3 piezas»), las acciones y «Quitar selección». Las destructivas abren un
`ConfirmDialog`.

## Cómo se usa en una app

```css
/* apps/web/src/main.css */
@import "@cargo-3d/ui/styles.css";
```

```tsx
import { Button, Logo, StatusBadge, appIcons } from "@cargo-3d/ui";

<Logo />
<Button variant="primary" icon={appIcons.optimize}>Optimizar carga</Button>
<StatusBadge status="warning">Eje trasero al 93 %</StatusBadge>
```

Utilidades de Tailwind disponibles: los colores del sistema (`bg-canvas`, `bg-card`, `text-ink`,
`text-ink-muted`, `bg-accent`, `border-line`, `bg-success-soft`…), la escala tipográfica, `font-sans`,
`font-mono`, el espaciado, los radios y `shadow-card`. **Las paletas de Tailwind (`bg-blue-500`…) no existen**:
si hace falta un color nuevo, se agrega como token con su función.

## Vista previa

La vista previa es una página HTML que muestra todo el sistema en claro y oscuro, lado a lado. No necesita
ninguna app ni servidor. Los comandos se corren en una terminal abierta en la raíz del proyecto (`C:\cargo-3d`).

**1. Compilar los estilos** (genera `packages/ui/preview/dist/preview.css`, que no se sube a git):

```bash
pnpm --filter @cargo-3d/ui build:preview
```

Debe terminar con `Done in …ms`.

**2. Abrir la página en el navegador:**

| Sistema | Comando |
|---|---|
| Windows (PowerShell) | `start packages\ui\preview\index.html` |
| macOS | `open packages/ui/preview/index.html` |
| Linux | `xdg-open packages/ui/preview/index.html` |

También se puede abrir con doble clic en el archivo. Los botones **Sistema / Claro / Oscuro** de arriba cambian
el tema de la página, y **Abrir modal** muestra el modal de confirmación.

**Mientras se cambian estilos**, conviene el modo que recompila solo cada vez que se guarda un archivo:

```bash
pnpm --filter @cargo-3d/ui preview
```

Se deja corriendo, se recarga el navegador (F5) después de cada cambio y se detiene con Ctrl+C.

**Si algo no se ve bien:**

| Síntoma | Causa | Solución |
|---|---|---|
| La página sale sin estilos | No se compiló o se borró `preview/dist/` | Repetir el paso 1 |
| Faltan los íconos | No están instaladas las dependencias (los íconos se toman de `node_modules`) | `pnpm install` en la raíz |
| La letra no es Inter Tight | Sin internet no cargan las fuentes de Google | Es normal: se usa la fuente del sistema |

## Reglas para mantenerlo

1. **Ningún color, tamaño, radio ni sombra escrito a mano** en una pantalla: todo sale de un token.
2. Un valor nuevo o distinto se cambia **en `tokens.css` y en el sistema oficial de claude.ai**, en el mismo
   cambio, y se anota aquí.
3. Un componente nuevo va primero en `components.css` (clase `vx-`) y luego, si hace falta, en React.
4. Contraste mínimo: 4.5:1 para texto y 3:1 para controles y bordes que indican algo. Se revisa en los dos temas.
5. Nada se comunica solo con color: estados con texto y punto, variaciones con flecha y signo.
6. Los componentes de `@cargo-3d/ui` no importan nada de la API, del motor ni de un estado global.

## Pendientes

- Agregar al sistema oficial los componentes propuestos (modal, toast, barra masiva, botón destructivo,
  estado en vivo) y los tokens `on-danger`, `scrim`, `shadow-float` y `--motion-*`.
- Componentes que faltan: campo de texto, tabla, panel del visor 3D.
- Móvil: conectar los mismos tokens a React Native (por ejemplo con NativeWind) en el sprint 3.2.
- Íconos de app para Wails y Expo a partir de la insignia.
- Servir las fuentes desde la app (sin internet) para el escritorio del patio.
