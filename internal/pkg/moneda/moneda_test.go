package moneda

import "testing"

func TestObtenerInfo(t *testing.T) {
	t.Parallel()

	casos := []struct {
		codigo    string
		decimales int
	}{
		{codigo: "USD", decimales: 2},
		{codigo: "JPY", decimales: 0},
		{codigo: "KWD", decimales: 3},
	}

	for _, caso := range casos {
		caso := caso
		t.Run(caso.codigo, func(t *testing.T) {
			t.Parallel()

			info, err := ObtenerInfo(caso.codigo)
			if err != nil {
				t.Fatalf("ObtenerInfo() error = %v", err)
			}
			if info.Decimales != caso.decimales {
				t.Fatalf("Decimales = %d, want %d", info.Decimales, caso.decimales)
			}
		})
	}
}

func TestObtenerInfoOficial(t *testing.T) {
	t.Parallel()

	info, err := ObtenerInfo("PEN")
	if err != nil {
		t.Fatalf("ObtenerInfo() error = %v", err)
	}
	if info.Render.Metodo != "Intl.NumberFormat" {
		t.Fatalf("Render.Metodo = %q", info.Render.Metodo)
	}
	if info.Render.Currency != "PEN" {
		t.Fatalf("Render.Currency = %q", info.Render.Currency)
	}
	if info.Render.MinimumFractionDigits != 2 || info.Render.MaximumFractionDigits != 2 {
		t.Fatalf("unexpected render fractions: %+v", info.Render)
	}
}

func TestListar(t *testing.T) {
	t.Parallel()

	lista, err := Listar()
	if err != nil {
		t.Fatalf("Listar() error = %v", err)
	}
	if len(lista) == 0 {
		t.Fatal("Listar() returned empty list")
	}

	encontroUSD := false
	for _, item := range lista {
		if item.Codigo == "USD" {
			encontroUSD = true
			break
		}
	}
	if !encontroUSD {
		t.Fatal("Listar() should include USD")
	}
}
