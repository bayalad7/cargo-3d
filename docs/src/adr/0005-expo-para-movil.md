# ADR-0005: React Native + Expo para la app móvil

- **Estado:** Aceptado (se implementa en la fase 3)
- **Fecha:** 2026-11-20

## Contexto

La app móvil es para el **patio**: el montacarguista o el jefe de patio sigue el plan paso a paso,
marca piezas cargadas, toma fotos de evidencia y escanea etiquetas, a veces con mala señal.

## Decisión

`apps/mobile` usa **React Native con Expo**. Comparte con la web los tipos, el cliente de la API y el
validador en TypeScript. Los planes se calculan en la API en Go; la app descarga el plan del día y lo
guarda en el teléfono para trabajar sin señal.

## Alternativas consideradas

| Opción | A favor | En contra |
|---|---|---|
| **Expo / React Native (elegida)** | UI nativa; cámara y almacenamiento fáciles; comparte TypeScript con la web | Una base de código móvil aparte de la web |
| Wails v3 móvil | Mismo Go del backend | Soporte móvil muy reciente; experiencia de "web en una ventana" |
| Flutter | Excelente UI móvil | Otro lenguaje (Dart) sin nada compartido |
| Nativo (Swift + Kotlin) | Máximo rendimiento | Dos apps que mantener |

## Consecuencias

- **Positivas:** una base de código para Android e iOS.
- **Negativas:** la app móvil no calcula planes sin conexión (los descarga ya calculados), lo cual es
  suficiente para el uso en patio.
- **Revertir:** al ser API-first (ADR-0008), cualquier tecnología móvil puede consumir la misma API.
