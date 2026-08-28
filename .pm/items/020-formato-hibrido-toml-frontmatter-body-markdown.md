+++
id = 20
title = 'Formato hibrido TOML frontmatter + body markdown'
description = 'Items con cuerpo markdown libre despues del frontmatter. Migracion transparente de .toml legacy. Spec preliminar en .pm/design.md y notas en esta descripcion. Decision del dueno: post-integracion (ya integrado).'
state = 'done'
type = 'feat'
created = '2026-08-22T23:29:35-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Motivación

El diseño original de `pj` contemplaba que los items fueran archivos markdown con frontmatter TOML delimitado por `+++`, permitiendo a humanos y agentes escribir notas técnicas, especificaciones y bitácoras libres tras los metadatos.
En las primeras fases del MVP, los items se almacenaron como archivos `.toml` puros sin cuerpo. Este item restituyó el formato híbrido definitivo.

## Decisiones de Implementación (2026-08-24)

- **Extensión**: Archivos `.md` con frontmatter TOML delimitado por `+++` al inicio y final.
- **Lectura retrocompatible**: `model.UnmarshalItem` y `store.LoadItems` cargan tanto archivos híbridos `.md` como archivos `.toml` legacy (asignando `Body` vacío).
- **Migración transparente al escribir**: Cualquier guardado (`Item.Save`) reescribe en formato `.md` y elimina de inmediato el archivo legacy `.toml` con el mismo ID.
- **Modelo**: Campo `Item.Body string `toml:"-"`` preservado en memoria y en ediciones con `pj edit`.
- **Delimitadores**: Se seleccionó `+++` (estándar TOML frontmatter en Hugo y Zola) para evitar ambigüedades con el delimitador YAML `---`.