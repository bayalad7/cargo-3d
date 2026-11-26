# Biblioteca de prompts por sesión

Este documento tiene **el prompt de cada sesión de trabajo** de las fases 1 a 3, en el orden del
[Plan de trabajo](plan.md). Cada sprint se divide en sesiones cortas con una sola meta.

## Cómo usar este documento

Todo se trabaja directo en `develop` (ver "Ramas" en `CLAUDE.md`). Tienes dos formas de usar cada sesión:

1. **Con la skill `/sesion` (recomendado):** escribe en Claude Code `/sesion 1.1.2` y Claude busca aquí
   la sesión `1.1.2`, la lee completa y la ejecuta empezando por un plan. No tienes que copiar nada.
2. **Copiando el prompt:** copia el bloque "Prompt" de la sesión y pégalo en Claude Code.

En ambos casos el ciclo es el mismo (detallado en [Flujo de trabajo con Claude Code](desarrollo/flujo/claude-code.md)):

```
git switch develop && git pull
claude
/model <modelo de la sesión>
/iniciar-sesion
/sesion <número>          ← revisas y apruebas el plan antes de que escriba código
... trabajo ...
/revisar                  ← opcional, antes de cerrar
/cerrar-sesion
git add -A && git commit  (con el mensaje que te propone)  && git push
```

## Formato de cada sesión

Cada sesión tiene un número `fase.sprint.sesión` (por ejemplo `1.1.2`) y estos campos:

- **Modelo:** el modelo y el esfuerzo sugeridos (`/model` y `/effort`).
- **Necesitas:** lo que debe existir antes de empezar.
- **Prompt:** el texto que se le da a Claude Code.
- **Terminada cuando:** cómo sabes que la sesión cumplió su meta.

> **Regla para editar este documento:** el encabezado de cada sesión debe empezar con `### ` seguido
> del número y un espacio (`### 1.1.2 — ...`). La skill `/sesion` busca exactamente ese formato.
> Si agregas o cambias sesiones, respétalo.

---

# Fase 1 — Prueba de concepto

## Sprint 1.1 — Modelos y contrato del motor

### 1.1.1 — Contrato del motor

- **Modelo:** `opus` / `medium`
- **Necesitas:** al menos 3 cargas reales y el catálogo de equipos en `docs/datos/` (Excel, CSV o fotos).
  Esa carpeta **no se sube a git** (está en `.gitignore`): son datos del cliente.

**Prompt:**

```
Sprint 1.1, sesión 1: contrato del motor. Aplica la regla 10 de CLAUDE.md:
el contrato va antes que el código.

Diseña contracts/engine.schema.json (JSON Schema 2020-12) con:
1. La entrada del motor: el equipo (contenedor, caja seca o plataforma con
   medidas internas y peso máximo) y la lista de piezas (id, descripción,
   largo, ancho, alto, peso, cantidad, orientaciones permitidas y, para el
   futuro, campos opcionales de apilabilidad y destino).
2. La salida: el LoadPlan con las colocaciones (pieza, posición x/y/z,
   orientación), las métricas (ocupación de volumen y de peso) y la lista
   de piezas que no cupieron, con el motivo.

Medidas en milímetros y pesos en kilogramos, siempre enteros. Revisa los
datos reales en docs/datos/ para que el esquema los cubra sin forzarlos.

Antes de escribir el archivo, preséntame: el esquema propuesto, un ejemplo
de entrada real, un ejemplo de salida y la explicación de cada campo.
No escribas código Go en esta sesión.
```

Cuando apruebes el diseño: *"Escribe el esquema y los dos ejemplos en contracts/examples/,
actualiza contracts/README.md y, si tomamos una decisión importante (por ejemplo, el sistema de
coordenadas), escribe el ADR."*

- **Terminada cuando:** existe `contracts/engine.schema.json`, los ejemplos validan contra el esquema
  y puedes explicar cada campo.

### 1.1.2 — Container y CargoItem

- **Modelo:** `opus` / `medium`
- **Necesitas:** sesión 1.1.1 terminada.

**Prompt:**

```
Sprint 1.1, sesión 2: implementa los tipos Container y CargoItem en
backend/internal/engine, siguiendo exactamente contracts/engine.schema.json.

- Usa int64 para milímetros y kilogramos, con la unidad en el nombre del
  campo (LengthMM, WeightKG).
- Agrega etiquetas json que coincidan con el contrato.
- Agrega el método Volume() a los dos tipos.
- Al llegar a MaxPayloadKG y a las medidas, explícame sobre el código dos temas
  que no me quedaron claros: quién calcula ese límite y por qué el motor
  no conoce la NOM-012, y hacia dónde se redondea cada valor.
- Comenta casi cada línea, como pide CLAUDE.md, explicando cada concepto
  de Go nuevo (struct, tipos numéricos, etiquetas, métodos) comparado con PHP.
- Pruebas de tabla para Volume() y para convertir de JSON a Go y de regreso
  usando contracts/examples/ (todos los input.json y output.json). Ojo con
  los valores por omisión (allowedOrientations = las 6, stackable = true,
  overhang = 0): un campo que no viene debe regresar sin venir, así que no
  se llenan al leer el JSON sino al calcular.
- Escribe el capítulo 02 de la guía de Go (docs/src/aprendizaje/go/
  02-structs-y-tipos.md) con Container y CargoItem como ejemplos, y enlázalo
  en SUMMARY.md.

Primero propón el plan: archivos, tipos, pruebas y contenido del capítulo.
```

- **Terminada cuando:** `go test ./backend/...` pasa, el capítulo 02 existe y puedes explicar por qué
  se usa `int64` y no `float64`.

### 1.1.3 — Orientation, Placement y LoadPlan

- **Modelo:** `opus` / `medium`
- **Necesitas:** sesión 1.1.2 terminada.

**Prompt:**

```
Sprint 1.1, sesión 3: implementa Orientation, Placement y LoadPlan según
contracts/engine.schema.json.

- Orientation representa las 6 formas de girar una caja alineada a los
  ejes. Necesito una función que, dada una pieza y una orientación,
  devuelva sus medidas rotadas (largo, ancho, alto).
- CargoItem debe poder decir qué orientaciones tiene permitidas.
- Placement une una pieza, su posición (x, y, z en mm desde la esquina del
  fondo, abajo, a la izquierda mirando del fondo hacia la puerta; ver
  ADR-0015) y su orientación.
- LoadPlan junta las colocaciones, las métricas y las piezas que no
  cupieron.
- Pruebas de tabla para las 6 orientaciones.
- Escribe el capítulo 03 de la guía (métodos y receptores) y explícame
  cuándo usar receptor por valor y cuándo por puntero, con estos tipos.

El sistema de coordenadas ya está en docs/src/conceptos/sistema-coordenadas.md:
sigue ese documento y corrígelo si el código muestra algo distinto.

Primero propón el plan.
```

- **Terminada cuando:** las pruebas de las 6 orientaciones pasan y el sistema de coordenadas está
  documentado en `docs/src/conceptos/`.

### 1.1.4 — Validación de la entrada

