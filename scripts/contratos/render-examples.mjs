#!/usr/bin/env node
/**
 * Dibuja una imagen de cada equipo de `contracts/examples/` y la guarda como `equipo.svg` en la carpeta del ejemplo.
 *
 * Cómo funciona: lee el `input.json` de cada ejemplo, toma su tramo, su puerta y su margen (overhang), y dibuja el
 * equipo en proyección isométrica con la MISMA escala, el mismo ángulo y la misma paleta (Voxa) para todos. Así un
 * 20' se ve más chico que una caja de 53', y la puerta queda siempre del mismo lado (x = largo, ADR-0015).
 *
 * Por qué así: el repositorio es público y sin licencia; las fotos de fabricantes no se pueden redistribuir. Un render
 * generado desde los datos no tiene ese problema y siempre coincide con las medidas del ejemplo.
 *
 * Uso: `node scripts/contratos/render-examples.mjs` (o `pnpm examples:render`). Vuelve a correrlo si cambia un ejemplo.
 */
import { readFileSync, readdirSync, writeFileSync } from "node:fs";
import { join, relative, sep } from "node:path";

const ROOT = join(import.meta.dirname, "..", "..", "contracts", "examples");
const CATALOG = JSON.parse(readFileSync(join(import.meta.dirname, "..", "..", "contracts", "equipment.json"), "utf8"));

/** Px por mm: la misma escala para todos, calculada para que quepa una caja de 53'. */
const SCALE = 0.028;
const WIDTH = 760;
const HEIGHT = 600;
const BASELINE = 498; // Píxel donde queda el punto más bajo del dibujo: todos «pisan» el mismo suelo.
const COS30 = Math.cos(Math.PI / 6);

/** Equipo del catálogo que representa cada ejemplo (por nombre de carpeta). */
const EQUIPMENT_ID = {
  "20dv-heavy": "20dv", "20hc-tall-racks": "20hc", "40dv-gma-pallets": "40dv", "40hc-tall-pallets": "40hc",
  "45hc-gma-pallets": "45hc", "45pw-epal-pallets": "45pw", "20rf-epal-pallets": "20rf", "40rh-epal-pallets": "40rh",
  "20fr-in-gauge": "20fr", "20ot-door": "20ot", "20ot-roof": "20ot", "40fr-oog": "40fr", "40ot-door": "40ot",
  "40ot-roof": "40ot", "40pf-wide-load": "40pf", "48dv-mixed": "48dry", "53dv-stacked": "53dry",
  "53dv-contract-coverage": "53dry", "48fb-steel": "48fb", "53fb-long-pipes": "48fb", "48rv-pallets": "48rv",
  "53rv-pallets": "53rv", "10dv-epal-pallets": "10dv", "45pw-reefer-epal-pallets": "45pw-reefer",
  "40fr-collapsible-steel": "40fr-collapsible", "20bk-coffee-sacks": "20bk",
};
/** Nombres de los equipos que no están en el catálogo o que el catálogo agrupa. */
const NAME_OVERRIDE = { "48cn-coils": "Caja Conestoga 48'", "53fb-long-pipes": "Plataforma flatbed 53'", "48fb-steel": "Plataforma flatbed 48'" };

/** Tipo de dibujo según la ruta del ejemplo. */
function kindOf(path) {
  const [mode, type, name] = path.split(sep);
  if (mode === "maritime" && type === "dry") return "container";
  if (mode === "maritime" && type === "reefer") return "reefer";
  if (mode === "maritime" && type === "bulk") return "bulk";
  if (name.includes("fr-")) return "flatrack";
  if (name.endsWith("ot-roof")) return "opentop-roof";
  if (name.includes("ot-")) return "opentop";
  if (name.includes("pf-")) return "platform";
  if (type === "curtain") return "curtain";
  if (type === "dry-van") return "van";
  if (type === "reefer-van") return "van-reefer";
  if (type === "flatbed") return "flatbed";
  throw new Error(`Tipo desconocido: ${path}`);
}

/** «16003» → «16 003», como pide la voz del proyecto. */
const fmt = (n) => String(n).replace(/\B(?=(\d{3})+(?!\d))/g, " ");

/** Proyección isométrica: x = largo (hacia la puerta), y = alto, z = ancho. Devuelve [u, v] sin desplazar. */
const project = ([x, y, z]) => [(x - z) * COS30 * SCALE, ((x + z) * 0.5 - y) * SCALE];

