+++
id = 2
title = 'Autocompletado de shell'
description = 'Generar scripts de autocompletado para fish/zsh con IDs de items'
state = 'todo'
type = 'feat'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T17:10:07-03:00'
+++

## Contexto y Alcance

Este item se encuentra actualmente en el backlog para **v2**, ya que en `AGENTS.md` se definió explícitamente como fuera del scope del MVP.
El objetivo es acelerar drásticamente el flujo de trabajo en la terminal proporcionando autocompletado inteligente tanto para subcomandos como para identificadores de items existentes en `.pm/items/`.

## Objetivos

1. **Subcomandos y Flags**: Completar comandos principales (`init`, `add`, `list`, `done`, `show`, `status`, `renum`, `edit`, `finish`) y sus flags (`-t`, `-d`, `--save`).
2. **Completado dinámico de IDs**: Para comandos que reciben un ID (`pj done <id>`, `pj edit <id>`), autocompletar con los IDs existentes y mostrar el título del item como descripción en la shell:
   - Ejemplo en fish: `pj done 1<TAB>` muestra `14  Integración con Kateto`, `15  Tests de integración end-to-end`.
3. **Shells soportadas**: Zsh, Fish y Bash.

## Diseño Técnico

- Cobra cuenta con generadores integrados (`rootCmd.GenFishCompletion`, `GenZshCompletion`, `GenBashCompletionV2`).
- Implementar un comando `pj completion [bash|zsh|fish]`.
- Registrar `ValidArgsFunction` en comandos con argumento `<id>` para cargar dinámicamente los items de `.pm/items/` y devolver pares `[id]	[title]`.

## Consideraciones

- **Performance**: El completado de shell debe responder en menos de 30-50ms para evitar latencia perceptible en la interacción.
- **Fail-safe**: Si `.pm/` no existe en el directorio de trabajo, retornar completados vacíos de forma limpia sin emitir errores en el prompt.