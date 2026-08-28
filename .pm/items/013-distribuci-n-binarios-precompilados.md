+++
id = 13
title = 'Distribución: binarios precompilados'
description = 'Releases con binarios para Linux/macOS/Windows'
state = 'done'
type = 'chore'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T17:22:00-03:00'
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

## Implementación Realizada (Rol DevOps Automator)

1. **Configuración GoReleaser (`.goreleaser.yaml`)**:
   - Compilación estática sin dependencias CGO (`CGO_ENABLED=0`) con flags `-trimpath -ldflags="-s -w"`.
   - Targets de arquitectura soportados:
     - Linux: `amd64`, `arm64`
     - macOS: `amd64` (Intel), `arm64` (Apple Silicon)
     - Windows: `amd64`
   - Empaquetado en `.tar.gz` (Unix) y `.zip` (Windows) con generación de `checksums.txt` (SHA256).
2. **GitHub Actions CI Pipeline (`.github/workflows/ci.yml`)**:
   - Ejecución en push a `master`/`main` y pull requests.
   - Matriz multi-plataforma: `ubuntu-latest`, `macos-latest`, `windows-latest`.
   - Verificación de dependencias (`go mod verify`), análisis estático (`go vet`), tests con detector de carreras de datos (`go test -v -race`), compilación estática y smoke test del binario generado.
3. **GitHub Actions Release Pipeline (`.github/workflows/release.yml`)**:
   - Automatización en tags `v*` y `workflow_dispatch`.
   - Gate de calidad: ejecuta suite de pruebas antes de publicar.
   - Orquestación con `goreleaser-action` para crear el GitHub Release y adjuntar artefactos con checksums.
4. **Script de Instalación (`install.sh`)**:
   - Script POSIX autónomo con auto-detección de sistema operativo (`uname -s`) y arquitectura (`uname -m`).
   - Descarga el tarball correspondiente y verifica su firma SHA256 contra `checksums.txt`.
   - Instalación segura en `/usr/local/bin` o `$HOME/.local/bin`.