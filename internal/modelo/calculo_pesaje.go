package modelo

import "errors"

// ErrCalculoPesaje agrupa validaciones de bruto/tara (paridad CalculoPesaje Java).
var (
	ErrBrutoInvalido     = errors.New("El peso bruto debe ser mayor a cero.")
	ErrTaraNegativa      = errors.New("La tara no puede ser negativa.")
	ErrTaraMayorQueBruto = errors.New("La tara no puede ser mayor al peso bruto.")
)

// Neto calcula peso neto en kilos.
func Neto(bruto, tara int) (int, error) {
	if err := Validar(bruto, tara); err != nil {
		return 0, err
	}
	return bruto - tara, nil
}

// Validar comprueba la pareja bruto/tara.
func Validar(bruto, tara int) error {
	if bruto <= 0 {
		return ErrBrutoInvalido
	}
	if tara < 0 {
		return ErrTaraNegativa
	}
	if tara > bruto {
		return ErrTaraMayorQueBruto
	}
	return nil
}