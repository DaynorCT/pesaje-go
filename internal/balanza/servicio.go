package balanza

import (
	"context"
	"sync"
	"time"
)

const IntervaloLectura = 200 * time.Millisecond

// Servicio lee peso en vivo en un goroutine (paridad ServicioBalanza Java).
type Servicio struct {
	lector LectorPeso
	parser *Parser

	mu         sync.Mutex
	conectado  bool
	ultimoPeso int
	ultimoErr  string

	cancel context.CancelFunc
	wg     sync.WaitGroup

	onPeso   func(int)
	onEstado func(conectado bool, msg string)
}

func NuevoServicio(lector LectorPeso, cfg Config) *Servicio {
	return &Servicio{
		lector: lector,
		parser: NuevoParser(cfg),
	}
}

func (s *Servicio) AlPeso(cb func(int)) { s.onPeso = cb }

func (s *Servicio) AlEstado(cb func(bool, string)) { s.onEstado = cb }

func (s *Servicio) Conectar() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conectado {
		return true
	}
	if !s.lector.Conectar() {
		s.ultimoErr = s.lector.UltimoError()
		s.avisarEstado(false, s.ultimoErr)
		return false
	}
	s.conectado = true
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.wg.Add(1)
	go s.bucle(ctx)
	s.avisarEstado(true, "conectada")
	return true
}

func (s *Servicio) Desconectar() {
	s.mu.Lock()
	if !s.conectado {
		s.mu.Unlock()
		return
	}
	s.conectado = false
	if s.cancel != nil {
		s.cancel()
	}
	s.mu.Unlock()
	s.wg.Wait()
	s.lector.Desconectar()
	s.ultimoPeso = 0
	s.avisarEstado(false, "desconectada")
}

func (s *Servicio) EstaConectado() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conectado
}

func (s *Servicio) UltimoPeso() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ultimoPeso
}

func (s *Servicio) UltimoError() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ultimoErr
}

func (s *Servicio) Cerrar() { s.Desconectar() }

func (s *Servicio) bucle(ctx context.Context) {
	defer s.wg.Done()
	t := time.NewTicker(IntervaloLectura)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			trama := s.lector.LeerTrama()
			if trama == "" {
				continue
			}
			if peso, ok := s.parser.Parsear(trama); ok {
				s.mu.Lock()
				s.ultimoPeso = peso
				s.mu.Unlock()
				if s.onPeso != nil {
					s.onPeso(peso)
				}
			}
		}
	}
}

func (s *Servicio) avisarEstado(conectado bool, msg string) {
	if s.onEstado != nil {
		s.onEstado(conectado, msg)
	}
}
