package balanza

const (
	PuertoComPorDefecto       = "COM1"
	BaudPorDefecto            = 9600
	DataBitsPorDefecto        = 8
	ParidadPorDefecto         = "N"
	StopBitsPorDefecto        = 1
	FlowControlPorDefecto     = "None"
	TimeoutMsPorDefecto       = 1000
	PosicionInicioPorDefecto  = 8
	PosicionLongitudPorDefecto = 5
)

// Config es la configuración de la báscula (paridad ConfigBalanza Java).
type Config struct {
	PuertoCOM         string
	BaudRate          int
	DataBits          int
	Paridad           string
	StopBits          int
	FlowControl       string
	TimeoutMs         int
	MetodoLectura     MetodoLectura
	PosicionInicio    int
	PosicionLongitud  int
}

func ConfigPorDefecto() Config {
	return Config{
		PuertoCOM:         PuertoComPorDefecto,
		BaudRate:          BaudPorDefecto,
		DataBits:          DataBitsPorDefecto,
		Paridad:           ParidadPorDefecto,
		StopBits:          StopBitsPorDefecto,
		FlowControl:       FlowControlPorDefecto,
		TimeoutMs:         TimeoutMsPorDefecto,
		MetodoLectura:     MetodoAuto,
		PosicionInicio:    PosicionInicioPorDefecto,
		PosicionLongitud:  PosicionLongitudPorDefecto,
	}
}
