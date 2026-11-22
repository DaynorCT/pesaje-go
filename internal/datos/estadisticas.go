package datos

import (
	"context"
	"database/sql"
)

// ResumenInicio cifras del dashboard (paridad InicioPane).
type ResumenInicio struct {
	PesajesTotal      int64
	PesajesHoy        int64
	VehiculosActivos  int64
	ProductosActivos  int64
	ChoferesActivos   int64
	ProveedoresActivos int64
}

func CargarResumenInicio(ctx context.Context, db *sql.DB) (ResumenInicio, error) {
	var r ResumenInicio
	var err error
	if err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pesajes`).Scan(&r.PesajesTotal); err != nil {
		return r, err
	}
	if err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pesajes WHERE date(fecha_creacion) = date('now', 'localtime')`).Scan(&r.PesajesHoy); err != nil {
		return r, err
	}
	if err = contar(ctx, db, `SELECT COUNT(*) FROM vehiculos WHERE estado = 'ACTIVO'`, &r.VehiculosActivos); err != nil {
		return r, err
	}
	if err = contar(ctx, db, `SELECT COUNT(*) FROM productos WHERE estado = 'ACTIVO'`, &r.ProductosActivos); err != nil {
		return r, err
	}
	if err = contar(ctx, db, `SELECT COUNT(*) FROM choferes WHERE estado = 'ACTIVO'`, &r.ChoferesActivos); err != nil {
		return r, err
	}
	err = contar(ctx, db, `SELECT COUNT(*) FROM proveedores WHERE estado = 'ACTIVO'`, &r.ProveedoresActivos)
	return r, err
}

func contar(ctx context.Context, db *sql.DB, q string, dest *int64) error {
	return db.QueryRowContext(ctx, q).Scan(dest)
}