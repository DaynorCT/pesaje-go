package main

import (
	"strings"

	"sistema-pesaje/internal/config"
	pesync "sistema-pesaje/internal/sync"
)

func (a *App) prepararSync(correo, contrasena string) {
	a.muSync.Lock()
	defer a.muSync.Unlock()
	if a.syncSvc != nil {
		a.syncSvc.Cerrar()
		a.syncSvc = nil
	}
	a.claveWeb = contrasena
	cfg := config.Cargar()
	if strings.TrimSpace(cfg.URL) == "" || a.db == nil {
		return
	}
	a.syncSvc = pesync.Crear(cfg, a.db.Conn(), correo, contrasena)
	if a.syncSvc != nil && cfg.AutoEnviar {
		a.syncSvc.IniciarAutomatico(pesync.IntervaloAutomaticoMinutos)
	}
}

func (a *App) cerrarSync() {
	a.muSync.Lock()
	defer a.muSync.Unlock()
	if a.syncSvc != nil {
		a.syncSvc.Cerrar()
		a.syncSvc = nil
	}
	a.claveWeb = ""
}

func (a *App) syncAlMostrar() {
	go func() {
		a.muSync.Lock()
		svc := a.syncSvc
		a.muSync.Unlock()
		if svc == nil {
			return
		}
		svc.SincronizarAhora()
	}()
}

func (a *App) textoSyncIntegrado() (texto, variante string) {
	texto = "Sin web"
	variante = "proceso"
	cfg := config.Cargar()
	if strings.TrimSpace(cfg.URL) == "" {
		return texto, variante
	}
	a.muSync.Lock()
	svc := a.syncSvc
	a.muSync.Unlock()
	if svc == nil {
		texto = "Sin web"
		return texto, variante
	}
	errores := svc.Errores()
	pendientes := svc.Pendientes()
	if errores > 0 {
		return "Con error", "inactivo"
	}
	if pendientes > 0 {
		return "Pend. " + itoa64(pendientes), "proceso"
	}
	return "Al día", "activo"
}

func itoa64(n int64) string {
	if n < 0 {
		n = -n
	}
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}

func (a *App) syncPendientesErrores() (pendientes, errores int) {
	a.muSync.Lock()
	svc := a.syncSvc
	a.muSync.Unlock()
	if svc == nil {
		return 0, 0
	}
	return int(svc.Pendientes()), int(svc.Errores())
}

func (a *App) syncEstadoInicio() string {
	cfg := config.Cargar()
	if strings.TrimSpace(cfg.URL) == "" {
		return "La sincronización no está configurada en este equipo."
	}
	a.muSync.Lock()
	svc := a.syncSvc
	a.muSync.Unlock()
	if svc == nil {
		return "Inicie sesión para activar la sincronización."
	}
	ult := svc.UltimoResultado()
	if ult.Mensaje != "" {
		return ult.Mensaje
	}
	if svc.Pendientes() == 0 && svc.Errores() == 0 {
		return "Al día con el servidor web."
	}
	return ""
}