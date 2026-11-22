package clave

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"crypto/sha256"
)

const (
	Prefijo       = "pbkdf2"
	PrefijoSHA1   = "sha1"
	Iteraciones   = 120_000
	LongitudSalt  = 16
	LongitudHash  = 32
)

// Hash genera un hash PBKDF2-HMAC-SHA256 en el formato del módulo Java.
func Hash(contrasena string) string {
	salt := make([]byte, LongitudSalt)
	_, _ = rand.Read(salt)
	derivado := derivar(contrasena, salt, Iteraciones)
	return fmt.Sprintf("%s:%d:%s:%s",
		Prefijo,
		Iteraciones,
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(derivado),
	)
}

func EsHashModerno(almacenado string) bool {
	return almacenado != "" && strings.HasPrefix(almacenado, Prefijo+":")
}

func EsHashLegadoSHA1(almacenado string) bool {
	return almacenado != "" && strings.HasPrefix(almacenado, PrefijoSHA1+":")
}

// Verificar comprueba la contraseña contra hash PBKDF2 o SHA1 legado.
func Verificar(contrasena, almacenado string) bool {
	if contrasena == "" || almacenado == "" {
		return false
	}
	if EsHashModerno(almacenado) {
		return verificarPBKDF2(contrasena, almacenado)
	}
	if EsHashLegadoSHA1(almacenado) {
		return verificarSHA1Legado(contrasena, almacenado)
	}
	return false
}

func verificarPBKDF2(contrasena, almacenado string) bool {
	partes := strings.Split(almacenado, ":")
	if len(partes) != 4 || partes[0] != Prefijo {
		return false
	}
	iteraciones, err := strconv.Atoi(partes[1])
	if err != nil || iteraciones <= 0 {
		return false
	}
	salt, err := base64.StdEncoding.DecodeString(partes[2])
	if err != nil {
		return false
	}
	esperado, err := base64.StdEncoding.DecodeString(partes[3])
	if err != nil {
		return false
	}
	derivado := derivar(contrasena, salt, iteraciones)
	return subtle.ConstantTimeCompare(derivado, esperado) == 1
}

func verificarSHA1Legado(contrasena, almacenado string) bool {
	cuerpo := almacenado[len(PrefijoSHA1)+1:]
	partes := strings.SplitN(cuerpo, ":", 2)
	if len(partes) != 2 {
		return false
	}
	hexBytes, err := base64.StdEncoding.DecodeString(partes[0])
	if err != nil {
		return false
	}
	saltBytes, err := base64.StdEncoding.DecodeString(partes[1])
	if err != nil {
		return false
	}
	hexEsperado := string(hexBytes)
	salt := string(saltBytes)
	hexCalculado := sha1Hex(salt + contrasena)
	return subtle.ConstantTimeCompare([]byte(hexEsperado), []byte(hexCalculado)) == 1
}

func sha1Hex(texto string) string {
	h := sha1.Sum([]byte(texto))
	var b strings.Builder
	b.Grow(len(h) * 2)
	for _, by := range h {
		b.WriteString(fmt.Sprintf("%02x", by))
	}
	return b.String()
}

func derivar(contrasena string, salt []byte, iteraciones int) []byte {
	return pbkdf2.Key([]byte(contrasena), salt, iteraciones, LongitudHash, sha256.New)
}
