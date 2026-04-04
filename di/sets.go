package di

import (
	"rentals-go/internal/controller"
	"rentals-go/internal/domain"
	"rentals-go/internal/repository"
	"rentals-go/internal/service"

	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	repository.NewAdminRepo,
	wire.Bind(new(domain.AdminRepository), new(*repository.AdminRepoEnt)),
	repository.NewEmpresaRepo,
	wire.Bind(new(domain.EmpresaRepository), new(*repository.EmpresaRepoEnt)),
	repository.NewUsuarioRepo,
	wire.Bind(new(domain.UsuarioRepository), new(*repository.UsuarioRepoEnt)),
	repository.NewRolRepo,
	wire.Bind(new(domain.RolRepository), new(*repository.RolRepoEnt)),
	repository.NewMembresiaRepo,
	wire.Bind(new(domain.MembresiaRepository), new(*repository.MembresiaRepoEnt)),
)

var ServiceSet = wire.NewSet(
	service.NewAdminService,
	service.NewMonedaService,
	service.NewUsuarioService,
	ProvideJWTSecret,
)

var ControllerSet = wire.NewSet(
	controller.NewAdminController,
	controller.NewMonedaController,
	controller.NewUsuarioController,
)
