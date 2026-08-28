+++
id = 16
title = 'Implementar pj edit'
description = 'Abrir item en $EDITOR para modificar título, descripción y estado manualmente'
state = 'done'
type = 'feat'
created = '2026-08-17T14:24:13-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Cierre de la Brecha B1 de `.pm/design.md` y especificación T1 en `.pm/spec-gap-closure.md`.
`pj edit <id>` formaba parte del scope MVP aprobado en `AGENTS.md` pero no había sido implementado en la primera iteración.

## Implementación

1. **Comando `pj edit <id>`** (`cmd/pm/main.go`):
   - Determina el editor preferido leyendo la variable de entorno `$EDITOR` (con fallback a `vi` o `nano`).
   - Maneja comandos de editor con argumentos o flags (ej. `code --wait` o `subl -w`).
   - Abre el archivo del item en el editor interactivo conectando `os.Stdin`, `os.Stdout` y `os.Stderr`.
2. **Validación y Resguardo**:
   - Al salir del editor, re-parsea el archivo con `model.UnmarshalItem`.
   - Valida que el estado pertenezca a `model.ValidStates`. Si es inválido, rechaza el guardado y emite un error explicativo sin alterar el archivo original.
   - Valida que el título no sea vacío.
   - Si el título cambió, recalcula el slug y renombra el archivo `%03d-<nuevo-slug>.md` eliminando el archivo anterior.
   - Actualiza automáticamente el campo `updated` con el timestamp actual.
   - Preserva íntegramente el cuerpo markdown (`Item.Body`).
3. **Capa de Almacenamiento** (`internal/store/store.go`):
   - Métodos `FindItem(id)`, `DeleteItem(id)` y `EditItem(id, edited)`.

## Tests

- Tests en `internal/store/edit_test.go` cubriendo:
  - Renombrado de archivo cuando cambia el título.
  - Rechazo de estados no válidos.
  - Fallo ante IDs inexistentes.
  - Actualización precisa del timestamp `updated`.