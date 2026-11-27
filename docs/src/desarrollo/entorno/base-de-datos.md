# Base de datos: conectarse a PostgreSQL de Docker

PostgreSQL 19 corre en el contenedor `cargo-3d-db-1` ([Levantar y detener servicios](servicios.md)).
Desde tu computadora te conectas como si estuviera instalado localmente, en el puerto **5433**.
Esta página explica cómo verlo con **pgAdmin** (interfaz gráfica) o con **psql** (terminal).

> **Hoy la base de datos está vacía.** Las primeras tablas llegan en el sprint 2.3, con las migraciones de
> `backend/internal/store/migrations/` (ADR-0007). Mientras tanto sirve para comprobar que la conexión funciona.

## 1. Encender PostgreSQL

pgAdmin solo se conecta si el contenedor está corriendo. Cualquiera de estas opciones lo deja encendido:

```bash
docker compose up -d db                 # solo PostgreSQL
./scripts/servicios/levantar.sh --dejar-db   # PostgreSQL + API con go run; al salir, PostgreSQL sigue encendido
```

Para confirmar: `docker compose ps` debe mostrar `cargo-3d-db-1` como `healthy`.

## 2. Instalar pgAdmin

Instálalo como aplicación de escritorio, **no** como contenedor de Docker: así no se suma otra imagen al disco
y no hay que resolver la red entre contenedores.

| Sistema | Comando |
|---|---|
| macOS | `brew install --cask pgadmin4` (o el instalador de pgadmin.org) |
| Windows | `winget install PostgreSQL.pgAdmin` |
| Linux | el repositorio oficial de pgadmin.org (`pgadmin4-desktop`) |

Usa **la versión más reciente**: cada versión de pgAdmin reconoce hasta cierta versión de PostgreSQL. Si una
versión vieja no reconoce PostgreSQL 19, muestra un aviso o falla al leer el catálogo; la solución es actualizar
pgAdmin, no PostgreSQL.

## 3. Registrar el servidor en pgAdmin

En pgAdmin: clic derecho en **Servers → Register → Server…**

| Pestaña | Campo | Valor |
|---|---|---|
| General | Name | `cargo-3d (desarrollo)` (cualquier nombre; solo lo ves tú) |
| Connection | Host name/address | `127.0.0.1` |
| Connection | Port | `5433` (o el `POSTGRES_PORT` de tu `.env`) |
| Connection | Maintenance database | `cargo3d` (`POSTGRES_DB`) |
| Connection | Username | `cargo3d` (`POSTGRES_USER`) |
| Connection | Password | vacía (en desarrollo no hay contraseña; ver ADR-0024) |
| Connection | Save password? | No hace falta |
| Parameters | SSL mode | `disable` (la conexión no sale de tu computadora) |

Los valores salen de `.env` y `docker-compose.yml`. Si cambiaste alguno en `.env`, usa el tuyo.
En desarrollo PostgreSQL usa `trust`: acepta conexiones sin contraseña, y solo desde tu computadora porque el
puerto escucha en `127.0.0.1` ([ADR-0024](../../adr/0024-desarrollo-sin-contrasena.md)). En producción sí hay contraseña,
y **nunca** se escribe en la documentación ni en el código.

Una vez conectado, las tablas estarán en **cargo-3d (desarrollo) → Databases → cargo3d → Schemas → public → Tables**.

### La misma conexión en una línea

Útil para `psql`, DBeaver, TablePlus o la extensión de PostgreSQL de VS Code:

```text
postgres://cargo3d:@127.0.0.1:5433/cargo3d?sslmode=disable
```

La API usa esta misma cadena en `DATABASE_URL`; dentro de Docker el host es `db` y el puerto `5432`
(ver `docker-compose.yml`), porque ahí habla de contenedor a contenedor.

## 4. Sin pgAdmin: psql dentro del contenedor

No hace falta instalar nada; `psql` viene en la imagen de PostgreSQL:

```bash
docker compose exec db psql -U cargo3d cargo3d
```

Comandos útiles dentro de `psql`: `\dt` (tablas), `\d nombre_tabla` (columnas), `\l` (bases de datos),
`\q` (salir).

## Reglas al usar pgAdmin

- **Mira, no cambies el esquema a mano.** Las tablas se crean y modifican solo con migraciones
  (`backend/internal/store/migrations/`). Un cambio hecho en pgAdmin no queda en git, no existe en otra
  computadora ni en producción, y desordena las migraciones siguientes. Si algo debe cambiar, va en una migración.
- **Datos de prueba sí**, pero recuerda que viven en el volumen `cargo-3d_db-data`: se borran si alguien
  corre `docker compose down -v`.
- **Datos reales de clientes no.** Las cargas reales viven en `docs/datos/` y no se cargan a la base de
  desarrollo hasta las pruebas de la v0.2.0 (P-022).

## Conectarse a la base de pruebas (`cargo3d-test`)

Las pruebas usan un proyecto de Compose aparte ([Pruebas en un proyecto aparte](servicios.md#pruebas-en-un-proyecto-aparte)).
Como usaría el mismo puerto 5433, levántalo en otro:

```bash
POSTGRES_PORT=5434 docker compose -p cargo3d-test up -d db
```

En pgAdmin registra un segundo servidor, `cargo-3d (pruebas)`, con puerto `5434` y los mismos datos.
Al terminar la prueba se borra (`down -v`); el servidor registrado en pgAdmin puede quedarse, solo no conectará.

## Problemas comunes

| Síntoma | Causa | Solución |
|---|---|---|
| `Connection refused` | El contenedor está apagado | `docker compose up -d db` |
| `password authentication failed` | Contraseña distinta a la del `.env`, o el volumen se creó con otra contraseña | Revisa `.env`. La contraseña y `POSTGRES_HOST_AUTH_METHOD` se fijan al **crear** el volumen; cambiarlas en `.env` después no las cambia en PostgreSQL. Recrea el volumen: `docker compose --profile api down -v` |
| Conecta a otra base o pide otro usuario | Usaste el puerto 5432 (PostgreSQL instalado en el sistema) | Usa el 5433 |
| `port is already allocated` al levantar | Otro programa ocupa el 5433 | Cambia `POSTGRES_PORT` en `.env` y en pgAdmin |
| Aviso de versión no soportada | pgAdmin viejo | Actualiza pgAdmin |
| `Databases` vacío o sin tablas | Normal antes del sprint 2.3 | — |
