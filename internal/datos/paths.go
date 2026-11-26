package datos

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	CarpetaDatos               = "SistemaPesaje"
	ArchivoBase                = "pesaje.db"
	ArchivoBackup              = "pesaje.db.backup"
	VariableEntornoRuta        = "PESAJE_DB_PATH"
	PropiedadSistemaRuta       = "pesaje.db.path"
	VersionEsquema             = 2
	EmailAdministrador         = "admin@sistema.com"
	ClaveAdministradorInicial  = "admin123"
)

// RutaResuelta devuelve la ruta efectiva de pesaje.db.
func RutaResuelta() string {
	if v := strings.TrimSpace(os.Getenv(PropiedadSistemaRuta)); v != "" {
		return v
	}
	if v := strings.TrimSpace(os.Getenv(VariableEntornoRuta)); v != "" {
		return v
	}
	return RutaPorDefecto()
}

// DirectorioDatos es la carpeta SistemaPesaje según el SO.
func DirectorioDatos() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		base := home
		if strings.TrimSpace(appData) != "" {
			base = appData
		}
		return filepath.Join(base, CarpetaDatos)
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", CarpetaDatos)
	default:
		return filepath.Join(home, ".local", "share", CarpetaDatos)
	}
}

func RutaPorDefecto() string {
	return filepath.Join(DirectorioDatos(), ArchivoBase)
}

// CarpetaReportes carpeta de PDF (paridad ReportesPane Java: Documents/SistemaPesaje).
func CarpetaReportes() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "Documents", CarpetaDatos)
	}
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Documents", CarpetaDatos)
	}
	return filepath.Join(home, "Documents", CarpetaDatos)
}