- **Modelo:** `opus` / `medium`
- **Necesitas:** sesión 1.1.3 terminada.

**Prompt:**

```
Sprint 1.1, sesión 4: validación de la entrada del motor.

Escribe una función que revise la entrada antes de calcular y devuelva
errores claros en español, nunca un panic:
- medidas en cero o negativas, pesos negativos, cantidades menores a 1;
- lista de piezas vacía o con ids repetidos;
- versión del contrato distinta y tramos (sections) encimados o fuera de
  rango (una puerta más grande que el interior NO es error: ver ADR-0016);
- un equipo con puerta (door) que además trae margen de sobredimensión
  (overhang mayor que cero): el margen es solo para tramos abiertos (ADR-0017);
- una pieza sin ninguna orientación permitida.

Una pieza que no cabe en ninguna orientación NO es error de validación:
sale en unplaced con EXCEEDS_DIMENSIONS o EXCEEDS_DOOR (ADR-0016).

Usa errores con contexto (fmt.Errorf con %w) y errores comparables con
errors.Is, para que la API pueda saber qué tipo de error fue. Si hay varios
problemas, repórtalos todos juntos, no solo el primero.

Pruebas de tabla con un caso por cada error y uno válido; usa también
contracts/examples/invalid/ (schema y semantic): todos deben dar error.
Escribe el capítulo 04 de la guía (errores como valores) comparándolo con try/catch
de PHP.

Primero propón el plan y la lista de errores que vas a definir.
```

- **Terminada cuando:** cada error tiene su prueba, ninguna entrada inválida provoca `panic`
  y entiendes la diferencia entre `errors.Is` y comparar textos.

### 1.1.5 — Casos simulados y catálogo de equipos

- **Modelo:** `opus` / `medium`
- **Necesitas:** sesión 1.1.4 terminada.

**Prompt:**

```
Sprint 1.1, sesión 5: convierte los casos de uso de docs/src/casos/ en
fixtures del motor: engine/fixtures/casos/cs-NNN/input.json, más una
variante por cada equipo o camino que el caso compare (por ejemplo
input-camino-b.json), siguiendo contracts/engine.schema.json. El catálogo
de equipos ya existe en contracts/equipment.json: usa sus medidas por id y
no lo cambies sin fuente. Los equipos simulados de los casos (camionetas,
rabón, tortón) van tal cual, marcados como simulados (ADR-0020).

Para cada fixture agrega un .expect.json: por ahora, que la entrada es
válida; además anota el resultado esperado del caso (piezas colocadas y
motivos) para usarlo desde el sprint 1.3. Los casos de la fase 2 solo se
convierten si el contrato 0.1.0 ya tiene sus campos. Escribe una prueba en
Go que recorra todos los fixtures, los lea, los valide contra el contrato y
pase la validación de entrada.

Si un caso no cabe en el contrato o tiene un dato dudoso, detente y
pregúntame antes de inventar o cambiar nada. Documenta en
engine/fixtures/README.md qué caso produce cada fixture.

Primero propón el plan.
```

- **Terminada cuando:** todos los casos de la fase 1 son fixtures válidos y la prueba los recorre.

### 1.1.6 — Cierre del sprint 1.1

- **Modelo:** `opus` / `high`
- **Necesitas:** sesiones 1.1.1 a 1.1.5 terminadas en `develop`.

**Prompt:**

```
/cerrar-sprint 1.1
```

- **Terminada cuando:** el checklist del sprint 1.1 en `plan.md` está completo.

## Sprint 1.2 — Geometría

### 1.2.1 — Cajas en el espacio

- **Modelo:** `opus` / `high`
- **Necesitas:** sprint 1.1 cerrado.

**Prompt:**

```
Sprint 1.2, sesión 1: geometría de cajas alineadas a los ejes.

Implementa en backend/internal/engine un tipo Box (posición mínima y
medidas) con estas funciones:
- Intersects: dos cajas se enciman (tocarse en una cara NO es encimarse);
- ContainedIn: una caja queda completamente dentro de otra;
- Contact: el área de contacto entre la cara de abajo de una caja y la de
  arriba de otra (lo usaremos para estabilidad en el sprint 2.1).

Todo con enteros. Pruebas de tabla con casos límite: cajas que se tocan
por una cara, por una arista, por un vértice, una dentro de otra, en el
borde exacto del contenedor y separadas por 1 mm.

Escribe en docs/src/conceptos/geometria.md la explicación con dibujos en
texto. Primero propón el plan y la lista de casos límite.
```

- **Terminada cuando:** todas las pruebas de casos límite pasan y la explicación tiene dibujos.

### 1.2.2 — Orientaciones permitidas en la práctica

- **Modelo:** `opus` / `medium`
- **Necesitas:** sesión 1.2.1 terminada.

**Prompt:**

```
Sprint 1.2, sesión 2: une Orientation con Box. Para una pieza y un punto,
genera las cajas candidatas en cada orientación permitida, descartando las
que se salen del equipo. Agrega casos de la vida real: "este lado arriba" (solo
giros sobre el eje vertical), piezas cúbicas (orientaciones repetidas que
hay que eliminar) y piezas largas.

Escribe el capítulo 05 de la guía (slices, maps y ciclos) usando la lista
de candidatas como ejemplo. Primero propón el plan.
```

- **Terminada cuando:** las piezas cúbicas no generan candidatas repetidas y el capítulo 05 existe.

### 1.2.3 — Fuzzing de la geometría

- **Modelo:** `opus` / `high`
- **Necesitas:** sesión 1.2.2 terminada.

**Prompt:**

```
Sprint 1.2, sesión 3: agrega pruebas de fuzzing nativas de Go
(go test -fuzz) para la geometría. Propiedades que deben cumplirse siempre:
- Intersects es simétrica: A con B da lo mismo que B con A;
- una caja siempre se interseca consigo misma si tiene volumen;
- si A está contenida en B y B en C, A está contenida en C;
- girar una pieza nunca cambia su volumen.

Explícame qué es el fuzzing, cómo se corre y qué hacer cuando encuentra un
caso que falla (se guarda en testdata/). Agrega el comando a CLAUDE.md.
Primero propón el plan.
```

- **Terminada cuando:** el fuzzing corre al menos 1 minuto por propiedad sin fallas y sabes correrlo.

### 1.2.4 — Punteros y cierre

- **Modelo:** `opus` / `medium`
- **Necesitas:** sesiones 1.2.1 a 1.2.3 terminadas.

**Prompt:**

```
Escribe el capítulo 06 de la guía de Go (punteros: cuándo *Container y
cuándo Container) con ejemplos del código actual del motor, y revisa si
algún lugar del motor usa punteros sin necesidad o copia estructuras
grandes donde convendría un puntero. No cambies código sin preguntarme.
```

Después: `/cerrar-sprint 1.2`.

- **Terminada cuando:** el capítulo 06 existe y el sprint 1.2 está cerrado.

## Sprint 1.3 — Primer algoritmo: puntos extremos

### 1.3.1 — Diseño del algoritmo

