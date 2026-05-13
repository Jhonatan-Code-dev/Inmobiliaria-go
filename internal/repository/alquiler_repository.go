package repository

import (
	"context"
	"fmt"
	"time"

	"rentals-go/ent"
	entCargo "rentals-go/ent/cargo"
	entCliente "rentals-go/ent/cliente"
	entContrato "rentals-go/ent/contrato"
	entMov "rentals-go/ent/movimientocaja"
	entPago "rentals-go/ent/pago"
	entPlantilla "rentals-go/ent/plantillacontrato"
	entPropiedad "rentals-go/ent/propiedad"
	entUnidad "rentals-go/ent/unidad"
	"rentals-go/internal/domain"
)

type AlquilerRepoEnt struct {
	client *ent.Client
}

func NewAlquilerRepo(client *ent.Client) *AlquilerRepoEnt {
	return &AlquilerRepoEnt{client: client}
}

func (r *AlquilerRepoEnt) ListarPaginado(ctx context.Context, filtros domain.AlquilerFiltros) ([]*domain.Alquiler, int, error) {
	query := r.client.Contrato.Query().
		Where(entContrato.EmpresaID(filtros.EmpresaID)).
		WithCliente().
		WithUnidad()

	if filtros.Estado != "" {
		query = query.Where(entContrato.EstadoEQ(entContrato.Estado(filtros.Estado)))
	}
	if filtros.UnidadID > 0 {
		query = query.Where(entContrato.UnidadIDEQ(filtros.UnidadID))
	}
	if filtros.Busqueda != "" {
		query = query.Where(
			entContrato.Or(
				entContrato.HasClienteWith(
					entCliente.Or(
						entCliente.NombresContainsFold(filtros.Busqueda),
						entCliente.ApellidosContainsFold(filtros.Busqueda),
					),
				),
				entContrato.HasUnidadWith(entUnidad.CodigoContainsFold(filtros.Busqueda)),
			),
		)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (filtros.Pagina - 1) * filtros.Limite
	if offset < 0 {
		offset = 0
	}

	list, err := query.
		Limit(filtros.Limite).
		Offset(offset).
		Order(ent.Desc(entContrato.FieldCreadoEn), ent.Desc(entContrato.FieldID)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*domain.Alquiler, 0, len(list))
	for _, item := range list {
		out = append(out, mapContratoEntity(item))
	}
	return out, total, nil
}

func (r *AlquilerRepoEnt) BuscarPorID(ctx context.Context, id int) (*domain.Alquiler, error) {
	item, err := r.client.Contrato.Query().
		Where(entContrato.IDEQ(id)).
		WithCliente().
		WithUnidad().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapContratoEntity(item), nil
}

func (r *AlquilerRepoEnt) ListarActivos(ctx context.Context, empresaID int) ([]*domain.Alquiler, error) {
	list, err := r.client.Contrato.Query().
		Where(
			entContrato.EmpresaIDEQ(empresaID),
			entContrato.EstadoIn(entContrato.EstadoActivo, entContrato.EstadoVencido),
		).
		WithCliente().
		WithUnidad().
		Order(ent.Asc(entContrato.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*domain.Alquiler, 0, len(list))
	for _, item := range list {
		out = append(out, mapContratoEntity(item))
	}
	return out, nil
}

func (r *AlquilerRepoEnt) Crear(ctx context.Context, alquiler *domain.Alquiler) (*domain.Alquiler, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollbackTx(tx)

	clienteExiste, err := tx.Cliente.Query().
		Where(entCliente.IDEQ(alquiler.ClienteID), entCliente.EmpresaID(alquiler.EmpresaID)).
		Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !clienteExiste {
		return nil, fmt.Errorf("%w: cliente no pertenece a la empresa", domain.ErrForbidden)
	}

	unidadActual, err := tx.Unidad.Query().
		Where(
			entUnidad.IDEQ(alquiler.UnidadID),
			entUnidad.HasPropiedadWith(entPropiedad.EmpresaIDEQ(alquiler.EmpresaID)),
		).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: unidad no pertenece a la empresa", domain.ErrForbidden)
	}
	if unidadActual.Estado != entUnidad.EstadoDisponible {
		return nil, fmt.Errorf("la unidad no está disponible")
	}

	updated, err := tx.Unidad.Update().
		Where(entUnidad.IDEQ(alquiler.UnidadID), entUnidad.EstadoEQ(entUnidad.EstadoDisponible)).
		SetEstado(entUnidad.EstadoOcupado).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	if updated == 0 {
		return nil, fmt.Errorf("la unidad ya fue ocupada por otro contrato")
	}

	item, err := tx.Contrato.Create().
		SetEmpresaID(alquiler.EmpresaID).
		SetClienteID(alquiler.ClienteID).
		SetUnidadID(alquiler.UnidadID).
		SetCodigo(alquiler.Codigo).
		SetTipo(entContrato.Tipo(alquiler.Tipo)).
		SetFechaInicio(alquiler.FechaInicio).
		SetNillableFechaFin(alquiler.FechaFin).
		SetDiaVencimiento(alquiler.DiaVencimiento).
		SetMoneda(alquiler.Moneda).
		SetMontoRenta(alquiler.MontoRentaCents).
		SetMontoDeposito(alquiler.MontoDepositoCts).
		SetMoraDiaria(alquiler.MoraDiariaCents).
		SetServiciosIncluidos(alquiler.ServiciosIncl).
		SetActivoParaCobro(alquiler.ActivoParaCobro).
		SetEstado(entContrato.Estado(alquiler.Estado)).
		SetNillableObservaciones(alquiler.Observaciones).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.BuscarPorID(ctx, item.ID)
}

func (r *AlquilerRepoEnt) Actualizar(ctx context.Context, alq *domain.Alquiler) (*domain.Alquiler, error) {
	item, err := r.client.Contrato.UpdateOneID(alq.ID).
		SetFechaInicio(alq.FechaInicio).
		SetNillableFechaFin(alq.FechaFin).
		SetDiaVencimiento(alq.DiaVencimiento).
		SetMontoRenta(alq.MontoRentaCents).
		SetMoraDiaria(alq.MoraDiariaCents).
		SetServiciosIncluidos(alq.ServiciosIncl).
		SetActivoParaCobro(alq.ActivoParaCobro).
		SetEstado(entContrato.Estado(alq.Estado)).
		SetNillableObservaciones(alq.Observaciones).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.BuscarPorID(ctx, item.ID)
}

func (r *AlquilerRepoEnt) Eliminar(ctx context.Context, id int) error {
	// Al eliminar un contrato, liberamos la unidad si estaba ocupada.
	alq, err := r.client.Contrato.Query().Where(entContrato.IDEQ(id)).WithUnidad().Only(ctx)
	if err != nil {
		return err
	}
	if alq.Edges.Unidad != nil {
		err = r.client.Unidad.UpdateOneID(alq.UnidadID).SetEstado(entUnidad.EstadoDisponible).Exec(ctx)
		if err != nil {
			return err
		}
	}
	return r.client.Contrato.DeleteOneID(id).Exec(ctx)
}

// --- Plantillas ---

func (r *AlquilerRepoEnt) ListarPlantillas(ctx context.Context, empresaID int) ([]*domain.PlantillaContrato, error) {
	list, err := r.client.PlantillaContrato.Query().
		Where(entPlantilla.EmpresaID(empresaID)).
		Order(ent.Asc(entPlantilla.FieldNombre)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*domain.PlantillaContrato, 0, len(list))
	for _, item := range list {
		out = append(out, mapPlantillaEntity(item))
	}
	return out, nil
}

func (r *AlquilerRepoEnt) ObtenerPlantilla(ctx context.Context, id int, empresaID int) (*domain.PlantillaContrato, error) {
	item, err := r.client.PlantillaContrato.Query().
		Where(entPlantilla.IDEQ(id), entPlantilla.EmpresaID(empresaID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapPlantillaEntity(item), nil
}

func (r *AlquilerRepoEnt) CrearPlantilla(ctx context.Context, p *domain.PlantillaContrato) (*domain.PlantillaContrato, error) {
	item, err := r.client.PlantillaContrato.Create().
		SetEmpresaID(p.EmpresaID).
		SetNombre(p.Nombre).
		SetContenido(p.Contenido).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapPlantillaEntity(item), nil
}

func (r *AlquilerRepoEnt) ActualizarPlantilla(ctx context.Context, p *domain.PlantillaContrato) (*domain.PlantillaContrato, error) {
	item, err := r.client.PlantillaContrato.UpdateOneID(p.ID).
		SetNombre(p.Nombre).
		SetContenido(p.Contenido).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapPlantillaEntity(item), nil
}

func (r *AlquilerRepoEnt) EliminarPlantilla(ctx context.Context, id int, empresaID int) error {
	_, err := r.client.PlantillaContrato.Delete().
		Where(entPlantilla.IDEQ(id), entPlantilla.EmpresaID(empresaID)).
		Exec(ctx)
	return err
}

func mapPlantillaEntity(item *ent.PlantillaContrato) *domain.PlantillaContrato {
	return &domain.PlantillaContrato{
		ID:        item.ID,
		EmpresaID: item.EmpresaID,
		Nombre:    item.Nombre,
		Contenido: item.Contenido,
		CreadoEn:  item.CreadoEn,
	}
}


type PagoAlquilerRepoEnt struct {
	client *ent.Client
}

func NewPagoAlquilerRepo(client *ent.Client) *PagoAlquilerRepoEnt {
	return &PagoAlquilerRepoEnt{client: client}
}

func (r *PagoAlquilerRepoEnt) Registrar(ctx context.Context, pago *domain.RegistroPagoAlquiler) (*domain.PagoAlquiler, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollbackTx(tx)

	contratoItem, err := tx.Contrato.Query().
		Where(entContrato.IDEQ(pago.ContratoID), entContrato.EmpresaIDEQ(pago.EmpresaID)).
		WithCliente().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: alquiler no pertenece a la empresa", domain.ErrForbidden)
	}

	periodoInicio := time.Date(pago.FechaPago.Year(), time.Month(pago.MesCorrespondiente), 1, 0, 0, 0, 0, time.UTC)
	periodoFin := periodoInicio.AddDate(0, 1, -1)
	fechaVencimiento := fechaConDiaSeguro(periodoInicio.Year(), periodoInicio.Month(), contratoItem.DiaVencimiento)

	cargoItem, err := tx.Cargo.Query().
		Where(
			entCargo.HasContratoWith(entContrato.IDEQ(contratoItem.ID)),
			entCargo.ConceptoEQ(entCargo.ConceptoRenta),
			entCargo.PeriodoInicioEQ(periodoInicio),
			entCargo.PeriodoFinEQ(periodoFin),
		).
		Only(ctx)
	if err != nil {
		cargoItem, err = tx.Cargo.Create().
			SetContratoID(contratoItem.ID).
			SetConcepto(entCargo.ConceptoRenta).
			SetDescripcion(fmt.Sprintf("Renta %02d/%d", pago.MesCorrespondiente, pago.FechaPago.Year())).
			SetMoneda(contratoItem.Moneda).
			SetPeriodoInicio(periodoInicio).
			SetPeriodoFin(periodoFin).
			SetFechaEmision(pago.FechaPago).
			SetFechaVencimiento(fechaVencimiento).
			SetMonto(contratoItem.MontoRenta).
			SetSaldo(contratoItem.MontoRenta).
			SetEstado(entCargo.EstadoPendiente).
			SetGeneradoAutomaticamente(true).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	if pago.MontoPagadoCents <= 0 {
		return nil, fmt.Errorf("monto_pagado debe ser mayor a 0")
	}
	if pago.MontoPagadoCents > cargoItem.Saldo {
		return nil, fmt.Errorf("monto_pagado excede el saldo pendiente del alquiler")
	}

	recibo := fmt.Sprintf("PAGO-%d-%d", contratoItem.ID, time.Now().UTC().UnixNano())
	clienteID := contratoItem.ClienteID
	pagoItem, err := tx.Pago.Create().
		SetEmpresaID(pago.EmpresaID).
		SetNillableClienteID(&clienteID).
		SetNillableContratoID(&contratoItem.ID).
		SetNumeroRecibo(recibo).
		SetFechaPago(pago.FechaPago).
		SetMoneda(contratoItem.Moneda).
		SetMontoTotal(pago.MontoPagadoCents).
		SetMetodo(entPago.Metodo(pago.MetodoPago)).
		SetNillableNotas(pago.Nota).
		SetEstado(entPago.EstadoConfirmado).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	_, err = tx.PagoAplicacion.Create().
		SetPagoID(pagoItem.ID).
		SetCargoID(cargoItem.ID).
		SetMoneda(contratoItem.Moneda).
		SetMontoAplicado(pago.MontoPagadoCents).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	nuevoSaldo := cargoItem.Saldo - pago.MontoPagadoCents
	estadoCargo := entCargo.EstadoParcial
	if nuevoSaldo == 0 {
		estadoCargo = entCargo.EstadoPagado
	}
	_, err = tx.Cargo.UpdateOneID(cargoItem.ID).
		SetSaldo(nuevoSaldo).
		SetEstado(estadoCargo).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	_, err = tx.MovimientoCaja.Create().
		SetEmpresaID(pago.EmpresaID).
		SetPagoID(pagoItem.ID).
		SetTipo("ingreso").
		SetConcepto(fmt.Sprintf("Pago alquiler %s", contratoItem.Codigo)).
		SetFechaMovimiento(pago.FechaPago).
		SetMoneda(contratoItem.Moneda).
		SetMonto(float64(pago.MontoPagadoCents) / 100).
		SetMetodo(entMov.Metodo(pago.MetodoPago)).
		SetNillableObservaciones(pago.Nota).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &domain.PagoAlquiler{
		ID:                 pagoItem.ID,
		EmpresaID:          pagoItem.EmpresaID,
		ContratoID:         contratoItem.ID,
		ClienteID:          &clienteID,
		NumeroRecibo:       pagoItem.NumeroRecibo,
		FechaPago:          pagoItem.FechaPago,
		Moneda:             pagoItem.Moneda,
		MontoPagado:        float64(pagoItem.MontoTotal) / 100,
		MontoPagadoCents:   pagoItem.MontoTotal,
		MetodoPago:         string(pagoItem.Metodo),
		Nota:               pagoItem.Notas,
		MesCorrespondiente: pago.MesCorrespondiente,
	}, nil
}

func (r *PagoAlquilerRepoEnt) Listar(ctx context.Context, filtros domain.PagoFiltros) ([]*domain.PagoAlquiler, int, error) {
	query := r.client.Pago.Query().
		Where(entPago.EmpresaIDEQ(filtros.EmpresaID), entPago.EstadoNEQ(entPago.EstadoAnulado)).
		WithCliente().
		WithContrato(func(q *ent.ContratoQuery) {
			q.WithUnidad()
		})

	if filtros.Busqueda != "" {
		query = query.Where(
			entPago.Or(
				entPago.NumeroReciboContainsFold(filtros.Busqueda),
				entPago.HasClienteWith(
					entCliente.Or(
						entCliente.NombresContainsFold(filtros.Busqueda),
						entCliente.ApellidosContainsFold(filtros.Busqueda),
						entCliente.DocumentoNumeroContainsFold(filtros.Busqueda),
					),
				),
				entPago.HasContratoWith(
					entContrato.HasUnidadWith(
						entUnidad.CodigoContainsFold(filtros.Busqueda),
					),
				),
			),
		)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (filtros.Pagina - 1) * filtros.Limite
	list, err := query.Limit(filtros.Limite).Offset(offset).Order(ent.Desc(entPago.FieldFechaPago)).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*domain.PagoAlquiler, 0, len(list))
	for _, item := range list {
		var clienteNombre string
		if item.Edges.Cliente != nil {
			clienteNombre = nombreClienteContrato(item.Edges.Cliente)
		}
		var unidadCodigo string
		if item.Edges.Contrato != nil && item.Edges.Contrato.Edges.Unidad != nil {
			unidadCodigo = item.Edges.Contrato.Edges.Unidad.Codigo
		}
		out = append(out, &domain.PagoAlquiler{
			ID:                 item.ID,
			EmpresaID:          item.EmpresaID,
			ContratoID:         ptrToInt(item.ContratoID),
			ClienteID:          item.ClienteID,
			ClienteNombre:      clienteNombre,
			UnidadCodigo:       unidadCodigo,
			NumeroRecibo:       item.NumeroRecibo,
			FechaPago:          item.FechaPago,
			Moneda:             item.Moneda,
			MontoPagado:        float64(item.MontoTotal) / 100,
			MontoPagadoCents:   item.MontoTotal,
			MetodoPago:         string(item.Metodo),
			Nota:               item.Notas,
		})
	}
	return out, total, nil
}

func (r *PagoAlquilerRepoEnt) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.PagoAlquiler, error) {
	item, err := r.client.Pago.Query().
		Where(entPago.IDEQ(id), entPago.EmpresaIDEQ(empresaID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.PagoAlquiler{
		ID:               item.ID,
		EmpresaID:        item.EmpresaID,
		ContratoID:       ptrToInt(item.ContratoID),
		ClienteID:        item.ClienteID,
		NumeroRecibo:     item.NumeroRecibo,
		FechaPago:        item.FechaPago,
		Moneda:           item.Moneda,
		MontoPagado:      float64(item.MontoTotal) / 100,
		MontoPagadoCents: item.MontoTotal,
		MetodoPago:       string(item.Metodo),
		Nota:             item.Notas,
	}, nil
}

func (r *PagoAlquilerRepoEnt) Eliminar(ctx context.Context, id int, empresaID int) error {
	// Anulación es compleja porque implica devolver el saldo al cargo.
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer rollbackTx(tx)

	pagoItem, err := tx.Pago.Query().
		Where(entPago.IDEQ(id), entPago.EmpresaIDEQ(empresaID)).
		WithAplicaciones().
		Only(ctx)
	if err != nil {
		return err
	}

	// Revertir aplicaciones
	for _, app := range pagoItem.Edges.Aplicaciones {
		err = tx.Cargo.UpdateOneID(app.CargoID).
			AddSaldo(app.MontoAplicado).
			SetEstado(entCargo.EstadoParcial). // Debería recalcularse, pero por ahora lo ponemos parcial
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	// Cambiar estado del pago a anulado si existe ese estado, sino eliminar aplicaciones y el pago.
	// En este esquema el estado del pago es 'confirmado', 'pendiente', 'anulado'.
	err = tx.Pago.UpdateOneID(id).SetEstado(entPago.EstadoAnulado).Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PagoAlquilerRepoEnt) Actualizar(ctx context.Context, pago *domain.PagoAlquiler) (*domain.PagoAlquiler, error) {
	item, err := r.client.Pago.UpdateOneID(pago.ID).
		SetMetodo(entPago.Metodo(pago.MetodoPago)).
		SetNillableNotas(pago.Nota).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.BuscarPorID(ctx, item.ID, item.EmpresaID)
}

func ptrToInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func (r *PagoAlquilerRepoEnt) ListarPendientesMesActual(ctx context.Context, empresaID int, now time.Time) ([]*domain.PagoPendiente, error) {
	contratos, err := r.client.Contrato.Query().
		Where(
			entContrato.EmpresaIDEQ(empresaID),
			entContrato.ActivoParaCobroEQ(true),
			entContrato.EstadoIn(entContrato.EstadoActivo, entContrato.EstadoVencido),
		).
		WithCliente().
		WithUnidad().
		All(ctx)
	if err != nil {
		return nil, err
	}

	inicio := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	fin := inicio.AddDate(0, 1, -1)
	out := make([]*domain.PagoPendiente, 0)

	for _, contratoItem := range contratos {
		cargoItem, err := r.client.Cargo.Query().
			Where(
				entCargo.HasContratoWith(entContrato.IDEQ(contratoItem.ID)),
				entCargo.ConceptoEQ(entCargo.ConceptoRenta),
				entCargo.PeriodoInicioEQ(inicio),
				entCargo.PeriodoFinEQ(fin),
			).
			Only(ctx)

		saldo := contratoItem.MontoRenta
		estado := string(contratoItem.Estado)
		fechaVenc := fechaConDiaSeguro(now.Year(), now.Month(), contratoItem.DiaVencimiento)
		if err == nil {
			saldo = cargoItem.Saldo
			estado = string(cargoItem.Estado)
			fechaVenc = cargoItem.FechaVencimiento
		}
		if saldo == 0 {
			continue
		}

		out = append(out, &domain.PagoPendiente{
			AlquilerID:       contratoItem.ID,
			Cliente:          nombreClienteContrato(contratoItem.Edges.Cliente),
			Unidad:           contratoItem.Edges.Unidad.Codigo,
			Monto:            float64(saldo) / 100,
			MontoCents:       saldo,
			FechaVencimiento: fechaVenc,
			Estado:           estado,
		})
	}

	return out, nil
}

func mapContratoEntity(item *ent.Contrato) *domain.Alquiler {
	var clienteNombre string
	if item.Edges.Cliente != nil {
		clienteNombre = nombreClienteContrato(item.Edges.Cliente)
	}
	var unidadCodigo string
	if item.Edges.Unidad != nil {
		unidadCodigo = item.Edges.Unidad.Codigo
	}
	return &domain.Alquiler{
		ID:               item.ID,
		EmpresaID:        item.EmpresaID,
		ClienteID:        item.ClienteID,
		UnidadID:         item.UnidadID,
		Codigo:           item.Codigo,
		Tipo:             string(item.Tipo),
		FechaInicio:      item.FechaInicio,
		FechaFin:         item.FechaFin,
		DiaVencimiento:   item.DiaVencimiento,
		Moneda:           item.Moneda,
		MontoRenta:       float64(item.MontoRenta) / 100,
		MontoRentaCents:  item.MontoRenta,
		MontoDeposito:    float64(item.MontoDeposito) / 100,
		MontoDepositoCts: item.MontoDeposito,
		MoraDiaria:       float64(item.MoraDiaria) / 100,
		MoraDiariaCents:  item.MoraDiaria,
		ServiciosIncl:    item.ServiciosIncluidos,
		ActivoParaCobro:  item.ActivoParaCobro,
		Estado:           string(item.Estado),
		Observaciones:    item.Observaciones,
		CreadoEn:         item.CreadoEn,
		ClienteNombre:    clienteNombre,
		UnidadCodigo:     unidadCodigo,
	}
}

func nombreClienteContrato(item *ent.Cliente) string {
	if item == nil {
		return ""
	}
	if item.Apellidos != nil && *item.Apellidos != "" {
		return item.Nombres + " " + *item.Apellidos
	}
	return item.Nombres
}

func fechaConDiaSeguro(year int, month time.Month, day int) time.Time {
	if day <= 0 {
		day = 1
	}
	ultimoDia := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	if day > ultimoDia {
		day = ultimoDia
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func rollbackTx(tx *ent.Tx) {
	if tx == nil {
		return
	}
	_ = tx.Rollback()
}
