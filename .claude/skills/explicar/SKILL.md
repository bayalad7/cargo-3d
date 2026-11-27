---
name: explicar
description: Explica a Benjamín un archivo, una función, un error del compilador o un concepto de Go o TypeScript, con el código real del proyecto.
argument-hint: "[archivo, función, error pegado o concepto]"
disable-model-invocation: true
---

# Explicar: $ARGUMENTS

Responde siempre en español. Benjamín es desarrollador y está aprendiendo Go (y TypeScript) con este
proyecto. **No modifiques archivos en esta skill.**

## Qué hacer

1. Si no hay nada que explicar, pregunta qué quiere entender y detente.
2. Identifica qué es lo que pidió:
   - **Un archivo o función:** léelo completo y explícalo de arriba abajo: para qué existe, qué recibe,
     qué devuelve, cómo se conecta con el resto y por qué está hecho así (cita el ADR si existe).
   - **Un error del compilador o de las pruebas:** explica primero qué significa el mensaje palabra por
     palabra, después por qué ocurre en este código y, al final, cómo se arregla. Si es un error que se
     arregla en una línea, muestra la línea, pero no la cambies.
   - **Un concepto:** explícalo con un ejemplo tomado del código real del proyecto (no uno genérico) y
     dile en qué archivo puede verlo.
3. Explica solo en términos de Go (o de TypeScript), sin comparar con otros lenguajes (ADR-0022). Describe lo que
   el compilador y el programa hacen de verdad; si hay duda, compruébalo con un programa en el directorio temporal
   antes de afirmarlo. Señala las trampas reales que apliquen (copias superficiales, `nil`, valores cero).
4. Si el concepto ya tiene capítulo en `docs/src/aprendizaje/go/`, enlázalo. Si no lo tiene y es
   importante, propón agregarlo o ampliarlo (no lo escribas sin permiso).
5. Termina con una pregunta corta para comprobar que se entendió (sin dar la respuesta).
