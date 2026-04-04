package service

import (
	"sort"

	"rentals-go/internal/domain"
	"rentals-go/internal/pkg/moneda"
)

type MonedaService struct{}

func NewMonedaService() *MonedaService {
	return &MonedaService{}
}

func (s *MonedaService) Obtener(codigo string) (*domain.MonedaInfo, error) {
	info, err := moneda.ObtenerInfo(codigo)
	if err != nil {
		return nil, err
	}
	return &domain.MonedaInfo{
		Codigo:     info.Codigo,
		Decimales:  info.Decimales,
		Incremento: info.Incremento,
		Regiones:   mapRegiones(info.Regiones),
		Render: domain.MonedaRender{
			Metodo:                info.Render.Metodo,
			Currency:              info.Render.Currency,
			MinimumFractionDigits: info.Render.MinimumFractionDigits,
			MaximumFractionDigits: info.Render.MaximumFractionDigits,
		},
	}, nil
}

func (s *MonedaService) Listar() ([]domain.MonedaInfo, error) {
	lista, err := moneda.Listar()
	if err != nil {
		return nil, err
	}

	out := make([]domain.MonedaInfo, 0, len(lista))
	for _, info := range lista {
		out = append(out, domain.MonedaInfo{
			Codigo:     info.Codigo,
			Decimales:  info.Decimales,
			Incremento: info.Incremento,
			Regiones:   mapRegiones(info.Regiones),
			Render: domain.MonedaRender{
				Metodo:                info.Render.Metodo,
				Currency:              info.Render.Currency,
				MinimumFractionDigits: info.Render.MinimumFractionDigits,
				MaximumFractionDigits: info.Render.MaximumFractionDigits,
			},
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Codigo < out[j].Codigo
	})

	return out, nil
}

func mapRegiones(regiones []moneda.RegionInfo) []domain.MonedaRegion {
	out := make([]domain.MonedaRegion, 0, len(regiones))
	for _, region := range regiones {
		out = append(out, domain.MonedaRegion{
			Codigo: region.Codigo,
			Nombre: region.Nombre,
		})
	}
	return out
}
