+++
id = 15
title = 'Tests de integración end-to-end'
description = 'Tests que simulan flujo completo de usuario con TUI y CLI'
state = 'done'
type = 'chore'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Motivación

Identificado como Brecha B7 en `.pm/design.md`.
Aunque los paquetes `model`, `store` y `tui` contaban con tests unitarios aislados, no existía una suite integral que validara el flujo de extremo a extremo (E2E) simulando la experiencia de un desarrollador en su terminal.

## Implementación Realizada

Implementado en commit `f4a84e5` (`internal/e2e/e2e_test.go`):

1. **Flujo de CLI (`TestE2E_FullCLIWorkflow`)**:
   - Inicializa un proyecto limpio con `pj init` en `t.TempDir()`.
   - Crea items usando `pj add` con y sin flags `--type` y `--description`.
   - Valida la salida de `pj list`, comprobando colores ANSI y el formato `[id] [tipo] Titulo — estado`.
   - Marca items como completados con `pj done <id>` y verifica actualización de timestamps y estado en disco.
   - Modifica items mediante `pj edit` simulado, verificando renombrado automático de slugs.
2. **Persistencia TUI (`TestE2E_TUIStatePersistence`)**:
   - Simula el ciclo de vida de un modelo Bubbletea `tuiModel`.
   - Envía eventos de teclado (`KeyMsg`): navegación con flechas, ciclado con barra espaciadora (`NextState`), seteo numérico directo con teclas `1`-`8`.
   - Comprueba que cada interacción persiste inmediatamente en el almacenamiento en disco (`item.Save`).
   - Verifica el cálculo dinámico de la barra de progreso.

## Verificación

```bash
go test -v ./internal/e2e/...
```
Todas las pruebas corren de manera determinista y en milisegundos sin requerir emulación de terminal externa.