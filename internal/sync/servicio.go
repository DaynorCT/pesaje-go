package sync

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	"sistema-pesaje/internal/config"
	"sistema-pesaje/internal/datos"
)

const IntervaloAutomaticoMinutos = 2

// Servicio orquesta el motor y la cola.
type Servicio struct {
	motor         *Motor
	outbox        *Outbox
	mu            sync.Mutex
	ultimo        ResultadoSync
	ultimoMensaje string
	stopAuto      chan struct{}
	autoActivo    bool
}

func Crear(cfg config.Sistema, db *sql.DB, correo, contrasena string) *Servicio {
	if db == nil {
		return nil
	}
	t := NuevoTransporteHTTP(cfg.TimeoutHTTP)
	sesion := NuevaSesion(t, cfg.URL, correo, contrasena)
	outbox := NuevoOutbox(db)
	aplicador := NuevoAplicador(db)
	motor := NuevoMotor(t, sesion, outbox, aplicador)
	return &Servicio{motor: motor, outbox: outbox}
}

func (s *Servicio) EncolarPesaje(ctx context.Context, db *sql.DB, pesajeID int64) error {
	if s == nil || db == nil {
		return nil
	}
	p, err := datos.CargarPesajeParaSync(ctx, db, pesajeID)
	if err != nil {
		return err
	}
	clientEventID := p.SyncEventID
	if clientEventID == "" {
		clientEventID = uuid.NewString()
	}
	envio := pesajeParaEnvio(p)
	payload := PayloadPesaje(clientEventID, envio)
	return s.outbox.Encolar(clientEventID, "pesajes", int(pesajeID), payload)
}

func pesajeParaEnvio(p datos.PesajeParaSync) PesajeParaEnvio {
	var chofer, prov, prod, orig, dest *int64
	if p.ChoferID.Valid {
		v := p.ChoferID.Int64
		chofer = &v
	}
	if p.ProveedorID.Valid {
		v := p.ProveedorID.Int64
		prov = &v
	}
	if p.ProductoID.Valid {
		v := p.ProductoID.Int64
		prod = &v
	}
	if p.OrigenID.Valid {
		v := p.OrigenID.Int64
		orig = &v
	}
	if p.DestinoID.Valid {
		v := p.DestinoID.Int64
		dest = &v
	}
	costo := 0
	if p.CostoBs.Valid {
		costo = int(p.CostoBs.Int64)
	}
	flete := 0
	if p.FleteBsPendiente.Valid {
		flete = int(p.FleteBsPendiente.Int64)
	}
	return PesajeParaEnvio{
		ID:               p.ID,
		SyncEventID:      p.SyncEventID,
		Guia:             p.Guia,
		VehiculoID:       p.VehiculoID,
		ChoferID:         chofer,
		ProveedorID:      prov,
		ProductoID:       prod,
		OrigenID:         orig,
		DestinoID:        dest,
		PesoBruto:        p.PesoBruto,
		Tara:             p.Tara,
		CostoBs:          costo,
		FleteBsPendiente: flete,
		EstadoBalanza:    p.EstadoBalanza,
		FechaCreacion:    p.FechaCreacion,
	}
}

func (s *Servicio) EnviarPesoVivo(peso int) bool {
	if s == nil {
		return false
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[sync] peso vivo: %v", r)
		}
	}()
	ok := s.motor.EnviarPesoVivo(peso)
	if !ok {
		s.mu.Lock()
		s.ultimoMensaje = "No se pudo publicar el peso en vivo en la web."
		s.mu.Unlock()
	}
	return ok
}

func (s *Servicio) SincronizarAhora() ResultadoSync {
	if s == nil {
		return ResultadoSync{Mensaje: "Servicio no disponible."}
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[sync] vuelta: %v", r)
		}
	}()
	res := s.motor.Sincronizar()
	s.mu.Lock()
	s.ultimo = res
	s.ultimoMensaje = res.Mensaje
	s.mu.Unlock()
	return res
}

func (s *Servicio) IniciarAutomatico(minutos int) {
	if s == nil {
		return
	}
	s.DetenerAutomatico()
	if minutos < 1 {
		minutos = IntervaloAutomaticoMinutos
	}
	s.stopAuto = make(chan struct{})
	s.autoActivo = true
	go func() {
		select {
		case <-time.After(30 * time.Minute):
		case <-s.stopAuto:
			return
		}
		ticker := time.NewTicker(time.Duration(minutos) * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.SincronizarAhora()
			case <-s.stopAuto:
				return
			}
		}
	}()
}

func (s *Servicio) DetenerAutomatico() {
	if s == nil {
		return
	}
	if s.stopAuto != nil {
		close(s.stopAuto)
		s.stopAuto = nil
	}
	s.autoActivo = false
}

func (s *Servicio) AutomaticoActivo() bool {
	return s != nil && s.autoActivo
}

func (s *Servicio) Pendientes() int64 {
	if s == nil {
		return 0
	}
	n, _ := s.outbox.PendientesCount()
	return n
}

func (s *Servicio) Errores() int64 {
	if s == nil {
		return 0
	}
	n, _ := s.outbox.ErroresCount()
	return n
}

func (s *Servicio) UltimoResultado() ResultadoSync {
	if s == nil {
		return ResultadoSync{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ultimo
}

func (s *Servicio) UltimoMensaje() string {
	if s == nil {
		return ""
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ultimoMensaje
}

func (s *Servicio) Cerrar() {
	s.DetenerAutomatico()
}