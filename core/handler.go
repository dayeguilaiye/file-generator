package core

import "errors"

type HandlerFunc func(generator *Generator, targetDir string, data interface{}) error

type Handler interface {
	GetHandleType() string
	GetHandlerFunc() HandlerFunc
}

var WrongDataTypeError = errors.New("can not transform data")
