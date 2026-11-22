package sync

// RespuestaHTTP resultado de una petición al servidor web.
type RespuestaHTTP struct {
	Codigo int
	Cuerpo string
	Error  string
}

func RespuestaDeCodigo(codigo int, cuerpo string) RespuestaHTTP {
	return RespuestaHTTP{Codigo: codigo, Cuerpo: cuerpo}
}

func RespuestaDeFallo(err string) RespuestaHTTP {
	return RespuestaHTTP{Codigo: 0, Error: err}
}

func (r RespuestaHTTP) Exito() bool {
	return r.Codigo >= 200 && r.Codigo < 300
}

func (r RespuestaHTTP) FalloDeRed() bool {
	return r.Codigo == 0
}

// Transporte abstrae HTTP (pruebas y producción).
type Transporte interface {
	Enviar(metodo, url, cuerpo string, cabeceras map[string]string) RespuestaHTTP
}

// EventoOutbox fila pendiente de envío.
type EventoOutbox struct {
	ClientEventID string
	Recurso       string
	EntidadID     int
	Payload       string
	Intentos      int
}

// ResultadoSync resume una vuelta de sincronización.
type ResultadoSync struct {
	Enviados     int
	Conflictos   int
	Errores      int
	Descargados  int
	Mensaje      string
}

func (r *ResultadoSync) SumarEnviado()     { r.Enviados++ }
func (r *ResultadoSync) SumarConflicto()   { r.Conflictos++ }
func (r *ResultadoSync) SumarError()       { r.Errores++ }
func (r *ResultadoSync) SumarDescargado()  { r.Descargados++ }

// Limpio indica que no quedaron errores en la vuelta.
func (r ResultadoSync) Limpio() bool {
	return r.Errores == 0
}