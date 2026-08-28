+++
id = 12
title = 'Changelog automático'
description = 'Generar release notes a partir de items completados entre versiones'
state = 'todo'
type = 'feat'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Item de backlog para v2.
Permite transformar automáticamente los items completados en un `CHANGELOG.md` estructurado y limpio siguiendo el estándar de *Keep a Changelog* y *Conventional Commits*.

## Objetivos

1. Subcomando `pj changelog`:
   - Agrupa los items en estado `done` y `closed`.
   - Organiza los cambios por tipo (`feat`, `fix`, `chore`, `docs`).
   - Muestra ID, título y descripción corta de cada item.
2. Soporte para rangos de versión:
   - Generar el changelog para una nueva versión (`--tag v0.2.0`).
   - Filtrar items cerrados desde el último tag de git o desde una fecha determinada.
3. Integración con `pj finish --save`.

## Ejemplo de Salida

```markdown
# Changelog

## [0.2.0] - 2026-08-28

### Features
- [#16] Implementar pj edit
- [#17] Tarea: pj show cicla estados y permite setear con 1-8
- [#19] Agregar tipo de tarea (feat, chore, fix, docs)
- [#20] Formato hibrido TOML frontmatter + body markdown

### Bug Fixes
- [#18] Fix test: TestTUIModel_MarkDoneAndSave espera toggle viejo

### Maintenance
- [#15] Tests de integración end-to-end
```