- **Modelo:** `opusplan` / `high`
- **Necesitas:** sprint 1.2 cerrado.

**Prompt:**

```
Sprint 1.3, sesión 1: diseño del algoritmo de puntos extremos, sin código.

Escribe docs/src/conceptos/puntos-extremos.md explicando paso a paso:
1. cómo se generan los puntos extremos al colocar una pieza;
2. cómo se elige el mejor punto y la mejor orientación (propón 2 o 3
   criterios de "mejor" y recomienda uno);
3. cómo se ordenan las piezas antes de colocarlas;
4. qué pasa cuando una pieza no cabe;
5. la complejidad aproximada y qué tan rápido debería ser con 100, 300
   y 1000 piezas.

Incluye un ejemplo a mano con 3 piezas en un contenedor pequeño, dibujado
en texto. Si hay decisiones de diseño, escribe el ADR.
```

- **Terminada cuando:** puedes explicar el algoritmo en el pizarrón con el ejemplo de 3 piezas.

### 1.3.2 — Ordenamiento de piezas

- **Modelo:** `sonnet` / `medium`
- **Necesitas:** sesión 1.3.1 terminada.

**Prompt:**

```
Sprint 1.3, sesión 2: implementa las estrategias de ordenamiento de piezas
descritas en docs/src/conceptos/puntos-extremos.md (volumen, área de base,
peso) como funciones intercambiables. Expande las cantidades (una pieza con
cantidad 5 se convierte en 5 unidades). El orden debe ser determinista:
con la misma entrada, siempre el mismo resultado (desempata por id).
Pruebas de tabla. Primero propón el plan.
```

- **Terminada cuando:** el ordenamiento es determinista y está probado.

### 1.3.3 — Colocación por puntos extremos

- **Modelo:** `opus` / `high`
- **Necesitas:** sesión 1.3.2 terminada.

**Prompt:**

```
Sprint 1.3, sesión 3: implementa la colocación por puntos extremos según
docs/src/conceptos/puntos-extremos.md. Función principal:
Optimize(input) (LoadPlan, error).

- Usa la geometría del sprint 1.2 para revisar que no haya traslapes.
- Cada punto extremo y cada orientación permitida se evalúan con el
  criterio elegido en el diseño.
- Las piezas que no caben van a la lista de no colocadas con su motivo.
- Pruebas: el ejemplo a mano de 3 piezas debe dar exactamente el
  resultado del documento; además, todos los fixtures de casos deben producir
  planes sin traslapes y dentro del equipo.

Comenta el algoritmo paso a paso enlazando a la sección del documento de
conceptos. Primero propón el plan.
```

- **Terminada cuando:** los fixtures de casos producen planes sin traslapes y con el resultado que dice cada caso.

### 1.3.4 — Peso máximo y métricas

- **Modelo:** `sonnet` / `medium`
- **Necesitas:** sesión 1.3.3 terminada.

**Prompt:**

```
Sprint 1.3, sesión 4: agrega la restricción de peso máximo del equipo
(una pieza que excede el peso restante no se coloca y queda con motivo
"peso") y calcula las métricas del LoadPlan: ocupación de volumen (%),
ocupación de peso (%), número de piezas colocadas y no colocadas.
Pruebas de tabla, incluido un caso donde el peso se acaba antes que el
espacio. Primero propón el plan.
```

- **Terminada cuando:** las métricas cuadran con un cálculo a mano en al menos un caso (por ejemplo CS-001).

### 1.3.5 — Cierre

- **Modelo:** `opus` / `high`

**Prompt:** `/cerrar-sprint 1.3`.

- **Terminada cuando:** el sprint 1.3 está cerrado.

## Sprint 1.4 — Validador independiente y banco de pruebas

### 1.4.1 — Paquete del validador en TypeScript

- **Modelo:** `sonnet` / `medium`
- **Necesitas:** sprint 1.3 cerrado; Node 26 LTS y pnpm 12 instalados (ADR-0011).

**Prompt:**

```
Sprint 1.4, sesión 1: crea el paquete packages/validator en TypeScript
estricto, con Vitest para pruebas, dentro del workspace de pnpm y de
Turborepo. Genera los tipos de TypeScript desde
contracts/engine.schema.json (propón la herramienta; si agrega una
dependencia nueva, escribe el ADR). Aún sin reglas de validación: solo la
estructura, los tipos y una prueba de humo. Explícame cada archivo de
configuración. Primero propón el plan.
```

- **Terminada cuando:** `pnpm test` corre en el paquete y los tipos salen del contrato.

### 1.4.2 — Reglas del validador

- **Modelo:** `opus` / `high`
- **Necesitas:** sesión 1.4.1 terminada.

**Prompt:**

```
Sprint 1.4, sesión 2: implementa en packages/validator la función
validatePlan(input, plan) que revise, SIN reutilizar código del motor:
- ninguna pieza se sale de su tramo ni de su margen (overhang), incluido un
  zMm negativo; el sectionId existe y outOfGauge es correcto (ADR-0017);
- ninguna pieza se encima con otra y ninguna queda flotando;
- cada pieza usa una orientación permitida, sus dx/dy/dz coinciden con ella
  y pasa por la puerta si la hay;
- el peso total no excede el máximo;
- toda pieza de la entrada aparece colocada o en no colocadas, una sola vez;
- el motivo de cada pieza no colocada es el que toca (ADR-0016) y las
  métricas cuadran.
Las mismas reglas están, provisionales, en scripts/contratos/check-examples.mjs:
cuando los ejemplos de contracts/examples/ pasen con el validador, borra ese
script, su línea en package.json y su mención en CLAUDE.md (P-010).
Debe devolver la lista completa de problemas encontrados, no solo el
primero. Pruebas con planes correctos y con planes rotos a propósito.
Primero propón el plan.
```

- **Terminada cuando:** cada regla tiene una prueba que la rompe y la detecta.

### 1.4.3 — Banco de pruebas con instancias BR

- **Modelo:** `opus` / `high`
- **Necesitas:** las instancias BR1 a BR15 de Bischoff y Ratcliff descargadas de la OR-Library en
  `docs/datos/br/` (Claude Code te ayuda a encontrarlas si no las tienes).

**Prompt:**

```
Sprint 1.4, sesión 3: convierte las instancias BR de docs/datos/br/ a
fixtures en engine/fixtures/br/ siguiendo el contrato (documenta el formato
original y la conversión en engine/fixtures/README.md). Luego crea una
herramienta que corra el motor en Go sobre TODOS los fixtures (casos y BR),
guarde los planes y los revise con packages/validator. Cualquier error del
validador hace fallar la prueba. Primero propón el plan.
```

- **Terminada cuando:** todos los planes del banco pasan el validador.

### 1.4.4 — Benchmarks y línea base

- **Modelo:** `sonnet` / `medium`
- **Necesitas:** sesión 1.4.3 terminada.

**Prompt:**

```
Sprint 1.4, sesión 4: agrega benchmarks con testing.B para Optimize sobre
fixtures pequeños, medianos y grandes. Genera una tabla con tiempo y
ocupación de volumen por instancia BR y por caso, y regístrala en la
bitácora como línea base. Explícame cómo leer la salida de go test -bench.
Primero propón el plan.
```

