package main

import (
	"strings"

	"sistema-pesaje/internal/config"
)

// CambiarModoSync alterna entre modo web y offline-first (persistido en config.json).
func (a *App) CambiarModoSync(modo string) AccionDTO {
	m := strings.ToLower(strings.TrimSpace(modo))
	switch m {
	case config.ModoSyncWeb, config.ModoSyncOfflineFirst:
	default:
		return AccionDTO{
			OK: false,
			Aviso: &AvisoDTO{
				Titulo:  "Modo sincronización",
				Mensaje: "Modo no válido. Use web u offline_first.",
				EsError: true,
			},
		}
	}
	cfg := config.Cargar()
	cfg.ModoSync = m
	if err := config.Guardar(cfg); err != nil {
		return AccionDTO{
			OK: false,
			Aviso: &AvisoDTO{Titulo: "Modo sincronización", Mensaje: err.Error(), EsError: true},
		}
	}
	est := a.construirEstado()
	return AccionDTO{OK: true, Estado: &est}
}

func modoSyncActual() string {
	return config.Cargar().ModoSync
}

func offlineFirstActivo() bool {
	return config.Cargar().EsOfflineFirst()
}
