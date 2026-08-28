+++
id = 19
title = 'Agregar tipo de tarea (feat, chore, fix, docs)'
description = 'Los items pueden tener un tipo: feat (feature), chore (tarea rutinaria), fix (bugfix), docs (documentación). Mostrar el tipo en pj list y pj show con color/icono distinto.'
state = 'done'
type = 'feat'
created = '2026-08-17T14:42:02-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Cierre de la Brecha B2 anunciada en `README.md` y especificada en `.pm/spec-gap-closure.md` (tarea T2).
Permite clasificar los items según su naturaleza técnica o funcional.

## Implementación Realizada

1. **Modelo (`internal/model/item.go`)**:
   - Se incorporó el campo `Type string `toml:"type,omitempty"`` en `model.Item`.
   - Se añadieron `model.ValidTypes = []string{"feat", "chore", "fix", "docs"}` y el validador `model.IsValidType`.
   - El campo es opcional: valores vacíos se ignoran en el frontmatter garantizando retrocompatibilidad con items existentes.
2. **CLI `pj add`**:
   - Flag `-t` / `--type` para definir el tipo al crear la tarea (ej. `pj add "Nueva feature" -t feat`).
   - Rechazo con mensaje de ayuda si se ingresa un tipo no reconocido.
3. **CLI `pj list`**:
   - Muestra el tipo entre corchetes con colores diferenciados:
     - `feat`: Verde (`\033[32m`)
     - `fix`: Rojo (`\033[31m`)
     - `chore`: Amarillo (`\033[33m`)
     - `docs`: Azul (`\033[34m`)
4. **TUI `pj show`**:
   - Renderizado del badge correspondiente en la lista con estilos de `lipgloss`.

## Pruebas

- Pruebas unitarias de serialización/deserialización TOML con y sin el campo `type` en `internal/model/type_test.go`.
- Pruebas de renderizado en TUI (`internal/tui/tui_test.go`).