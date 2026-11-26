# Levantar y detener servicios

Los servicios del proyecto corren en Docker, en un solo grupo llamado `cargo-3d`:

| Contenedor | Imagen | Puerto en tu computadora | Cuándo corre |
|---|---|---|---|
| `cargo-3d-db-1` | `postgres:19` | `127.0.0.1:5433` | Siempre que trabajes con datos |
| `cargo-3d-api-1` | `cargo-3d-api` (se construye del `backend/Dockerfile`) | `127.0.0.1:8080` | Solo con `--profile api` |

Requisito: haber copiado `.env.example` a `.env` ([Instalación del entorno](instalacion.md#3-docker)).
Para ver los datos de PostgreSQL con pgAdmin o `psql`: [Base de datos](base-de-datos.md).

## Día a día: el script

Mientras programas el backend, PostgreSQL corre en Docker y la API con `go run` (recompila en segundos y se
depura desde el editor). **Un script lo hace todo:** abre Docker si está cerrado, crea `.env` si falta, apaga
el contenedor de la API, levanta PostgreSQL, espera a que esté listo y corre la API; al detenerla con `Ctrl+C`,
apaga PostgreSQL.

| Cómo | Windows | macOS / Linux |
| --- | --- | --- |
| VS Code | `Ctrl+Shift+B` (tarea "Levantar servicios") | `Cmd+Shift+B` / `Ctrl+Shift+B` |
| Doble clic | `scripts\servicios\levantar.cmd` | `scripts/servicios/levantar.command` (Finder) |
| Terminal | `.\scripts\servicios\levantar.ps1` | `./scripts/servicios/levantar.sh` |
| Dejar PostgreSQL encendido al salir | `levantar.ps1 -DejarDb` | `levantar.sh --dejar-db` |
| Apagar todo | `detener.ps1` / `detener.cmd` o la tarea "Detener servicios" | `detener.sh` / `detener.command` |

En macOS, la primera vez que hagas doble clic en un `.command`, macOS puede pedir permiso:
clic derecho → **Abrir**.

> **Para usar pgAdmin** mientras programas, levanta con `--dejar-db` (o `docker compose up -d db`):
> así PostgreSQL sigue encendido aunque detengas la API.

Lo que hace el script por dentro:

- `docker compose up -d db` levanta solo PostgreSQL (el mismo contenedor `cargo-3d-db-1`, no uno nuevo)
  y la API se corre con `go run ./backend/cmd/api`.
- No mezcles las dos formas al mismo tiempo: `go run` y el contenedor de la API usan el puerto 8080.
  Antes de `go run`, apaga el grupo con `docker compose --profile api down`.
- Las dos formas usan el mismo volumen `cargo-3d_db-data`: la base de datos es la misma.

## Todo dentro de Docker

Para probar la imagen de producción de la API, siempre desde la raíz del repo:

```bash
docker compose --profile api up -d --build   # levanta PostgreSQL y la API
docker compose --profile api ps              # db "healthy", api "running"
docker compose --profile api down            # apaga todo; los datos se conservan en el volumen
```

En Docker Desktop aparece **un solo grupo** con dos contenedores:

```
▾ cargo-3d
    ├─ cargo-3d-db-1    postgres:19    (127.0.0.1:5433)
    └─ cargo-3d-api-1   cargo-3d-api   (127.0.0.1:8080)
```

Para liberar la imagen de la API y la caché de construcción:
`docker compose --profile api down --rmi local` y `docker builder prune -f`.

## Reglas para no llenar el disco

- **Usa siempre `docker compose` desde la raíz del repo**, nunca `docker run` ni `docker build` sueltos:
  esos crean contenedores e imágenes fuera del grupo `cargo-3d`.
- El nombre del grupo está fijo en `docker-compose.yml` (`name: cargo-3d`), así que es el mismo aunque
  la carpeta se llame distinto o estés en otra computadora; volver a correr el comando reutiliza lo que ya existe.
- Incluye siempre `--profile api` al levantar **y** al apagar; sin él, el contenedor de la API queda fuera del comando.
- Los contenedores no arrancan solos al abrir Docker: solo corren cuando los levantas.
- **Puerto 5433:** PostgreSQL de Docker usa el 5433 de tu computadora para no chocar con un PostgreSQL
  instalado en el sistema (que ocupa el 5432). Si el 5433 también está ocupado, cambia `POSTGRES_PORT` en `.env`.

## Pruebas en un proyecto aparte

Para probar algo sin tocar la base de datos de desarrollo se usa el proyecto `cargo3d-test`, y **todo lo que
se crea se borra al terminar**, pase o falle (reglas completas en «Reglas de Docker» de `CLAUDE.md`):

```bash
docker system df && docker image ls                    # 1. anotar lo que ya existe
docker compose -p cargo3d-test up -d db                 # 2. probar (o: --profile api up -d --build)
docker compose -p cargo3d-test --profile api down -v --rmi all --remove-orphans   # 3. limpiar
docker builder prune -f
docker system df                                        # 4. confirmar que quedó como al inicio
```

Nunca uses `docker system prune -a --volumes`: borraría también el volumen con tu base de datos de desarrollo.
