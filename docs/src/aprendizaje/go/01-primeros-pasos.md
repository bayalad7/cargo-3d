# 01 — Primeros pasos: herramientas, módulos y el servidor `/health`

## Instalar Go

Descarga el instalador de go.dev/dl (versión 1.27 o superior), o usa `winget install GoLang.Go` en Windows
y `brew install go` en macOS. Instala una sola herramienta: `go`.
Con ella compilas, pruebas, formateas y administras dependencias. No hay equivalente a Composer aparte.

## Comandos que usarás a diario

| Comando | Qué hace | Equivalente aproximado en PHP |
|---|---|---|
| `go run ./backend/cmd/api` | Compila y ejecuta el servidor | `php -S` / `php artisan serve` |
| `go test ./backend/...` | Corre todas las pruebas (`...` = todas las subcarpetas) | `phpunit` |
| `go vet ./backend/...` | Busca errores comunes | PHPStan |
| `gofmt -w backend` | Formatea el código (formato único, sin debates) | PHP-CS-Fixer |
| `go doc ./backend/internal/engine` | Muestra la documentación de un paquete | phpDocumentor |
| `go mod tidy` | Ajusta las dependencias de `go.mod` | `composer install` |
| `go build -o api ./backend/cmd/api` | Produce el binario ejecutable | (PHP no compila) |

## Los archivos del backend

### `go.work` (raíz)

Agrupa los módulos de Go del monorepo. Por eso puedes correr `go test ./backend/...` desde la raíz.

### `backend/go.mod`

Define el módulo `github.com/bayalad7/cargo-3d/backend`. Esa ruta es el prefijo de todos los imports
internos, por ejemplo `github.com/bayalad7/cargo-3d/backend/internal/engine`.

### `backend/cmd/api/main.go`

El punto de entrada. Fíjate en cuatro ideas nuevas, todas comentadas en el archivo:

1. `package main` + `func main()`: así Go sabe que esto produce un ejecutable.
2. `os.Getenv("PORT")`: la configuración viene de variables de entorno, no del código.
3. `go func() { ... }()`: una **goroutine**; el servidor atiende peticiones en paralelo mientras
   `main` espera la orden de apagado.
4. `defer stop()`: se ejecuta al salir de la función, como un `finally`.

### `backend/internal/httpapi/health.go`

Una ruta mínima, `GET /health`:

```go
type HealthResponse struct {
	Status        string `json:"status"`
	EngineVersion string `json:"engineVersion"`
}
```

- `type ... struct` es como una clase sin métodos. Los campos con mayúscula son públicos.
- Lo que va entre comillas invertidas son **etiquetas**: le dicen al codificador JSON qué nombre usar.

```go
if err := json.NewEncoder(w).Encode(resp); err != nil {
	http.Error(w, "no se pudo generar la respuesta", http.StatusInternalServerError)
}
```

Así se ven los errores en Go: la función devuelve un `error`; si no es `nil`, algo falló y se maneja ahí mismo.
No hay `try/catch`.

### `backend/internal/httpapi/health_test.go`

Las pruebas viven junto al código, en archivos `_test.go`. `httptest` simula peticiones sin red:
la prueba corre en milisegundos.

## Ejercicio

1. Corre `go test ./backend/...`. Deben pasar 3 pruebas.
2. Corre `go run ./backend/cmd/api` y abre `http://localhost:8080/health` en el navegador.
3. Agrega un campo `Name string \`json:"name"\`` a `HealthResponse` con el valor `"cargo-3d"`,
   actualiza `contracts/openapi.yaml` primero (regla 10) y luego el código. Corre las pruebas.
4. Rompe a propósito el `if Version == ""` de `doc_test.go` cambiándolo a `!=`, lee el error y regrésalo.
