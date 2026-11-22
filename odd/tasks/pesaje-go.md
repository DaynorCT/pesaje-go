# Sistema de Pesaje — Go + Wails (paridad con Java)

Portar el módulo `:pesaje` (Java/JavaFX) a **Go + Wails v2**, manteniendo el mismo producto para Windows 10/11.

## Specs

S1. "quiero q este siusteam lo pasae a go tal como esta sus funcioanldias diseño colores stylo tipografia etc me dejo entender"
S2. Misma funcionalidad que el Java actual documentado en `odd/tasks/pesaje-java.md` y `pesaje/README.md`: login, modo integrado, shell con navegación, catálogos, pesaje con báscula RS232, boletas/reportes PDF, sync con la web, misma base SQLite y `config.json`, mismo flujo (un lanzador, sistema completo con contraseña).
S3. Misma apariencia: paleta `Tema.java` (#9a3412 primario, fondos #f7f8fa/#ffffff, tipografía base 13px System, clases de `pesaje.css` y compartidas de `jacobitus.css`); logo `assets/images/logo-pesaje.jpg`; layout de login, shell, modo integrado alineado a `design/modo-integrado/` y pantallas JavaFX.
S4. El módulo Java `:pesaje` queda en el repo como referencia hasta paridad verificada; el código Go vive en `pesaje-go/`.
S5. Target SO: Windows 10 y 11; entrega `.exe` (y opcional instalador NSIS/Inno); **no** Windows 7.
S6. Reutilizar `%APPDATA%\SistemaPesaje` (y paths equivalentes): misma `pesaje.db`, mismos hashes PBKDF2 y contrato sync con `sistema-pesaje-next`.

## Tasks

| ID | Spec | Tarea | Ruta | Commit |
|----|------|-------|------|--------|
| T1 | S1,S3,S4,S5 | Scaffold Wails, copiar assets/CSS (web), tema CSS variables desde `Tema.java`, README build Windows | pesaje-go/ | `cd970ae` |
| T2 | S2,S6 | SQLite + esquema + seed admin + tests auth (paridad `AuthServiceTest`) | inline | |
| T3 | S2,S3 | Pantalla login (HTML/CSS paridad `LoginPane` + logo) | inline | |
| T4 | S2,S3 | Modo integrado 480×680, always-on-top, barra superior (paridad `ModoIntegradoPane`) | inline | |
| T5 | S2,S3 | Shell + navegación + menú usuario (paridad `ShellPane`/`Navegacion`) | inline | |
| T6 | S2 | Báscula RS232 + parser + tests | inline | |
| T7 | S2 | Módulo pesaje (captura, grilla, quick-add) | inline | |
| T8 | S2 | Catálogos ABM + empresas/logo + usuarios | inline | |
| T9 | S2 | Dashboard + reportes PDF + vista previa | inline | |
| T10 | S2,S6 | Sync motor/outbox (paridad tests Java) | worker | |
| T11 | S5 | Build `sistema-pesaje.exe` + guía instalación Windows | inline | |
| T12 | S1,S2 | Verificación: tests Go + checklist paridad visual por pantalla | inline | |

## Log

L1. "quiero q este siusteam lo pasae a go tal como esta sus funcioanldias diseño colores stylo tipografia etc me dejo entender ??"
L2. Usuario abrió `/Users/jaru/dev/firma/pesaje-go` como carpeta de proyecto en Cursor para continuar la migración (T1).

