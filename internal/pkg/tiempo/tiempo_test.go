package tiempo

import (
	"errors"
	"testing"
	"time"
)

func TestParsearEnZonaYConvertirMantieneUTC(t *testing.T) {
	t.Parallel()

	casos := []struct {
		nombre      string
		zona        string
		entrada     string
		layout      string
		esperadoUTC string
	}{
		{
			nombre:      "Peru",
			zona:        "America/Lima",
			layout:      "2006-01-02 15:04",
			entrada:     "2026-04-03 17:30",
			esperadoUTC: "2026-04-03T22:30:00Z",
		},
		{
			nombre:      "Espana",
			zona:        "Europe/Madrid",
			layout:      "2006-01-02 15:04",
			entrada:     "2026-08-10 09:15",
			esperadoUTC: "2026-08-10T07:15:00Z",
		},
		{
			nombre:      "Japon",
			zona:        "Asia/Tokyo",
			layout:      "2006-01-02 15:04",
			entrada:     "2026-12-01 08:00",
			esperadoUTC: "2026-11-30T23:00:00Z",
		},
		{
			nombre:      "NuevaZelanda",
			zona:        "Pacific/Auckland",
			layout:      "2006-01-02 15:04",
			entrada:     "2026-01-20 18:45",
			esperadoUTC: "2026-01-20T05:45:00Z",
		},
	}

	for _, caso := range casos {
		caso := caso
		t.Run(caso.nombre, func(t *testing.T) {
			t.Parallel()

			utc, err := ParsearEnZona(caso.layout, caso.entrada, caso.zona)
			if err != nil {
				t.Fatalf("ParsearEnZona() error = %v", err)
			}
			if got := utc.Format(time.RFC3339); got != caso.esperadoUTC {
				t.Fatalf("UTC inesperado: got %s want %s", got, caso.esperadoUTC)
			}

			local, err := EnZona(utc, caso.zona)
			if err != nil {
				t.Fatalf("EnZona() error = %v", err)
			}
			if got := local.Format(caso.layout); got != caso.entrada {
				t.Fatalf("hora local inesperada: got %s want %s", got, caso.entrada)
			}
		})
	}
}

func TestConfigurarZonaSistema(t *testing.T) {
	t.Parallel()

	if got := ZonaSistema(); got != "UTC" {
		t.Fatalf("ZonaSistema() = %s, want %s", got, "UTC")
	}

	ahora, err := AhoraSistema()
	if err != nil {
		t.Fatalf("AhoraSistema() error = %v", err)
	}
	if ahora.Location().String() != "UTC" {
		t.Fatalf("location = %s, want %s", ahora.Location().String(), "UTC")
	}
}

func TestEnZonaConDST(t *testing.T) {
	t.Parallel()

	casos := []struct {
		nombre   string
		utc      string
		zona     string
		esperado string
	}{
		{
			nombre:   "Madrid verano",
			utc:      "2026-07-01T12:00:00Z",
			zona:     "Europe/Madrid",
			esperado: "2026-07-01T14:00:00+02:00",
		},
		{
			nombre:   "Madrid invierno",
			utc:      "2026-01-01T12:00:00Z",
			zona:     "Europe/Madrid",
			esperado: "2026-01-01T13:00:00+01:00",
		},
	}

	for _, caso := range casos {
		caso := caso
		t.Run(caso.nombre, func(t *testing.T) {
			t.Parallel()

			instante, err := time.Parse(time.RFC3339, caso.utc)
			if err != nil {
				t.Fatalf("time.Parse() error = %v", err)
			}

			got, err := EnZona(instante, caso.zona)
			if err != nil {
				t.Fatalf("EnZona() error = %v", err)
			}
			if got.Format(time.RFC3339) != caso.esperado {
				t.Fatalf("EnZona() = %s, want %s", got.Format(time.RFC3339), caso.esperado)
			}
		})
	}
}

func TestParsearEnZonaConZonaInvalida(t *testing.T) {
	t.Parallel()

	if _, err := ParsearEnZona(time.RFC3339, "2026-01-01T12:00:00", "Mars/Olympus"); err == nil {
		t.Fatal("expected error for invalid timezone")
	}
}

func TestParsearEnZonaRechazaHoraInexistentePorDST(t *testing.T) {
	t.Parallel()

	_, err := ParsearEnZona("2006-01-02 15:04", "2026-03-29 02:30", "Europe/Madrid")
	if !errors.Is(err, ErrHoraInexistente) {
		t.Fatalf("expected ErrHoraInexistente, got %v", err)
	}
}

func TestParsearEnZonaRechazaHoraAmbiguaPorDST(t *testing.T) {
	t.Parallel()

	_, err := ParsearEnZona("2006-01-02 15:04", "2026-10-25 02:30", "Europe/Madrid")
	if !errors.Is(err, ErrHoraAmbigua) {
		t.Fatalf("expected ErrHoraAmbigua, got %v", err)
	}
}
