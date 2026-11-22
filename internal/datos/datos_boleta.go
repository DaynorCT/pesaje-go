package datos

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// DatosBoleta datos para PDF (paridad reportes.DatosBoleta Java).
type DatosBoleta struct {
	ID               int64
	Guia             string
	Fecha            string
	Hora             string
	ChoferNombre     string
	ChoferCI         string
	ChoferLicencia   string
	VehiculoPlaca    string
	VehiculoTipo     string
	ProveedorNombre  string
	ProductoNombre   string
	OrigenNombre     string
	DestinoNombre    string
	Operador         string
	PesoBruto        int
	Tara             int
	PesoNeto         int
	CostoBs          int
	FleteBs          int
}

const selectDatosBoleta = `
SELECT p.id, p.guia, p.fecha_creacion,
       TRIM(COALESCE(pe.nombre, '') || ' ' || COALESCE(pe.apellido_paterno, '')
            || ' ' || COALESCE(pe.apellido_materno, '')) AS chofer_nombre,
       COALESCE(pe.ci, '') AS chofer_ci,
       COALESCE(ch.licencia, '') AS chofer_licencia,
       COALESCE(v.placa, '') AS vehiculo_placa,
       COALESCE(v.tipo_vehiculo, '') AS vehiculo_tipo,
       TRIM(COALESCE(pr.nombre, '') || ' ' || COALESCE(pr.apellido_paterno, '')) AS proveedor_nombre,
       COALESCE(prod.nombre, '') AS producto_nombre,
       COALESCE(o.nombre, '') AS origen_nombre,
       COALESCE(d.nombre, '') AS destino_nombre,
       COALESCE(ps.nombre, '') AS operador,
       COALESCE(p.peso_bruto, 0), COALESCE(p.tara, 0), COALESCE(p.peso_neto, 0),
       COALESCE(p.costo_bs, 0), COALESCE(p.flete_bs_pendiente, 0)
FROM pesajes p
LEFT JOIN vehiculos v ON v.id = p.vehiculo_id
LEFT JOIN choferes ch ON ch.id = p.chofer_id
LEFT JOIN personas pe ON pe.id = ch.persona_id
LEFT JOIN proveedores ppr ON ppr.id = p.proveedor_id
LEFT JOIN personas pr ON pr.id = ppr.persona_id
LEFT JOIN productos prod ON prod.id = p.producto_id
LEFT JOIN origenes o ON o.id = p.id_origen
LEFT JOIN destinos d ON d.id = p.id_destino
LEFT JOIN personas ps ON ps.id = p.pesador_id`

func CargarDatosBoleta(ctx context.Context, db *sql.DB, pesajeID int64) (*DatosBoleta, error) {
	q := selectDatosBoleta + ` WHERE p.id = ?`
	return scanDatosBoleta(db.QueryRowContext(ctx, q, pesajeID))
}

func CargarDatosReporte(ctx context.Context, db *sql.DB, desde, hasta, placa string) ([]DatosBoleta, error) {
	q := strings.Builder{}
	q.WriteString(selectDatosBoleta)
	q.WriteString(` WHERE p.estado = 'ACTIVO'`)
	args := []any{}
	if strings.TrimSpace(desde) != "" {
		q.WriteString(` AND date(p.fecha_creacion) >= date(?)`)
		args = append(args, strings.TrimSpace(desde))
	}
	if strings.TrimSpace(hasta) != "" {
		q.WriteString(` AND date(p.fecha_creacion) <= date(?)`)
		args = append(args, strings.TrimSpace(hasta))
	}
	if strings.TrimSpace(placa) != "" {
		q.WriteString(` AND UPPER(v.placa) LIKE ?`)
		args = append(args, "%"+strings.ToUpper(strings.TrimSpace(placa))+"%")
	}
	q.WriteString(` ORDER BY p.id DESC`)
	rows, err := db.QueryContext(ctx, q.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DatosBoleta
	for rows.Next() {
		d, err := scanDatosBoletaRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *d)
	}
	return out, rows.Err()
}

func NombreEmpresaActiva(ctx context.Context, db *sql.DB) string {
	var n sql.NullString
	_ = db.QueryRowContext(ctx,
		`SELECT nombre_empresa FROM empresas WHERE estado = 'ACTIVO' ORDER BY id DESC LIMIT 1`).Scan(&n)
	return n.String
}

func scanDatosBoleta(row *sql.Row) (*DatosBoleta, error) {
	if row == nil {
		return nil, fmt.Errorf("pesaje no encontrado")
	}
	var d DatosBoleta
	var guia sql.NullString
	var fc string
	err := row.Scan(
		&d.ID, &guia, &fc,
		&d.ChoferNombre, &d.ChoferCI, &d.ChoferLicencia,
		&d.VehiculoPlaca, &d.VehiculoTipo, &d.ProveedorNombre,
		&d.ProductoNombre, &d.OrigenNombre, &d.DestinoNombre, &d.Operador,
		&d.PesoBruto, &d.Tara, &d.PesoNeto, &d.CostoBs, &d.FleteBs,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pesaje no encontrado")
	}
	if err != nil {
		return nil, err
	}
	d.Guia = guia.String
	d.Fecha, d.Hora = partesFechaHora(fc)
	return &d, nil
}

func scanDatosBoletaRow(rows *sql.Rows) (*DatosBoleta, error) {
	var d DatosBoleta
	var guia sql.NullString
	var fc string
	err := rows.Scan(
		&d.ID, &guia, &fc,
		&d.ChoferNombre, &d.ChoferCI, &d.ChoferLicencia,
		&d.VehiculoPlaca, &d.VehiculoTipo, &d.ProveedorNombre,
		&d.ProductoNombre, &d.OrigenNombre, &d.DestinoNombre, &d.Operador,
		&d.PesoBruto, &d.Tara, &d.PesoNeto, &d.CostoBs, &d.FleteBs,
	)
	if err != nil {
		return nil, err
	}
	d.Guia = guia.String
	d.Fecha, d.Hora = partesFechaHora(fc)
	return &d, nil
}

func partesFechaHora(fechaCompleta string) (fecha, hora string) {
	if len(fechaCompleta) >= 16 {
		fecha = fechaCompleta[8:10] + "/" + fechaCompleta[5:7] + "/" + fechaCompleta[0:4]
		hora = fechaCompleta[11:16]
	}
	return fecha, hora
}