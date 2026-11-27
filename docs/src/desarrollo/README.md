# Desarrollo

Todo lo necesario para **programar** Cargo 3D: preparar la computadora, levantar los servicios, ver la base de
datos y trabajar con git y Claude Code. Lo que explica **el producto** (conceptos, contrato, casos) vive en
otras secciones; la organización completa está en el [ADR-0021](../adr/0021-organizacion-de-la-documentacion.md).

## Por dónde empezar

| Si quieres… | Lee |
|---|---|
| Preparar una computadora nueva (una sola vez) | [Instalación del entorno](entorno/instalacion.md) |
| Encender o apagar PostgreSQL y la API | [Levantar y detener servicios](entorno/servicios.md) |
| Ver la base de datos con pgAdmin o `psql` | [Base de datos](entorno/base-de-datos.md) |
| Saber cómo trabajar con Claude Code y los comandos `/sesion`, `/cerrar-sesion`… | [Flujo con Claude Code](flujo/claude-code.md) |
| Hacer commits, push o liberar una versión | [Ramas, commits y versiones](flujo/ramas-y-commits.md) |
| Nombrar algo, elegir unidades o escribir un error | [Convenciones](convenciones.md) |
| Entender el Go del proyecto | [Aprendizaje → Go](../aprendizaje/go/00-introduccion.md) |

## Qué va en esta carpeta

```text
desarrollo/
  README.md           esta página
  entorno/            preparar y operar la computadora: instalación, servicios, base de datos
  flujo/              cómo se trabaja: Claude Code, git, versiones
  convenciones.md     reglas de nombres, unidades, errores y comentarios
```

Una página nueva sobre **cómo desarrollar** (una herramienta, un procedimiento, un comando del día a día) va aquí,
en `entorno/` o `flujo/`, y se agrega a la tabla de arriba y a `SUMMARY.md`.
