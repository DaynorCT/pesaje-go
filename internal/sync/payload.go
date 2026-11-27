package sync

import (
	"strings"
)

// PesajeParaEnvio datos para POST de pesaje finalizado al centralizador web.
type PesajeParaEnvio struct {
	ID                int64
	SyncEventID       string
	Guia              string
	VehiculoID        int64
	ChoferID          *int64
	ProveedorID       *int64
	ProductoID        *int64
	OrigenID          *int64
	DestinoID         *int64
	PesoBruto         int
	Tara              int
	PesoNeto          int
	CostoBs           int
	FleteBsPendiente  int
	Estado            string
	EstadoBalanza     string
	PesadorID         *int64
	FechaCreacion     string
	FechaFinalizacion string
}

// PayloadPesaje JSON del pesaje finalizado para el centralizador (sin ids de cola).
func PayloadPesaje(clientEventID string, p PesajeParaEnvio) string {
	_ = clientEventID
	neto := p.PesoNeto
	if neto <= 0 && p.PesoBruto > p.Tara {
		neto = p.PesoBruto - p.Tara
	}
	estado := strings.TrimSpace(p.Estado)
	if estado == "" {
		estado = "ACTIVO"
	}
	balanza := strings.TrimSpace(p.EstadoBalanza)
	if balanza == "" {
		balanza = "FINALIZADO"
	}
	cuerpo := map[string]any{
		"pesaje_id":           p.ID,
		"vehiculo_id":         p.VehiculoID,
		"chofer_id":           referenciaID(p.ChoferID),
		"proveedor_id":        referenciaID(p.ProveedorID),
		"producto_id":         referenciaID(p.ProductoID),
		"origen_id":           referenciaID(p.OrigenID),
		"destino_id":          referenciaID(p.DestinoID),
		"peso_bruto":          p.PesoBruto,
		"tara":                p.Tara,
		"peso_neto":           neto,
		"costo_bs":            p.CostoBs,
		"flete_bs_pendiente":  p.FleteBsPendiente,
		"estado":              estado,
		"estado_balanza":      balanza,
		"fecha_creacion":      p.FechaCreacion,
		"fecha_finalizacion":  p.FechaFinalizacion,
		"operador_id":         referenciaID(p.PesadorID),
	}
	guia := strings.TrimSpace(p.Guia)
	if guia != "" {
		cuerpo["guia"] = guia
		cuerpo["numero_boleta"] = guia
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
