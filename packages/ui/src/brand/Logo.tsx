import { useId, type CSSProperties, type SVGProps } from "react";

/**
 * LogoVariant son las piezas del logo Vóxel.
 *
 * - `logo`: insignia + palabra VOXA (encabezado, documentos). Mínimo 20 px de alto.
 * - `badge`: solo la insignia (ícono de app, avatar). Mínimo 24 px.
 * - `mark`: el cubo sin cuadro (visor 3D, pantallas de carga, marcas de agua).
 * - `favicon`: versión simplificada para 16–20 px (pestaña del navegador).
 */
export type LogoVariant = "logo" | "badge" | "mark" | "favicon";

/** LogoProps son las propiedades de `Logo`. */
export interface LogoProps extends Omit<SVGProps<SVGSVGElement>, "children" | "viewBox"> {
  /** Qué pieza dibujar. Por omisión `logo`. */
  variant?: LogoVariant;
  /** Alto en píxeles; el ancho sale de la proporción. Por omisión 32. */
  height?: number;
  /** Nombre para lectores de pantalla. Pasa `""` si el logo es decorativo. */
  title?: string;
}

/** Cubo de 2 × 2 × 2 en isométrico (coordenadas del sistema oficial, lienzo de 64). */
const CUBE_TOP = [
  "32,6.22 41.59,11.75 32,17.29 22.41,11.75",
  "20.3,12.96 29.89,18.5 20.3,24.04 10.71,18.5",
  "43.7,12.96 53.29,18.5 43.7,24.04 34.11,18.5",
];
const CUBE_VOXEL = "32,19.71 41.59,25.25 32,30.79 22.41,25.25";
const CUBE_FRONT = [
  "9.65,20.32 19.25,25.86 19.25,36.93 9.65,31.39",
  "9.65,33.82 19.25,39.36 19.25,50.43 9.65,44.89",
  "21.35,27.07 30.95,32.61 30.95,43.68 21.35,38.14",
  "21.35,40.57 30.95,46.11 30.95,57.18 21.35,51.64",
];
const CUBE_SIDE = [
  "33.05,32.61 42.65,27.07 42.65,38.14 33.05,43.68",
  "33.05,46.11 42.65,40.57 42.65,51.64 33.05,57.18",
  "44.75,25.86 54.35,20.32 54.35,31.39 44.75,36.93",
  "44.75,39.36 54.35,33.82 54.35,44.89 44.75,50.43",
];

/** Colores de cada cara, leídos de los tokens `--logo-*` de tokens.css. */
interface CubeColors {
  top: string;
  voxel: string;
  front: string;
  side: string;
}

/** Valor claro de cada token, por si la página todavía no carga tokens.css. */
const FALLBACK: Readonly<Record<string, string>> = {
  "--logo-badge": "#FFC400",
  "--logo-badge-top": "#64500B",
  "--logo-badge-side": "#3A300E",
  "--logo-badge-front": "#111111",
  "--logo-badge-voxel": "#FFFFFF",
  "--logo-word": "#111111",
  "--logo-mark-top": "#4D4D4A",
  "--logo-mark-front": "#111111",
  "--logo-mark-side": "#2F2F2E",
  "--logo-mark-voxel": "#FFC400",
};

const token = (name: string): string => `var(${name}, ${FALLBACK[name] ?? "currentColor"})`;

const fill = (name: string): CSSProperties => ({ fill: token(name) });

function Cube({ top, voxel, front, side }: CubeColors) {
  return (
    <>
      <g style={fill(top)}>
        {CUBE_TOP.map((points) => (
          <polygon key={points} points={points} />
        ))}
      </g>
      <polygon points={CUBE_VOXEL} style={fill(voxel)} />
      <g style={fill(front)}>
        {CUBE_FRONT.map((points) => (
          <polygon key={points} points={points} />
        ))}
      </g>
      <g style={fill(side)}>
        {CUBE_SIDE.map((points) => (
          <polygon key={points} points={points} />
        ))}
      </g>
    </>
  );
}

const BADGE_CUBE: CubeColors = {
  top: "--logo-badge-top",
  voxel: "--logo-badge-voxel",
  front: "--logo-badge-front",
  side: "--logo-badge-side",
};

