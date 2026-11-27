# Flujo de trabajo con Claude Code

Esta guía explica **cómo se trabaja cada día** en Cargo 3D con Claude Code: los comandos propios del
proyecto (skills), la biblioteca de prompts y el ciclo completo de una sesión, desde abrir la terminal
hasta el `git push`.

## Las tres piezas

| Pieza | Dónde vive | Para qué sirve |
| --- | --- | --- |
| **`CLAUDE.md`** | Raíz del repo | La memoria del proyecto. Claude Code lo lee solo al arrancar: stack, reglas, estado actual. |
| **Biblioteca de prompts** | [`docs/src/prompts.md`](../../prompts.md) | El prompt de cada sesión de cada sprint (fases 1 a 3), con modelo, requisitos y criterio de terminado. |
| **Skills (comandos propios)** | `.claude/skills/<nombre>/SKILL.md` | Comandos que escribes con `/` en Claude Code para los pasos que se repiten: arrancar, ejecutar una sesión, revisar, cerrar. |

`CLAUDE.md` dice **cómo** se trabaja; `prompts.md` dice **qué** se hace en cada sesión; las skills
**automatizan** lo que siempre es igual.

## Los comandos del proyecto

Claude Code carga solo las skills de la carpeta `.claude/skills/` del repo. Al escribir `/` las verás en el
menú junto a los comandos integrados. Todas están marcadas para que **solo tú** las ejecutes: Claude no las
dispara por su cuenta.

| Comando | Cuándo usarlo | Qué hace | ¿Modifica archivos? |
| --- | --- | --- | --- |
| `/iniciar-sesion [meta]` | Al abrir Claude Code, siempre | Lee `CLAUDE.md`, el plan, los pendientes y la última bitácora; revisa la rama y los cambios pendientes; te dice dónde vamos y qué sesión sigue | No |
| `/buscador <sprint o tema>` | **Opcional.** Después de `/iniciar-sesion` y antes de `/sesion`, solo si el sprint necesita datos logísticos que todavía no tienes (si ya están en `docs/datos/`, `contracts/` o `engine/fixtures/`, sáltalo) | Lanza el agente `investigador-logistico` (`.claude/agents/`) en un contexto aparte: busca en fuentes primarias (ISO, NOM-012, CTU Code, VGM, EPAL), cita cada dato con su fuente y nivel de confianza y entrega catálogos, fórmulas, casos borde, escenarios sintéticos y el destino propuesto | No |
| `/sesion 1.1.2` | Para ejecutar una sesión de `prompts.md` | Busca la sesión, revisa requisitos, que estés en `develop` y el modelo; **propone un plan y espera tu aprobación**; luego implementa y verifica | Sí, después de que apruebes |
| `/revisar [enfoque]` | Antes de cerrar una sesión con cambios importantes | Revisión como ingeniero senior de todo lo que cambió contra `main`; hallazgos ordenados por gravedad | No |
| `/cerrar-sesion [nota]` | Al terminar cada sesión | Corre `gofmt`, `go vet` y pruebas; actualiza bitácora, registro de pendientes (`docs/src/pendientes.md`), CHANGELOG, Estado de `CLAUDE.md` y checklist del plan; propone el mensaje de commit | Sí (solo documentación) |
| `/cerrar-sprint 1.1` | Al terminar todas las sesiones de un sprint | Verifica con evidencia cada tarea del sprint; si todo está, lo marca y actualiza el Estado; en el último sprint de una fase revisa el criterio de la fase | Sí (plan, Estado, bitácora) |
| `/explicar <algo>` | Cuando no entiendes algo | Explica un archivo, una función, un error pegado o un concepto, con el código real del proyecto | No |
| `/nuevo-adr <título>` | Cuando se toma una decisión de arquitectura | Crea el ADR con el número siguiente desde la plantilla, lo enlaza y espera tu aceptación | Sí |
| `/cambiar-contrato <cambio>` | Antes de cambiar la API o el motor | Propone el cambio en `contracts/`, la versión nueva y lo que habrá que tocar; lo escribe al aprobarlo y **se detiene antes del código** | Sí (solo `contracts/`) |

Detalles que conviene saber:

- **Ninguna skill hace commit ni push.** Siempre te muestran los comandos y tú los ejecutas. Así revisas
  todo antes de que entre al historial.
- **Algunas skills cambian de modelo solas** durante su ejecución: `/cerrar-sesion` usa Sonnet (es trabajo
  rutinario); `/revisar` y `/cerrar-sprint` usan Opus. Al terminar, la conversación
  vuelve al modelo que tenías.
- **Las skills leen el estado de git al arrancar** (rama, cambios, últimos commits) para no trabajar a ciegas.
  Por eso deben ejecutarse desde la raíz del repo.
- **Para ver o editar una skill**, abre su `SKILL.md`: es Markdown normal con una cabecera de configuración.
  Si cambias una, regístralo en la bitácora como cualquier otro cambio del proyecto.

## El ciclo de una sesión, paso a paso

Ejemplo real: la sesión **1.1.2 — Container y CargoItem**.

### 1. Actualiza `develop` (en tu terminal, fuera de Claude Code)

```bash
git switch develop
git pull
```

