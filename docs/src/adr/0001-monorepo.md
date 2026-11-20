# ADR-0001: Monorepo con Go workspace, pnpm y Turborepo

- **Estado:** Aceptado
- **Fecha:** 2026-11-20

## Contexto

El producto tiene un backend en Go (API y motor) y varias interfaces en TypeScript (web, escritorio,
móvil) que comparten contratos, el visor 3D, el validador y el cliente de la API. Si cada pieza viviera
en su propio repositorio, un cambio en el contrato obligaría a coordinar versiones entre varios repos.

## Decisión

Un solo repositorio con:

- **Go workspace** (`go.work`) para los módulos de Go (por ahora solo `backend/`).
- **pnpm workspaces** para los paquetes de TypeScript (`apps/*`, `packages/*`).
- **Turborepo** para orquestar builds y pruebas de TypeScript, con caché.
- **`contracts/`** en la raíz: OpenAPI y el esquema del motor, compartidos por todos.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Monorepo (elegida)** | Un cambio de contrato se ve al instante en backend y pantallas; una sola CI | Requiere orquestar dos ecosistemas |
| Repos separados (backend y frontend) | Cada uno independiente | Contratos duplicados o publicados como paquete; sincronización manual |
| Nx | Muy completo | Más complejo de lo que necesitamos |

## Consecuencias

- **Positivas:** un solo commit puede cambiar contrato, backend y pantalla juntos.
- **Negativas:** hay dos herramientas de construcción (go y pnpm); se mitiga con scripts en `package.json`.
- **Revertir:** `backend/` es un módulo de Go autocontenido y puede salir a su propio repo.
