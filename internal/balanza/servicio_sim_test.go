package balanza_test

import (
	"sync"
	"testing"
	"time"

	"sistema-pesaje/internal/balanza"
)

func TestServicioConFuenteSimuladaRecibePesoVivo(t *testing.T) {
	fuente := balanza.NuevaFuenteSimuladaSemilla(42)
	cfg := balanza.ConfigPorDefecto()
	svc := balanza.NuevoServicio(fuente, cfg)

	var mu sync.Mutex
	var callback int
	svc.AlPeso(func(p int) {
		mu.Lock()
		callback = p
		mu.Unlock()
	})

	if !svc.Conectar() {
		t.Fatal("conectar")
	}
	defer svc.Cerrar()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		p := callback
		mu.Unlock()
		if p >= balanza.PesoMinimoSim && p <= balanza.PesoMaximoSim {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	mu.Lock()
	p := callback
	mu.Unlock()
	if p < balanza.PesoMinimoSim || p > balanza.PesoMaximoSim {
		t.Fatalf("no llegó peso simulado; último callback=%d", p)
	}
	if svc.UltimoPeso() != p {
		t.Fatalf("UltimoPeso=%d callback=%d", svc.UltimoPeso(), p)
	}
	if !svc.EstaConectado() {
		t.Fatal("debe seguir conectado")
	}
}

func TestServicioSimuladaDesconectarLimpia(t *testing.T) {
	svc := balanza.NuevoServicio(balanza.NuevaFuenteSimulada(), balanza.ConfigPorDefecto())
	if !svc.Conectar() {
		t.Fatal("conectar")
	}
	svc.Cerrar()
	if svc.EstaConectado() {
		t.Fatal("debe quedar desconectado")
	}
	if svc.UltimoPeso() != 0 {
		t.Fatalf("peso residual=%d", svc.UltimoPeso())
	}
}
