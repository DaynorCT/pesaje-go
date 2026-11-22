package sync

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TransporteHTTP implementación real con net/http.
type TransporteHTTP struct {
	client *http.Client
}

func NuevoTransporteHTTP(timeoutMs int) *TransporteHTTP {
	if timeoutMs <= 0 {
		timeoutMs = 4000
	}
	d := time.Duration(timeoutMs) * time.Millisecond
	return &TransporteHTTP{
		client: &http.Client{Timeout: d},
	}
}

func (t *TransporteHTTP) Enviar(metodo, url, cuerpo string, cabeceras map[string]string) RespuestaHTTP {
	req, err := http.NewRequest(metodo, url, nil)
	if err != nil {
		return RespuestaDeFallo(err.Error())
	}
	if cuerpo != "" {
		req.Body = io.NopCloser(strings.NewReader(cuerpo))
		req.ContentLength = int64(len(cuerpo))
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	for k, v := range cabeceras {
		req.Header.Set(k, v)
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return RespuestaDeFallo(err.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return RespuestaDeCodigo(resp.StatusCode, string(body))
}

// TransporteFalso para pruebas del motor.
type TransporteFalso struct {
	respuestas map[string]RespuestaHTTP
	defecto    RespuestaHTTP
	Peticiones []string
}

func NuevoTransporteFalso() *TransporteFalso {
	return &TransporteFalso{
		respuestas: make(map[string]RespuestaHTTP),
		defecto:    RespuestaDeCodigo(200, "{}"),
	}
}

func (f *TransporteFalso) Responder(metodoRuta string, r RespuestaHTTP) {
	f.respuestas[metodoRuta] = r
}

func (f *TransporteFalso) Enviar(metodo, url, cuerpo string, cabeceras map[string]string) RespuestaHTTP {
	clave := metodo + " " + url
	f.Peticiones = append(f.Peticiones, clave)
	if r, ok := f.respuestas[clave]; ok {
		return r
	}
	return f.defecto
}

func (f *TransporteFalso) ContieneRuta(sub string) bool {
	for _, p := range f.Peticiones {
		if strings.Contains(p, sub) {
			return true
		}
	}
	return false
}

var _ Transporte = (*TransporteHTTP)(nil)

// String corto para depuración (no usado en producción).
func (r RespuestaHTTP) String() string {
	if r.Error != "" {
		return fmt.Sprintf("HTTP %d (%s)", r.Codigo, r.Error)
	}
	return fmt.Sprintf("HTTP %d", r.Codigo)
}