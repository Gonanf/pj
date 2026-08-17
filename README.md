# 📓 pj — Project Journal

CLI + TUI en Go para gestionar tareas y proyectos directamente desde tu repositorio de Git. Sin bases de datos ni servidores: cada proyecto almacena su estado en la carpeta `.pm/` mediante archivos TOML/Markdown versionables.

---

## 📦 Instalación

### Con `go install`

```bash
go install github.com/chaos/pj/cmd/pm@latest
```

*(Asegurate de tener `$GOPATH/bin` o `~/go/bin` en tu `$PATH`)*

### Compilando desde el código fuente

```bash
git clone https://github.com/chaos/pj.git
cd pj
go build -o pj ./cmd/pm
```

---

## 🚀 Inicio rápido

```bash
# 1. Inicializar pj en el repositorio
pj init

# 2. Agregar tareas
pj add "Diseñar API REST" -d "Definir endpoints y modelos"
pj add "Configurar pipeline de CI/CD"

# 3. Listar tareas en la terminal
pj list

# 4. Abrir la interfaz interactiva (TUI)
pj show

# 5. Marcar una tarea como terminada
pj done 1

# 6. Verificar el estado del proyecto
pj status

# 7. Renumerar tareas por fecha de creación
pj renum
```

---

## 💻 Comandos

### `pj init`

Inicializa el directorio `.pm/` con su archivo de configuración `project.toml` y el directorio `items/`.

```bash
$ pj init
Initialized pj in /home/user/repo/.pm
```

### `pj add [titulo] [-d descripción]`

Crea un nuevo item con ID autoincremental y estado inicial `todo`.

```bash
$ pj add "Implementar login con OAuth" -d "Soporte para GitHub y Google"
Added [#1] Implementar login con OAuth
```

### `pj list`

Muestra un listado compacto de todos los items coloreados según su estado.

```bash
$ pj list
[1] Implementar login con OAuth — todo
[2] Diseñar base de datos — in progress
[3] Setup inicial del proyecto — done
```

### `pj show`

Abre la interfaz TUI interactiva basada en Bubbletea con barra de progreso agrupada por estado.

**Controles dentro de la TUI:**

- `↑` / `k` o `↓` / `j`: Navegar entre items.
- `Espacio`: Ciclar al siguiente estado.
- `1` a `8`: Asignar estado numéricamente directo.
- `q` / `Esc` / `Ctrl+C`: Salir de la TUI.

### `pj done <id>`

Cambia de forma directa el estado del item especificado a `done`.

```bash
$ pj done 1
[#1] marked as done
```

### `pj status`

Ejecuta un chequeo de salud sobre el proyecto. Detecta colisiones de IDs (por ejemplo, luego de un merge conflict en Git).

```bash
$ pj status
Project healthy
```

### `pj renum`

Reordena y renumera correlativamente todos los IDs de los items según su fecha de creación (`created`).

```bash
$ pj renum
Renumbered 3 items
```

---

## 📊 Estados de las tareas

Cada item transita por los siguientes estados:

| Estado | Descripción |
|---|---|
| `todo` | Tarea pendiente |
| `in progress` | Tarea en desarrollo activo |
| `testing` | En etapa de pruebas o validación |
| `blocked` | Bloqueada por alguna dependencia |
| `done` | Tarea completada |
| `closed` | Tarea cerrada y archivada |
| `in specification` | En definición de requerimientos / diseño |
| `discarded` | Tarea descartada o cancelada |

---

## 🏷️ Tipos de tarea *(Próximamente)*

Se encuentra planificado el soporte para clasificar los items por tipo:

- `feat`: Nueva funcionalidad.
- `chore`: Mantenimiento, dependencias o tareas generales.
- `fix`: Corrección de errores.
- `docs`: Documentación técnica o de usuario.

---

## 📁 Estructura de datos

Todos los datos se guardan en texto plano dentro del repositorio, facilitando revisiones y merges en Git:

```
.pm/
├── project.toml           # Metadatos del proyecto
└── items/
    ├── 001-disenar-api.toml
    └── 002-configurar-ci.toml
```

---

## 🤝 Contribuir

1. Hacé un fork del repositorio.
2. Creá una rama con tu cambio (`git checkout -b feature/nueva-mejora`).
3. Ejecutá los tests para validar: `go test ./...`.
4. Hacé commit de tus cambios (`git commit -m 'feat: agrega soporte para...'`).
5. Hacé push a la rama y abrí un Pull Request.

---

## 📄 Licencia

// TODO
