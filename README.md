# Sistema de Pesaje — Go + Wails

Port del módulo Java `:pesaje` (`../jacobitus-escritorio/pesaje`) a **Go + Wails v2**, con paridad funcional y visual (ver `odd/tasks/pesaje-go.md`).

**Target:** Windows 10 y 11 (`.exe`); desarrollo también en macOS.

## Requisitos

| Herramienta | Versión |
|-------------|---------|
| [Go](https://go.dev/dl/) | 1.21+ |
| [Node.js](https://nodejs.org/) | 18+ |
| [Wails CLI](https://wails.io/docs/gettingstarted/installation/) | v2.17+ |

### Windows 10/11 (build del `.exe`)

1. Instalar Go y Node.
2. `go install github.com/wailsapp/wails/v2/cmd/wails@latest` y agregar `%USERPROFILE%\go\bin` al `PATH`.
3. **WebView2** — en Windows 11 ya está; en Windows 10 instalar [Evergreen WebView2](https://developer.microsoft.com/en-us/microsoft-edge/webview2/).
4. Opcional: [NSIS](https://nsis.sourceforge.io/) si querés instalador (`wails build -nsis`).

En macOS solo hace falta tener Xcode Command Line Tools para compilar el host nativo.

## Desarrollo

Desde la raíz de este repo (`pesaje-go/`):

```bash
export PATH="$HOME/go/bin:$PATH"   # macOS/Linux
wails dev
```

Frontend (Vite) en caliente; bindings Go vía Wails.

### Báscula en desarrollo (macOS sin puerto COM)

```bash
export PESAJE_BALANZA_SIM=1
wails dev
```

O configurá `puerto_com` = `SIM` en `config_balanza`. En Windows usa el puerto real (`COM1`, etc.).

## Build producción

```bash
wails build
```

| SO | Salida |
|----|--------|
| Windows amd64 | `build/bin/sistema-pesaje.exe` |
| macOS arm64 | `build/bin/sistema-pesaje.app` |

Cross-compile desde macOS hacia Windows (requiere toolchain mingw si aplica):

```bash
wails build -platform windows/amd64
```

## Datos y paridad con Java

- Referencia de comportamiento: `../jacobitus-escritorio/odd/tasks/pesaje-java.md`.
- Misma carpeta de datos de usuario: `%APPDATA%\SistemaPesaje` (Windows), `~/Library/Application Support/SistemaPesaje` (macOS), `pesaje.db` y `config.json` compatibles.
- Estilos: variables en `frontend/src/css/tema.css` (desde `Tema.java`); `pesaje.css` y `jacobitus.css` web generados desde los recursos JavaFX (`scripts/javafx-css-to-web.py`).

## Estructura

| Ruta | Rol |
|------|-----|
| `main.go`, `app.go` | Host Wails |
| `frontend/src/` | UI web (TypeScript + CSS) |
| `frontend/src/assets/images/logo-pesaje.jpg` | Logo del producto |
| `odd/tasks/pesaje-go.md` | Spec ODD y tareas T1–T12 |

