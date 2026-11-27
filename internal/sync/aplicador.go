package sync

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"sistema-pesaje/internal/datos"
)

// AplicadorDelta aplica filas descargadas del servidor.
type AplicadorDelta struct {
	db *sql.DB
}

func NuevoAplicador(db *sql.DB) *AplicadorDelta {
	return &AplicadorDelta{db: db}
}

func (a *AplicadorDelta) AplicarCatalogo(recurso string, fila map[string]any) (bool, error) {
	id := entero(fila, "id")
	if id == nil || *id <= 0 {
		return false, nil
	}
	versionServidor := 0
	if v := entero(fila, "version"); v != nil {
		versionServidor = *v
	}
	versionLocal, err := a.versionLocal(recurso, *id)
	if err != nil {
		return false, err
	}
	pendiente, err := a.pendienteLocal(recurso, *id)
	if err != nil {
		return false, err
	}
	accion := ResolverVersion(versionServidor, versionLocal, pendiente)
	if accion == AccionNada {
		return false, nil
	}
	if accion == AccionConflicto {
		raw, _ := json.Marshal(fila)
		a.registrarConflicto(recurso, *id,
			fmt.Sprintf("version local %d / version servidor %d", versionLocal, versionServidor),
			string(raw))
	}
	ahora := datos.Ahora()
	switch recurso {
	case "vehiculos":
		return a.aplicarVehiculo(*id, fila, versionServidor, ahora)
	case "productos", "origenes", "destinos":
		return a.aplicarCatalogoSimple(recurso, *id, fila, versionServidor, ahora)
	case "choferes":
		return a.aplicarChofer(*id, fila, versionServidor, ahora)
	case "proveedores":
		return a.aplicarProveedor(*id, fila, versionServidor, ahora)
	default:
		return false, nil
	}
}

