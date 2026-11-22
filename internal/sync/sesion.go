package sync

import (
	"strings"
)

// Sesion mantiene el token de acceso contra el servidor web.
type Sesion struct {
	transporte    Transporte
	urlBase       string
	correo        string
	contrasena    string
	accessToken   string
	refreshToken  string
	ultimoError   string
}

func NuevaSesion(t Transporte, urlBase, correo, contrasena string) *Sesion {
	return &Sesion{
		transporte: t,
		urlBase:    strings.TrimRight(urlBase, "/"),
		correo:     correo,
		contrasena: contrasena,
	}
}

func (s *Sesion) URLBase() string { return s.urlBase }

func (s *Sesion) UltimoError() string { return s.ultimoError }

func (s *Sesion) TieneToken() bool {
	return s.accessToken != ""
}

func (s *Sesion) AsegurarToken() bool {
	if s.accessToken != "" {
		return true
	}
	return s.iniciarSesion()
}

func (s *Sesion) iniciarSesion() bool {
	cuerpo := escribirJSON(mapaValores("usuario", s.correo, "contrasena", s.contrasena))
	resp := s.transporte.Enviar("POST", s.urlBase+"/api/auth/token", cuerpo, nil)
	if !resp.Exito() {
		s.ultimoError = "No se pudo iniciar sesion en la web: " + resp.String()
		s.accessToken = ""
		s.refreshToken = ""
		return false
	}
	raiz, err := parseJSON(resp.Cuerpo)
	if err != nil || raiz == nil {
		s.ultimoError = "La web no devolvio un token de acceso."
		return false
	}
	token := texto(raiz, "access_token")
	if token == "" {
		s.ultimoError = "La web no devolvio un token de acceso."
		return false
	}
	s.accessToken = token
	s.refreshToken = texto(raiz, "refresh_token")
	s.ultimoError = ""
	return true
}

func (s *Sesion) renovar() bool {
	if s.refreshToken == "" {
		return s.iniciarSesion()
	}
	cab := map[string]string{"Cookie": "refresh_token=" + s.refreshToken}
	resp := s.transporte.Enviar("POST", s.urlBase+"/api/auth/refresh", "", cab)
	raiz, _ := parseJSON(resp.Cuerpo)
	token := texto(raiz, "access_token")
	if resp.Exito() && token != "" {
		s.accessToken = token
		s.ultimoError = ""
		return true
	}
	s.accessToken = ""
	return s.iniciarSesion()
}

func (s *Sesion) Invalidar() {
	s.accessToken = ""
}

func (s *Sesion) Cabeceras() map[string]string {
	h := make(map[string]string)
	if s.accessToken != "" {
		h["Cookie"] = "access_token=" + s.accessToken
	}
	return h
}