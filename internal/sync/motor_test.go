package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sistema-pesaje/internal/datos"
)

func prepararMotor(t *testing.T) (*datos.Database, func()) {
	t.Helper()
	dir := t.TempDir()
	ruta := filepath.Join(dir, "pesaje.db")
	t.Setenv(datos.PropiedadSistemaRuta, ruta)
	datos.CerrarInstancia()
	base, err := datos.ObtenerInstancia()
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	return base, func() {
		datos.CerrarInstancia()
		_ = os.RemoveAll(dir)
	}
}

func TestPideTokenAntesDeEnviar(t *testing.T) {
	transporte := NuevoTransporteFalso()
	transporte.Responder("POST http://servidor/api/auth/token",
		RespuestaDeCodigo(200, `{"access_token":"abc","refresh_token":"r1"}`))
	sesion := NuevaSesion(transporte, "http://servidor", "admin@sistema.com", "clave")
	if !sesion.AsegurarToken() {
		t.Fatal("token")
	}
	if !sesion.TieneToken() {
		t.Fatal("sin token")
	}
	h := sesion.Cabeceras()["Cookie"]
	if len(h) < len("access_token=") || h[len("access_token="):] != "abc" {
		t.Fatalf("cookie %q", h)
	}
}

func TestSinTokenValidoNoIntentaEnviar(t *testing.T) {
	base, limpiar := prepararMotor(t)
	defer limpiar()
	transporte := NuevoTransporteFalso()
	transporte.Responder("POST http://servidor/api/auth/token",
		RespuestaDeCodigo(401, `{"message":"Credenciales invalidas"}`))
	outbox := NuevoOutbox(base.Conn())
	motor := NuevoMotor(transporte, NuevaSesion(transporte, "http://servidor", "a@b.c", "x"), outbox, NuevoAplicador(base.Conn()))
	res := motor.Sincronizar()
	if res.Enviados != 0 {
		t.Fatalf("enviados %d", res.Enviados)
	}
	if res.Mensaje == "" {
		t.Fatal("mensaje vacio")
	}
}

func TestElPayloadDePesajeNoLlevaGuia(t *testing.T) {
	payload := PayloadPesaje("evento-1", PesajeParaEnvio{
		ID:            125,
		SyncEventID:   "evento-1",
		Guia:          "PESO-2026-000001",
		VehiculoID:    7,
		PesoBruto:     10000,
		Tara:          4000,
		PesoNeto:      6000,
		Estado:        "ACTIVO",
		EstadoBalanza: "FINALIZADO",
	})
	if strings.Contains(payload, "sync_event_id") || strings.Contains(payload, "client_event_id") {
		t.Fatalf("ids de cola en payload: %s", payload)
	}
	if payload == "" || payload[0] != '{' {
		t.Fatal("json")
	}
	if !strings.Contains(payload, "pesaje_id") {
		t.Fatalf("falta pesaje_id: %s", payload)
	}
}

func TestPublicaPesoEnVivoEnLaWeb(t *testing.T) {
	base, limpiar := prepararMotor(t)
	defer limpiar()
	transporte := NuevoTransporteFalso()
	transporte.Responder("POST http://servidor/api/auth/token",
		RespuestaDeCodigo(200, `{"access_token":"abc","refresh_token":"r1"}`))
	transporte.Responder("POST http://servidor/api/pesajes/ultimo-peso",
		RespuestaDeCodigo(200, `{"success":true}`))
	motor := NuevoMotor(transporte, NuevaSesion(transporte, "http://servidor", "a@b.c", "clave"),
		NuevoOutbox(base.Conn()), NuevoAplicador(base.Conn()))
	if !motor.EnviarPesoVivo(12450) {
		t.Fatal("envio vivo")
	}
	if !transporte.ContieneRuta("/api/pesajes/ultimo-peso") {
		t.Fatal("no llamo ultimo-peso")
	}
}