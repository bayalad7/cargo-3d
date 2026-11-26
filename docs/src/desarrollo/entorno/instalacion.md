# Instalación del entorno


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
Se hace **una sola vez** por computadora, antes de programar. Después, el día a día está en
[Levantar y detener servicios](servicios.md) y la conexión a PostgreSQL en [Base de datos](base-de-datos.md).


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
El proyecto se desarrolla igual en **Windows, macOS y Linux** (también en Mac con Apple Silicon).
Elige la columna de tu sistema; los comandos del proyecto (`go`, `pnpm`, `docker compose`) son los mismos en los tres.


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
| Herramienta | Windows (nativo) | macOS | Linux |
| --- | --- | --- | --- |
| Gestor de paquetes | `winget` (viene con Windows) | [Homebrew](https://brew.sh) | el de tu distribución |
| Git | `winget install Git.Git` (incluye Git Bash) | `xcode-select --install` o `brew install git` | `apt install git` |
| Go | `winget install GoLang.Go` | `brew install go` | binario de go.dev/dl |
| Docker | Docker Desktop (motor WSL2) | Docker Desktop u OrbStack (más ligero) | Docker Engine + plugin Compose |
| Node.js | `winget install Schniz.fnm`, luego `fnm install 26` | `brew install fnm`, luego `fnm install 26` | fnm o nvm |
| mdBook | `winget install Rustlang.mdBook` | `brew install mdbook` | binario de sus releases en GitHub |


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
Notas:


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
- **Windows:** se trabaja nativo con PowerShell o Git Bash; WSL2 es opcional (Docker Desktop lo usa por dentro).
  Si la ruta de tu usuario tiene espacios o acentos, clona el repo en una ruta corta como `C:\cargo-3d`.
- **Fin de línea:** `.gitattributes` y `.editorconfig` fijan LF para todos. Si clonaste antes de que existiera
  `.gitattributes`, corre una vez `git rm -r --cached . -q && git reset --hard` (con todo ya en commit).
- **Mac con Apple Silicon:** las imágenes de Docker del proyecto existen para ARM (`arm64`) y corren nativas.


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
### 1. Git y GitHub


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
1. Crea en GitHub el repositorio vacío `cargo-3d` (privado).
2. Entra a la carpeta `cargo-3d` y súbelo: `git remote add origin https://github.com/bayalad7/cargo-3d.git` y `git push -u origin main`.
3. Crea la rama `develop` desde `main` (`main` sigue como rama por defecto). El CI y la publicación
   ya están en `.github/workflows/` (ADR-0012 y ADR-0013).
4. Opcional: en Settings → Rules marca los checks del CI como obligatorios para que un CI en rojo bloquee el merge.
5. Si usas otro usuario u otro nombre de repo, cambia la ruta del módulo en `backend/go.mod` y en los imports (búscala con `grep -r bayalad7`).


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
### 2. Go


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
1. Instala Go 1.27 o superior desde go.dev/dl.
2. Verifica: `go version`.


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
### 3. Docker


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
Docker Desktop (Windows o macOS), OrbStack (macOS) o Docker Engine (Linux), con Docker Compose.


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
1. En la raíz del repo, copia `.env.example` a `.env` y cambia la contraseña
   (PowerShell: `Copy-Item .env.example .env`; macOS/Linux: `cp .env.example .env`).
2. Cómo usarlo sin llenar el disco: [Levantar y detener servicios](servicios.md).


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
### 4. Node.js y pnpm (desde el sprint 1.4: validador en TypeScript)


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
Node 26 LTS con fnm o nvm. Desde Node 25, Corepack ya no viene incluido: instala pnpm con
`npm install -g pnpm@12`. pnpm lee el campo `packageManager` de `package.json` y usa
exactamente la versión del proyecto.


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
### 5. mdBook (sitio de documentación)


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
`winget install Rustlang.mdBook` en Windows, `brew install mdbook` en macOS o el binario de sus releases en GitHub en Linux.
Después de instalarlo, **abre una terminal nueva** para que reconozca el comando `mdbook`.


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
### 6. Editor


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
| Extensión de VS Code | Para qué |
| --- | --- |
| Go (oficial, usa gopls) | Autocompletado, errores en vivo, pruebas y depuración. Indispensable. |
| Error Lens | Errores en la misma línea. |
| Docker | Administrar contenedores desde el editor. |
| ESLint y Prettier | TypeScript (desde la fase 1). |
| YAML | Editar `openapi.yaml` y `docker-compose.yml` con validación. |


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
### 7. Claude Code


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
Instálalo con la guía oficial en [docs.claude.com](https://docs.claude.com), abre una terminal en la raíz
del repo y ejecuta `claude`. Lee `CLAUDE.md` en automático.


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
### 8. Verificación


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
| Comando | Debe mostrar |
| --- | --- |
| `go test ./backend/...` | `ok` en `engine` y en `httpapi` |
| `go vet ./backend/...` | Nada (sin errores) |
| `docker compose --profile api up -d --build` y `docker compose --profile api ps` | `cargo-3d-db-1` como `healthy` y `cargo-3d-api-1` como `running` |
| Abrir `localhost:8080/health` | `{"status":"ok","engineVersion":"0.1.0"}` (la API desde Docker) |
| `docker compose --profile api down` | Los dos contenedores se detienen y se eliminan; el volumen se conserva |
| `mdbook serve docs --open` | El sitio de documentación |


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).


### 9. pgAdmin (opcional)

Para ver la base de datos con una interfaz gráfica. Instalación y conexión en [Base de datos](base-de-datos.md).
Para programar el backend a diario usa el script de [Levantar y detener servicios](servicios.md#día-a-día-el-script).
