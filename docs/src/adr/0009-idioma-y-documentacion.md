# ADR-0009: Idioma del código y política de documentación

- **Estado:** Aceptado
- **Fecha:** 2026-11-20

## Contexto

La mayor parte del código la escribe Claude; Benjamín tiene que entenderlo, mantenerlo y aprender Go en
el proceso. Claude no recuerda perfectamente entre sesiones: el repositorio tiene que ser la memoria del proyecto.

## Decisión

1. **Código en inglés, comentarios y documentación en español** (detalle en `desarrollo/convenciones.md`).
2. **Comentarios de Go muy detallados al inicio** (casi cada línea), reduciéndose por etapas.
3. **Cinco capas de documentación:** código (comentarios de `go doc` y TSDoc), ADRs, conceptos,
   guía de aprendizaje de Go comparada con PHP, y bitácora/changelog, publicadas con **mdBook**.
4. **`CLAUDE.md`** en la raíz como resumen vivo del proyecto para cada sesión.
5. **Nada se da por terminado sin documentación, `gofmt`, `go vet` y pruebas.**
6. ~~**Puntos de control:** al terminar un módulo, Benjamín lo explica con sus palabras.~~ Sustituido por el
   [ADR-0018](0018-sin-puntos-de-control.md) (2026-11-25).

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Inglés + español (elegida)** | Estándar de industria en código; documentación clara | Mezcla de idiomas |
| Todo en español | Máxima claridad | Choca con librerías y con `go doc` |
| Todo en inglés | Uniforme | Fricción de comprensión |

## Consecuencias

- **Positivas:** cualquier sesión retoma el contexto leyendo `CLAUDE.md`; Benjamín puede mantener el sistema solo.
- **Negativas:** documentar toma tiempo extra en cada entrega (es intencional).
