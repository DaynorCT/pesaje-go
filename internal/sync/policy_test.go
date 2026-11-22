package sync

import (
	"testing"
	"time"
)

func TestLaEsperaEntreReintentosCreceYTieneTope(t *testing.T) {
	if EsperaBackoffSegundos(0) != 30 {
		t.Fatalf("intentos 0: got %d", EsperaBackoffSegundos(0))
	}
	if EsperaBackoffSegundos(1) != 60 {
		t.Fatalf("intentos 1")
	}
	if EsperaBackoffSegundos(2) != 120 {
		t.Fatalf("intentos 2")
	}
	if EsperaBackoffSegundos(20) != 300 {
		t.Fatalf("tope")
	}
	if EsperaBackoffSegundos(-3) < 30 {
		t.Fatalf("negativo")
	}
}

func TestElProximoIntentoSeCalculaSobreLaHoraActual(t *testing.T) {
	base := time.Date(2026, 11, 20, 10, 0, 0, 0, time.Local)
	if got := CalcularProximoIntento(0, base, 0); got != "2026-11-20 10:00:30" {
		t.Fatalf("0: %q", got)
	}
	if got := CalcularProximoIntento(0, base, 7); got != "2026-11-20 10:00:37" {
		t.Fatalf("jitter: %q", got)
	}
	if got := CalcularProximoIntento(1, base, 0); got != "2026-11-20 10:01:00" {
		t.Fatalf("1: %q", got)
	}
}

func TestQueErroresSeReintentanYCualesQuedanParaRevisar(t *testing.T) {
	if !EsErrorPermanente(400) || !EsErrorPermanente(409) || !EsErrorPermanente(422) {
		t.Fatal("4xx permanentes")
	}
	if EsErrorPermanente(401) || EsErrorPermanente(429) || EsErrorPermanente(500) || EsErrorPermanente(0) || EsErrorPermanente(200) {
		t.Fatal("no permanentes")
	}
}

func TestElCodigoCeroEsUnFalloDeRed(t *testing.T) {
	if !EsFalloDeRed(0) {
		t.Fatal("0 red")
	}
	if EsFalloDeRed(200) || EsFalloDeRed(404) || EsFalloDeRed(500) {
		t.Fatal("otros")
	}
}

func TestLosCatalogosNoSeBajanEnCadaVuelta(t *testing.T) {
	if !DebeSincronizarCatalogos(0, 5) {
		t.Fatal("0")
	}
	if DebeSincronizarCatalogos(1, 5) {
		t.Fatal("1")
	}
	if !DebeSincronizarCatalogos(5, 5) || !DebeSincronizarCatalogos(10, 5) {
		t.Fatal("multiplos")
	}
	if !DebeSincronizarCatalogos(7, 1) {
		t.Fatal("cada 1")
	}
}

func TestComoSeResuelveUnaFilaDeCatalogoQueVieneDelServidor(t *testing.T) {
	if ResolverVersion(3, -1, false) != AccionAplicar {
		t.Fatal("local ausente")
	}
	if ResolverVersion(5, 3, false) != AccionAplicar {
		t.Fatal("servidor mas nuevo")
	}
	if ResolverVersion(3, 5, false) != AccionNada {
		t.Fatal("local mas nuevo")
	}
	if ResolverVersion(4, 4, false) != AccionNada {
		t.Fatal("iguales")
	}
	if ResolverVersion(5, 3, true) != AccionConflicto {
		t.Fatal("conflicto")
	}
	if ResolverVersion(3, 3, true) != AccionNada {
		t.Fatal("igual pendiente")
	}
}