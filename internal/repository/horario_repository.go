package repository

import (
	"context"
	"fmt"

	"rentals-go/ent"
	"rentals-go/ent/horario"
	"rentals-go/internal/domain"
)

type HorarioRepoEnt struct {
	client *ent.Client
}

func NewHorarioRepo(client *ent.Client) *HorarioRepoEnt {
	return &HorarioRepoEnt{client: client}
}

func (r *HorarioRepoEnt) BuscarPorUsuario(ctx context.Context, usuarioID int, empresaID int) (*domain.Horario, error) {
	h, err := r.client.Horario.Query().
		Where(
			horario.UsuarioIDEQ(usuarioID),
			horario.EmpresaIDEQ(empresaID),
		).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("horario no encontrado para el usuario")
		}
		return nil, fmt.Errorf("error al buscar horario: %v", err)
	}

	return mapHorarioEntToDomain(h), nil
}

func (r *HorarioRepoEnt) Crear(ctx context.Context, h *domain.Horario) (*domain.Horario, error) {
	nuevo, err := r.client.Horario.Create().
		SetEmpresaID(h.EmpresaID).
		SetUsuarioID(h.UsuarioID).
		SetHoraEntrada(h.HoraEntrada).
		SetHoraSalida(h.HoraSalida).
		SetToleranciaMinutos(h.ToleranciaMinutos).
		SetDiasLaborables(h.DiasLaborables).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("error al crear horario: %v", err)
	}

	return mapHorarioEntToDomain(nuevo), nil
}

func (r *HorarioRepoEnt) Actualizar(ctx context.Context, h *domain.Horario) (*domain.Horario, error) {
	actualizado, err := r.client.Horario.UpdateOneID(h.ID).
		SetHoraEntrada(h.HoraEntrada).
		SetHoraSalida(h.HoraSalida).
		SetToleranciaMinutos(h.ToleranciaMinutos).
		SetDiasLaborables(h.DiasLaborables).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("error al actualizar horario: %v", err)
	}

	return mapHorarioEntToDomain(actualizado), nil
}

func mapHorarioEntToDomain(h *ent.Horario) *domain.Horario {
	if h == nil {
		return nil
	}
	return &domain.Horario{
		ID:                 h.ID,
		EmpresaID:          h.EmpresaID,
		UsuarioID:          h.UsuarioID,
		HoraEntrada:        h.HoraEntrada,
		HoraSalida:         h.HoraSalida,
		ToleranciaMinutos:  h.ToleranciaMinutos,
		DiasLaborables:     h.DiasLaborables,
	}
}
