package reportes

import (
	"testing"

	"sistema-pesaje/internal/datos"
)

func TestGenerarBoletaProducePDF(t *testing.T) {
	d := &datos.DatosBoleta{
		ID: 1, Guia: "G-1", Fecha: "22/11/2026", Hora: "10:00",
		VehiculoPlaca: "ABC123", PesoBruto: 1000, Tara: 200, PesoNeto: 800,
	}
	cfg := &datos.BoletaConfig{Marca: "TEST", TituloSuperior: "BOLETA"}
	b, err := GenerarBoleta(d, cfg, datos.LogoSistemaJPEG())
	if err != nil || len(b) < 4 || string(b[:4]) != "%PDF" {
		t.Fatalf("pdf boleta invalido: len=%d err=%v", len(b), err)
	}
	if len(datos.LogoSistemaJPEG()) == 0 {
		t.Fatal("logo del sistema no embebido")
	}
	if len(b) < 8000 {
		t.Fatalf("pdf boleta sin imagen esperada: len=%d", len(b))
	}
}

func TestGenerarReporteProducePDF(t *testing.T) {
	b, err := GenerarReporte(nil, "Empresa")
	if err != nil || len(b) < 4 || string(b[:4]) != "%PDF" {
		t.Fatalf("pdf reporte invalido: err=%v", err)
	}
}