+++
id = 18
title = 'Fix test: TestTUIModel_MarkDoneAndSave espera toggle viejo'
description = 'El test asume espacio toggle (todo→done→todo) pero ahora cicla (todo→in progress→testing→...). Actualizar test para reflejar NextState.'
state = 'done'
type = 'fix'
created = '2026-08-17T14:40:19-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Diagnóstico

Tras implementar el ciclado de estados completo en la TUI (#017), la suite de pruebas unitarias falló en `internal/tui/tui_test.go`:
- El test existente `TestTUIModel_MarkDoneAndSave` asumía el comportamiento binario antiguo (pulsar espacio sobre un item `todo` lo convertía en `done`, y una segunda pulsación lo regresaba a `todo`).
- Al cambiar a `NextState`, pulsar espacio sobre un item en `todo` lo traslada a `in progress`.

## Solución Aplicada

1. Se refactorizó `TestTUIModel_MarkDoneAndSave` para verificar el nuevo comportamiento:
   - Espacio avanza de `todo` a `in progress`.
   - Se verificaron las transiciones a lo largo de todo el arreglo `model.ValidStates`.
2. Se agregaron tests específicos para los atajos numéricos `1`-`8`.
3. Resultado: `go test ./internal/tui/...` volvió a quedar completamente en verde.