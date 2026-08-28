+++
id = 6
title = 'Metodología Agile'
description = 'Product backlog, sprints, velocity tracking'
state = 'todo'
type = 'feat'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Descartado para el MVP en `AGENTS.md` y en el dictamen fundacional de agosto 2026: *"Sprints, velocity, burndown = YAGNI"*.
El propósito de este item es documentar los requisitos mínimos en caso de que en futuras etapas se desee incorporar soporte de iteraciones ágiles livianas.

## Objetivos

1. Agrupación temporal de tareas en Sprints o Milestones.
2. Cálculo de métricas básicas de velocidad (cantidad de tareas o puntos cerrados por ciclo).
3. Comandos de ciclo de vida del sprint:
   - `pj sprint new "Sprint 1" --start 2026-09-01 --end 2026-09-15`
   - `pj sprint close`
   - `pj list --sprint current`

## Diseño Técnico

- Guardar metadata de sprints en `.pm/sprints.toml` o en `.pm/sprints/<id>.toml`.
- Los items referencian su sprint mediante `sprint = "sprint-1"` en el frontmatter.
- Resumen en `pj finish` con métricas de cierre del sprint (carry-over vs completados).

## Justificación de Postergación

`pj` busca ser la alternativa minimalista frente a la sobrecarga de herramientas como Jira o Linear. Mantener el núcleo simple y centrado en el backlog directo en git es la máxima prioridad del MVP.