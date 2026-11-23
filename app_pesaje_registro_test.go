package main

import (
	"context"
	"strings"
	"testing"

	"sistema-pesaje/internal/datos"
)

const placaPesajeE2E = "E2E-SIM-001"

func sesionAdministrador(t *testing.T, app *App) {
	t.Helper()
	if app.db == nil {
		t.Fatal("base no disponible")
	}
	var id, personaID int64
	var email, nombre, rol string
	err := app.db.Conn().QueryRowContext(context.Background(), `
SELECT u.id, u.persona_id, u.email, u.rol, p.nombre
FROM usuarios u INNER JOIN personas p ON p.id = u.persona_id
WHERE u.email = ?`, datos.EmailAdministrador).
		Scan(&id, &personaID, &email, &rol, &nombre)
	if err != nil {
		t.Fatalf("usuario admin: %v", err)
	}
	app.sesion = &UsuarioDTO{
		ID:        id,
		PersonaID: personaID,
		Email:     email,
		Nombre:    nombre,
		Rol:       rol,
	}
}

func vehiculoIDPorPlaca(t *testing.T, app *App, placa string) int64 {
	t.Helper()
	opc, err := app.PesajeOpciones()
	if err != nil {
		t.Fatalf("opciones: %v", err)
	}
	for _, v := range opc.Vehiculos {
		if strings.EqualFold(v.Label, placa) {
			return v.ID
		}
	}
	t.Fatalf("vehículo %q no encontrado en opciones", placa)
	return 0
}

func TestPesajeRegistrarConBalanzaSimuladaYPersisteEnSQLite(t *testing.T) {
	app, limpiar := prepararAppConBalanzaSim(t)
	defer limpiar()
	sesionAdministrador(t, app)

	if r := app.PesajeAlternarBalanza(); !r.OK {
		t.Fatalf("conectar báscula: %+v", r.Aviso)
	}
	bruto := esperarPesoVivo(t, app)
	const tara = 500
	if bruto <= tara {
		t.Fatalf("bruto simulado=%d debe ser > tara=%d", bruto, tara)
	}

	if cr := app.PesajeCrearVehiculo(placaPesajeE2E, "Camión", tara); !cr.OK {
		t.Fatalf("crear vehículo: %+v", cr.Aviso)
	}
	vehID := vehiculoIDPorPlaca(t, app, placaPesajeE2E)

	reg := app.PesajeRegistrar(RegistrarPesajeDTO{
		VehiculoID: vehID,
		PesoBruto:  bruto,
		Tara:       tara,
	})
	if !reg.OK {
		if reg.EstadoPesaje != nil && reg.EstadoPesaje.Mensaje != "" {
			t.Fatalf("registrar: %s", reg.EstadoPesaje.Mensaje)
		}
		t.Fatalf("registrar: %+v", reg.Aviso)
	}

	filas, err := app.PesajeListarGrilla(20)
	if err != nil {
		t.Fatalf("grilla: %v", err)
	}
	if len(filas) != 1 {
		t.Fatalf("filas=%d want 1", len(filas))
	}
	f := filas[0]
	if f.Placa != placaPesajeE2E {
		t.Fatalf("placa=%q", f.Placa)
	}
	if f.PesoBruto != bruto || f.Tara != tara || f.PesoNeto != bruto-tara {
		t.Fatalf("pesos bruto=%d tara=%d neto=%d", f.PesoBruto, f.Tara, f.PesoNeto)
	}
	if f.EstadoBalanza != "EN_PROCESO" {
		t.Fatalf("estado_balanza=%q", f.EstadoBalanza)
	}

	var n int
	if err := app.db.Conn().QueryRow(`SELECT COUNT(*) FROM pesajes`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("pesajes en DB=%d", n)
	}
}

func TestPesajeRegistrarTrasCapturaBalanzaSimulada(t *testing.T) {
	app, limpiar := prepararAppConBalanzaSim(t)
	defer limpiar()
	sesionAdministrador(t, app)

	if r := app.PesajeAlternarBalanza(); !r.OK {
		t.Fatalf("conectar: %+v", r.Aviso)
	}
	_ = esperarPesoVivo(t, app)

	captura := app.PesajeCapturarBruto()
	if !captura.OK || captura.EstadoPesaje == nil {
		t.Fatal("capturar bruto")
	}
	bruto := captura.EstadoPesaje.PesoBruto
	tara := bruto / 5
	if tara < 100 {
		tara = 100
	}
	app.PesajeActualizarTara(tara)

	if cr := app.PesajeCrearVehiculo("E2E-CAP-01", "Volquete", tara); !cr.OK {
		t.Fatalf("crear vehículo: %+v", cr.Aviso)
	}
	vehID := vehiculoIDPorPlaca(t, app, "E2E-CAP-01")

	est := app.PesajeEstado()
	reg := app.PesajeRegistrar(RegistrarPesajeDTO{
		VehiculoID: vehID,
		PesoBruto:  est.PesoBruto,
		Tara:       est.PesoTara,
	})
	if !reg.OK {
		t.Fatalf("registrar: %+v", reg)
	}
	if reg.EstadoPesaje == nil || reg.EstadoPesaje.Mensaje == "" {
		t.Fatal("mensaje de éxito esperado")
	}
	if reg.EstadoPesaje.PesoBruto != 0 || reg.EstadoPesaje.PesoTara != 0 {
		t.Fatal("UI debe limpiar bruto/tara tras guardar")
	}

	filas, _ := app.PesajeListarGrilla(5)
	if len(filas) != 1 || filas[0].PesoNeto != est.PesoNeto {
		t.Fatalf("grilla %+v neto pane %d", filas, est.PesoNeto)
	}
}
