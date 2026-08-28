+++
id = 23
title = 'allow to set descriptions interactively'
description = 'Set item descriptions interactively with $EDITOR in a Git-like workflow'
state = 'done'
type = 'feat'
created = '2026-08-28T17:13:04-03:00'
updated = '2026-08-28T17:20:00-03:00'
+++

## Contexto y Alcance

Permite ingresar y redactar descripciones de items interactivamente utilizando `$EDITOR` (al igual que `git commit`), facilitando la escritura de descripciones multilínea, notas o listas en markdown sin depender únicamente del flag `-d`.

## Implementación Realizada

1. **Flag Interactivo (`-i` / `--interactive`) en `pj add` (`cmd/pm/main.go`)**:
   - `pj add [title] -i`: Abre `$EDITOR` con una plantilla temporal comentada (`# Lines starting with '#' will be ignored`).
   - `pj add -i` (sin título): Abre `$EDITOR` solicitando el título en la primera línea y la descripción en las subsiguientes (idéntico al comportamiento de git commit para el subject y el body).
   - Ignora líneas de comentario que comiencen con `#`, preservando el formato y saltos de línea del usuario.
   - En caso de título con convención git (`feat:`, etc.), se auto-detecta el tipo de tarea automáticamente.
2. **Tests**:
   - `cmd/pm/add_test.go`: `TestParseInteractiveDescription` y `TestParseInteractiveTitleAndDescription`.
   - `internal/e2e/e2e_test.go`: `TestAddInteractiveToEndToEnd` validando `pj add "Title" -i` y `pj add -i` con editores mockeados.