/**
 * Lienzo: junta polígonos y líneas en coordenadas 3D y al final los desplaza para que el dibujo quede centrado y
 * apoyado en BASELINE.
 */
class Scene {
  constructor() { this.shapes = []; }
  poly(points, cls) { this.shapes.push({ type: "poly", points: points.map(project), cls }); }
  line(a, b, cls) { this.shapes.push({ type: "line", points: [project(a), project(b)], cls }); }
  /** Caja alineada a los ejes con sus tres caras visibles (arriba, lado z = máx. y frente x = máx.). */
  box([x0, y0, z0], [dx, dy, dz], { top = "top", side = "side", front = "front" } = {}) {
    const [x1, y1, z1] = [x0 + dx, y0 + dy, z0 + dz];
    this.poly([[x0, y1, z0], [x1, y1, z0], [x1, y1, z1], [x0, y1, z1]], top);
    this.poly([[x0, y0, z1], [x1, y0, z1], [x1, y1, z1], [x0, y1, z1]], side);
    this.poly([[x1, y0, z0], [x1, y1, z0], [x1, y1, z1], [x1, y0, z1]], front);
  }
  /** Rueda en el costado visible (plano z), como polígono de 24 lados. */
  wheel(x, y, z, r) {
    const pts = Array.from({ length: 24 }, (_, i) => {
      const t = (i / 24) * 2 * Math.PI;
      return [x + r * Math.cos(t), y + r * Math.sin(t), z];
    });
    this.poly(pts, "wheel");
    this.poly(pts.map(([px, py, pz]) => [x + (px - x) * 0.45, y + (py - y) * 0.45, pz]), "hub");
  }
  render() {
    const all = this.shapes.flatMap((s) => s.points);
    const minU = Math.min(...all.map((p) => p[0]));
    const maxU = Math.max(...all.map((p) => p[0]));
    const maxV = Math.max(...all.map((p) => p[1]));
    const du = (WIDTH - (maxU - minU)) / 2 - minU;
    const dv = BASELINE - maxV;
    const pt = ([u, v]) => `${(u + du).toFixed(1)},${(v + dv).toFixed(1)}`;
    return this.shapes.map((s) => s.type === "poly"
      ? `<polygon class="${s.cls}" points="${s.points.map(pt).join(" ")}"/>`
      : `<polyline class="${s.cls}" points="${s.points.map(pt).join(" ")}"/>`).join("\n    ");
  }
}

/** Puerta en la cara x = L, centrada a lo ancho, con dos hojas y barras de cierre. */
function door(scene, L, W, d) {
  const z0 = (W - d.widthMm) / 2;
  const z1 = z0 + d.widthMm;
  scene.poly([[L, 0, z0], [L, d.heightMm, z0], [L, d.heightMm, z1], [L, 0, z1]], "door");
  scene.line([L, 0, (z0 + z1) / 2], [L, d.heightMm, (z0 + z1) / 2], "door-line");
  for (const f of [0.18, 0.38, 0.62, 0.82]) {
    const z = z0 + f * d.widthMm;
    scene.line([L, 80, z], [L, d.heightMm - 80, z], "door-bar");
  }
}

/** Escotilla de carga en el techo (contenedor bulk), como círculo en el plano y = H. */
function hatch(scene, x, y, z, r) {
  const pts = Array.from({ length: 24 }, (_, i) => {
    const t = (i / 24) * 2 * Math.PI;
    return [x + r * Math.cos(t), y, z + r * Math.sin(t)];
  });
  scene.poly(pts, "opening");
}

/** Corrugado de los contenedores marítimos en el costado visible. */
function corrugation(scene, L, H, W, step = 300) {
  for (let x = step; x < L - 100; x += step) scene.line([x, 60, W], [x, H - 60, W], "rib");
}

/** Envolvente punteada: lo que se permite rebasar (carga sobredimensionada) o, con `always`, el alto útil sin paredes. */
function envelope(scene, s, always = false) {
  const o = { topMm: 0, leftMm: 0, rightMm: 0, ...s.overhang };
  if (!always && !o.topMm && !o.leftMm && !o.rightMm) return;
  const [x0, x1] = [0, s.lengthMm];
  const [z0, z1] = [-o.leftMm, s.widthMm + o.rightMm];
  const y1 = s.heightMm + o.topMm;
  const top = [[x0, y1, z0], [x1, y1, z0], [x1, y1, z1], [x0, y1, z1], [x0, y1, z0]];
  scene.shapes.push({ type: "line", points: top.map(project), cls: "env" });
  for (const [x, z] of [[x1, z0], [x1, z1], [x0, z1]]) scene.line([x, 0, z], [x, y1, z], "env");
  scene.line([x1, 0, z1], [x0, 0, z1], "env");
  scene.line([x1, 0, z0], [x1, 0, z1], "env");
}

