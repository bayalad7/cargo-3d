# Casos de uso

Situaciones de carga **simuladas** a partir de lo que pasa en la vida real: qué equipo se usa, qué piezas llegan, qué
reglas aplican y qué resultado se espera del motor. Sustituyen a las cargas reales del patio mientras se construye la
app (ADR-0020). Cuando haya una versión que use una empresa (v0.2.0), cada caso se contrasta con cargas reales
siguiendo [Pruebas con cargas reales](pruebas-reales.md).

## Cómo están organizados

```
casos/
  maritimo/      seco · refrigerado · especial (open top, flat rack, plataforma ISO)
  carretera/     caja-seca · refrigerada · plataforma · rigidos (camionetas, rabón, tortón) · sobredimensionada
  restricciones/ apilado · multidestino · peso-por-eje · orientacion
  limites/       peso · vgm
  operacion/     equipo-alterado · cambios-en-patio
```

Cada caso es un archivo `cs-NNN-<tema>.md`. El número `CS-NNN` no se reutiliza; si un caso se descarta, se queda
con estado `Descartado` y la razón en su historial. Para una subcategoría nueva basta con crear la carpeta y agregar el
caso al índice de su categoría y a `SUMMARY.md`.

## Qué lleva un caso

| Sección | Contenido |
|---|---|
| Tabla inicial | Estado, categoría, equipo, fase en que aplica, motivos de rechazo esperados y ejemplo del contrato relacionado |
| Situación | La historia en dos o tres líneas, como la contaría quien carga |
| Equipo | Las medidas que usa el caso y de dónde salen: id de `contracts/equipment.json` o equipo de la empresa simulado |
| Piezas | Una fila por tipo de pieza con los campos del contrato (`allowedOrientations`, `stackable`, `deliverySequence`) |
| Resultado esperado | Cuántas piezas se colocan, cuáles no y con qué motivo |
| Por qué | La cuenta que lo demuestra |
| Qué validar con datos reales | Lo que hay que comprobar en el periodo de pruebas |
| Historial | Fecha, estado y nota de cada cambio |

## Estados

| Estado | Significado |
|---|---|
| `Simulado` | Escrito con fichas públicas y datos simulados; todavía no se comparó con una carga real |
| `Validado en real` | Una carga real parecida dio el resultado esperado (se cita el registro de pruebas, sin datos del cliente) |
| `Corregido` | La carga real mostró otra cosa; el caso se ajustó y el historial dice qué cambió y por qué |
| `Descartado` | El caso no pasa en la práctica; se deja con la razón |

## Datos simulados

- Lo que viene del catálogo usa el id de `contracts/equipment.json` y sus fuentes.
- Lo simulado se marca siempre con **simulado**: medidas de camionetas, rabones y tortones, taras de tractos y la útil
  de equipos sin ficha. Son valores típicos del mercado mexicano, no de una unidad concreta.
- Nada de nombres de empresas, clientes, placas ni direcciones.

## De los casos a las pruebas del motor

En la sesión 1.1.5 cada caso se convierte en `engine/fixtures/casos/cs-NNN/input.json` (y una variante por cada
equipo o camino que el caso compara) con su `.expect.json`. Los casos de la fase 2 entran como fixtures cuando llegue
su sesión; mientras, sirven de especificación.

## Índice

- [Marítimo](maritimo/README.md): CS-001, 002, 004 a 010.
- [Carretera](carretera/README.md): CS-011 a 020, 028 y 029.
- [Restricciones](restricciones/README.md): CS-021 a 024 (fase 2).
- [Límites](limites/README.md): CS-003 y 025.
- [Operación](operacion/README.md): CS-026 y 027.
