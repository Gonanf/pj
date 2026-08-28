Tarea: agregar tipos de tarea (feat/chore/fix/docs) al repo pj (Go + cobra + bubbletea).

## Contexto
Repo: /home/chaos/proyectos/pj/.worktrees/item-types (branch feat/item-types, ya creado).
Leé AGENTS.md en la raíz del repo: es la biblia del proyecto. Spec completa en .pm/spec-gap-closure.md (tarea T2).

## Qué construir
Nuevo campo opcional Type en Item:

- Valores válidos: 'feat', 'chore', 'fix', 'docs'. Vacío/ausente = válido (backward compatible: los items TOML existentes sin type no rompen).
- `pj add "titulo" -d "desc" -t feat`: flag --type/-t en addCmd. Valor inválido → error listando valores válidos.
- `pj list`: cuando el item tiene tipo, mostrarlo entre corchetes después del ID: `[3] [fix] Corregir crash — todo`. Color fijo por tipo: feat verde (\033[32m), fix rojo (\033[31m), chore amarillo (\033[33m), docs azul (\033[34m).
- TUI (`pj show`): misma columna [tipo] en el render de items.
- Struct Item en internal/model/item.go: agregar campo `Type string \`toml:"type"\`` (omitempty no aplica en toml v2; campo vacío simplemente se omite o queda vacío — verificar round-trip con test).

## Archivos
- Tocar: internal/model/item.go (campo Type), cmd/pm/main.go SOLO las partes de addCmd flag y listCmd render, internal/tui/tui.go (render).
- NO tocar internal/store/store.go (otro branch paralelo lo está modificando para pj edit). Si necesitás persistir, Item.Save ya funciona con el nuevo campo vía marshaling.

## Tests TDD primero
- internal/model/type_test.go: TestTypeRoundTrip (marshal/unmarshal con y sin type), TestInvalidTypeRejected (validador IsValidType).
- internal/tui/tui_test.go: TestRenderShowsTypeBracket, TestRenderWithoutTypeOmitsBracket.
- Usar t.TempDir() como los tests existentes.

## Reglas
- TDD: rojo → verde antes de commit. go build -o pj ./cmd/pm && go test ./... siempre verde.
- Commits bite-sized forward. Nombres en inglés.
- Al terminar: git commit y avisar; NO pushear.

PROTOCOLO MEMPALACE (obligatorio):
1. ANTES: mempalace search "item types" --wing pj
2. DESPUÉS: mcp__mempalace__mempalace_add_drawer (qué hiciste y por qué) + bitácora en diario.
