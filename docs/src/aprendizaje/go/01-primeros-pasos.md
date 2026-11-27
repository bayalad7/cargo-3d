# 01 — Primeros pasos: herramientas, módulos y el servidor `/health`

## Instalar Go

Descarga el instalador de go.dev/dl (versión 1.27 o superior), o usa `winget install GoLang.Go` en Windows
y `brew install go` en macOS. Instala una sola herramienta: `go`.
Con ella compilas, pruebas, formateas y administras dependencias; no hace falta otra herramienta.

## Comandos que usarás a diario

| Comando | Qué hace |
|---|---|
| `go run ./backend/cmd/api` | Compila en una carpeta temporal y ejecuta el servidor |
| `go test ./backend/...` | Corre todas las pruebas (`...` = todas las subcarpetas) |
| `go vet ./backend/...` | Busca errores que compilan pero casi seguro están mal, como un `Printf` con argumentos que no cuadran |
| `gofmt -w backend` | Formatea el código (un solo formato para todo Go) |
| `go doc ./backend/internal/engine` | Muestra la documentación de un paquete |
| `go mod tidy` | Reescribe `go.mod` y `go.sum`: agrega lo que el código importa y quita lo que ya no |
| `go mod download` | Descarga las dependencias de `go.mod` sin cambiarlo |
| `go build -o api ./backend/cmd/api` | Produce el binario ejecutable |

## Los archivos del backend

### `go.work` (raíz)

Agrupa los módulos de Go del monorepo. Por eso puedes correr `go test ./backend/...` desde la raíz.

### `backend/go.mod`

Define el módulo `github.com/bayalad7/cargo-3d/backend`. Esa ruta es el prefijo de todos los imports
internos, por ejemplo `github.com/bayalad7/cargo-3d/backend/internal/engine`. Ese prefijo es un nombre, no una
descarga: Go lo resuelve en tu disco (ver [Módulos e imports](01b-modulos-e-imports.md)).

### `backend/cmd/api/main.go`

El punto de entrada. Fíjate en cuatro ideas nuevas, todas comentadas en el archivo:

1. `package main` + `func main()`: así Go sabe que esto produce un ejecutable.
2. `os.Getenv("PORT")`: la configuración viene de variables de entorno, no del código.
3. `go func() { ... }()`: una **goroutine**; el servidor atiende peticiones en paralelo mientras
   `main` espera la orden de apagado.
4. `defer stop()`: deja programada la llamada a `stop()` para cuando `main` termine, ya sea con `return` o por un
   *panic*. La excepción es `os.Exit`: termina el programa en el acto y **no** corre los `defer`.

### `backend/internal/httpapi/health.go`

Una ruta mínima, `GET /health`:

```go
type HealthResponse struct {
	Status        string `json:"status"`
	EngineVersion string `json:"engineVersion"`
}
```

- `type ... struct` define un tipo con campos ([capítulo 02](02-structs-y-tipos.md)). Los campos con mayúscula se exportan.
- Lo que va entre comillas invertidas son **etiquetas**: le dicen al codificador JSON qué nombre usar.

```go
if err := json.NewEncoder(w).Encode(resp); err != nil {
	http.Error(w, "no se pudo generar la respuesta", http.StatusInternalServerError)
}
```

Así se ven los errores en Go: la función devuelve un `error`; si no es `nil`, algo falló y se maneja ahí mismo.

### `backend/internal/httpapi/health_test.go`

Las pruebas viven junto al código, en archivos `_test.go`. `httptest` simula peticiones sin red:
la prueba corre en milisegundos.

## Ejercicio

1. Corre `go test ./backend/...`. Deben pasar 3 pruebas.
2. Corre `go run ./backend/cmd/api` y abre `http://localhost:8080/health` en el navegador.
3. Agrega un campo `Name string \`json:"name"\`` a `HealthResponse` con el valor `"cargo-3d"`,
   actualiza `contracts/openapi.yaml` primero (regla 10) y luego el código. Corre las pruebas.
4. Rompe a propósito el `if Version == ""` de `doc_test.go` cambiándolo a `!=`, lee el error y regrésalo.
