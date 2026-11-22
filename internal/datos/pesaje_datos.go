package datos

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Opcion id + etiqueta para combos (paridad ComboOpcion).
type Opcion struct {
	ID    int64
	Label string
}

// DetalleVehiculo datos derivados al elegir placa.
type DetalleVehiculo struct {
	Tipo string
	Tara int
}

// DetalleChofer licencia al elegir chofer.
type DetalleChofer struct {
	Licencia string
}

// FilaPesaje fila de la grilla del módulo pesaje.
type FilaPesaje struct {
	ID                 int64
	Guia               string
	NombreChofer       string
	PlacaVehiculo      string
	LicenciaChofer     string
	TipoVehiculo       string
	NombreProveedor    string
	NombreProducto     string
	Origen             string
	Destino            string
	CostoBs            sql.NullInt64
	FleteBsPendiente   sql.NullInt64
	FechaCreacion      string
	PesoBruto          sql.NullInt64
	Tara               sql.NullInt64
	PesoNeto           sql.NullInt64
	Estado             string
	EstadoBalanza      string
}

func (f FilaPesaje) Fecha() string { return SoloFecha(f.FechaCreacion) }
func (f FilaPesaje) Hora() string  { return SoloHora(f.FechaCreacion) }

// PesajeEdicion datos completos de un pesaje para cargar el formulario.
type PesajeEdicion struct {
	ID                 int64
	VehiculoID         int64
	ChoferID           *int64
	ProveedorID        *int64
	ProductoID         *int64
	OrigenID           *int64
	DestinoID          *int64
	PesoBruto          int
	Tara               int
	CostoBs            int
	FleteBsPendiente   int
	EstadoBalanza      string
}

// NuevoPesaje datos para insertar un pesaje desde la pantalla operativa.
type NuevoPesaje struct {
	VehiculoID         int64
	ChoferID           *int64
	ProveedorID        *int64
	ProductoID         *int64
	OrigenID           *int64
	DestinoID          *int64
	PesoBruto          int
	Tara               int
	CostoBs            int
	FleteBsPendiente   int
	PesadorID          int64
	EstadoBalanza      string
}

const listarPesajesSQL = `
SELECT p.id, p.guia,
       TRIM(COALESCE(cper.nombre, '') || ' ' || COALESCE(cper.apellido_paterno, '')
            || ' ' || COALESCE(cper.apellido_materno, '')) AS chofer,
       v.placa AS placa, COALESCE(c.licencia, '') AS licencia,
       COALESCE(v.tipo_vehiculo, '') AS tipo,
       COALESCE(NULLIF(TRIM(COALESCE(pper.nombre, '') || ' '
            || COALESCE(pper.apellido_paterno, '') || ' '
            || COALESCE(pper.apellido_materno, '')), ''), pr.nombre_empresa, '') AS proveedor,
       COALESCE(prod.nombre, '') AS producto,
       COALESCE(o.nombre, '') AS origen, COALESCE(d.nombre, '') AS destino,
       p.costo_bs, p.flete_bs_pendiente, p.fecha_creacion,
       p.peso_bruto, p.tara, p.peso_neto, p.estado, p.estado_balanza
FROM pesajes p
LEFT JOIN vehiculos v ON v.id = p.vehiculo_id
LEFT JOIN choferes c ON c.id = p.chofer_id
LEFT JOIN personas cper ON cper.id = c.persona_id
LEFT JOIN proveedores pr ON pr.id = p.proveedor_id
LEFT JOIN personas pper ON pper.id = pr.persona_id
LEFT JOIN productos prod ON prod.id = p.producto_id
LEFT JOIN origenes o ON o.id = p.id_origen
LEFT JOIN destinos d ON d.id = p.id_destino
WHERE p.estado IN ('ACTIVO', 'INACTIVO')
ORDER BY p.id DESC LIMIT ?`

func ListarVehiculosActivos(ctx context.Context, db *sql.DB) ([]Opcion, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT id, placa FROM vehiculos WHERE estado = 'ACTIVO' ORDER BY placa`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOpciones(rows)
}

func ObtenerDetalleVehiculo(ctx context.Context, db *sql.DB, id int64) (DetalleVehiculo, bool, error) {
	var d DetalleVehiculo
	var tipo sql.NullString
	err := db.QueryRowContext(ctx,
		`SELECT COALESCE(tipo_vehiculo, ''), COALESCE(tara, 0) FROM vehiculos WHERE id = ?`, id).
		Scan(&tipo, &d.Tara)
	if err == sql.ErrNoRows {
		return d, false, nil
	}
	if err != nil {
		return d, false, err
	}
	d.Tipo = tipo.String
	return d, true, nil
}

func ListarChoferesActivos(ctx context.Context, db *sql.DB) ([]Opcion, error) {
	rows, err := db.QueryContext(ctx, `
SELECT c.id,
       TRIM(COALESCE(p.nombre, '') || ' ' || COALESCE(p.apellido_paterno, '')
            || ' ' || COALESCE(p.apellido_materno, ''))
FROM choferes c
JOIN personas p ON p.id = c.persona_id
WHERE c.estado = 'ACTIVO'
ORDER BY p.nombre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOpciones(rows)
}

