# Sistema de Pesaje (Go) — Windows 10 y 11

Este documento cubre **cómo generar** `sistema-pesaje.exe` y **cómo instalarlo** en equipos de producción. No se admite Windows 7.

El ejecutable es **autónomo** (Go + WebView2): no hace falta instalar Java ni Node en el cliente.

## Requisitos en el equipo de desarrollo

| Herramienta | Notas |
|-------------|--------|
| [Go](https://go.dev/dl/) 1.21+ | `go version` |
| [Node.js](https://nodejs.org/) 18+ | Solo para compilar el frontend |
| Wails CLI v2.17+ | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` y `%USERPROFILE%\go\bin` en el `PATH` |
| WebView2 (runtime) | En **Windows 11** ya viene; en **Windows 10** instalar [WebView2 Evergreen](https://developer.microsoft.com/microsoft-edge/webview2/) |
| NSIS 3.x (opcional) | Solo si querés instalador (`wails build -nsis` o `scripts\build-windows.ps1 -Nsis`) |

### Compilar en Windows (recomendado)

Clonar o copiar la carpeta `pesaje-go`, abrir **PowerShell** o **cmd** y ejecutar:

```powershell
cd pesaje-go
go install github.com/wailsapp/wails/v2/cmd/wails@latest
.\scripts\build-windows.ps1
```

Salida: `build\bin\sistema-pesaje.exe`

Con instalador (requiere `makensis` en el `PATH`):

```powershell
.\scripts\build-windows.ps1 -Nsis
```

### Compilar en macOS o Linux (cross-build)

Desde la raíz del proyecto:

```bash
export PATH="$HOME/go/bin:$PATH"
./scripts/build-windows.sh
```

Wails compila el binario Windows **sin** toolchain MinGW adicional en la mayoría de entornos recientes. El `.exe` queda en `build/bin/sistema-pesaje.exe`.

Instalador NSIS desde macOS/Linux solo funciona si `makensis` está instalado; si no, use el `.exe` portable.

## Instalación en el cliente (usuario final)

### Opción A — Ejecutable portable

1. Copie `sistema-pesaje.exe` a una carpeta de su elección (por ejemplo `C:\Program Files\SistemaPesaje\` o el escritorio).
2. En **Windows 10**, si al abrir aparece un error de WebView, instale [WebView2 Evergreen](https://developer.microsoft.com/microsoft-edge/webview2/).
3. Ejecute `sistema-pesaje.exe`. En el **primer arranque** la aplicación crea la base local y crea el usuario administrador (`admin.com`) con contraseña inicial `admin123` (también en el aviso de login; cámbiela después).

### Opción B — Instalador NSIS

Si generó el instalador con `-nsis`:

1. Ejecute `sistema-pesaje-amd64-installer.exe` (nombre puede variar según versión).
2. Siga el asistente; el acceso directo apunta al mismo binario empaquetado por Wails.

### Datos de la aplicación

| Elemento | Ubicación en Windows |
|----------|----------------------|
| Base SQLite `pesaje.db` | `%APPDATA%\SistemaPesaje\pesaje.db` |
| `config.json` (URL web, sync, always-on-top) | `%APPDATA%\SistemaPesaje\config.json` |
| PDF de reportes y boletas | `%USERPROFILE%\Documents\SistemaPesaje\` |

Una instalación Go **reutiliza** la misma carpeta que el sistema Java/Lazarus anterior: puede copiar `pesaje.db` y `config.json` existentes sin migración adicional.

Ejemplo de `config.json` (también en `docs/config.ejemplo.json`):

```json
{
  "url": "https://su-servidor-pesaje.ejemplo",
  "autoenviar": true,
  "modo_operacion": "INTEGRADO",
  "alwaysontop": true,
  "timeout_http": 4000,
  "ciclo_catalogos": 5
}
```

Si `url` está vacía, el escritorio trabaja solo en local (sin pill de sincronización activa).

### Báscula RS232

1. Conecte la báscula al puerto COM (por ejemplo `COM3`).
2. En el sistema completo → **Balanza**, configure puerto, baudios y protocolo (o `SIM` para pruebas).
3. La configuración se guarda en la base local (`config_balanza`).

### Desinstalación

- **Portable:** elimine la carpeta del `.exe` y, si lo desea, `%APPDATA%\SistemaPesaje` para borrar datos.
- **NSIS:** use “Agregar o quitar programas” del instalador generado por Wails; los datos en `%APPDATA%` no se borran automáticamente (comportamiento habitual para no perder pesajes).

## Verificación rápida tras instalar

1. Login con el correo configurado en la base (por defecto `admin@sistema.com` y la clave mostrada en el primer inicio).
2. Modo integrado: conectar báscula, capturar peso, opcional “Enviar a la web” si hay `url` en `config.json`.
3. Menú → **Sistema completo** (pide contraseña) → módulo Pesaje y guardar un registro.
4. Comprobar que `sync_outbox` envía cuando hay servidor (icono “Al día” / “Pend. N” en la barra).

## Solución de problemas

| Síntoma | Qué revisar |
|---------|-------------|
| No abre / pantalla en blanco | WebView2 en Windows 10; antivirus bloqueando el `.exe` |
| “No se pudo abrir el puerto COM…” | Puerto correcto, cable, permisos, que ningún otro programa use el COM |
| Sync “Sin web” | `config.json` sin `url` o sesión no iniciada |
| Credenciales web fallan | Mismo usuario/contraseña que el portal; `timeout_http` muy bajo en redes lentas |

## Más información

- Desarrollo diario: `README.md` (`wails dev`).
- Especificación de migración: `odd/tasks/pesaje-go.md`.
- Referencia Java (MSI con jpackage): `../jacobitus-escritorio/pesaje/README-WINDOWS.md`.