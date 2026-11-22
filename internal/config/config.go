package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"sistema-pesaje/internal/datos"
)

// Sistema es la configuración persistida en config.json (paridad ConfigSistema Java).
type Sistema struct {
	URL            string `json:"url"`
	APIKey         string `json:"apikey"`
	AutoEnviar     bool   `json:"autoenviar"`
	ModoOperacion  string `json:"modo_operacion"`
	AlwaysOnTop    bool   `json:"alwaysontop"`
	TimeoutHTTP    int    `json:"timeout_http"`
	CicloCatalogos int    `json:"ciclo_catalogos"`
}

const archivo = "config.json"

// Cargar lee config.json desde la carpeta SistemaPesaje; si no existe, devuelve valores por defecto.
func Cargar() Sistema {
	def := Sistema{
		URL:            "http://localhost:3000",
		AutoEnviar:     true,
		ModoOperacion:  "INTEGRADO",
		TimeoutHTTP:    4000,
		CicloCatalogos: 5,
	}
	ruta := filepath.Join(datos.DirectorioDatos(), archivo)
	b, err := os.ReadFile(ruta)
	if err != nil {
		return def
	}
	var leida Sistema
	if json.Unmarshal(b, &leida) != nil {
		return def
	}
	if strings.TrimSpace(leida.URL) == "" {
		leida.URL = def.URL
	}
	if leida.TimeoutHTTP <= 0 {
		leida.TimeoutHTTP = def.TimeoutHTTP
	}
	if leida.CicloCatalogos <= 0 {
		leida.CicloCatalogos = def.CicloCatalogos
	}
	if leida.ModoOperacion == "" {
		leida.ModoOperacion = def.ModoOperacion
	}
	return leida
}
