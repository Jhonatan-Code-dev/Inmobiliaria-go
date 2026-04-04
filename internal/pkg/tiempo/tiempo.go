package tiempo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	_ "time/tzdata"
)

const zonaUTC = "UTC"

var (
	ErrHoraAmbigua    = errors.New("la hora local es ambigua en la zona horaria indicada")
	ErrHoraInexistente = errors.New("la hora local no existe en la zona horaria indicada")
)

// PuntoTemporal representa un instante listo para persistirse y presentarse.
type PuntoTemporal struct {
	UTC         time.Time
	Zona        string
	Local       time.Time
	Despliegue  string
}

// ZonaSistema retorna la zona fija del sistema. La politica global siempre es UTC.
func ZonaSistema() string {
	return zonaUTC
}

// AhoraUTC entrega siempre el tiempo actual en UTC para persistencia.
func AhoraUTC() time.Time {
	return time.Now().UTC()
}

// AhoraSistema entrega el tiempo actual expresado en la zona global.
func AhoraSistema() (time.Time, error) {
	return EnZona(AhoraUTC(), ZonaSistema())
}

// ParsearEnZona interpreta una fecha local en la zona indicada y la normaliza a UTC.
func ParsearEnZona(layout, valor, zona string) (time.Time, error) {
	loc, err := cargarZona(zona)
	if err != nil {
		return time.Time{}, err
	}
	instante, err := time.ParseInLocation(layout, valor, loc)
	if err != nil {
		return time.Time{}, err
	}
	return validarHoraLocalExacta(layout, valor, instante, loc)
}

// NormalizarUTC fuerza cualquier time.Time a UTC antes de guardarlo.
func NormalizarUTC(valor time.Time) time.Time {
	if valor.IsZero() {
		return valor
	}
	return valor.UTC()
}

// EnZona convierte un instante UTC a la zona indicada.
func EnZona(valor time.Time, zona string) (time.Time, error) {
	loc, err := cargarZona(zona)
	if err != nil {
		return time.Time{}, err
	}
	return NormalizarUTC(valor).In(loc), nil
}

// Describir construye una vista completa lista para mostrar o depurar.
func Describir(valor time.Time, zona string) (PuntoTemporal, error) {
	utc := NormalizarUTC(valor)
	local, err := EnZona(utc, zona)
	if err != nil {
		return PuntoTemporal{}, err
	}
	return PuntoTemporal{
		UTC:        utc,
		Zona:       local.Location().String(),
		Local:      local,
		Despliegue: local.Format(time.RFC3339),
	}, nil
}

func cargarZona(zona string) (*time.Location, error) {
	objetivo := strings.TrimSpace(zona)
	if objetivo == "" {
		objetivo = ZonaSistema()
	}
	loc, err := time.LoadLocation(objetivo)
	if err != nil {
		return nil, fmt.Errorf("no se pudo cargar la zona horaria %q: %w", objetivo, err)
	}
	return loc, nil
}

func validarHoraLocalExacta(layout, valor string, instante time.Time, loc *time.Location) (time.Time, error) {
	paso := detectarPaso(layout)
	utcBase := instante.UTC()
	coincidencias := make([]time.Time, 0, 2)

	for delta := -6 * time.Hour; delta <= 6*time.Hour; delta += paso {
		candidatoUTC := utcBase.Add(delta)
		if candidatoUTC.In(loc).Format(layout) == valor {
			coincidencias = append(coincidencias, candidatoUTC.UTC())
			if len(coincidencias) > 1 {
				return time.Time{}, ErrHoraAmbigua
			}
		}
	}

	if len(coincidencias) == 0 {
		return time.Time{}, ErrHoraInexistente
	}
	return coincidencias[0], nil
}

func detectarPaso(layout string) time.Duration {
	if strings.Contains(layout, "05") || strings.Contains(layout, ".000") || strings.Contains(layout, ".999") {
		return time.Second
	}
	return time.Minute
}
