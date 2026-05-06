package di

import (
	"rentals-go/config/env"
	"rentals-go/ent"
	"rentals-go/internal/controller"
)

// App contiene solo las dependencias base del nuevo proyecto
type App struct {
	Config      *env.Config
	EntClient   *ent.Client
	AdminCtrl      *controller.AdminController
	MonedaCtrl     *controller.MonedaController
	UsuarioCtrl    *controller.UsuarioController
	GastoCtrl      *controller.GastoController
	ClienteCtrl    *controller.ClienteController
	InmuebleCtrl   *controller.InmuebleController
	AlquilerCtrl   *controller.AlquilerController
	StaffCtrl      *controller.StaffController
	CargoCtrl      *controller.CargoController
	ServicioCtrl   *controller.ServicioMedicionController
	TicketCtrl     *controller.TicketController
	DashboardCtrl  *controller.DashboardController
}
