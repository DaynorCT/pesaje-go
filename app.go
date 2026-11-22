package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"sistema-pesaje/internal/auth"
	"sistema-pesaje/internal/balanza"
	"sistema-pesaje/internal/config"
	"sistema-pesaje/internal/datos"
	pesync "sistema-pesaje/internal/sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App expone bindings Wails al frontend.
type App struct {
	ctx  context.Context
	auth *auth.Service
	db   *datos.Database

	muIntegrado sync.Mutex
	integrado   estadoIntegrado
	muPesaje    sync.Mutex
	pesajeUI    estadoPesajeUI
	sesion      *UsuarioDTO
	balanzaSvc  *balanza.Servicio
	muSync      sync.Mutex
	syncSvc     *pesync.Servicio
	claveWeb    string
}

type estadoPesajeUI struct {
	pesoBruto  int
	pesoTara   int
	mensaje    string
	esError    bool
}

type estadoIntegrado struct {
	conectado     bool
	pesoVivo      int
	pesoCapturado int
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	base, err := datos.ObtenerInstancia()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[SistemaPesaje] No se pudo abrir la base: %v\n", err)
		return
	}
	a.db = base
	a.auth = auth.NewService(base)
	runtime.WindowSetSize(ctx, 460, 520)
}

func (a *App) ProductVersion() string {
	return "0.1.0"
}

type LoginRespuesta struct {
	Exito   bool        `json:"exito"`
	Estado  string      `json:"estado"`
	Mensaje string      `json:"mensaje"`
	Usuario *UsuarioDTO `json:"usuario,omitempty"`
}

type UsuarioDTO struct {
	ID        int64  `json:"id"`
	PersonaID int64  `json:"personaId"`
	Email     string `json:"email"`
	Nombre    string `json:"nombre"`
	Rol       string `json:"rol"`
}

type EstadoIntegradoDTO struct {
	BalanzaConectada bool   `json:"balanzaConectada"`
	BalanzaTexto     string `json:"balanzaTexto"`
	SyncTexto        string `json:"syncTexto"`
	SyncVariante     string `json:"syncVariante"`
	PesoVivo         int    `json:"pesoVivo"`
	PesoCapturado    int    `json:"pesoCapturado"`
	TextoConectar    string `json:"textoConectar"`
	PuedeCapturar    bool   `json:"puedeCapturar"`
	PuedeEnviar      bool   `json:"puedeEnviar"`
	SyncDisponible   bool   `json:"syncDisponible"`
}

type AvisoDTO struct {
	Titulo  string `json:"titulo"`
	Mensaje string `json:"mensaje"`
	EsError bool   `json:"esError"`
}

type AccionDTO struct {
	OK           bool                  `json:"ok"`
	Aviso        *AvisoDTO             `json:"aviso,omitempty"`
	Estado       *EstadoIntegradoDTO   `json:"estado,omitempty"`
	EstadoPesaje *EstadoPesajePaneDTO  `json:"estadoPesaje,omitempty"`
}

func (a *App) AvisoLoginInicial() string {
	if a.db == nil {
		return ""
	}
	clave := a.db.ClaveAdministrador()
	if clave == "" {
		return ""
	}
	return fmt.Sprintf(
		"Primer inicio: usuario %s con la contraseña %s. Cámbiela después de entrar.",
		 datos.EmailAdministrador, clave,
	)
}

func (a *App) Autenticar(email, contrasena string) LoginRespuesta {
	if a.auth == nil {
		return LoginRespuesta{
			Exito:   false,
			Estado:  string(auth.EstadoError),
			Mensaje: "No se pudo iniciar sesión.",
		}
	}
	res := a.auth.Autenticar(a.ctx, email, contrasena)
	if !res.EsExito() {
		return LoginRespuesta{
			Exito:   false,
			Estado:  string(res.Estado),
			Mensaje: mensajeAutenticacion(res.Estado),
		}
	}
	u := res.Usuario
	dto := &UsuarioDTO{ID: u.ID, PersonaID: u.PersonaID, Email: u.Email, Nombre: u.Nombre, Rol: u.Rol}
	a.sesion = dto
	a.muIntegrado.Lock()
	a.integrado = estadoIntegrado{}
	a.muIntegrado.Unlock()

	a.prepararSync(email, contrasena)
	a.activarVentanaIntegrado()
	a.syncAlMostrar()

	return LoginRespuesta{
		Exito:   true,
		Estado:  string(auth.EstadoExito),
		Usuario: dto,
	}
}

