package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"sistema-pesaje/internal/datos"
	"sistema-pesaje/internal/modelo"
)

// OpcionDTO opción de combo para el frontend.
type OpcionDTO struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

type PesajeOpcionesDTO struct {
	Vehiculos   []OpcionDTO `json:"vehiculos"`
	Choferes    []OpcionDTO `json:"choferes"`
	Proveedores []OpcionDTO `json:"proveedores"`
	Productos   []OpcionDTO `json:"productos"`
	Origenes    []OpcionDTO `json:"origenes"`
	Destinos    []OpcionDTO `json:"destinos"`
}

type VehiculoDetalleDTO struct {
	Tipo string `json:"tipo"`
	Tara int    `json:"tara"`
}

type ChoferDetalleDTO struct {
	Licencia string `json:"licencia"`
}

type EstadoPesajePaneDTO struct {
	BalanzaConectada bool   `json:"balanzaConectada"`
	TextoConectar    string `json:"textoConectar"`
	PesoVivo         int    `json:"pesoVivo"`
	PesoBruto        int    `json:"pesoBruto"`
	PesoTara         int    `json:"pesoTara"`
	PesoNeto         int    `json:"pesoNeto"`
	PuedeCapturar    bool   `json:"puedeCapturar"`
	Mensaje          string `json:"mensaje"`
	MensajeEsError   bool   `json:"mensajeEsError"`
}

type FilaPesajeDTO struct {
	ID              int64  `json:"id"`
	Guia            string `json:"guia"`
	NombreChofer    string `json:"nombreChofer"`
	Placa           string `json:"placa"`
	Licencia        string `json:"licencia"`
	TipoVehiculo    string `json:"tipoVehiculo"`
	NombreProveedor string `json:"nombreProveedor"`
	Producto        string `json:"producto"`
	Origen          string `json:"origen"`
	Destino         string `json:"destino"`
	CostoBs         int    `json:"costoBs"`
	FleteBs         int    `json:"fleteBs"`
	Fecha           string `json:"fecha"`
	Hora            string `json:"hora"`
	PesoBruto       int    `json:"pesoBruto"`
	Tara            int    `json:"tara"`
	PesoNeto        int    `json:"pesoNeto"`
	Estado          string `json:"estado"`
	EstadoBalanza   string `json:"estadoBalanza"`
}

type PesajeEdicionDTO struct {
	ID               int64  `json:"id"`
	VehiculoID       int64  `json:"vehiculoId"`
	ChoferID         *int64 `json:"choferId"`
	ProveedorID      *int64 `json:"proveedorId"`
	ProductoID       *int64 `json:"productoId"`
	OrigenID         *int64 `json:"origenId"`
	DestinoID        *int64 `json:"destinoId"`
	PesoBruto        int    `json:"pesoBruto"`
	Tara             int    `json:"tara"`
	CostoBs          int    `json:"costoBs"`
	FleteBsPendiente int    `json:"fleteBsPendiente"`
	EstadoBalanza    string `json:"estadoBalanza"`
}

type ActualizarPesajeDTO struct {
	ID               int64  `json:"id"`
	VehiculoID       int64  `json:"vehiculoId"`
	ChoferID         *int64 `json:"choferId"`
	ProveedorID      *int64 `json:"proveedorId"`
	ProductoID       *int64 `json:"productoId"`
	OrigenID         *int64 `json:"origenId"`
	DestinoID        *int64 `json:"destinoId"`
	PesoBruto        int    `json:"pesoBruto"`
	Tara             int    `json:"tara"`
	CostoBs          int    `json:"costoBs"`
	FleteBsPendiente int    `json:"fleteBsPendiente"`
	TaraManual       bool   `json:"taraManual"`
}

