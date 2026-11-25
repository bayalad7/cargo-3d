# Pruebas con cargas reales

Procedimiento para el **periodo de pruebas** que empieza con la versión v0.2.0 (MVP, fin de la fase 2): cada empresa
que use Cargo 3D registra sus cargas reales y lo que pasó en el patio, y eso confirma o corrige los
[casos simulados](README.md).

## Principios

1. **Lo real se queda privado.** Los registros con nombres, placas, pedidos o fotos de una empresa viven en
   `docs/datos/pruebas/<empresa>/` (no se sube a git). Al repositorio solo llegan casos anonimizados.
2. **Cada carga real se compara con un plan.** Sin plan de Cargo 3D no hay prueba: la idea es medir el motor.
3. **Toda diferencia termina en algo:** un caso corregido, un caso nuevo, una restricción nueva, una medida de equipo
   corregida o un pendiente (`P-NNN`).

## Antes de empezar con una empresa

| Paso | Qué se hace | Dónde queda |
|---|---|---|
| 1 | Ficha de la empresa: tipo de operación (marítima, carretera o ambas), rutas y tipos de camino | `docs/datos/pruebas/<empresa>/empresa.md` |
| 2 | Alta de sus equipos con medidas propias: cada medida dice si es **medida**, de **tarjeta de circulación** o de **ficha** (ADR-0019) | En la app; copia en `equipos.csv` de la misma carpeta |
| 3 | Reglas del patio: qué nunca va abajo, qué nunca se encima, holguras, orden por destino | `reglas.md` de la misma carpeta |
| 4 | Versión de la app y del motor con que se prueba | `empresa.md` |

## Por cada carga real

Una fila en `docs/datos/pruebas/<empresa>/registro.csv`:

| Columna | Qué se anota |
|---|---|
| `fecha`, `id_plan` | Día y plan de Cargo 3D (versión del motor incluida en el plan) |
| `equipo` | Id del equipo de la empresa usado |
| `piezas_plan`, `piezas_cargadas` | Cuántas propuso el plan y cuántas se cargaron de verdad |
| `cupo_como_plan` | sí / no / con cambios |
| `cambios_en_patio` | Qué se movió, giró o se dejó y por qué (texto corto) |
| `medidas_distintas` | Piezas cuya medida o peso real no coincidió con la declarada |
| `tiempo_planeacion_min`, `tiempo_carga_min` | Para comparar contra cómo planeaban antes |
| `peso_bascula_kg` | Si se pesó: bruto y por eje |
| `fotos` | Nombres de archivo en `fotos/` (antes, durante y al cerrar) |
| `caso` | `CS-NNN` relacionado, si lo hay |
| `resultado` | `confirma`, `corrige`, `caso-nuevo` o `restriccion-nueva` |

## Qué hacer con cada resultado

| Resultado | Acción |
|---|---|
| `confirma` | El caso pasa a `Validado en real` y su historial cita la fecha y la empresa por un código (`EMP-01`), nunca por nombre |
| `corrige` | Se ajusta el caso (`Corregido`) y su fixture; si cambió una medida del catálogo, se corrige con su fuente |
| `caso-nuevo` | Se escribe un `CS-NNN` nuevo, anonimizado, en la subcategoría que le toca |
| `restriccion-nueva` | Se abre un pendiente para el motor (regla del patio que el motor no conoce) |

## Revisión semanal

Cada semana del periodo de pruebas, con `/iniciar-sesion` se revisa `registro.csv` de cada empresa y se deja en la
bitácora (tema «Pruebas con cargas reales») cuántas cargas hubo, cuántas cupieron como el plan y qué casos cambiaron.
Esa cifra es la que mide si el producto está listo para la fase 3.
