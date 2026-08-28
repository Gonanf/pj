+++
id = 10
title = 'Blame de git por tarea'
description = 'Asignar autor a item y vincular líneas de código específicas'
state = 'todo'
type = 'feat'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Marcado como no-MVP en `AGENTS.md` (*"Blame de git por tarea"*).
Permite conectar el historial de control de versiones con las tareas documentadas en `.pm/`, respondiendo a preguntas como *"¿qué código se escribió para resolver la tarea #X?"* o *"¿qué tarea originó este cambio?"*.

## Objetivos

1. **Vincular commits con items**:
   - Estandarizar convención de commit messages: `feat(#10): ...` o `[#10] ...`.
2. **Subcomando `pj blame <id>`**:
   - Inspeccionar el log de git ejecutando internamente `git log --grep="\[#<id>\]" --stat`.
   - Listar commits, autores, fechas y archivos modificados para la tarea.
3. **Subcomando inverso `pj origin <file>:<line>`**:
   - Ejecutar `git blame -L <line>,<line> <file>` y extraer la referencia de tarea en el mensaje de commit.

## Diseño Técnico

- Evitar persistir SHAs de commit en el frontmatter de los items; esto generaría commits circulares y constantes conflictos de merge.
- La vinculación debe ser resuelta de manera dinámica en runtime interactuando con el comando `git` del sistema (`os/exec`).