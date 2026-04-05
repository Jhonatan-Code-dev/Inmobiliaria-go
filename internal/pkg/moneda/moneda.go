package moneda

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"rentals-go/internal/pkg/tiempo"
)

// Info describe metadatos monetarios globales derivados de ISO/CLDR.
type Info struct {
	Codigo     string
	Decimales  int
	Incremento int
	Regiones   []RegionInfo
	Render     RenderInfo
}

type RegionInfo struct {
	Codigo string
	Nombre string
}

type RenderInfo struct {
	Metodo                string
	Currency              string
	MinimumFractionDigits int
	MaximumFractionDigits int
}

// NormalizarCodigo asegura ISO 4217 en mayusculas.
func NormalizarCodigo(codigo string) string {
	return strings.ToUpper(strings.TrimSpace(codigo))
}

// ValidarCodigo verifica que la moneda exista en ISO 4217.
func ValidarCodigo(codigo string) error {
	_, err := currency.ParseISO(NormalizarCodigo(codigo))
	if err != nil {
		return fmt.Errorf("codigo de moneda invalido %q: %w", codigo, err)
	}
	return nil
}

// ObtenerMonedaPorPais intenta resolver la moneda por defecto para un codigo de pais (ISO 3166-1 alpha-2).
func ObtenerMonedaPorPais(pais string) string {
	pais = strings.ToUpper(strings.TrimSpace(pais))
	r, err := language.ParseRegion(pais)
	if err != nil {
		return ""
	}
	unit, _ := currency.FromRegion(r)
	if unit.String() == "XXX" {
		return ""
	}
	return unit.String()
}

// ObtenerInfo retorna precision y simbolos de una moneda.
func ObtenerInfo(codigo string) (Info, error) {
	unit, scale, increment, err := parseUnit(codigo)
	if err != nil {
		return Info{}, err
	}
	return Info{
		Codigo:     unit.String(),
		Decimales:  scale,
		Incremento: increment,
		Regiones:   regionesDeUnidad(unit),
		Render: RenderInfo{
			Metodo:                "Intl.NumberFormat",
			Currency:              unit.String(),
			MinimumFractionDigits: scale,
			MaximumFractionDigits: scale,
		},
	}, nil
}

// Listar retorna las monedas de curso legal vigentes conocidas por CLDR/ISO.
func Listar() ([]Info, error) {
	iter := currency.Query()
	seen := make(map[string]struct{})
	out := make([]Info, 0, 180)

	for iter.Next() {
		codigo := iter.Unit().String()
		if _, ok := seen[codigo]; ok {
			continue
		}
		seen[codigo] = struct{}{}

		info, err := ObtenerInfo(codigo)
		if err != nil {
			return nil, err
		}
		out = append(out, info)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Codigo < out[j].Codigo
	})

	return out, nil
}

func parseUnit(codigo string) (currency.Unit, int, int, error) {
	unit, err := currency.ParseISO(NormalizarCodigo(codigo))
	if err != nil {
		return currency.Unit{}, 0, 0, fmt.Errorf("codigo de moneda invalido %q: %w", codigo, err)
	}
	scale, increment := currency.Standard.Rounding(unit)
	return unit, scale, increment, nil
}

func regionesDeUnidad(unit currency.Unit) []RegionInfo {
	iter := currency.Query(currency.Date(fechaConsulta()), currency.NonTender)
	seen := make(map[string]struct{})
	regiones := make([]RegionInfo, 0, 8)
	namer := display.Regions(language.Spanish)

	for iter.Next() {
		if iter.Unit() != unit || !iter.IsTender() {
			continue
		}
		region := iter.Region()
		codigo := region.String()
		if codigo == "ZZ" {
			continue
		}
		if _, ok := seen[codigo]; ok {
			continue
		}
		seen[codigo] = struct{}{}
		regiones = append(regiones, RegionInfo{
			Codigo: codigo,
			Nombre: namer.Name(region),
		})
	}

	sort.Slice(regiones, func(i, j int) bool {
		return regiones[i].Codigo < regiones[j].Codigo
	})

	return regiones
}

func fechaConsulta() time.Time {
	return tiempo.AhoraUTC()
}
