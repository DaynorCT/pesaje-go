package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"sistema-pesaje/internal/config"
	"sistema-pesaje/internal/datos"
	"sistema-pesaje/internal/reportes"
)

type ResumenInicioDTO struct {
	PesajesTotal       int64  `json:"pesajesTotal"`
	PesajesHoy         int64  `json:"pesajesHoy"`
	VehiculosActivos   int64  `json:"vehiculosActivos"`
	ProductosActivos   int64  `json:"productosActivos"`
	ChoferesActivos    int64  `json:"choferesActivos"`
	ProveedoresActivos int64  `json:"proveedoresActivos"`
	SyncConfigurada    bool   `json:"syncConfigurada"`
	SyncEstado         string `json:"syncEstado"`
	SyncPendientes     int    `json:"syncPendientes"`
	SyncErrores        int    `json:"syncErrores"`
}

type BoletaConfigDTO struct {
	ID              int64  `json:"id"`
	Salida          string `json:"salida"`
	Direccion       string `json:"direccion"`
	Celular1        string `json:"celular1"`
	Celular2        string `json:"celular2"`
	Ciudad          string `json:"ciudad"`
	TituloSuperior  string `json:"tituloSuperior"`
	Marca           string `json:"marca"`
	TituloDocumento string `json:"tituloDocumento"`
	Acreditacion    string `json:"acreditacion"`
}

type PdfResultadoDTO struct {
	OK      bool   `json:"ok"`
	Ruta    string `json:"ruta"`
	Mensaje string `json:"mensaje"`
	EsError bool   `json:"esError"`
}

func (a *App) InicioResumen() (ResumenInicioDTO, error) {
	out := ResumenInicioDTO{
		SyncEstado: a.syncEstadoInicio(),
	}
	cfg := config.Cargar()
	if strings.TrimSpace(cfg.URL) != "" {
		out.SyncConfigurada = true
	}
	pend, syncErr := a.syncPendientesErrores()
	out.SyncPendientes = pend
	out.SyncErrores = syncErr
	if a.db == nil {
		return out, nil
	}
	r, err := datos.CargarResumenInicio(context.Background(), a.db.Conn())
	if err != nil {
		return out, err
	}
	out.PesajesTotal = r.PesajesTotal
	out.PesajesHoy = r.PesajesHoy
	out.VehiculosActivos = r.VehiculosActivos
	out.ProductosActivos = r.ProductosActivos
	out.ChoferesActivos = r.ChoferesActivos
	out.ProveedoresActivos = r.ProveedoresActivos
	return out, nil
}

func (a *App) BoletaConfigObtener() (BoletaConfigDTO, error) {
	if a.db == nil {
		return BoletaConfigDTO{}, nil
	}
	b, err := datos.ObtenerBoletaConfig(context.Background(), a.db.Conn())
	if err != nil || b == nil {
		return BoletaConfigDTO{}, err
	}
	return boletaToDTO(*b), nil
}

func (a *App) BoletaConfigGuardar(dto BoletaConfigDTO) AbmResultadoDTO {
	if a.db == nil {
		return AbmResultadoDTO{OK: false, Mensaje: "Base no disponible.", EsError: true}
	}
	err := datos.GuardarBoletaConfig(context.Background(), a.db.Conn(), datos.BoletaConfig{
		ID: dto.ID, Salida: dto.Salida, Direccion: dto.Direccion,
		Celular1: dto.Celular1, Celular2: dto.Celular2, Ciudad: dto.Ciudad,
		TituloSuperior: dto.TituloSuperior, Marca: dto.Marca,
		TituloDocumento: dto.TituloDocumento, Acreditacion: dto.Acreditacion,
	})
	if err != nil {
		return AbmResultadoDTO{OK: false, Mensaje: err.Error(), EsError: true}
	}
	return AbmResultadoDTO{OK: true, Mensaje: "Configuración de boleta guardada."}
}

