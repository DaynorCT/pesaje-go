package datos

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ResumenInicio cifras del dashboard (paridad InicioPane).
type ResumenInicio struct {
	PesajesTotal         int64
	PesajesHoy           int64
	PesajesEnProceso     int64
	PesajesFinalizadosHoy int64
	VehiculosActivos     int64
	ProductosActivos     int64
	ChoferesActivos      int64
	ProveedoresActivos   int64
	Ultimos7Dias         []PuntoDia
}

// PuntoDia un día para gráficos de inicio.
type PuntoDia struct {
	Fecha    string `json:"fecha"`
	Cantidad int64  `json:"cantidad"`
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
	if err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pesajes WHERE estado = 'ACTIVO' AND estado_balanza = 'EN_PROCESO'`).Scan(&r.PesajesEnProceso); err != nil {
		return r, err
	}
	if err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pesajes WHERE estado_balanza = 'FINALIZADO' AND date(fecha_modificacion) = date('now', 'localtime')`).Scan(&r.PesajesFinalizadosHoy); err != nil {
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
	if err = contar(ctx, db, `SELECT COUNT(*) FROM proveedores WHERE estado = 'ACTIVO'`, &r.ProveedoresActivos); err != nil {
		return r, err
	}
	r.Ultimos7Dias, err = cargarPesajesUltimosDias(ctx, db, 7)
	return r, err
}

func cargarPesajesUltimosDias(ctx context.Context, db *sql.DB, dias int) ([]PuntoDia, error) {
	if dias <= 0 {
		dias = 7
	}
	q := fmt.Sprintf(`
SELECT date(fecha_creacion) AS d, COUNT(*) AS n
FROM pesajes
WHERE date(fecha_creacion) >= date('now', '-%d days', 'localtime')
GROUP BY d
ORDER BY d`, dias-1)
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	porDia := map[string]int64{}
	for rows.Next() {
		var d string
		var n int64
		if err := rows.Scan(&d, &n); err != nil {
			return nil, err
		}
		porDia[d] = n
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	loc := time.Local
	hoy := time.Now().In(loc)
	out := make([]PuntoDia, 0, dias)
	for i := dias - 1; i >= 0; i-- {
		t := hoy.AddDate(0, 0, -i)
		iso := t.Format("2006-01-02")
		out = append(out, PuntoDia{
			Fecha:    t.Format("02/01"),
			Cantidad: porDia[iso],
		})
	}
	return out, nil
}

func contar(ctx context.Context, db *sql.DB, q string, dest *int64) error {
	return db.QueryRowContext(ctx, q).Scan(dest)
}