# 21 nov 2026 · Contrato del motor

**Sesión:** 1.1.1 · **Área:** contratos

Sprint 1.1, sesión 1 del plan. Modelo: `opus` / `medium`. Se hizo con datos **sintéticos** (todavía no hay cargas
reales en `docs/datos/`); la sesión 1.1.5 ajusta el esquema cuando lleguen.

## Qué se hizo

- Dos investigaciones con `/buscador` (respaldadas en `docs/datos/`, que no se sube a git): medidas, puertas, tara y
  carga útil de contenedores ISO, cajas secas y plataformas (NOM-012-SCT-2-2017), pallets EPAL y GMA, y los
  equipos que faltaban (45' HC, 20' RF, open top, flat rack, caja seca 48', caja refrigerada).
- `contracts/engine.schema.json` (JSON Schema 2020-12, versión 0.1.0): `EngineInput` y `LoadPlan`, camelCase,
  milímetros y kilogramos enteros, `description` en español en cada campo.
- 15 ejemplos sintéticos en `contracts/examples/<modo>/<tipo>/<caso>/{input,output}.json`: contenedores secos,
  reefer, open top (por puerta y por techo), flat rack con carga sobredimensionada, caja seca 48' y 53', caja
  refrigerada 53' y plataformas 48' y 53'. Entre todos aparecen los 6 motivos de rechazo.
- `contracts/examples/README.md`: lista de equipos y de partes del contrato que faltan por cubrir, y paso a paso
  para agregar un caso. `contracts/README.md`: convenciones, versiones y cómo validar.
- ADR-0015 (coordenadas, orientaciones y tramos), ADR-0016 (qué rechaza el motor, peso efectivo y redondeo) y
  ADR-0017 (carga sobredimensionada y open top), enlazados en `SUMMARY.md`.
- Ajustes a `prompts.md`: 1.1.2 (valores por omisión en la ida y vuelta JSON), 1.1.3 (de dónde se mira la izquierda)
  y 1.1.4 (qué es error de validación y qué no).
- Revisión con `/revisar`: sin bloqueantes; se corrigieron 5 hallazgos (ocupación de volumen mayor a 1 con OOG,
  margen solo en tramos abiertos, `20ot-roof` centrado, valores por omisión y límites de tamaño).

## Decisiones

- **Coordenadas:** `y` hacia arriba, mano derecha como Three.js; origen en la esquina del fondo, abajo, a la
  izquierda mirando del fondo hacia la puerta (ADR-0015).
- **Tramos:** el equipo es una lista `sections`; en 0.1.0 solo se acepta uno (`maxItems: 1`).
- **Lo que no cabe no es error:** sale en `unplaced` con motivo. La validación de la 1.1.4 es solo para entradas
  mal formadas (ADR-0016).
- **Peso efectivo (`maxPayloadKg`):** lo calcula la API o el catálogo con la NOM-012; el motor solo recibe el número.
- **Redondeo conservador:** capacidades del equipo hacia abajo, piezas hacia arriba, mínimo 1 kg (ADR-0016).
- **Carga sobredimensionada:** `Section.overhang` (arriba y a cada lado, nunca al frente ni atrás), `outOfGauge` por
  pieza y `outOfGaugeCount`; el margen es solo para tramos sin puerta (ADR-0017).
- **Open top:** dos entradas del catálogo, «por puerta» y «por techo» (tramo centrado con `xMm` 109 y `zMm` 49).
- **Jerarquía de equipos** (modo / tipo / subtipo): vive solo en el catálogo, no en el motor.
- El contrato queda en **0.1.0** porque todavía no tenía commit; no se tocó `engine.Version` ni Go.

## Verificación

- `gofmt`, `go vet` y `go test ./backend/...` sin errores (no hubo cambios de Go).
- ajv 8 en modo estricto (temporal, fuera del repo): los 15 pares validan contra el esquema y 8 entradas inválidas
  fallan como se espera. Un script aparte revisó las reglas físicas (traslapes, tramo y margen, puerta, apoyo, peso,
  motivos de rechazo y métricas).

## Pendientes

Detalle y cómo retomarlos en [Pendientes](../pendientes.md).

- P-001: cargas reales y catálogo del cliente (Benjamín, antes de 1.1.5).
- P-002: revisar el catálogo del cliente fila por fila (1.1.5).
- P-003: medidas de rabón y tortón del patio (Benjamín).
- P-004: confirmar valores supuestos (margen OOG, altura útil del 20' RF, cargas útiles, cama de plataforma).
- P-005: Tabla A-2 de la NOM-012 contra el DOF (antes de 2.1).
- P-006: casos con datos listos (open top 40', flat rack 20', caja refrigerada 48').
- P-007: punto de control del contrato (Benjamín, antes de 1.1.2).
- P-008 y P-009: validación de la entrada y ejemplos inválidos (1.1.4).
- P-010: validador independiente con las reglas físicas (1.4.2).
- P-011: conceptos de coordenadas y de equipos y normas (1.1.3 y 1.1.5).
- P-012: varios tramos por equipo (decisión futura).
- P-013: fichas de equipos sin datos (Benjamín).
- P-014 a P-017: pendientes de la sesión de Voxa que quedaron registrados (componentes «propuesta», conexión de `@cargo-3d/ui`, tokens móviles y Node 26).
- P-018: atributo `arguments` no válido en la skill `/sesion`.

## Registro de pendientes

Al cierre se creó `docs/src/pendientes.md`: desde ahora todo lo abierto se registra ahí con id, responsable, fecha límite y cómo retomarlo, y `/cerrar-sesion`, `/iniciar-sesion` y `/sesion` lo usan.

---
