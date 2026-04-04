//go:build wireinject
// +build wireinject

package di

import "github.com/google/wire"

func InitializeApp() (*App, error) {
	wire.Build(
		ConfigSet,
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