/** Ruedas de un semirremolque (tándem atrás) y patas de apoyo adelante. */
function roadChassis(scene, L, W, deck = 1200) {
  const r = 520;
  const y = -deck + r;
  scene.box([0, -260, 0], [L, 260, W], { top: "frame", side: "frame", front: "frame" });
  for (const x of [L - 1500, L - 2800]) scene.wheel(x, y, W, r);
  scene.line([L * 0.3, -260, W - 150], [L * 0.3, -deck, W - 150], "leg");
}

function draw(kind, input) {
  const c = input.container;
  const s = c.sections[0];
  const [L, H, W] = [s.lengthMm, s.heightMm, s.widthMm];
  const scene = new Scene();
  switch (kind) {
    case "container":
    case "reefer":
      scene.box([0, 0, 0], [L, H, W], kind === "reefer" ? { top: "top", side: "white", front: "white" } : {});
      if (kind === "container") corrugation(scene, L, H, W);
      else for (const y of [H * 0.33, H * 0.66]) scene.line([150, y, W], [L - 150, y, W], "rib");
      if (c.door) door(scene, L, W, c.door);
      break;
    case "bulk":
      scene.box([0, 0, 0], [L, H, W]);
      corrugation(scene, L, H, W);
      for (const f of [0.2, 0.5, 0.8]) hatch(scene, L * f, H, W / 2, 250);
      if (c.door) door(scene, L, W, c.door);
      break;
    case "opentop":
    case "opentop-roof":
      scene.box([0, 0, 0], [L, H, W], { top: kind === "opentop-roof" ? "opening" : "tarp" });
      corrugation(scene, L, H, W);
      for (let x = 600; x < L; x += 1200) scene.line([x, H, 0], [x, H, W], "bow");
      if (c.door) door(scene, L, W, c.door);
      break;
    case "flatrack":
      envelope(scene, s);
      scene.box([0, -300, 0], [L, 300, W]);
      scene.box([0, 0, 0], [220, H, W], { top: "post", side: "post", front: "post" });
      scene.box([L - 220, 0, 0], [220, H, W], { top: "post", side: "post", front: "post" });
      break;
    case "platform":
      envelope(scene, s);
      scene.box([0, -650, 0], [L, 650, W]);
      for (const x of [0, L - 180]) scene.box([x, -650, 0], [180, 650, W], { top: "post", side: "post", front: "post" });
      break;
    case "flatbed":
      roadChassis(scene, L, W, 1478);
      scene.box([0, -120, 0], [L, 120, W]);
      envelope(scene, s, true); // Alto útil hasta los 4 250 mm de la NOM-012.
      break;
    case "curtain":
      roadChassis(scene, L, W);
      scene.box([0, 0, 0], [L, H, W], { side: "tarp", top: "top", front: "front" });
      for (let x = 400; x < L; x += 400) scene.line([x, 40, W], [x, H - 40, W], "fold");
      break;
    case "van":
    case "van-reefer":
      if (kind === "van-reefer") scene.box([-450, H - 1300, W * 0.2], [450, 1000, W * 0.6], { top: "unit", side: "unit", front: "unit" });
      roadChassis(scene, L, W);
      scene.box([0, 0, 0], [L, H, W], kind === "van-reefer" ? { top: "top", side: "white", front: "white" } : {});
      for (const y of [140, H - 140]) scene.line([0, y, W], [L, y, W], "rib");
      if (c.door) door(scene, L, W, c.door);
      break;
    default:
      throw new Error(kind);
  }
  return scene.render();
}

