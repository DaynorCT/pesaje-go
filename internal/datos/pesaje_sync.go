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
	PesoNeto         int
	CostoBs          sql.NullInt64
	FleteBsPendiente sql.NullInt64
	Estado           string
	EstadoBalanza    string
	PesadorID        sql.NullInt64
	FechaCreacion    string
	FechaModificacion string
}

// CargarPesajeParaSync obtiene un pesaje por id local.
func CargarPesajeParaSync(ctx context.Context, db *sql.DB, id int64) (PesajeParaSync, error) {
	var p PesajeParaSync
	var guia sql.NullString
	err := db.QueryRowContext(ctx, `
SELECT id, sync_event_id, guia, vehiculo_id, chofer_id, proveedor_id, producto_id, id_origen, id_destino,
       peso_bruto, tara, peso_neto, costo_bs, flete_bs_pendiente, estado, estado_balanza, pesador_id,
       fecha_creacion, fecha_modificacion
FROM pesajes WHERE id = ?`, id).Scan(
		&p.ID, &p.SyncEventID, &guia, &p.VehiculoID,
		&p.ChoferID, &p.ProveedorID, &p.ProductoID, &p.OrigenID, &p.DestinoID,
		&p.PesoBruto, &p.Tara, &p.PesoNeto, &p.CostoBs, &p.FleteBsPendiente, &p.Estado, &p.EstadoBalanza, &p.PesadorID,
		&p.FechaCreacion, &p.FechaModificacion,
	)
	if err != nil {
		return PesajeParaSync{}, err
	}
	p.Guia = guia.String
	return p, nil
}