- **Terminada cuando:** la línea base está en la bitácora.

### 1.4.5 — Interfaces, pruebas y cierre

- **Modelo:** `opus` / `medium`

**Prompt:**

```
Escribe los capítulos 07 (interfaces) y 08 (pruebas de tabla, benchmarks y
fuzzing) de la guía de Go con ejemplos del código actual. En el 07, muestra
cómo las estrategias de ordenamiento podrían ser una interfaz y compáralo
con interface/implements de PHP.
```

Después: `/cerrar-sprint 1.4`.

## Sprint 1.5 — Endpoint de optimización

### 1.5.1 — Contrato del endpoint

- **Modelo:** `opus` / `medium`

**Prompt:**

```
/cambiar-contrato agregar POST /plans/optimize que recibe la entrada del
motor y devuelve el LoadPlan, con errores 400 (entrada inválida, con la
lista de problemas), 413 (demasiadas piezas o cuerpo muy grande) y 504
(el cálculo superó el tiempo máximo)
```

Esta sesión solo cambia el contrato (`/cambiar-contrato` se detiene antes del código).

- **Terminada cuando:** `contracts/openapi.yaml` describe el endpoint con ejemplos.

### 1.5.2 — Implementación del endpoint

- **Modelo:** `opusplan` / `medium`
- **Necesitas:** sesión 1.5.1 terminada.

**Prompt:**

```
Sprint 1.5, sesión 2: implementa POST /plans/optimize en
backend/internal/httpapi según contracts/openapi.yaml.
- Límite de tamaño del cuerpo y de número de piezas, configurables por
  variables de entorno con valores por defecto razonables.
- Tiempo máximo de cálculo con context.WithTimeout; el motor debe
  respetar la cancelación.
- Traduce los errores de validación del motor a 400 con la lista de
  problemas, usando errors.Is.
- Pruebas con httptest para cada código de respuesta.
- Agrega las variables nuevas a .env.example y a docker-compose.yml.
Primero propón el plan.
```

- **Terminada cuando:** las pruebas cubren 200, 400, 413 y 504.

### 1.5.3 — Tipos y cliente generados

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Sprint 1.5, sesión 3: genera packages/types y packages/api-client desde
contracts/openapi.yaml (propón la herramienta y escribe el ADR si agrega
una dependencia). Agrega un script de pnpm que los regenere y una
verificación en la CI que falle si el código generado no está al día con
el contrato. Primero propón el plan.
```

- **Terminada cuando:** cambiar el contrato y regenerar actualiza los tipos sin editar nada a mano.

### 1.5.4 — HTTP en Go y cierre

- **Modelo:** `opus` / `medium`

**Prompt:**

```
Escribe el capítulo 10 de la guía de Go (HTTP con la librería estándar y
JSON) con el endpoint /plans/optimize como ejemplo: enrutador, handlers,
decodificar y codificar JSON, códigos de estado y context.
```

Después: `/cerrar-sprint 1.5`.

## Sprint 1.6 — Visor web mínimo

### 1.6.1 — Estructura de la web

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Sprint 1.6, sesión 1: crea apps/web con Vite, React y TypeScript estricto
dentro del workspace. Una sola pantalla vacía con el nombre del producto,
que llame a GET /health y muestre la versión del motor. La URL de la API
viene de una variable de entorno de Vite. Conecta el sistema de diseño
@cargo-3d/ui (estilos, Logo y Button) desde esta primera pantalla, como pide
el ADR-0014 (P-015). Explícame la estructura de una
app de React como si viniera de PHP con plantillas. Primero propón el plan.
```

- **Terminada cuando:** la web muestra la versión del motor que responde la API.

### 1.6.2 — Paquete del visor 3D

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Sprint 1.6, sesión 2: crea packages/viewer-3d con React Three Fiber. Un
componente que reciba un LoadPlan (tipo generado del contrato) y dibuje el
equipo como caja de alambre y cada pieza como caja sólida con color por
tipo, etiqueta con id y peso al pasar el mouse, y controles para girar y
acercar. Respeta el sistema de coordenadas documentado en
docs/src/conceptos/. Agrega una página de ejemplo con un plan de
contracts/examples/. Primero propón el plan.
```

- **Terminada cuando:** el ejemplo se ve igual que el dibujo del sistema de coordenadas.

### 1.6.3 — Captura de la carga y cálculo

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Sprint 1.6, sesión 3: en apps/web agrega la captura de la carga: elegir
el equipo del catálogo, capturar piezas en una tabla y pegar filas
copiadas desde Excel (tabuladores). Botón "Calcular" que llama a
POST /plans/optimize con el api-client generado, muestra los errores 400
de forma entendible y dibuja el resultado con viewer-3d junto con las
métricas. Primero propón el plan y los estados de la pantalla (vacío,
calculando, error, resultado).
```

- **Terminada cuando:** pegas un caso simulado desde Excel y lo ves acomodado.

### 1.6.4 — Validador en la web y revisión

- **Modelo:** `opus` / `medium`

**Prompt:**

```
Sprint 1.6, sesión 4: después de recibir el plan, revísalo con
packages/validator en la web y muestra una alerta clara si el validador
encuentra problemas (no debería pasar, pero si pasa queremos saberlo).
Luego revisa toda la web como un ingeniero senior: accesibilidad básica,
manejo de errores de red y textos en español.
```

Después: `/cerrar-sprint 1.6`.

## Sprint 1.7 — Integración continua, Docker y demo

