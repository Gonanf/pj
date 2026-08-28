+++
id = 11
title = 'Documentación centralizada en .pm/docs/'
description = 'Carpeta .pm/docs/ para archivos markdown referenciados por items'
state = 'todo'
type = 'docs'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Originalmente postergado en `AGENTS.md`, pero en la práctica el proyecto ya utiliza archivos de diseño y especificación dentro de `.pm/` (como `.pm/design.md` y `.pm/spec-gap-closure.md`).
Este item formaliza y estructura la subcarpeta `.pm/docs/`.

## Objetivos

1. Establecer una convención para documentos de arquitectura, RFDs (Request for Discussion) y manuales técnicos asociados a tareas del backlog.
2. Soportar enlaces relativos en el frontmatter de items:
   ```toml
docs = ["architecture.md", "rfcs/auth.md"]
```
3. Proveer subcomandos CLI para interactuar con la documentación:
   - `pj doc list`: Enumera los documentos en `.pm/docs/`.
   - `pj doc view <name>`: Abre el documento en `$PAGER` (o `cat`).
   - `pj doc edit <name>`: Abre el documento en `$EDITOR`.

## Ventajas

- Mantiene toda la documentación viva junto al código en el repositorio git, visible tanto para humanos como para agentes de IA (Kateto, OpenCode, Agy).
- Evita la dispersión de documentación en wikis externas desactualizadas.