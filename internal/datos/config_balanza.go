package datos

import (
	"database/sql"
	"os"
	"strings"

	"sistema-pesaje/internal/balanza"
)

// ObtenerConfigBalanza devuelve la fila activa o valores por defecto.
func (d *Database) ObtenerConfigBalanza() balanza.Config {
	def := balanza.ConfigPorDefecto()
	row := d.Conn().QueryRow(`SELECT puerto_com, baudrate, databits, paridad, stopbits, flowcontrol, timeout_ms,
		metodo_lectura, posicion_inicio, posicion_longitud FROM config_balanza ORDER BY id DESC LIMIT 1`)
	var puerto, paridad, flow, metodo sql.NullString
	var baud, data, stop, timeout, inicio, longitud sql.NullInt64
	err := row.Scan(&puerto, &baud, &data, &paridad, &stop, &flow, &timeout, &metodo, &inicio, &longitud)
	if err == sql.ErrNoRows || err != nil {
		return def
	}
	if puerto.Valid {
		def.PuertoCOM = puerto.String
	}
	if baud.Valid {
		def.BaudRate = int(baud.Int64)
	}
	if data.Valid {
		def.DataBits = int(data.Int64)
	}
	if paridad.Valid {
		def.Paridad = paridad.String
	}
	if stop.Valid {
		def.StopBits = int(stop.Int64)
	}
	if flow.Valid {
		def.FlowControl = flow.String
	}
	if timeout.Valid {
		def.TimeoutMs = int(timeout.Int64)
	}
	if metodo.Valid {
		def.MetodoLectura = balanza.MetodoDesdeTexto(metodo.String)
	}
	if inicio.Valid {
		def.PosicionInicio = int(inicio.Int64)
	}
	if longitud.Valid {
		def.PosicionLongitud = int(longitud.Int64)
	}
	return def
}

// UsarBalanzaSimulada indica simulación (desarrollo o puerto SIM).
func UsarBalanzaSimulada(cfg balanza.Config) bool {
	if os.Getenv("PESAJE_BALANZA_SIM") == "1" {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(cfg.PuertoCOM), "SIM")
}
