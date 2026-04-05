package controller

import (
	"net/http"

	"rentals-go/internal/domain"
	"rentals-go/internal/service"

	"github.com/gofiber/fiber/v2"
)

// DTOs moved to dtos.go

type MonedaController struct {
	svc *service.MonedaService
}

func NewMonedaController(svc *service.MonedaService) *MonedaController {
	return &MonedaController{svc: svc}
}

// ListarMonedas godoc
// @Summary Catalogo global de monedas
// @Description Lista las monedas ISO 4217 actualmente soportadas por el backend.
// @Tags catalogos
// @Produce json
// @Success 200 {array} monedaResponse
// @Failure 500 {object} errorResponse
// @Router /catalogos/monedas [get]
func (h *MonedaController) Listar(c *fiber.Ctx) error {
	lista, err := h.svc.Listar()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	resp := make([]monedaResponse, 0, len(lista))
	for _, info := range lista {
		resp = append(resp, mapMoneda(info))
	}

	return c.JSON(resp)
}

// ObtenerMoneda godoc
// @Summary Detalle de moneda
// @Description Retorna precision y simbolos de una moneda ISO 4217.
// @Tags catalogos
// @Produce json
// @Param codigo path string true "Codigo ISO 4217"
// @Success 200 {object} monedaResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /catalogos/monedas/{codigo} [get]
func (h *MonedaController) Obtener(c *fiber.Ctx) error {
	info, err := h.svc.Obtener(c.Params("codigo"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(mapMoneda(*info))
}

func mapMoneda(info domain.MonedaInfo) monedaResponse {
	return monedaResponse{
		Codigo:     info.Codigo,
		Decimales:  info.Decimales,
		Incremento: info.Incremento,
		Regiones:   mapRegionesResponse(info.Regiones),
		Render: monedaRenderResponse{
			Metodo:                info.Render.Metodo,
			Currency:              info.Render.Currency,
			MinimumFractionDigits: info.Render.MinimumFractionDigits,
			MaximumFractionDigits: info.Render.MaximumFractionDigits,
		},
	}
}

func mapRegionesResponse(regiones []domain.MonedaRegion) []monedaRegionResponse {
	out := make([]monedaRegionResponse, 0, len(regiones))
	for _, region := range regiones {
		out = append(out, monedaRegionResponse{
			Codigo: region.Codigo,
			Nombre: region.Nombre,
		})
	}
	return out
}
