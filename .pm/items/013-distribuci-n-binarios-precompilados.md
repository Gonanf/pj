+++
id = 13
title = 'Distribución: binarios precompilados'
description = 'Releases con binarios para Linux/macOS/Windows'
state = 'todo'
type = 'chore'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Identificado como Brecha B6 en `.pm/design.md` (*"Distribución sin decidir: binarios precompilados vs go install"*).
El objetivo es permitir que cualquier desarrollador pueda instalar y usar `pj` inmediatamente sin requerir una instalación previa del toolchain de Go.

## Objetivos

1. Configurar GoReleaser (`.goreleaser.yaml`) para cross-compilación automática de binarios estáticos (`CGO_ENABLED=0`).
2. Targets de arquitectura:
   - Linux: `amd64`, `arm64`
   - macOS: `amd64` (Intel), `arm64` (Apple Silicon)
   - Windows: `amd64`
3. Workflow de GitHub Actions que compile y publique releases automáticas en cada tag `v*`.
4. Script de instalación one-liner:
   ```bash
curl -fsSL https://raw.githubusercontent.com/chaos/pj/master/install.sh | sh
```
5. Publicación en Homebrew tap o Arch AUR como canales secundarios.

## Verificaciones Requeridas

- Confirmar que los paquetes TUI de Bubbletea y Lipgloss rendericen correctamente en terminales de Windows (Windows Terminal) y no arrojen errores con colores ANSI o secuencias VT100.