func (a *App) EstadoIntegrado() EstadoIntegradoDTO {
	return a.construirEstado()
}

func (a *App) construirEstado() EstadoIntegradoDTO {
	a.muIntegrado.Lock()
	defer a.muIntegrado.Unlock()
	cfg := config.Cargar()
	syncDisp := strings.TrimSpace(cfg.URL) != ""
	textoSync, syncVar := a.textoSyncIntegrado()
	if !syncDisp {
		textoSync = "Sin web"
		syncVar = "proceso"
	}
	balanzaTxt := "desconectada"
	textoBtn := "Conectar balanza"
	if a.integrado.conectado {
		balanzaTxt = "conectada"
		textoBtn = "Desconectar balanza"
	}
	return EstadoIntegradoDTO{
		BalanzaConectada: a.integrado.conectado,
		BalanzaTexto:     balanzaTxt,
		SyncTexto:        textoSync,
		SyncVariante:     syncVar,
		PesoVivo:         a.integrado.pesoVivo,
		PesoCapturado:    a.integrado.pesoCapturado,
		TextoConectar:    textoBtn,
		PuedeCapturar:    a.integrado.conectado && a.integrado.pesoVivo > 0,
		PuedeEnviar:      a.integrado.pesoCapturado > 0,
		SyncDisponible:   syncDisp,
	}
}

func (a *App) AlternarBalanza() AccionDTO {
	a.muIntegrado.Lock()
	conectado := a.integrado.conectado
	a.muIntegrado.Unlock()
	if conectado {
		a.detenerBalanza()
		a.muIntegrado.Lock()
		a.integrado.conectado = false
		a.integrado.pesoVivo = 0
		a.integrado.pesoCapturado = 0
		est := a.construirEstadoUnlocked()
		a.muIntegrado.Unlock()
		return AccionDTO{OK: true, Estado: &est}
	}
	if a.db == nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Balanza", Mensaje: "Base de datos no disponible.", EsError: true}}
	}
	cfg := a.db.ObtenerConfigBalanza()
	if !a.iniciarBalanza(cfg) {
		msg := "No se pudo abrir el puerto " + cfg.PuertoCOM + "."
		if a.balanzaSvc != nil {
			if e := a.balanzaSvc.UltimoError(); e != "" {
				msg = e
			}
		}
		if datos.UsarBalanzaSimulada(cfg) {
			msg = "No se pudo iniciar la báscula simulada."
		}
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Balanza", Mensaje: msg}}
	}
	est := a.construirEstado()
	return AccionDTO{OK: true, Estado: &est}
}

func (a *App) construirEstadoUnlocked() EstadoIntegradoDTO {
	cfg := config.Cargar()
	syncDisp := strings.TrimSpace(cfg.URL) != ""
	textoSync, syncVar := a.textoSyncIntegrado()
	if !syncDisp {
		textoSync = "Sin web"
		syncVar = "proceso"
	}
	balanzaTxt := "desconectada"
	textoBtn := "Conectar balanza"
	if a.integrado.conectado {
		balanzaTxt = "conectada"
		textoBtn = "Desconectar balanza"
	}
	return EstadoIntegradoDTO{
		BalanzaConectada: a.integrado.conectado,
		BalanzaTexto:     balanzaTxt,
		SyncTexto:        textoSync,
		SyncVariante:     syncVar,
		PesoVivo:         a.integrado.pesoVivo,
		PesoCapturado:    a.integrado.pesoCapturado,
		TextoConectar:    textoBtn,
		PuedeCapturar:    a.integrado.conectado && a.integrado.pesoVivo > 0,
		PuedeEnviar:      a.integrado.pesoCapturado > 0,
		SyncDisponible:   syncDisp,
	}
}