type RegistrarPesajeDTO struct {
	VehiculoID       int64  `json:"vehiculoId"`
	ChoferID         *int64 `json:"choferId"`
	ProveedorID      *int64 `json:"proveedorId"`
	ProductoID       *int64 `json:"productoId"`
	OrigenID         *int64 `json:"origenId"`
	DestinoID        *int64 `json:"destinoId"`
	PesoBruto        int    `json:"pesoBruto"`
	Tara             int    `json:"tara"`
	CostoBs          int    `json:"costoBs"`
	FleteBsPendiente int    `json:"fleteBsPendiente"`
	TaraManual       bool   `json:"taraManual"`
}

// PesajeCerrar detiene la báscula al salir del módulo (paridad PesajePane.cerrar).
func (a *App) PesajeCerrar() {
	a.detenerBalanza()
	a.muIntegrado.Lock()
	a.integrado.conectado = false
	a.integrado.pesoVivo = 0
	a.muIntegrado.Unlock()
}

func (a *App) PesajeEstado() EstadoPesajePaneDTO {
	return a.construirEstadoPesajePane()
}

func (a *App) construirEstadoPesajePane() EstadoPesajePaneDTO {
	a.muIntegrado.Lock()
	conectado := a.integrado.conectado
	pesoVivo := a.integrado.pesoVivo
	a.muIntegrado.Unlock()

	a.muPesaje.Lock()
	bruto := a.pesajeUI.pesoBruto
	tara := a.pesajeUI.pesoTara
	msg := a.pesajeUI.mensaje
	errMsg := a.pesajeUI.esError
	a.muPesaje.Unlock()

	neto := 0
	if n, e := modelo.Neto(bruto, tara); e == nil {
		neto = n
	}

	textoBtn := "Conectar"
	if conectado {
		textoBtn = "Desconectar"
	}
	return EstadoPesajePaneDTO{
		BalanzaConectada: conectado,
		TextoConectar:    textoBtn,
		PesoVivo:         pesoVivo,
		PesoBruto:        bruto,
		PesoTara:         tara,
		PesoNeto:         neto,
		PuedeCapturar:    conectado && pesoVivo > 0,
		Mensaje:          msg,
		MensajeEsError:   errMsg,
	}
}

func (a *App) pesajeEstablecerMensaje(msg string, esError bool) {
	a.muPesaje.Lock()
	a.pesajeUI.mensaje = msg
	a.pesajeUI.esError = esError
	a.muPesaje.Unlock()
}

// PesajeAlternarBalanza conecta o desconecta la báscula en el módulo operativo.
func (a *App) PesajeAlternarBalanza() AccionDTO {
	r := a.AlternarBalanza()
	est := a.construirEstadoPesajePane()
	r.EstadoPesaje = &est
	return r
}

func (a *App) PesajeCapturarBruto() AccionDTO {
	a.muIntegrado.Lock()
	vivo := a.integrado.pesoVivo
	conectado := a.integrado.conectado
	a.muIntegrado.Unlock()
	if !conectado || vivo <= 0 {
		a.pesajeEstablecerMensaje("Todavía no hay una lectura de peso válida.", true)
		est := a.construirEstadoPesajePane()
		return AccionDTO{OK: false, EstadoPesaje: ptrEstadoPesaje(est)}
	}
	a.muPesaje.Lock()
	a.pesajeUI.pesoBruto = vivo
	a.pesajeUI.mensaje = ""
	a.pesajeUI.esError = false
	a.muPesaje.Unlock()
	est := a.construirEstadoPesajePane()
	return AccionDTO{OK: true, EstadoPesaje: ptrEstadoPesaje(est)}
}

func (a *App) PesajeCapturarTara() AccionDTO {
	a.muIntegrado.Lock()
	vivo := a.integrado.pesoVivo
	conectado := a.integrado.conectado
	a.muIntegrado.Unlock()
	if !conectado || vivo <= 0 {
		a.pesajeEstablecerMensaje("Todavía no hay una lectura de peso válida.", true)
		est := a.construirEstadoPesajePane()
		return AccionDTO{OK: false, EstadoPesaje: ptrEstadoPesaje(est)}
	}
	a.muPesaje.Lock()
	a.pesajeUI.pesoTara = vivo
	a.pesajeUI.mensaje = ""
	a.pesajeUI.esError = false
	a.muPesaje.Unlock()
	est := a.construirEstadoPesajePane()
	return AccionDTO{OK: true, EstadoPesaje: ptrEstadoPesaje(est)}
}

