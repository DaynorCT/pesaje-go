package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"sistema-pesaje/internal/datos"
)

func prepararAppConBalanzaSim(t *testing.T) (*App, func()) {
	t.Helper()
	dir := t.TempDir()
	ruta := filepath.Join(dir, "pesaje.db")
	t.Setenv(datos.PropiedadSistemaRuta, ruta)
	datos.CerrarInstancia()
	base, err := datos.ObtenerInstancia()
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	_, err = base.Conn().Exec(`INSERT INTO config_balanza (
		puerto_com, baudrate, databits, paridad, stopbits, flowcontrol, timeout_ms,
		metodo_lectura, posicion_inicio, posicion_longitud
	) VALUES ('SIM', 9600, 8, 'N', 1, 'None', 1000, 'AUTO', 8, 5)`)
	if err != nil {
		t.Fatalf("config báscula: %v", err)
	}
	app := NewApp()
	app.db = base
	return app, func() {
		app.PesajeCerrar()
		datos.CerrarInstancia()
		_ = os.RemoveAll(dir)
	}
}

func esperarPesoVivo(t *testing.T, app *App) int {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		est := app.PesajeEstado()
		if est.BalanzaConectada && est.PesoVivo > 0 {
			return est.PesoVivo
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("timeout esperando peso vivo de báscula simulada")
	return 0
}

func TestModuloPesajeConectaBalanzaSimuladaYCaptura(t *testing.T) {
	app, limpiar := prepararAppConBalanzaSim(t)
	defer limpiar()

	r := app.PesajeAlternarBalanza()
	if !r.OK {
		t.Fatalf("conectar: %+v", r.Aviso)
	}
	if r.EstadoPesaje == nil || !r.EstadoPesaje.BalanzaConectada {
		t.Fatal("debe quedar conectada")
	}

	vivo := esperarPesoVivo(t, app)

	bruto := app.PesajeCapturarBruto()
	if !bruto.OK || bruto.EstadoPesaje == nil || bruto.EstadoPesaje.PesoBruto != vivo {
		t.Fatalf("captura bruto: ok=%v bruto=%v vivo=%d", bruto.OK, bruto.EstadoPesaje, vivo)
	}

	tara := app.PesajeCapturarTara()
	if !tara.OK || tara.EstadoPesaje == nil || tara.EstadoPesaje.PesoTara <= 0 {
		t.Fatal("captura tara")
	}
	if tara.EstadoPesaje.PesoNeto != bruto.EstadoPesaje.PesoBruto-tara.EstadoPesaje.PesoTara {
		t.Fatalf("neto=%d esperado=%d", tara.EstadoPesaje.PesoNeto, bruto.EstadoPesaje.PesoBruto-tara.EstadoPesaje.PesoTara)
	}

	desc := app.PesajeAlternarBalanza()
	if !desc.OK || desc.EstadoPesaje == nil || desc.EstadoPesaje.BalanzaConectada {
		t.Fatal("desconectar")
	}
}

func TestModoIntegradoConBalanzaSimulada(t *testing.T) {
	app, limpiar := prepararAppConBalanzaSim(t)
	defer limpiar()

	r := app.AlternarBalanza()
	if !r.OK || r.Estado == nil || !r.Estado.BalanzaConectada {
		t.Fatalf("integrado conectar: %+v", r)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		est := app.construirEstado()
		if est.PesoVivo > 0 && est.PuedeCapturar {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	final := app.construirEstado()
	if final.PesoVivo <= 0 {
		t.Fatal("integrado sin peso vivo")
	}

	cap := app.CapturarPeso()
	if !cap.OK || cap.Estado == nil || cap.Estado.PesoCapturado <= 0 {
		t.Fatalf("capturar integrado: %+v", cap)
	}
}
