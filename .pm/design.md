# pj — Design Doc (discovery make-project)

**Fecha:** 2026-08-22
**Origen:** pipeline `make-project`, fase 0 (discovery) consolidada por Kateto.
**Fuentes:** AGENTS.md, README.md, git log del repo, backlog `.pm/items/` (18 items), blogs `juicio-idea-project-manager.md` (2026-08-11), `brainstorming-project-manager.md` (2026-08-17), `metodologia-soloto-control-tower.md`.

---

## 1. Qué es pj

CLI + TUI en Go para gestionar proyectos directamente desde el repo de git.
Cada proyecto guarda su estado en `.pm/` (markdown + TOML frontmatter, un
archivo por item). Sin DB, sin server, un binario estático.

**Stack:** Go 1.21+ / cobra / bubbletea / go-toml.
**Usuario objetivo:** dev joven que trabaja en terminal, proyectos chicos o
equipos pequeños. Git recomendado pero opcional; CLI opcional (los archivos
se pueden editar a mano). NO para proyectos grandes.

**Propuesta de valor:** eliminar el cambio de contexto (no abrir Jira/Linear/
Notion para actualizar estado). Todo vive en el repo, cero config, cero
fricción. "Project manager para devs que odian los project managers."
Competidor real: no Jira, sino el TODO.md desorganizado y la memoria.

## 2. Decisiones fundacionales (juicio 2026-08-11, aprobado con condiciones)

- **Fuente de verdad por-proyecto:** la carpeta `.pm/`, no el agente ni una app.
  El todo_list de Kateto es scratchpad global cross-project; no se pisan.
- **Dos escritores, una fuente:** humano y Kateto escriben en `.pm/`; git
  resuelve los merges. No hay conflict resolution programado.
- **Formato:** TOML frontmatter legible a mano. Trade-off aceptado.
- **Scope día 1:** backlog + work items + progreso. Sprints, velocity,
  burndown = YAGNI (condición 1 del veredicto).
- **Relación con Conquest/OpenProject:** no reemplaza a Conquest, lo extiende.
  Meta final: reemplazar OpenProject mediante un parser `.pm/` → work packages
  vía Conquest (condición 3).

## 3. Estado actual verificado (2026-08-22)

Verificado por Kateto directamente, no self-report:

- Build OK (`go build -o pj ./cmd/pm`)
- Tests OK: model, store, tui (todos en verde)
- `pj status`: Project healthy
- 13 commits de features + docs + MemPalace methodology
- Comandos implementados: init, add, list, done, show (TUI con ciclado de
  estados y seteo numérico 1-8), status, renum
- 8 estados: todo, in progress, testing, blocked, done, closed,
  in specification, discarded

### Backlog: 18 items, solo 2 done

Los 16 pendientes son casi todo FUTURO/no-MVP según AGENTS.md (dependencias,
metodologías, recordatorios CalDAV/webhook, blame de git, monorepo).
Si alguien propone eso, va a backlog, no se construye.

## 4. Brechas identificadas (discovery)

| # | Brecha | Origen | Prioridad |
|---|--------|--------|-----------|
| B1 | `pj edit <id>` falta pese a estar en el scope MVP aprobado del AGENTS.md | item #16 | Alta |
| B2 | Tipos de tarea (feat/chore/fix/docs) anunciados como "Próximamente" en README | item #19 | Alta |
| B3 | Sin verificación real: el blog cierra con "usar pj en proyectos reales" y hace 5 días que nadie lo usó fuera de este repo | brainstorming §próximo paso | Alta |
| B4 | Promesa "registro histórico al finalizar el proyecto" sin implementar: no existe comando de cierre/resumen (`pj finish` / `pj summary`) | brainstorming resp. 5 | Media |
| B5 | Condición 3 del juicio sin tocar: parser `.pm/` → OpenProject vía Conquest (puente estratégico) | juicio condición 3 | Media (estratégica) |
| B6 | Distribución sin decidir (binarios precompilados vs go install) | item #13 | Media |
| B7 | Tests e2e | item #15 | Baja |

## 5. Camino elegido (OK del usuario, 2026-08-22)

- B1 + B2 (gap de spec) → spec en `.pm/spec-gap-closure.md` (T1, T2).
- B4 (`pj finish`) → incluido en la spec como T3. Aclaración del dueño:
  muchos proyectos nunca terminan; finish es aditivo y reversible, no
  bloquea ni cambia estados.
- **B5 descartada** por decisión del dueño: el parser OpenProject/Conquest
  no tiene sentido. La condición 3 del juicio queda anulada.
- B3 (verificación real) → T4 de la spec: dogfooding en este repo y en otro.
- B6/B7 quedan en backlog sin fecha.

Spec completa: `.pm/spec-gap-closure.md`. Siguiente paso: delegar a OpenCode
cuando Chaos dé el OK final.

## 6. Reglas de trabajo (vigentes, de AGENTS.md)

TDD en todo, bite-sized commits, commits forward (sin rebase/amend), comandos
en inglés, fail fast, protocolo MemPalace obligatorio para harness externo
(search antes de codear, add_drawer después, diario al cerrar).
Coding lo ejecuta OpenCode/Agy/Codex; Kateto orquesta y especifica.

## 7. Próximo paso

Gate de fase 0 presentado al usuario (2026-08-22). No se construye nada hasta
elegir camino (opción 1, 2, 3 u otra).
