+++
id = 3
title = 'Dependencias entre items'
description = 'Grafo DAG con detección de ciclos y validación con pm validate'
state = 'todo'
type = 'feat'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Item postergado a backlog según `AGENTS.md` (*"Dependencias entre items: grafo DAG, ciclos, estimación automática"*).
Permite declarar que un item no puede comenzar o completarse hasta que otros items especificados hayan sido resueltos.

## Objetivos

1. Declarar dependencias en el frontmatter del item mediante `depends_on = [1, 2]`.
2. Validar que el grafo de dependencias sea un Grafo Acíclico Dirigido (DAG), previniendo referencias circulares (ej. 3 -> 4 -> 3).
3. Integrar la verificación en `pj status` o mediante un subcomando dedicado `pj validate`.
4. Reflejar items bloqueados en `pj list` y `pj show` (visualización de dependencias pendientes).

## Diseño Técnico

- **Modelo**: Extender `model.Item` con el campo opcional `DependsOn []int `toml:"depends_on,omitempty"``.
- **Algoritmo de Detección de Ciclos**:
  - Implementar orden topológico con el algoritmo de Kahn o búsqueda en profundidad (DFS) con coloración de nodos (blanco: no visitado, gris: en proceso, negro: visitado).
  - Si se detecta un nodo gris durante el recorrido, reconstruir el ciclo exacto para reportar un error amigable:
    `Error: ciclo de dependencias detectado: #3 -> #4 -> #7 -> #3`.
- **Validación de integridad**:
  - Detectar referencias a IDs inexistentes en `.pm/items/`.
  - Advertir sobre dependencias que apunten a items descartados (`discarded`).
- **Integración con TUI (`pj show`)**:
  - Mostrar en detalle: *Bloqueado por: #1, #2*.
  - Advertencia al marcar `done` un item con dependencias abiertas.