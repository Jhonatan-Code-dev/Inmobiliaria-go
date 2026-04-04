package repository

import (
	"context"

	"rentals-go/ent"
)

type MembresiaRepoEnt struct {
	client *ent.Client
}

func NewMembresiaRepo(client *ent.Client) *MembresiaRepoEnt {
	return &MembresiaRepoEnt{client: client}
}

func (r *MembresiaRepoEnt) AsignarPrincipal(ctx context.Context, empresaID, usuarioID, rolID int) error {
	_, err := r.client.EmpresaUsuario.Create().
		SetEmpresaID(empresaID).
		SetUsuarioID(usuarioID).
		SetRolID(rolID).
		SetPrincipal(true).
		Save(ctx)
	return err
}