func ptrEstadoPesaje(e EstadoPesajePaneDTO) *EstadoPesajePaneDTO {
	x := e
	return &x
}

func (a *App) PesajeOpciones() (PesajeOpcionesDTO, error) {
	if a.db == nil {
		return PesajeOpcionesDTO{}, nil
	}
	ctx := context.Background()
	db := a.db.Conn()
	veh, _ := datos.ListarVehiculosActivos(ctx, db)
	cho, _ := datos.ListarChoferesActivos(ctx, db)
	prov, _ := datos.ListarProveedoresActivos(ctx, db)
	prod, _ := datos.ListarCatalogoActivos(ctx, db, "productos")
	orig, _ := datos.ListarCatalogoActivos(ctx, db, "origenes")
	dest, _ := datos.ListarCatalogoActivos(ctx, db, "destinos")
	return PesajeOpcionesDTO{
		Vehiculos:   mapOpciones(veh),
		Choferes:    mapOpciones(cho),
		Proveedores: mapOpciones(prov),
		Productos:   mapOpciones(prod),
		Origenes:    mapOpciones(orig),
		Destinos:    mapOpciones(dest),
	}, nil
}

func mapOpciones(in []datos.Opcion) []OpcionDTO {
	out := make([]OpcionDTO, len(in))
	for i, o := range in {
		out[i] = OpcionDTO{ID: o.ID, Label: o.Label}
	}
	return out
}

func (a *App) PesajeDetalleVehiculo(id int64) (VehiculoDetalleDTO, bool) {
	if a.db == nil || id <= 0 {
		return VehiculoDetalleDTO{}, false
	}
	d, ok, _ := datos.ObtenerDetalleVehiculo(context.Background(), a.db.Conn(), id)
	return VehiculoDetalleDTO{Tipo: d.Tipo, Tara: d.Tara}, ok
}

func (a *App) PesajeDetalleChofer(id int64) (ChoferDetalleDTO, bool) {
	if a.db == nil || id <= 0 {
		return ChoferDetalleDTO{}, false
	}
	d, ok, _ := datos.ObtenerDetalleChofer(context.Background(), a.db.Conn(), id)
	return ChoferDetalleDTO{Licencia: d.Licencia}, ok
}

func (a *App) PesajeListarGrilla(limite int) ([]FilaPesajeDTO, error) {
	if a.db == nil {
		return nil, nil
	}
	filas, err := datos.ListarPesajes(context.Background(), a.db.Conn(), limite)
	if err != nil {
		return nil, err
	}
	out := make([]FilaPesajeDTO, len(filas))
	for i, f := range filas {
		out[i] = filaToDTO(f)
	}
	return out, nil
}

func filaToDTO(f datos.FilaPesaje) FilaPesajeDTO {
	return FilaPesajeDTO{
		ID:              f.ID,
		Guia:            f.Guia,
		NombreChofer:    f.NombreChofer,
		Placa:           f.PlacaVehiculo,
		Licencia:        f.LicenciaChofer,
		TipoVehiculo:    f.TipoVehiculo,
		NombreProveedor: f.NombreProveedor,
		Producto:        f.NombreProducto,
		Origen:          f.Origen,
		Destino:         f.Destino,
		CostoBs:         nullInt(f.CostoBs),
		FleteBs:         nullInt(f.FleteBsPendiente),
		Fecha:           f.Fecha(),
		Hora:            f.Hora(),
		PesoBruto:       nullInt(f.PesoBruto),
		Tara:            nullInt(f.Tara),
		PesoNeto:        nullInt(f.PesoNeto),
		Estado:          f.Estado,
		EstadoBalanza:   f.EstadoBalanza,
	}
}

