package balanza

import (
	"strings"
	"time"

	"go.bug.st/serial"
)

// LectorSerie lee el puerto serie (paridad LectorSerie Java).
type LectorSerie struct {
	cfg   Config
	port  serial.Port
	ok    bool
	err   string
	buf   strings.Builder
}

func NuevoLectorSerie(cfg Config) *LectorSerie {
	return &LectorSerie{cfg: cfg}
}

func (l *LectorSerie) Conectar() bool {
	if l.ok {
		return true
	}
	mode := &serial.Mode{
		BaudRate: l.cfg.BaudRate,
		DataBits: l.cfg.DataBits,
		StopBits: stopBits(l.cfg.StopBits),
		Parity:   paridad(l.cfg.Paridad),
	}
	p, err := serial.Open(l.cfg.PuertoCOM, mode)
	if err != nil {
		l.err = err.Error()
		return false
	}
	_ = p.SetReadTimeout(time.Duration(l.cfg.TimeoutMs) * time.Millisecond)
	l.port = p
	l.ok = true
	l.err = ""
	return true
}

func (l *LectorSerie) Desconectar() {
	if l.port != nil {
		_ = l.port.Close()
		l.port = nil
	}
	l.ok = false
	l.buf.Reset()
}

func (l *LectorSerie) EstaConectado() bool { return l.ok }

func (l *LectorSerie) UltimoError() string { return l.err }

func (l *LectorSerie) LeerTrama() string {
	if !l.ok || l.port == nil {
		return ""
	}
	b := make([]byte, 256)
	n, err := l.port.Read(b)
	if err != nil || n == 0 {
		return ""
	}
	l.buf.Write(b[:n])
	s := l.buf.String()
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		linea := strings.TrimSpace(s[:i])
		resto := s[i+1:]
		if strings.HasPrefix(resto, "\n") && len(resto) > 0 {
			resto = resto[1:]
		}
		l.buf.Reset()
		l.buf.WriteString(resto)
		return linea
	}
	return ""
}

func paridad(p string) serial.Parity {
	switch strings.ToUpper(strings.TrimSpace(p)) {
	case "E":
		return serial.EvenParity
	case "O":
		return serial.OddParity
	default:
		return serial.NoParity
	}
}

func stopBits(n int) serial.StopBits {
	switch n {
	case 2:
		return serial.TwoStopBits
	case 3:
		return serial.OnePointFiveStopBits
	default:
		return serial.OneStopBit
	}
}
