package auth

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"sistema-pesaje/internal/clave"
	"sistema-pesaje/internal/datos"
)

const consultaUsuario = `
SELECT u.id, u.persona_id, u.email, u.password_hash, u.rol, u.estado,
       p.nombre, p.apellido_paterno, p.apellido_materno
FROM usuarios u
INNER JOIN personas p ON p.id = u.persona_id
WHERE u.email = ?`

type filaUsuario struct {
	id              int64
	personaID       int64
	email           string
	passwordHash    string
	rol             string
	estado          string
	nombre          string
	apellidoPaterno sql.NullString
	apellidoMaterno sql.NullString
}

// Service autentica contra la tabla usuarios (paridad AuthService Java).
type Service struct {
	db *datos.Database
}

func NewService(db *datos.Database) *Service {
	return &Service{db: db}
}

func (s *Service) Autenticar(ctx context.Context, email, contrasena string) ResultadoAutenticacion {
	email = strings.TrimSpace(email)
	if email == "" || contrasena == "" {
		return fallo(EstadoEmailInvalido)
	}
	fila, ok, err := s.buscarUsuario(ctx, email)
	if err != nil {
		return fallo(EstadoError)
	}
	if !ok {
		return fallo(EstadoEmailInvalido)
	}
	if !strings.EqualFold(fila.estado, "ACTIVO") {
		return fallo(EstadoUsuarioInactivo)
	}
	if !clave.Verificar(contrasena, fila.passwordHash) {
		return fallo(EstadoContrasenaInvalida)
	}
	s.actualizarHashSiLegado(ctx, fila.id, contrasena, fila.passwordHash)
	s.registrarUltimoLogin(ctx, fila.id)
	return exito(fila.aSesion())
}

func (s *Service) VerificarContrasena(ctx context.Context, email, contrasena string) bool {
	email = strings.TrimSpace(email)
	if email == "" || contrasena == "" {
		return false
	}
	fila, ok, err := s.buscarUsuario(ctx, email)
	if err != nil || !ok {
		return false
	}
	if !strings.EqualFold(fila.estado, "ACTIVO") {
		return false
	}
	if !clave.Verificar(contrasena, fila.passwordHash) {
		return false
	}
	s.actualizarHashSiLegado(ctx, fila.id, contrasena, fila.passwordHash)
	return true
}

func (s *Service) buscarUsuario(ctx context.Context, email string) (filaUsuario, bool, error) {
	var f filaUsuario
	err := s.db.Conn().QueryRowContext(ctx, consultaUsuario, email).Scan(
		&f.id, &f.personaID, &f.email, &f.passwordHash, &f.rol, &f.estado,
		&f.nombre, &f.apellidoPaterno, &f.apellidoMaterno,
	)
	if err == sql.ErrNoRows {
		return filaUsuario{}, false, nil
	}
	if err != nil {
		return filaUsuario{}, false, err
	}
	return f, true, nil
}

func (f filaUsuario) aSesion() UsuarioSesion {
	nombre := f.nombre
	if f.apellidoPaterno.Valid && f.apellidoPaterno.String != "" {
		nombre += " " + f.apellidoPaterno.String
	}
	if f.apellidoMaterno.Valid && f.apellidoMaterno.String != "" {
		nombre += " " + f.apellidoMaterno.String
	}
	return UsuarioSesion{
		ID:        f.id,
		PersonaID: f.personaID,
		Email:     f.email,
		Nombre:    strings.TrimSpace(nombre),
		Rol:       f.rol,
	}
}

func (s *Service) actualizarHashSiLegado(ctx context.Context, usuarioID int64, contrasena, hashAlmacenado string) {
	if clave.EsHashModerno(hashAlmacenado) || !clave.EsHashLegadoSHA1(hashAlmacenado) {
		return
	}
	_, _ = s.db.Conn().ExecContext(ctx,
		"UPDATE usuarios SET password_hash = ? WHERE id = ?",
		clave.Hash(contrasena), usuarioID)
}

func (s *Service) registrarUltimoLogin(ctx context.Context, usuarioID int64) {
	_, _ = s.db.Conn().ExecContext(ctx,
		"UPDATE usuarios SET ultimo_login = ? WHERE id = ?",
		time.Now().Unix(), usuarioID)
}
