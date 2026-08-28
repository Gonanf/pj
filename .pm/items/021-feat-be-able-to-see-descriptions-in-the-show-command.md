+++
id = 21
title = 'feat: be able to see descriptions in the show command'
description = 'Display item descriptions in pj show TUI and allow toggling with key d'
state = 'done'
type = 'feat'
created = '2026-08-28T17:11:02-03:00'
updated = '2026-08-28T17:19:00-03:00'
+++

## Contexto y Alcance

Permite visualizar la descripción de cada tarea directamente dentro de la interfaz interactiva `pj show`, sin tener que abrir el editor ni consultar los archivos individuales en `.pm/items/`.

## Implementación Realizada

1. **Renderizado en TUI (`internal/tui/tui.go`)**:
   - Para cualquier tarea con descripción no vacía, se muestra indentada debajo del título con estilo diferenciado (`descriptionStyle`: color atenuado e itálica).
   - Soporte multilínea preservando saltos de línea de la descripción.
2. **Interactividad (Tecla `d`)**:
   - Al presionar `d`, se conmuta la visibilidad de las descripciones (`showDescriptions = !showDescriptions`).
   - El footer de ayuda se actualizó para indicar `d toggle desc`.
3. **Tests (`internal/tui/tui_test.go`)**:
   - `TestRenderShowsDescriptionAndToggles`: comprueba que la descripción esté visible inicialmente, se oculte al pulsar `d`, y reaparezca al volver a pulsar `d`.

