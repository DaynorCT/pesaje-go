package main

import (
	"strings"

	"sistema-pesaje/internal/config"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	ventanaLoginAncho     = 460
	ventanaLoginAlto      = 520
	ventanaLoginMinAncho  = 420
	ventanaLoginMinAlto   = 480
	ventanaLoginMaxAncho  = 520
	ventanaLoginMaxAlto   = 640

	ventanaEscritorioAncho    = 1280
	ventanaEscritorioAlto     = 800
	ventanaEscritorioMinAncho = 800
	ventanaEscritorioMinAlto  = 600
)

func (a *App) aplicarTamanoEscritorio() {
	runtime.WindowSetMaxSize(a.ctx, 0, 0)
	runtime.WindowSetMinSize(a.ctx, ventanaEscritorioMinAncho, ventanaEscritorioMinAlto)
	runtime.WindowSetSize(a.ctx, ventanaEscritorioAncho, ventanaEscritorioAlto)
}

func (a *App) activarVentanaLogin() {
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
	runtime.WindowSetMaxSize(a.ctx, ventanaLoginMaxAncho, ventanaLoginMaxAlto)
	runtime.WindowSetMinSize(a.ctx, ventanaLoginMinAncho, ventanaLoginMinAlto)
	runtime.WindowSetSize(a.ctx, ventanaLoginAncho, ventanaLoginAlto)
}

func (a *App) activarVentanaShell() {
	a.aplicarTamanoEscritorio()
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
}

// activarVentanaIntegrado: mismo tamaño que el sistema completo; mantiene always-on-top de config.
func (a *App) activarVentanaIntegrado() {
	a.aplicarTamanoEscritorio()
	cfg := config.Cargar()
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