func nullInt(v sql.NullInt64) int {
	if !v.Valid {
		return 0
	}
	return int(v.Int64)
}

func (a *App) PesajeCargarEdicion(id int64) (PesajeEdicionDTO, error) {
	if a.db == nil {
		return PesajeEdicionDTO{}, fmt.Errorf("base no disponible")
	}
	p, ok, err := datos.ObtenerPesajeEdicion(context.Background(), a.db.Conn(), id)
	if err != nil {
		return PesajeEdicionDTO{}, err
	}
	if !ok {
		return PesajeEdicionDTO{}, fmt.Errorf("pesaje no encontrado")
	}
	a.muPesaje.Lock()
	a.pesajeUI.pesoBruto = p.PesoBruto
	a.pesajeUI.pesoTara = p.Tara
	a.pesajeUI.mensaje = fmt.Sprintf("Editando pesaje #%d.", id)
	a.pesajeUI.esError = false
	a.muPesaje.Unlock()
	return PesajeEdicionDTO{
		ID:               p.ID,
		VehiculoID:       p.VehiculoID,
		ChoferID:         p.ChoferID,
		ProveedorID:      p.ProveedorID,
		ProductoID:       p.ProductoID,
		OrigenID:         p.OrigenID,
		DestinoID:        p.DestinoID,
		PesoBruto:        p.PesoBruto,
		Tara:             p.Tara,
		CostoBs:          p.CostoBs,
		FleteBsPendiente: p.FleteBsPendiente,
		EstadoBalanza:    p.EstadoBalanza,
	}, nil
}

// PesajeCancelarEdicion abandona la edición en pantalla y limpia pesos del panel.
func (a *App) PesajeCancelarEdicion() EstadoPesajePaneDTO {
	a.muPesaje.Lock()
	a.pesajeUI.pesoBruto = 0
	a.pesajeUI.pesoTara = 0
	a.pesajeUI.mensaje = ""
	a.pesajeUI.esError = false
	a.muPesaje.Unlock()
	return a.construirEstadoPesajePane()
}

func (a *App) PesajeActualizarRegistro(dto ActualizarPesajeDTO) AccionDTO {
	if a.db == nil || a.sesion == nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Pesaje", Mensaje: "Sesión o base no disponible.", EsError: true}}
	}
	if dto.ID <= 0 {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Pesaje", Mensaje: "Identificador inválido.", EsError: true}}
	}
	if dto.VehiculoID <= 0 {
		a.pesajeEstablecerMensaje("Seleccione un vehículo (placa).", true)
		est := a.construirEstadoPesajePane()
		return AccionDTO{OK: false, EstadoPesaje: ptrEstadoPesaje(est)}
	}
	bruto := dto.PesoBruto
	tara := dto.Tara
	if err := modelo.Validar(bruto, tara); err != nil {
		a.pesajeEstablecerMensaje(err.Error(), true)
		est := a.construirEstadoPesajePane()
		return AccionDTO{OK: false, EstadoPesaje: ptrEstadoPesaje(est)}
	}
	existente, ok, err := datos.ObtenerPesajeEdicion(context.Background(), a.db.Conn(), dto.ID)
	if err != nil || !ok {
		msg := "No se encontró el pesaje."
		if err != nil {
			msg = err.Error()
		}
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Pesaje", Mensaje: msg, EsError: true}}
	}
	reg := datos.PesajeEdicion{
		ID:               dto.ID,
		VehiculoID:       dto.VehiculoID,
		ChoferID:         dto.ChoferID,
		ProveedorID:      dto.ProveedorID,
		ProductoID:       dto.ProductoID,
		OrigenID:         dto.OrigenID,
		DestinoID:        dto.DestinoID,
		PesoBruto:        bruto,
		Tara:             tara,
		CostoBs:          dto.CostoBs,
		FleteBsPendiente: dto.FleteBsPendiente,
		EstadoBalanza:    existente.EstadoBalanza,
	}
	if err := datos.ActualizarPesajeRegistro(context.Background(), a.db.Conn(), reg, a.sesion.ID); err != nil {
		a.pesajeEstablecerMensaje(err.Error(), true)
		est := a.construirEstadoPesajePane()
		return AccionDTO{OK: false, EstadoPesaje: ptrEstadoPesaje(est)}
	}
	a.muPesaje.Lock()
	a.pesajeUI.pesoBruto = 0
	a.pesajeUI.pesoTara = 0
	a.pesajeUI.mensaje = "Pesaje actualizado correctamente."
	a.pesajeUI.esError = false
	a.muPesaje.Unlock()
	est := a.construirEstadoPesajePane()
	return AccionDTO{OK: true, EstadoPesaje: ptrEstadoPesaje(est)}
}

