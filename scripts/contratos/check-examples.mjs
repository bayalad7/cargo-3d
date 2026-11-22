/**
 * Revisa las reglas físicas de los ejemplos de `contracts/examples/`.
 *
 * Uso (desde la raíz del repo):  node scripts/contratos/check-examples.mjs
 *
 * Cómo funciona: busca cada carpeta con `input.json` y `output.json` y comprueba lo que JSON Schema no puede
 * expresar (traslapes, apoyo, peso, motivos, métricas…). Termina con código 1 si algún ejemplo tiene problemas.
 *
 * Por qué así: es provisional y no usa dependencias, para no adelantar el sprint 1.4. NO valida contra el esquema
 * (eso pide ajv y un ADR). Cuando exista `packages/validator` (sesión 1.4.2), este archivo se borra (P-010).
 */
import { readFileSync, readdirSync } from "node:fs";
import { dirname, join, relative } from "node:path";
import { fileURLToPath } from "node:url";

const EXAMPLES_DIR = join(dirname(fileURLToPath(import.meta.url)), "..", "..", "contracts", "examples");

/** Tolerancia para comparar ocupaciones, que se guardan con 4 decimales. */
const TOLERANCE = 1e-4;

/** Medidas de una pieza ya girada. Las letras van sobre x, z, y: en ESE orden (ADR-0015). */
const ROTATIONS = {
  LWH: (l, w, h) => ({ dx: l, dz: w, dy: h }),
  WLH: (l, w, h) => ({ dx: w, dz: l, dy: h }),
  LHW: (l, w, h) => ({ dx: l, dz: h, dy: w }),
  HLW: (l, w, h) => ({ dx: h, dz: l, dy: w }),
  WHL: (l, w, h) => ({ dx: w, dz: h, dy: l }),
  HWL: (l, w, h) => ({ dx: h, dz: w, dy: l }),
};

const readJson = (path) => JSON.parse(readFileSync(path, "utf8"));
const rotate = (item, orientation) => ROTATIONS[orientation](item.lengthMm, item.widthMm, item.heightMm);
const allowedOf = (item) => item.allowedOrientations ?? Object.keys(ROTATIONS);

/** Cuánto se enciman dos intervalos; 0 si solo se tocan. */
const overlap = (aMin, aSize, bMin, bSize) => Math.max(0, Math.min(aMin + aSize, bMin + bSize) - Math.max(aMin, bMin));

/** Todas las carpetas, a cualquier profundidad, que tienen un `input.json`. */
function findCases(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const path = join(dir, entry.name);
    if (entry.isDirectory()) return findCases(path);
    return entry.name === "input.json" ? [dir] : [];
  });
}

/** Motivo que le toca a una pieza no colocada, o `null` si depende del hueco o del tiempo. */
function expectedReason(item, container, section, margin, loadedKg) {
  const inside = allowedOf(item).map((o) => rotate(item, o)).filter(({ dx, dy, dz }) =>
    dx <= section.lengthMm &&
    dz <= section.widthMm + margin.leftMm + margin.rightMm &&
    dy <= section.heightMm + margin.topMm);
  if (inside.length === 0) return "EXCEEDS_DIMENSIONS";
  const { door } = container;
  if (door && !inside.some(({ dy, dz }) => dz <= door.widthMm && dy <= door.heightMm)) return "EXCEEDS_DOOR";
  if (item.weightKg > container.maxPayloadKg) return "EXCEEDS_MAX_WEIGHT";
  if (item.weightKg > container.maxPayloadKg - loadedKg) return "WEIGHT_LIMIT_REACHED";
  return null;
}

