# Índice

[Introducción](introduccion.md)
[Plan de trabajo](plan.md)
[Pendientes](pendientes.md)
[Biblioteca de prompts por sesión](prompts.md)

# Desarrollo

- [Desarrollo](desarrollo/README.md)
  - [Instalación del entorno](desarrollo/entorno/instalacion.md)
  - [Levantar y detener servicios](desarrollo/entorno/servicios.md)
  - [Base de datos: pgAdmin y psql](desarrollo/entorno/base-de-datos.md)
  - [Flujo con Claude Code](desarrollo/flujo/claude-code.md)
  - [Ramas, commits y versiones](desarrollo/flujo/ramas-y-commits.md)
  - [Convenciones](desarrollo/convenciones.md)

# Conceptos

- [Visión general del sistema](conceptos/vision-general.md)
- [El problema del acomodo 3D](conceptos/problema-bin-packing.md)
- [Sistema de coordenadas](conceptos/sistema-coordenadas.md)
- [Equipos y normas](conceptos/equipos-y-normas.md)
- [Equipos personalizados](conceptos/equipos-personalizados.md)

# Catálogos

- [Catálogo de equipos](catalogos/equipos.md)

# Casos de uso

- [Casos de uso](casos/README.md)
  - [Pruebas con cargas reales](casos/pruebas-reales.md)
  - [Marítimo](casos/maritimo/README.md)
    - [CS-001 · 20' Dry lleno justo por espacio](casos/maritimo/seco/cs-001-20dv-justo-por-espacio.md)
    - [CS-002 · 20' Dry lleno justo por peso](casos/maritimo/seco/cs-002-20dv-justo-por-peso.md)
    - [CS-004 · Tarimas de 2,5 m: solo caben en High Cube](casos/maritimo/seco/cs-004-40hc-tarimas-altas.md)
    - [CS-005 · 45' pallet wide: 3 europallets a lo ancho](casos/maritimo/seco/cs-005-45pw-europallets.md)
    - [CS-006 · Reefer 20': tarima más alta que el interior](casos/maritimo/refrigerado/cs-006-20rf-tarima-sobre-alto-interior.md)
    - [CS-007 · Reefer 40' estándar contra High Cube](casos/maritimo/refrigerado/cs-007-40rf-contra-40rh.md)
    - [CS-008 · Flat rack 40': carga más alta que los postes](casos/maritimo/especial/cs-008-40fr-sobredimensionada.md)
    - [CS-009 · Open top 20': máquina que solo entra por el techo](casos/maritimo/especial/cs-009-20ot-carga-por-techo.md)
    - [CS-010 · Plataforma ISO 40': carga más ancha que el piso](casos/maritimo/especial/cs-010-40pf-carga-ancha.md)
  - [Carretera](casos/carretera/README.md)
    - [CS-011 · Caja seca 53': la misma carga por camino ET y por camino B](casos/carretera/caja-seca/cs-011-53-dos-niveles-y-tipo-de-camino.md)
    - [CS-012 · Caja seca 53': cabe adentro pero no pasa por la puerta](casos/carretera/caja-seca/cs-012-53-pieza-mas-ancha-que-la-puerta.md)
    - [CS-013 · Caja refrigerada 48': se acaba el peso antes que el espacio](casos/carretera/refrigerada/cs-013-48rv-limite-de-peso.md)
    - [CS-014 · Caja refrigerada 53': puerta más baja que el interior](casos/carretera/refrigerada/cs-014-53rv-rack-mas-alto-que-la-puerta.md)
    - [CS-015 · Plataforma 53': atados de tubo de 12 m](casos/carretera/plataforma/cs-015-53fb-tubos-largos.md)
    - [CS-016 · Plataforma plana 40': la ficha no da la altura de la cama](casos/carretera/plataforma/cs-016-40fb-sin-altura-de-cama.md)
    - [CS-017 · Camioneta de 1,5 t: paquetería](casos/carretera/rigidos/cs-017-camioneta-1-5t-paqueteria.md)
    - [CS-018 · Camioneta de 3,5 t: línea blanca](casos/carretera/rigidos/cs-018-camioneta-3-5t-linea-blanca.md)
    - [CS-019 · Rabón (C2): la ruta baja la carga permitida](casos/carretera/rigidos/cs-019-rabon-tipo-de-camino.md)
    - [CS-020 · Tortón (C3): 8 o 10 llantas cambian la carga](casos/carretera/rigidos/cs-020-torton-cemento-llantas.md)
    - [CS-028 · Plataforma 48': carga de 3,2 m de ancho con permiso de la NOM-040](casos/carretera/sobredimensionada/cs-028-plataforma-carga-ancha-con-permiso.md)
    - [CS-029 · Plataforma escalonada 53': dos tramos](casos/carretera/plataforma/cs-029-escalonada-dos-tramos.md)
  - [Restricciones](casos/restricciones/README.md)
    - [CS-021 · Pantallas frágiles: nada encima](casos/restricciones/apilado/cs-021-fragil-no-debajo.md)
    - [CS-022 · Tres entregas en una caja: lo último que se carga sale primero](casos/restricciones/multidestino/cs-022-tres-destinos-lifo.md)
    - [CS-023 · Rollos de acero: el peso no puede quedar en un solo extremo](casos/restricciones/peso-por-eje/cs-023-rollos-de-acero-concentrados.md)
    - [CS-024 · «Este lado arriba»: acostada cabría, parada no](casos/restricciones/orientacion/cs-024-este-lado-arriba.md)
  - [Límites](casos/limites/README.md)
    - [CS-003 · Pieza más pesada que todo el equipo](casos/limites/peso/cs-003-pieza-mas-pesada-que-el-equipo.md)
    - [CS-025 · VGM de un 40' Dry](casos/limites/vgm/cs-025-vgm-40dv.md)
  - [Operación](casos/operacion/README.md)
    - [CS-026 · Caja seca forrada por dentro: el catálogo ya no sirve](casos/operacion/equipo-alterado/cs-026-caja-forrada-por-dentro.md)
    - [CS-027 · Llega una tarima más alta de lo que se declaró](casos/operacion/cambios-en-patio/cs-027-tarima-mas-alta-de-lo-declarado.md)

# Decisiones de arquitectura (ADR)

- [Plantilla](adr/plantilla.md)
- [0001 — Monorepo](adr/0001-monorepo.md)
- [0002 — Go para el backend](adr/0002-go-para-el-backend.md)
- [0003 — React + React Three Fiber](adr/0003-interfaz-react-three-fiber.md)
- [0004 — Wails para escritorio](adr/0004-wails-para-escritorio.md)
- [0005 — Expo para móvil](adr/0005-expo-para-movil.md)
- [0006 — API con la librería estándar](adr/0006-api-con-libreria-estandar.md)
- [0007 — PostgreSQL con pgx y sqlc](adr/0007-postgresql-pgx-sqlc.md)
- [0008 — API-first con OpenAPI](adr/0008-api-first-openapi.md)
- [0009 — Idioma y documentación](adr/0009-idioma-y-documentacion.md)
- [0010 — Docker](adr/0010-docker.md)
- [0011 — Versiones de las tecnologías](adr/0011-versiones.md)
- [0012 — Ramas y versiones](adr/0012-ramas-y-versiones.md)
- [0013 — CI y publicación de versiones](adr/0013-ci-y-release.md)
- [0014 — Sistema de diseño Voxa](adr/0014-sistema-de-diseno.md)
- [0015 — Geometría del contrato del motor](adr/0015-geometria-del-contrato-del-motor.md)
- [0016 — Qué rechaza el motor y qué reporta](adr/0016-que-rechaza-el-motor.md)
- [0017 — Carga sobredimensionada y open top](adr/0017-carga-sobredimensionada-y-open-top.md)
- [0018 — Sin puntos de control](adr/0018-sin-puntos-de-control.md)
- [0019 — Equipos personalizados](adr/0019-equipos-personalizados.md)
- [0020 — Casos simulados](adr/0020-casos-simulados.md)
- [0021 — Organización de la documentación](adr/0021-organizacion-de-la-documentacion.md)
- [0022 — Go explicado solo con Go](adr/0022-go-sin-comparaciones.md)
- [0023 — Lectura y validación de la entrada](adr/0023-lectura-y-validacion-de-la-entrada.md)
- [0024 — PostgreSQL sin contraseña en desarrollo](adr/0024-desarrollo-sin-contrasena.md)

# Experiencia de usuario

- [Sistema de diseño Voxa](ux/sistema-diseno.md)

# Aprendizaje

- [Go para este proyecto](aprendizaje/go/00-introduccion.md)
  - [Primeros pasos](aprendizaje/go/01-primeros-pasos.md)
  - [Módulos e imports](aprendizaje/go/01b-modulos-e-imports.md)
  - [Structs y tipos numéricos](aprendizaje/go/02-structs-y-tipos.md)
  - [Métodos y receptores](aprendizaje/go/03-metodos-y-receptores.md)
  - [Errores como valores](aprendizaje/go/04-errores-como-valores.md)

# Bitácora

- [Índice](bitacora/README.md)
  - [Noviembre 2026]()
    - [26 nov · Lectura y validación de la entrada (1.1.4)](bitacora/2026-11-26-validacion-de-la-entrada.md)
    - [26 nov · Orientaciones, colocaciones y plan de carga (1.1.3)](bitacora/2026-11-26-orientaciones-y-plan-de-carga.md)
    - [26 nov · Tipos del motor y Go explicado solo con Go (1.1.2)](bitacora/2026-11-26-tipos-del-motor-y-go-puro.md)
    - [26 nov · Documentación de desarrollo y pgAdmin](bitacora/2026-11-26-documentacion-de-desarrollo.md)
    - [25 nov · Casos simulados y ajustes al flujo](bitacora/2026-11-25-casos-simulados-y-ajustes.md)
    - [23 nov · Catálogo de equipos y correcciones de P-004](bitacora/2026-11-23-catalogo-de-equipos.md)
    - [21 nov · Recorrido de pendientes](bitacora/2026-11-21-recorrido-de-pendientes.md)
    - [21 nov · Contrato del motor (1.1.1)](bitacora/2026-11-21-contrato-del-motor.md)
    - [21 nov · Sistema de diseño Voxa](bitacora/2026-11-21-sistema-de-diseno-voxa.md)
    - [20 nov · Prompts por sesión y skills](bitacora/2026-11-20-prompts-y-skills.md)
    - [20 nov · Ramas, CI y flujo en develop](bitacora/2026-11-20-ramas-ci-y-flujo.md)
    - [20 nov · Docker y entorno de desarrollo](bitacora/2026-11-20-docker-y-entorno.md)
    - [20 nov · Cimientos del proyecto](bitacora/2026-11-20-cimientos.md)