### 1.7.1 — GitHub Actions

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Sprint 1.7, sesión 1: .github/workflows/ci.yml ya revisa Go (gofmt, go vet,
go test -race) y construye la imagen Docker en cada push (ADR-0013).
Agrégale un job de TypeScript: pnpm install, lint, pruebas y build de los
paquetes y de apps/web, y una verificación de que los tipos generados están
al día con los contratos. Explícame cada paso nuevo del archivo.
Primero propón el plan.
```

- **Terminada cuando:** el push a `develop` muestra todos los jobs del CI en verde.

### 1.7.2 — Demo publicada

- **Modelo:** `sonnet` / `high`

**Prompt:**

```
Sprint 1.7, sesión 2: ayúdame a publicar la demo. Compara 2 o 3 opciones
para el backend en contenedor (por ejemplo, un VPS con Docker Compose o un
servicio administrado de contenedores) y 2 opciones de hosting estático
para la web, con costo mensual aproximado en pesos y complejidad. Cuando
elija, escribe el ADR y los archivos y pasos de despliegue en
docs/src/operacion/despliegue.md. No ejecutes nada que cobre o publique
sin preguntarme.
```

- **Terminada cuando:** abres la demo desde tu celular.

### 1.7.3 — Cierre de la fase 1

- **Modelo:** `opus` / `high`

**Prompt:**

```
/cerrar-sprint 1.7
```

Y al final:

```
Revisa el criterio de terminado de la fase 1 en docs/src/plan.md punto por
punto, con evidencia (resultados de pruebas, tiempos medidos, línea base).
Dime claramente si la fase 1 está terminada o qué falta. Si está
terminada, actualiza la hoja de ruta y la sección Estado de CLAUDE.md.
```

### 1.7.4 — Publicar la versión v0.1.0

- **Modelo:** `sonnet` / `medium`
- **Necesitas:** la fase 1 terminada (sesión anterior) y el CI en verde en `develop`.

**Prompt:**

```
Prepara la versión v0.1.0 siguiendo el ADR-0013. En CHANGELOG.md mueve
todo lo de "## [Sin publicar]" a "## [0.1.0] - <fecha de hoy>", deja una
sección "## [Sin publicar]" vacía arriba y revisa que el texto describa lo
que entrega esta versión. Es la primera versión: revisa también que README.md
describa lo que ya funciona. Corre gofmt, go vet y go test.
No hagas commit, merge, etiqueta ni push: dame los comandos exactos en orden
(commit en develop, merge --ff-only a main, etiqueta v0.1.0, push de main
y de la etiqueta, regreso a develop) y lo que debo revisar en GitHub
(la ejecución de release.yml, la imagen en ghcr.io y el GitHub Release).
```

- **Terminada cuando:** existe el GitHub Release `v0.1.0` con sus notas y la imagen
  `ghcr.io/bayalad7/cargo-3d/api:0.1.0` está publicada.

---

# Fase 2 — MVP y piloto

## Sprint 2.1 — Restricciones del patio

### 2.1.1 — Diseño de las restricciones

- **Modelo:** `opusplan` / `high`

**Prompt:**

```
Sprint 2.1, sesión 1: diseño de las restricciones, sin código. Propón una
interfaz Constraint de Go que permita agregar restricciones al motor sin
modificar el algoritmo. Explica cómo cada restricción (apoyo, apilabilidad,
orden de descarga, peso por eje) se conecta con la evaluación de puntos
extremos, qué datos nuevos necesita el contrato y qué costo de rendimiento
tiene. Escribe docs/src/conceptos/restricciones.md y el ADR.
```

- **Terminada cuando:** el ADR está aceptado y sabes qué campos nuevos lleva el contrato.

### 2.1.2 — Apoyo y estabilidad

- **Modelo:** `opus` / `high`

**Prompt:**

```
Sprint 2.1, sesión 2: primero /cambiar-contrato si hace falta el
porcentaje mínimo de apoyo. Luego implementa la restricción de apoyo:
ninguna pieza puede quedar con menos de ese porcentaje de su base apoyada
(en el piso o sobre otras piezas), usando Contact del sprint 1.2. Pruebas
de tabla y agrega la regla también al validador de TypeScript.
Primero propón el plan.
```

### 2.1.3 — Apilabilidad

- **Modelo:** `opus` / `high`

**Prompt:**

```
Sprint 2.1, sesión 3: restricción de apilabilidad: piezas no apilables
(nada encima), peso máximo que soporta cada pieza encima y piezas frágiles
siempre arriba. Cambia el contrato primero si hace falta. Implementa en el
motor y en el validador, con pruebas de tabla y con al menos un caso de
docs/src/casos/restricciones/ que lo use. Primero propón el plan.
```

### 2.1.4 — Orden de descarga por destino

- **Modelo:** `opus` / `high`

**Prompt:**

```
Sprint 2.1, sesión 4: orden de descarga por destino. Las piezas de la
última parada van al fondo y las de la primera, cerca de la puerta, sin
que una pieza de una parada posterior bloquee a una anterior. Cambia el
contrato primero. Implementa en el motor y en el validador, con pruebas.
Primero propón el plan y define exactamente qué significa "bloquear".
```

### 2.1.5 — Peso por eje (NOM-012)

- **Modelo:** `opus` / `high`

**Prompt:**

```
Sprint 2.1, sesión 5: reparto de peso por eje en cajas secas y
plataformas. Las tablas A-2 y B-2 de la NOM-012 ya están leídas del DOF en
docs/src/conceptos/equipos-y-normas.md; antes de programar, revisa en
dof.gob.mx si hubo una modificación después del 2026-11-21 (había una
programada para 2026) y preséntame las fuentes; no programes límites
que no puedas citar. Propón cómo modelar ejes y distancias en el
catálogo de equipos, cambia el contrato, e implementa en el motor y el
validador. Primero propón el plan.
```

- **Terminada cuando:** cada límite programado tiene su fuente citada en la documentación.

### 2.1.6 — Equipos con varios tramos

- **Modelo:** `opus` / `high`

**Prompt:**

```
Sprint 2.1, sesión 6: plataformas con varios tramos (escalonada, doble caída,
cama baja y RGN; P-012; decide si la nodriza entra). Sube sections.maxItems en contracts/engine.schema.json,
decide en un ADR cómo se aplica la puerta y si una pieza puede apoyarse en dos
tramos, agrega ejemplos en contracts/examples/road/ y el caso inválido de tramos
encimados, y haz que el motor, el validador y el peso por eje los respeten.
Primero propón el plan.
```

- **Terminada cuando:** hay un ejemplo válido de cada tipo y el validador no encuentra errores.

### 2.1.7 — Cierre

**Prompt:** `/cerrar-sprint 2.1`.

## Sprint 2.2 — Metaheurística en paralelo

### 2.2.1 — Diseño de la metaheurística

- **Modelo:** `opus` / `xhigh`

**Prompt:**

```
Sprint 2.2, sesión 1: compara algoritmo genético, búsqueda local y
recocido simulado para mejorar el orden y las orientaciones sobre los
puntos extremos. Para cada uno: cómo funciona, qué tan fácil es de
entender y mantener, cómo se paraleliza y qué mejora esperar según la
literatura. Recomienda uno, escribe docs/src/conceptos/metaheuristica.md
y el ADR. Sin código.
```

### 2.2.2 — Implementación secuencial

- **Modelo:** `opus` / `high`

**Prompt:**

```
Sprint 2.2, sesión 2: implementa la metaheurística elegida, primero SIN
concurrencia. Usa un generador aleatorio con semilla configurable para que
los resultados se puedan repetir. Tiempo límite respetando context.
Compara contra la línea base del sprint 1.4 y registra los resultados en la
bitácora. Primero propón el plan.
```

### 2.2.3 — Paralelo con goroutines

- **Modelo:** `opus` / `xhigh` (o `fable` si se atora)

**Prompt:**

```
Sprint 2.2, sesión 3: paraleliza la metaheurística con goroutines usando
todos los núcleos disponibles, sin compartir estado mutable entre ellas.
Con la misma semilla, el resultado debe ser el mismo. Todas las pruebas
deben pasar con go test -race. Explícame cada goroutine, canal y uso de
context antes de escribirlo. Primero propón el plan.
```

- **Terminada cuando:** `go test -race` pasa y el resultado es repetible con la misma semilla.

### 2.2.4 — Concurrencia y cierre

**Prompt:**

```
Escribe el capítulo 09 de la guía de Go (goroutines, canales y context)
con la metaheurística como ejemplo, incluidos los errores típicos de
concurrencia y cómo los detecta -race.
```

Después: `/cerrar-sprint 2.2`.

## Sprint 2.3 — PostgreSQL y autenticación

### 2.3.1 — Esquema de la base de datos

- **Modelo:** `opus` / `high`

**Prompt:**

```
Sprint 2.3, sesión 1: diseña el esquema de PostgreSQL para empresas,
usuarios, catálogo de equipos, catálogo de piezas frecuentes y planes de
carga (el plan completo en JSONB). Cada tabla de datos lleva la empresa a
la que pertenece. Compara 2 herramientas de migraciones para Go y escribe
el ADR. Muéstrame el diagrama de tablas en texto antes de escribir SQL.
```

### 2.3.2 — Migraciones, sqlc y pgx

- **Modelo:** `sonnet` / `high`

**Prompt:**

```
Sprint 2.3, sesión 2: escribe las migraciones, configura sqlc y pgx, y
conecta el backend a PostgreSQL con DATABASE_URL (con un pool de
conexiones y un /health que revise la base de datos). Las pruebas de base
de datos corren contra el PostgreSQL de docker compose. Explícame qué
genera sqlc y por qué no usamos un ORM. Primero propón el plan.
```

### 2.3.3 — Guardar y consultar planes

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
/cambiar-contrato agregar endpoints para guardar un plan calculado,
listarlos con paginación, consultar uno y duplicarlo
```

