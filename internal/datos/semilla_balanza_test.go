package datos

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSemillaConfigBalanzaSIMEnBaseNueva(t *testing.T) {
	t.Setenv("PESAJE_BALANZA_SIM", "1")
	dir := t.TempDir()
	ruta := filepath.Join(dir, "pesaje.db")
	t.Setenv(PropiedadSistemaRuta, ruta)
	CerrarInstancia()

	base, err := ObtenerInstancia()
	if err != nil {
		t.Fatalf("abrir: %v", err)
	}
	defer func() {
		CerrarInstancia()
		_ = os.RemoveAll(dir)
	}()

	cfg := base.ObtenerConfigBalanza()
	if cfg.PuertoCOM != "SIM" {
		t.Fatalf("puerto=%q want SIM", cfg.PuertoCOM)
	}
}

func TestSemillaConfigBalanzaNoSinEnv(t *testing.T) {
	t.Setenv("PESAJE_BALANZA_SIM", "")
	dir := t.TempDir()
	ruta := filepath.Join(dir, "pesaje.db")
	t.Setenv(PropiedadSistemaRuta, ruta)
	CerrarInstancia()

	base, err := ObtenerInstancia()
	if err != nil {
		t.Fatalf("abrir: %v", err)
	}
	defer func() {
		CerrarInstancia()
		_ = os.RemoveAll(dir)
	}()

	var n int
	if err := base.Conn().QueryRow(`SELECT COUNT(*) FROM config_balanza`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("sin env no debe insertar config_balanza, rows=%d", n)
	}
}
