package balanza

import "strings"

type MetodoLectura int

const (
	MetodoAuto MetodoLectura = iota
	MetodoPosicion
)

func MetodoDesdeTexto(texto string) MetodoLectura {
	t := strings.TrimSpace(strings.ToUpper(texto))
	if t == "POSICION" {
		return MetodoPosicion
	}
	return MetodoAuto
}

func (m MetodoLectura) Texto() string {
	if m == MetodoPosicion {
		return "POSICION"
	}
	return "AUTO"
}