func ObtenerDetalleChofer(ctx context.Context, db *sql.DB, id int64) (DetalleChofer, bool, error) {
	var lic sql.NullString
	err := db.QueryRowContext(ctx, `SELECT COALESCE(licencia, '') FROM choferes WHERE id = ?`, id).
		Scan(&lic)
	if err == sql.ErrNoRows {
		return DetalleChofer{}, false, nil
	}
	if err != nil {
		return DetalleChofer{}, false, err
	}
	return DetalleChofer{Licencia: lic.String}, true, nil
}

func ListarProveedoresActivos(ctx context.Context, db *sql.DB) ([]Opcion, error) {
	rows, err := db.QueryContext(ctx, `
SELECT pr.id,
       COALESCE(NULLIF(TRIM(pr.nombre_empresa), ''),
                TRIM(COALESCE(p.nombre, '') || ' ' || COALESCE(p.apellido_paterno, '')))
FROM proveedores pr
JOIN personas p ON p.id = pr.persona_id
WHERE pr.estado = 'ACTIVO'
ORDER BY pr.nombre_empresa, p.nombre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOpciones(rows)
}

func ListarCatalogoActivos(ctx context.Context, db *sql.DB, tabla string) ([]Opcion, error) {
	tabla = strings.ToLower(strings.TrimSpace(tabla))
	switch tabla {
	case "productos", "origenes", "destinos":
	default:
		return nil, fmt.Errorf("tabla de catálogo no permitida: %s", tabla)
	}
	q := fmt.Sprintf(`SELECT id, nombre FROM %s WHERE estado = 'ACTIVO' ORDER BY nombre`, tabla)
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOpciones(rows)
}

func ListarPesajes(ctx context.Context, db *sql.DB, limite int) ([]FilaPesaje, error) {
	if limite <= 0 {
		limite = 50
	}
	rows, err := db.QueryContext(ctx, listarPesajesSQL, limite)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []FilaPesaje
	for rows.Next() {
		var f FilaPesaje
		var guia sql.NullString
		if err := rows.Scan(
			&f.ID, &guia, &f.NombreChofer, &f.PlacaVehiculo, &f.LicenciaChofer, &f.TipoVehiculo,
			&f.NombreProveedor, &f.NombreProducto, &f.Origen, &f.Destino,
			&f.CostoBs, &f.FleteBsPendiente, &f.FechaCreacion,
			&f.PesoBruto, &f.Tara, &f.PesoNeto, &f.Estado, &f.EstadoBalanza,
		); err != nil {
			return nil, err
		}
		f.Guia = guia.String
		out = append(out, f)
	}
	return out, rows.Err()
}

func ObtenerPesajeEdicion(ctx context.Context, db *sql.DB, id int64) (PesajeEdicion, bool, error) {
	var p PesajeEdicion
	var chofer, prov, prod, orig, dest sql.NullInt64
	var bruto, tara, costo, flete sql.NullInt64
	err := db.QueryRowContext(ctx, `
SELECT vehiculo_id, chofer_id, proveedor_id, producto_id, id_origen, id_destino,
       peso_bruto, tara, costo_bs, flete_bs_pendiente, estado_balanza
FROM pesajes WHERE id = ?`, id).Scan(
		&p.VehiculoID, &chofer, &prov, &prod, &orig, &dest,
		&bruto, &tara, &costo, &flete, &p.EstadoBalanza,
	)
	if err == sql.ErrNoRows {
		return p, false, nil
	}
	if err != nil {
		return p, false, err
	}
	p.ID = id
	p.ChoferID = nullInt64Ptr(chofer)
	p.ProveedorID = nullInt64Ptr(prov)
	p.ProductoID = nullInt64Ptr(prod)
	p.OrigenID = nullInt64Ptr(orig)
	p.DestinoID = nullInt64Ptr(dest)
	if bruto.Valid {
		p.PesoBruto = int(bruto.Int64)
	}
	if tara.Valid {
		p.Tara = int(tara.Int64)
	}
	if costo.Valid {
		p.CostoBs = int(costo.Int64)
	}
	if flete.Valid {
		p.FleteBsPendiente = int(flete.Int64)
	}
	return p, true, nil
}

func nullInt64Ptr(v sql.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}
	id := v.Int64
	return &id
}

func ActualizarPesajeRegistro(ctx context.Context, db *sql.DB, p PesajeEdicion, usuarioID int64) error {
	neto := p.PesoBruto - p.Tara
	_, err := db.ExecContext(ctx, `
UPDATE pesajes SET
  vehiculo_id = ?, chofer_id = ?, proveedor_id = ?, producto_id = ?,
  id_origen = ?, id_destino = ?,
  peso_bruto = ?, tara = ?, peso_neto = ?,
  costo_bs = ?, flete_bs_pendiente = ?,
  usuario_modificacion = ?, fecha_modificacion = ?
WHERE id = ?`,
		p.VehiculoID, nullableID(p.ChoferID), nullableID(p.ProveedorID), nullableID(p.ProductoID),
		nullableID(p.OrigenID), nullableID(p.DestinoID),
		p.PesoBruto, p.Tara, neto, p.CostoBs, p.FleteBsPendiente,
		usuarioID, Ahora(), p.ID,
	)
	return err
}

func GuardarPesaje(ctx context.Context, db *sql.DB, p NuevoPesaje, usuarioID int64) (int64, error) {
	neto := p.PesoBruto - p.Tara
	fecha := Ahora()
	evento := uuid.NewString()
	res, err := db.ExecContext(ctx, `
INSERT INTO pesajes (
  vehiculo_id, chofer_id, proveedor_id, producto_id, id_origen, id_destino,
  peso_bruto, tara, peso_neto, costo_bs, flete_bs_pendiente, pesador_id,
  estado, estado_balanza, sync_event_id,
  usuario_creacion, usuario_modificacion, fecha_creacion, fecha_modificacion
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'ACTIVO', ?, ?, ?, ?, ?, ?)`,
		p.VehiculoID, nullableID(p.ChoferID), nullableID(p.ProveedorID), nullableID(p.ProductoID),
		nullableID(p.OrigenID), nullableID(p.DestinoID),
		p.PesoBruto, p.Tara, neto, p.CostoBs, p.FleteBsPendiente, p.PesadorID,
		p.EstadoBalanza, evento,
		usuarioID, usuarioID, fecha, fecha,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func CambiarEstadoRegistroPesaje(ctx context.Context, db *sql.DB, id int64, estado string, usuarioID int64) error {
	_, err := db.ExecContext(ctx, `
UPDATE pesajes SET estado = ?, usuario_modificacion = ?, fecha_modificacion = ? WHERE id = ?`,
		estado, usuarioID, Ahora(), id)
	return err
}

func FinalizarPesaje(ctx context.Context, db *sql.DB, id, usuarioID int64) error {
	_, err := db.ExecContext(ctx, `
UPDATE pesajes SET estado_balanza = 'FINALIZADO', usuario_modificacion = ?, fecha_modificacion = ?
WHERE id = ? AND estado_balanza = 'EN_PROCESO' AND peso_bruto > 0 AND tara > 0`,
		usuarioID, Ahora(), id)
	return err
}

func CrearVehiculo(ctx context.Context, db *sql.DB, placa, tipo string, tara int) (int64, error) {
	fecha := Ahora()
	res, err := db.ExecContext(ctx, `
INSERT INTO vehiculos (placa, tipo_vehiculo, tara, estado, pendiente_push, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, 'ACTIVO', 1, ?, ?)`, placa, tipo, tara, fecha, fecha)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func CrearChofer(ctx context.Context, db *sql.DB, nombre, apPat, apMat, ci, tel, licencia string) (int64, error) {
	fecha := Ahora()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, `
INSERT INTO personas (nombre, apellido_paterno, apellido_materno, ci, telefono, estado, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, 'ACTIVO', ?, ?)`, nombre, apPat, apMat, ci, tel, fecha, fecha)
	if err != nil {
		return 0, err
	}
	personaID, _ := res.LastInsertId()
	res, err = tx.ExecContext(ctx, `
INSERT INTO choferes (persona_id, licencia, telefono, estado, pendiente_push, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, 'ACTIVO', 1, ?, ?)`, personaID, licencia, tel, fecha, fecha)
	if err != nil {
		return 0, err
	}
	choferID, _ := res.LastInsertId()
	return choferID, tx.Commit()
}

func CrearProveedor(ctx context.Context, db *sql.DB, nombre, apPat, apMat, ci, tel, empresa, descripcion string) (int64, error) {
	fecha := Ahora()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, `
INSERT INTO personas (nombre, apellido_paterno, apellido_materno, ci, telefono, estado, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, 'ACTIVO', ?, ?)`, nombre, apPat, apMat, ci, tel, fecha, fecha)
	if err != nil {
		return 0, err
	}
	personaID, _ := res.LastInsertId()
	res, err = tx.ExecContext(ctx, `
INSERT INTO proveedores (persona_id, nombre_empresa, descripcion, estado, pendiente_push, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, 'ACTIVO', 1, ?, ?)`, personaID, empresa, descripcion, fecha, fecha)
	if err != nil {
		return 0, err
	}
	provID, _ := res.LastInsertId()
	return provID, tx.Commit()
}

func scanOpciones(rows *sql.Rows) ([]Opcion, error) {
	var out []Opcion
	for rows.Next() {
		var o Opcion
		if err := rows.Scan(&o.ID, &o.Label); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

func nullableID(id *int64) any {
	if id == nil {
		return nil
	}
	return *id
}