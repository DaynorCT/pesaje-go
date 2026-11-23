package datos

import (
	"testing"

	"sistema-pesaje/internal/balanza"
)

func TestUsarBalanzaSimuladaPuertoSIM(t *testing.T) {
	cfg := balanza.ConfigPorDefecto()
	cfg.PuertoCOM = "SIM"
	if !UsarBalanzaSimulada(cfg) {
		t.Fatal("puerto SIM debe activar simulación")
	}
}

func TestUsarBalanzaSimuladaVariableEntorno(t *testing.T) {
	t.Setenv("PESAJE_BALANZA_SIM", "1")
	cfg := balanza.ConfigPorDefecto()
	cfg.PuertoCOM = "COM1"
	if !UsarBalanzaSimulada(cfg) {
		t.Fatal("PESAJE_BALANZA_SIM debe activar simulación")
	}
}

func TestUsarBalanzaSimuladaHardwarePorDefecto(t *testing.T) {
	t.Setenv("PESAJE_BALANZA_SIM", "")
	cfg := balanza.ConfigPorDefecto()
	if UsarBalanzaSimulada(cfg) {
		t.Fatal("COM1 sin env no debe forzar simulación")
	}
}
