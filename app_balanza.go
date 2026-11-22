package main

import (
	"os"
	"strings"

	"sistema-pesaje/internal/balanza"
	"sistema-pesaje/internal/datos"
)

func (a *App) detenerBalanza() {
	if a.balanzaSvc != nil {
		a.balanzaSvc.Cerrar()
		a.balanzaSvc = nil
	}
}

func (a *App) crearLector(cfg balanza.Config) balanza.LectorPeso {
	if a.db != nil && datos.UsarBalanzaSimulada(cfg) {
		return balanza.NuevaFuenteSimulada()
	}
	if strings.EqualFold(cfg.PuertoCOM, "SIM") {
		return balanza.NuevaFuenteSimulada()
	}
	// En macOS sin puerto, permitir sim con variable de entorno al desarrollar.
	if os.Getenv("PESAJE_BALANZA_SIM") == "1" {
		return balanza.NuevaFuenteSimulada()
	}
	return balanza.NuevoLectorSerie(cfg)
}

func (a *App) iniciarBalanza(cfg balanza.Config) bool {
	a.detenerBalanza()
	lector := a.crearLector(cfg)
	svc := balanza.NuevoServicio(lector, cfg)
	svc.AlPeso(func(peso int) {
		a.muIntegrado.Lock()
		a.integrado.pesoVivo = peso
		a.muIntegrado.Unlock()
	})
	if !svc.Conectar() {
		a.detenerBalanza()
		return false
	}
	a.balanzaSvc = svc
	a.muIntegrado.Lock()
		a.integrado.conectado = true
		a.muIntegrado.Unlock()
	return true
}
