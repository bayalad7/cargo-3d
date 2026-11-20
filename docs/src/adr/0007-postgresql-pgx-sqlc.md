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
