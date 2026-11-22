package balanza

// LectorPeso lee tramas de la báscula.
type LectorPeso interface {
	Conectar() bool
	Desconectar()
	EstaConectado() bool
	LeerTrama() string
	UltimoError() string
}