func (a *App) PesajeRegistrar(dto RegistrarPesajeDTO) AccionDTO {
	if a.db == nil || a.sesion == nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Pesaje", Mensaje: "Sesión o base no disponible.", EsError: true}}
	}
	if dto.VehiculoID <= 0 {
		a.pesajeEstablecerMensaje("Seleccione un vehículo (placa).", true)
		est := a.construirEstadoPesajePane()
		return AccionDTO{OK: false, EstadoPesaje: ptrEstadoPesaje(est)}
	}
	bruto := dto.PesoBruto
	tara := dto.Tara
	if err := modelo.Validar(bruto, tara); err != nil {
		a.pesajeEstablecerMensaje(err.Error(), true)
		est := a.construirEstadoPesajePane()
		return AccionDTO{OK: false, EstadoPesaje: ptrEstadoPesaje(est)}
	}
	id, err := datos.GuardarPesaje(context.Background(), a.db.Conn(), datos.NuevoPesaje{
		VehiculoID:         dto.VehiculoID,
		ChoferID:           dto.ChoferID,
		ProveedorID:        dto.ProveedorID,
		ProductoID:         dto.ProductoID,
		OrigenID:           dto.OrigenID,
		DestinoID:          dto.DestinoID,
		PesoBruto:          bruto,
		Tara:               tara,
		CostoBs:            dto.CostoBs,
		FleteBsPendiente:   dto.FleteBsPendiente,
		PesadorID:          a.sesion.PersonaID,
		EstadoBalanza:      "EN_PROCESO",
	}, a.sesion.ID)
	if err != nil {
		a.pesajeEstablecerMensaje(err.Error(), true)
		est := a.construirEstadoPesajePane()
		return AccionDTO{OK: false, EstadoPesaje: ptrEstadoPesaje(est)}
	}
	a.muSync.Lock()
	svc := a.syncSvc
	a.muSync.Unlock()
	if svc != nil {
		_ = svc.EncolarPesaje(context.Background(), a.db.Conn(), id)
	}
	a.muPesaje.Lock()
	a.pesajeUI.pesoBruto = 0
	a.pesajeUI.pesoTara = 0
	a.pesajeUI.mensaje = "Pesaje guardado correctamente."
	a.pesajeUI.esError = false
	a.muPesaje.Unlock()
	est := a.construirEstadoPesajePane()
	return AccionDTO{OK: true, EstadoPesaje: ptrEstadoPesaje(est)}
}

func (a *App) PesajeCambiarEstado(id int64, activar bool) AccionDTO {
	if a.db == nil || a.sesion == nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Pesaje", Mensaje: "Sesión no disponible.", EsError: true}}
	}
	nuevo := "INACTIVO"
	if activar {
		nuevo = "ACTIVO"
	}
	if err := datos.CambiarEstadoRegistroPesaje(context.Background(), a.db.Conn(), id, nuevo, a.sesion.ID); err != nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Estado", Mensaje: err.Error(), EsError: true}}
	}
	return AccionDTO{OK: true}
}

