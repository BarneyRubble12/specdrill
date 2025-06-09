package di

import (
	"github.com/BarneyRubble12/specdrill/internal/core/executor"
	"github.com/BarneyRubble12/specdrill/internal/core/parser"
)

// Container holds all the application dependencies
type Container struct {
	Parser   *parser.Parser
	Executor *executor.Executor
}

// NewContainer creates a new application container
func NewContainer(
	parser *parser.Parser,
	executor *executor.Executor,
) *Container {
	return &Container{
		Parser:   parser,
		Executor: executor,
	}
}
