package repository

import (
	"context"
	"fmt"
	"time"

	"rentals-go/ent"
	entCargo "rentals-go/ent/cargo"
	entCliente "rentals-go/ent/cliente"
	entContrato "rentals-go/ent/contrato"
	entGasto "rentals-go/ent/gasto"
	entPago "rentals-go/ent/pago"
	entPropiedad "rentals-go/ent/propiedad"
	entTicket "rentals-go/ent/ticket"
	entUnidad "rentals-go/ent/unidad"
	"rentals-go/internal/domain"
)

// DashboardRepoEnt implementa domain.DashboardRepository con Ent ORM.
type DashboardRepoEnt struct {
	client *ent.Client
}

func NewDashboardRepo(client *ent.Client) *DashboardRepoEnt {
	return &DashboardRepoEnt{client: client}
}

// centavosAFloat convierte int64 centavos a float64 con 2 decimales.
func centavosAFloat(cents int64) float64 {
	return float64(cents) / 100.0
}

// ─────────────────────────────────────────
// ResumenGeneral
// ─────────────────────────────────────────

func (r *DashboardRepoEnt) ResumenGeneral(ctx context.Context, empresaID int, ahora time.Time) (*domain.ResumenGeneral, error) {
	inicioMes := time.Date(ahora.Year(), ahora.Month(), 1, 0, 0, 0, 0, time.UTC)
	finMes := inicioMes.AddDate(0, 1, 0).Add(-time.Nanosecond)

	totalPropiedades, err := r.client.Propiedad.Query().
		Where(entPropiedad.EmpresaID(empresaID)).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("propiedades: %w", err)
	}

	totalUnidades, err := r.client.Unidad.Query().
		Where(entUnidad.HasPropiedadWith(entPropiedad.EmpresaID(empresaID))).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("unidades totales: %w", err)
	}

	unidadesOcupadas, err := r.client.Unidad.Query().
		Where(
			entUnidad.HasPropiedadWith(entPropiedad.EmpresaID(empresaID)),
			entUnidad.EstadoEQ(entUnidad.EstadoOcupado),
		).Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("unidades ocupadas: %w", err)
	}

	contratosActivos, err := r.client.Contrato.Query().
		Where(entContrato.EmpresaID(empresaID), entContrato.EstadoEQ(entContrato.EstadoActivo)).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("contratos activos: %w", err)
	}

	contratosBorrador, err := r.client.Contrato.Query().
		Where(entContrato.EmpresaID(empresaID), entContrato.EstadoEQ(entContrato.EstadoBorrador)).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("contratos borrador: %w", err)
	}

	contratosVencidos, err := r.client.Contrato.Query().
		Where(entContrato.EmpresaID(empresaID), entContrato.EstadoEQ(entContrato.EstadoVencido)).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("contratos vencidos: %w", err)
	}

	// Ingresos del mes (MontoTotal es int64 en centavos)
	pagosMes, err := r.client.Pago.Query().
		Where(
			entPago.EmpresaID(empresaID),
			entPago.EstadoEQ(entPago.EstadoConfirmado),
			entPago.FechaPagoGTE(inicioMes),
			entPago.FechaPagoLTE(finMes),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("pagos mes: %w", err)
	}
	var ingresosEsteMesCents int64
	for _, p := range pagosMes {
		ingresosEsteMesCents += p.MontoTotal
	}

	// Gastos del mes (Monto es float64 en gasto)
	gastosMes, err := r.client.Gasto.Query().
		Where(
			entGasto.EmpresaID(empresaID),
			entGasto.FechaGTE(inicioMes),
			entGasto.FechaLTE(finMes),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("gastos mes: %w", err)
	}
	var gastosEsteMes float64
	for _, g := range gastosMes {
		gastosEsteMes += g.Monto
	}

	// Cargos pendientes vencidos (Monto, Saldo en int64 centavos)
	cargosMorosos, err := r.client.Cargo.Query().
		Where(
			entCargo.HasContratoWith(entContrato.EmpresaID(empresaID)),
			entCargo.EstadoIn(entCargo.EstadoPendiente, entCargo.EstadoParcial),
			entCargo.FechaVencimientoLT(ahora),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("cargos morosos: %w", err)
	}
	var montoPendienteCents int64
	morososSet := map[int]struct{}{}
	for _, c := range cargosMorosos {
		montoPendienteCents += c.Saldo
		morososSet[c.ContratoID] = struct{}{}
	}

	ticketsAbiertos, err := r.client.Ticket.Query().
		Where(entTicket.EmpresaID(empresaID), entTicket.EstadoEQ(entTicket.EstadoAbierto)).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("tickets abiertos: %w", err)
	}

	ticketsEnProgreso, err := r.client.Ticket.Query().
		Where(entTicket.EmpresaID(empresaID), entTicket.EstadoEQ(entTicket.EstadoEnProgreso)).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("tickets en progreso: %w", err)
	}

	totalClientes, _ := r.client.Cliente.Query().
		Where(entCliente.EmpresaID(empresaID)).
		Count(ctx)

	var tasaOcupacion float64
	if totalUnidades > 0 {
		tasaOcupacion = float64(unidadesOcupadas) / float64(totalUnidades) * 100
	}

	ingresosEsteMes := centavosAFloat(ingresosEsteMesCents)
	montoPendiente := centavosAFloat(montoPendienteCents)

	return &domain.ResumenGeneral{
		TotalPropiedades:  totalPropiedades,
		TotalUnidades:     totalUnidades,
		UnidadesOcupadas:  unidadesOcupadas,
		UnidadesLibres:    totalUnidades - unidadesOcupadas,
		TasaOcupacion:     tasaOcupacion,
		ContratosActivos:  contratosActivos,
		ContratosBorrador: contratosBorrador,
		ContratosVencidos: contratosVencidos,
		IngresosEsteMes:   ingresosEsteMes,
		GastosEsteMes:     gastosEsteMes,
		BalanceNeto:       ingresosEsteMes - gastosEsteMes,
		TotalMorosos:      len(morososSet),
		MontoPendiente:    montoPendiente,
		TicketsAbiertos:   ticketsAbiertos,
		TicketsEnProgreso: ticketsEnProgreso,
		TotalClientes:     totalClientes,
		Mes:               int(ahora.Month()),
		Anio:              ahora.Year(),
	}, nil
}

