# Spec — Gap closure v1

**Fecha:** 2026-08-22
**Estado:** pendiente de OK de Chaos antes de mandar a OpenCode.
**Alcance:** cerrar las brechas B1, B2, B4 y B3 del design.md.
**Fuera de scope explícito:** parser OpenProject/Conquest (descartado por
decisión del dueño), dependencias, recordatorios, autocompletado, monorepo.

---

## T1 — `pj edit <id>`

Cierre del gap: estaba en el scope MVP aprobado y nunca se implementó.

- `pj edit <id>` abre el archivo TOML del item en `$EDITOR` (fallback `vi`).
- Al salir: re-parsear el archivo. Validar que `state` sea un estado válido
  (`model.IsValidState`) y que el título no quede vacío. Si es inválido,
  error claro y sin escribir nada.
- Si cambia el título: renombrar el archivo (nuevo slug), manteniendo ID.
  Actualizar siempre `updated`.
- Edge cases: ID inexistente → error; `$EDITOR` con flags (ej
  `"code --wait"`) → parsear como command + args.

**Tests:** edición de título renombra archivo; estado inválido rechazado;
ID inexistente; updated se actualiza.

## T2 — Tipos de tarea (feat/chore/fix/docs)

Anunciados en README como "Próximamente".

- Nuevo campo opcional `type = 'feat' | 'chore' | 'fix' | 'docs'` en Item.
  Vacío/ausente = válido (items viejos no rompen — backward compatible).
- `pj add "titulo" -t feat` setea el tipo al crear.
- `pj list`: mostrar tipo entre corchetes cuando exista, ej
  `[3] [fix] Corregir crash — todo`, color fijo por tipo (feat verde,
  fix rojo, chore amarillo, docs azul).
- TUI: misma columna en el render de items.
- No romper archivos TOML existentes sin campo type.

**Tests:** round-trip TOML con y sin type; flag -t con valor inválido
rechaza; render de list con y sin tipos.

## T3 — `pj finish`

Registro histórico de cierre. Aclaración del dueño: muchos proyectos nunca
terminan, así que finish NO borra, NO archiva, NO bloquea nada después.
Es puramente aditivo y reversible.

- `pj finish` imprime un resumen de cierre:
  - nombre del proyecto (de `.pm/project.toml`)
  - total de items por estado (counts)
  - % completitud sobre items no-descartados
  - duración: fecha más antigua de `created` vs hoy
  - lista de items que quedaron abiertos (todo/in progress/blocked/testing)
    marcados como "no terminado"
- Flag opcional `--save`: escribe el resumen en `.pm/SUMMARY.md`
  (sobrescribe). Sin flag, solo stdout. El proyecto queda usable igual
  después de finish (se puede seguir agregando items).
- No cambia estados automáticamente. Los items abiertos quedan como están;
  si el dueño quiere cerrarlos, usa `pj done` / la TUI antes o después.

**Tests:** counts correctos con items mixtos; --save crea SUMMARY.md;
finish no modifica ningún item.

## T4 — Verificación real (B3)

Después de T1-T3:

1. Usar pj para gestionar el propio desarrollo: los items de esta spec se
   cargan en `.pm/items/` con `pj add` y se van cerrando con `pj done`.
2. Dogfooding mínimo en otro repo (el que sea, ej portafolio): init + 3 items
   + edit + done. Confirmar que el flujo no traba.

---

## Orden y reglas

- Orden: T1 → T2 → T3 → T4. Un commit bite-sized por sub-feature, TDD.
- Coding: OpenCode (principal), con protocolo MemPalace de AGENTS.md.
- Tests en verde antes de cada commit; commits forward, sin rebase/amend.
- Comandos en inglés, mensajes de error claros.
