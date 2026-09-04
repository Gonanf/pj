# pj

> Proyecto de Gonanf — colección personal.
> **Lenguaje principal (GitHub):** Go · **URL:** https://github.com/Gonanf/pj

## Qué es

Este repositorio forma parte de la colección de **Gonanf / Gabriel Solotorevsky** clonada en `/run/media/chaos/terciario/proyectos/pj`.

> **Nota:** README original preservado abajo en la sección "README original".

- **Path absoluto:** `/run/media/chaos/terciario/proyectos/pj`
- **Estado git:** último commit `2026-08-28 chore(ci): support workflow_dispatch tag creation and set default repo in install.sh`
- **Archivos (aprox):** 87
- **Stack detectado:** Go (go.mod)

## Stack

- Go (go.mod)

## Estructura

```
pj/
.goreleaser.yaml
.opencode/
  .opencode/opencode.json
.pm/
  .pm/design.md
  .pm/items
  .pm/project.toml
  .pm/spec-gap-closure.md
.prompts/
  .prompts/item-types.md
  .prompts/pj-edit.md
  .prompts/pj-finish.md
AGENTS.md
README.md
cmd/
  cmd/pm
go.mod
go.sum
install.sh
internal/
```

## Cómo correr

> Instrucciones genéricas según el stack detectado. Ajustar según el repo.

```bash
go run ./...
go build -o app .
```

## Estado

- **Último commit:** `2026-08-28 chore(ci): support workflow_dispatch tag creation and set default repo in install.sh`
- **Clonado en:** `/run/media/chaos/terciario/proyectos/pj`
- **Exclusiones del lote:** Forks, Workmatch, el-hornero-digital, mali/meli, Sherut (no tocados por consigna)

## Docs

- `docs/overview.md` — descripción extendida y guía rápida (generado en este lote)

## README original (preservado)

> Contenido previo de README.md recortado a 2000 chars para referencia:

```markdown
# 📓 pj — Project Journal

CLI + TUI en Go para gestionar tareas y proyectos directamente desde tu repositorio de Git. Sin bases de datos ni servidores: cada proyecto almacena su estado en la carpeta `.pm/` mediante archivos TOML/Markdown versionables.

---

## 📦 Instalación

### Con `go install`

```bash
go install github.com/chaos/pj/cmd/pm@latest
```

*(Asegurate de tener `$GOPATH/bin` o `~/go/bin` en tu `$PATH`)*

### Compilando desde el código fuente

```bash
git clone https://github.com/chaos/pj.git
cd pj
go build -o pj ./cmd/pm
```

---

## 🚀 Inicio rápido

```bash
# 1. Inicializar pj en el repositorio
pj init

# 2. Agregar tareas
pj add "Diseñar API REST" -d "Definir endpoints y modelos"
pj add "Configurar pipeline de CI/CD"

# 3. Listar tareas en la terminal
pj list

# 4. Abrir la interfaz interactiva (TUI)
pj show

# 5. Marcar una tarea como terminada
pj done 1

# 6. Verificar el estado del proyecto
pj status

# 7. Renumerar tareas por fecha de creación
pj renum
```

---

## 💻 Comandos

### `pj init`

Inicializa el directorio `.pm/` con su archivo de configuración `project.toml` y el directorio `items/`.

```bash
$ pj init
Initialized pj in /home/user/repo/.pm
```

### `pj add [titulo] [-d descripción]`

Crea un nuevo item con ID autoincremental y estado inicial `todo`.

```bash
$ pj add "Implementar login con OAuth" -d "Soporte para GitHub y Google"
Added [#1] Implementar login con OAuth
```

### `pj list`

Muestra un listado compacto de todos los items coloreados según su estado.

```bash
$ pj list
[1] Implementar login con OAuth — todo
[2] Diseñar base de datos — in progress
[3] Setup inicial del proyecto — done
```

### `pj show`

Abre la interfaz TUI interactiva basada en Bubbletea con barra de progreso agrupada por estado.

**Controles dentro de la TUI:**
```

---
*README generado/mejorado automáticamente el 2026-09-04 con inspección de repo (opencode/agy pattern: lectura de estructura, lenguaje y entrypoints). No se modificó código, solo documentación.*
*Autor original: Gonanf — https://github.com/Gonanf/pj*