Todo se trabaja directo en `develop` (ADR-0012): sin ramas de trabajo ni pull requests.

### 2. Abre Claude Code y elige el modelo

```bash
claude
```

Dentro de Claude Code:

```
/model opus
/effort medium
```

El modelo y el esfuerzo sugeridos vienen en cada sesión de `prompts.md`. Si se te olvida, `/sesion` te lo recuerda.

### 3. Arranca la sesión

```
/iniciar-sesion
```

Lee su resumen. Si dice algo que no cuadra con lo que tú sabes (por ejemplo, que una sesión anterior no
está terminada), corrígelo antes de seguir: esa información viene de la bitácora y de `CLAUDE.md`.

### 4. Ejecuta la sesión

```
/sesion 1.1.2
```

Claude Code:

1. encuentra la sesión en `prompts.md`;
2. revisa que existan los requisitos y que estés en `develop`;
3. **te presenta un plan**: archivos, pruebas, documentación y conceptos nuevos de Go.

**Este es el momento más importante de la sesión.** Lee el plan con calma. Puedes:

- aprobarlo: *"Adelante."*
- corregirlo: *"Usa int64 también para la cantidad"*, *"No agregues todavía el método X"*;
- pedir alternativas: *"Dame otra opción para el sistema de coordenadas y compáralas."*

No escribe código hasta que apruebes.

### 5. Mientras trabaja

- Claude Code te pedirá permiso para editar archivos y correr comandos. Lee qué va a hacer antes de aceptar.
- Si algo no lo entiendes, interrumpe: *"Espera. Explícame esta función antes de seguir."* o usa `/explicar`.
- Si quieres escribir tú una parte: *"La función Volume la escribo yo. Dime qué debe cumplir y después revísala."*
- Si la sesión se alarga mucho o empieza a confundirse, ciérrala (paso 7) y abre otra con `/clear`.

### 6. Revisa (recomendado en sesiones con código)

```
/revisar
```

Pide que se corrijan los hallazgos bloqueantes e importantes: *"Corrige los hallazgos 1 y 3."*

### 7. Cierra la sesión

```
/cerrar-sesion
```

Si las pruebas fallan, se detiene y te explica. Si todo pasa, actualiza la documentación y te propone el
mensaje de commit. Revisa el diff de la documentación también.

### 8. Commit y push (en tu terminal)

```bash
git add -A
git commit            # pega el mensaje propuesto
git push
```

En GitHub (pestaña *Actions*), revisa que el CI quede en verde.

### 9. Siguiente sesión

Vuelves al paso 1. Lo recomendable es **una sesión nueva de Claude Code por sesión de trabajo**
(sal con `/exit` o usa `/clear`): cada sesión arranca limpia y con el contexto que dicen `CLAUDE.md` y la bitácora.

## Al terminar un sprint

Cuando todas las sesiones del sprint están en `develop`:

```bash
git switch develop && git pull
claude
```

```
/model opus
/cerrar-sprint 1.1
/cerrar-sesion
```

`/cerrar-sprint` no deja cerrar un sprint con tareas sin evidencia o con pruebas fallando.

## Buenas prácticas

- **Una meta por sesión.** Si en medio de una sesión surge otra idea, anótala como pendiente; no la mezcles.
- **El plan primero, siempre.** Aunque la tarea parezca sencilla, revisar el plan cuesta un minuto y evita
  horas de código en la dirección equivocada.
- **El contrato primero.** Si un cambio toca la API o el motor, empieza con `/cambiar-contrato`.
- **No aceptes lo que no entiendes.** Si después de `/explicar` sigue sin quedar claro, el problema es la
  documentación: pide que se mejore.
- **Datos reales.** Cuando una sesión dependa de tus datos (cargas, equipos, reglas del patio), tenlos en
  `docs/datos/` antes de empezar.
- **El repositorio es la memoria.** Lo que no quedó en la bitácora o en `CLAUDE.md`, la siguiente sesión no
  lo sabe. Por eso `/cerrar-sesion` no es opcional.
- **Nada pendiente se queda en el chat.** Todo lo abierto (datos que faltan, valores supuestos, decisiones
  aplazadas) va a [Pendientes](../../pendientes.md) con un id `P-NNN` y el «cómo retomarlo». `/iniciar-sesion` te
  muestra primero los que bloquean la siguiente sesión.

## Si algo no funciona

| Problema | Solución |
| --- | --- |
| No aparecen los comandos al escribir `/` | Abre Claude Code desde la raíz del repo (donde está `CLAUDE.md`). Revisa que exista `.claude/skills/<nombre>/SKILL.md` y que la primera línea del archivo sea `---`. Reinicia Claude Code. |
| Una skill se detiene al arrancar | Las skills leen el estado de git al iniciar. Confirma que estás dentro del repo y que `git status` funciona. |
| `/sesion` no encuentra el número | Revisa en `prompts.md` que el encabezado sea exactamente `### 1.1.2 — ...` (tres almohadillas, espacio, número, espacio). |
| La documentación de Claude Code cambió | Las skills siguen el formato de la documentación oficial de Claude Code (noviembre de 2026). Si algo cambia, pídele a Claude Code que revise la documentación actual y adapte los archivos. |
