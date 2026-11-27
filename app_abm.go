package main

import (
	"context"
	"strings"

	"sistema-pesaje/internal/balanza"
	"sistema-pesaje/internal/datos"
)

type AbmFilaDTO struct {
	ID     int64             `json:"id"`
	Estado string            `json:"estado"`
	Celdas map[string]string `json:"celdas"`
}

type AbmGuardarDTO struct {
	Entidad string            `json:"entidad"`
	ID      int64             `json:"id"`
	Campos  map[string]string `json:"campos"`
}

type AbmResultadoDTO struct {
	OK      bool     `json:"ok"`
	Mensaje string   `json:"mensaje"`
	EsError bool     `json:"esError"`
}

type BalanzaConfigDTO struct {
	PuertoCOM        string `json:"puertoCom"`
	BaudRate         int    `json:"baudRate"`
	DataBits         int    `json:"dataBits"`
	Paridad          string `json:"paridad"`
	StopBits         int    `json:"stopBits"`
	FlowControl      string `json:"flowControl"`
	TimeoutMs        int    `json:"timeoutMs"`
	MetodoLectura    string `json:"metodoLectura"`
	PosicionInicio   int    `json:"posicionInicio"`
	PosicionLongitud int    `json:"posicionLongitud"`
}

func (a *App) AbmListar(entidad string) ([]AbmFilaDTO, error) {
	if a.db == nil {
		return nil, nil
	}
	filas, err := datos.AbmListar(context.Background(), a.db.Conn(), entidad)
	if err != nil {
		return nil, err
	}
	out := make([]AbmFilaDTO, len(filas))
	for i, f := range filas {
		celdas := f.Celdas
		if celdas == nil {
			celdas = map[string]string{}
		}
		out[i] = AbmFilaDTO{ID: f.ID, Estado: f.Estado, Celdas: celdas}
	}
	return out, nil
}

func (a *App) AbmGuardar(dto AbmGuardarDTO) AbmResultadoDTO {
	if a.db == nil {
		return AbmResultadoDTO{OK: false, Mensaje: "Base no disponible.", EsError: true}
	}
	entidad := strings.TrimSpace(dto.Entidad)
	if entidad == "usuarios" && dto.ID > 0 && a.sesion != nil {
		if dto.ID == a.sesion.ID && strings.EqualFold(strings.TrimSpace(dto.Campos["estado"]), "INACTIVO") {
			return AbmResultadoDTO{OK: false, Mensaje: "No puede inactivar su propio usuario.", EsError: true}
		}
	}
	if entidad == "usuarios" && dto.ID <= 0 {
		if strings.TrimSpace(dto.Campos["clave"]) == "" {
			return AbmResultadoDTO{OK: false, Mensaje: "El campo Contraseña es obligatorio.", EsError: true}
		}
	}
	id, err := datos.AbmGuardarConID(context.Background(), a.db.Conn(), entidad, dto.ID, dto.Campos)
	if err != nil {
		return AbmResultadoDTO{OK: false, Mensaje: err.Error(), EsError: true}
	}
	switch entidad {
	case "productos", "origenes", "destinos", "vehiculos", "choferes", "proveedores":
		a.encolarCatalogoWeb(entidad, id)
	}
	msg := "Registro actualizado."
	if dto.ID <= 0 {
		switch entidad {
		case "vehiculos":
			msg = "Vehículo creado."
		case "choferes":
			msg = "Chofer creado."
		case "proveedores":
			msg = "Proveedor creado."
		case "empresas":
			msg = "Empresa creada."
		case "usuarios":
			msg = "Usuario creado."
		default:
			msg = "Registro creado."
		}
	}
	return AbmResultadoDTO{OK: true, Mensaje: msg}
}

func (a *App) BalanzaObtenerConfig() BalanzaConfigDTO {
	cfg := balanza.ConfigPorDefecto()
	if a.db != nil {
		cfg = a.db.ObtenerConfigBalanza()
	}
	return balanzaToDTO(cfg)
}

func (a *App) BalanzaGuardarConfig(dto BalanzaConfigDTO) AbmResultadoDTO {
	if a.db == nil {
		return AbmResultadoDTO{OK: false, Mensaje: "Base no disponible.", EsError: true}
	}
	cfg := balanza.Config{
		PuertoCOM:        strings.TrimSpace(dto.PuertoCOM),
		BaudRate:         dto.BaudRate,
		DataBits:         dto.DataBits,
		Paridad:          dto.Paridad,
		StopBits:         dto.StopBits,
		FlowControl:      dto.FlowControl,
		TimeoutMs:        dto.TimeoutMs,
		MetodoLectura:    balanza.MetodoDesdeTexto(dto.MetodoLectura),
		PosicionInicio:   dto.PosicionInicio,
		PosicionLongitud: dto.PosicionLongitud,
	}
	if err := datos.GuardarConfigBalanza(context.Background(), a.db.Conn(), cfg); err != nil {
		return AbmResultadoDTO{OK: false, Mensaje: err.Error(), EsError: true}
	}
	return AbmResultadoDTO{OK: true, Mensaje: "Configuración de balanza guardada."}
}

func balanzaToDTO(cfg balanza.Config) BalanzaConfigDTO {
	return BalanzaConfigDTO{
		PuertoCOM:        cfg.PuertoCOM,
		BaudRate:         cfg.BaudRate,
		DataBits:         cfg.DataBits,
		Paridad:          cfg.Paridad,
		StopBits:         cfg.StopBits,
		FlowControl:      cfg.FlowControl,
		TimeoutMs:        cfg.TimeoutMs,
		MetodoLectura:    cfg.MetodoLectura.Texto(),
		PosicionInicio:   cfg.PosicionInicio,
		PosicionLongitud: cfg.PosicionLongitud,
	}
}
