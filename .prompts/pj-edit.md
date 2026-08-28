Tarea: implementar `pj edit <id>` en el repo pj (Go + cobra).

## Contexto
Repo: /home/chaos/proyectos/pj/.worktrees/pj-edit (branch feat/pj-edit, ya creado).
Leé AGENTS.md en la raíz del repo: es la biblia del proyecto. La spec completa está en .pm/spec-gap-closure.md (tarea T1).

## Qué construir
`pj edit <id>` abre el archivo TOML del item en $EDITOR.

- $EDITOR con fallback a `vi`. Soportar editores con flags (ej "code --wait"): parsear como comando + args.
- Al cerrar el editor: re-parsear el TOML. Validar state con model.IsValidState y título no vacío. Inválido → error claro, no escribir nada.
- Si cambia el título: renombrar el archivo con el nuevo slug, manteniendo ID. Borrar el archivo viejo.
- Actualizar campo updated siempre.
- ID inexistente → error.

## Archivos
- Tocar: internal/store/store.go (agregar FindItem(id) / DeleteItem si faltan), cmd/pm/main.go (nuevo cobra command editCmd).
- NO tocar internal/model/item.go ni internal/tui/ (otros branches paralelos los están modificando).
- Tests TDD primero: internal/store/edit_test.go — TestEditRenamesFileOnTitleChange, TestEditRejectsInvalidState, TestEditFailsOnUnknownID, TestEditUpdatesTimestamp. Usar t.TempDir() como hacen los tests existentes de store.

## Reglas
- TDD: test rojo → verde antes de commit. go build -o pj ./cmd/pm && go test ./... siempre verde.
- Commits bite-sized forward (sin rebase/amend). Nombres en inglés.
- Al terminar: git commit y avisar; NO pushear.

PROTOCOLO MEMPALACE (obligatorio):
1. ANTES de codear: mempalace search "pj edit" --wing pj
2. DESPUÉS: mcp__mempalace__mempalace_add_drawer registrando qué hiciste y por qué (room lessons/architecture según corresponda) y bitácora en diario.
