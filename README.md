# Sistema de Pesaje — Go + Wails

Port del módulo Java `:pesaje` (`../jacobitus-escritorio/pesaje`) a **Go + Wails v2**, con paridad funcional y visual.

**Target:** Windows 10 y 11 (`.exe`); desarrollo también en macOS.

| Documento | Para qué sirve |
|-----------|----------------|
| [`odd/tasks/pesaje-go.md`](odd/tasks/pesaje-go.md) | Spec **ODD**: requisitos S1–S6, tareas T1–T12, bitácora |
| [`docs/T12-VERIFICACION.md`](docs/T12-VERIFICACION.md) | Checklist paridad Java vs Go + mapa de tests |
| [`README-WINDOWS.md`](README-WINDOWS.md) | Build e instalación en clientes Windows |
| [`docs/balanza-pruebas.ejemplo.json`](docs/balanza-pruebas.ejemplo.json) | Báscula simulada (`SIM`) sin hardware |

---

## Arquitectura

Aplicación **escritorio** en tres capas. El frontend solo habla con Go vía **bindings Wails** (`App` en `main`); la lógica de negocio y SQLite viven en Go.

```
┌─────────────────────────────────────────────────────────────────┐
│  Presentación (TypeScript + CSS)   frontend/src/                │
│  login · modo-integrado · shell · modulos/*.ts                  │
└────────────────────────────┬────────────────────────────────────┘
                             │ Wails Bindings (JSON/RPC)
┌────────────────────────────▼────────────────────────────────────┐
│  Aplicación / adaptadores        package main                    │
│  app.go          sesión, login, DTOs, modo integrado              │
│  app_shell.go    ventanas, navegación, sesión expuesta           │
│  app_pesaje.go   módulo operativo de pesaje                       │
│  app_balanza.go  conectar báscula, elegir lector SIM/COM        │
│  app_abm.go      catálogos y config báscula                       │
│  app_reportes.go dashboard, PDFs                                   │
│  app_sync.go     cola sync hacia la web                           │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│  Dominio + infraestructura       internal/                        │
│  auth            login PBKDF2                                    │
│  datos           SQLite (pesaje.db), ABM, pesajes, outbox          │
│  modelo          validación bruto/tara/neto (kg)                  │
│  balanza         RS232, parser, servicio en vivo, simulador SIM   │
│  reportes        generación PDF (boletas / reportes)               │
│  sync            política + motor HTTP hacia sistema web          │
│  config          config.json (%APPDATA% / Application Support)   │
│  clave           hashes compatibles con Java                      │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│  Runtime externo                                                │
│  Báscula RS232 (COMx)  o  FuenteSimulada (puerto SIM / env)       │
│  Servidor web (opcional, config URL + apikey)                     │
└─────────────────────────────────────────────────────────────────┘
```

### Flujo de peso (báscula)

1. **Config** en SQLite (`config_balanza`) o defaults (`internal/balanza/config.go`).
2. `app_balanza.crearLector` elige **serie real** vs **`FuenteSimulada`** (`SIM`, `PESAJE_BALANZA_SIM=1`).
3. `balanza.Servicio` lee tramas cada ~200 ms → `Parser` → peso en **kg** (enteros, paridad Java).
4. Callback actualiza `integrado.pesoVivo`; la UI captura bruto/tara y `datos.GuardarPesaje` persiste en `pesajes`.

### Pantallas ↔ archivos Go

| Pantalla / módulo UI | Frontend | Bindings Go |
|----------------------|----------|-------------|
| Login | `login.ts` | `app.go` (`Autenticar`, …) |
| Modo integrado | `modo-integrado.ts` | `app.go` + `app_balanza.go` |
| Shell / menú | `shell.ts`, `navegacion.ts` | `app_shell.go` |
| Pesaje | `modulos/pesaje.ts` | `app_pesaje.go` |
| Báscula (ABM) | `modulos/balanza.ts` | `app_abm.go` + `app_balanza.go` |
| Catálogos | `modulos/catalogos.ts`, `abm.ts` | `app_abm.go` |
| Inicio / dashboard | `modulos/inicio.ts` | `app_reportes.go` (métricas) |
| Reportes / boleta | `modulos/reportes.ts` | `app_reportes.go` |
| Sync | (integrado + shell) | `app_sync.go` |

