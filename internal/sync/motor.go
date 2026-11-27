package sync

import (
	"log"
	"strconv"
	"time"
)

// RecursosCatalogo que se descargan de la web.
var RecursosCatalogo = []string{
	"vehiculos", "productos", "origenes", "destinos", "choferes", "proveedores",
}

// Motor ejecuta push/pull con el servidor.
type Motor struct {
	transporte Transporte
	sesion     *Sesion
	outbox     *Outbox
	aplicador  *AplicadorDelta
}

func NuevoMotor(t Transporte, sesion *Sesion, outbox *Outbox, aplicador *AplicadorDelta) *Motor {
	return &Motor{transporte: t, sesion: sesion, outbox: outbox, aplicador: aplicador}
}

func (m *Motor) Sincronizar() ResultadoSync {
	resultado := ResultadoSync{}
	if !m.sesion.AsegurarToken() {
		resultado.Mensaje = m.sesion.UltimoError()
		return resultado
	}
	m.procesarOutbox(&resultado)
	m.descargarCatalogos(&resultado)
	m.descargarPesajes(&resultado)
	return resultado
}

func (m *Motor) EnviarPesoVivo(peso int) bool {
	if peso < 0 {
		return false
	}
	if !m.sesion.AsegurarToken() {
		return false
	}
	cuerpo := escribirJSON(mapaValores("peso", peso))
	url := m.sesion.URLBase() + "/api/pesajes/ultimo-peso"
	resp := m.transporte.Enviar("POST", url, cuerpo, m.sesion.Cabeceras())
	if resp.Codigo == 401 {
		m.sesion.Invalidar()
		if m.sesion.AsegurarToken() {
			resp = m.transporte.Enviar("POST", url, cuerpo, m.sesion.Cabeceras())
		}
	}
	return resp.Exito()
}

func (m *Motor) procesarOutbox(resultado *ResultadoSync) {
	eventos, err := m.outbox.Pendientes()
	if err != nil {
		resultado.Mensaje = "No se pudo leer la cola local."
		return
	}
	for _, evento := range eventos {
		func() {
			defer func() {
				if r := recover(); r != nil {
					intentos := evento.Intentos + 1
					prox := CalcularProximoIntento(intentos, time.Now(), 0)
					_ = m.outbox.MarcarError(evento.ClientEventID, "Error inesperado", false, intentos, &prox)
					resultado.SumarError()
					resultado.Mensaje = "Un evento quedo con error y se reintentara mas tarde."
				}
			}()
			m.enviarEvento(evento, resultado)
		}()
	}
}

func (m *Motor) enviarEvento(evento EventoOutbox, resultado *ResultadoSync) {
	ruta := "/api/sync/" + evento.Recurso
	if evento.Recurso == "pesajes" {
		ruta = "/api/pesajes"
	}
	url := m.sesion.URLBase() + ruta
	resp := m.transporte.Enviar("POST", url, evento.Payload, m.sesion.Cabeceras())
	if resp.Codigo == 401 {
		m.sesion.Invalidar()
		if m.sesion.AsegurarToken() {
			resp = m.transporte.Enviar("POST", url, evento.Payload, m.sesion.Cabeceras())
		}
	}
	if resp.Exito() {
		m.aplicarConfirmacion(evento, resp)
		webID := webIDDe(resp)
		_ = m.outbox.MarcarEnviado(evento.ClientEventID, webID)
		if evento.Recurso != "pesajes" && evento.EntidadID > 0 {
			_ = m.aplicador.LimpiarPendientePush(evento.Recurso, evento.EntidadID)
		}
		resultado.SumarEnviado()
		return
	}
	if EsErrorPermanente(resp.Codigo) {
		intentos := evento.Intentos + 1
		_ = m.outbox.MarcarError(evento.ClientEventID,
			"HTTP "+itoa(resp.Codigo)+" "+resp.Cuerpo, true, intentos, nil)
		resultado.SumarError()
		resultado.Mensaje = "Hay eventos que el servidor rechazo y necesitan revision."
		return
	}
	intentos := evento.Intentos + 1
	prox := CalcularProximoIntento(intentos, time.Now(), 0)
	_ = m.outbox.MarcarError(evento.ClientEventID,
		"HTTP "+itoa(resp.Codigo)+" "+resp.Cuerpo, false, intentos, &prox)
	resultado.SumarError()
	resultado.Mensaje = "No se pudo enviar un evento; se reintentara mas tarde."
	if resp.FalloDeRed() {
		resultado.Mensaje = "Sin conexion con la web. La cola sigue esperando."
		return
	}
}

func (m *Motor) descargarCatalogos(resultado *ResultadoSync) {
	for _, recurso := range RecursosCatalogo {
		m.descargar(recurso, resultado, false)
	}
}

func (m *Motor) descargarPesajes(resultado *ResultadoSync) {
	m.descargar("pesajes", resultado, true)
}

func (m *Motor) descargar(recurso string, resultado *ResultadoSync, esPesaje bool) {
	desde, _ := m.aplicador.LeerCursor(recurso)
	ruta := "/api/sync/" + recurso + "?since=" + desde
	url := m.sesion.URLBase() + ruta
	resp := m.transporte.Enviar("GET", url, "", m.sesion.Cabeceras())
	if resp.Codigo == 401 {
		m.sesion.Invalidar()
		if m.sesion.AsegurarToken() {
			resp = m.transporte.Enviar("GET", url, "", m.sesion.Cabeceras())
		}
	}
	if !resp.Exito() {
		resultado.Mensaje = "No se pudo descargar " + recurso + " (" + resp.String() + ")."
		return
	}
	raiz, err := parseJSON(resp.Cuerpo)
	if err != nil || raiz == nil {
		return
	}
	for _, fila := range arrayDatos(raiz) {
		var aplicado bool
		var err error
		if esPesaje {
			aplicado, err = m.aplicador.AplicarPesaje(fila)
		} else {
			aplicado, err = m.aplicador.AplicarCatalogo(recurso, fila)
		}
		if err != nil {
			resultado.SumarError()
			resultado.Mensaje = "Una fila de " + recurso + " no se pudo aplicar: " + err.Error()
			continue
		}
		if aplicado {
			resultado.SumarDescargado()
		}
	}
	cursor := texto(raiz, "cursor")
	if cursor != "" {
		_ = m.aplicador.GuardarCursor(recurso, cursor)
	}
}

func (m *Motor) aplicarConfirmacion(evento EventoOutbox, resp RespuestaHTTP) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[sync] confirmacion %s: %v", evento.ClientEventID, r)
		}
	}()
	raiz, err := parseJSON(resp.Cuerpo)
	if err != nil || raiz == nil {
		return
	}
	datosNodo := raiz
	if d, ok := raiz["data"].(map[string]any); ok {
		datosNodo = d
	}
	if evento.Recurso != "pesajes" {
		return
	}
	if evento.EntidadID == 0 {
		return
	}
	id := entero(datosNodo, "id")
	guia := texto(datosNodo, "guia")
	if id != nil {
		_ = m.aplicador.ActualizarPesajeSincronizado(evento.EntidadID, id, guia)
	}
}

func webIDDe(resp RespuestaHTTP) *int {
	raiz, err := parseJSON(resp.Cuerpo)
	if err != nil || raiz == nil {
		return nil
	}
	if d, ok := raiz["data"].(map[string]any); ok {
		return entero(d, "id")
	}
	return entero(raiz, "id")
}

func itoa(n int) string {
	return strconv.Itoa(n)
}