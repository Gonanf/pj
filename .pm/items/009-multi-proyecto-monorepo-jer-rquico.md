+++
id = 9
title = 'Multi-proyecto / Monorepo jerárquico'
description = 'Proyecto → apps → tareas, CLI unificado que navega todo'
state = 'todo'
type = 'feat'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Postergado a v2 en `AGENTS.md` (*"Multi-proyecto / monorepo jerárquico"*).
Apunta a repositorios que albergan múltiples aplicaciones o paquetes (ej. `apps/api`, `apps/web`, `packages/shared`) permitiendo gestionar el journal de manera granular por subproyecto o de forma consolidada.

## Objetivos

1. **Resolución jerárquica de `.pm/`**:
   - Similar al comportamiento de `git`, buscar la carpeta `.pm/` en el directorio actual y ascender hacia la raíz si no se encuentra.
2. **Espacios de trabajo (Workspaces)**:
   - Configuración en `.pm/project.toml` en la raíz del monorepo:
     ```toml
[workspace]
members = ["apps/*", "packages/*"]
```
3. **Comandos con ámbito**:
   - `pj list --all`: Muestra todos los items del monorepo con prefijo `[app]`.
   - `pj list --app web`: Filtra exclusivamente los items de esa sub-aplicación.

## Consideraciones

- Riesgo de fragmentación o IDs conflictivos si cada subproyecto mantiene su propia secuencia de numeración (`001`, `002`).
- Mantener la convención de que un proyecto simple siga funcionando exactamente igual con un único `.pm/` sin complejidad adicional.