// ─────────────────────────────────────────
// ResumenOcupacion
// ─────────────────────────────────────────

func (r *DashboardRepoEnt) ResumenOcupacion(ctx context.Context, empresaID int) (*domain.ResumenOcupacion, error) {
	propiedades, err := r.client.Propiedad.Query().
		Where(entPropiedad.EmpresaID(empresaID)).
		WithUnidades().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("ocupacion: %w", err)
	}

	var porPropiedad []domain.OcupacionPropiedad
	totalUnidades, totalOcupadas := 0, 0

	for _, p := range propiedades {
		ocupadas := 0
		for _, u := range p.Edges.Unidades {
			if u.Estado == entUnidad.EstadoOcupado {
				ocupadas++
			}
		}
		total := len(p.Edges.Unidades)
		totalUnidades += total
		totalOcupadas += ocupadas

		var tasa float64
		if total > 0 {
			tasa = float64(ocupadas) / float64(total) * 100
		}
		porPropiedad = append(porPropiedad, domain.OcupacionPropiedad{
			PropiedadID:   p.ID,
			Nombre:        p.Nombre,
			Direccion:     p.Direccion,
			TotalUnidades: total,
			Ocupadas:      ocupadas,
			Libres:        total - ocupadas,
			TasaOcupacion: tasa,
		})
	}

	var tasaGlobal float64
	if totalUnidades > 0 {
		tasaGlobal = float64(totalOcupadas) / float64(totalUnidades) * 100
	}

	if porPropiedad == nil {
		porPropiedad = []domain.OcupacionPropiedad{}
	}

	return &domain.ResumenOcupacion{
		TotalUnidades: totalUnidades,
		TotalOcupadas: totalOcupadas,
		TotalLibres:   totalUnidades - totalOcupadas,
		TasaGlobal:    tasaGlobal,
		PorPropiedad:  porPropiedad,
	}, nil
}

// ─────────────────────────────────────────
// ResumenMorosidad
// ─────────────────────────────────────────

