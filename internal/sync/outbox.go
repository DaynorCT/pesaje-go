package sync

import (
	"context"
	"database/sql"
	"strings"

	"sistema-pesaje/internal/datos"
)

const (
	EstadoPendiente = "PENDIENTE"
	EstadoEnviado   = "ENVIADO"
	EstadoError     = "ERROR"
	LimitePorVuelta = 50
)

// Outbox cola durable hacia la web.
type Outbox struct {
	db *sql.DB
}

func NuevoOutbox(db *sql.DB) *Outbox {
	return &Outbox{db: db}
}

func (o *Outbox) Encolar(clientEventID, recurso string, entidadID int, payload string) error {
	ahora := datos.Ahora()
	_, err := o.db.ExecContext(context.Background(), `
INSERT INTO sync_outbox (client_event_id, recurso, entidad_id, payload, estado, intentos, creado_en, actualizado_en)
VALUES (?, ?, ?, ?, ?, 0, ?, ?)`,
		clientEventID, recurso, entidadID, payload, EstadoPendiente, ahora, ahora,
	)
	return err
}

func (o *Outbox) Pendientes() ([]EventoOutbox, error) {
	ahora := datos.Ahora()
	rows, err := o.db.QueryContext(context.Background(), `
SELECT client_event_id, recurso, entidad_id, payload, intentos FROM sync_outbox
WHERE estado = ? AND (proximo_intento IS NULL OR proximo_intento <= ?)
ORDER BY id ASC LIMIT ?`,
		EstadoPendiente, ahora, LimitePorVuelta,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []EventoOutbox
	for rows.Next() {
		var e EventoOutbox
		var entidad sql.NullInt64
		if err := rows.Scan(&e.ClientEventID, &e.Recurso, &entidad, &e.Payload, &e.Intentos); err != nil {
			return nil, err
		}
		if entidad.Valid {
			e.EntidadID = int(entidad.Int64)
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (o *Outbox) MarcarEnviado(clientEventID string, webID *int) error {
	ahora := datos.Ahora()
	_, err := o.db.ExecContext(context.Background(), `
UPDATE sync_outbox SET estado = ?, web_id = ?, ultimo_error = NULL, proximo_intento = NULL, actualizado_en = ? WHERE client_event_id = ?`,
		EstadoEnviado, nullableWebID(webID), ahora, clientEventID,
	)
	return err
}

func (o *Outbox) MarcarError(clientEventID, mensaje string, permanente bool, intentos int, proximoIntento *string) error {
	ahora := datos.Ahora()
	msg := recortar(mensaje)
	if permanente {
		_, err := o.db.ExecContext(context.Background(), `
UPDATE sync_outbox SET estado = ?, intentos = ?, ultimo_error = ?, proximo_intento = NULL, actualizado_en = ? WHERE client_event_id = ?`,
			EstadoError, intentos, msg, ahora, clientEventID,
		)
		return err
	}
	_, err := o.db.ExecContext(context.Background(), `
UPDATE sync_outbox SET intentos = ?, ultimo_error = ?, proximo_intento = ?, actualizado_en = ? WHERE client_event_id = ?`,
		intentos, msg, proximoIntento, ahora, clientEventID,
	)
	return err
}

func (o *Outbox) PendientesCount() (int64, error) {
	var n int64
	err := o.db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM sync_outbox WHERE estado = ?`, EstadoPendiente).Scan(&n)
	return n, err
}

func (o *Outbox) ErroresCount() (int64, error) {
	var n int64
	err := o.db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM sync_outbox WHERE estado = ?`, EstadoError).Scan(&n)
	return n, err
}

func (o *Outbox) UltimoError() string {
	var msg sql.NullString
	err := o.db.QueryRowContext(context.Background(), `
SELECT ultimo_error FROM sync_outbox WHERE ultimo_error IS NOT NULL ORDER BY id DESC LIMIT 1`).Scan(&msg)
	if err != nil || !msg.Valid {
		return ""
	}
	return msg.String
}

func recortar(mensaje string) string {
	if len(mensaje) > 190 {
		return mensaje[:190]
	}
	return strings.TrimSpace(mensaje)
}

func nullableWebID(id *int) any {
	if id == nil {
		return nil
	}
	return *id
}