+++
id = 4
title = 'Estimación automática por dependencias'
description = 'Calcular camino crítico y tiempos mínimos por dependencias'
state = 'todo'
type = 'feat'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Extensión analítica de la funcionalidad de dependencias (#003). Marcado como no-MVP en `AGENTS.md`.
Apunta a proveer estimaciones del proyecto calculando la ruta crítica sobre el grafo de dependencias de tareas.

## Objetivos

1. Soportar estimación de esfuerzo/duración por tarea en el frontmatter (ej. `estimate = "2d"`, `estimate = "4h"`).
2. Calcular el Camino Crítico (Critical Path Method - CPM) identificando la secuencia más larga de tareas dependientes que determina la duración mínima del proyecto.
3. Identificar la holgura (*slack*) de las tareas no críticas.
4. Exponer los resultados en el CLI (`pj critical-path` o flag en `pj finish`).

## Diseño Técnico

- **Frontmatter**:
  ```toml
estimate = "1d"
```
- **Cálculo CPM**:
  - Forward pass: calcular Early Start (ES) y Early Finish (EF) para cada nodo del DAG.
  - Backward pass: calcular Late Start (LS) y Late Finish (LF).
  - Holgura (Slack) = $LS - ES$. Los items con $Slack = 0$ conforman la ruta crítica.
- **Salida CLI**:
  - Renderizado tabular mostrando duración total estimada del proyecto y lista ordenada de tareas en la ruta crítica con formato visual destacado.

## Filosofía de Producto

- Evitar convertir `pj` en una herramienta de gestión pesada tipo MS Project.
- Mantener la sintaxis de estimación opcional y legible a simple vista por humanos.