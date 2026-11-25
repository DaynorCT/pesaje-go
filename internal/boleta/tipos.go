package boleta

import "sistema-pesaje/internal/datos"

// Datos campos para el PDF (paridad sistema-pesaje-go-web/internal/boleta).
type Datos struct {
	ID              int64
	Guia            string
	Fecha           string
	Hora            string
	ChoferNombre    string
	ChoferCI        string
	ChoferLicencia  string
	VehiculoPlaca   string
	VehiculoTipo    string
	ProveedorNombre string
	ProductoNombre  string
	OrigenNombre    string
	DestinoNombre   string
	Operador        string
	PesoBruto       int
	Tara            int
	PesoNeto        int
	CostoBs         int
	FleteBs         int
}

// Config textos de encabezado boleta.
type Config struct {
	Salida          string
	Direccion       string
	Celular1        string
	Celular2        string
	Ciudad          string
	TituloSuperior  string
	Marca           string
	TituloDocumento string
	Acreditacion    string
}

func DatosDesde(d *datos.DatosBoleta) *Datos {
	if d == nil {
		return nil
	}
	return &Datos{
		ID: d.ID, Guia: d.Guia, Fecha: d.Fecha, Hora: d.Hora,
		ChoferNombre: d.ChoferNombre, ChoferCI: d.ChoferCI, ChoferLicencia: d.ChoferLicencia,
		VehiculoPlaca: d.VehiculoPlaca, VehiculoTipo: d.VehiculoTipo,
		ProveedorNombre: d.ProveedorNombre, ProductoNombre: d.ProductoNombre,
		OrigenNombre: d.OrigenNombre, DestinoNombre: d.DestinoNombre, Operador: d.Operador,
		PesoBruto: d.PesoBruto, Tara: d.Tara, PesoNeto: d.PesoNeto,
		CostoBs: d.CostoBs, FleteBs: d.FleteBs,
	}
}

func ConfigDesde(c *datos.BoletaConfig) *Config {
	if c == nil {
		return nil
	}
	return &Config{
		Salida: c.Salida, Direccion: c.Direccion, Celular1: c.Celular1, Celular2: c.Celular2,
		Ciudad: c.Ciudad, TituloSuperior: c.TituloSuperior, Marca: c.Marca,
		TituloDocumento: c.TituloDocumento, Acreditacion: c.Acreditacion,
	}
}