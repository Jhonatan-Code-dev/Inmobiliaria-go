package domain

type MonedaInfo struct {
	Codigo     string
	Decimales  int
	Incremento int
	Regiones   []MonedaRegion
	Render     MonedaRender
}

type MonedaRegion struct {
	Codigo string
	Nombre string
}

type MonedaRender struct {
	Metodo                string
	Currency              string
	MinimumFractionDigits int
	MaximumFractionDigits int
}
