package di

import (
	"github.com/BarneyRubble12/specdrill/internal/core/executor"
	"github.com/BarneyRubble12/specdrill/internal/core/parser"
	"github.com/google/wire"
)

// ProviderSet is a Wire provider set for the application
var ProviderSet = wire.NewSet(
	// Core components
	parser.NewOpenAPIParser,
	executor.NewHTTPExecutor,
	wire.Bind(new(parser.Parser), new(*parser.OpenAPIParser)),
	wire.Bind(new(executor.Executor), new(*executor.HTTPExecutor)),
)
