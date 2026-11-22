package main

import (
	"context"
	"fmt"
	"os"

	"sistema-pesaje/internal/auth"
	"sistema-pesaje/internal/datos"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App expone bindings Wails al frontend.
type App struct {
	ctx  context.Context
	auth *auth.Service
	db   *datos.Database
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

// ProductVersion es la versión mostrada en login.
func (a *App) ProductVersion() string {
	return "0.1.0"
}

// LoginRespuesta es el resultado del intento de inicio de sesión.
type LoginRespuesta struct {
	Exito   bool        `json:"exito"`
	Estado  string      `json:"estado"`
	Mensaje string      `json:"mensaje"`
	Usuario *UsuarioDTO `json:"usuario,omitempty"`
}

// UsuarioDTO datos de sesión para el frontend.
type UsuarioDTO struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

// AvisoLoginInicial aviso de contraseña generada en el primer arranque (misma sesión de base).
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

// Autenticar inicia sesión local (paridad AuthService Java).
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
	runtime.WindowSetMinSize(a.ctx, 420, 600)
	runtime.WindowSetMaxSize(a.ctx, 0, 0)
	runtime.WindowSetSize(a.ctx, 480, 680)
	u := res.Usuario
	return LoginRespuesta{
		Exito:  true,
		Estado: string(auth.EstadoExito),
		Usuario: &UsuarioDTO{
			ID:     u.ID,
			Email:  u.Email,
			Nombre: u.Nombre,
			Rol:    u.Rol,
		},
	}
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
