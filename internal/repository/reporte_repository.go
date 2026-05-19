package repository

import (
	"context"
	"fmt"
	"time"

	"rentals-go/ent"
	entGasto "rentals-go/ent/gasto"
	entPago "rentals-go/ent/pago"
	entPropiedad "rentals-go/ent/propiedad"
	entTicket "rentals-go/ent/ticket"
	entUnidad "rentals-go/ent/unidad"
	entContrato "rentals-go/ent/contrato"
	"rentals-go/internal/domain"
)

type ReporteRepoEnt struct {
	client *ent.Client
}

func NewReporteRepo(client *ent.Client) *ReporteRepoEnt {
	return &ReporteRepoEnt{client: client}
}

// helper cents to float
func centsToFloat(cents int64) float64 {
	return float64(cents) / 100.0
}

// 1. ObtenerIngresosGastos
func (r *ReporteRepoEnt) ObtenerIngresosGastos(ctx context.Context, empresaID int, desde, hasta time.Time) (*domain.ReporteIngresosGastos, error) {
	// Query all confirmed payments in the range
	pagos, err := r.client.Pago.Query().
		Where(
			entPago.EmpresaID(empresaID),
			entPago.EstadoEQ(entPago.EstadoConfirmado),
			entPago.FechaPagoGTE(desde),
			entPago.FechaPagoLTE(hasta),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al consultar pagos: %w", err)
	}

	// Query all expenses in the range
	gastos, err := r.client.Gasto.Query().
		Where(
			entGasto.EmpresaID(empresaID),
			entGasto.FechaGTE(desde),
			entGasto.FechaLTE(hasta),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al consultar gastos: %w", err)
	}

	// Group in-memory for speed and reliability
	ingresosPorMes := map[string]int64{}
	gastosPorMes := map[string]float64{}

	var totalIngresosCents int64
	var totalGastos float64

	for _, p := range pagos {
		periodoKey := p.FechaPago.Format("2006-01")
		ingresosPorMes[periodoKey] += p.MontoTotal
		totalIngresosCents += p.MontoTotal
	}

	for _, g := range gastos {
		periodoKey := g.Fecha.Format("2006-01")
		gastosPorMes[periodoKey] += g.Monto
		totalGastos += g.Monto
	}

	// Generate complete series month by month between desde and hasta
	var serie []domain.PuntoIngresoGasto
	actual := time.Date(desde.Year(), desde.Month(), 1, 0, 0, 0, 0, time.UTC)
	fin := time.Date(hasta.Year(), hasta.Month(), 1, 0, 0, 0, 0, time.UTC)

	for !actual.After(fin) {
		periodoKey := actual.Format("2006-01")
		ing := centsToFloat(ingresosPorMes[periodoKey])
		gst := gastosPorMes[periodoKey]

		serie = append(serie, domain.PuntoIngresoGasto{
			Periodo:  periodoKey,
			Ingresos: ing,
			Gastos:   gst,
			Balance:  ing - gst,
		})

		actual = actual.AddDate(0, 1, 0)
	}

	totalIngresos := centsToFloat(totalIngresosCents)

	return &domain.ReporteIngresosGastos{
		Desde:         desde.Format("2006-01-02"),
		Hasta:         hasta.Format("2006-01-02"),
		TotalIngresos: totalIngresos,
		TotalGastos:   totalGastos,
		BalanceNeto:   totalIngresos - totalGastos,
		Serie:         serie,
	}, nil
}

