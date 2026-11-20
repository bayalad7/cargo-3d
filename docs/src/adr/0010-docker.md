# ADR-0010: Docker para la base de datos, el backend y el despliegue

- **Estado:** Aceptado
- **Fecha:** 2026-11-20

## Contexto

Benjamín ya usa Docker en otros proyectos. Hace falta una forma repetible de levantar PostgreSQL en
desarrollo, de construir el backend igual en cualquier máquina y de desplegarlo en un servidor.

## Decisión

Docker se usa donde aporta y no donde estorba:

| Dónde | ¿Docker? | Cómo |
|---|---|---|
| PostgreSQL en desarrollo | **Sí** | `docker compose up -d db` (`docker-compose.yml`) |
| Backend mientras programas | **No** (opcional) | `go run ./backend/cmd/api`: más rápido y más fácil de depurar |
| Backend en producción | **Sí** | `backend/Dockerfile` multi-etapa: compila en `golang` y corre en `distroless/static` como usuario sin privilegios |
| Probar todo junto | **Sí** | `docker compose --profile api up --build` |
| Integración continua | **Sí** | La CI construye la imagen en cada push |
| Web en desarrollo (Vite) | **No** | Corre local; recarga en caliente más rápida |
| Web en producción | **No necesario** | Son archivos estáticos: hosting estático o servidos por un contenedor ligero si se prefiere un solo servidor |
| Móvil y escritorio | **No** | Se compilan con sus propias herramientas (Expo, Wails) |

Reglas: versiones fijas de imágenes; secretos solo en `.env` (nunca en el repo ni en la imagen);
puertos publicados solo en `127.0.0.1` en desarrollo; datos de PostgreSQL en un volumen.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Docker donde aporta (elegida)** | Entorno repetible; despliegue de una imagen pequeña | Aprender Compose y multi-etapa (Benjamín ya conoce lo básico) |
| Docker para todo, incluido el desarrollo de Go y de la web | Entorno idéntico al de producción | Recompilación y recarga más lentas; depuración más difícil |
| Sin Docker | Nada que aprender | PostgreSQL instalado a mano en cada máquina; despliegues frágiles |

## Consecuencias

- **Positivas:** con Go, la imagen de producción es muy pequeña (un binario estático, sin sistema operativo
  completo), arranca rápido y tiene poca superficie de ataque.
- **Un solo grupo:** el proyecto de Compose tiene nombre fijo (`name: cargo-3d`), así que `db` y `api` siempre
  aparecen juntos en un grupo y nunca se duplican contenedores ni volúmenes. Los contenedores no se reinician solos.
  Se mantienen separados (y no en un solo contenedor) para poder reiniciar, actualizar y respaldar cada uno
  por su cuenta; en producción PostgreSQL puede ser un servicio administrado.
- **Negativas:** la imagen distroless no tiene shell: para depurar en producción se usan logs, no `docker exec`.
- **Revertir:** el binario de Go corre sin Docker en cualquier Linux; Docker es comodidad, no dependencia.