func (a *App) CapturarPeso() AccionDTO {
	a.muIntegrado.Lock()
	defer a.muIntegrado.Unlock()
	if !a.integrado.conectado {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Balanza", Mensaje: "Conecte la balanza primero."}}
	}
	if a.integrado.pesoVivo <= 0 {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Peso", Mensaje: "Peso inválido."}}
	}
	a.integrado.pesoCapturado = a.integrado.pesoVivo
	est := a.construirEstadoUnlocked()
	return AccionDTO{OK: true, Estado: &est}
}

func (a *App) EnviarPeso() AccionDTO {
	a.muIntegrado.Lock()
	peso := a.integrado.pesoCapturado
	a.muIntegrado.Unlock()
	if peso <= 0 {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Peso", Mensaje: "Capture el peso primero."}}
	}
	cfg := config.Cargar()
	if strings.TrimSpace(cfg.URL) == "" {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Sincronización", Mensaje: "Servicio de sincronización no disponible."}}
	}
	a.muSync.Lock()
	svc := a.syncSvc
	a.muSync.Unlock()
	if svc == nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Sincronización", Mensaje: "Servicio de sincronización no disponible."}}
	}
	ok := svc.EnviarPesoVivo(peso)
	est := a.construirEstado()
	if ok {
		return AccionDTO{
			OK:    true,
			Aviso: &AvisoDTO{Titulo: "Envío", Mensaje: "Peso " + formatearPeso(peso) + " enviado a la web."},
			Estado: ptrEstado(est),
		}
	}
	return AccionDTO{
		OK: false,
		Aviso: &AvisoDTO{Titulo: "Envío", Mensaje: "No se pudo enviar el peso. Revise la conexión con el sistema web."},
		Estado: ptrEstado(est),
	}
}

func (a *App) SincronizarAhora() AccionDTO {
	cfg := config.Cargar()
	if strings.TrimSpace(cfg.URL) == "" {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Sincronización", Mensaje: "La sincronización no está configurada."}}
	}
	a.muSync.Lock()
	svc := a.syncSvc
	a.muSync.Unlock()
	if svc == nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Sincronización", Mensaje: "Inicie sesión para sincronizar."}}
	}
	r := svc.SincronizarAhora()
	est := a.construirEstado()
	dto := AccionDTO{OK: true, Estado: ptrEstado(est)}
	if r.Mensaje != "" {
		dto.Aviso = &AvisoDTO{Titulo: "Sincronización", Mensaje: r.Mensaje}
	}
	return dto
}

func (a *App) SolicitarSistemaCompleto(contrasena string) AccionDTO {
	if a.sesion == nil || a.auth == nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Sesión", Mensaje: "No hay sesión activa.", EsError: true}}
	}
	if !a.auth.VerificarContrasena(a.ctx, a.sesion.Email, contrasena) {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Acceso", Mensaje: "Contraseña incorrecta.", EsError: true}}
	}
	a.detenerBalanza()
	a.muIntegrado.Lock()
	a.integrado = estadoIntegrado{}
	a.muIntegrado.Unlock()
	a.activarVentanaShell()
	return AccionDTO{OK: true}
}

func (a *App) CerrarSesion() {
	a.detenerBalanza()
	a.muIntegrado.Lock()
	a.integrado = estadoIntegrado{}
	a.muIntegrado.Unlock()
	a.sesion = nil
	a.cerrarSync()
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
	runtime.WindowSetMaxSize(a.ctx, 520, 640)
	runtime.WindowSetMinSize(a.ctx, 420, 480)
	runtime.WindowSetSize(a.ctx, 460, 520)
}

func ptrEstado(e EstadoIntegradoDTO) *EstadoIntegradoDTO {
	x := e
	return &x
}

func formatearPeso(peso int) string {
	return fmt.Sprintf("%d kg", peso)
}

func mensajeAutenticacion(estado auth.EstadoAutenticacion) string {
	switch estado {
	case auth.EstadoEmailInvalido:
		return "El correo no existe en la base local."
	case auth.EstadoContrasenaInvalida:
		return "La contraseña no es correcta."
	case auth.EstadoUsuarioInactivo:
		return "El usuario está inactivo. Consulte al administrador."
	default:
		return "No se pudo iniciar sesión."
	}
}
