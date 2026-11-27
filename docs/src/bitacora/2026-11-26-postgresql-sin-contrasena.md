# 26 nov 2026 · PostgreSQL sin contraseña en desarrollo

**Sesión:** Fuera del plan · **Área:** infra, docs

## Qué se hizo

- El commit `ec8ee2b` (ADR-0024) dejó la base de datos de desarrollo sin contraseña: `.env.example` trae
  `POSTGRES_PASSWORD=` vacía y `POSTGRES_HOST_AUTH_METHOD=trust`, y `docker-compose.yml` usa `scram-sha-256` si esa
  variable no llega. Actualizó los avisos de `levantar.sh` y `levantar.ps1`, la guía de pgAdmin y `CLAUDE.md`.
- La revisión de ese commit (`/revisar`, durante la sesión 1.1.5) encontró detalles y se corrigieron:
  - `README.md` todavía decía «cambia la contraseña».
  - Sin `.env`, `docker compose` ya no se detiene con un mensaje propio: el contenedor arranca y se detiene. Se explicó
    en el comentario de `docker-compose.yml`, en `servicios.md` y en el ADR.
  - El ADR y la guía decían que `trust` solo acepta a la propia computadora; acepta también a cualquier contenedor de la
    red de Compose. Se corrigieron los dos.

## Decisiones

- ADR-0024: en desarrollo, PostgreSQL sin contraseña (`trust`); en producción siempre con contraseña y sin `trust`.

## Verificación

- `docker compose config` (solo lee el archivo, no crea nada) dice que `docker-compose.yml` es válido.
- **No se probó** levantar PostgreSQL con `docker compose up -d db` ni conectar sin contraseña: ver P-024.

## Pendientes

- P-024: ahora incluye comprobar la conexión sin contraseña.
