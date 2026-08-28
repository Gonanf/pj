+++
id = 22
title = 'feat: allow to set the type with the git methodology'
description = 'Auto-detect item type from conventional commit prefixes (e.g. feat:, fix(scope):)'
state = 'done'
type = 'feat'
created = '2026-08-28T17:12:01-03:00'
updated = '2026-08-28T17:18:00-03:00'
+++

## Contexto y Alcance

Permite que `pj add` detecte automáticamente el tipo de tarea (`feat`, `fix`, `chore`, `docs`) siguiendo la convención de Git (Conventional Commits) cuando el usuario no especifica explícitamente el flag `-t` / `--type`.

## Implementación Realizada

1. **Extractor de Tipo (`internal/model/item.go`)**:
   - Función `model.DetectType(title string) string`.
   - Soporta patrones como `feat: ...`, `feat(scope): ...`, `fix!: ...`, `fix(scope)!: ...`, etc.
   - Case-insensitive, normalizado a minúsculas y validado contra `model.ValidTypes`.
2. **CLI `pj add` (`cmd/pm/main.go`)**:
   - Si `-t` no fue provisto, infiere el tipo invocando `model.DetectType(title)`.
   - Si se especifica `-t`, dicho valor toma precedencia.
3. **Tests**:
   - `internal/model/type_test.go`: `TestDetectType` con diversos formatos y variantes.
   - `internal/e2e/e2e_test.go`: validación integral de `pj add "feat(cli): ..."` y `pj add "chore: ..."`.


