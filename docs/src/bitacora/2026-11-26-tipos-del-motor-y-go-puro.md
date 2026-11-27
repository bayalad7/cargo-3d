# 26 nov 2026 · Tipos del motor y Go explicado solo con Go

**Sesión:** 1.1.2 · **Área:** motor, docs

## Qué se hizo

- **Tipos del motor** en `backend/internal/engine/`, siguiendo `contracts/engine.schema.json` campo por campo:
  - `input.go`: `EngineInput`.
  - `container.go`: `Container`, `Section`, `Overhang`, `OverhangValues` y `Door`, con `Volume()` (del tramo y del equipo,
    sin el margen) y `Section.Margin()`.
  - `cargo_item.go`: `CargoItem` con `Volume()` de una unidad e `IsStackable()`, y el tipo `Orientation` con sus 6
    constantes (girar las medidas queda para la 1.1.3).
  - Los campos donde el 0 o el `false` significan algo (`Stackable`, `TareKG`, `MaxLoadOnTopKG`, `Overhang.*`) y los
    objetos opcionales (`Door`, `Overhang`) son punteros: `nil` es «no vino». Los valores por omisión se calculan con
    métodos, no al leer el JSON.
  - En los comentarios, sobre el código: quién calcula `MaxPayloadKG` y por qué el motor no conoce la NOM-012, y hacia
    dónde se redondea cada valor (ADR-0016).
- **Pruebas:** de tabla para `Volume()`, `Margin()` e `IsStackable()` con los extremos del contrato, y
  `contract_examples_test.go`, que pasa los 27 `input.json` de JSON a Go y de regreso en modo estricto, revisa que no
  quede contenido después del objeto y compara por valores.
- **Capítulo 02** de la guía de Go (`aprendizaje/go/02-structs-y-tipos.md`): structs, valores cero, qué se copia y qué se
  comparte, `int64` contra `float64`, etiquetas `json`, «no vino» contra «vino en cero», métodos, tipos con nombre y las
  dos decisiones de `MaxPayloadKG` y redondeo.
- **Revisión (`/revisar`) y correcciones:** se corrigió lo que la guía decía mal de `Orientation` (un texto escrito tal
  cual sí compila), se aclaró que `Volume()` solo es confiable con entrada validada (se desborda sin aviso), la prueba de
  ida y vuelta ahora detecta contenido de más y se corrigió «sprint 2.1.5» por «sesión 2.1.5».
- **Go explicado solo con Go (ADR-0022):** Benjamín pidió quitar toda comparación con otros lenguajes. Se reescribieron
  los capítulos 00, 01, 01b y 02, los comentarios del backend, la skill `/explicar`, la regla 6 de `CLAUDE.md` y los
  prompts que pedían comparar. Al reescribir se corrigieron errores reales de Go en lo que ya estaba escrito: un receptor
  por valor sí puede modificar datos compartidos, `defer` no corre si se llama a `os.Exit`, `0.1 + 0.2` escrito como
  constante da `0.3`, y las equivalencias de herramientas que no lo eran.
- Las pruebas usan `new(int64(400))` y `new(false)` (Go 1.26 en adelante); se quitaron las funciones auxiliares.

## Decisiones

- ADR-0022: el código, los comentarios, la guía, las skills y los prompts explican Go solo en términos de Go y lo que se
  afirma se comprueba antes. Sustituye la parte de la comparación del punto 3 del ADR-0009.
- Opción A para los `output.json`: la ida y vuelta de `LoadPlan` se prueba en la 1.1.3, cuando exista el tipo.
- Se dejaron sin cambiar los registros históricos que mencionan otro lenguaje (ADR-0002 y bitácoras anteriores).

## Verificación

- `gofmt -l backend` sin salida; `go vet ./backend/...` y `go test ./backend/...` sin errores; `mdbook build docs` sin
  errores. Las 27 entradas de `contracts/examples/` van y vuelven sin cambios.
- Los comportamientos de Go que se afirman (copias superficiales, `nil`, constantes, `slices.Clone`, `new(expr)`) se
  probaron con programas aparte en el directorio temporal.
- No se usó Docker.

## Pendientes

- P-025: orientaciones vacías y validación antes de cualquier cálculo (agendado en la 1.1.3 y la 1.1.4).
- Hallazgos de la revisión que no se corrigieron: la descripción vacía no regresa con `omitempty` (se documentó en el
  código) y una prueba que solo multiplica números; ninguno bloquea nada.