/** Devuelve la lista completa de problemas de un ejemplo (vacía si está bien). */
function checkCase(input, plan) {
  const problems = [];
  const { container } = input;
  const items = new Map(input.items.map((item) => [item.id, item]));
  const sections = new Map(container.sections.map((section) => [section.id, section]));
  const marginOf = (section) => ({ topMm: 0, leftMm: 0, rightMm: 0, ...section.overhang });

  const placedCount = new Map();
  const seenUnits = new Set();
  let loadedKg = 0;
  let loadedVolume = 0;
  let outOfGauge = 0;

  if (plan.containerId !== container.id) problems.push("containerId no es el de la entrada");

  plan.placements.forEach((p, i) => {
    const item = items.get(p.itemId);
    const section = sections.get(p.sectionId);
    if (!item) return problems.push(`#${i}: la pieza ${p.itemId} no está en la entrada`);
    if (!section) return problems.push(`#${i}: el tramo ${p.sectionId} no existe`);
    const margin = marginOf(section);

    const unit = `${p.itemId}/${p.instance}`;
    if (seenUnits.has(unit)) problems.push(`#${i}: la unidad ${unit} está repetida`);
    if (p.instance >= item.quantity) problems.push(`#${i}: instance ${p.instance} fuera de la cantidad`);
    seenUnits.add(unit);

    if (!allowedOf(item).includes(p.orientation)) problems.push(`#${i}: orientación no permitida`);
    const { dx, dy, dz } = rotate(item, p.orientation);
    if (dx !== p.dxMm || dy !== p.dyMm || dz !== p.dzMm) problems.push(`#${i}: medidas distintas a su orientación`);

    const top = section.yMm + section.heightMm;
    const right = section.zMm + section.widthMm;
    if (p.xMm < section.xMm || p.xMm + dx > section.xMm + section.lengthMm ||
        p.yMm < section.yMm || p.yMm + dy > top + margin.topMm ||
        p.zMm < section.zMm - margin.leftMm || p.zMm + dz > right + margin.rightMm) {
      problems.push(`#${i}: fuera del tramo y de su margen`);
    }
    const exceeds = p.yMm + dy > top || p.zMm < section.zMm || p.zMm + dz > right;
    if (exceeds !== (p.outOfGauge === true)) problems.push(`#${i}: outOfGauge debería ser ${exceeds}`);
    if (exceeds) outOfGauge++;

    const { door } = container;
    if (door && (dz > door.widthMm || dy > door.heightMm)) problems.push(`#${i}: no pasa por la puerta`);

    if (p.yMm > section.yMm) {
      const base = plan.placements
        .filter((q) => q !== p && q.yMm + q.dyMm === p.yMm)
        .reduce((sum, q) => sum + overlap(p.xMm, dx, q.xMm, q.dxMm) * overlap(p.zMm, dz, q.zMm, q.dzMm), 0);
      if (base < dx * dz) problems.push(`#${i}: flotando (apoyo ${((100 * base) / (dx * dz)).toFixed(0)} %)`);
    }

    loadedKg += item.weightKg;
    loadedVolume += dx * dy * dz;
    placedCount.set(p.itemId, (placedCount.get(p.itemId) ?? 0) + 1);
  });

  plan.placements.forEach((a, i) => plan.placements.slice(i + 1).forEach((b, k) => {
    if (overlap(a.xMm, a.dxMm, b.xMm, b.dxMm) && overlap(a.yMm, a.dyMm, b.yMm, b.dyMm) &&
        overlap(a.zMm, a.dzMm, b.zMm, b.dzMm)) problems.push(`#${i} y #${i + 1 + k} se enciman`);
  }));

  if (loadedKg > container.maxPayloadKg) problems.push("excede maxPayloadKg");

  const unplacedCount = new Map();
  for (const { itemId, quantity, reason } of plan.unplaced) {
    unplacedCount.set(itemId, (unplacedCount.get(itemId) ?? 0) + quantity);
    const item = items.get(itemId);
    if (!item) { problems.push(`unplaced: ${itemId} no está en la entrada`); continue; }
    // Con varios tramos habría que probar cada uno; el contrato 0.1.0 solo acepta uno.
    const [section] = container.sections;
    const expected = expectedReason(item, container, section, marginOf(section), loadedKg);
    const ok = expected ? reason === expected : ["NO_SPACE_AVAILABLE", "TIME_LIMIT_REACHED"].includes(reason);
    if (!ok) problems.push(`${itemId}: motivo ${reason}, se esperaba ${expected ?? "hueco o tiempo"}`);
  }
  for (const { id, quantity } of input.items) {
    if ((placedCount.get(id) ?? 0) + (unplacedCount.get(id) ?? 0) !== quantity) {
      problems.push(`${id}: colocadas + no colocadas ≠ cantidad`);
    }
  }

  const m = plan.metrics;
  const capacity = container.sections.reduce((sum, s) => sum + s.lengthMm * s.widthMm * s.heightMm, 0);
  const totalUnplaced = [...unplacedCount.values()].reduce((sum, n) => sum + n, 0);
  const exact = { placedCount: plan.placements.length, unplacedCount: totalUnplaced, outOfGaugeCount: outOfGauge,
    loadedWeightKg: loadedKg, loadedVolumeMm3: loadedVolume };
  for (const [name, value] of Object.entries(exact)) {
    if (m[name] !== value) problems.push(`metrics.${name} es ${m[name]}, se esperaba ${value}`);
  }
  if (Math.abs(m.volumeUtilization - loadedVolume / capacity) > TOLERANCE) problems.push("volumeUtilization incorrecta");
  if (Math.abs(m.weightUtilization - loadedKg / container.maxPayloadKg) > TOLERANCE) problems.push("weightUtilization incorrecta");
  return problems;
}

let failed = 0;
const cases = findCases(EXAMPLES_DIR).sort();
for (const dir of cases) {
  const name = relative(EXAMPLES_DIR, dir).replaceAll("\\", "/");
  const problems = checkCase(readJson(join(dir, "input.json")), readJson(join(dir, "output.json")));
  if (problems.length > 0) failed++;
  console.log(`${problems.length > 0 ? "FALLA" : "ok   "} ${name}`);
  problems.forEach((problem) => console.log(`        ${problem}`));
}
console.log(`\n${cases.length - failed} de ${cases.length} ejemplos correctos.`);
// Sin ejemplos también es una falla: significaría que se movió la carpeta.
process.exit(failed > 0 || cases.length === 0 ? 1 : 0);