func (a *AplicadorDelta) AplicarPesaje(fila map[string]any) (bool, error) {
	id := entero(fila, "id")
	if id == nil || *id <= 0 {
		return false, nil
	}
	clientEventID := texto(fila, "client_event_id")
	guia := texto(fila, "guia")
	ahora := datos.Ahora()

	vehiculoID := entero(fila, "vehiculo_id")
	if vehiculoID == nil || *vehiculoID <= 0 {
		raw, _ := json.Marshal(fila)
		a.registrarConflicto("pesajes", *id,
			"el servidor envio un pesaje sin vehiculo: no se pudo guardar aca", string(raw))
		return false, nil
	}

	if clientEventID != "" {
		var n int
		err := a.db.QueryRowContext(context.Background(),
			`SELECT COUNT(*) FROM pesajes WHERE sync_event_id = ?`, clientEventID).Scan(&n)
		if err == nil && n > 0 {
			return false, nil
		}
	}

	bruto := valorEntero(fila, "peso_bruto", 0)
	tara := valorEntero(fila, "tara", 0)
	neto := valorEntero(fila, "peso_neto", 0)
	estadoBalanza := texto(fila, "estado_balanza")
	if estadoBalanza == "" {
		estadoBalanza = "EN_PROCESO"
	}

	var existe int
	err := a.db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM pesajes WHERE web_id = ?`, *id).Scan(&existe)
	if err != nil {
		return false, err
	}
	if existe > 0 {
		_, err = a.db.ExecContext(context.Background(), `
UPDATE pesajes SET guia = ?, peso_bruto = ?, tara = ?, peso_neto = ?, estado_balanza = ?, fecha_modificacion = ? WHERE web_id = ?`,
			guia, bruto, tara, neto, estadoBalanza, ahora, *id)
		return err == nil, err
	}

	costo := valorEntero(fila, "costo_bs", 0)
	flete := valorEntero(fila, "flete_bs_pendiente", 0)
	_, err = a.db.ExecContext(context.Background(), `
INSERT INTO pesajes (id, guia, web_id, sync_event_id, vehiculo_id, chofer_id, proveedor_id, producto_id,
  id_origen, id_destino, peso_bruto, tara, peso_neto, costo_bs, flete_bs_pendiente, estado, estado_balanza,
  origen, sincronizado, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'ACTIVO', ?, 'WEB', 1, ?, ?)`,
		-*id, guia, *id, nullString(clientEventID),
		*vehiculoID, enteroNull(fila, "chofer_id"), enteroNull(fila, "proveedor_id"),
		enteroNull(fila, "producto_id"), enteroNull(fila, "id_origen"), enteroNull(fila, "id_destino"),
		bruto, tara, neto, costo, flete, estadoBalanza, ahora, ahora)
	return err == nil, err
}

func (a *AplicadorDelta) ActualizarPesajeSincronizado(pesajeID int, webID *int, guia string) error {
	ahora := datos.Ahora()
	if guia != "" {
		_, err := a.db.ExecContext(context.Background(), `
UPDATE pesajes SET sincronizado = 1, sync_error = NULL, guia = ?, web_id = ?, fecha_modificacion = ? WHERE id = ?`,
			guia, nullIntPtr(webID), ahora, pesajeID)
		return err
	}
	_, err := a.db.ExecContext(context.Background(), `
UPDATE pesajes SET sincronizado = 1, sync_error = NULL, web_id = ?, fecha_modificacion = ? WHERE id = ?`,
		nullIntPtr(webID), ahora, pesajeID)
	return err
}

func (a *AplicadorDelta) LeerCursor(recurso string) (string, error) {
	var cur sql.NullString
	err := a.db.QueryRowContext(context.Background(), `SELECT cursor FROM sync_cursor WHERE recurso = ?`, recurso).Scan(&cur)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return cur.String, nil
}

func (a *AplicadorDelta) GuardarCursor(recurso, valor string) error {
	ahora := datos.Ahora()
	res, err := a.db.ExecContext(context.Background(), `
UPDATE sync_cursor SET cursor = ?, ultima_sync = ? WHERE recurso = ?`, valor, ahora, recurso)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		_, err = a.db.ExecContext(context.Background(), `
INSERT INTO sync_cursor (recurso, cursor, ultima_sync) VALUES (?, ?, ?)`, recurso, valor, ahora)
	}
	return err
}

func (a *AplicadorDelta) registrarConflicto(recurso string, entidadID int, detalle, servidorJSON string) {
	_, _ = a.db.ExecContext(context.Background(), `
INSERT INTO sync_conflictos (recurso, entidad_id, cliente_json, servidor_json, detectado_en) VALUES (?, ?, ?, ?, ?)`,
		recurso, entidadID, detalle, servidorJSON, datos.Ahora())
}

func (a *AplicadorDelta) versionLocal(recurso string, id int) (int, error) {
	var v sql.NullInt64
	err := a.db.QueryRowContext(context.Background(),
		fmt.Sprintf("SELECT version FROM %s WHERE id = ?", recurso), id).Scan(&v)
	if err == sql.ErrNoRows {
		return -1, nil
	}
	if err != nil {
		return -1, err
	}
	if !v.Valid {
		return -1, nil
	}
	return int(v.Int64), nil
}

func (a *AplicadorDelta) LimpiarPendientePush(recurso string, id int) error {
	tablas := map[string]bool{
		"productos": true, "origenes": true, "destinos": true,
		"vehiculos": true, "choferes": true, "proveedores": true,
	}
	if !tablas[recurso] {
		return nil
	}
	_, err := a.db.ExecContext(context.Background(),
		fmt.Sprintf(`UPDATE %s SET pendiente_push = 0, fecha_modificacion = ? WHERE id = ?`, recurso),
		datos.Ahora(), id)
	return err
}

func (a *AplicadorDelta) pendienteLocal(recurso string, id int) (bool, error) {
	var v sql.NullInt64
	err := a.db.QueryRowContext(context.Background(),
		fmt.Sprintf("SELECT pendiente_push FROM %s WHERE id = ?", recurso), id).Scan(&v)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return v.Valid && v.Int64 != 0, nil
}

func (a *AplicadorDelta) aplicarVehiculo(id int, fila map[string]any, version int, ahora string) (bool, error) {
	placa := texto(fila, "placa")
	tipo := texto(fila, "tipo_vehiculo")
	tara := valorEntero(fila, "tara", 0)
	estado := estadoFila(fila)
	existe, err := a.existeFila("vehiculos", id)
	if err != nil {
		return false, err
	}
	if existe {
		_, err = a.db.ExecContext(context.Background(), `
UPDATE vehiculos SET placa = ?, tipo_vehiculo = ?, tara = ?, estado = ?, version = ?, pendiente_push = 0, fecha_modificacion = ? WHERE id = ?`,
			placa, tipo, tara, estado, version, ahora, id)
		return err == nil, err
	}
	_, err = a.db.ExecContext(context.Background(), `
INSERT INTO vehiculos (id, placa, tipo_vehiculo, tara, estado, version, pendiente_push, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, ?, 0, ?, ?)`,
		id, placa, tipo, tara, estado, version, ahora, ahora)
	return err == nil, err
}

func (a *AplicadorDelta) aplicarCatalogoSimple(tabla string, id int, fila map[string]any, version int, ahora string) (bool, error) {
	nombre := texto(fila, "nombre")
	descripcion := texto(fila, "descripcion")
	estado := estadoFila(fila)
	existe, err := a.existeFila(tabla, id)
	if err != nil {
		return false, err
	}
	if existe {
		_, err = a.db.ExecContext(context.Background(), fmt.Sprintf(`
UPDATE %s SET nombre = ?, descripcion = ?, estado = ?, version = ?, pendiente_push = 0, fecha_modificacion = ? WHERE id = ?`, tabla),
			nombre, descripcion, estado, version, ahora, id)
		return err == nil, err
	}
	_, err = a.db.ExecContext(context.Background(), fmt.Sprintf(`
INSERT INTO %s (id, nombre, descripcion, estado, version, pendiente_push, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, 0, ?, ?)`, tabla),
		id, nombre, descripcion, estado, version, ahora, ahora)
	return err == nil, err
}

func (a *AplicadorDelta) aplicarChofer(id int, fila map[string]any, version int, ahora string) (bool, error) {
	personaID, err := a.resolverPersona(fila, ahora)
	if err != nil || personaID == nil {
		return false, err
	}
	licencia := texto(fila, "licencia")
	telefono := texto(fila, "telefono")
	estado := estadoFila(fila)
	existe, err := a.existeFila("choferes", id)
	if err != nil {
		return false, err
	}
	if existe {
		_, err = a.db.ExecContext(context.Background(), `
UPDATE choferes SET persona_id = ?, licencia = ?, telefono = ?, estado = ?, version = ?, pendiente_push = 0, fecha_modificacion = ? WHERE id = ?`,
			*personaID, licencia, telefono, estado, version, ahora, id)
		return err == nil, err
	}
	_, err = a.db.ExecContext(context.Background(), `
INSERT INTO choferes (id, persona_id, licencia, telefono, estado, version, pendiente_push, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, ?, 0, ?, ?)`,
		id, *personaID, licencia, telefono, estado, version, ahora, ahora)
	return err == nil, err
}

func (a *AplicadorDelta) aplicarProveedor(id int, fila map[string]any, version int, ahora string) (bool, error) {
	personaID, err := a.resolverPersona(fila, ahora)
	if err != nil || personaID == nil {
		return false, err
	}
	nombreEmpresa := texto(fila, "nombre_empresa")
	descripcion := texto(fila, "descripcion")
	estado := estadoFila(fila)
	existe, err := a.existeFila("proveedores", id)
	if err != nil {
		return false, err
	}
	if existe {
		_, err = a.db.ExecContext(context.Background(), `
UPDATE proveedores SET persona_id = ?, nombre_empresa = ?, descripcion = ?, estado = ?, version = ?, pendiente_push = 0, fecha_modificacion = ? WHERE id = ?`,
			*personaID, nombreEmpresa, descripcion, estado, version, ahora, id)
		return err == nil, err
	}
	_, err = a.db.ExecContext(context.Background(), `
INSERT INTO proveedores (id, persona_id, nombre_empresa, descripcion, estado, version, pendiente_push, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, ?, 0, ?, ?)`,
		id, *personaID, nombreEmpresa, descripcion, estado, version, ahora, ahora)
	return err == nil, err
}

func (a *AplicadorDelta) existeFila(tabla string, id int) (bool, error) {
	var n int
	err := a.db.QueryRowContext(context.Background(),
		fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ?", tabla), id).Scan(&n)
	return n > 0, err
}

func (a *AplicadorDelta) resolverPersona(fila map[string]any, ahora string) (*int, error) {
	nombre := texto(fila, "nombre")
	if nombre == "" {
		nombre = texto(fila, "nombre_empresa")
	}
	apPat := texto(fila, "apellido_paterno")
	apMat := texto(fila, "apellido_materno")
	ci := texto(fila, "ci")
	telefono := texto(fila, "telefono")

	if nombre != "" {
		var id int64
		err := a.db.QueryRowContext(context.Background(), `
SELECT id FROM personas WHERE nombre = ? AND IFNULL(apellido_paterno, '') = ? AND IFNULL(apellido_materno, '') = ? LIMIT 1`,
			nombre, apPat, apMat).Scan(&id)
		if err == nil {
			i := int(id)
			return &i, nil
		}
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	res, err := a.db.ExecContext(context.Background(), `
INSERT INTO personas (nombre, apellido_paterno, apellido_materno, ci, telefono, estado, fecha_creacion, fecha_modificacion)
VALUES (?, ?, ?, ?, ?, 'ACTIVO', ?, ?)`,
		nombre, nullString(apPat), nullString(apMat), nullString(ci), nullString(telefono), ahora, ahora)
	if err != nil {
		return nil, err
	}
	lid, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	i := int(lid)
	return &i, nil
}

func estadoFila(fila map[string]any) string {
	e := texto(fila, "estado")
	if e == "" {
		return "ACTIVO"
	}
	return e
}

func valorEntero(fila map[string]any, campo string, def int) int {
	if v := entero(fila, campo); v != nil {
		return *v
	}
	return def
}

func enteroNull(fila map[string]any, campo string) any {
	if v := entero(fila, campo); v != nil && *v > 0 {
		return *v
	}
	return nil
}

func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nullIntPtr(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}