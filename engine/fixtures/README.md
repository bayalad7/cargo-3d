# Casos de prueba del motor

Cargas de ejemplo en JSON, **independientes del lenguaje** del motor.
Cada archivo sigue `contracts/engine.schema.json` e incluye lo que debe cumplirse:

- sin traslapes entre piezas,
- todas las piezas dentro del equipo,
- peso máximo respetado,
- ocupación mínima esperada (línea base).

Las pruebas revisan **propiedades**, no posiciones exactas, así que sirven para cualquier motor.

Se llenan en la fase 1:

- `casos/` — casos de uso simulados de `docs/src/casos/` (sesión 1.1.5, ADR-0020).
- `br/` — instancias académicas BR1 a BR15 (sprint 1.4).