func boletaToDTO(b datos.BoletaConfig) BoletaConfigDTO {
	return BoletaConfigDTO{
		ID: b.ID, Salida: b.Salida, Direccion: b.Direccion,
		Celular1: b.Celular1, Celular2: b.Celular2, Ciudad: b.Ciudad,
		TituloSuperior: b.TituloSuperior, Marca: b.Marca,
		TituloDocumento: b.TituloDocumento, Acreditacion: b.Acreditacion,
	}
}

type FiltroReporteDTO struct {
	Desde string `json:"desde"`
	Hasta string `json:"hasta"`
	Placa string `json:"placa"`
}

func (a *App) ReportesListar(filtro FiltroReporteDTO, limite int) ([]FilaPesajeDTO, error) {
	if a.db == nil {
		return nil, nil
	}
	filas, err := datos.ListarPesajesFiltrado(context.Background(), a.db.Conn(), filtro.Desde, filtro.Hasta, filtro.Placa, limite)
	if err != nil {
		return nil, err
	}
	out := make([]FilaPesajeDTO, len(filas))
	for i, f := range filas {
		out[i] = filaToDTO(f)
	}
	return out, nil
}

func (a *App) GenerarReportePdf(filtro FiltroReporteDTO) PdfResultadoDTO {
	if a.db == nil {
		return PdfResultadoDTO{OK: false, Mensaje: "Base no disponible.", EsError: true}
	}
	items, err := datos.CargarDatosReporte(context.Background(), a.db.Conn(), filtro.Desde, filtro.Hasta, filtro.Placa)
	if err != nil {
		return PdfResultadoDTO{OK: false, Mensaje: err.Error(), EsError: true}
	}
	empresa := datos.NombreEmpresaActiva(context.Background(), a.db.Conn())
	pdf, err := reportes.GenerarReporte(items, empresa)
	if err != nil {
		return PdfResultadoDTO{OK: false, Mensaje: err.Error(), EsError: true}
	}
	nombre := fmt.Sprintf("reporte-pesaje-%s.pdf", time.Now().Format("2006-01-02"))
	return a.guardarPdf(nombre, pdf, "Reporte generado")
}

func (a *App) GenerarBoletaPdf(pesajeID int64) PdfResultadoDTO {
	if a.db == nil {
		return PdfResultadoDTO{OK: false, Mensaje: "Base no disponible.", EsError: true}
	}
	d, err := datos.CargarDatosBoleta(context.Background(), a.db.Conn(), pesajeID)
	if err != nil {
		return PdfResultadoDTO{OK: false, Mensaje: err.Error(), EsError: true}
	}
	cfg, _ := datos.ObtenerBoletaConfig(context.Background(), a.db.Conn())
	pdf, err := reportes.GenerarBoleta(d, cfg)
	if err != nil {
		return PdfResultadoDTO{OK: false, Mensaje: err.Error(), EsError: true}
	}
	guia := d.Guia
	if guia == "" {
		guia = fmt.Sprintf("%d", pesajeID)
	}
	nombre := fmt.Sprintf("boleta-%s.pdf", guia)
	return a.guardarPdf(nombre, pdf, "Boleta generada")
}

func (a *App) guardarPdf(nombre string, pdf []byte, okMsg string) PdfResultadoDTO {
	dir := datos.CarpetaReportes()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return PdfResultadoDTO{OK: false, Mensaje: err.Error(), EsError: true}
	}
	dest := filepath.Join(dir, nombre)
	if err := os.WriteFile(dest, pdf, 0o644); err != nil {
		return PdfResultadoDTO{OK: false, Mensaje: err.Error(), EsError: true}
	}
	return PdfResultadoDTO{OK: true, Ruta: dest, Mensaje: okMsg + ": " + dest}
}

func (a *App) AbrirArchivo(ruta string) AbmResultadoDTO {
	ruta = strings.TrimSpace(ruta)
	if ruta == "" {
		return AbmResultadoDTO{OK: false, Mensaje: "Ruta vacía.", EsError: true}
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", ruta)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", ruta)
	default:
		cmd = exec.Command("xdg-open", ruta)
	}
	if err := cmd.Start(); err != nil {
		return AbmResultadoDTO{OK: false, Mensaje: err.Error(), EsError: true}
	}
	return AbmResultadoDTO{OK: true, Mensaje: "Abriendo archivo..."}
}