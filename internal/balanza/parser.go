package balanza

import "strconv"

const minimoDigitos = 2

// Parser interpreta tramas según la configuración (paridad ParserPeso Java).
type Parser struct {
	cfg Config
}

func NuevoParser(cfg Config) *Parser {
	return &Parser{cfg: cfg}
}

func (p *Parser) Parsear(trama string) (int, bool) {
	if p.cfg.MetodoLectura == MetodoPosicion {
		v, ok := p.parsearPorPosicion(trama)
		return v, ok
	}
	v, ok := p.parsearAutomatico(trama)
	return v, ok
}

func (p *Parser) PesoEnKilos(trama string) int {
	v, ok := p.Parsear(trama)
	if !ok {
		return 0
	}
	return v
}

func (p *Parser) parsearAutomatico(trama string) (int, bool) {
	indice := 0
	for indice < len(trama) {
		if !esDigito(trama[indice]) {
			indice++
			continue
		}
		fin := indice
		for fin < len(trama) && esDigito(trama[fin]) {
			fin++
		}
		if fin < len(trama) && trama[fin] == '.' && fin+1 < len(trama) && esDigito(trama[fin+1]) {
			fin++
			for fin < len(trama) && esDigito(trama[fin]) {
				fin++
			}
		}
		secuencia := trama[indice:fin]
		if cuentaDigitos(secuencia) >= minimoDigitos {
			return aKilos(secuencia)
		}
		indice = fin
	}
	return 0, false
}

func (p *Parser) parsearPorPosicion(trama string) (int, bool) {
	inicio := p.cfg.PosicionInicio
	longitud := p.cfg.PosicionLongitud
	if inicio < 1 || longitud < 1 {
		return 0, false
	}
	if inicio+longitud-1 > len(trama) {
		return 0, false
	}
	ventana := trama[inicio-1 : inicio-1+longitud]
	return aKilos(ventana)
}

func aKilos(token string) (int, bool) {
	valor := trimASCII(token)
	if valor == "" {
		return 0, false
	}
	punto := -1
	for i, c := range valor {
		if c == '.' {
			punto = i
			break
		}
	}
	var enteros, decimales string
	if punto >= 0 {
		enteros = valor[:punto]
		decimales = valor[punto+1:]
	} else {
		enteros = valor
	}
	if enteros == "" || !esSoloDigitos(enteros) {
		return 0, false
	}
	if decimales != "" && !esSoloDigitos(decimales) {
		return 0, false
	}
	n, err := strconv.Atoi(enteros)
	if err != nil {
		return 0, false
	}
	return n, true
}

func esDigito(b byte) bool { return b >= '0' && b <= '9' }

func esSoloDigitos(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !esDigito(s[i]) {
			return false
		}
	}
	return true
}

func cuentaDigitos(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if esDigito(s[i]) {
			n++
		}
	}
	return n
}

func trimASCII(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
