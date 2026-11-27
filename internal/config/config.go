package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"sistema-pesaje/internal/datos"
)

const (
	ModoSyncWeb          = "web"
	ModoSyncOfflineFirst = "offline_first"
)

// Sistema es la configuración persistida en config.json (paridad ConfigSistema Java).
type Sistema struct {
	URL            string `json:"url"`
	APIKey         string `json:"apikey"`
	AutoEnviar     bool   `json:"autoenviar"`
	ModoOperacion  string `json:"modo_operacion"`
	ModoSync       string `json:"modo_sync"`
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
		ModoSync:       ModoSyncOfflineFirst,
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
	leida.ModoSync = normalizarModoSync(leida.ModoSync)
	return leida
}

func normalizarModoSync(m string) string {
	m = strings.ToLower(strings.TrimSpace(m))
	switch m {
	case ModoSyncOfflineFirst, "offline-first", "offlinefirst":
		return ModoSyncOfflineFirst
	default:
		return ModoSyncWeb
	}
}

// EsOfflineFirst indica cola con catálogos antes que pesajes (centralizador offline-first).
func (s Sistema) EsOfflineFirst() bool {
	return normalizarModoSync(s.ModoSync) == ModoSyncOfflineFirst
}

// Guardar persiste config.json en la carpeta de datos.
func Guardar(s Sistema) error {
	s.ModoSync = normalizarModoSync(s.ModoSync)
	ruta := filepath.Join(datos.DirectorioDatos(), archivo)
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ruta, b, 0o644)
}