### Datos en disco (paridad Java)

| Archivo | Ubicación típica |
|---------|------------------|
| `pesaje.db` | `%APPDATA%\SistemaPesaje\` (Win) · `~/Library/Application Support/SistemaPesaje/` (mac) |
| `config.json` | Misma carpeta |
| Esquema SQL | `internal/datos/schema.sql` |

### ODD (cómo se gobierna el proyecto)

Este repo usa **Organic Driven Development**: la verdad funcional está en `odd/tasks/pesaje-go.md` (specs **S#**, tareas **T#**, **Log**). El README es la **vista de arquitectura**; T12 es la **evidencia de paridad** con Java.

**Gentle Shell / Pi** (harness del agente en Cursor) **no** forma parte del binario de pesaje: sirve para desarrollar y correr tests; la app en producción es solo `sistema-pesaje.exe` + datos locales.

### Verificación automática

```bash
./scripts/verify.sh    # go test ./... + build frontend
# o:
go test ./... -count=1
```

| Paquete | Tests destacados |
|---------|------------------|
| `internal/auth` | Login, admin sembrado |
| `internal/balanza` | Parser, simulador, servicio + SIM |
| `internal/datos` | Reglas `SIM`, semilla `PESAJE_BALANZA_SIM` |
| `internal/modelo` | Cálculo neto |
| `internal/reportes` | PDF |
| `internal/sync` | Política y motor |
| `package main` (raíz) | E2E báscula simulada + **registro de pesaje** en SQLite |

Detalle y checklist visual: [`docs/T12-VERIFICACION.md`](docs/T12-VERIFICACION.md).

---

## Requisitos

| Herramienta | Versión |
|-------------|---------|
| [Go](https://go.dev/dl/) | 1.21+ |
| [Node.js](https://nodejs.org/) | 18+ |
| [Wails CLI](https://wails.io/docs/gettingstarted/installation/) | v2.17+ |

### Windows 10/11 (build del `.exe`)

Guía completa: [`README-WINDOWS.md`](README-WINDOWS.md).

```bash
./scripts/build-windows.sh          # macOS/Linux → build/bin/sistema-pesaje.exe
# o en PowerShell: .\scripts\build-windows.ps1
```

## Desarrollo

Desde la raíz (`pesaje-go/`):

```bash
export PATH="$HOME/go/bin:$PATH"   # macOS/Linux
wails dev
```

Frontend (Vite) en caliente; bindings Go vía Wails.

### Báscula en desarrollo (sin pesador real)

```bash
export PESAJE_BALANZA_SIM=1
wails dev
```

O `puerto_com` = **`SIM`** en ABM Balanza (`docs/balanza-pruebas.ejemplo.json`). En el **primer** arranque con esa variable, la base nueva puede quedar con `SIM` automáticamente. Pesos en pantalla: **kg** (como Java). En Windows con hardware: puerto COM real (`COM1`, `COM3`, …).

## Build producción

```bash
wails build
```

| SO | Salida |
|----|--------|
| Windows amd64 | `build/bin/sistema-pesaje.exe` |
| macOS arm64 | `build/bin/sistema-pesaje.app` |

Cross-compile desde macOS:

```bash
./scripts/build-windows.sh
```

## Estructura de carpetas

| Ruta | Rol |
|------|-----|
| `main.go` | Entrada Wails, embed `frontend/dist` |
| `app*.go` | Bindings por módulo (ver tabla arriba) |
| `frontend/src/` | UI web (TypeScript + CSS) |
| `internal/*` | Dominio e infraestructura (sin dependencia de Wails) |
| `assets/images/logo-pesaje.jpg` | Logo canónico (UI + PDF) |
| `scripts/` | Build Windows, `verify.sh`, utilidades CSS desde JavaFX |
| `odd/tasks/pesaje-go.md` | Spec ODD migración |
| `docs/` | Verificación T12, ejemplos de config |

## Estilos y paridad visual

- Referencia Java: `../jacobitus-escritorio/odd/tasks/pesaje-java.md`.
- Variables de tema: `frontend/src/css/tema.css` (desde `Tema.java`).
- `pesaje.css` / `jacobitus.css` web: script `scripts/javafx-css-to-web.py`.
