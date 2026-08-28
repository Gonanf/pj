+++
id = 8
title = 'Recordatorios con Webhook'
description = 'Disparar webhook configurable al vencer deadline'
state = 'todo'
type = 'feat'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Listado como no-MVP en `AGENTS.md`.
Permite la integración con canales de chat (Discord, Slack, Matrix, Telegram) o sistemas de observabilidad cuando ocurren eventos clave en el proyecto.

## Objetivos

1. Enviar un payload JSON mediante HTTP POST a una URL configurada cuando:
   - Se aproxima o vence la fecha límite (`deadline`) de un item.
   - Un item cambia de estado a `blocked` o `done`.
   - Se ejecuta el comando de cierre `pj finish`.
2. Integración adaptable a GitHub Actions / CI o cron local.

## Diseño Técnico

- **Arquitectura**: Dado que `pj` es una herramienta CLI sin demonio en segundo plano (*"Sin DB, sin server. Un solo binario estático"*), la invocación del webhook no puede depender de un cron interno permanente.
- **Mecanismos de ejecución**:
  1. Trigger explícito en comandos CLI (ej. `pj notify` o flag `--notify` en `pj done`).
  2. Integración en pre-push hook de git o pipeline CI (ej. job programado que ejecute `pj status --check-deadlines`).
- **Payload estándar**:
  ```json
{
  "event": "item.deadline_approaching",
  "project": "pj",
  "item": {
    "id": 8,
    "title": "Recordatorios con Webhook",
    "state": "todo",
    "due": "2026-09-01"
  },
  "timestamp": "2026-08-28T15:45:00Z"
}
```