function svg(path, kind, input) {
  const name = path.split(sep).at(-1);
  const c = input.container;
  const s = c.sections[0];
  const eq = CATALOG.equipment.find((e) => e.id === EQUIPMENT_ID[name]);
  let title = NAME_OVERRIDE[name] ?? eq?.name ?? name;
  if (kind === "opentop-roof") title += " · carga por el techo";
  else if (kind === "opentop") title += " · carga por la puerta";
  const code = eq?.isoCode ? `ISO ${eq.isoCode} · ` : "";
  const parts = [`Interior ${fmt(s.lengthMm)} × ${fmt(s.widthMm)} × ${fmt(s.heightMm)} mm`];
  if (c.door) parts.push(`Puerta ${fmt(c.door.widthMm)} × ${fmt(c.door.heightMm)} mm`);
  parts.push(`Carga útil ${fmt(c.maxPayloadKg)} kg`);
  const bar = 1000 * SCALE * COS30; // 1 m sobre el eje x, ya proyectado en horizontal.
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${WIDTH} ${HEIGHT}" width="${WIDTH}" height="${HEIGHT}" role="img" aria-label="${title}">
  <title>${title}</title>
  <style>
    .bg{fill:#efefea}.card{fill:#ffffff;stroke:#deded7}
    .top{fill:#ffffff}.side{fill:#deded7}.front{fill:#c9c9c1}.white{fill:#f7f7f3}
    .top,.side,.front,.white,.tarp,.opening,.post,.frame,.unit,.wheel,.hub,.door{stroke:#111111;stroke-width:1.2;stroke-linejoin:round}
    .tarp{fill:#fff6d6}.opening{fill:#ffc400}.post{fill:#5c5c57}.frame{fill:#5c5c57}.unit{fill:#5c5c57}
    .wheel{fill:#111111}.hub{fill:#a6a69e}.door{fill:#ffc400}
    .rib,.bow,.fold{fill:none;stroke:#5c5c57;stroke-width:.8}
    .door-line,.door-bar{fill:none;stroke:#111111;stroke-width:1}.leg{fill:none;stroke:#111111;stroke-width:3}
    .env{fill:none;stroke:#8a6500;stroke-width:1.2;stroke-dasharray:5 4}
    .title{font:600 20px 'Inter Tight',system-ui,sans-serif;fill:#111111}
    .sub{font:400 13px 'Inter Tight',system-ui,sans-serif;fill:#5c5c57}
    .data{font:400 12px 'IBM Plex Mono',ui-monospace,monospace;fill:#111111}
    .scale{stroke:#111111;stroke-width:2}
    @media (prefers-color-scheme: dark){
      .bg{fill:#0f0f0e}.card{fill:#1a1a18;stroke:#2e2e2b}
      .top{fill:#3a3a36}.side{fill:#2e2e2b}.front{fill:#242421}.white{fill:#4a4a45}
      .top,.side,.front,.white,.tarp,.opening,.post,.frame,.unit,.wheel,.hub,.door{stroke:#f4f4ef}
      .tarp{fill:#433914}.post,.frame,.unit{fill:#a6a69e}.wheel{fill:#0f0f0e}.hub{fill:#5c5c57}
      .rib,.bow,.fold{stroke:#a6a69e}.door-line,.door-bar,.leg,.scale{stroke:#111111}.leg,.scale{stroke:#f4f4ef}
      .env{stroke:#ffc400}.title,.data{fill:#f4f4ef}.sub{fill:#a6a69e}
    }
  </style>
  <rect class="bg" width="${WIDTH}" height="${HEIGHT}"/>
  <rect class="card" x="12" y="12" width="${WIDTH - 24}" height="${HEIGHT - 24}" rx="12"/>
  <text class="title" x="36" y="52">${title}</text>
  <text class="sub" x="36" y="74">${code}Ejemplo ${path.split(sep).join("/")}</text>
  <g>
    ${draw(kind, input)}
  </g>
  <line class="scale" x1="36" y1="${BASELINE + 34}" x2="${(36 + bar).toFixed(1)}" y2="${BASELINE + 34}"/>
  <text class="sub" x="${(42 + bar).toFixed(1)}" y="${BASELINE + 38}">1 m · misma escala en todos los ejemplos</text>
  <text class="data" x="36" y="${HEIGHT - 32}">${parts.join("  ·  ")}</text>
</svg>
`;
}

/** Todas las carpetas con `input.json`, salvo los casos inválidos (no son equipos). */
function findCases(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const path = join(dir, entry.name);
    if (entry.isDirectory()) return entry.name === "invalid" ? [] : findCases(path);
    return entry.name === "input.json" ? [dir] : [];
  });
}

let count = 0;
for (const dir of findCases(ROOT)) {
  const rel = relative(ROOT, dir);
  const input = JSON.parse(readFileSync(join(dir, "input.json"), "utf8"));
  writeFileSync(join(dir, "equipo.svg"), svg(rel, kindOf(rel), input));
  count += 1;
}
console.log(`${count} imágenes escritas (equipo.svg en cada ejemplo).`);
