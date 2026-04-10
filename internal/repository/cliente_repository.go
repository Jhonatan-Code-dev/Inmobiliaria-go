package repository

import (
	"context"

	"rentals-go/ent"
	entCliente "rentals-go/ent/cliente"
	"rentals-go/internal/domain"
)

type ClienteRepoEnt struct {
	client *ent.Client
}

func NewClienteRepo(client *ent.Client) *ClienteRepoEnt {
	return &ClienteRepoEnt{client: client}
}

func (r *ClienteRepoEnt) ListarPaginado(ctx context.Context, filtros domain.ClienteFiltros) ([]*domain.Cliente, int, error) {
	query := r.client.Cliente.Query().Where(entCliente.EmpresaID(filtros.EmpresaID))

	if filtros.Busqueda != "" {
		query = query.Where(
			entCliente.Or(
				entCliente.NombresContainsFold(filtros.Busqueda),
				entCliente.ApellidosContainsFold(filtros.Busqueda),
				entCliente.DocumentoNumeroContainsFold(filtros.Busqueda),
			),
		)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	pagina := filtros.Pagina
	if pagina <= 0 {
		pagina = 1
	}
	limite := filtros.Limite
	if limite <= 0 {
		limite = 10
	}

	offset := (pagina - 1) * limite
	
	list, err := query.
		Limit(limite).
		Offset(offset).
		Order(ent.Desc(entCliente.FieldCreadoEn), ent.Desc(entCliente.FieldID)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*domain.Cliente, 0, len(list))
	for _, e := range list {
		out = append(out, mapClienteEntity(e))
	}
	return out, total, nil
}

func (r *ClienteRepoEnt) BuscarPorID(ctx context.Context, id int) (*domain.Cliente, error) {
	e, err := r.client.Cliente.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapClienteEntity(e), nil
}

func (r *ClienteRepoEnt) Crear(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error) {
	builder := r.client.Cliente.Create().
		SetEmpresaID(c.EmpresaID).
		SetTipoIdentificacionID(c.TipoIdentificacionID).
		SetDocumentoNumero(c.DocumentoNumero).
		SetNombres(c.Nombres).
		SetNillableApellidos(c.Apellidos).
		SetNillableCorreo(c.Correo).
		SetNillableFechaNacimiento(c.FechaNacimiento).
		SetNillableNacionalidad(c.Nacionalidad).
		SetNillableDireccion(c.Direccion).
		SetNillableContactoEmergencia(c.ContactoEmergencia).
		SetNillableTelefonoEmergencia(c.TelefonoEmergencia).
		SetNillableNotas(c.Notas).
		SetEstado(entCliente.Estado(c.Estado))

	e, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapClienteEntity(e), nil
}

func (r *ClienteRepoEnt) Actualizar(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error) {
	builder := r.client.Cliente.UpdateOneID(c.ID).
		SetTipoIdentificacionID(c.TipoIdentificacionID).
		SetDocumentoNumero(c.DocumentoNumero).
		SetNombres(c.Nombres).
		SetNillableApellidos(c.Apellidos).
		SetNillableCorreo(c.Correo).
		SetNillableFechaNacimiento(c.FechaNacimiento).
		SetNillableNacionalidad(c.Nacionalidad).
		SetNillableDireccion(c.Direccion).
		SetNillableContactoEmergencia(c.ContactoEmergencia).
		SetNillableTelefonoEmergencia(c.TelefonoEmergencia).
		SetNillableNotas(c.Notas).
		SetEstado(entCliente.Estado(c.Estado))

	e, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapClienteEntity(e), nil
}

func (r *ClienteRepoEnt) Eliminar(ctx context.Context, id int) error {
	return r.client.Cliente.DeleteOneID(id).Exec(ctx)
}

func mapClienteEntity(e *ent.Cliente) *domain.Cliente {
	return &domain.Cliente{
		ID:                   e.ID,
		EmpresaID:            e.EmpresaID,
		TipoIdentificacionID: e.TipoIdentificacionID,
		DocumentoNumero:      e.DocumentoNumero,
		Nombres:              e.Nombres,
		Apellidos:            e.Apellidos,
		Correo:               e.Correo,
		FechaNacimiento:      e.FechaNacimiento,
		Nacionalidad:         e.Nacionalidad,
		Direccion:            e.Direccion,
		ContactoEmergencia:   e.ContactoEmergencia,
		TelefonoEmergencia:   e.TelefonoEmergencia,
		Notas:                e.Notas,
		Estado:               string(e.Estado),
		CreadoEn:             e.CreadoEn,
	}
}
