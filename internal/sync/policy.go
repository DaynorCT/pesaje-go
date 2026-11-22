package sync

import (
	"time"

	"sistema-pesaje/internal/datos"
)

const (
	EsperaBaseSegundos    = 30
	EsperaMaximaSegundos  = 300
	EsperaMaximaExponente = 6
)

// AccionVersion indica qué hacer con una fila de catálogo del servidor.
type AccionVersion int

const (
	AccionNada AccionVersion = iota
	AccionAplicar
	AccionConflicto
)

// EsperaBackoffSegundos crece con tope (paridad SyncPolicy Java).
func EsperaBackoffSegundos(intentos int) int {
	exp := intentos
	if exp < 0 {
		exp = 0
	}
	if exp > EsperaMaximaExponente {
		exp = EsperaMaximaExponente
	}
	espera := EsperaBaseSegundos * (1 << exp)
	if espera > EsperaMaximaSegundos {
		return EsperaMaximaSegundos
	}
	return espera
}

// CalcularProximoIntento devuelve fecha-hora local en formato SQLite.
func CalcularProximoIntento(intentos int, ahora time.Time, jitterSegundos int) string {
	espera := int64(EsperaBackoffSegundos(intentos))
	if jitterSegundos > 0 {
		espera += int64(jitterSegundos)
	}
	return ahora.Add(time.Duration(espera) * time.Second).Format(datos.FormatoFechaHora)
}

// EsErrorPermanente clasifica códigos HTTP 4xx (salvo 401 y 429).
func EsErrorPermanente(codigo int) bool {
	return codigo >= 400 && codigo < 500 && codigo != 401 && codigo != 429
}

// EsFalloDeRed indica que no hubo respuesta del servidor.
func EsFalloDeRed(codigo int) bool {
	return codigo == 0
}

// DebeSincronizarCatalogos controla cada cuántas vueltas se bajan catálogos.
func DebeSincronizarCatalogos(contador, cada int) bool {
	if cada <= 1 {
		return true
	}
	return contador%cada == 0
}

// ResolverVersion aplica la política LWW del servidor.
func ResolverVersion(versionServidor, versionLocal int, pendientePush bool) AccionVersion {
	if versionLocal < 0 {
		return AccionAplicar
	}
	if versionServidor > versionLocal {
		if pendientePush {
			return AccionConflicto
		}
		return AccionAplicar
	}
	return AccionNada
}