// 2. ObtenerDistribucionMetodosPago
func (r *ReporteRepoEnt) ObtenerDistribucionMetodosPago(ctx context.Context, empresaID int, desde, hasta time.Time) ([]domain.DistribucionMetodoPago, error) {
	pagos, err := r.client.Pago.Query().
		Where(
			entPago.EmpresaID(empresaID),
			entPago.EstadoEQ(entPago.EstadoConfirmado),
			entPago.FechaPagoGTE(desde),
			entPago.FechaPagoLTE(hasta),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al consultar pagos: %w", err)
	}

	type agrupado struct {
		totalCents int64
		cantidad   int
	}

	totales := map[string]*agrupado{}
	var granTotalCents int64

	for _, p := range pagos {
		metodoStr := string(p.Metodo)
		if metodoStr == "" {
			metodoStr = "otro"
		}
		if _, ok := totales[metodoStr]; !ok {
			totales[metodoStr] = &agrupado{}
		}
		totales[metodoStr].totalCents += p.MontoTotal
		totales[metodoStr].cantidad++
		granTotalCents += p.MontoTotal
	}

	var resultado []domain.DistribucionMetodoPago
	for metodo, val := range totales {
		total := centsToFloat(val.totalCents)
		porcentaje := 0.0
		if granTotalCents > 0 {
			porcentaje = (float64(val.totalCents) / float64(granTotalCents)) * 100.0
		}

		resultado = append(resultado, domain.DistribucionMetodoPago{
			Metodo:        metodo,
			Total:         total,
			CantidadPagos: val.cantidad,
			Porcentaje:    porcentaje,
		})
	}

	if resultado == nil {
		resultado = []domain.DistribucionMetodoPago{}
	}

	return resultado, nil
}

// 3. ObtenerDistribucionCategoriasGastos
func (r *ReporteRepoEnt) ObtenerDistribucionCategoriasGastos(ctx context.Context, empresaID int, desde, hasta time.Time) ([]domain.DistribucionCategoriaGasto, error) {
	gastos, err := r.client.Gasto.Query().
		Where(
			entGasto.EmpresaID(empresaID),
			entGasto.FechaGTE(desde),
			entGasto.FechaLTE(hasta),
		).
		WithTipoPago().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al consultar gastos: %w", err)
	}

	type agrupado struct {
		total      float64
		cantidad   int
		tipoPagoID int
	}

	totales := map[string]*agrupado{}
	var granTotal float64

	for _, g := range gastos {
		nombreCat := "Sin Categoría"
		tipoPagoID := 0
		if g.Edges.TipoPago != nil {
			nombreCat = g.Edges.TipoPago.Nombre
			tipoPagoID = g.Edges.TipoPago.ID
		}

		if _, ok := totales[nombreCat]; !ok {
			totales[nombreCat] = &agrupado{tipoPagoID: tipoPagoID}
		}
		totales[nombreCat].total += g.Monto
		totales[nombreCat].cantidad++
		granTotal += g.Monto
	}

	var resultado []domain.DistribucionCategoriaGasto
	for cat, val := range totales {
		porcentaje := 0.0
		if granTotal > 0 {
			porcentaje = (val.total / granTotal) * 100.0
		}

		resultado = append(resultado, domain.DistribucionCategoriaGasto{
			TipoPagoID:     val.tipoPagoID,
			Categoria:      cat,
			Total:          val.total,
			CantidadGastos: val.cantidad,
			Porcentaje:     porcentaje,
		})
	}

	if resultado == nil {
		resultado = []domain.DistribucionCategoriaGasto{}
	}

	return resultado, nil
}

