# 25 nov 2026 · Casos simulados y ajustes al flujo

**Sesión:** Fuera del plan · **Área:** docs, contratos

## Qué se hizo

- **Sin puntos de control** (ADR-0018): se borró la skill `/punto-de-control`, la regla 9 de `CLAUDE.md`, sus
  menciones en el plan, `prompts.md` (17 sesiones de cierre), `flujo-claude-code.md`, `/cerrar-sprint`, `/sesion` y la
  guía de Go, y las dos secciones de punto de control de la bitácora del 21 de noviembre.
- **Casos de uso simulados** (ADR-0020): 29 casos en `docs/src/casos/<categoría>/<subcategoría>/cs-NNN-<tema>.md`
  (marítimo, carretera, restricciones, límites y operación), con equipo, piezas, resultado esperado y su cuenta.
  Incluyen uno justo por espacio (CS-001), uno justo por peso (CS-002), varios donde algo no cabe, camionetas, rabón y
  tortón simulados (CS-017 a CS-020), el mismo viaje por distintos tipos de camino (CS-011, CS-019), carga
  sobredimensionada en carretera con permiso de la NOM-040 (CS-028) y una plataforma escalonada de dos tramos (CS-029).
- **Pruebas con cargas reales:** procedimiento en `docs/src/casos/pruebas-reales.md` para el periodo de pruebas de la
  v0.2.0 (registro privado por empresa en `docs/datos/pruebas/`, casos `Validado en real` o `Corregido`).
- **Equipos personalizados** (ADR-0019 y `conceptos/equipos-personalizados.md`): equipo estándar, plantilla por clase
  y equipo de la empresa con versiones; se construye en las sesiones 2.3 y 2.4.
- **Plan y prompts:** las puertas de las fases 0 y 1 usan casos simulados; la sesión 1.1.5 convierte los casos en
  fixtures (`engine/fixtures/casos/`); el sprint 2.6 pasa a «Despliegue y pruebas con cargas reales». Se quitó el nombre
  del cliente previsto de todos los documentos que se suben a git.
- **Bitácora:** un archivo por tema (`AAAA-MM-DD-tema.md`), índice en `bitacora/README.md` y agrupada por mes en el
  sitio. Los dos archivos por día se separaron en siete sin perder contenido; se agregó la entrada del 23 de noviembre.
- **Imágenes de los ejemplos:** `equipo.svg` en los 27 casos de `contracts/examples/`, generados con
  `scripts/contratos/render-examples.mjs` (`pnpm examples:render`): misma escala, ángulo y paleta Voxa, claro y oscuro.

- **Ejemplos y catálogo:** cuatro ejemplos nuevos con medidas de ficha (`10dv-epal-pallets`, `45pw-reefer-epal-pallets`,
  `40fr-collapsible-steel` y `20bk-coffee-sacks`, ruta nueva `maritime/bulk/`). Se retiraron del catálogo el tanque ISO 20'
  y el chasis portacontenedor (su altura de 1 219 mm pasó a la nota de la plataforma ISO 40'). Cada equipo de
  `contracts/equipment.json` lleva ahora `examples` o `exampleNote`, y `equipos.md` una tabla «Ejemplos del proyecto por equipo».
- **Corrección:** cuatro ejemplos corregidos el 23-nov (`20rf`, `40rh`, `40fr-oog`, `40pf`) tenían descripciones de más de
  200 caracteres; el validador de reglas físicas no lo detecta, el esquema sí.

## Decisiones

- Imágenes generadas en lugar de fotos: el repositorio es público y sin licencia.
- Sin ejemplo y con la razón anotada: 40' reefer estándar (sin puerta publicada), flat rack 20' plegable (mismas medidas
  que el 20' fijo), plataforma ISO 20' (sin chasis de 20' con fuente) y plataforma plana 40' (sin altura de cama).
- El periodo de pruebas con cargas reales empieza en la v0.2.0 (MVP).
- P-007 se cerró como cancelado (la regla de pendientes no permite borrar ids).
- Los CSV sintéticos del patio se movieron a `docs/datos/sprint 1.1.1/respaldo-patio-sintetico/`; la hoja de
  medición de `docs/plantillas/patio/` queda para el periodo de pruebas.

## Verificación

- `mdbook build docs` sin errores.
- `node scripts/contratos/check-examples.mjs`: 27 de 27; ajv (temporal, fuera del repo): los 27 pares cumplen el esquema.
- Renders revisados en Edge headless.
- No se cambió código de Go.

## Pendientes

- P-001, P-003 y P-007: cerrados. P-004 y P-013: pasan a la sesión 2.6. P-019: agregado al registro (faltaba).
  P-022 y P-023: nuevos.
