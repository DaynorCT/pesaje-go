package datos

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"sistema-pesaje/internal/balanza"
	"sistema-pesaje/internal/clave"
)

// AbmFila fila de listado genérico.
type AbmFila struct {
	ID     int64
	Estado string
	Celdas map[string]string
}

// GuardarConfigBalanza persiste configuración RS232.
func GuardarConfigBalanza(ctx context.Context, db *sql.DB, cfg balanza.Config) error {
	fecha := Ahora()
	metodo := cfg.MetodoLectura.Texto()
	var id int64
	err := db.QueryRowContext(ctx, `SELECT id FROM config_balanza ORDER BY id DESC LIMIT 1`).Scan(&id)
	if err == sql.ErrNoRows {
		_, err = db.ExecContext(ctx, `
INSERT INTO config_balanza (puerto_com, baudrate, databits, paridad, stopbits, flowcontrol, timeout_ms,
 metodo_lectura, posicion_inicio, posicion_longitud, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			cfg.PuertoCOM, cfg.BaudRate, cfg.DataBits, cfg.Paridad, cfg.StopBits, cfg.FlowControl,
			cfg.TimeoutMs, metodo, cfg.PosicionInicio, cfg.PosicionLongitud, fecha)
		return err
	}
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, `
UPDATE config_balanza SET puerto_com=?, baudrate=?, databits=?, paridad=?, stopbits=?, flowcontrol=?,
 timeout_ms=?, metodo_lectura=?, posicion_inicio=?, posicion_longitud=?, fecha_modificacion=? WHERE id=?`,
		cfg.PuertoCOM, cfg.BaudRate, cfg.DataBits, cfg.Paridad, cfg.StopBits, cfg.FlowControl,
		cfg.TimeoutMs, metodo, cfg.PosicionInicio, cfg.PosicionLongitud, fecha, id)
	return err
}

// AbmListar devuelve filas según entidad (paridad paneles Java).
func AbmListar(ctx context.Context, db *sql.DB, entidad string) ([]AbmFila, error) {
	switch entidad {
	case "productos", "origenes", "destinos":
		return listarCatalogo(ctx, db, entidad)
	case "vehiculos":
		return listarVehiculos(ctx, db)
	case "choferes":
		return listarChoferes(ctx, db)
	case "proveedores":
		return listarProveedores(ctx, db)
	case "empresas":
		return listarEmpresas(ctx, db)
	case "usuarios":
		return listarUsuarios(ctx, db)
	default:
		return nil, fmt.Errorf("entidad ABM desconocida: %s", entidad)
	}
}

// AbmGuardar crea o actualiza un registro.
func AbmGuardar(ctx context.Context, db *sql.DB, entidad string, id int64, campos map[string]string) error {
	_, err := AbmGuardarConID(ctx, db, entidad, id, campos)
	return err
}

// AbmGuardarConID devuelve el id del registro creado o actualizado.
func AbmGuardarConID(ctx context.Context, db *sql.DB, entidad string, id int64, campos map[string]string) (int64, error) {
	if campos == nil {
		campos = map[string]string{}
	}
	switch entidad {
	case "productos", "origenes", "destinos":
		return guardarCatalogo(ctx, db, entidad, id, campos)
	case "vehiculos":
		return guardarVehiculo(ctx, db, id, campos)
	case "choferes":
		return guardarChofer(ctx, db, id, campos)
	case "proveedores":
		return guardarProveedorABM(ctx, db, id, campos)
	case "empresas":
		return id, guardarEmpresa(ctx, db, id, campos)
	case "usuarios":
		return id, guardarUsuario(ctx, db, id, campos)
	default:
		return 0, fmt.Errorf("entidad ABM desconocida: %s", entidad)
	}
}

func listarCatalogo(ctx context.Context, db *sql.DB, tabla string) ([]AbmFila, error) {
	if err := validarTablaCatalogo(tabla); err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT id, nombre, COALESCE(descripcion,''), estado FROM %s ORDER BY id DESC`, tabla)
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AbmFila
	for rows.Next() {
		var id int64
		var nombre, desc, estado string
		if err := rows.Scan(&id, &nombre, &desc, &estado); err != nil {
			return nil, err
		}
		out = append(out, AbmFila{
			ID: id, Estado: estado,
			Celdas: map[string]string{"nombre": nombre, "descripcion": desc},
		})
	}
	return out, rows.Err()
}

func guardarCatalogo(ctx context.Context, db *sql.DB, tabla string, id int64, c map[string]string) (int64, error) {
	if err := validarTablaCatalogo(tabla); err != nil {
		return 0, err
	}
	nombre := strings.TrimSpace(c["nombre"])
	if nombre == "" {
		return 0, fmt.Errorf("el campo Nombre es obligatorio")
	}
	desc := strings.TrimSpace(c["descripcion"])
	fecha := Ahora()
	if id <= 0 {
		res, err := db.ExecContext(ctx, fmt.Sprintf(`
INSERT INTO %s (nombre, descripcion, estado, pendiente_push, fecha_creacion, fecha_modificacion)
VALUES (?, ?, 'ACTIVO', 1, ?, ?)`, tabla), nombre, desc, fecha, fecha)
		if err != nil {
			return 0, err
		}
		nid, _ := res.LastInsertId()
		return nid, nil
	}
	estado := strings.TrimSpace(c["estado"])
	if estado == "" {
		estado = "ACTIVO"
	}
	_, err := db.ExecContext(ctx, fmt.Sprintf(`
UPDATE %s SET nombre=?, descripcion=?, estado=?, pendiente_push=1, fecha_modificacion=? WHERE id=?`, tabla),
		nombre, desc, estado, fecha, id)
	return id, err
}

func listarVehiculos(ctx context.Context, db *sql.DB) ([]AbmFila, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, placa, COALESCE(tipo_vehiculo,''), tara, estado FROM vehiculos ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AbmFila
	for rows.Next() {
		var id int64
		var placa, tipo, estado string
		var tara int
		if err := rows.Scan(&id, &placa, &tipo, &tara, &estado); err != nil {
			return nil, err
		}
		out = append(out, AbmFila{
			ID: id, Estado: estado,
			Celdas: map[string]string{"placa": placa, "tara": fmt.Sprintf("%d", tara), "tipo": tipo},
		})
	}
	return out, rows.Err()
}

func guardarVehiculo(ctx context.Context, db *sql.DB, id int64, c map[string]string) (int64, error) {
	placa := strings.TrimSpace(c["placa"])
	if placa == "" {
		return 0, fmt.Errorf("el campo Placa es obligatorio")
	}
	tipo := strings.TrimSpace(c["tipo"])
	tara, err := parsearEntero(c["tara"], "Tara")
	if err != nil {
		return 0, err
	}
	fecha := Ahora()
	if id <= 0 {
		res, err := db.ExecContext(ctx, `
INSERT INTO vehiculos (placa, tipo_vehiculo, tara, estado, pendiente_push, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, 'ACTIVO', 1, ?, ?)`, placa, tipo, tara, fecha, fecha)
		if err != nil {
			return 0, err
		}
		nid, _ := res.LastInsertId()
		return nid, nil
	}
	estado := strings.TrimSpace(c["estado"])
	if estado == "" {
		estado = "ACTIVO"
	}
	_, err = db.ExecContext(ctx, `
UPDATE vehiculos SET placa=?, tipo_vehiculo=?, tara=?, estado=?, pendiente_push=1, fecha_modificacion=? WHERE id=?`,
		placa, tipo, tara, estado, fecha, id)
	return id, err
}

func listarChoferes(ctx context.Context, db *sql.DB) ([]AbmFila, error) {
	rows, err := db.QueryContext(ctx, `
SELECT c.id, p.nombre, p.apellido_paterno, p.apellido_materno, p.ci, c.telefono, c.licencia, c.estado
FROM choferes c JOIN personas p ON p.id = c.persona_id ORDER BY c.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AbmFila
	for rows.Next() {
		var id int64
		var nom, apP, apM, ci, tel, lic, estado sql.NullString
		if err := rows.Scan(&id, &nom, &apP, &apM, &ci, &tel, &lic, &estado); err != nil {
			return nil, err
		}
		out = append(out, AbmFila{
			ID: id, Estado: estado.String,
			Celdas: map[string]string{
				"nombre":            nom.String,
				"apellido_paterno":  apP.String,
				"apellido_materno":  apM.String,
				"ci":                ci.String,
				"telefono":          tel.String,
				"licencia":          lic.String,
				"nombre_completo":   strings.TrimSpace(nom.String + " " + apP.String + " " + apM.String),
			},
		})
	}
	return out, rows.Err()
}

func guardarChofer(ctx context.Context, db *sql.DB, id int64, c map[string]string) (int64, error) {
	nombre := strings.TrimSpace(c["nombre"])
	lic := strings.TrimSpace(c["licencia"])
	if nombre == "" {
		return 0, fmt.Errorf("el campo Nombre es obligatorio")
	}
	if lic == "" {
		return 0, fmt.Errorf("el campo Licencia es obligatorio")
	}
	apP := strings.TrimSpace(c["apellido_paterno"])
	apM := strings.TrimSpace(c["apellido_materno"])
	ci := strings.TrimSpace(c["ci"])
	tel := strings.TrimSpace(c["telefono"])
	fecha := Ahora()
	if id <= 0 {
		return CrearChofer(ctx, db, nombre, apP, apM, ci, tel, lic)
	}
	estado := strings.TrimSpace(c["estado"])
	if estado == "" {
		estado = "ACTIVO"
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var personaID int64
	if err := tx.QueryRowContext(ctx, `SELECT persona_id FROM choferes WHERE id=?`, id).Scan(&personaID); err != nil {
		return 0, err
	}
	_, err = tx.ExecContext(ctx, `
UPDATE personas SET nombre=?, apellido_paterno=?, apellido_materno=?, ci=?, telefono=?, estado=?, fecha_modificacion=? WHERE id=?`,
		nombre, apP, apM, ci, tel, estado, fecha, personaID)
	if err != nil {
		return 0, err
	}
	_, err = tx.ExecContext(ctx, `
UPDATE choferes SET licencia=?, telefono=?, estado=?, pendiente_push=1, fecha_modificacion=? WHERE id=?`,
		lic, tel, estado, fecha, id)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func listarProveedores(ctx context.Context, db *sql.DB) ([]AbmFila, error) {
	rows, err := db.QueryContext(ctx, `
SELECT pr.id, p.nombre, p.apellido_paterno, p.apellido_materno, p.ci, p.telefono,
       COALESCE(pr.nombre_empresa,''), COALESCE(pr.descripcion,''), pr.estado
FROM proveedores pr JOIN personas p ON p.id = pr.persona_id ORDER BY pr.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AbmFila
	for rows.Next() {
		var id int64
		var nom, apP, apM, ci, tel, empresa, desc, estado sql.NullString
		if err := rows.Scan(&id, &nom, &apP, &apM, &ci, &tel, &empresa, &desc, &estado); err != nil {
			return nil, err
		}
		nc := strings.TrimSpace(nom.String + " " + apP.String + " " + apM.String)
		out = append(out, AbmFila{
			ID: id, Estado: estado.String,
			Celdas: map[string]string{
				"nombre": nom.String, "apellido_paterno": apP.String, "apellido_materno": apM.String,
				"ci": ci.String, "telefono": tel.String, "empresa": empresa.String, "descripcion": desc.String,
				"nombre_completo": nc,
			},
		})
	}
	return out, rows.Err()
}

func guardarProveedorABM(ctx context.Context, db *sql.DB, id int64, c map[string]string) (int64, error) {
	nombre := strings.TrimSpace(c["nombre"])
	empresa := strings.TrimSpace(c["empresa"])
	if nombre == "" && empresa == "" {
		return 0, fmt.Errorf("indique el nombre o la empresa (al menos uno)")
	}
	if id <= 0 {
		return CrearProveedor(ctx, db, nombre, c["apellido_paterno"], c["apellido_materno"], c["ci"], c["telefono"], empresa, c["descripcion"])
	}
	estado := strings.TrimSpace(c["estado"])
	if estado == "" {
		estado = "ACTIVO"
	}
	fecha := Ahora()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var personaID int64
	if err := tx.QueryRowContext(ctx, `SELECT persona_id FROM proveedores WHERE id=?`, id).Scan(&personaID); err != nil {
		return 0, err
	}
	_, err = tx.ExecContext(ctx, `
UPDATE personas SET nombre=?, apellido_paterno=?, apellido_materno=?, ci=?, telefono=?, estado=?, fecha_modificacion=? WHERE id=?`,
		nombre, c["apellido_paterno"], c["apellido_materno"], c["ci"], c["telefono"], estado, fecha, personaID)
	if err != nil {
		return 0, err
	}
	_, err = tx.ExecContext(ctx, `
UPDATE proveedores SET nombre_empresa=?, descripcion=?, estado=?, pendiente_push=1, fecha_modificacion=? WHERE id=?`,
		empresa, c["descripcion"], estado, fecha, id)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func listarEmpresas(ctx context.Context, db *sql.DB) ([]AbmFila, error) {
	rows, err := db.QueryContext(ctx, `
SELECT id, nombre_empresa, COALESCE(actividad_economica,''), COALESCE(correo_electronico,''),
       COALESCE(telefono,''), estado, COALESCE(logo,'') FROM empresas ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AbmFila
	for rows.Next() {
		var id int64
		var nom, act, mail, tel, estado, logo string
		if err := rows.Scan(&id, &nom, &act, &mail, &tel, &estado, &logo); err != nil {
			return nil, err
		}
		out = append(out, AbmFila{
			ID: id, Estado: estado,
			Celdas: map[string]string{
				"nombre": nom, "actividad": act, "correo": mail, "telefono": tel, "logo": logo,
			},
		})
	}
	return out, rows.Err()
}

func guardarEmpresa(ctx context.Context, db *sql.DB, id int64, c map[string]string) error {
	nombre := strings.TrimSpace(c["nombre"])
	if nombre == "" {
		return fmt.Errorf("el campo Nombre es obligatorio")
	}
	act := strings.TrimSpace(c["actividad"])
	mail := strings.TrimSpace(c["correo"])
	tel := strings.TrimSpace(c["telefono"])
	logo := c["logo"]
	fecha := Ahora()
	if id <= 0 {
		_, err := db.ExecContext(ctx, `
INSERT INTO empresas (nombre_empresa, actividad_economica, correo_electronico, telefono, logo, estado, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, 'ACTIVO', ?, ?)`, nombre, act, mail, tel, logo, fecha, fecha)
		return err
	}
	estado := strings.TrimSpace(c["estado"])
	if estado == "" {
		estado = "ACTIVO"
	}
	_, err := db.ExecContext(ctx, `
UPDATE empresas SET nombre_empresa=?, actividad_economica=?, correo_electronico=?, telefono=?, logo=?, estado=?, fecha_modificacion=? WHERE id=?`,
		nombre, act, mail, tel, logo, estado, fecha, id)
	return err
}

func listarUsuarios(ctx context.Context, db *sql.DB) ([]AbmFila, error) {
	rows, err := db.QueryContext(ctx, `
SELECT u.id, p.nombre, p.apellido_paterno, p.apellido_materno, p.ci, u.email, u.rol, p.telefono, u.estado
FROM usuarios u JOIN personas p ON p.id = u.persona_id ORDER BY u.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AbmFila
	for rows.Next() {
		var id int64
		var nom, apP, apM, ci, email, rol, tel, estado sql.NullString
		if err := rows.Scan(&id, &nom, &apP, &apM, &ci, &email, &rol, &tel, &estado); err != nil {
			return nil, err
		}
		nc := strings.TrimSpace(nom.String + " " + apP.String + " " + apM.String)
		out = append(out, AbmFila{
			ID: id, Estado: estado.String,
			Celdas: map[string]string{
				"nombre": nom.String, "apellido_paterno": apP.String, "apellido_materno": apM.String,
				"ci": ci.String, "correo": email.String, "rol": rol.String, "telefono": tel.String,
				"nombre_completo": nc,
			},
		})
	}
	return out, rows.Err()
}

func guardarUsuario(ctx context.Context, db *sql.DB, id int64, c map[string]string) error {
	nombre := strings.TrimSpace(c["nombre"])
	correo := strings.TrimSpace(c["correo"])
	if nombre == "" {
		return fmt.Errorf("el campo Nombre es obligatorio")
	}
	if correo == "" {
		return fmt.Errorf("el campo Correo es obligatorio")
	}
	rol := strings.TrimSpace(c["rol"])
	if rol == "" {
		rol = "usuario"
	}
	clavePlana := c["clave"]
	apP := strings.TrimSpace(c["apellido_paterno"])
	apM := strings.TrimSpace(c["apellido_materno"])
	ci := strings.TrimSpace(c["ci"])
	tel := strings.TrimSpace(c["telefono"])
	fecha := Ahora()
	if id <= 0 {
		if strings.TrimSpace(clavePlana) == "" {
			return fmt.Errorf("el campo Contraseña es obligatorio")
		}
		hash := clave.Hash(clavePlana)
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		res, err := tx.ExecContext(ctx, `
INSERT INTO personas (nombre, apellido_paterno, apellido_materno, ci, correo, telefono, estado, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, ?, 'ACTIVO', ?, ?)`, nombre, apP, apM, ci, correo, tel, fecha, fecha)
		if err != nil {
			return err
		}
		personaID, _ := res.LastInsertId()
		_, err = tx.ExecContext(ctx, `
INSERT INTO usuarios (persona_id, email, password_hash, rol, estado, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, ?, 'ACTIVO', ?, ?)`, personaID, correo, hash, rol, fecha, fecha)
		if err != nil {
			return err
		}
		return tx.Commit()
	}
	estado := strings.TrimSpace(c["estado"])
	if estado == "" {
		estado = "ACTIVO"
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var personaID int64
	if err := tx.QueryRowContext(ctx, `SELECT persona_id FROM usuarios WHERE id=?`, id).Scan(&personaID); err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
UPDATE personas SET nombre=?, apellido_paterno=?, apellido_materno=?, ci=?, correo=?, telefono=?, estado=?, fecha_modificacion=? WHERE id=?`,
		nombre, apP, apM, ci, correo, tel, estado, fecha, personaID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(clavePlana) != "" {
		hash := clave.Hash(clavePlana)
		_, err = tx.ExecContext(ctx, `
UPDATE usuarios SET email=?, rol=?, estado=?, password_hash=?, fecha_modificacion=? WHERE id=?`,
			correo, rol, estado, hash, fecha, id)
	} else {
		_, err = tx.ExecContext(ctx, `
UPDATE usuarios SET email=?, rol=?, estado=?, fecha_modificacion=? WHERE id=?`,
			correo, rol, estado, fecha, id)
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}

func parsearEntero(texto, etiqueta string) (int, error) {
	texto = strings.TrimSpace(texto)
	if texto == "" {
		return 0, nil
	}
	var n int
	if _, err := fmt.Sscanf(texto, "%d", &n); err != nil {
		return 0, fmt.Errorf("%s: valor numérico inválido", etiqueta)
	}
	return n, nil
}

func validarTablaCatalogo(tabla string) error {
	switch tabla {
	case "productos", "origenes", "destinos":
		return nil
	default:
		return fmt.Errorf("tabla de catálogo no permitida: %s", tabla)
	}
}