package repository

import (
	"context"

	"rentals-go/ent"
	entTipo "rentals-go/ent/tipoidentificacion"
	"rentals-go/internal/domain"
)

type TipoIdentificacionRepoEnt struct {
	client *ent.Client
}

func NewTipoIdentificacionRepo(client *ent.Client) *TipoIdentificacionRepoEnt {
	return &TipoIdentificacionRepoEnt{client: client}
}

func (r *TipoIdentificacionRepoEnt) ListarActivos(ctx context.Context) ([]*domain.TipoIdentificacion, error) {
	list, err := r.client.TipoIdentificacion.
		Query().
		Where(entTipo.ActivoEQ(true)).
		Order(ent.Asc(entTipo.FieldNombre)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*domain.TipoIdentificacion, 0, len(list))
	for _, item := range list {
		out = append(out, &domain.TipoIdentificacion{
			ID:     item.ID,
			Codigo: item.Codigo,
			Nombre: item.Nombre,
			Pais:   item.Pais,
			Activo: item.Activo,
		})
	}

	return out, nil
}

func (r *TipoIdentificacionRepoEnt) ExisteActivo(ctx context.Context, id int) (bool, error) {
	return r.client.TipoIdentificacion.
		Query().
		Where(
			entTipo.IDEQ(id),
			entTipo.ActivoEQ(true),
		).
		Exist(ctx)
}
