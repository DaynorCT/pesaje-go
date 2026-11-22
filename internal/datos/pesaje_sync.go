package datos

import (
	"context"
	"database/sql"
)

// PesajeParaSync fila local para encolar en sync_outbox.
type PesajeParaSync struct {
	ID               int64
	SyncEventID      string
	Guia             string
	VehiculoID       int64
	ChoferID         sql.NullInt64
	ProveedorID      sql.NullInt64
	ProductoID       sql.NullInt64
	OrigenID         sql.NullInt64
	DestinoID        sql.NullInt64
	PesoBruto        int
	Tara             int
	CostoBs          sql.NullInt64
	FleteBsPendiente sql.NullInt64
	EstadoBalanza    string
	FechaCreacion    string
}

// CargarPesajeParaSync obtiene un pesaje por id local.
func CargarPesajeParaSync(ctx context.Context, db *sql.DB, id int64) (PesajeParaSync, error) {
	var p PesajeParaSync
	var guia sql.NullString
	err := db.QueryRowContext(ctx, `
SELECT id, sync_event_id, guia, vehiculo_id, chofer_id, proveedor_id, producto_id, id_origen, id_destino,
       peso_bruto, tara, costo_bs, flete_bs_pendiente, estado_balanza, fecha_creacion
FROM pesajes WHERE id = ?`, id).Scan(
		&p.ID, &p.SyncEventID, &guia, &p.VehiculoID,
		&p.ChoferID, &p.ProveedorID, &p.ProductoID, &p.OrigenID, &p.DestinoID,
		&p.PesoBruto, &p.Tara, &p.CostoBs, &p.FleteBsPendiente, &p.EstadoBalanza, &p.FechaCreacion,
	)
	if err != nil {
		return PesajeParaSync{}, err
	}
	p.Guia = guia.String
	return p, nil
}