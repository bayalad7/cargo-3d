# 21 nov 2026 · Recorrido de pendientes

**Sesión:** 1.1.1 · **Área:** docs

Fuera del plan: se revisaron uno por uno P-001 a P-018.

## Qué se hizo

- **Cerrados:** P-005, P-006, P-009, P-011, P-014, P-017 y P-018. **Agendados** en el prompt de su sesión: P-002, P-008, P-010,
  P-012, P-015 y P-016. **Abiertos:** P-001, P-003, P-004, P-007, P-013 y P-019 (nuevo).
- 8 ejemplos más (23 en total) y 17 entradas inválidas. `scripts/contratos/check-examples.mjs` rescata la revisión física
  que vivía en una carpeta temporal.
- Dos conceptos nuevos, formatos para el patio y sesión 2.1.6; la de cierre del sprint 2.1 pasó a ser la 2.1.7.

## Decisiones

- Catálogo del cliente: la tara está restada dos veces; **no** son límites de carretera, como se había concluido en la
  1.1.1. Filas mal nombradas: se renombran y en los reefers se usan fichas públicas.
- El plan cambió: los equipos con varios tramos entran en el sprint 2.1 (fase 2).
- El registro de pendientes tiene una tabla nueva, «Agendados», que no cuenta como hecho.

## Correcciones a lo anterior

- `40fr-oog` traía un margen de 200 mm por lado, que rebasaba los 2 600 mm de la NOM-012: ahora 196.
- P-018 no era un error, y su arreglo propuesto habría roto la skill `/sesion`.

## Verificación

- `check-examples.mjs`: 23 de 23 ejemplos; 12 errores sembrados a propósito, todos detectados.
- ajv (temporal, fuera del repo): los 23 pares cumplen el esquema; los 15 de `invalid/schema` fallan y los 2 de
  `invalid/semantic` pasan el esquema, como se espera.

## Pendientes

- P-014: el primer intento de publicar la v5 de Voxa lo bloqueó una regla de permisos; se publicó después, a petición de
  Benjamín, desde `docs/datos/voxa-v5/`. Cerrado.
- P-019: copiar `docs/datos/` y `.env` antes de cambiar de computadora.