func (r *DashboardRepoEnt) ResumenMorosidad(ctx context.Context, empresaID int, ahora time.Time) (*domain.ResumenMorosidad, error) {
	cargos, err := r.client.Cargo.Query().
		Where(
			entCargo.HasContratoWith(
				entContrato.EmpresaID(empresaID),
				entContrato.EstadoEQ(entContrato.EstadoActivo),
			),
			entCargo.EstadoIn(entCargo.EstadoPendiente, entCargo.EstadoParcial),
			entCargo.FechaVencimientoLT(ahora),
		).
		WithContrato(func(q *ent.ContratoQuery) {
			q.WithCliente()
			q.WithUnidad(func(uq *ent.UnidadQuery) {
				uq.WithPropiedad()
			})
		}).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("morosidad: %w", err)
	}

	type morData struct {
		clienteID        int
		nombre           string
		unidadCodigo     string
		propiedadNombre  string
		saldoCents       int64
		diasVencido      int
		fechaVencimiento time.Time
	}
	mapa := map[int]*morData{}

	for _, c := range cargos {
		cid := c.ContratoID
		dias := int(ahora.Sub(c.FechaVencimiento).Hours() / 24)
		if _, ok := mapa[cid]; !ok {
			clienteNombre := ""
			clienteID := 0
			unidadCodigo := ""
			propiedadNombre := ""
			if c.Edges.Contrato != nil {
				if c.Edges.Contrato.Edges.Cliente != nil {
					cl := c.Edges.Contrato.Edges.Cliente
					clienteID = cl.ID
					apellidos := ""
					if cl.Apellidos != nil {
						apellidos = " " + *cl.Apellidos
					}
					clienteNombre = cl.Nombres + apellidos
				}
				if c.Edges.Contrato.Edges.Unidad != nil {
					u := c.Edges.Contrato.Edges.Unidad
					unidadCodigo = u.Codigo
					if u.Edges.Propiedad != nil {
						propiedadNombre = u.Edges.Propiedad.Nombre
					}
				}
			}
			mapa[cid] = &morData{
				clienteID:        clienteID,
				nombre:           clienteNombre,
				unidadCodigo:     unidadCodigo,
				propiedadNombre:  propiedadNombre,
				fechaVencimiento: c.FechaVencimiento,
			}
		}
		mapa[cid].saldoCents += c.Saldo
		if dias > mapa[cid].diasVencido {
			mapa[cid].diasVencido = dias
			mapa[cid].fechaVencimiento = c.FechaVencimiento
		}
	}

	morosos := make([]domain.ClienteMoroso, 0, len(mapa))
	var montoTotalCents int64
	for contratoID, d := range mapa {
		morosos = append(morosos, domain.ClienteMoroso{
			ClienteID:        d.clienteID,
			NombreCompleto:   d.nombre,
			UnidadCodigo:     d.unidadCodigo,
			PropiedadNombre:  d.propiedadNombre,
			ContratoID:       contratoID,
			MontoPendiente:   centavosAFloat(d.saldoCents),
			DiasVencido:      d.diasVencido,
			FechaVencimiento: d.fechaVencimiento,
		})
		montoTotalCents += d.saldoCents
	}

	return &domain.ResumenMorosidad{
		TotalMorosos: len(morosos),
		MontoTotal:   centavosAFloat(montoTotalCents),
		Morosos:      morosos,
	}, nil
}

// ─────────────────────────────────────────
// ReporteFinanciero
// ─────────────────────────────────────────

