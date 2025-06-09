//go:build wireinject
// +build wireinject

package di

import (
	"github.com/BarneyRubble12/specdrill/internal/core/executor"
	"github.com/BarneyRubble12/specdrill/internal/core/parser"
	"github.com/google/wire"
)

// ProviderSet is a Wire provider set for the application
var ProviderSet = wire.NewSet(
	parser.NewParser,
	executor.NewExecutor,
)

// InitializeContainer creates a new application container with all dependencies
func InitializeContainer() (*Container, error) {
	wire.Build(
		ProviderSet,
		NewContainer,
	)
	return nil, nil
}
