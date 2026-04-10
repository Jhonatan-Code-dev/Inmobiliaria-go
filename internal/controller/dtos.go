package controller

import (
	"time"

	"rentals-go/internal/domain"
)

type errorResponse struct {
	Message string `json:"message"`
}

type empresaResponse struct {
	ID             int       `json:"id"`
	Nombre         string    `json:"nombre"`
	Pais           string    `json:"pais,omitempty"`
	Moneda         string    `json:"moneda"`
	MaximoUsuarios int       `json:"maximo_usuarios"`
	Estado         bool      `json:"estado"`
	Vencimiento    time.Time `json:"vencimiento,omitempty"`
	CreadoEn       time.Time `json:"creado_en"`
}

type empresaListItemResponse struct {
	ID          int       `json:"id"`
	Nombre      string    `json:"nombre"`
	Pais        string    `json:"pais,omitempty"`
	Estado      bool      `json:"estado"`
	Vencimiento time.Time `json:"vencimiento,omitempty"`
}

type paginadorResponse struct {
	Total        int `json:"total"`
	Paginas      int `json:"paginas"`
	Pagina       int `json:"pagina"`
	PaginaActual int `json:"pagina_actual"`
	PorPagina    int `json:"por_pagina"`
}

type listadoEmpresasResponse struct {
	Datos      []empresaListItemResponse `json:"datos"`
	Paginacion paginadorResponse         `json:"paginacion"`
}

type monedaResponse struct {
	Codigo     string                 `json:"codigo"`
	Decimales  int                    `json:"decimales"`
	Incremento int                    `json:"incremento"`
	Regiones   []monedaRegionResponse `json:"regiones,omitempty"`
	Render     monedaRenderResponse   `json:"render"`
}

type monedaRegionResponse struct {
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}

type monedaRenderResponse struct {
	Metodo                string `json:"metodo"`
	Currency              string `json:"currency"`
	MinimumFractionDigits int    `json:"minimum_fraction_digits"`
	MaximumFractionDigits int    `json:"maximum_fraction_digits"`
}

type tipoIdentificacionResponse struct {
	ID     int     `json:"id"`
	Codigo string  `json:"codigo"`
	Nombre string  `json:"nombre"`
	Pais   *string `json:"pais"`
	Activo bool    `json:"activo"`
}

type usuarioResponse struct {
	ID        int `json:"id"`
	Usuario   string `json:"usuario"`
	EmpresaID int `json:"empresa_id"`
}

func mapEmpresaResponse(e *domain.Empresa) *empresaResponse {
	if e == nil {
		return nil
	}
	resp := &empresaResponse{
		ID:             e.ID,
		Nombre:         e.Nombre,
		Pais:           e.Pais,
		Moneda:         e.Moneda,
		MaximoUsuarios: e.MaximoUsuarios,
		Estado:         e.Estado,
		Vencimiento:    e.Vencimiento,
		CreadoEn:       e.CreadoEn,
	}
	return resp
}

type empresaDetalleResponse struct {
	ID             int       `json:"id"`
	Nombre         string    `json:"nombre"`
	Pais           string    `json:"pais,omitempty"`
	Moneda         string    `json:"moneda"`
	MaximoUsuarios int       `json:"maximo_usuarios"`
	Estado         bool      `json:"estado"`
	Vencimiento    time.Time `json:"vencimiento,omitempty"`
	CreadoEn       time.Time `json:"creado_en"`
}

func mapEmpresaDetalleResponse(e *domain.Empresa) *empresaDetalleResponse {
	if e == nil {
		return nil
	}
	resp := &empresaDetalleResponse{
		ID:             e.ID,
		Nombre:         e.Nombre,
		Pais:           e.Pais,
		Moneda:         e.Moneda,
		MaximoUsuarios: e.MaximoUsuarios,
		Estado:         e.Estado,
		Vencimiento:    e.Vencimiento,
		CreadoEn:       e.CreadoEn,
	}
	return resp
}
