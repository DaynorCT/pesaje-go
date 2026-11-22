# T12 — Verificación (tests + paridad visual)

Referencia Java: `jacobitus-escritorio/pesaje` (`LoginPane`, `ModoIntegradoPane`, `ShellPane`, módulos en `ui/modulos/`).

**Leyenda**

| Símbolo | Significado |
|---------|-------------|
| ✅ | Paridad aceptada (código + prueba manual o test Go) |
| ⚠️ | Parcial — funciona con diferencias documentadas |
| 🔲 | Revisión manual pendiente (comparar lado a lado en Windows) |

**Automatizado** (ejecutar desde la raíz):

```bash
./scripts/verify.sh
```

Última corrida documentada: 2026-11-22 — `go test ./...` OK; paquetes con tests: `auth`, `balanza`, `modelo`, `reportes`, `sync`.

| Paquete | Tests | Cobertura aprox. (stmt) | Paridad Java |
|---------|-------|-------------------------|--------------|
| `internal/auth` | `service_test.go` | ~88% | `AuthServiceTest` |
| `internal/balanza` | `parser_test.go` | ~33% | `ParserPeso` / tramas |
| `internal/modelo` | `calculo_pesaje_test.go` | 100% | `CalculoPesaje` |
| `internal/reportes` | `pdf_test.go` | ~61% | `BoletaPdf` / reportes |
| `internal/sync` | `policy_test.go`, `motor_test.go` | ~16% motor+aplicador vía integración | `SyncPolicyTest`, `MotorSyncTest` |

---

## Global (tema y datos)

| Ítem | Ref. Java | Go/Wails | Estado |
|------|-----------|----------|--------|
| Paleta `#9a3412`, fondos `#f7f8fa` / `#ffffff` | `Tema.java` | `frontend/src/css/tema.css` | ✅ |
| Tipografía 13px system | `pesaje.css` | `--pesaje-font-size: 13px` | ✅ |
| Logo `logo-pesaje.jpg` | recursos JavaFX | login, modo integrado, shell | ✅ |
| `pesaje.db` en `%APPDATA%\SistemaPesaje` | `Database.java` | `internal/datos/paths.go` | ✅ |
| `config.json` mismo contrato | `ConfigSistema.java` | `internal/config` + `docs/config.ejemplo.json` | ✅ |
| PBKDF2 contraseñas | `AuthService` | `internal/auth` + tests | ✅ |

---

## Pantallas

### Login (`LoginPane` → `frontend/src/login.ts`)

| Ítem | Estado |
|------|--------|
| Logo centrado, campos correo/clave, botón entrar | ✅ |
| Aviso primer inicio (admin sembrado) | ✅ |
| Ventana ~460×520 antes de entrar | ✅ |
| Mensajes error (correo, clave, inactivo) | ✅ |
| 🔲 Mismos márgenes/espaciado pixel a pixel vs JavaFX | 🔲 |

### Modo integrado (`ModoIntegradoPane` → `modo-integrado.ts`)

| Ítem | Estado |
|------|--------|
| Ventana 480×680, `alwaysontop` desde config | ✅ |
| Barra: logo, balanza, pill sync, sync ahora, menú | ✅ |
| Tarjeta peso vivo, conectar, capturar, enviar web | ✅ |
| Pill sync: Sin web / Al día / Pend. N / Con error | ✅ |
| Menú: sistema completo, balanza (shell), cerrar sesión | ⚠️ “Configurar balanza” abre shell (no mini-diálogo) |
| Sync inicial al mostrar + auto cada 2 min | ✅ |
| 🔲 Comparar con `design/modo-integrado/` artboards | 🔲 |

### Shell (`ShellPane` / `Navegacion` → `shell.ts`)

| Ítem | Estado |
|------|--------|
| 1280×800 tras contraseña sistema completo | ✅ |
| Cabecera logo + menú principal + submenús Catálogo/Config | ✅ |
| Menú usuario: web, volver integrado, cerrar sesión | ✅ |
| 🔲 Ancho columna navegación vs JavaFX | 🔲 |

### Inicio (`InicioPane` → `modulos/inicio.ts`)

| Ítem | Estado |
|------|--------|
| Contadores pesajes, vehículos, productos, etc. | ✅ |
| Sync: configurada, pendientes, errores, mensaje | ✅ |
| Botón sincronizar ahora | ✅ |

### Pesaje (`PesajePane` → `modulos/pesaje.ts`)

| Ítem | Estado |
|------|--------|
| Combos vehículo/chofer/proveedor/producto/origen/destino | ✅ |
| Peso vivo, capturar bruto/tara, neto calculado | ✅ |
| Grilla recientes, estados, finalizar | ✅ |
| Quick-add chofer/vehículo/proveedor | ✅ |
| Boleta → PDF + abrir visor SO | ⚠️ Sin vista previa embebida (Java `PrevisualizadorPdf`) |
| Encolar pesaje en sync al guardar | ✅ |
| 🔲 Layout columna formulario vs grilla | 🔲 |

### ABM catálogos (`PanelAbm` → `modulos/abm.ts` + `catalogos.ts`)

| Módulo | Estado |
|--------|--------|
| Vehículos, choferes, proveedores | ✅ |
| Productos, orígenes, destinos | ✅ |
| Usuarios (roles, contraseña alta) | ✅ |
| Empresas (datos texto) | ⚠️ Campo **logo** en BD; formulario Go sin carga de imagen |
| `pendiente_push` en guardar | ✅ local; push catálogo vía sync pull principalmente |

### Reportes (`ReportesPane` → `modulos/reportes.ts`)

| Ítem | Estado |
|------|--------|
| Filtros fechas / entidades | ✅ |
| Generar PDF + abrir archivo | ✅ |

### Boleta config (`BoletaConfigPane` → `boleta-config.ts`)

| Ítem | Estado |
|------|--------|
| Edición textos encabezado boleta | ✅ |
| PDF boleta landscape datos pesaje | ✅ |
| Logo empresa en PDF | ✅ `assets/images/logo-pesaje.jpg` embebido (paridad `LogoEmpresa.dataUrlDelSistema`) |

### Balanza (`BalanzaPane` → `modulos/balanza.ts`)

| Ítem | Estado |
|------|--------|
| Formulario puerto COM, baudios, paridad, etc. | ✅ |
| Lectura serie + parser (tests) | ✅ |
| Simulada `SIM` / `PESAJE_BALANZA_SIM` | ✅ |

---

## Funcionalidad transversal

| Ítem | Estado |
|------|--------|
| Un lanzador: integrado → contraseña → shell | ✅ |
| Báscula RS232 Windows | 🔲 Probar con hardware real en Win10/11 |
| Sync push/pull/outbox | ✅ tests + manual con servidor |
| Encolar solo pesajes (como Java `encolarPesaje`) | ✅ |
| Build `sistema-pesaje.exe` | ✅ T11 |

---

## Cierre de paridad (S4)

Cuando los ítems 🔲 estén revisados en Windows y los ⚠️ aceptados o corregidos, se puede declarar **paridad verificada** y retirar el módulo Java de referencia según `pesaje-go.md` S4.

**Brechas conocidas a priorizar (opcional post-T12)**

1. Vista previa PDF embebida en pantalla (si se exige 1:1 con Java).
2. Carga de logo en ABM Empresas (columna `empresas.logo`; boleta usa logo del sistema embebido).
4. Afinado pixel-perfect de login / integrado / shell.