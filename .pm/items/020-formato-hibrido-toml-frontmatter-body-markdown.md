+++
id = 20
title = 'Formato hibrido TOML frontmatter + body markdown'
description = 'Items con cuerpo markdown libre despues del frontmatter. Migracion transparente de .toml legacy. Spec preliminar en .pm/design.md y notas en esta descripcion. Decision del dueno: post-integracion (ya integrado).'
state = 'done'
created = '2026-08-22T23:29:35-03:00'
updated = '2026-08-24T13:18:48-03:00'
+++

## Decisiones de implementacion (2026-08-24)

- Extension nueva: .md con frontmatter TOML entre delimitadores +++.
- Lectura backward compatible: los .toml legacy siguen cargando (body vacio).
- Migracion transparente al escribir: cualquier Save reescribe en formato
  hibrido .md y elimina el archivo viejo del mismo ID (demostrado: este
  archivo migro solo con `pj done 20`).
- Campo nuevo Item.Body (toml:"-"), preservado por pj edit (el editor abre
  el archivo completo frontmatter+body).
- Grieta abierta: design.md no tenia spec detallada del formato; estas
  decisiones fueron tomadas por el harness como minimo razonable.
