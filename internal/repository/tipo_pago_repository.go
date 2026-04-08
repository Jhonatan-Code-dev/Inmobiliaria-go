package repository

import (
	"context"

	"rentals-go/ent"
	"rentals-go/internal/domain"
)

type TipoPagoRepoEnt struct {
	client *ent.Client
}

func NewTipoPagoRepo(client *ent.Client) *TipoPagoRepoEnt {
	return &TipoPagoRepoEnt{client: client}
}

func (r *TipoPagoRepoEnt) Listar(ctx context.Context) ([]*domain.TipoPago, error) {
	list, err := r.client.TipoPago.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*domain.TipoPago, 0, len(list))
	for _, tp := range list {
		out = append(out, &domain.TipoPago{
			ID:     tp.ID,
			Nombre: tp.Nombre,
		})
	}

	return out, nil
}
