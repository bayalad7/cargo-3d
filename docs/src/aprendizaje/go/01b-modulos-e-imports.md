# 01b — Módulos e imports: ¿por qué `github.com/...` si el código está en mi disco?

Complemento de [Primeros pasos](01-primeros-pasos.md). Responde una duda que aparece al leer el primer
archivo del backend: en `backend/internal/httpapi/health.go`
y en `backend/cmd/api/main.go` hay imports como este:

```go
import "github.com/bayalad7/cargo-3d/backend/internal/engine"
```

Parece que Go descarga el paquete del repositorio de GitHub. **No es así:** lo toma de tu disco.

## La idea clave: es un nombre, no una URL

En `backend/go.mod` está la línea:

```go
module github.com/bayalad7/cargo-3d/backend
```

Significa: «este módulo **se llama** `github.com/bayalad7/cargo-3d/backend`, y su raíz es la carpeta donde
está este `go.mod`», es decir, `backend/`.

Un **módulo** es un conjunto de paquetes que se versionan juntos (un `go.mod`). Un **paquete** es una
carpeta con archivos `.go` que comparten `package nombre`. La ruta de import de un paquete es
`<nombre del módulo>/<carpeta dentro del módulo>`.

## Qué hace Go al encontrar el import

```text
import "github.com/bayalad7/cargo-3d/backend/internal/engine"
                │
                ▼
¿Con qué empieza la ruta?
   │
   ├─ Con el nombre de un módulo del workspace (go.work → backend/go.mod)
   │     → quita el prefijo → internal/engine → lee backend/internal/engine/   (tu disco)
   │
   ├─ Con un paquete de la librería estándar (net/http, os…)
   │     → lo toma de la instalación de Go
   │
   └─ Con otra cosa (github.com/jackc/pgx/v5…)
         → caché ~/go/pkg/mod (se descarga de internet solo la primera vez)
```

Para nuestro caso:

| Paso | Valor |
|---|---|
| Import | `github.com/bayalad7/cargo-3d/backend/internal/engine` |
| Nombre del módulo (`go.mod`) | `github.com/bayalad7/cargo-3d/backend` |
| Lo que sobra | `internal/engine` |
| Carpeta que se compila | `backend/internal/engine/` (`doc.go`) |

No hay red de por medio. Puedes comprobarlo: apaga el wifi y corre `go build ./backend/...`; compila igual.
Y si cambias algo en `backend/internal/engine/doc.go`, el cambio se ve de inmediato en `httpapi`, sin
`push` ni descargas.

## Comparación con PHP (Composer y PSR-4)

Es exactamente el autoload PSR-4 de Composer:

```json
// composer.json
"autoload": {
    "psr-4": { "Bayalad7\\Cargo3d\\": "src/" }
}
```

```php
use Bayalad7\Cargo3d\Engine\Version; // PHP NO va a Packagist:
                                     // cambia el prefijo por src/ y abre src/Engine/Version.php
```

| PHP (Composer) | Go |
|---|---|
| `composer.json` → `"name": "bayalad7/cargo-3d"` | `go.mod` → `module github.com/bayalad7/cargo-3d/backend` |
| Prefijo PSR-4 `Bayalad7\Cargo3d\` → `src/` | Prefijo del módulo → carpeta del `go.mod` (`backend/`) |
| `use Bayalad7\Cargo3d\Engine\...` | `import "github.com/bayalad7/cargo-3d/backend/internal/engine"` |
| Paquete de terceros: `composer require` → `vendor/` | Paquete de terceros: `go get` → `~/go/pkg/mod` |
| `composer.lock` | `go.sum` (huellas de cada dependencia descargada) |

La diferencia principal: en PHP el namespace y el nombre del paquete de Packagist son cosas separadas; en
Go **son lo mismo**. El nombre del módulo sirve como prefijo de imports y, si alguien de fuera lo pide, como
dirección para descargarlo.

## Entonces, ¿por qué tiene forma de URL?

Es una **convención** para que los nombres sean únicos en el mundo, como `com.empresa.proyecto` en Java.
Nadie más puede publicar un módulo `github.com/bayalad7/...` porque esa cuenta es tuya.

La URL solo se usa cuando **otro proyecto** quiere tu código: si alguien escribe
`go get github.com/bayalad7/cargo-3d/backend`, Go sí va a ese repositorio. Dentro de nuestro monorepo, nunca.

El módulo podría llamarse `cargo3d/backend` y todo compilaría igual; usamos la ruta completa porque es la
práctica estándar y no cuesta nada.

## Qué pinta `go.work`

`go.work` en la raíz dice `use ./backend`: le avisa a Go que ese módulo vive en esa
carpeta, para poder correr `go test ./backend/...` desde la raíz. Si algún día hay otro módulo de Go en el
monorepo y uno importa al otro, `go.work` hace que se usen las copias locales en vez de descargar versiones
publicadas (como un `"repositories": [{"type": "path"}]` de Composer).

## Bonus: la carpeta `internal/`

`internal` tiene un significado especial para el compilador: un paquete dentro de `internal/` **solo** lo
pueden importar paquetes que estén bajo la carpeta padre de `internal/` (aquí, `backend/`). Aunque alguien
de fuera conozca la ruta `github.com/bayalad7/cargo-3d/backend/internal/engine`, el compilador le negará el
import. Es como si en PHP pudieras marcar un namespace completo como `private`.

Por eso el motor, la API y el acceso a la base de datos viven en `backend/internal/`: son piezas del
servidor, no una librería pública.

## Cuándo sí se descargará algo

Hoy todos los imports son de la librería estándar o de nuestro propio módulo, así que `go.mod` no tiene
sección `require` y no existe `go.sum`. La primera dependencia externa será `pgx` en el sprint 2.3
(ADR-0007). Ese día:

1. `go get github.com/jackc/pgx/v5` la descarga a `~/go/pkg/mod`.
2. `go.mod` gana una línea `require github.com/jackc/pgx/v5 vX.Y.Z`.
3. Aparece `go.sum` con las huellas, para que nadie pueda cambiar el código descargado sin que se note.

## Ejercicio

1. Abre `backend/go.mod` y cambia temporalmente el nombre del módulo a `cargo3d/backend`. Corre
   `go build ./backend/...`: falla, porque los imports ya no empiezan con el nombre del módulo y Go los
   busca como dependencias externas.
2. Cambia también los imports de `main.go` y `health.go` a `cargo3d/backend/internal/...` y vuelve a
   compilar: funciona. Así compruebas que es solo un nombre.
3. Regresa todo como estaba (`git checkout backend`) y confirma con `go test ./backend/...`.
