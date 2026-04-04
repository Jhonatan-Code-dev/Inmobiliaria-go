package repository

import (
	"context"

	"rentals-go/ent"
	"rentals-go/ent/empresa"
	"rentals-go/internal/domain"
)

type EmpresaRepoEnt struct {
	client *ent.Client
}

func NewEmpresaRepo(client *ent.Client) *EmpresaRepoEnt {
	return &EmpresaRepoEnt{client: client}
}

func (r *EmpresaRepoEnt) Listar(ctx context.Context) ([]*domain.Empresa, error) {
	list, err := r.client.Empresa.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*domain.Empresa, 0, len(list))
	for _, e := range list {
		out = append(out, mapEmpresaEntity(e))
	}
	return out, nil
}

func (r *EmpresaRepoEnt) BuscarPorID(ctx context.Context, id int) (*domain.Empresa, error) {
	e, err := r.client.Empresa.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapEmpresaEntity(e), nil
}

func (r *EmpresaRepoEnt) Crear(ctx context.Context, emp *domain.Empresa) (*domain.Empresa, error) {
	e, err := r.client.Empresa.Create().
		SetNombre(emp.Nombre).
		SetNillableDocumentoFiscal(nilIfEmpty(emp.DocumentoFiscal)).
		SetNillableCorreo(nilIfEmpty(emp.Correo)).
		SetNillableTelefono(nilIfEmpty(emp.Telefono)).
		SetNillableDireccion(nilIfEmpty(emp.Direccion)).
		SetNillableCiudad(nilIfEmpty(emp.Ciudad)).
		SetNillablePais(nilIfEmpty(emp.Pais)).
		SetMoneda(emp.Moneda).
		SetMaximoUsuarios(defaultInt(emp.MaximoUsuarios, 1)).
		SetEstado(empresa.Estado(defaultStringValue(emp.Estado, "activa"))).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapEmpresaEntity(e), nil
}

func (r *EmpresaRepoEnt) Actualizar(ctx context.Context, emp *domain.Empresa) (*domain.Empresa, error) {
	e, err := r.client.Empresa.UpdateOneID(emp.ID).
		SetNombre(emp.Nombre).
		SetNillableDocumentoFiscal(nilIfEmpty(emp.DocumentoFiscal)).
		SetNillableCorreo(nilIfEmpty(emp.Correo)).
		SetNillableTelefono(nilIfEmpty(emp.Telefono)).
		SetNillableDireccion(nilIfEmpty(emp.Direccion)).
		SetNillableCiudad(nilIfEmpty(emp.Ciudad)).
		SetNillablePais(nilIfEmpty(emp.Pais)).
		SetMoneda(emp.Moneda).
		SetMaximoUsuarios(defaultInt(emp.MaximoUsuarios, 1)).
		SetEstado(empresa.Estado(defaultStringValue(emp.Estado, "activa"))).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapEmpresaEntity(e), nil
}

func (r *EmpresaRepoEnt) Eliminar(ctx context.Context, id int) error {
	return r.client.Empresa.DeleteOneID(id).Exec(ctx)
}

func mapEmpresaEntity(e *ent.Empresa) *domain.Empresa {
	return &domain.Empresa{
		ID:              e.ID,
		Nombre:          e.Nombre,
		DocumentoFiscal: ptrToString(e.DocumentoFiscal),
		Correo:          ptrToString(e.Correo),
		Telefono:        ptrToString(e.Telefono),
		Direccion:       ptrToString(e.Direccion),
		Ciudad:          ptrToString(e.Ciudad),
		Pais:            ptrToString(e.Pais),
		Moneda:          e.Moneda,
		MaximoUsuarios:  e.MaximoUsuarios,
		Estado:          string(e.Estado),
	}
}

func defaultInt(val, def int) int {
	if val <= 0 {
		return def
	}
	return val
}

func defaultStringValue(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
