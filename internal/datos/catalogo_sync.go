package datos

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

var tablasPendientePush = map[string]bool{
	"productos": true, "origenes": true, "destinos": true,
	"vehiculos": true, "choferes": true, "proveedores": true,
}

// CargarCatalogoParaSync arma el JSON de push para un id local.
func CargarCatalogoParaSync(ctx context.Context, db *sql.DB, recurso string, id int64) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("id inválido")
	}
	recurso = strings.TrimSpace(recurso)
	var cuerpo map[string]any
	var err error
	switch recurso {
	case "productos", "origenes", "destinos":
		cuerpo, err = filaCatalogoSimple(ctx, db, recurso, id)
	case "vehiculos":
		cuerpo, err = filaVehiculo(ctx, db, id)
	case "choferes":
		cuerpo, err = filaChofer(ctx, db, id)
	case "proveedores":
		cuerpo, err = filaProveedor(ctx, db, id)
	default:
		return "", fmt.Errorf("recurso no soportado: %s", recurso)
	}
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(cuerpo)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ListarIDsConPendientePush devuelve ids locales marcados para subir.
func ListarIDsConPendientePush(ctx context.Context, db *sql.DB, tabla string) ([]int64, error) {
	if !tablasPendientePush[tabla] {
		return nil, fmt.Errorf("tabla no soportada: %s", tabla)
	}
	rows, err := db.QueryContext(ctx, fmt.Sprintf(`SELECT id FROM %s WHERE pendiente_push = 1`, tabla))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func filaCatalogoSimple(ctx context.Context, db *sql.DB, tabla string, id int64) (map[string]any, error) {
	if err := validarTablaCatalogo(tabla); err != nil {
		return nil, err
	}
	var nombre, desc, estado string
	err := db.QueryRowContext(ctx, fmt.Sprintf(
		`SELECT nombre, COALESCE(descripcion,''), estado FROM %s WHERE id = ?`, tabla), id).
		Scan(&nombre, &desc, &estado)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"id": id, "nombre": nombre, "descripcion": desc, "estado": estado,
	}, nil
}

func filaVehiculo(ctx context.Context, db *sql.DB, id int64) (map[string]any, error) {
	var placa, tipo, estado string
	var tara int
	err := db.QueryRowContext(ctx, `SELECT placa, COALESCE(tipo_vehiculo,''), tara, estado FROM vehiculos WHERE id = ?`, id).
		Scan(&placa, &tipo, &tara, &estado)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"id": id, "placa": placa, "tipo_vehiculo": tipo, "tara": tara, "estado": estado,
	}, nil
}

func filaChofer(ctx context.Context, db *sql.DB, id int64) (map[string]any, error) {
	var lic, tel, estado string
	var nom, apP, apM, ci sql.NullString
	err := db.QueryRowContext(ctx, `
SELECT c.licencia, c.telefono, c.estado, p.nombre, p.apellido_paterno, p.apellido_materno, p.ci
FROM choferes c JOIN personas p ON p.id = c.persona_id WHERE c.id = ?`, id).
		Scan(&lic, &tel, &estado, &nom, &apP, &apM, &ci)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"id":                id,
		"licencia":          lic,
		"telefono":          tel,
		"estado":            estado,
		"nombre":            nom.String,
		"apellido_paterno":  apP.String,
		"apellido_materno":  apM.String,
		"ci":                ci.String,
	}, nil
}

func filaProveedor(ctx context.Context, db *sql.DB, id int64) (map[string]any, error) {
	var empresa, desc, estado string
	var nom, apP, apM, ci sql.NullString
	err := db.QueryRowContext(ctx, `
SELECT pr.nombre_empresa, COALESCE(pr.descripcion,''), pr.estado, p.nombre, p.apellido_paterno, p.apellido_materno, p.ci
FROM proveedores pr JOIN personas p ON p.id = pr.persona_id WHERE pr.id = ?`, id).
		Scan(&empresa, &desc, &estado, &nom, &apP, &apM, &ci)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"id":                id,
		"nombre_empresa":    empresa,
		"descripcion":       desc,
		"estado":            estado,
		"nombre":            nom.String,
		"apellido_paterno":  apP.String,
		"apellido_materno":  apM.String,
		"ci":                ci.String,
	}, nil
}
