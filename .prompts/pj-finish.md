Tarea: implementar `pj finish` en el repo pj (Go + cobra).

## Contexto
Repo: /home/chaos/proyectos/pj/.worktrees/pj-finish (branch feat/pj-finish, ya creado).
Leé AGENTS.md en la raíz del repo: es la biblia del proyecto. Spec completa en .pm/spec-gap-closure.md (tarea T3).

## Qué construir
`pj finish` imprime resumen de cierre del proyecto. IMPORTANTE: es puramente aditivo e informativo — NO cambia ningún estado, NO archiva, NO bloquea. El proyecto sigue usable igual después.

Output:
- Nombre del proyecto (de .pm/project.toml)
- Total de items por estado (counts de los 8 estados)
- % completitud sobre items no-descartados (done+closed sobre total-no-descartados)
- Duración: fecha created más antigua vs hoy (ej "12 días")
- Lista de items abiertos (todo/in progress/blocked/testing/in specification) con su ID y título, marcados "no terminado"

Flag opcional --save: escribe el mismo resumen en .pm/SUMMARY.md (sobrescribe). Sin flag, solo stdout.

## Archivos
- Tocar: cmd/pm/main.go (nuevo finishCmd usando store.LoadItems, read-only), internal/store/store.go SOLO si falta un helper de lectura (preferir LoadItems existente).
- NO tocar internal/model/item.go ni internal/tui/ (branches paralelos los modifican). No modificar NewItem ni Save.

## Tests TDD primero
- internal/store/finish_test.go (o cmd test si corresponde): TestFinishCountsMixedStates, TestFinishSaveWritesSummaryMD, TestFinishDoesNotModifyItems (los archivos TOML quedan byte-idénticos después de finish), TestFinishEmptyProject (0 items → manejar división por cero del %).
- Usar t.TempDir() como los tests existentes.

## Reglas
- TDD: rojo → verde antes de commit. go build -o pj ./cmd/pm && go test ./... siempre verde.
- Commits bite-sized forward. Nombres en inglés.
- Al terminar: git commit y avisar; NO pushear.

PROTOCOLO MEMPALACE (obligatorio):
1. ANTES: mempalace search "pj finish summary" --wing pj
2. DESPUÉS: mcp__mempalace__mempalace_add_drawer (qué hiciste y por qué) + bitácora en diario.