Luego, en la misma rama: *"Implementa esos endpoints con sqlc, con pruebas contra la base de datos."*

### 2.3.4 — Autenticación

- **Modelo:** `opus` / `high`

**Prompt:**

```
Sprint 2.3, sesión 4: autenticación con correo y contraseña. Compara
sesiones con cookie contra tokens, recomienda una para web y móvil y
escribe el ADR. Contraseñas con un algoritmo de hash moderno, límite de
intentos y middleware que ponga la empresa del usuario en el context.
Cambia el contrato primero. Primero propón el plan.
```

### 2.3.5 — Separación de datos entre empresas

- **Modelo:** `opus` / `high`

**Prompt:**

```
Sprint 2.3, sesión 5: asegúrate de que ninguna consulta pueda leer o
modificar datos de otra empresa. Revisa cada consulta de sqlc y cada
endpoint, y escribe pruebas que intenten acceder a datos de otra empresa
en cada endpoint. Este es el error más grave posible en el producto:
revísalo como auditor de seguridad.
```

- **Terminada cuando:** existe una prueba de acceso cruzado por cada endpoint y todas fallan como deben.

### 2.3.6 — PostgreSQL en Go y cierre

**Prompt:** *"Escribe el capítulo 11 de la guía de Go (PostgreSQL con pgx y sqlc)."*
Después: `/cerrar-sprint 2.3`.

## Sprint 2.4 — La web como producto

### 2.4.1 — Importación de Excel

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Sprint 2.4, sesión 1: importación de archivos .xlsx y .csv en la web con
mapeo de columnas ("esta columna es el largo"), detección de unidades
sospechosas (por ejemplo, medidas que parecen centímetros) y vista previa
antes de importar. Usa mis archivos reales de docs/datos/ para probar.
Primero propón el plan y la librería (con ADR si es dependencia nueva).
```

### 2.4.2 — Editor manual

- **Modelo:** `opus` / `medium`

**Prompt:**

```
Sprint 2.4, sesión 2: editor manual en el visor 3D: seleccionar una pieza,
moverla, girarla o quitarla. Cada cambio se revisa al instante con
packages/validator y los problemas se marcan en rojo en la pieza. Botón
para guardar el plan editado. Primero propón el plan de interacción.
```

### 2.4.3 — Historial y duplicar

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Sprint 2.4, sesión 3: pantalla de historial de planes con búsqueda y
filtros por fecha y equipo, abrir un plan guardado y duplicarlo para una
carga parecida. Usa los endpoints del sprint 2.3. Primero propón el plan.
```

### 2.4.4 — Varios equipos por carga

- **Modelo:** `opusplan` / `high`

**Prompt:**

```
Sprint 2.4, sesión 4: cuando la carga no cabe en un equipo, el motor debe
repartirla en varios (del mismo tipo o de una lista de tipos disponibles)
intentando usar el menor número. Cambia el contrato primero, luego el
motor, el validador, la API y el visor (pestañas por equipo).
Primero propón el plan.
```

Después: `/cerrar-sprint 2.4`.

## Sprint 2.5 — El plan en papel

### 2.5.1 — Dónde generar el PDF

- **Modelo:** `opus` / `medium`

**Prompt:**

```
Sprint 2.5, sesión 1: compara generar el PDF del plan en la web contra
generarlo en el backend de Go: calidad de la vista 3D, mantenimiento,
dependencias y posibilidad de mandarlo por WhatsApp o correo. Recomienda
una opción y escribe el ADR.
```

