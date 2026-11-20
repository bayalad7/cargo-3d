# El problema del acomodo 3D (bin packing)

## Qué problema resolvemos

Tenemos un **espacio** (un contenedor de 40', una caja seca de 53') y una **lista de piezas**
con medidas, peso y reglas. Hay que decidir **dónde va cada pieza y en qué orientación** para:

- meter el máximo posible (o todo, usando el menor número de contenedores),
- sin pasarse del peso permitido,
- respetando las reglas reales: qué se puede apilar, qué no se puede voltear,
  qué se descarga primero, cómo repartir el peso.

## Por qué es difícil

Es un problema **NP-difícil**: el número de combinaciones crece de forma explosiva.
Con solo 20 piezas y 6 orientaciones posibles, el número de órdenes y giros supera por mucho
lo que cualquier computadora podría probar. Nadie lo resuelve de forma perfecta en la práctica;
se usan **heurísticas**: reglas inteligentes que dan soluciones muy buenas en segundos.

## La estrategia que usaremos (fase 1)

**Puntos extremos (Extreme Points):**

1. Empezamos con un solo punto disponible: la esquina del fondo del contenedor `(0, 0, 0)`.
2. Ordenamos las piezas (por ejemplo, de mayor a menor volumen).
3. Para cada pieza, probamos cada punto disponible y cada orientación permitida, y elegimos
   la posición que mejor encaja (menos espacio desperdiciado, apoyo estable).
4. Al colocar la pieza, aparecen **nuevos puntos** en sus esquinas: arriba, a un lado y al frente.
5. Repetimos hasta acabar las piezas o el espacio.

En fases posteriores agregaremos una **metaheurística** encima (por ejemplo, un algoritmo
genético que prueba distintos órdenes de piezas y se queda con el mejor resultado).

## Restricciones que el motor tendrá que entender

| Restricción | Ejemplo real | Fase |
|---|---|---|
| Dimensiones y no traslape | Ninguna pieza se sale ni se encima con otra | 1 |
| Orientaciones permitidas | "Este lado arriba" | 1 |
| Peso máximo del contenedor | Carga útil de un 40' HC | 1 |
| Apoyo / estabilidad | Nada queda "flotando" | 2 |
| Apilabilidad y peso sobre una pieza | No poner tarimas pesadas sobre cajas frágiles | 2 |
| Orden de descarga por destino | Lo de la última parada va al fondo | 2 |
| Reparto de peso por eje | Pesos por eje según la NOM-012 | 2 |
| Centro de gravedad | Evitar cargas cargadas a un lado | 3 |

> Al implementar cada restricción se documenta aquí cómo funciona y se enlaza al código.
