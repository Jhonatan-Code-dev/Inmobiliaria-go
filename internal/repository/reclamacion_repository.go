package repository

import (
	"context"
	"time"

	"rentals-go/ent"
	entReclamacion "rentals-go/ent/reclamacion"
	"rentals-go/internal/domain"
)

type ReclamacionRepoEnt struct {
	client *ent.Client
}

func NewReclamacionRepo(client *ent.Client) *ReclamacionRepoEnt {
	return &ReclamacionRepoEnt{client: client}
}

func (r *ReclamacionRepoEnt) ListarPaginado(ctx context.Context, empresaID int, pag, limite int) ([]*domain.Reclamacion, int, error) {
	query := r.client.Reclamacion.Query()
	if empresaID > 0 {
		query = query.Where(entReclamacion.EmpresaID(empresaID))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (pag - 1) * limite
	if offset < 0 {
		offset = 0
	}

	list, err := query.
		Limit(limite).
		Offset(offset).
		Order(ent.Desc(entReclamacion.FieldCreadoEn)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*domain.Reclamacion, 0, len(list))
	for _, e := range list {
		out = append(out, mapReclamacionEntity(e))
	}
	return out, total, nil
}

func (r *ReclamacionRepoEnt) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.Reclamacion, error) {
	query := r.client.Reclamacion.Query().Where(entReclamacion.IDEQ(id))
	if empresaID > 0 {
		query = query.Where(entReclamacion.EmpresaIDEQ(empresaID))
	}

	e, err := query.Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapReclamacionEntity(e), nil
}

func (r *ReclamacionRepoEnt) Crear(ctx context.Context, rec *domain.Reclamacion) (*domain.Reclamacion, error) {
	builder := r.client.Reclamacion.Create().
		SetCodigo(rec.Codigo).
		SetEmpresaID(rec.EmpresaID).
		SetNombres(rec.Nombres).
		SetApellidos(rec.Apellidos).
		SetTipoDocumento(rec.TipoDocumento).
		SetNumeroDocumento(rec.NumeroDocumento).
		SetTelefono(rec.Telefono).
		SetEmail(rec.Email).
		SetDireccion(rec.Direccion).
		SetMenorEdad(rec.MenorEdad).
		SetNombreApoderado(rec.NombreApoderado).
		SetTipoBien(rec.TipoBien).
		SetMontoReclamado(rec.MontoReclamado).
		SetDescripcionBien(rec.DescripcionBien).
		SetTipoReclamacion(rec.TipoReclamacion).
		SetDetalleReclamacion(rec.DetalleReclamacion).
		SetPedidoConsumidor(rec.PedidoConsumidor).
		SetEstado(rec.Estado)

	e, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapReclamacionEntity(e), nil
}

func (r *ReclamacionRepoEnt) Actualizar(ctx context.Context, rec *domain.Reclamacion) (*domain.Reclamacion, error) {
	builder := r.client.Reclamacion.UpdateOneID(rec.ID).
		SetEstado(rec.Estado).
		SetRespuestaDetalle(rec.RespuestaDetalle)

	if rec.RespondidoEn != nil {
		builder.SetRespondidoEn(*rec.RespondidoEn)
	} else {
		builder.ClearRespondidoEn()
	}

	e, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapReclamacionEntity(e), nil
}

func (r *ReclamacionRepoEnt) Eliminar(ctx context.Context, id int, empresaID int) error {
	query := r.client.Reclamacion.Delete().Where(entReclamacion.IDEQ(id))
	if empresaID > 0 {
		query = query.Where(entReclamacion.EmpresaIDEQ(empresaID))
	}
	_, err := query.Exec(ctx)
	return err
}

func mapReclamacionEntity(e *ent.Reclamacion) *domain.Reclamacion {
	var respondidoEn *time.Time
	if e.RespondidoEn != nil {
		respondidoEn = e.RespondidoEn
	}

	return &domain.Reclamacion{
		ID:                 e.ID,
		Codigo:             e.Codigo,
		EmpresaID:          e.EmpresaID,
		Nombres:            e.Nombres,
		Apellidos:          e.Apellidos,
		TipoDocumento:      e.TipoDocumento,
		NumeroDocumento:    e.NumeroDocumento,
		Telefono:           e.Telefono,
		Email:              e.Email,
		Direccion:          e.Direccion,
		MenorEdad:          e.MenorEdad,
		NombreApoderado:    e.NombreApoderado,
		TipoBien:           e.TipoBien,
		MontoReclamado:     e.MontoReclamado,
		DescripcionBien:    e.DescripcionBien,
		TipoReclamacion:    e.TipoReclamacion,
		DetalleReclamacion: e.DetalleReclamacion,
		PedidoConsumidor:   e.PedidoConsumidor,
		Estado:             e.Estado,
		RespuestaDetalle:   e.RespuestaDetalle,
		RespondidoEn:       respondidoEn,
		CreadoEn:           e.CreadoEn,
	}
}
