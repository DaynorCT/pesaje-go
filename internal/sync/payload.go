package sync

import (
	"strings"
)

// PesajeParaEnvio datos mínimos para armar el JSON del servidor.
type PesajeParaEnvio struct {
	ID              int64
	SyncEventID     string
	Guia            string
	VehiculoID      int64
	ChoferID        *int64
	ProveedorID     *int64
	ProductoID      *int64
	OrigenID        *int64
	DestinoID       *int64
	PesoBruto       int
	Tara            int
	CostoBs         int
	FleteBsPendiente int
	EstadoBalanza   string
	FechaCreacion   string
}

// PayloadPesaje JSON sin guía (la asigna el servidor).
func PayloadPesaje(clientEventID string, p PesajeParaEnvio) string {
	cuerpo := map[string]any{
		"client_event_id": clientEventID,
		"vehiculo_id":     p.VehiculoID,
		"chofer_id":       referenciaID(p.ChoferID),
		"proveedor_id":    referenciaID(p.ProveedorID),
		"producto_id":     referenciaID(p.ProductoID),
		"id_origen":       referenciaID(p.OrigenID),
		"id_destino":      referenciaID(p.DestinoID),
		"peso_bruto":      p.PesoBruto,
		"tara":            p.Tara,
		"costo_bs":        p.CostoBs,
		"flete_bs_pendiente": p.FleteBsPendiente,
		"estado_balanza":  p.EstadoBalanza,
		"fecha_creacion":  p.FechaCreacion,
	}
	return escribirJSON(cuerpo)
}

func referenciaID(id *int64) any {
	if id != nil && *id > 0 {
		return *id
	}
	return nil
}

// PayloadSinGuia comprueba que no se filtre la guía al servidor.
func PayloadSinGuia(payload string) bool {
	return !strings.Contains(payload, "guia")
}