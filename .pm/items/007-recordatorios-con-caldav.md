+++
id = 7
title = 'Recordatorios con CalDAV'
description = 'Sincronizar fechas límite con calendario via CalDAV'
state = 'todo'
type = 'feat'
created = '2026-08-17T14:21:53-03:00'
updated = '2026-08-28T15:44:38-03:00'
+++

## Contexto y Alcance

Funcionalidad clasificada como no-MVP en `AGENTS.md`.
Permite a desarrolladores que utilizan clientes de calendario estándar (Apple Calendar, Nextcloud, Thunderbird, Fastmail) recibir recordatorios de fechas límite directamente en sus dispositivos.

## Objetivos

1. Soporte de campo `due = "2026-09-15"` o `deadline` en el frontmatter de items.
2. Comando de sincronización `pj sync caldav` que traduzca items con fecha límite a objetos `VTODO` de iCalendar (RFC 5545).
3. Sincronización con un servidor CalDAV estándar (RFC 4791).

## Diseño Técnico

- **Configuración**: Archivo local no versionado `.pm/caldav.toml` (o variables de entorno `PJ_CALDAV_URL`, `PJ_CALDAV_USER`, `PJ_CALDAV_PASS`) para evitar comitear credenciales en el repo git.
- **Mapeo de Campos**:
  - `UID`: `pj-<project>-<item-id>`
  - `SUMMARY`: Título del item
  - `DESCRIPTION`: Descripción del item y path local
  - `STATUS`: Mapeado desde `state` (`NEEDS-ACTION`, `COMPLETED`, `CANCELLED`).
  - `DUE`: Fecha límite del item.

## Seguridad y Privacidad

- Prohibido almacenar contraseñas o tokens en archivos trackeados por git. El hook de gitleaks del proyecto bloqueará cualquier intento de commit con credenciales.