# ADR-0007: PostgreSQL con pgx y sqlc

- **Estado:** Aceptado (se implementa en la fase 2)
- **Fecha:** 2026-11-20

## Contexto

Hay que guardar datos relacionales (empresas, usuarios, catálogos) y documentos complejos
(un plan de carga con cientos de colocaciones). En el futuro puede haber rutas y ubicaciones.

## Decisión

- **PostgreSQL** como base de datos, en Docker para desarrollo (ADR-0010).
- **pgx** como controlador de PostgreSQL para Go.
- **sqlc**: se escribe SQL normal en archivos `.sql` y sqlc genera funciones de Go tipadas.
  No hay ORM.
- Migraciones con archivos SQL versionados (la herramienta se elige en un ADR al iniciar el sprint 2.3).

### Dónde vive el código

Todo lo relacionado con PostgreSQL vive en un solo paquete de Go, `backend/internal/store`.
Las carpetas se crean al empezar el sprint 2.3, no antes.

```
backend/
  sqlc.yaml                     configuración de sqlc (lee queries/ y migrations/, escribe sqlcgen/)
  internal/
    store/                      paquete de acceso a datos: el único que importa pgx
      migrations/               cambios al esquema, versionados: 0001_crear_empresas.up.sql / .down.sql
      queries/                  el SQL que se escribe a mano, un archivo por tema: plans.sql, companies.sql
      sqlcgen/                  código de Go generado por sqlc: no se edita a mano
      store.go                  conexión (pool de pgx) y funciones que usa httpapi
```

Reglas:

1. **Solo `internal/store` habla con la base de datos.** Ningún otro paquete importa pgx ni `sqlcgen`.
2. **El motor no sabe que existe PostgreSQL.** `internal/httpapi` llama a `store` para leer o guardar
   y a `engine` para calcular; `engine` nunca importa `store`.

   ```
   httpapi ──► store ──► PostgreSQL
      │
      └──────► engine   (sin base de datos)
   ```
3. **`sqlcgen/` es código generado**, no se edita a mano: si hay que cambiarlo, se cambia
   el `.sql` de `queries/` (o una migración) y se vuelve a correr `sqlc generate`. Sí se sube al repo,
   para que el proyecto compile sin tener sqlc instalado.
4. **Una migración ya aplicada no se edita**: los cambios al esquema van en una migración nueva.
5. La cadena de conexión llega por la variable de entorno `DATABASE_URL` (ya definida en `docker-compose.yml`).

**Por qué así:** el SQL queda junto al Go que lo usa y es fácil de encontrar. `internal/` impide que algo
fuera del backend dependa de la base de datos, y separar `sqlcgen/` deja claro qué se escribe a mano y
qué se genera.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **pgx + sqlc (elegida)** | SQL que Benjamín ya sabe escribir; código de Go generado y tipado; errores de consulta detectados al compilar | Paso de generación de código |
| GORM | ORM popular | Magia y sorpresas de rendimiento; esconde el SQL |
| `database/sql` a mano | Sin herramientas | Mucho código repetitivo y propenso a errores |
| MariaDB | Conocida | JSON menos potente; sin equivalente maduro a PostGIS |

## Consecuencias

- **Positivas:** el plan de carga completo se guarda como JSONB; las consultas son SQL legible.
- **Negativas:** hay que correr `sqlc generate` cuando cambian las consultas.
- **Revertir:** el SQL es estándar; cambiar de herramienta en Go no toca la base de datos.
