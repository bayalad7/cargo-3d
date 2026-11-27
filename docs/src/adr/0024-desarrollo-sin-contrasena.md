# ADR-0024: PostgreSQL sin contraseña en desarrollo

- **Estado:** Aceptado
- **Fecha:** 2026-11-26

## Contexto

Al instalar el proyecto en otra computadora, `levantar.sh` copió `.env.example` a `.env` y la base de datos quedó con
la contraseña de ejemplo (`cambia-esta-contrasena`). Esa contraseña no protege nada, porque está publicada en el repo,
y además hay que escribirla en pgAdmin y en cada herramienta.

La imagen `postgres:19` no arranca con `POSTGRES_PASSWORD` vacía, salvo que `POSTGRES_HOST_AUTH_METHOD=trust`
diga que no pida contraseña. Las dos variables solo se aplican cuando se crea el volumen.

## Decisión

En desarrollo PostgreSQL no lleva contraseña: `.env.example` trae `POSTGRES_PASSWORD=` vacía y
`POSTGRES_HOST_AUTH_METHOD=trust`. Si `docker-compose.yml` no recibe la segunda variable, usa `scram-sha-256` y
se exige contraseña. Producción siempre tiene contraseña y nunca usa `trust`.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| Contraseña de ejemplo publicada | Igual a producción | No protege nada y hay que escribirla en cada herramienta |
| Contraseña distinta en cada computadora | Algo más de protección | Hay que recordarla y escribirla, aunque solo se entra desde la propia computadora |
| **`trust` sin contraseña (elegida)** | Cero fricción; pgAdmin y `psql` conectan directo | Cualquier programa de la computadora puede entrar como cualquier usuario |

## Consecuencias

- **Positivas:** instalar en una computadora nueva no pide pasos con la contraseña; `DATABASE_URL` queda como
  `postgres://cargo3d:@localhost:5433/cargo3d?sslmode=disable`.
- **Negativas / costos:** el riesgo queda limitado porque el puerto solo escucha en `127.0.0.1`. Un `.env` de
  desarrollo nunca debe llegar a un servidor.
- **Qué tendríamos que hacer para revertirla:** poner una contraseña en `POSTGRES_PASSWORD`, quitar
  `POSTGRES_HOST_AUTH_METHOD` de `.env` y recrear el volumen (`docker compose --profile api down -v`).