func (a *App) PesajeFinalizar(id int64) AccionDTO {
	if a.db == nil || a.sesion == nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Pesaje", Mensaje: "Sesión no disponible.", EsError: true}}
	}
	if err := datos.FinalizarPesaje(context.Background(), a.db.Conn(), id, a.sesion.ID); err != nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Finalizar", Mensaje: err.Error(), EsError: true}}
	}
	return AccionDTO{OK: true}
}

func (a *App) PesajeBoleta(pesajeID int64) AccionDTO {
	res := a.GenerarBoletaPdf(pesajeID)
	if !res.OK {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Boleta", Mensaje: res.Mensaje, EsError: true}}
	}
	_ = a.AbrirArchivo(res.Ruta)
	return AccionDTO{
		OK: true,
		Aviso: &AvisoDTO{Titulo: "Boleta", Mensaje: res.Mensaje},
	}
}

func (a *App) PesajeActualizarTara(tara int) EstadoPesajePaneDTO {
	a.muPesaje.Lock()
	a.pesajeUI.pesoTara = tara
	a.muPesaje.Unlock()
	return a.construirEstadoPesajePane()
}

func (a *App) PesajeCrearVehiculo(placa, tipo string, tara int) AccionDTO {
	placa = strings.TrimSpace(placa)
	if placa == "" {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Nuevo vehículo", Mensaje: "La placa es obligatoria.", EsError: true}}
	}
	if a.db == nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Nuevo vehículo", Mensaje: "Base no disponible.", EsError: true}}
	}
	if _, err := datos.CrearVehiculo(context.Background(), a.db.Conn(), placa, strings.TrimSpace(tipo), tara); err != nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Nuevo vehículo", Mensaje: err.Error(), EsError: true}}
	}
	a.pesajeEstablecerMensaje("Vehículo creado.", false)
	return AccionDTO{OK: true}
}

func (a *App) PesajeCrearChofer(nombre, apPat, apMat, ci, tel, licencia string) AccionDTO {
	if strings.TrimSpace(nombre) == "" {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Nuevo chofer", Mensaje: "El nombre es obligatorio.", EsError: true}}
	}
	if strings.TrimSpace(licencia) == "" {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Nuevo chofer", Mensaje: "La licencia es obligatoria.", EsError: true}}
	}
	if a.db == nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Nuevo chofer", Mensaje: "Base no disponible.", EsError: true}}
	}
	if _, err := datos.CrearChofer(context.Background(), a.db.Conn(), nombre, apPat, apMat, ci, tel, licencia); err != nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Nuevo chofer", Mensaje: err.Error(), EsError: true}}
	}
	a.pesajeEstablecerMensaje("Chofer creado.", false)
	return AccionDTO{OK: true}
}

func (a *App) PesajeCrearProveedor(nombre, apPat, apMat, ci, tel, empresa, descripcion string) AccionDTO {
	if strings.TrimSpace(nombre) == "" && strings.TrimSpace(empresa) == "" {
		return AccionDTO{
			OK:    false,
			Aviso: &AvisoDTO{Titulo: "Nuevo proveedor", Mensaje: "Indique el nombre o la empresa (al menos uno).", EsError: true},
		}
	}
	if a.db == nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Nuevo proveedor", Mensaje: "Base no disponible.", EsError: true}}
	}
	if _, err := datos.CrearProveedor(context.Background(), a.db.Conn(), nombre, apPat, apMat, ci, tel, empresa, descripcion); err != nil {
		return AccionDTO{OK: false, Aviso: &AvisoDTO{Titulo: "Nuevo proveedor", Mensaje: err.Error(), EsError: true}}
	}
	a.pesajeEstablecerMensaje("Proveedor creado.", false)
	return AccionDTO{OK: true}
}