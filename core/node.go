package core

// Node is the basic unit of the file tree.
type Node struct {
	// Type is the type of the node, which can be "dir", "template", "copy", etc.
	Type string
	// Data is the data of the node, which should be handled by the handler.
	Data interface{}
}
