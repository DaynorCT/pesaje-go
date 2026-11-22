package balanza

import (
	"fmt"
	"math/rand/v2"
)

const (
	PesoMinimoSim = 1000
	PesoMaximoSim = 5000
)

// FuenteSimulada para pruebas y desarrollo (no usar en producción con hardware real).
type FuenteSimulada struct {
	rng       *rand.Rand
	conectada bool
	ultimo    int
	semilla   uint64
}

func NuevaFuenteSimulada() *FuenteSimulada {
	return &FuenteSimulada{rng: rand.New(rand.NewPCG(1, 2))}
}

func NuevaFuenteSimuladaSemilla(semilla uint64) *FuenteSimulada {
	return &FuenteSimulada{rng: rand.New(rand.NewPCG(semilla, semilla^0x9e3779b97f4a7c15)), semilla: semilla}
}

func (f *FuenteSimulada) Conectar() bool {
	f.conectada = true
	return true
}

func (f *FuenteSimulada) Desconectar() { f.conectada = false }

func (f *FuenteSimulada) EstaConectado() bool { return f.conectada }

func (f *FuenteSimulada) UltimoError() string { return "" }

func (f *FuenteSimulada) LeerTrama() string {
	if !f.conectada {
		return ""
	}
	f.ultimo = PesoMinimoSim + f.rng.IntN(PesoMaximoSim-PesoMinimoSim+1)
	return fmt.Sprintf("ST: %d kg", f.ultimo)
}

func (f *FuenteSimulada) UltimoPeso() int { return f.ultimo }