// 4. ObtenerRentabilidadPropiedades
func (r *ReporteRepoEnt) ObtenerRentabilidadPropiedades(ctx context.Context, empresaID int, desde, hasta time.Time) ([]domain.RentabilidadPropiedad, error) {
	// Get all properties with their units
	propiedades, err := r.client.Propiedad.Query().
		Where(entPropiedad.EmpresaID(empresaID)).
		WithUnidades().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al consultar propiedades: %w", err)
	}

	// Query all company-level expenses in range to allocate them
	gastosEmpresa, err := r.client.Gasto.Query().
		Where(
			entGasto.EmpresaID(empresaID),
			entGasto.FechaGTE(desde),
			entGasto.FechaLTE(hasta),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al consultar gastos totales: %w", err)
	}

	var totalGastosEmpresa float64
	for _, g := range gastosEmpresa {
		totalGastosEmpresa += g.Monto
	}

	// Calculate total company units for pro-rata allocation
	totalUnidadesEmpresa, err := r.client.Unidad.Query().
		Where(entUnidad.HasPropiedadWith(entPropiedad.EmpresaID(empresaID))).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al contar unidades totales: %w", err)
	}

	var resultado []domain.RentabilidadPropiedad

	for _, p := range propiedades {
		// Calculate occupancy rate
		totalUnidades := len(p.Edges.Unidades)
		unidadesOcupadas := 0
		for _, u := range p.Edges.Unidades {
			if u.Estado == entUnidad.EstadoOcupado {
				unidadesOcupadas++
			}
		}

		var ocupacionPct float64
		if totalUnidades > 0 {
			ocupacionPct = (float64(unidadesOcupadas) / float64(totalUnidades)) * 100.0
		}

		// Calculate property income (confirmed payments linked to contracts on units of this property)
		pagosPropiedad, err := r.client.Pago.Query().
			Where(
				entPago.EmpresaID(empresaID),
				entPago.EstadoEQ(entPago.EstadoConfirmado),
				entPago.FechaPagoGTE(desde),
				entPago.FechaPagoLTE(hasta),
				entPago.HasContratoWith(
					entContrato.HasUnidadWith(
						entUnidad.PropiedadID(p.ID),
					),
				),
			).All(ctx)
		if err != nil {
			return nil, fmt.Errorf("error al calcular pagos de propiedad %s: %w", p.Nombre, err)
		}

		var ingresosCents int64
		for _, py := range pagosPropiedad {
			ingresosCents += py.MontoTotal
		}
		ingresos := centsToFloat(ingresosCents)

		// Pro-rata expense allocation
		gastoPropio := 0.0
		if totalUnidadesEmpresa > 0 && totalUnidades > 0 {
			gastoPropio = (float64(totalUnidades) / float64(totalUnidadesEmpresa)) * totalGastosEmpresa
		}

		resultado = append(resultado, domain.RentabilidadPropiedad{
			PropiedadID:      p.ID,
			Nombre:           p.Nombre,
			Direccion:        p.Direccion,
			TotalUnidades:    totalUnidades,
			UnidadesOcupadas: unidadesOcupadas,
			TasaOcupacionPct: ocupacionPct,
			Ingresos:         ingresos,
			Gastos:           gastoPropio,
			Rentabilidad:     ingresos - gastoPropio,
		})
	}

	if resultado == nil {
		resultado = []domain.RentabilidadPropiedad{}
	}

	return resultado, nil
}

// 5. ObtenerResumenMantenimiento
func (r *ReporteRepoEnt) ObtenerResumenMantenimiento(ctx context.Context, empresaID int, desde, hasta time.Time) (*domain.ResumenMantenimientoReporte, error) {
	tickets, err := r.client.Ticket.Query().
		Where(
			entTicket.EmpresaID(empresaID),
			entTicket.CreadoEnGTE(desde),
			entTicket.CreadoEnLTE(hasta),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al consultar tickets: %w", err)
	}

	var porEstado domain.TicketsPorEstado
	var porPrioridad domain.TicketsPorPrioridad

	for _, t := range tickets {
		// Group by State
		switch t.Estado {
		case entTicket.EstadoAbierto:
			porEstado.Abierto++
		case entTicket.EstadoEnProgreso:
			porEstado.EnProgreso++
		case entTicket.EstadoResuelto:
			porEstado.Resuelto++
		case entTicket.EstadoCerrado:
			// cerrados act as resueltos/anulados depending on operational rules, but we'll map cerrados to resuelto/anulado count or we can extend domain to include cerrado.
			// Let's check state values: "abierto", "en_progreso", "resuelto", "cerrado"
			// Wait, domain has Abierto, EnProgreso, Resuelto, Anulado. Let's map "cerrado" to "resuelto" or keep it separated.
			// Actually, let's map cerrado to Resuelto since resolved and closed are the same for metric dashboards. Let's do that!
			porEstado.Resuelto++
		default:
			porEstado.Anulado++
		}

		// Group by Priority
		switch t.Prioridad {
		case entTicket.PrioridadBaja:
			porPrioridad.Baja++
		case entTicket.PrioridadMedia:
			porPrioridad.Media++
		case entTicket.PrioridadAlta:
			porPrioridad.Alta++
		}
	}

	return &domain.ResumenMantenimientoReporte{
		TotalTickets: len(tickets),
		PorEstado:    porEstado,
		PorPrioridad: porPrioridad,
	}, nil
}
