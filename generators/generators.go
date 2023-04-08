package generators

import (
	"github.com/dayeguilaiye/file-generator/core"
	"github.com/dayeguilaiye/file-generator/handler"
)

// DefaultGenerator returns a new Generator with default handlers.
func DefaultGenerator() *core.Generator {
	g := core.NewGenerator()
	g.Handle(handler.NewCopyHandler())
	g.Handle(handler.NewDirHandler())
	g.Handle(handler.NewGoTemplateHandler())
	g.Handle(handler.NewReplaceHandler())
	g.Handle(handler.NewTgzHandler())
	return g
}