function Badge() {
  return (
    <>
      <rect width="64" height="64" rx="14" style={fill("--logo-badge")} />
      <g transform="translate(10 10) scale(.6875)">
        <Cube {...BADGE_CUBE} />
      </g>
    </>
  );
}

/** Palabra VOXA: trazos geométricos de grosor constante, recortados a 40 de alto. */
function Wordmark({ clipId }: { clipId: string }) {
  return (
    <g transform="translate(86 12) scale(.8333)">
      <clipPath id={clipId}>
        <rect x="0" y="4" width="170" height="40" />
      </clipPath>
      <g
        clipPath={`url(#${clipId})`}
        fill="none"
        strokeWidth="8"
        strokeLinejoin="miter"
        strokeMiterlimit="10"
        style={{ stroke: token("--logo-word") }}
      >
        <path d="M1.9 -2 L18 46 L34.1 -2" />
        <circle cx="62" cy="24" r="16" />
        <path d="M88 -2 L124 50 M124 -2 L88 50" />
        <path d="M127.9 50 L144 2 L160.1 50 M136 32 H152" />
      </g>
    </g>
  );
}

function Favicon() {
  return (
    <>
      <rect width="64" height="64" rx="14" style={fill("--logo-badge")} />
      <g transform="translate(8 8) scale(.75)">
        <polygon points="32,6.35 53.06,18.5 32,30.65 10.94,18.5" style={fill("--logo-badge-top")} />
        <polygon points="32,19.85 41.36,25.25 32,30.65 22.64,25.25" style={fill("--logo-badge-voxel")} />
        <polygon points="9.77,20.52 30.83,32.67 30.83,56.98 9.77,44.83" style={fill("--logo-badge-front")} />
        <polygon points="33.17,32.67 54.23,20.52 54.23,44.83 33.17,56.98" style={fill("--logo-badge-side")} />
      </g>
    </>
  );
}

const VIEWBOX: Readonly<Record<LogoVariant, { width: number; height: number }>> = {
  logo: { width: 222, height: 64 },
  badge: { width: 64, height: 64 },
  mark: { width: 64, height: 64 },
  favicon: { width: 64, height: 64 },
};

/**
 * Logo dibuja el logo oficial de Voxa, el Vóxel: un cubo de 2 × 2 × 2 en isométrico con
 * la pieza que el motor acaba de colocar resaltada.
 *
 * Cómo funciona: es un SVG en línea cuyos colores salen de los tokens `--logo-*`, así que
 * la palabra y el isotipo cambian solos entre claro y oscuro; la insignia es amarilla en
 * los dos temas. Para archivos sueltos (favicon, documentos) están los `.svg` de `brand/`.
 *
 * Por qué así: un solo componente evita elegir a mano la versión clara u oscura y
 * respeta las reglas del logo (no girar, no recolorear, proporción fija entre insignia
 * y palabra). Tamaño mínimo: 20 px de alto el logotipo y 24 px la insignia; abajo de
 * eso, usa `variant="favicon"`.
 */
export function Logo({ variant = "logo", height = 32, title = "Voxa", ...rest }: LogoProps) {
  // El id del recorte debe ser único si hay varios logos en la página.
  const clipId = `voxa-word-${useId().replace(/[^a-zA-Z0-9_-]/g, "")}`;
  const box = VIEWBOX[variant];
  const width = Math.round((height * box.width) / box.height);
  const decorative = title === "";

  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox={`0 0 ${box.width} ${box.height}`}
      width={width}
      height={height}
      {...(decorative ? { "aria-hidden": true } : { role: "img", "aria-label": title })}
      {...rest}
    >
      {variant === "logo" && (
        <>
          <Badge />
          <Wordmark clipId={clipId} />
        </>
      )}
      {variant === "badge" && <Badge />}
      {variant === "favicon" && <Favicon />}
      {variant === "mark" && (
        <Cube top="--logo-mark-top" voxel="--logo-mark-voxel" front="--logo-mark-front" side="--logo-mark-side" />
      )}
    </svg>
  );
}
