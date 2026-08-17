# AGENTS.md — pj (Project Journal)

> **Para harness (OpenCode/Agy) y colaboradores:** Este archivo es la BIBLIA del proyecto. Leelo antes de escribir una sola línea de código.

## ¿Qué es pj?

CLI + TUI escrito en Go para gestionar proyectos directamente desde el repo de git. Cada proyecto tiene una carpeta `.pm/` con items en markdown + TOML frontmatter. Sin DB, sin server. Un solo binario estático.

**Ejecutable:** `pj`
**Stack:** Go 1.21+ / cobra (CLI) / bubbletea (TUI) / go-toml (datos)
**Estado actual:** MVP en desarrollo (H1 completado)

---

## Scope MVP (APROBADO)

Lo que SÍ vamos a construir:

| Comando | Qué hace |
|---------|----------|
| `pj init` | Crea `.pm/project.toml` y `.pm/items/` en el directorio actual |
| `pj add [titulo] [descripción]` | Crea nuevo item con ID automático (1,2,3...) |
| `pj list` | Output compacto con colores por estado |
| `pj done <id>` | Marca item como done |
| `pj show` | TUI con barra de progreso e items agrupados por estado |
| `pj status` | Health check (detecta IDs duplicados) |
| `pj renum` | Re-numerar IDs por fecha de creación |
| `pj edit <id>` | Abre item en `$EDITOR` |

**Estados de items:** `todo`, `in progress`, `testing`, `done`, `blocked`, `discarded`

**Formato de datos:** Markdown + TOML frontmatter, un archivo por item en `.pm/items/`
```
.pm/
  project.toml
  items/
    001-planificar.toml
    002-codear.toml
```

---

## Lo que NO es MVP (FUTURO, BACKLOG)

- Dependencias entre items (grafo DAG, ciclos, estimación automática)
- Metodologías (waterfall vs agile)
- Recordatorios (webhook, CalDAV)
- Multi-proyecto / monorepo jerárquico
- Blame de git por tarea
- Autocompletado de shell (para v2)
- Documentación centralizada en `.pm/docs/`

Si un harness propone algo de esto, **lo manda a backog**, no se construye.

---

## MemPalace — Memoria a largo plazo

El proyecto tiene un wing en MemPalace (`wing: pj`). Toda decisión, lección y arquitectura vive acá. **Nunca se pierde.**

### Rooms estándar (metodología MemPalace)

| Room | Propósito | Cuándo escribir |
|------|-----------|-----------------|
| `architecture` | Diseño, stack, archivos clave, decisiones de arquitectura | Cuando se define o cambia el diseño |
| `decisions` | Por qué X sobre Y (trade-offs, alternativas rechazadas) | Cuando se toma una decisión con pros/tras |
| `lessons` | Errores, gotchas, "no hacer X de nuevo" | Cuando algo se rompe o se aprende un patrón |
| `diario` | Bitácora de sesión: qué se hizo, qué importa | **Después de cada sesión de trabajo** |

### Protocolo OBLIGATORIO para harness externo (OpenCode/Agy/Codex)

Cada vez que delegues a un harness, incluí esto en el context:

> **PROTOCOLO MEMPALACE:**
> 1. ANTES de codear: ejecutá `mempalace search "<tema>" --wing pj` para aprender patrones existentes
> 2. DESPUÉS de codear: llamá a `mcp__mempalace__mempalace_add_drawer` para registrar qué hiciste y por qué
> 3. Usá room `architecture` para diseño, `decisions` para trade-offs, `lessons` para errores, `diario` para bitácora

### Search antes de responder

Si el usuario pregunta algo que ya se discutió en otra sesión:
```bash
mempalace search "<tema>" --wing pj
```
Si no encontrás nada, decís "no encontré contexto previo" — **nunca inventes**.


---

## Reglas de trabajo

1. **TDD en todo.** Tests antes que código. `go test ./...` siempre verde antes de commit.
2. **Bite-sized commits.** Un commit por feature o cambio lógico.
3. **No reescribir historial.** Commits forward, no rebase ni amend en main.
4. **Nombres en inglés.** Los comandos están en inglés (`pj list`, no `pj lista`).
5. **La carpeta `.pm/` es fuente de verdad.** El agente es reader y writer. Git resuelve merges.
6. **Fail fast.** Si `pj status` detecta conflicto, bloquea hasta que corras `pj renum`.

---

## Estructura del repo

```
pj/
├── AGENTS.md          # Este archivo
├── cmd/
│   └── pm/
│       └── main.go    # Entrypoint (cobra root)
├── internal/
│   ├── model/
│   │   └── item.go    # Struct Item + TOML marshaling
│   ├── store/
│   │   └── store.go   # ProjectStore (leer/escribir .pm/)
│   └── tui/
│       └── tui.go     # Bubbletea TUI
├── go.mod
├── go.sum
└── README.md
```

---

## Cómo verificar que todo anda

```bash
go build -o pj ./cmd/pm   # Compila
go test ./...              # Tests
./pj init                  # Inicializa en el directorio actual
./pj add "Test item"       # Agrega item
./pj list                  # Lista items
./pj done 1                # Marca done
./pj show                  # Abre TUI
./pj status                # Health check
```

---

## Contacto y decisiones

- Dueño del proyecto: Chaos
- Agent principal: Kateto (opina, no escribe código directamente)
- Cuando haya grietas de diseño importantes, se discuten en blog antes de codear