### 2.5.2 — PDF del plan de carga

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Sprint 2.5, sesión 2: genera el PDF del plan de carga según el ADR: portada
con datos de la carga y métricas, vista 3D, paso a paso por capas o en el
orden de carga, lista de piezas con su posición y pesos totales y por eje.
Tamaño carta, legible impreso en blanco y negro. Primero propón el diseño
de cada página.
```

Después: `/cerrar-sprint 2.5`.

## Sprint 2.6 — Despliegue y pruebas con cargas reales

### 2.6.1 — Despliegue de producción

- **Modelo:** `sonnet` / `high`

**Prompt:**

```
Sprint 2.6, sesión 1: prepara el despliegue de producción en un VPS con
Docker Compose: API, PostgreSQL, un proxy con HTTPS automático, variables
de entorno de producción y firewall. Escribe
docs/src/operacion/produccion.md paso a paso. No ejecutes nada en el
servidor sin mostrarme el comando y esperar mi confirmación.
```

### 2.6.2 — Respaldos y restauración

- **Modelo:** `sonnet` / `high`

**Prompt:**

```
Sprint 2.6, sesión 2: respaldos diarios de PostgreSQL fuera del servidor,
con retención de 30 días, y un procedimiento de restauración escrito en
docs/src/operacion/respaldos.md. Guíame para hacer una restauración de
prueba en una base de datos aparte y regístrala en la bitácora.
```

- **Terminada cuando:** restauraste un respaldo real y está anotado en la bitácora.

### 2.6.3 — Material de capacitación

- **Modelo:** `opus` / `medium`

**Prompt:**

```
Sprint 2.6, sesión 3: escribe una guía de uso de 2 páginas para quien
planea las cargas en la empresa piloto (sin tecnicismos) y prepara el
registro de pruebas de docs/src/casos/pruebas-reales.md (carpeta privada
docs/datos/pruebas/<empresa>/ con empresa.md, equipos.csv, reglas.md y
registro.csv).
```

### 2.6.4 — Seguimiento semanal del piloto

- **Modelo:** `sonnet` / `medium`

**Prompt (una vez por semana durante 4 semanas):**

```
Seguimiento de las pruebas con cargas reales, semana N. Lee el
registro.csv de cada empresa en docs/datos/pruebas/. Resume tiempos,
cargas rehechas y comentarios; aplica docs/src/casos/pruebas-reales.md
(caso validado, corregido o nuevo, anonimizado) y convierte cada problema
en una tarea priorizada; regístralo en la bitácora.
No cambies código en esta sesión.
```

### 2.6.5 — Cierre de la fase 2

- **Modelo:** `opus` / `high`

**Prompt:**

```
/cerrar-sprint 2.6
```

Y al final: *"Revisa el criterio de terminado de la fase 2 punto por punto con evidencia del piloto y dime si
está terminada. Si lo está, actualiza la hoja de ruta y el Estado de CLAUDE.md."*

### 2.6.6 — Publicar la versión v0.2.0

- **Modelo:** `sonnet` / `medium`
- **Necesitas:** la fase 2 terminada (sesión anterior) y el CI en verde en `develop`.

**Prompt:**

```
Prepara la versión v0.2.0 siguiendo el ADR-0013. En CHANGELOG.md mueve
todo lo de "## [Sin publicar]" a "## [0.2.0] - <fecha de hoy>", deja una
sección "## [Sin publicar]" vacía arriba y revisa que el texto describa lo
que entrega esta versión. Es la versión que usan las empresas en el
periodo de pruebas: incluye en las notas cómo actualizar el servidor del piloto. Corre gofmt, go vet y go test.
No hagas commit, merge, etiqueta ni push: dame los comandos exactos en orden
(commit en develop, merge --ff-only a main, etiqueta v0.2.0, push de main
y de la etiqueta, regreso a develop) y lo que debo revisar en GitHub
(la ejecución de release.yml, la imagen en ghcr.io y el GitHub Release).
```

- **Terminada cuando:** existe el GitHub Release `v0.2.0` con sus notas y la imagen
  `ghcr.io/bayalad7/cargo-3d/api:0.2.0` está publicada.

---

# Fase 3 — Producto y lanzamiento

## Bloque 3.1 — Multiempresa y permisos

### 3.1.1 — Roles y permisos

- **Modelo:** `opus` / `high`

**Prompt:**

```
Bloque 3.1, sesión 1: roles por empresa (administrador, planeador,
operador de patio) con una tabla de qué puede hacer cada uno. Cambia el
contrato primero, implementa la revisión de permisos en el middleware y
escribe pruebas que intenten cada acción con cada rol. Primero propón la
tabla de permisos.
```

### 3.1.2 — Invitaciones y cuentas

- **Modelo:** `sonnet` / `high`

**Prompt:**

```
Bloque 3.1, sesión 2: invitaciones por correo con enlace que expira,
recuperación de contraseña y registro de actividad (quién creó, editó o
borró cada plan). Propón el proveedor de correo (con ADR). Cambia el
contrato primero. Primero propón el plan.
```

### 3.1.3 — Límites por suscripción

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Bloque 3.1, sesión 3: límites por plan de suscripción (usuarios y planes
calculados por mes), con mensajes claros cuando se alcanzan y sin bloquear
el acceso a planes ya guardados. Primero propón el plan.
```

### 3.1.4 — Cierre del bloque 3.1

- **Modelo:** `opus` / `high`

**Prompt:**

```
/cerrar-sprint 3.1
```

## Bloque 3.2 — App móvil de patio

### 3.2.1 — Estructura y sesión

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Bloque 3.2, sesión 1: crea apps/mobile con Expo y TypeScript dentro del
workspace, con inicio de sesión de operador de patio usando el api-client
generado. Genera los tokens de Voxa para React Native a partir de
packages/ui/src/styles/tokens.css (fuente única, sin copiar valores a mano) y
los íconos de la app desde packages/ui/src/brand/ (P-016). Explícame la
diferencia entre React para web y React Native.
Primero propón el plan.
```

### 3.2.2 — Plan del día y paso a paso

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Bloque 3.2, sesión 2: pantalla con los planes del día y vista de paso a
paso: qué pieza sigue, dónde va (vista 3D simple o diagrama por capas) y
botón grande para marcarla como cargada. Pensado para usarse con guantes
y bajo el sol: botones grandes y alto contraste. Primero propón las
pantallas.
```

### 3.2.3 — Evidencia y modo sin conexión

- **Modelo:** `opusplan` / `high`

**Prompt:**

```
Bloque 3.2, sesión 3: fotos de evidencia por paso y del sello del
contenedor, y modo sin conexión: el plan se guarda en el teléfono, los
avances y fotos se encolan y se sincronizan al recuperar señal sin
duplicarse. Cambia el contrato para la sincronización. Diseña primero
cómo se resuelven los conflictos.
```

### 3.2.4 — Publicación

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Bloque 3.2, sesión 4: guíame para publicar en Google Play (y después App
Store): cuentas necesarias, íconos, capturas, política de privacidad y
build con Expo. Escribe docs/src/operacion/publicacion-movil.md. No
publiques nada sin mi confirmación.
```

### 3.2.5 — Cierre del bloque 3.2

- **Modelo:** `opus` / `high`

**Prompt:**

```
/cerrar-sprint 3.2
```

## Bloque 3.3 — Cobro y CFDI

### 3.3.1 — Elección de pasarela y PAC

- **Modelo:** `opus` / `high`

**Prompt:**

```
Bloque 3.3, sesión 1: compara Stripe, Mercado Pago y Conekta para
suscripciones en pesos (tarjetas mexicanas, OXXO o transferencia,
comisiones vigentes, facilidad de integración en Go) y 2 o 3 PAC para
emitir CFDI por API. Busca información actual y cita las fuentes.
Recomienda y escribe los ADR. Sin código.
```

### 3.3.2 — Suscripciones

- **Modelo:** `opus` / `high`

**Prompt:**

```
Bloque 3.3, sesión 2: integra la pasarela elegida: planes, periodo de
prueba, cobro recurrente, cancelación sin llamar a nadie y webhooks
verificados con su firma. Todo en modo de pruebas de la pasarela.
Primero propón el plan y los estados de una suscripción.
```

### 3.3.3 — Facturación CFDI

- **Modelo:** `opus` / `high`

**Prompt:**

```
Bloque 3.3, sesión 3: emisión automática de CFDI con el PAC elegido al
confirmarse cada pago, con datos fiscales del cliente, reintentos si el PAC
falla y envío del PDF y XML por correo. Todo en ambiente de pruebas del
PAC. Antes de programar, dime qué datos fiscales necesito confirmar con mi
contador.
```

### 3.3.4 — Cierre del bloque 3.3

- **Modelo:** `opus` / `high`

**Prompt:**

```
/cerrar-sprint 3.3
```

## Bloque 3.4 — Operación profesional

### 3.4.1 — Monitoreo y ambientes

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Bloque 3.4, sesión 1: monitoreo de errores en backend, web y móvil (propón
servicio con ADR), y ambientes separados de desarrollo, pruebas y
producción con sus variables. Documenta en docs/src/operacion/.
```

### 3.4.2 — Revisión de seguridad

- **Modelo:** `opus` / `max`

