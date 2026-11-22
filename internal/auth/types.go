package auth

import "strings"

// EstadoAutenticacion es el resultado de un intento de inicio de sesión.
type EstadoAutenticacion string

const (
	EstadoExito              EstadoAutenticacion = "EXITO"
	EstadoEmailInvalido      EstadoAutenticacion = "EMAIL_INVALIDO"
	EstadoContrasenaInvalida EstadoAutenticacion = "CONTRASENA_INVALIDA"
	EstadoUsuarioInactivo    EstadoAutenticacion = "USUARIO_INACTIVO"
	EstadoError              EstadoAutenticacion = "ERROR"
)

// UsuarioSesion representa un usuario con sesión iniciada.
type UsuarioSesion struct {
	ID        int64
	PersonaID int64
	Email     string
	Nombre    string
	Rol       string
}

func (u UsuarioSesion) EsAdministrador() bool {
	return strings.EqualFold(u.Rol, "administrador")
}

// ResultadoAutenticacion agrupa estado y usuario (si hubo éxito).
type ResultadoAutenticacion struct {
	Estado  EstadoAutenticacion
	Usuario *UsuarioSesion
}

func (r ResultadoAutenticacion) EsExito() bool {
	return r.Estado == EstadoExito
}

func exito(usuario UsuarioSesion) ResultadoAutenticacion {
	u := usuario
	return ResultadoAutenticacion{Estado: EstadoExito, Usuario: &u}
}

func fallo(estado EstadoAutenticacion) ResultadoAutenticacion {
	return ResultadoAutenticacion{Estado: estado, Usuario: nil}
}
