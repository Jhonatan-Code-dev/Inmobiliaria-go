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
	repository.NewGastoRepo,
	wire.Bind(new(domain.GastoRepository), new(*repository.GastoRepoEnt)),
	repository.NewClienteRepo,
	wire.Bind(new(domain.ClienteRepository), new(*repository.ClienteRepoEnt)),
	repository.NewInmuebleRepo,
	wire.Bind(new(domain.InmuebleRepository), new(*repository.InmuebleRepoEnt)),
	repository.NewAlquilerRepo,
	wire.Bind(new(domain.AlquilerRepository), new(*repository.AlquilerRepoEnt)),
	repository.NewPagoAlquilerRepo,
	wire.Bind(new(domain.PagoAlquilerRepository), new(*repository.PagoAlquilerRepoEnt)),
	repository.NewMovimientoCajaRepo,
	wire.Bind(new(domain.MovimientoCajaRepository), new(*repository.MovimientoCajaRepoEnt)),
	repository.NewTipoPagoRepo,
	wire.Bind(new(domain.TipoPagoRepository), new(*repository.TipoPagoRepoEnt)),
	repository.NewTipoIdentificacionRepo,
	wire.Bind(new(domain.TipoIdentificacionRepository), new(*repository.TipoIdentificacionRepoEnt)),
)

var ServiceSet = wire.NewSet(
	service.NewAdminService,
	service.NewMonedaService,
	service.NewUsuarioService,
	service.NewGastoService,
	service.NewClienteService,
	service.NewInmuebleService,
	service.NewAlquilerService,
	service.NewPagoAlquilerService,
	wire.Bind(new(domain.GastoService), new(*service.GastoService)),
	wire.Bind(new(domain.ClienteService), new(*service.ClienteService)),
	wire.Bind(new(domain.InmuebleService), new(*service.InmuebleService)),
	wire.Bind(new(domain.AlquilerService), new(*service.AlquilerService)),
	wire.Bind(new(domain.PagoAlquilerService), new(*service.PagoAlquilerService)),
	ProvideJWTSecret,
)

var ControllerSet = wire.NewSet(
	controller.NewAdminController,
	controller.NewMonedaController,
	controller.NewUsuarioController,
	controller.NewGastoController,
	controller.NewClienteController,
	controller.NewInmuebleController,
	controller.NewAlquilerController,
)
