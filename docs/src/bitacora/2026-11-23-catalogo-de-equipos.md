# 23 nov 2026 · Catálogo de equipos y correcciones de P-004

**Sesión:** 1.1.5 (parte de datos, adelantada) · **Área:** contratos, docs

## Qué se hizo

- Catálogo canónico de 28 equipos en `contracts/equipment.json` y su versión legible con foto real en
  `docs/src/catalogos/equipos.md` (preparados con una investigación externa con fuentes; commits `398880b` y `272a604`).
- Revisión fila por fila de ambos archivos: tres diferencias resueltas con Benjamín (alto interior de la 48rv a `null`,
  volumen publicado del 45' pallet wide y bruto de la 53rv en blanco).
- Plataforma plana 40' con la ficha Fruehauf PM-HT2-40W (la altura de cama sigue sin publicar).
- Ejemplos de `contracts/examples/` corregidos con fuente (auditoría G1 a G5): reefers con la serie Hapag-Lloyd,
  cajas refrigeradas con el máximo Alzper, alto derivado 2 383 de la plataforma ISO 40', 400 mm como límite interno
  del flat rack, puerta del 45' pallet wide seco y 30 000 kg de cajas secas como capacidad estructural.
- `catalogos/equipos.md` publicado en el sitio; hoja de medición §4 y §5 con los valores vigentes.
- Revisión de los CSV del patio: solo tenían casos sintéticos copiados de los ejemplos (base de la decisión del
  2026-11-25, ADR-0020).

## Decisiones

- `contracts/equipment.json` es la fuente única del catálogo; ningún valor cambia sin fuente.
- Los reefers de los ejemplos usan la serie Hapag completa y su alto se etiqueta como alto interior, no línea roja.
- Catálogo del cliente convertido en privado (`docs/datos/`), sin cambiar el catálogo público.

## Verificación

- `python -m json.tool contracts/equipment.json` sin errores; `node scripts/contratos/check-examples.mjs`: 23 de 23.

## Pendientes

- P-002: cerrado. P-004 y P-013: actualizados. P-020 y P-021: nuevos.