func (r *DashboardRepoEnt) ReporteFinanciero(ctx context.Context, empresaID int, desde, hasta time.Time) (*domain.ReporteFinanciero, error) {
	pagos, err := r.client.Pago.Query().
		Where(
			entPago.EmpresaID(empresaID),
			entPago.EstadoEQ(entPago.EstadoConfirmado),
			entPago.FechaPagoGTE(desde),
			entPago.FechaPagoLTE(hasta),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("reporte pagos: %w", err)
	}

	gastos, err := r.client.Gasto.Query().
		Where(
			entGasto.EmpresaID(empresaID),
			entGasto.FechaGTE(desde),
			entGasto.FechaLTE(hasta),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("reporte gastos: %w", err)
	}

	type mesKey = string
	ingrMes := map[mesKey]int64{}
	gastMes := map[mesKey]float64{}

	var totalIngCents int64
	var totalGast float64

	for _, p := range pagos {
		k := fmt.Sprintf("%d-%02d", p.FechaPago.Year(), p.FechaPago.Month())
		ingrMes[k] += p.MontoTotal
		totalIngCents += p.MontoTotal
	}
	for _, g := range gastos {
		k := fmt.Sprintf("%d-%02d", g.Fecha.Year(), g.Fecha.Month())
		gastMes[k] += g.Monto
		totalGast += g.Monto
	}

	var serie []domain.PuntoFinanciero
	cur := time.Date(desde.Year(), desde.Month(), 1, 0, 0, 0, 0, time.UTC)
	fin := time.Date(hasta.Year(), hasta.Month(), 1, 0, 0, 0, 0, time.UTC)
	for !cur.After(fin) {
		k := fmt.Sprintf("%d-%02d", cur.Year(), cur.Month())
		ing := centavosAFloat(ingrMes[k])
		gst := gastMes[k]
		serie = append(serie, domain.PuntoFinanciero{
			Periodo:  k,
			Ingresos: ing,
			Gastos:   gst,
			Balance:  ing - gst,
		})
		cur = cur.AddDate(0, 1, 0)
	}

	if serie == nil {
		serie = []domain.PuntoFinanciero{}
	}

	totalIngresos := centavosAFloat(totalIngCents)
	return &domain.ReporteFinanciero{
		Desde:         desde.Format("2006-01-02"),
		Hasta:         hasta.Format("2006-01-02"),
		TotalIngresos: totalIngresos,
		TotalGastos:   totalGast,
		BalanceNeto:   totalIngresos - totalGast,
		Serie:         serie,
	}, nil
}

// ─────────────────────────────────────────
// ContratosProximosVencer
// ─────────────────────────────────────────

func (r *DashboardRepoEnt) ContratosProximosVencer(ctx context.Context, empresaID int, dias int, ahora time.Time) ([]domain.ContratoProximoVencer, error) {
	limite := ahora.AddDate(0, 0, dias)
	contratos, err := r.client.Contrato.Query().
		Where(
			entContrato.EmpresaID(empresaID),
			entContrato.EstadoEQ(entContrato.EstadoActivo),
			entContrato.FechaFinNotNil(),
			entContrato.FechaFinGTE(ahora),
			entContrato.FechaFinLTE(limite),
		).
		WithCliente().
		WithUnidad(func(q *ent.UnidadQuery) {
			q.WithPropiedad()
		}).
		Order(ent.Asc(entContrato.FieldFechaFin)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("contratos vencer: %w", err)
	}

	result := make([]domain.ContratoProximoVencer, 0, len(contratos))
	for _, c := range contratos {
		clienteNombre := ""
		if c.Edges.Cliente != nil {
			apellidos := ""
			if c.Edges.Cliente.Apellidos != nil {
				apellidos = " " + *c.Edges.Cliente.Apellidos
			}
			clienteNombre = c.Edges.Cliente.Nombres + apellidos
		}
		unidadCodigo, propiedadNombre := "", ""
		if c.Edges.Unidad != nil {
			unidadCodigo = c.Edges.Unidad.Codigo
			if c.Edges.Unidad.Edges.Propiedad != nil {
				propiedadNombre = c.Edges.Unidad.Edges.Propiedad.Nombre
			}
		}
		diasRestantes := int(c.FechaFin.Sub(ahora).Hours() / 24)
		result = append(result, domain.ContratoProximoVencer{
			ContratoID:      c.ID,
			Codigo:          c.Codigo,
			ClienteNombre:   clienteNombre,
			UnidadCodigo:    unidadCodigo,
			PropiedadNombre: propiedadNombre,
			FechaFin:        *c.FechaFin,
			DiasRestantes:   diasRestantes,
			MontoRenta:      centavosAFloat(c.MontoRenta),
		})
	}
	return result, nil
}

// ─────────────────────────────────────────
// EstadoCuentaCliente
// ─────────────────────────────────────────

