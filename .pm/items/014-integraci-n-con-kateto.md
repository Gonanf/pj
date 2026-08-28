+++
id = 14
title = 'Integración con Kateto'
description = 'Kateto lee/escribe .pm/ como fuente de verdad del proyecto'
state = 'todo'
type = 'feat'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto Fundacional

Principio rector de `AGENTS.md` y del juicio fundacional del 2026-08-11:
*"Dos escritores, una fuente: humano y Kateto escriben en .pm/; git resuelve los merges. No hay conflict resolution programado."*

Kateto actúa como el orquestador principal de proyectos del usuario Chaos. Para coordinar el trabajo sin desincronizaciones, `.pm/` debe ser la única Fuente de Verdad (*Single Source of Truth*).

## Objetivos

1. **Lectura y Escritura Autónoma**:
   - Kateto puede inspeccionar el estado del proyecto consultando `.pm/items/` o ejecutando `./pj list` / `./pj status`.
   - Kateto puede agregar nuevas tareas mediante `./pj add "<titulo>" -t <tipo> -d "<descripcion>"`.
   - Kateto enriquece el cuerpo markdown de los items con especificaciones técnicas detalladas y criterios de aceptación antes de delegar a harnesses ejecutores (OpenCode, Agy, Codex).
2. **Concurrencia y Resolución**:
   - Ningún agente mantiene estado en memoria o base de datos externa; todo se lee del disco.
   - Las discrepancias entre humano y agente se resuelven a nivel de commits en git.
   - Si se detectan IDs en conflicto tras un merge, `pj status` bloquea y `pj renum` re-ordena la secuencia cronológicamente.
3. **Protocolo MemPalace**:
   - Todo contexto estratégico se sincroniza con el wing `pj` en MemPalace (rooms: `architecture`, `decisions`, `lessons`, `diario`).

## Reglas de Interacción

- Kateto nunca borra items sin autorización explícita; los items obsoletos se marcan con estado `discarded`.
- Los cambios realizados por Kateto se commitean con mensajes descriptivos bajo la convención establecida.