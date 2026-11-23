package balanza

import (
	"strings"
	"testing"
)

func TestFuenteSimuladaConectaYEmiteTramasParseables(t *testing.T) {
	f := NuevaFuenteSimuladaSemilla(99)
	if !f.Conectar() {
		t.Fatal("conectar")
	}
	defer f.Desconectar()

	trama := f.LeerTrama()
	if trama == "" {
		t.Fatal("trama vacía")
	}
	if !strings.HasPrefix(trama, "ST:") {
		t.Fatalf("trama inesperada: %q", trama)
	}
	p := f.UltimoPeso()
	if p < PesoMinimoSim || p > PesoMaximoSim {
		t.Fatalf("peso %d fuera de [%d,%d]", p, PesoMinimoSim, PesoMaximoSim)
	}
	v, ok := NuevoParser(ConfigPorDefecto()).Parsear(trama)
	if !ok || v != p {
		t.Fatalf("parser: got %d %v want %d", v, ok, p)
	}
}

func TestFuenteSimuladaDesconectadaNoLee(t *testing.T) {
	f := NuevaFuenteSimulada()
	if f.LeerTrama() != "" {
		t.Fatal("sin conectar no debe haber trama")
	}
}

func TestFuenteSimuladaSemillaDeterministica(t *testing.T) {
	a := NuevaFuenteSimuladaSemilla(7)
	b := NuevaFuenteSimuladaSemilla(7)
	_ = a.Conectar()
	_ = b.Conectar()
	if a.LeerTrama() != b.LeerTrama() {
		t.Fatal("misma semilla debe dar la misma trama")
	}
}
