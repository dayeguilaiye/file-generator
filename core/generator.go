package core

import "context"

// Generator is the main struct of the package.
type Generator struct {
	// handlers is the map of the handlers.
	handlers map[string]HandlerFunc
	Context  context.Context
}

// Handle registers a handler.
func (g *Generator) Handle(handler Handler) {
	g.handlers[handler.GetHandleType()] = handler.GetHandlerFunc()
}

// NewGenerator returns a new empty Generator.
func NewGenerator() *Generator {
	return &Generator{
		handlers: make(map[string]HandlerFunc),
	}
}

func (g *Generator) Generate(target string, node Node) error {
	return g.handlers[node.Type](g, target, node.Data)
}
