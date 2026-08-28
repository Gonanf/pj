+++
id = 5
title = 'Metodología Waterfall'
description = 'Fases del proyecto: requerimientos, diseño, implementación, testing, deployment'
state = 'todo'
type = 'feat'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Clasificado como funcionalidad futura (no-MVP) en `AGENTS.md`.
Permite estructurar proyectos que requieren fases estrictas y secuenciales, comunes en proyectos de infraestructura, hardware o desarrollo regulado.

## Objetivos

1. Configurar fases del ciclo de vida del proyecto en `.pm/project.toml`:
   ```toml
methodology = "waterfall"
phases = ["requirements", "design", "implementation", "verification", "deployment"]
```
2. Asociar items a una fase específica en su frontmatter (`phase = "design"`).
3. Reglas de validación y gating:
   - Alertar o impedir avanzar una fase si existen items bloqueantes o incompletos en la fase previa.
4. Vistas en CLI y TUI filtradas y ordenadas por fase cronológica.

## Consideraciones de Diseño

- La filosofía central de `pj` es ser ligero y sin fricción. Cualquier soporte de metodologías formales debe ser estrictamente opt-in mediante configuración en `project.toml`.
- Sin esa configuración, el comportamiento por defecto de `pj` se mantiene libre y flexible.