**Prompt:**

```
Bloque 3.4, sesión 2: revisión de seguridad completa del proyecto como un
auditor: dependencias con vulnerabilidades conocidas, manejo de contraseñas
y sesiones, límites de peticiones, validación de entradas, separación de
datos entre empresas, secretos, encabezados HTTP y configuración de Docker.
Entrega una lista priorizada de hallazgos con cómo reproducir cada uno. No
corrijas nada todavía.
```

Luego, por cada hallazgo importante, una sesión: *"Corrige el hallazgo N con su prueba."*

### 3.4.3 — Recuperación ante desastres

- **Modelo:** `sonnet` / `high`

**Prompt:**

```
Bloque 3.4, sesión 3: escribe el plan de recuperación ante desastres
(servidor perdido, base de datos dañada, secreto filtrado) con pasos
concretos y tiempos objetivo, y guíame en un simulacro de recuperación
completa en un servidor nuevo. Regístralo en la bitácora.
```

### 3.4.4 — Cierre del bloque 3.4

- **Modelo:** `opus` / `high`

**Prompt:**

```
/cerrar-sprint 3.4
```

## Bloque 3.5 — Escritorio con Wails (solo si un cliente lo pide)

### 3.5.1 — App de escritorio

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Bloque 3.5, sesión 1: crea apps/desktop con Wails v3 (fija una versión
exacta) reutilizando apps/web. El motor de internal/engine se llama
directamente desde Go para calcular sin conexión. Genera los íconos de la app
desde packages/ui/src/brand/ (P-016). Revisa la documentación
actual de Wails antes de empezar. Primero propón el plan.
```

### 3.5.2 — Instaladores

- **Modelo:** `sonnet` / `medium`

**Prompt:**

```
Bloque 3.5, sesión 2: instaladores firmados para Windows (y macOS si hace
falta) construidos en la CI. Explícame qué certificados se necesitan y
cuánto cuestan antes de configurar nada.
```

### 3.5.3 — Cierre del bloque 3.5

- **Modelo:** `opus` / `high`
- **Solo si se hizo el bloque 3.5.**

**Prompt:**

```
/cerrar-sprint 3.5
```

## Bloque 3.6 — Salida al mercado

### 3.6.1 — Página de producto

- **Modelo:** `opus` / `medium`

**Prompt:**

```
Bloque 3.6, sesión 1: página de producto con propuesta de valor, demo,
precios en pesos y un caso real con números del piloto (solo con permiso
de la empresa piloto). Textos en español claro para gente de logística,
sin tecnicismos. Primero propón la estructura y los textos.
```

### 3.6.2 — Documentos legales

- **Modelo:** `opus` / `medium`

**Prompt:**

```
Bloque 3.6, sesión 2: borradores de aviso de privacidad, términos y
condiciones y contrato de servicio para revisión de un abogado. Deben dejar
claro que el plan de carga es una sugerencia y que la carga física y el
cumplimiento de la NOM son responsabilidad de quien carga. Márcalos como
BORRADOR PARA REVISIÓN LEGAL.
```

### 3.6.3 — Beta y soporte

- **Modelo:** `opus` / `medium`

**Prompt:**

```
Bloque 3.6, sesión 3: plan de beta cerrada con 5 a 10 empresas del puerto:
mensaje de invitación, proceso de alta, encuesta de seguimiento y base de
preguntas frecuentes para el soporte por WhatsApp con horario definido.
```

### 3.6.4 — Cierre de la fase 3

- **Modelo:** `opus` / `high`

**Prompt:**

```
/cerrar-sprint 3.6
```

Y al final:

```
Revisa el criterio de terminado de la fase 3 en docs/src/plan.md punto por
punto con evidencia. Si está terminada, actualiza la hoja de ruta, el
Estado de CLAUDE.md y propón el plan de la fase 4 con el mismo formato.
```

### 3.6.5 — Publicar la versión v1.0.0

- **Modelo:** `sonnet` / `medium`
- **Necesitas:** la fase 3 terminada (sesión anterior) y el CI en verde en `develop`.

**Prompt:**

```
Prepara la versión v1.0.0 siguiendo el ADR-0013. En CHANGELOG.md mueve
todo lo de "## [Sin publicar]" a "## [1.0.0] - <fecha de hoy>", deja una
sección "## [Sin publicar]" vacía arriba y revisa que el texto describa lo
que entrega esta versión. Es el lanzamiento: revisa que README.md, la página de
producto y los documentos legales coincidan con lo que entrega. Corre gofmt, go vet y go test.
No hagas commit, merge, etiqueta ni push: dame los comandos exactos en orden
(commit en develop, merge --ff-only a main, etiqueta v1.0.0, push de main
y de la etiqueta, regreso a develop) y lo que debo revisar en GitHub
(la ejecución de release.yml, la imagen en ghcr.io y el GitHub Release).
```

- **Terminada cuando:** existe el GitHub Release `v1.0.0` con sus notas y la imagen
  `ghcr.io/bayalad7/cargo-3d/api:1.0.0` está publicada.

---

# Prompts de apoyo (cualquier momento)

| Situación | Prompt |
| --- | --- |
| No entiendes un error del compilador | `/explicar` seguido del error pegado |
| No entiendes un archivo o función | `/explicar backend/internal/engine/placement.go` |
| Quieres aprender un concepto de Go | `/explicar goroutines con el código del motor` |
| Claude Code propone algo que no te convence | *"Detente. Dame 2 alternativas con pros y contras antes de seguir."* |
| La sesión se alargó o se confundió | `/cerrar-sesion`, luego `/clear` y empieza otra con `/iniciar-sesion` |
| Algo que funcionaba se rompió | *"Corre las pruebas, encuentra qué commit lo rompió con git log y git diff, y explícame la causa antes de arreglar."* |
| Quieres escribir tú el código | *"Voy a escribir yo la función X. Dime qué debe hacer y qué pruebas debe pasar, pero no me des el código. Después revísalo."* |
| Tomaste una decisión importante | `/nuevo-adr título de la decisión` |
| Hay que cambiar la API o el motor | `/cambiar-contrato descripción del cambio` |
| Falla algo urgente en una versión ya liberada | *"Hay un error en la versión vX.Y.Z: <qué pasa>. Reprodúcelo con una prueba que falle, corrígelo en develop y prepara la versión de parche vX.Y.Z+1 como en el ADR-0013 (CHANGELOG, pruebas y los comandos de merge, etiqueta y push, sin ejecutarlos)."* |
| Instalar una tecnología nueva o subir de versión | *"Revisa en sus fuentes oficiales la última versión estable de <tecnología> (y la LTS de Node), compárala con el ADR-0011, dime qué cambios incompatibles trae y qué archivos hay que tocar. Si es una versión mayor, va en su propio commit con las pruebas pasando. Actualiza la tabla del ADR-0011."* |
| Probar algo en Docker | *"Pruébalo siguiendo las Reglas de Docker de CLAUDE.md (proyecto cargo3d-test) y al terminar elimina todo lo que creaste y dime cuánto espacio se liberó."* |
