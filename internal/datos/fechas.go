package datos

import "time"

const formatoFechaHora = "2006-01-02 15:04:05"

// FormatoFechaHora es el formato usado en SQLite y sync.
const FormatoFechaHora = formatoFechaHora

// Ahora devuelve la marca de tiempo del sistema en formato SQLite.
func Ahora() string {
	return time.Now().Format(formatoFechaHora)
}

// SoloFecha formatea dd/MM/yyyy para grillas.
func SoloFecha(fechaHora string) string {
	if len(fechaHora) < 10 {
		return ""
	}
	if len(fechaHora) >= 16 {
		return fechaHora[8:10] + "/" + fechaHora[5:7] + "/" + fechaHora[0:4]
	}
	return fechaHora[:10]
}

// SoloHora devuelve HH:mm.
func SoloHora(fechaHora string) string {
	if len(fechaHora) < 16 {
		return ""
	}
	return fechaHora[11:16]
}