func (r *DashboardRepoEnt) EstadoCuentaCliente(ctx context.Context, empresaID, clienteID int) (*domain.EstadoCuentaCliente, error) {
	cliente, err := r.client.Cliente.Query().
		Where(
			entCliente.ID(clienteID),
			entCliente.EmpresaID(empresaID),
		).
		WithContratos(func(q *ent.ContratoQuery) {
			q.Where(entContrato.EmpresaID(empresaID))
			q.WithCargos()
		}).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("cliente estado cuenta %d: %w", clienteID, err)
	}

	nombreCompleto := cliente.Nombres
	if cliente.Apellidos != nil {
		nombreCompleto += " " + *cliente.Apellidos
	}

	var totalCargadoCents, totalSaldoCents int64
	cargosResumen := make([]domain.CargoResumen, 0)

	for _, contrato := range cliente.Edges.Contratos {
		for _, cargo := range contrato.Edges.Cargos {
			totalCargadoCents += cargo.Monto
			totalSaldoCents += cargo.Saldo
			cargosResumen = append(cargosResumen, domain.CargoResumen{
				CargoID:          cargo.ID,
				Concepto:         string(cargo.Concepto),
				Monto:            centavosAFloat(cargo.Monto),
				Saldo:            centavosAFloat(cargo.Saldo),
				Estado:           string(cargo.Estado),
				FechaVencimiento: cargo.FechaVencimiento,
			})
		}
	}

	totalCargado := centavosAFloat(totalCargadoCents)
	saldoPendiente := centavosAFloat(totalSaldoCents)
	totalPagado := totalCargado - saldoPendiente

	return &domain.EstadoCuentaCliente{
		ClienteID:      cliente.ID,
		NombreCompleto: nombreCompleto,
		Documento:      cliente.DocumentoNumero,
		Correo:         cliente.Correo,
		TotalCargado:   totalCargado,
		TotalPagado:    totalPagado,
		SaldoPendiente: saldoPendiente,
		Cargos:         cargosResumen,
	}, nil
}

// ─────────────────────────────────────────
// TopUnidades
// ─────────────────────────────────────────

func (r *DashboardRepoEnt) TopUnidades(ctx context.Context, empresaID int, desde, hasta time.Time, limite int) ([]domain.TopUnidad, error) {
	pagos, err := r.client.Pago.Query().
		Where(
			entPago.EmpresaID(empresaID),
			entPago.EstadoEQ(entPago.EstadoConfirmado),
			entPago.FechaPagoGTE(desde),
			entPago.FechaPagoLTE(hasta),
		).
		WithContrato(func(q *ent.ContratoQuery) {
			q.WithUnidad(func(uq *ent.UnidadQuery) {
				uq.WithPropiedad()
			})
		}).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("top unidades: %w", err)
	}

	type unidadData struct {
		id              int
		codigo          string
		propiedadNombre string
		totalCents      int64
		count           int
	}
	mapa := map[int]*unidadData{}

	for _, p := range pagos {
		if p.Edges.Contrato == nil || p.Edges.Contrato.Edges.Unidad == nil {
			continue
		}
		u := p.Edges.Contrato.Edges.Unidad
		if _, ok := mapa[u.ID]; !ok {
			propNombre := ""
			if u.Edges.Propiedad != nil {
				propNombre = u.Edges.Propiedad.Nombre
			}
			mapa[u.ID] = &unidadData{
				id:              u.ID,
				codigo:          u.Codigo,
				propiedadNombre: propNombre,
			}
		}
		mapa[u.ID].totalCents += p.MontoTotal
		mapa[u.ID].count++
	}

	result := make([]domain.TopUnidad, 0, len(mapa))
	for _, d := range mapa {
		result = append(result, domain.TopUnidad{
			UnidadID:        d.id,
			Codigo:          d.codigo,
			PropiedadNombre: d.propiedadNombre,
			TotalIngresos:   centavosAFloat(d.totalCents),
			TotalPagos:      d.count,
		})
	}

	// Ordenar por total descendente
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].TotalIngresos > result[i].TotalIngresos {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	if limite > 0 && len(result) > limite {
		result = result[:limite]
	}
	return result, nil
}
