# 20 nov 2026 · Prompts por sesión y skills de Claude Code

**Sesión:** Fase 0 · **Área:** docs

## Prompts y skills de Claude Code

### Qué se hizo

- Biblioteca de prompts por sesión (`docs/src/prompts.md`): cada sprint de las fases 1 a 3 dividido en
  sesiones numeradas (`fase.sprint.sesión`), cada una con modelo sugerido, rama, requisitos, prompt y
  criterio de terminado. Incluye una tabla de prompts de apoyo para cualquier momento.
- Nueve skills de Claude Code en `.claude/skills/`: `iniciar-sesion`, `sesion`, `cerrar-sesion`,
  `punto-de-control`, `revisar`, `explicar`, `nuevo-adr`, `cambiar-contrato` y `cerrar-sprint`.
- Guía del flujo de trabajo (`docs/src/flujo-claude-code.md`) con el ciclo completo de una sesión.
- `CLAUDE.md`, el plan, el índice, el CHANGELOG y `.gitignore` actualizados.

### Decisiones

- Se usan **skills** (`.claude/skills/<nombre>/SKILL.md`), el formato recomendado por la documentación
  oficial de Claude Code; los archivos en `.claude/commands/` siguen funcionando, pero quedaron como
  formato heredado.
- Todas las skills tienen `disable-model-invocation: true`: solo Benjamín las ejecuta.
- Ninguna skill hace commit ni push.
- `/sesion` siempre propone un plan y espera aprobación antes de escribir código.

### Pendientes

1. Probar que las 9 skills aparezcan al escribir `/` en Claude Code desde la raíz del repo.
2. Tener los datos reales en `docs/datos/` y arrancar con `/iniciar-sesion` y `/sesion 1.1.1`.

## Biblioteca de prompts completada

- Sesiones nuevas: publicar v0.1.0 (1.7.4), v0.2.0 (2.6.6) y v1.0.0 (3.6.5); cierre con punto de control de cada
  bloque de la fase 3 (3.1.4, 3.2.5, 3.3.4, 3.4.4, 3.5.3 y dentro de 3.6.4).
- Puntos de control agregados en los sprints 1.6, 1.7, 2.4, 2.5 y 2.6 (la skill `/cerrar-sprint` los exige).
- Prompts de apoyo nuevos: arreglo urgente de una versión, instalar o subir versiones (ADR-0011) y pruebas en Docker.
- `docs/datos/` queda fuera de git; los fixtures reales se suben anonimizados (regla nueva en `CLAUDE.md`).
- CHANGELOG sin entradas duplicadas; esta bitácora absorbe la de "prompts y skills" (un archivo por día).
