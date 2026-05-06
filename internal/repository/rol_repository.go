package repository

import (
	"context"

	"rentals-go/ent"
	"rentals-go/ent/rol"
	"rentals-go/internal/domain"
)

type RolRepoEnt struct {
	client *ent.Client
}

func NewRolRepo(client *ent.Client) *RolRepoEnt {
	return &RolRepoEnt{client: client}
}

func (r *RolRepoEnt) BuscarPorNombre(ctx context.Context, nombre string) (*domain.Rol, error) {
	item, err := r.client.Rol.Query().Where(rol.NombreEQ(nombre)).First(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Rol{
		ID:          item.ID,
		Nombre:      item.Nombre,
		Descripcion: ptrToString(item.Descripcion),
	}, nil
}

func (r *RolRepoEnt) Crear(ctx context.Context, rolIn *domain.Rol) (*domain.Rol, error) {
	item, err := r.client.Rol.Create().
		SetNombre(rolIn.Nombre).
		SetNillableDescripcion(nilIfEmpty(rolIn.Descripcion)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Rol{
		ID:          item.ID,
		Nombre:      item.Nombre,
		Descripcion: ptrToString(item.Descripcion),
	}, nil
}
func (r *RolRepoEnt) Listar(ctx context.Context) ([]*domain.Rol, error) {
	items, err := r.client.Rol.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	var results []*domain.Rol
	for _, item := range items {
		results = append(results, &domain.Rol{
			ID:          item.ID,
			Nombre:      item.Nombre,
			Descripcion: ptrToString(item.Descripcion),
		})
	}
	return results, nil
}
