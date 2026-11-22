package main

import (
	"strings"

	"sistema-pesaje/internal/config"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) activarVentanaShell() {
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
	runtime.WindowSetMinSize(a.ctx, 800, 600)
	runtime.WindowSetMaxSize(a.ctx, 0, 0)
	runtime.WindowSetSize(a.ctx, 1280, 800)
}

func (a *App) activarVentanaIntegrado() {
	cfg := config.Cargar()
	runtime.WindowSetMinSize(a.ctx, 420, 600)
	runtime.WindowSetMaxSize(a.ctx, 0, 0)
	runtime.WindowSetSize(a.ctx, 480, 680)
	runtime.WindowSetAlwaysOnTop(a.ctx, cfg.AlwaysOnTop)
}

// SesionActual expone el usuario con sesión iniciada al frontend.
func (a *App) SesionActual() *UsuarioDTO {
	if a.sesion == nil {
		return nil
	}
	u := *a.sesion
	return &u
}

// VolverAModoIntegrado restaura tamaño y always-on-top del modo báscula.
func (a *App) VolverAModoIntegrado() {
	a.activarVentanaIntegrado()
}

// AbrirSistemaWeb abre la URL configurada en el navegador.
func (a *App) AbrirSistemaWeb() AccionDTO {
	cfg := config.Cargar()
	url := strings.TrimSpace(cfg.URL)
	if url == "" {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Sistema web", Mensaje: "No hay URL configurada. Revise config.json.", EsError: true}}
	}
	runtime.BrowserOpenURL(a.ctx, url)
	return AccionDTO{OK: true}
}
