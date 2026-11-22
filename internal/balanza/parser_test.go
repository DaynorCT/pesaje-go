package balanza_test

import (
	"testing"

	"sistema-pesaje/internal/balanza"
)

const tramaConPrefijo = "PESO: 12345 kg"

func configAuto() balanza.Config {
	c := balanza.ConfigPorDefecto()
	c.MetodoLectura = balanza.MetodoAuto
	return c
}

func configPosicion(inicio, longitud int) balanza.Config {
	c := balanza.ConfigPorDefecto()
	c.MetodoLectura = balanza.MetodoPosicion
	c.PosicionInicio = inicio
	c.PosicionLongitud = longitud
	return c
}

func parsearAuto(trama string) (int, bool) {
	return balanza.NuevoParser(configAuto()).Parsear(trama)
}

func parsearPosicion(trama string, inicio, longitud int) (int, bool) {
	return balanza.NuevoParser(configPosicion(inicio, longitud)).Parsear(trama)
}

func TestAutoExtraeCuandoArrancaConNumero(t *testing.T) {
	v, ok := parsearAuto("1234 kg")
	if !ok || v != 1234 {
		t.Fatalf("got %d %v", v, ok)
	}
}

func TestAutoConPrefijo(t *testing.T) {
	v, ok := parsearAuto(tramaConPrefijo)
	if !ok || v != 12345 {
		t.Fatal()
	}
}

func TestAutoTruncaDecimal(t *testing.T) {
	v, ok := parsearAuto("ST: 1234.5 kg")
	if !ok || v != 1234 {
		t.Fatal()
	}
}

func TestAutoSinDigitos(t *testing.T) {
	_, ok := parsearAuto("SIN RESPUESTA DE LA BALANZA")
	if ok {
		t.Fatal()
	}
}

func TestAutoTramaCorta(t *testing.T) {
	if _, ok := parsearAuto("7"); ok {
		t.Fatal("un dígito")
	}
	if _, ok := parsearAuto(""); ok {
		t.Fatal()
	}
	if _, ok := parsearAuto("   "); ok {
		t.Fatal()
	}
}

func TestAutoIgnoraUnDigito(t *testing.T) {
	v, ok := parsearAuto("1 2345 kg")
	if !ok || v != 2345 {
		t.Fatal()
	}
}

func TestAutoPrimeraSecuencia(t *testing.T) {
	v, ok := parsearAuto("1000 2000 kg")
	if !ok || v != 1000 {
		t.Fatal()
	}
}

func TestPosicionVentana(t *testing.T) {
	v, ok := parsearPosicion(tramaConPrefijo, 7, 5)
	if !ok || v != 12345 {
		t.Fatal()
	}
}

func TestPosicionBaseUno(t *testing.T) {
	v, ok := parsearPosicion(tramaConPrefijo, 8, 4)
	if !ok || v != 2345 {
		t.Fatal()
	}
}

func TestPosicionFueraDeRango(t *testing.T) {
	if _, ok := parsearPosicion("1234", 8, 5); ok {
		t.Fatal()
	}
}

func TestPosicionInicioInvalido(t *testing.T) {
	if _, ok := parsearPosicion(tramaConPrefijo, 0, 5); ok {
		t.Fatal()
	}
}

func TestMetodoConfigurado(t *testing.T) {
	auto, okA := balanza.NuevoParser(configAuto()).Parsear(tramaConPrefijo)
	_, okP := balanza.NuevoParser(configPosicion(1, 4)).Parsear(tramaConPrefijo)
	if !okA || auto != 12345 {
		t.Fatal()
	}
	if okP {
		t.Fatal("PESO no es número")
	}
}
