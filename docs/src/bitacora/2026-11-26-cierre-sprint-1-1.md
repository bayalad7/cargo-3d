# 26 nov 2026 · Cierre del sprint 1.1

**Sesión:** 1.1.6 · **Área:** motor, docs

## Qué se hizo

Se revisó el checklist del sprint 1.1 (Modelos y contrato del motor) contra el repositorio. Todas las tareas tienen
evidencia:

| Tarea | Evidencia |
|---|---|
| `contracts/engine.schema.json` | Contrato v0.1.0, 27 ejemplos válidos y 17 inválidos en `contracts/examples/` (ADR-0015, 0016 y 0017; sesión 1.1.1) |
| Tipos en `internal/engine`, en mm y kg `int64` | `input.go`, `container.go`, `cargo_item.go`, `orientation.go`, `placement.go`, `load_plan.go`; ida y vuelta de los 27 `input.json` y los 27 `output.json`; `schema_test.go` compara campos, obligatorios y rangos con el esquema (sesiones 1.1.2 y 1.1.3) |
| Validación con errores claros, nunca `panic` | `decode.go`, `validate.go`, `errors.go`; los 17 inválidos dan cada uno su categoría; fuzz de 4,8 millones de entradas sin `panic` (ADR-0023; sesión 1.1.4) |
| Casos simulados como fixtures | 36 fixtures de 26 casos en `engine/fixtures/casos/`, revisados por `fixtures_test.go` (sesión 1.1.5) |
| Guía de Go, capítulos 02 a 04 | `aprendizaje/go/02-structs-y-tipos.md`, `03-metodos-y-receptores.md` y `04-errores-como-valores.md` |

## Mediciones

- **Pruebas:** 37 pruebas de primer nivel y 425 en total contando subpruebas, todas pasan
  (`go test ./backend/... -count=1 -v`). Una prueba fuzz (`FuzzDecodeInput`) con los ejemplos como semilla.
- **Cobertura:** `internal/engine` 96,4 %, `internal/httpapi` 85,7 % (`go test ./backend/... -cover`).
- **Tiempo de las pruebas:** menos de 1 s para `internal/engine`.
- **Código:** 20 archivos Go en `internal/engine`, unas 3 260 líneas con comentarios (nivel máximo).
- **Ocupación y tiempos del motor:** no aplica todavía; el algoritmo llega en el sprint 1.3.

## Decisiones

- ADR-0022 (Go explicado solo con Go), ADR-0023 (lectura y validación de la entrada) y ADR-0024 (PostgreSQL sin
  contraseña en desarrollo) se tomaron durante el sprint.

## Revisión con Benjamín

- **Comodidad con Go:** cómodo; entiende el código y los capítulos 02 a 04. El nivel de comentarios sigue en máximo por
  ahora y se puede bajar un poco más adelante.
- **Ritmo:** el sprint se hizo entre el 20 y el 26 de noviembre; el ritmo del plan le parece realista y se sigue igual.

## Verificación

- `gofmt -l backend` sin salida; `go vet ./backend/...` y `go test ./backend/...` sin errores.
- `packages/ui` no tiene pruebas, solo `lint` (`tsc`), y no se corrió porque `node` no está instalado en esta
  computadora; el paquete no cambió en el sprint.

## Pendientes

- Ninguno nuevo. Siguen abiertos P-004, P-013, P-019 a P-024 y P-028; agendados P-010, P-012, P-015, P-016 y P-026.
