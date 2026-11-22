package modelo

import "testing"

func TestNetoEsDiferenciaBrutoTara(t *testing.T) {
	n, err := Neto(10000, 2000)
	if err != nil || n != 8000 {
		t.Fatalf("neto=%d err=%v", n, err)
	}
	n, err = Neto(5000, 5000)
	if err != nil || n != 0 {
		t.Fatalf("neto=%d err=%v", n, err)
	}
}

func TestBrutoCeroRechazado(t *testing.T) {
	_, err := Neto(0, 0)
	if err != ErrBrutoInvalido {
		t.Fatalf("esperaba ErrBrutoInvalido, obtuvo %v", err)
	}
}

func TestTaraNegativaRechazada(t *testing.T) {
	_, err := Neto(10000, -1)
	if err != ErrTaraNegativa {
		t.Fatalf("esperaba ErrTaraNegativa, obtuvo %v", err)
	}
}

func TestTaraMayorBrutoRechazada(t *testing.T) {
	_, err := Neto(10000, 12000)
	if err != ErrTaraMayorQueBruto {
		t.Fatalf("esperaba ErrTaraMayorQueBruto, obtuvo %v", err)
	}
}