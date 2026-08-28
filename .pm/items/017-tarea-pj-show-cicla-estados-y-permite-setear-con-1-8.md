+++
id = 17
title = 'Tarea: pj show cicla estados y permite setear con 1-8'
description = 'El TUI de pj show debe permitir: espacio cicla al siguiente estado, teclas 1-8 setean estado directo. Estados: todo, in progress, testing, blocked, done, closed, in specification, discarded. Barra de progreso cuenta closed.'
state = 'done'
type = 'feat'
created = '2026-08-17T14:39:57-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Motivación

En la primera versión del TUI (`pj show`), la barra espaciadora solo alternaba binariamente entre `todo` y `done`.
Con la expansión del ciclo de vida a 8 estados (`todo`, `in progress`, `testing`, `blocked`, `done`, `closed`, `in specification`, `discarded`), era imperativo permitir transiciones ágiles sin salir de la interfaz gráfica de terminal.

## Implementación Realizada

1. **Ciclado Secuencial (Barra Espaciadora)**:
   - Al pulsar `espacio`, la tarea seleccionada avanza al siguiente estado según la función `model.NextState(current)`.
2. **Atajos Numéricos Directos (`1`-`8`)**:
   - Teclas `1` a `8` asignan directamente el estado correspondiente:
     - `1`: `todo`
     - `2`: `in progress`
     - `3`: `testing`
     - `4`: `blocked`
     - `5`: `done`
     - `6`: `closed`
     - `7`: `in specification`
     - `8`: `discarded`
3. **Barra de Progreso y Métricas**:
   - La barra cuenta tanto items en estado `done` como en `closed` como completados.
   - Los items con estado `discarded` se excluyen del cómputo total para evitar distorsiones en el porcentaje de progreso.
4. **Persistencia Inmediata**:
   - Cada transición persiste instantáneamente en disco llamando a `item.Save(m.dir)`.

## Verificación

- Verificado en `internal/tui/tui_test.go` y en suites E2E.