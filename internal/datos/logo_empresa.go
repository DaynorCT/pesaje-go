package datos

import (
	"encoding/base64"
	"strings"

	"sistema-pesaje/assets"
)

const (
	LogoLimiteBytes   = 2 * 1024 * 1024
	marcaBase64       = ";base64,"
	prefijoPNG        = "data:image/png;base64,"
	prefijoJPEG       = "data:image/jpeg;base64,"
)

// LogoSistemaJPEG bytes del logo embebido del repo (paridad LogoEmpresa.RECURSO_SISTEMA).
func LogoSistemaJPEG() []byte {
	if len(assets.LogoPesajeJPEG) == 0 {
		return nil
	}
	out := make([]byte, len(assets.LogoPesajeJPEG))
	copy(out, assets.LogoPesajeJPEG)
	return out
}

// LogoSistemaDataURL data URL del logo del sistema para PDF/empresas.
func LogoSistemaDataURL() string {
	b := LogoSistemaJPEG()
	if len(b) == 0 {
		return ""
	}
	if empiezaCon(b, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return prefijoPNG + base64.StdEncoding.EncodeToString(b)
	}
	if empiezaCon(b, []byte{0xFF, 0xD8, 0xFF}) {
		return prefijoJPEG + base64.StdEncoding.EncodeToString(b)
	}
	return ""
}

// EsDataURLLogo indica un data URL con imagen base64.
func EsDataURLLogo(valor string) bool {
	return valor != "" && strings.Contains(valor, marcaBase64)
}

// DecodificarLogoDataURL devuelve bytes de imagen o nil.
func DecodificarLogoDataURL(dataURL string) []byte {
	if !EsDataURLLogo(dataURL) {
		return nil
	}
	i := strings.Index(dataURL, marcaBase64)
	if i < 0 {
		return nil
	}
	b, err := base64.StdEncoding.DecodeString(strings.TrimSpace(dataURL[i+len(marcaBase64):]))
	if err != nil {
		return nil
	}
	return b
}

func empiezaCon(bytes, firma []byte) bool {
	if len(bytes) < len(firma) {
		return false
	}
	for i := range firma {
		if bytes[i] != firma[i] {
			return false
		}
	}
	return true
}