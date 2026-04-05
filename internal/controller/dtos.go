package controller

import (
	"rentals-go/internal/domain"
	"rentals-go/internal/pkg/moneda"
)

type errorResponse struct {
	Message string `json:"message"`
}

type empresaResponse struct {
	ID             int             `json:"id"`
	Nombre         string          `json:"nombre"`
	Pais           string          `json:"pais,omitempty"`
	Moneda         string          `json:"moneda"`
	MonedaInfo     *monedaResponse `json:"moneda_info,omitempty"`
	MaximoUsuarios int             `json:"maximo_usuarios"`
	Estado         string          `json:"estado"`
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

type usuarioResponse struct {
	ID      int    `json:"id"`
	Usuario string `json:"usuario"`
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
	}
	if info, err := moneda.ObtenerInfo(e.Moneda); err == nil {
		mapped := monedaResponse{
			Codigo:     info.Codigo,
			Decimales:  info.Decimales,
			Incremento: info.Incremento,
			Regiones:   nil, // Usually not needed in empresa detail
			Render: monedaRenderResponse{
				Metodo:                info.Render.Metodo,
				Currency:              info.Render.Currency,
				MinimumFractionDigits: info.Render.MinimumFractionDigits,
				MaximumFractionDigits: info.Render.MaximumFractionDigits,
			},
		}
		resp.MonedaInfo = &mapped
	}
	return resp
}
